package cardsvc

import (
	"bytes"
	"context"
	"encoding/hex"
	"testing"
)

const testMasterKeyHex = "2b7e151628aed2a6abf7158809cf4f3c"

var testPAN = "4000001234567899"

func testService(t *testing.T) *Service {
	t.Helper()
	key, err := hex.DecodeString(testMasterKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	svc, err := NewService(NewMemoryStore(), Config{IssuerMasterKey: key})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := svc.AddBinRange(ctx, BinRange{BIN: "400000", Low: 0, High: 9999999999, Currency: "840", Product: "classic"}); err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestCMACKnownVector(t *testing.T) {
	key, _ := hex.DecodeString("2b7e151628aed2a6abf7158809cf4f3c")
	cases := []struct {
		msg string
		tag string
	}{
		{"", "bb1d6929e95937287fa37d129b756746"},
		{"6bc1bee22e409f96e93d7e117393172a", "070a16b46b4d4144f79bdd9dd04a287c"},
		{"6bc1bee22e409f96e93d7e117393172aae2d8a571e03ac9c9eb76fac45af8e5130c81c46a35ce411", "dfa66747de9ae63030ca32611497c827"},
	}
	for _, c := range cases {
		msg, _ := hex.DecodeString(c.msg)
		tag, err := cmac(key, msg)
		if err != nil {
			t.Fatal(err)
		}
		if got := hex.EncodeToString(tag); got != c.tag {
			t.Fatalf("cmac(%q) = %s, want %s", c.msg, got, c.tag)
		}
	}
}

func TestLuhn(t *testing.T) {
	if !ValidLuhn(testPAN) {
		t.Fatalf("pan %s should be Luhn-valid", testPAN)
	}
	if ValidLuhn("4000001234567890") {
		t.Fatal("4000001234567890 should fail Luhn")
	}
	if ValidLuhn("123") {
		t.Fatal("short pan should fail Luhn")
	}
}

func TestDeriveCardKeyDeterministic(t *testing.T) {
	key, _ := hex.DecodeString(testMasterKeyHex)
	k1, err := DeriveCardKey(key, testPAN, 1)
	if err != nil {
		t.Fatal(err)
	}
	k2, _ := DeriveCardKey(key, testPAN, 1)
	if !bytes.Equal(k1, k2) {
		t.Fatal("card key derivation must be deterministic")
	}
	k3, _ := DeriveCardKey(key, "4000001234567888", 1)
	if bytes.Equal(k1, k3) {
		t.Fatal("different PANs must derive different keys")
	}
}

func TestARQCVerifyValid(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	card, err := svc.CreateCard(ctx, testPAN, "3012", "")
	if err != nil {
		t.Fatal(err)
	}
	data := ARQCData{Amount: 250000, Currency: "840", STAN: "100000", Date: "260816", ATC: 7, UN: "1234"}
	arqc, err := svc.ComputeARQC(ctx, card.Ref, data)
	if err != nil {
		t.Fatal(err)
	}
	valid, err := svc.VerifyARQC(ctx, card.Ref, data, arqc)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("cryptogram should verify")
	}
}

func TestARQCVerifyTamperedAmount(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	card, _ := svc.CreateCard(ctx, testPAN, "3012", "")
	data := ARQCData{Amount: 250000, Currency: "840", STAN: "100000", Date: "260816", ATC: 1, UN: "1234"}
	arqc, _ := svc.ComputeARQC(ctx, card.Ref, data)
	// The cryptogram was computed over a different amount: must fail.
	attacker := data
	attacker.Amount = 1
	valid, err := svc.VerifyARQC(ctx, card.Ref, attacker, arqc)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("tampered cryptogram must not verify")
	}
}

func TestARQCReplayRejected(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	card, _ := svc.CreateCard(ctx, testPAN, "3012", "")
	data := ARQCData{Amount: 250000, Currency: "840", STAN: "100000", Date: "260816", ATC: 3, UN: "1234"}
	arqc, _ := svc.ComputeARQC(ctx, card.Ref, data)
	if valid, _ := svc.VerifyARQC(ctx, card.Ref, data, arqc); !valid {
		t.Fatal("first presentation must verify")
	}
	// Replaying the same cryptogram (same ATC) must be rejected.
	if valid, _ := svc.VerifyARQC(ctx, card.Ref, data, arqc); valid {
		t.Fatal("replayed cryptogram must be rejected")
	}
	// An older ATC must also be rejected.
	old := data
	old.ATC = 2
	oldArqc, _ := svc.ComputeARQC(ctx, card.Ref, old)
	if valid, _ := svc.VerifyARQC(ctx, card.Ref, old, oldArqc); valid {
		t.Fatal("cryptogram with stale ATC must be rejected")
	}
}

