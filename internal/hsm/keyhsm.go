package hsm

import (
	"context"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// KeyType classifies keys by function; the scheme separates PIN, MAC, and
// transport keys (docs/17 §17.5).
type KeyType string

// Key types.
const (
	KeyTypeZMK  KeyType = "ZMK"  // zone master key (between network and members)
	KeyTypeKEK  KeyType = "KEK"  // key-encrypting key
	KeyTypePIN  KeyType = "PIN"  // PIN encryption key
	KeyTypeMAC  KeyType = "MAC"  // message authentication key
	KeyTypeData KeyType = "DATA" // data-encryption key
)

// Key algorithms.
const (
	AlgAES  = "A" // AES-128 (16 bytes)
	Alg3DES = "T" // 2-key 3DES (16 bytes)
)

// Key statuses.
const (
	KeyActive  = "active"
	KeyRetired = "retired"
)

// Key is a handle to wrapped key material inside the HSM. Clear key material
// never leaves the HSM; only the AES-key-wrapped copy under the local master
// KEK is stored (docs/17 §17.6).
type Key struct {
	ID        string
	Type      KeyType
	Alg       string
	KVN       int
	Status    string
	Label     string
	Wrapped   []byte
	CreatedAt time.Time
}

// AuditEvent records every key and crypto event (docs/17 §17.6).
type AuditEvent struct {
	Time   time.Time
	Op     string
	KeyID  string
	Actor  string
	Detail string
}

// HSM simulates a Secure Cryptographic Device: keys are stored wrapped under
// a master KEK, sensitive operations (PIN translate/verify, MAC) run here, and
// key ceremonies require dual control (M-of-N) approval.
type HSM struct {
	mu         sync.Mutex
	masterKEK  []byte
	keys       map[string]*Key
	custodians map[string]bool
	threshold  int
	audit      []AuditEvent
	zeroized   bool
	log        *slog.Logger
}

// NewHSM boots the HSM with the local master KEK and the custodian set.
func NewHSM(masterKEK []byte, threshold int, custodians []string) (*HSM, error) {
	if len(masterKEK) != 16 {
		return nil, errors.New("hsm: master KEK must be 16 bytes")
	}
	if threshold < 1 || threshold > len(custodians) {
		return nil, fmt.Errorf("hsm: threshold %d invalid for %d custodians", threshold, len(custodians))
	}
	cs := map[string]bool{}
	for _, c := range custodians {
		cs[c] = true
	}
	return &HSM{
		masterKEK:  append([]byte(nil), masterKEK...),
		keys:       map[string]*Key{},
		custodians: cs,
		threshold:  threshold,
		log:        slog.Default(),
	}, nil
}

// CreateKey runs a dual-control ceremony to generate a new key of the given
// type and wraps it under the master KEK (docs/17 §17.3).
func (h *HSM) CreateKey(ctx context.Context, typ KeyType, alg, label string, approvers []string) (Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return Key{}, err
	}
	if alg == "" {
		alg = defaultAlg(typ)
	}
	size := 16
	if alg == Alg3DES {
		size = 16
	}
	material := make([]byte, size)
	if _, err := rand.Read(material); err != nil {
		return Key{}, err
	}
	wrapped, err := wrapAESKey(h.masterKEK, material)
	if err != nil {
		return Key{}, err
	}
	k := &Key{
		ID:        fmt.Sprintf("K-%s-%d", typ, time.Now().UTC().UnixNano()),
		Type:      typ,
		Alg:       alg,
		KVN:       1,
		Status:    KeyActive,
		Label:     label,
		Wrapped:   wrapped,
		CreatedAt: time.Now().UTC(),
	}
	h.keys[k.ID] = k
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "key.generate", KeyID: k.ID, Actor: strings.Join(approvers, "+"), Detail: fmt.Sprintf("%s %s kvn=%d", typ, alg, k.KVN)})
	h.log.Info("key generated", "id", k.ID, "type", typ, "label", label)
	return *k, nil
}

