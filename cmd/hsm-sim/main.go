// Command hsm-sim exercises the Clara Network HSM simulation: key ceremonies
// with dual control, TR-31-style key blocks (AES key wrap), ISO 9564 PIN
// blocks and PIN verification, ISO 9797-1 retail MACs with tamper detection,
// key rotation, the audit trail, and a full zeroize.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/hsm"
)

const masterKEKHex = "0102030405060708090a0b0c0d0e0f10"

var custodians = []string{"alice", "bob", "carol"}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	h, err := hsm.NewHSM(mustBytes(masterKEKHex), 2, custodians)
	if err != nil {
		logger.Error("boot failed", "err", err)
		os.Exit(1)
	}
	logger.Info("HSM booted", "threshold", 2, "custodians", len(custodians))

	// Key ceremonies: every sensitive operation needs M-of-N approval.
	zmk, err := h.CreateKey(ctx, hsm.KeyTypeZMK, "", "clara net <-> BankA", []string{"alice", "bob"})
	if err != nil {
		logger.Error("ceremony failed", "err", err)
		os.Exit(1)
	}
	pinKey, err := h.CreateKey(ctx, hsm.KeyTypePIN, "", "acquirer PIN key", []string{"alice", "bob"})
	if err != nil {
		logger.Error("ceremony failed", "err", err)
		os.Exit(1)
	}
	macKey, err := h.CreateKey(ctx, hsm.KeyTypeMAC, "", "issuer link MAC", []string{"alice", "bob"})
	if err != nil {
		logger.Error("ceremony failed", "err", err)
		os.Exit(1)
	}
	logger.Info("keys generated",
		"zmk", zmk.ID, "pin", pinKey.ID, "mac", macKey.ID)

	// PIN verification (ISO-0): the PIN never leaves the HSM.
	pan := "4000001234567890"
	block, err := h.ComputePINBlock(ctx, pinKey.ID, pan, "4321", hsm.PinFormat0)
	if err != nil {
		logger.Error("pin block failed", "err", err)
		os.Exit(1)
	}
	ok, err := h.VerifyPIN(ctx, pinKey.ID, block, pan, "4321", hsm.PinFormat0)
	if err != nil {
		logger.Error("pin verify failed", "err", err)
		os.Exit(1)
	}
	logger.Info("ISO-0 PIN verify", "block", hexString(block), "accepted", ok)

	ok, err = h.VerifyPIN(ctx, pinKey.ID, block, pan, "0000", hsm.PinFormat0)
	if err != nil {
		logger.Error("pin verify failed", "err", err)
		os.Exit(1)
	}
	logger.Warn("wrong PIN", "accepted", ok)

	// PIN translation ISO-0 -> ISO 9564-1:2017 format 4 (AES).
	block4, err := h.ComputePINBlock(ctx, pinKey.ID, pan, "4321", hsm.PinFormat4)
	if err != nil {
		logger.Error("format-4 block failed", "err", err)
		os.Exit(1)
	}
	ok, err = h.VerifyPIN(ctx, pinKey.ID, block4, pan, "4321", hsm.PinFormat4)
	if err != nil {
		logger.Error("format-4 verify failed", "err", err)
		os.Exit(1)
	}
	logger.Info("ISO-4 PIN verify", "block", hexString(block4), "accepted", ok)

	// Retail MAC (ISO 9797-1 alg 3) with tamper detection.
	msg := []byte("0100B2308221400000000004000000004123456780120012345678")
	mac, err := h.ComputeMAC(ctx, macKey.ID, msg)
	if err != nil {
		logger.Error("mac failed", "err", err)
		os.Exit(1)
	}
	good, err := h.VerifyMAC(ctx, macKey.ID, msg, mac)
	if err != nil {
		logger.Error("mac verify failed", "err", err)
		os.Exit(1)
	}
	tampered := append([]byte(nil), msg...)
	tampered[10] ^= 0x01
	bad, err := h.VerifyMAC(ctx, macKey.ID, tampered, mac)
	if err != nil {
		logger.Error("mac verify failed", "err", err)
		os.Exit(1)
	}
	logger.Info("retail MAC", "mac", hexString(mac), "intact", good, "tampered", bad)

	// Key distribution: export the MAC key to a member under a shared KEK.
	// The KEK is loaded into both HSMs in a ceremony, simulating out-of-band
	// exchange.
	const sharedKEK = "808182838485868788898a8b8c8d8e8f"
	kek, err := h.LoadKeyMaterial(ctx, hsm.KeyTypeKEK, hsm.AlgAES, "BankA KEK", sharedKEK, []string{"bob", "carol"})
	if err != nil {
		logger.Error("kek load failed", "err", err)
		os.Exit(1)
	}
	peer, err := hsm.NewHSM(mustBytes(masterKEKHex), 2, custodians)
	if err != nil {
		logger.Error("peer boot failed", "err", err)
		os.Exit(1)
	}
	peerKEK, err := peer.LoadKeyMaterial(ctx, hsm.KeyTypeKEK, hsm.AlgAES, "BankA KEK", sharedKEK, []string{"bob", "carol"})
	if err != nil {
		logger.Error("peer kek load failed", "err", err)
		os.Exit(1)
	}
	blockStr, err := h.ExportKeyBlock(ctx, macKey.ID, kek.ID, []string{"alice", "carol"})
	if err != nil {
		logger.Error("export failed", "err", err)
		os.Exit(1)
	}
	imported, err := peer.ImportKeyBlock(ctx, blockStr, peerKEK.ID, []string{"alice", "bob"})
	if err != nil {
		logger.Error("import failed", "err", err)
		os.Exit(1)
	}
	peerMAC, err := peer.ComputeMAC(ctx, imported.ID, msg)
	if err != nil {
		logger.Error("peer mac failed", "err", err)
		os.Exit(1)
	}
	logger.Info("key distribution", "block", blockStr[:16]+"...", "imported", imported.ID,
		"kvn", imported.KVN, "peer_mac_matches", equal(mac, peerMAC))

	// Rotation: the old MAC key retires, the new one produces different MACs.
	rotated, err := h.RotateKey(ctx, macKey.ID, []string{"alice", "carol"})
	if err != nil {
		logger.Error("rotation failed", "err", err)
		os.Exit(1)
	}
	rotMAC, err := h.ComputeMAC(ctx, rotated.ID, msg)
	if err != nil {
		logger.Error("rotated mac failed", "err", err)
		os.Exit(1)
	}
	logger.Info("key rotated", "new_kvn", rotated.KVN, "mac_changed", !equal(mac, rotMAC))

	// Audit trail.
	events, err := h.AuditLog(ctx)
	if err != nil {
		logger.Error("audit failed", "err", err)
		os.Exit(1)
	}
	for _, e := range events {
		logger.Info("audit", "op", e.Op, "key", e.KeyID, "actor", e.Actor, "detail", e.Detail)
	}

	// Zeroize: dual-control destruction of all material.
	if err := h.Zeroize(ctx, []string{"alice", "bob"}); err != nil {
		logger.Error("zeroize failed", "err", err)
		os.Exit(1)
	}
	if _, err := h.ComputeMAC(ctx, rotated.ID, msg); err != nil {
		logger.Warn("post-zeroize operation rejected", "err", err)
	}
	logger.Info("hsm zeroized; demo complete")
}

func mustBytes(hex string) []byte {
	b, err := hexBytes(hex)
	if err != nil {
		panic(err)
	}
	return b
}

func hexBytes(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("odd hex length %d", len(s))
	}
	out := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		hi, ok := hexNibble(s[i])
		if !ok {
			return nil, fmt.Errorf("invalid hex char %q", s[i])
		}
		lo, ok := hexNibble(s[i+1])
		if !ok {
			return nil, fmt.Errorf("invalid hex char %q", s[i+1])
		}
		out[i/2] = hi<<4 | lo
	}
	return out, nil
}

func hexNibble(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	default:
		return 0, false
	}
}

const hexDigits = "0123456789abcdef"

func hexString(b []byte) string {
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexDigits[v>>4]
		out[i*2+1] = hexDigits[v&0x0F]
	}
	return string(out)
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
