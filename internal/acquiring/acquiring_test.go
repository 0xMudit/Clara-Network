package acquiring

import (
	"context"
	"testing"
)

func testAcquirer(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryStore())
}

func seedLists(t *testing.T, store Store) {
	t.Helper()
	ctx := context.Background()
	if err := store.SaveMatchEntries(ctx, []MatchEntry{
		{MerchantName: "Quantum Trading LLC", TaxID: "12-3456789", Reason: "excessive chargebacks"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveOfacEntries(ctx, []OfacEntry{
		{Name: "Vladimir Petrov", Program: "SDN"},
	}); err != nil {
		t.Fatal(err)
	}
}

func baseApp() Application {
	return Application{
		MerchantName: "Bob's Grocery",
		DBA:          "Bob's Grocery",
		TaxID:        "98-7654321",
		Principals:   []string{"Robert Smith"},
		MCCs:         []string{"5411"},
		CreditScore:  720,
		Volume:       10_000_000,
	}
}

func TestBoardLowRiskApproved(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	seedLists(t, svc.store)
	m, d, err := svc.Board(ctx, baseApp())
	if err != nil {
		t.Fatal(err)
	}
	if d.Status != StatusActive {
		t.Fatalf("status = %s, want active", d.Status)
	}
	if m.RiskTier != TierLow {
		t.Fatalf("tier = %s, want low", m.RiskTier)
	}
	if m.ReserveRateBPS != 0 || m.FundingDelayDays != 0 {
		t.Fatalf("low-risk should have no reserve/delay, got %d/%d", m.ReserveRateBPS, m.FundingDelayDays)
	}
	if m.Status != StatusActive {
		t.Fatalf("merchant status = %s", m.Status)
	}
}

func TestBoardMATCHHitDeclined(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	seedLists(t, svc.store)
	app := baseApp()
	app.MerchantName = "Quantum Trading LLC"
	app.TaxID = "12-3456789"
	_, d, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	if d.Status != StatusDeclined {
		t.Fatalf("status = %s, want declined", d.Status)
	}
	found := false
	for _, r := range d.Reasons {
		if contains(r, "MATCH") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected MATCH reason, got %+v", d.Reasons)
	}
}

func TestBoardOFACHitDeclined(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	seedLists(t, svc.store)
	app := baseApp()
	app.Principals = []string{"Vladimir Petrov"}
	_, d, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	if d.Status != StatusDeclined {
		t.Fatalf("status = %s, want declined", d.Status)
	}
	found := false
	for _, r := range d.Reasons {
		if contains(r, "OFAC") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected OFAC reason, got %+v", d.Reasons)
	}
}

func TestBoardHighRiskRequiresEnhancedDD(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	app := baseApp()
	app.MerchantName = "Golden Spins Casino"
	app.MCCs = []string{"7995"}
	app.EnhancedDD = false

	_, d, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	if d.Status != StatusDeclined {
		t.Fatalf("status = %s, want declined without enhanced DD", d.Status)
	}

	app.EnhancedDD = true
	m, d2, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	if d2.Status != StatusActive {
		t.Fatalf("status = %s, want approved with enhanced DD", d2.Status)
	}
	if m.RiskTier != TierHigh {
		t.Fatalf("tier = %s, want high", m.RiskTier)
	}
	if m.ReserveRateBPS != 1000 || m.FundingDelayDays != 2 {
		t.Fatalf("high-risk mitigations = %d/%d, want 1000/2", m.ReserveRateBPS, m.FundingDelayDays)
	}
}

func TestBoardCreditBelowThreshold(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	app := baseApp()
	app.CreditScore = 600
	_, d, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	if d.Status != StatusDeclined {
		t.Fatalf("status = %s, want declined", d.Status)
	}
}

func TestSettleBatchLowRisk(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	m, _, err := svc.Board(ctx, baseApp())
	if err != nil {
		t.Fatal(err)
	}
	eng := NewFundingEngine(svc.store, svc.policy)
	line, err := eng.SettleBatch(ctx, m.ID, 1_000_000) // 10,000.00
	if err != nil {
		t.Fatal(err)
	}
	// 5411 MDR 150 bps (1.50%) + 1.00 fixed fee = 151.00.
	if line.Fees != 15_000+100 {
		t.Fatalf("fees = %d, want 15100", line.Fees)
	}
	if line.ReserveHold != 0 {
		t.Fatalf("low-risk reserve = %d, want 0", line.ReserveHold)
	}
	if line.Net != 1_000_000-15_100 {
		t.Fatalf("net = %d", line.Net)
	}
}

func TestSettleBatchHighRiskReserve(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	app := baseApp()
	app.MerchantName = "Golden Spins Casino"
	app.MCCs = []string{"7995"}
	app.EnhancedDD = true
	app.Volume = 20_000_000
	m, _, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	eng := NewFundingEngine(svc.store, svc.policy)
	line, err := eng.SettleBatch(ctx, m.ID, 500_000)
	if err != nil {
		t.Fatal(err)
	}
	// 7995 MDR 400 bps (4.00%) + 1.00; reserve 1000 bps = 10%.
	if line.ReserveHold != 50_000 {
		t.Fatalf("reserve = %d, want 50000", line.ReserveHold)
	}
	if line.Net != 500_000-20_100-50_000 {
		t.Fatalf("net = %d", line.Net)
	}
	after, err := svc.GetMerchant(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.ReserveBalance != 50_000 {
		t.Fatalf("merchant reserve balance = %d, want 50000", after.ReserveBalance)
	}
}

func TestSettleBatchOverLimit(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	app := baseApp()
	app.MCCs = []string{"7995"}
	app.EnhancedDD = true
	m, _, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	eng := NewFundingEngine(svc.store, svc.policy)
	if _, err := eng.SettleBatch(ctx, m.ID, 10_000_000); err == nil {
		t.Fatal("expected batch over the transaction limit to fail")
	}
}

func TestReserveRelease(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	app := baseApp()
	app.MerchantName = "Golden Spins Casino"
	app.MCCs = []string{"7995"}
	app.EnhancedDD = true
	app.Volume = 2_000_000 // reserve cap = 10% * 2,000,000 = 200,000
	m, _, err := svc.Board(ctx, app)
	if err != nil {
		t.Fatal(err)
	}
	eng := NewFundingEngine(svc.store, svc.policy)
	for i := 0; i < 5; i++ {
		if _, err := eng.SettleBatch(ctx, m.ID, 500_000); err != nil {
			t.Fatal(err)
		}
	}
	// 5 x 50,000 reserve = 250,000; cap 200,000 -> release 50,000.
	line, err := eng.ReleaseReserve(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if line.Net != 50_000 {
		t.Fatalf("released = %d, want 50000", line.Net)
	}
	after, err := svc.GetMerchant(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.ReserveBalance != 200_000 {
		t.Fatalf("reserve after release = %d, want 200000", after.ReserveBalance)
	}
	if _, err := eng.ReleaseReserve(ctx, m.ID); err == nil {
		t.Fatal("expected second release to fail at the cap")
	}
}

func TestTerminate(t *testing.T) {
	ctx := context.Background()
	svc := testAcquirer(t)
	m, _, err := svc.Board(ctx, baseApp())
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Terminate(ctx, m.ID, "fraud"); err != nil {
		t.Fatal(err)
	}
	after, err := svc.GetMerchant(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusTerminated {
		t.Fatalf("status = %s, want terminated", after.Status)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