// RotateKey runs a dual-control ceremony to roll a key: the old version is
// retired (still usable for verify) and a new active version is created
// (docs/17 §17.5).
func (h *HSM) RotateKey(ctx context.Context, id string, approvers []string) (Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return Key{}, err
	}
	old, ok := h.keys[id]
	if !ok {
		return Key{}, fmt.Errorf("hsm: unknown key %s", id)
	}
	material := make([]byte, 16)
	if _, err := rand.Read(material); err != nil {
		return Key{}, err
	}
	wrapped, err := wrapAESKey(h.masterKEK, material)
	if err != nil {
		return Key{}, err
	}
	old.Status = KeyRetired
	k := &Key{
		ID:        old.ID + "-" + fmt.Sprint(old.KVN+1),
		Type:      old.Type,
		Alg:       old.Alg,
		KVN:       old.KVN + 1,
		Status:    KeyActive,
		Label:     old.Label,
		Wrapped:   wrapped,
		CreatedAt: time.Now().UTC(),
	}
	h.keys[k.ID] = k
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "key.rotate", KeyID: old.ID, Actor: strings.Join(approvers, "+"), Detail: fmt.Sprintf("kvn %d -> %d", old.KVN, k.KVN)})
	h.log.Info("key rotated", "old", old.ID, "new", k.ID, "kvn", k.KVN)
	return *k, nil
}

// LoadKeyMaterial loads clear key material into the HSM during a dual-control
// ceremony (for example a pre-shared KEK), wrapping it under the master KEK.
func (h *HSM) LoadKeyMaterial(ctx context.Context, typ KeyType, alg, label, materialHex string, approvers []string) (Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return Key{}, err
	}
	material, err := hexBytes(materialHex)
	if err != nil {
		return Key{}, err
	}
	if len(material) != 16 {
		return Key{}, errors.New("hsm: clear key material must be 16 bytes")
	}
	wrapped, err := wrapAESKey(h.masterKEK, material)
	if err != nil {
		return Key{}, err
	}
	k := &Key{
		ID:        fmt.Sprintf("K-%s-%d", typ, time.Now().UTC().UnixNano()),
		Type:      typ,
		Alg:       alg,
		KVN:       1,
		Status:    KeyActive,
		Label:     label,
		Wrapped:   wrapped,
		CreatedAt: time.Now().UTC(),
	}
	h.keys[k.ID] = k
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "key.load-clear", KeyID: k.ID, Actor: strings.Join(approvers, "+"), Detail: fmt.Sprintf("%s %s", typ, label)})
	h.log.Info("clear key loaded", "id", k.ID, "type", typ, "label", label)
	return *k, nil
}

// KeyInfo returns a key handle by ID.
func (h *HSM) KeyInfo(ctx context.Context, id string) (Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	k, ok := h.keys[id]
	if !ok {
		return Key{}, fmt.Errorf("hsm: unknown key %s", id)
	}
	cp := *k
	cp.Wrapped = append([]byte(nil), k.Wrapped...)
	return cp, nil
}

// ExportKeyBlock returns a TR-31-style wrapped key block for transport to a
// member under the given KEK (docs/17 §17.3). Clear keys never leave the HSM.
func (h *HSM) ExportKeyBlock(ctx context.Context, id, kekID string, approvers []string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return "", err
	}
	k, ok := h.keys[id]
	if !ok {
		return "", fmt.Errorf("hsm: unknown key %s", id)
	}
	kek, ok := h.keys[kekID]
	if !ok || kek.Type != KeyTypeKEK {
		return "", fmt.Errorf("hsm: unknown KEK %s", kekID)
	}
	material, err := unwrapAESKey(h.masterKEK, k.Wrapped)
	if err != nil {
		return "", err
	}
	kekMaterial, err := unwrapAESKey(h.masterKEK, kek.Wrapped)
	if err != nil {
		return "", err
	}
	block, err := wrapAESKey(kekMaterial, material)
	if err != nil {
		return "", err
	}
	header := "C0" + usageCode(k.Type) + k.Alg + fmt.Sprintf("%02d", k.KVN)
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "key.export", KeyID: id, Actor: strings.Join(approvers, "+"), Detail: "wrapped under " + kekID})
	return header + ":" + base64.StdEncoding.EncodeToString(block), nil
}

