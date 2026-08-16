package instant

import (
	"context"
	"strings"
	"testing"
	"time"
)

func testEngine(t *testing.T) *Engine {
	t.Helper()
	e := NewEngine(Config{Currency: "USD"})
	if _, err := e.SetPosition("BANK-A", 100000); err != nil {
		t.Fatal(err)
	}
	if _, err := e.SetPosition("BANK-B", 50000); err != nil {
		t.Fatal(err)
	}
	return e
}

func pay(msgID, txID, sender, beneficiary string, amount int64) Payment {
	return Payment{
		MsgID: msgID, InstrID: msgID, EndToEndID: msgID + "-E2E", TxID: txID,
		Sender: sender, Beneficiary: beneficiary,
		SenderIBAN: "DE89370400440532013000", BeneficiaryIBAN: "GB29NWBK60161331926819",
		AmountMinor: amount, Currency: "USD", Remittance: "invoice 42",
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}

func TestHappyPathSettlesInRealTime(t *testing.T) {
	e := testEngine(t)
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "BANK-B", 2500))
	if res.Status != StatusACSC || !res.Final {
		t.Fatalf("expected ACSC final, got %+v", res)
	}
	if res.Positions["BANK-A"] != 97500 || res.Positions["BANK-B"] != 52500 {
		t.Fatalf("positions wrong: %v", res.Positions)
	}
	h := e.History()
	if len(h) != 1 || h[0].Status != StatusACSC || h[0].Reason != "" {
		t.Fatalf("history wrong: %+v", h)
	}
}

func TestInsufficientFundsRejected(t *testing.T) {
	e := testEngine(t)
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "BANK-B", 999999))
	if res.Status != StatusRJCT || res.Reason != ReasonInsufficientFunds {
		t.Fatalf("expected RJCT AC04, got %+v", res)
	}
	// Funds never move on rejection.
	if p, _ := e.Position("BANK-A"); p != 100000 {
		t.Fatalf("sender position changed on reject: %d", p)
	}
}

func TestUnknownBeneficiaryRejected(t *testing.T) {
	e := testEngine(t)
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "NO-SUCH", 2500))
	if res.Reason != ReasonAccount {
		t.Fatalf("expected AC01, got %+v", res)
	}
	if p, _ := e.Position("BANK-A"); p != 100000 {
		t.Fatalf("sender position changed on reject: %d", p)
	}
}

func TestUnknownSenderRejected(t *testing.T) {
	e := testEngine(t)
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "NO-SUCH", "BANK-B", 2500))
	if res.Reason != ReasonForbidden {
		t.Fatalf("expected AG01, got %+v", res)
	}
}

func TestSelfTransferRejected(t *testing.T) {
	e := testEngine(t)
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "BANK-A", 2500))
	if res.Reason != ReasonFormat {
		t.Fatalf("expected FF01 for self transfer, got %+v", res)
	}
}

func TestCurrencyMismatchRejected(t *testing.T) {
	e := testEngine(t)
	p := pay("M-1", "TX-1", "BANK-A", "BANK-B", 2500)
	p.Currency = "EUR"
	res := e.Transfer(context.Background(), p)
	if res.Reason != ReasonFormat {
		t.Fatalf("expected FF01 for currency mismatch, got %+v", res)
	}
}

func TestMissingFieldsRejected(t *testing.T) {
	e := testEngine(t)
	p := pay("M-1", "TX-1", "BANK-A", "BANK-B", 2500)
	p.TxID = ""
	if res := e.Transfer(context.Background(), p); res.Reason != ReasonFormat {
		t.Fatalf("expected FF01 for missing TxId, got %+v", res)
	}
}

