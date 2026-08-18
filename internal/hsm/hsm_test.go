package hsm

import (
	"context"
	"strings"
	"testing"
)

const testMasterKEK = "0102030405060708090a0b0c0d0e0f10"

var custodians = []string{"alice", "bob", "carol"}

func newTestHSM(t *testing.T) *HSM {
	t.Helper()
	h, err := NewHSM(mustBytes(testMasterKEK), 2, custodians)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func mustBytes(hex string) []byte {
	b, err := hexBytes(hex)
	if err != nil {
		panic(err)
	}
	return b
}

func TestKeyWrapRFC3394KAT(t *testing.T) {
	if err := keyWrapKnownAnswer(); err != nil {
		t.Fatal(err)
	}
}

func TestKeyWrapWrongKEKFails(t *testing.T) {
	key := mustBytes(rfc3394Key)
	kek := mustBytes(rfc3394KEK)
	wrapped, err := wrapAESKey(kek, key)
	if err != nil {
		t.Fatal(err)
	}
	wrong := mustBytes("00000000000000000000000000000000")
	if _, err := unwrapAESKey(wrong, wrapped); err == nil {
		t.Fatal("expected unwrap to fail with the wrong KEK")
	}
}

func TestRetailMACKAT(t *testing.T) {
	// Independent vectors (computed with .NET DES): K1=0102030405060708,
	// K2=0807060504030201. Single-block and chained-CBC cases.
	key := mustBytes("01020304050607080807060504030201")
	mac, err := RetailMAC(key, mustBytes("0001020304050607"))
	if err != nil {
		t.Fatal(err)
	}
	if !equal(mac, mustBytes("563e741d0030117a")) {
		t.Fatalf("single-block MAC = %x", mac)
	}
	mac, err = RetailMAC(key, mustBytes("000102030405060708090a0b0c0d0e0f"))
	if err != nil {
		t.Fatal(err)
	}
	if !equal(mac, mustBytes("c3bb2da66e7188cd")) {
		t.Fatalf("two-block MAC = %x", mac)
	}
}

func TestRetailMACTamperDetected(t *testing.T) {
	key := mustBytes("01020304050607080807060504030201")
	data := mustBytes("48656c6c6f20436c61726121") // "Hello Clara!"
	mac, err := RetailMAC(key, data)
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyMAC(key, data, mac) {
		t.Fatal("valid MAC rejected")
	}
	tampered := append([]byte(nil), data...)
	tampered[len(tampered)-1] ^= 0x01
	if VerifyMAC(key, tampered, mac) {
		t.Fatal("tampered message accepted")
	}
	if VerifyMAC(key, data, append(append([]byte(nil), mac...), 0x00)) {
		t.Fatal("tampered MAC accepted")
	}
}

func TestPinBlockISO0(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	pinKey, err := h.CreateKey(ctx, KeyTypePIN, "", "atm link", []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	pan := "4000001234567890"
	block, err := h.ComputePINBlock(ctx, pinKey.ID, pan, "1234", PinFormat0)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := h.VerifyPIN(ctx, pinKey.ID, block, pan, "1234", PinFormat0)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("correct PIN rejected")
	}
	ok, err = h.VerifyPIN(ctx, pinKey.ID, block, pan, "9999", PinFormat0)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("wrong PIN accepted")
	}
	// A different PAN must not decode to the same PIN (PAN-bound block).
	ok, err = h.VerifyPIN(ctx, pinKey.ID, block, "4111111111111111", "1234", PinFormat0)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("ISO-0 block verified under a different PAN")
	}
}

func TestPinBlockFormat4(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	pinKey, err := h.CreateKey(ctx, KeyTypePIN, AlgAES, "format4 link", []string{"alice", "carol"})
	if err != nil {
		t.Fatal(err)
	}
	pan := "4000001234567890"
	block, err := h.ComputePINBlock(ctx, pinKey.ID, pan, "123456", PinFormat4)
	if err != nil {
		t.Fatal(err)
	}
	if len(block) != 16 {
		t.Fatalf("format 4 block = %d bytes, want 16", len(block))
	}
	ok, err := h.VerifyPIN(ctx, pinKey.ID, block, pan, "123456", PinFormat4)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("correct PIN rejected in format 4")
	}
	ok, err = h.VerifyPIN(ctx, pinKey.ID, block, pan, "654321", PinFormat4)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("wrong PIN accepted in format 4")
	}
}