// ImportKeyBlock loads a TR-31-style key block received from a peer, unwraps
// it with the local KEK, and re-wraps it under the master KEK.
func (h *HSM) ImportKeyBlock(ctx context.Context, block, kekID string, approvers []string) (Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return Key{}, err
	}
	kek, ok := h.keys[kekID]
	if !ok || kek.Type != KeyTypeKEK {
		return Key{}, fmt.Errorf("hsm: unknown KEK %s", kekID)
	}
	parts := strings.SplitN(block, ":", 2)
	if len(parts) != 2 || len(parts[0]) != 6 || parts[0][0:2] != "C0" {
		return Key{}, errors.New("hsm: malformed key block header")
	}
	header, b64 := parts[0], parts[1]
	kekMaterial, err := unwrapAESKey(h.masterKEK, kek.Wrapped)
	if err != nil {
		return Key{}, err
	}
	payload, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Key{}, err
	}
	material, err := unwrapAESKey(kekMaterial, payload)
	if err != nil {
		return Key{}, err
	}
	wrapped, err := wrapAESKey(h.masterKEK, material)
	if err != nil {
		return Key{}, err
	}
	var kvn int
	if _, err := fmt.Sscanf(header[4:6], "%d", &kvn); err != nil {
		return Key{}, err
	}
	k := &Key{
		ID:        fmt.Sprintf("K-%s-%d", typeFromUsage(header[2:3]), time.Now().UTC().UnixNano()),
		Type:      typeFromUsage(header[2:3]),
		Alg:       header[3:4],
		KVN:       kvn,
		Status:    KeyActive,
		Label:     "imported",
		Wrapped:   wrapped,
		CreatedAt: time.Now().UTC(),
	}
	h.keys[k.ID] = k
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "key.import", KeyID: k.ID, Actor: strings.Join(approvers, "+"), Detail: "from " + header})
	h.log.Info("key imported", "id", k.ID, "type", k.Type, "kvn", k.KVN)
	return *k, nil
}

// VerifyPIN checks an ISO 9564 PIN block inside the HSM in constant time; the
// PIN is never returned to the caller (docs/17 §17.2).
func (h *HSM) VerifyPIN(ctx context.Context, pinKeyID string, block []byte, pan, expectedPIN string, format int) (bool, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.usable(); err != nil {
		return false, err
	}
	k, ok := h.keys[pinKeyID]
	if !ok || k.Type != KeyTypePIN {
		return false, fmt.Errorf("hsm: unknown PIN key %s", pinKeyID)
	}
	material, err := unwrapAESKey(h.masterKEK, k.Wrapped)
	if err != nil {
		return false, err
	}
	var got string
	switch format {
	case PinFormat0:
		if len(block) != 8 {
			return false, errors.New("hsm: ISO-0 block must be 8 bytes")
		}
		clear, err := tdesECB(material, block, false)
		if err != nil {
			return false, err
		}
		got, err = pinFromFormat0(clear, pan)
		if err != nil {
			return false, err
		}
	case PinFormat4:
		got, err = pinFromFormat4(hexString(material), block)
		if err != nil {
			return false, err
		}
	default:
		return false, fmt.Errorf("hsm: unsupported PIN format %d", format)
	}
	okVerify := verifyPIN(got, expectedPIN)
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "pin.verify", KeyID: pinKeyID, Actor: "hsm", Detail: fmt.Sprintf("format %d result=%t", format, okVerify)})
	return okVerify, nil
}

// ComputePINBlock builds an ISO 9564 PIN block inside the HSM (used to
// translate an incoming block to an outgoing format).
func (h *HSM) ComputePINBlock(ctx context.Context, pinKeyID, pan, pin string, format int) ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.usable(); err != nil {
		return nil, err
	}
	k, ok := h.keys[pinKeyID]
	if !ok || k.Type != KeyTypePIN {
		return nil, fmt.Errorf("hsm: unknown PIN key %s", pinKeyID)
	}
	material, err := unwrapAESKey(h.masterKEK, k.Wrapped)
	if err != nil {
		return nil, err
	}
	var block []byte
	switch format {
	case PinFormat0:
		clear, err := format0Block(pin, pan)
		if err != nil {
			return nil, err
		}
		block, err = tdesECB(material, clear, true)
		if err != nil {
			return nil, err
		}
	case PinFormat4:
		block, err = format4Block(hexString(material), pin, pan)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("hsm: unsupported PIN format %d", format)
	}
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "pin.translate", KeyID: pinKeyID, Actor: "hsm", Detail: fmt.Sprintf("format %d", format)})
	return block, nil
}

// ComputeMAC computes a retail MAC over message data inside the HSM.
func (h *HSM) ComputeMAC(ctx context.Context, macKeyID string, data []byte) ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.usable(); err != nil {
		return nil, err
	}
	k, ok := h.keys[macKeyID]
	if !ok || k.Type != KeyTypeMAC {
		return nil, fmt.Errorf("hsm: unknown MAC key %s", macKeyID)
	}
	material, err := unwrapAESKey(h.masterKEK, k.Wrapped)
	if err != nil {
		return nil, err
	}
	mac, err := RetailMAC(material, data)
	if err != nil {
		return nil, err
	}
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "mac.compute", KeyID: macKeyID, Actor: "hsm", Detail: fmt.Sprintf("%d bytes", len(data))})
	return mac, nil
}

