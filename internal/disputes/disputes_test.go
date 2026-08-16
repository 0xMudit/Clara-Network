package disputes

import (
	"context"
	"testing"
	"time"
)

func testService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryStore())
}

func seedTx(t *testing.T, svc *Service, ref, merchant string, amount int64, credit bool) {
	t.Helper()
	if err := svc.RecordTransaction(context.Background(), MonitoredTransaction{
		RefID: ref, MerchantID: merchant, AmountMinor: amount, Currency: "840", IsCredit: credit,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatal(err)
	}
}

func file(t *testing.T, svc *Service, ref, merchant, holder, code string, amount int64) *Dispute {
	t.Helper()
	d, err := svc.File(context.Background(), FileRequest{
		RefID: ref, MerchantID: merchant, Cardholder: holder, Amount: amount, Currency: "840", ReasonCode: code,
	})
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestFileValid(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T1", "M-grocery", 5000, false)
	d := file(t, svc, "T1", "M-grocery", "Amy", "4837", 5000)
	if d.Stage != StageFiled || d.Status != StageFiled {
		t.Fatalf("stage = %s, want filed", d.Stage)
	}
	// Fraud category has a 20-day response window.
	if d.ResponseDue.Sub(d.FiledAt) < 19*24*time.Hour || d.ResponseDue.Sub(d.FiledAt) > 21*24*time.Hour {
		t.Fatalf("response window = %v, want ~20 days", d.ResponseDue.Sub(d.FiledAt))
	}
}

func TestFileUnknownTransaction(t *testing.T) {
	svc := testService(t)
	if _, err := svc.File(context.Background(), FileRequest{
		RefID: "NOPE", MerchantID: "M-grocery", Cardholder: "Amy", Amount: 5000, Currency: "840", ReasonCode: "4837",
	}); err == nil {
		t.Fatal("expected error for unknown transaction")
	}
}

func TestFileAssociatedTransactionCreditRejected(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T2", "M-grocery", 5000, true) // already refunded
	if _, err := svc.File(context.Background(), FileRequest{
		RefID: "T2", MerchantID: "M-grocery", Cardholder: "Bob", Amount: 5000, Currency: "840", ReasonCode: "4837",
	}); err == nil {
		t.Fatal("expected associated-transaction rejection for a credited transaction")
	}
}

func TestFileAmountMismatch(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T3", "M-grocery", 5000, false)
	if _, err := svc.File(context.Background(), FileRequest{
		RefID: "T3", MerchantID: "M-grocery", Cardholder: "Cara", Amount: 9999, Currency: "840", ReasonCode: "4837",
	}); err == nil {
		t.Fatal("expected amount-mismatch error")
	}
}

func TestRuleAcquirerWinsWithEvidence(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T4", "M-grocery", 8000, false)
	d := file(t, svc, "T4", "M-grocery", "Dana", "4837", 8000)
	if _, err := svc.Represent(context.Background(), d.ID, []string{"chip", "3ds", "receipt"}); err != nil {
		t.Fatal(err)
	}
	ruled, err := svc.Rule(context.Background(), d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ruled.Decision != DecisionRejected || ruled.Winner != WinnerAcquirer {
		t.Fatalf("decision = %s winner = %s, want rejected/acquirer", ruled.Decision, ruled.Winner)
	}
	if ruled.DisputeFee != svc.Config().DisputeFee {
		t.Fatalf("dispute fee = %d, want %d (issuer pays)", ruled.DisputeFee, svc.Config().DisputeFee)
	}
}

func TestRuleChargebackStands(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T5", "M-online", 12000, false)
	d := file(t, svc, "T5", "M-online", "Eve", "4870", 12000)
	// Evidence missing "delivery" -> acquirer loses.
	if _, err := svc.Represent(context.Background(), d.ID, []string{"receipt"}); err != nil {
		t.Fatal(err)
	}
	ruled, err := svc.Rule(context.Background(), d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ruled.Decision != DecisionAccepted || ruled.Winner != WinnerIssuer {
		t.Fatalf("decision = %s winner = %s, want accepted/issuer", ruled.Decision, ruled.Winner)
	}
}

func TestArbitration(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T6", "M-online", 20000, false)
	d := file(t, svc, "T6", "M-online", "Frank", "4841", 20000)
	if _, err := svc.Represent(context.Background(), d.ID, []string{"terms", "receipt"}); err != nil {
		t.Fatal(err)
	}
	ruled, err := svc.Rule(context.Background(), d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ruled.Decision != DecisionRejected {
		t.Fatalf("expected acquirer win before escalation, got %s", ruled.Decision)
	}
	if _, err := svc.Escalate(context.Background(), d.ID); err != nil {
		t.Fatal(err)
	}
	final, err := svc.Arbitrate(context.Background(), d.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	if final.Winner != WinnerIssuer || final.Decision != DecisionAccepted {
		t.Fatalf("winner = %s decision = %s, want issuer/accepted", final.Winner, final.Decision)
	}
	if final.ArbitrationFee != svc.Config().ArbitrationFee {
		t.Fatalf("arbitration fee = %d", final.ArbitrationFee)
	}
}

func TestOverdueSLA(t *testing.T) {
	svc := testService(t)
	seedTx(t, svc, "T7", "M-grocery", 3000, false)
	seedTx(t, svc, "T8", "M-grocery", 4000, false)
	d1 := file(t, svc, "T7", "M-grocery", "Grace", "4837", 3000)
	d2 := file(t, svc, "T8", "M-grocery", "Henry", "4831", 4000)
	// Backdate d1 past its response deadline to simulate an elapsed SLA.
	store := svc.store.(*MemoryStore)
	d1.ResponseDue = time.Now().UTC().Add(-time.Hour)
	if err := store.SaveDispute(context.Background(), *d1); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Represent(context.Background(), d2.ID, []string{"receipt", "avs"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Rule(context.Background(), d2.ID); err != nil {
		t.Fatal(err)
	}
	overdue, err := svc.Overdue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(overdue) != 1 || overdue[0].ID != d1.ID {
		t.Fatalf("expected only %s overdue, got %+v", d1.ID, overdue)
	}
}

func TestMonitorRatio(t *testing.T) {
	svc := testService(t)
	// 200 transactions for a grocery merchant, no disputes -> normal.
	for i := 0; i < 200; i++ {
		seedTx(t, svc, "G"+itoa(i), "M-grocery", 1000, false)
	}
	ratio, status, err := svc.MonitorRatio(context.Background(), "M-grocery")
	if err != nil {
		t.Fatal(err)
	}
	if status != "normal" || ratio != 0 {
		t.Fatalf("ratio = %v status = %s, want 0/normal", ratio, status)
	}

	// A high-risk merchant: 50 txs, 3 of them disputed and standing -> 6% excessive.
	for i := 0; i < 50; i++ {
		seedTx(t, svc, "H"+itoa(i), "M-casino", 5000, false)
	}
	for i := 0; i < 3; i++ {
		ref := "H" + itoa(10+i)
		d := file(t, svc, ref, "M-casino", "Ivy"+itoa(i), "4837", 5000)
		if _, err := svc.Represent(context.Background(), d.ID, []string{"receipt"}); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Rule(context.Background(), d.ID); err != nil {
			t.Fatal(err)
		}
	}
	ratio, status, err = svc.MonitorRatio(context.Background(), "M-casino")
	if err != nil {
		t.Fatal(err)
	}
	if ratio != 6.0 || status != "excessive" {
		t.Fatalf("ratio = %v status = %s, want 6/excessive", ratio, status)
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