func TestDualControlRequired(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	if _, err := h.CreateKey(ctx, KeyTypeMAC, "", "single approver", []string{"alice"}); err == nil {
		t.Fatal("expected ceremony to fail with a single approver")
	}
	if _, err := h.CreateKey(ctx, KeyTypeMAC, "", "unauthorized", []string{"alice", "mallory"}); err == nil {
		t.Fatal("expected ceremony to fail with an unknown custodian")
	}
	if _, err := h.CreateKey(ctx, KeyTypeMAC, "", "duplicates", []string{"alice", "alice"}); err == nil {
		t.Fatal("expected ceremony to fail with duplicate approvals")
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	macKey, err := h.CreateKey(ctx, KeyTypeMAC, "", "issuer mac", []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	// A shared KEK is loaded into both HSMs during a ceremony (clear load),
	// simulating out-of-band distribution between the network and the member.
	const sharedKEK = "808182838485868788898a8b8c8d8e8f"
	hKEK, err := h.LoadKeyMaterial(ctx, KeyTypeKEK, AlgAES, "acquirer-kek", sharedKEK, []string{"bob", "carol"})
	if err != nil {
		t.Fatal(err)
	}
	peer := newTestHSM(t)
	peerKEK, err := peer.LoadKeyMaterial(ctx, KeyTypeKEK, AlgAES, "acquirer-kek", sharedKEK, []string{"bob", "carol"})
	if err != nil {
		t.Fatal(err)
	}

	block, err := h.ExportKeyBlock(ctx, macKey.ID, hKEK.ID, []string{"alice", "carol"})
	if err != nil {
		t.Fatal(err)
	}
	imported, err := peer.ImportKeyBlock(ctx, block, peerKEK.ID, []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	if imported.Type != KeyTypeMAC || imported.KVN != macKey.KVN {
		t.Fatalf("imported = %+v, want MAC kvn %d", imported, imported.KVN)
	}
	// Both HSMs hold the same key material: MACs agree on both sides.
	msg := []byte("interchange message for the issuer link")
	m1, err := h.ComputeMAC(ctx, macKey.ID, msg)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := peer.ComputeMAC(ctx, imported.ID, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(m1, m2) {
		t.Fatalf("MACs differ after import: %x vs %x", m1, m2)
	}
}

func TestKeyRotationRetiresOldVersion(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	macKey, err := h.CreateKey(ctx, KeyTypeMAC, "", "issuer mac", []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	rotated, err := h.RotateKey(ctx, macKey.ID, []string{"alice", "carol"})
	if err != nil {
		t.Fatal(err)
	}
	if rotated.KVN != macKey.KVN+1 {
		t.Fatalf("kvn = %d, want %d", rotated.KVN, macKey.KVN+1)
	}
	// Old key must now be retired.
	oldKey, err := h.KeyInfo(ctx, macKey.ID)
	if err != nil {
		t.Fatal(err)
	}
	if oldKey.Status != KeyRetired {
		t.Fatalf("old key status = %s, want retired", oldKey.Status)
	}
	// Old and new keys must produce different MACs for the same data.
	data := []byte("rotate me")
	m1, err := h.ComputeMAC(ctx, macKey.ID, data)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := h.ComputeMAC(ctx, rotated.ID, data)
	if err != nil {
		t.Fatal(err)
	}
	if equal(m1, m2) {
		t.Fatal("rotated key must not produce the old MAC")
	}
}

func TestMACComputationAndAudit(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	macKey, err := h.CreateKey(ctx, KeyTypeMAC, "", "acquirer mac", []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("0210<bitmap>..MESSAGE PAYLOAD..")
	mac, err := h.ComputeMAC(ctx, macKey.ID, msg)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := h.VerifyMAC(ctx, macKey.ID, msg, mac)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("valid MAC rejected")
	}
	ok, err = h.VerifyMAC(ctx, macKey.ID, []byte("tampered"), mac)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("tampered message accepted")
	}
	// The audit trail captured the key ceremony and the MAC operations.
	audit, err := h.AuditLog(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ops := []string{}
	for _, e := range audit {
		ops = append(ops, e.Op)
	}
	for _, want := range []string{"key.generate", "mac.compute", "mac.verify"} {
		if !strings.Contains(strings.Join(ops, " "), want) {
			t.Fatalf("audit trail missing %s: %+v", want, ops)
		}
	}
}

func TestZeroize(t *testing.T) {
	h := newTestHSM(t)
	ctx := context.Background()
	k, err := h.CreateKey(ctx, KeyTypePIN, "", "to be destroyed", []string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Zeroize(ctx, []string{"alice", "bob"}); err != nil {
		t.Fatal(err)
	}
	if _, err := h.ComputePINBlock(ctx, k.ID, "4000001234567890", "1234", PinFormat0); err == nil {
		t.Fatal("expected operation on zeroized HSM to fail")
	}
	if _, err := h.CreateKey(ctx, KeyTypeKEK, "", "after", []string{"alice", "bob"}); err == nil {
		t.Fatal("expected key creation on zeroized HSM to fail")
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