func TestCreateCard(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	card, err := svc.CreateCard(ctx, testPAN, "3012", "platinum")
	if err != nil {
		t.Fatal(err)
	}
	if card.BIN != "400000" {
		t.Fatalf("bin = %s, want 400000", card.BIN)
	}
	if card.Status != StatusActive {
		t.Fatalf("status = %s", card.Status)
	}
	if card.PANMask != "400000******7899" {
		t.Fatalf("mask = %s", card.PANMask)
	}
	if card.Product != "platinum" {
		t.Fatalf("product = %s", card.Product)
	}
	if len(card.UDK) != cmacSize {
		t.Fatalf("udk length = %d", len(card.UDK))
	}
}

func TestCreateCardRejectsUnissuedBIN(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	// 400111... is outside the issued 400000 range.
	_, err := svc.CreateCard(ctx, "4001111234567899", "3012", "")
	if err == nil {
		t.Fatal("expected unissued BIN to be rejected")
	}
}

func TestCreateCardRejectsBadLuhn(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	_, err := svc.CreateCard(ctx, "4000001234567890", "3012", "")
	if err == nil {
		t.Fatal("expected bad Luhn to be rejected")
	}
}

func TestTokenizeDetokenize(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	card, _ := svc.CreateCard(ctx, testPAN, "3012", "")
	vault := NewTokenVault(NewMemoryStore())

	tok, err := vault.Tokenize(ctx, svc, testPAN)
	if err != nil {
		t.Fatal(err)
	}
	if !ValidLuhn(tok.Number) {
		t.Fatal("token must be Luhn-valid")
	}
	if len(tok.PAR) != parLength {
		t.Fatalf("PAR length = %d", len(tok.PAR))
	}

	// Tokenizing again returns the same token + PAR (token reuse).
	again, err := vault.Tokenize(ctx, svc, testPAN)
	if err != nil {
		t.Fatal(err)
	}
	if again.Number != tok.Number || again.PAR != tok.PAR {
		t.Fatal("tokenize must be stable for a PAN")
	}

	hash, err := vault.Detokenize(ctx, tok.Number)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(hash, card.PANHash) {
		t.Fatal("detokenize returned wrong PAN hash")
	}
}

func TestTokenizeUnknownCardRejected(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	vault := NewTokenVault(NewMemoryStore())
	_, err := vault.Tokenize(ctx, svc, "4000001111222233") // Luhn? not the point
	if err == nil {
		t.Fatal("expected tokenize of unknown card to be rejected")
	}
}

func TestProvision(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	if _, err := svc.CreateCard(ctx, testPAN, "3012", ""); err != nil {
		t.Fatal(err)
	}
	vault := NewTokenVault(NewMemoryStore())
	tok, _ := vault.Tokenize(ctx, svc, testPAN)

	p, err := vault.Provision(ctx, tok.Number, "device-42", "TRID001")
	if err != nil {
		t.Fatal(err)
	}
	if p.Token != tok.Number || p.PAR != tok.PAR {
		t.Fatal("provision payload must carry the token and PAR")
	}
	if p.DeviceID != "device-42" {
		t.Fatalf("device = %s", p.DeviceID)
	}
}

func TestPARsAreUnique(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	for _, pan := range []string{"4000001234567899", "4000001234567873"} {
		if !ValidLuhn(pan) {
			t.Fatalf("test pan %s must be Luhn-valid", pan)
		}
	}
	vault := NewTokenVault(NewMemoryStore())
	seen := map[string]bool{}
	for _, pan := range []string{"4000001234567899", "4000001234567873"} {
		if _, err := svc.CreateCard(ctx, pan, "3012", ""); err != nil {
			t.Fatal(err)
		}
		tok, err := vault.Tokenize(ctx, svc, pan)
		if err != nil {
			t.Fatal(err)
		}
		if seen[tok.PAR] {
			t.Fatal("PARs must be unique per PAN")
		}
		seen[tok.PAR] = true
	}
}

func TestFindBinRange(t *testing.T) {
	svc := testService(t)
	ctx := context.Background()
	r, ok, err := svc.FindBinRange(ctx, testPAN)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || r.BIN != "400000" || r.Currency != "840" {
		t.Fatalf("unexpected range %+v ok=%v", r, ok)
	}
	if _, ok, _ := svc.FindBinRange(ctx, "5000001234567899"); ok {
		t.Fatal("500000 must not be in an issued range")
	}
}