// VerifyMAC verifies a retail MAC in constant time.
func (h *HSM) VerifyMAC(ctx context.Context, macKeyID string, data, mac []byte) (bool, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.usable(); err != nil {
		return false, err
	}
	k, ok := h.keys[macKeyID]
	if !ok || k.Type != KeyTypeMAC {
		return false, fmt.Errorf("hsm: unknown MAC key %s", macKeyID)
	}
	material, err := unwrapAESKey(h.masterKEK, k.Wrapped)
	if err != nil {
		return false, err
	}
	okVerify := VerifyMAC(material, data, mac)
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "mac.verify", KeyID: macKeyID, Actor: "hsm", Detail: fmt.Sprintf("result=%t", okVerify)})
	return okVerify, nil
}

// AuditLog returns the ordered audit trail (docs/17 §17.6).
func (h *HSM) AuditLog(ctx context.Context) ([]AuditEvent, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]AuditEvent, len(h.audit))
	copy(out, h.audit)
	return out, nil
}

// ListKeys lists the key handles (wrapped material only).
func (h *HSM) ListKeys(ctx context.Context) ([]Key, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	var out []Key
	for _, k := range h.keys {
		cp := *k
		cp.Wrapped = append([]byte(nil), k.Wrapped...)
		out = append(out, cp)
	}
	return out, nil
}

// Zeroize destroys the master KEK and every wrapped key after a dual-control
// ceremony; the HSM becomes unusable (docs/17 §17.3).
func (h *HSM) Zeroize(ctx context.Context, approvers []string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.approve(approvers); err != nil {
		return err
	}
	h.auditLocked(AuditEvent{Time: time.Now().UTC(), Op: "hsm.zeroize", Actor: strings.Join(approvers, "+")})
	h.masterKEK = nil
	h.keys = map[string]*Key{}
	h.zeroized = true
	h.log.Warn("HSM zeroized")
	return nil
}

func (h *HSM) approve(approvers []string) error {
	if h.zeroized {
		return errors.New("hsm: zeroized")
	}
	if len(approvers) < h.threshold {
		return fmt.Errorf("hsm: ceremony needs %d distinct custodians, got %d", h.threshold, len(approvers))
	}
	seen := map[string]bool{}
	for _, a := range approvers {
		if !h.custodians[a] {
			return fmt.Errorf("hsm: %s is not a registered custodian", a)
		}
		seen[a] = true
	}
	if len(seen) < h.threshold {
		return fmt.Errorf("hsm: ceremony needs %d distinct custodians, got %d", h.threshold, len(seen))
	}
	return nil
}

func (h *HSM) usable() error {
	if h.zeroized {
		return errors.New("hsm: zeroized")
	}
	return nil
}

func (h *HSM) auditLocked(e AuditEvent) {
	h.audit = append(h.audit, e)
}

func defaultAlg(typ KeyType) string {
	switch typ {
	case KeyTypePIN:
		return Alg3DES
	case KeyTypeMAC:
		return Alg3DES
	default:
		return AlgAES
	}
}

func usageCode(typ KeyType) string {
	switch typ {
	case KeyTypeZMK:
		return "Z"
	case KeyTypeKEK:
		return "K"
	case KeyTypePIN:
		return "P"
	case KeyTypeMAC:
		return "M"
	default:
		return "D"
	}
}

func typeFromUsage(c string) KeyType {
	switch c {
	case "Z":
		return KeyTypeZMK
	case "K":
		return KeyTypeKEK
	case "P":
		return KeyTypePIN
	case "M":
		return KeyTypeMAC
	default:
		return KeyTypeData
	}
}

func hexString(b []byte) string {
	const hexDigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexDigits[v>>4]
		out[i*2+1] = hexDigits[v&0x0F]
	}
	return string(out)
}

// tdesECB encrypts (or decrypts) a single 8-byte block with 2-key 3DES
// (K3 = K1); 16-byte keys are expanded to 24 bytes.
func tdesECB(key, block []byte, encrypt bool) ([]byte, error) {
	key24 := append(append([]byte(nil), key...), key[:8]...)
	cipher, err := des.NewTripleDESCipher(key24)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 8)
	if encrypt {
		cipher.Encrypt(out, block)
	} else {
		cipher.Decrypt(out, block)
	}
	return out, nil
}