func TestTimeoutReleasesReservation(t *testing.T) {
	e := NewEngine(Config{Currency: "USD", Timeout: 50 * time.Millisecond, Forward: func(ctx context.Context, _ Payment) error {
		<-ctx.Done()
		return ctx.Err()
	}})
	if _, err := e.SetPosition("BANK-A", 100000); err != nil {
		t.Fatal(err)
	}
	if _, err := e.SetPosition("BANK-B", 50000); err != nil {
		t.Fatal(err)
	}
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "BANK-B", 2500))
	if res.Reason != ReasonNoAnswer {
		t.Fatalf("expected NOAS, got %+v", res)
	}
	if p, _ := e.Position("BANK-A"); p != 100000 {
		t.Fatalf("reservation not released: %d", p)
	}
	if p, _ := e.Position("BANK-B"); p != 50000 {
		t.Fatalf("beneficiary position changed: %d", p)
	}
}

func TestForwardFailureReleasesReservation(t *testing.T) {
	e := NewEngine(Config{Currency: "USD", Forward: func(context.Context, Payment) error {
		return context.Canceled
	}})
	if _, err := e.SetPosition("BANK-A", 100000); err != nil {
		t.Fatal(err)
	}
	if _, err := e.SetPosition("BANK-B", 50000); err != nil {
		t.Fatal(err)
	}
	res := e.Transfer(context.Background(), pay("M-1", "TX-1", "BANK-A", "BANK-B", 2500))
	if res.Reason != ReasonForbidden {
		t.Fatalf("expected AG01, got %+v", res)
	}
	if p, _ := e.Position("BANK-A"); p != 100000 {
		t.Fatalf("reservation not released: %d", p)
	}
}

func TestPacs008RoundTrip(t *testing.T) {
	p := pay("M-42", "TX-42", "BANK-A", "BANK-B", 12345)
	raw, err := BuildPacs008(p)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParsePacs008(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.MsgID != p.MsgID || got.TxID != p.TxID || got.EndToEndID != p.EndToEndID ||
		got.Sender != p.Sender || got.Beneficiary != p.Beneficiary ||
		got.AmountMinor != 12345 || got.Currency != "USD" ||
		got.SenderIBAN != p.SenderIBAN || got.BeneficiaryIBAN != p.BeneficiaryIBAN ||
		got.Remittance != p.Remittance {
		t.Fatalf("round trip mismatch:\nwant %+v\n got %+v", p, got)
	}
}

func TestPacs008RejectsWrongNamespace(t *testing.T) {
	raw := []byte(`<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.009.001.08"><FICdtTrf/></Document>`)
	if _, err := ParsePacs008(raw); err == nil {
		t.Fatal("expected error for non-pacs.008 document")
	}
}

func TestPacs002BuildsStatus(t *testing.T) {
	raw, err := BuildPacs002(StatusReport{
		MsgID: "S-1", OriginalMsgID: "M-1", OriginalMsgName: "pacs.008.001.09",
		EndToEndID: "M-1-E2E", TxID: "TX-1", Status: StatusRJCT, Reason: ReasonInsufficientFunds,
		CreatedAt: "2026-08-16T00:00:00.000Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !strings.Contains(s, "pacs.002") || !strings.Contains(s, "RJCT") || !strings.Contains(s, ReasonInsufficientFunds) {
		t.Fatalf("status report missing elements: %s", s)
	}
}

func TestConcurrentTransfersReserveAtomically(t *testing.T) {
	e := NewEngine(Config{Currency: "USD"})
	if _, err := e.SetPosition("BANK-A", 100000); err != nil {
		t.Fatal(err)
	}
	if _, err := e.SetPosition("BANK-B", 50000); err != nil {
		t.Fatal(err)
	}
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			msgID := "M-" + string(rune('a'+i))
			res := e.Transfer(context.Background(), pay(msgID, "TX-"+msgID, "BANK-A", "BANK-B", 20000))
			done <- res.Status == StatusACSC
		}(i)
	}
	ok := 0
	for i := 0; i < 10; i++ {
		if <-done {
			ok++
		}
	}
	if ok != 5 {
		t.Fatalf("expected exactly 5 of 10 to settle (50,000.00 capacity), got %d", ok)
	}
	if p, _ := e.Position("BANK-A"); p != 0 {
		t.Fatalf("sender position should be exhausted, got %d", p)
	}
	if p, _ := e.Position("BANK-B"); p != 150000 {
		t.Fatalf("beneficiary should hold all funds, got %d", p)
	}
}
