package ledger

import (
	"context"
	"testing"
)

func newTestLedger(t *testing.T) *Ledger {
	t.Helper()
	l := NewLedger(NewMemoryStore())
	ctx := context.Background()
	for _, a := range []Account{
		{ID: "M:ACQ-A", Type: AccountLiability, Currency: "840"},
		{ID: "M:ISS-C", Type: AccountLiability, Currency: "840"},
		{ID: "CASH", Type: AccountAsset, Currency: "840"},
		{ID: "INCOME:FEES", Type: AccountIncome, Currency: "840"},
	} {
		if err := l.EnsureAccount(ctx, a); err != nil {
			t.Fatalf("ensure account %s: %v", a.ID, err)
		}
	}
	return l
}

func TestPostBalanced(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	err := l.Post(ctx, "cyc:ACQ-A", []Entry{
		{AccountID: "M:ACQ-A", Direction: Debit, Amount: 12100, Currency: "840"},
		{AccountID: "CASH", Direction: Credit, Amount: 12100, Currency: "840"},
	})
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	bal, err := l.Balance(ctx, "M:ACQ-A")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if bal != -12100 {
		t.Fatalf("M:ACQ-A balance = %d, want -12100", bal)
	}
	if bal, _ = l.Balance(ctx, "CASH"); bal != 12100 {
		t.Fatalf("CASH balance = %d, want 12100", bal)
	}
}

func TestPostUnbalancedRejected(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	err := l.Post(ctx, "j1", []Entry{
		{AccountID: "M:ACQ-A", Direction: Debit, Amount: 100, Currency: "840"},
		{AccountID: "CASH", Direction: Credit, Amount: 90, Currency: "840"},
	})
	if err == nil {
		t.Fatal("expected unbalanced posting to be rejected")
	}
}

func TestPostCrossCurrencyBalanced(t *testing.T) {
	ctx := context.Background()
	l := NewLedger(NewMemoryStore())
	if err := l.EnsureAccount(ctx, Account{ID: "A", Currency: "840"}); err != nil {
		t.Fatal(err)
	}
	if err := l.EnsureAccount(ctx, Account{ID: "B", Currency: "978"}); err != nil {
		t.Fatal(err)
	}
	// Equal within each currency: 100 USD both ways, 50 EUR both ways.
	err := l.Post(ctx, "xr", []Entry{
		{AccountID: "A", Direction: Debit, Amount: 100, Currency: "840"},
		{AccountID: "A", Direction: Credit, Amount: 100, Currency: "840"},
		{AccountID: "B", Direction: Debit, Amount: 50, Currency: "978"},
		{AccountID: "B", Direction: Credit, Amount: 50, Currency: "978"},
	})
	if err != nil {
		t.Fatalf("cross-currency balanced posting rejected: %v", err)
	}
}

func TestPostUnknownAccountRejected(t *testing.T) {
	ctx := context.Background()
	l := NewLedger(NewMemoryStore())
	err := l.Post(ctx, "j2", []Entry{
		{AccountID: "NOPE", Direction: Debit, Amount: 5, Currency: "840"},
		{AccountID: "CASH", Direction: Credit, Amount: 5, Currency: "840"},
	})
	if err == nil {
		t.Fatal("expected unknown account to be rejected")
	}
}

func TestPostSingleEntryRejected(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	err := l.Post(ctx, "j3", []Entry{
		{AccountID: "CASH", Direction: Debit, Amount: 5, Currency: "840"},
	})
	if err == nil {
		t.Fatal("expected single-entry journal to be rejected")
	}
}

func TestPostNonPositiveAmountRejected(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	err := l.Post(ctx, "j4", []Entry{
		{AccountID: "M:ACQ-A", Direction: Debit, Amount: -5, Currency: "840"},
		{AccountID: "CASH", Direction: Credit, Amount: 5, Currency: "840"},
	})
	if err == nil {
		t.Fatal("expected negative amount to be rejected")
	}
}

func TestReversalRestoresBalance(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	post := func(entries ...Entry) error {
		return l.Post(ctx, "rev", entries)
	}
	if err := post(
		Entry{AccountID: "M:ISS-C", Direction: Debit, Amount: 5000, Currency: "840"},
		Entry{AccountID: "CASH", Direction: Credit, Amount: 5000, Currency: "840"},
	); err != nil {
		t.Fatal(err)
	}
	if err := post(
		Entry{AccountID: "CASH", Direction: Debit, Amount: 5000, Currency: "840"},
		Entry{AccountID: "M:ISS-C", Direction: Credit, Amount: 5000, Currency: "840"},
	); err != nil {
		t.Fatal(err)
	}
	if bal, _ := l.Balance(ctx, "M:ISS-C"); bal != 0 {
		t.Fatalf("after reversal M:ISS-C = %d, want 0", bal)
	}
	if bal, _ := l.Balance(ctx, "CASH"); bal != 0 {
		t.Fatalf("after reversal CASH = %d, want 0", bal)
	}
}

func TestTotalBalanceZero(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	// Full four-party settlement: ACQ-A and ACQ-B pay, ISS-C and ISS-D
	// receive, scheme keeps fees as income. The ledger must net to zero.
	postings := []struct {
		ref string
		es  []Entry
	}{
		{"cyc:ACQ-A", []Entry{
			{AccountID: "M:ACQ-A", Direction: Debit, Amount: 12100, Currency: "840"},
			{AccountID: "CASH", Direction: Credit, Amount: 12100, Currency: "840"},
		}},
		{"cyc:ACQ-B", []Entry{
			{AccountID: "M:ACQ-B", Direction: Debit, Amount: 9400, Currency: "840"},
			{AccountID: "CASH", Direction: Credit, Amount: 9400, Currency: "840"},
		}},
		{"cyc:ISS-C", []Entry{
			{AccountID: "CASH", Direction: Debit, Amount: 16200, Currency: "840"},
			{AccountID: "M:ISS-C", Direction: Credit, Amount: 16200, Currency: "840"},
		}},
		{"cyc:ISS-D", []Entry{
			{AccountID: "CASH", Direction: Debit, Amount: 5150, Currency: "840"},
			{AccountID: "M:ISS-D", Direction: Credit, Amount: 5150, Currency: "840"},
		}},
		{"cyc:CLARA", []Entry{
			{AccountID: "CASH", Direction: Debit, Amount: 150, Currency: "840"},
			{AccountID: "INCOME:FEES", Direction: Credit, Amount: 150, Currency: "840"},
		}},
	}
	for _, acct := range []string{"M:ACQ-B", "M:ISS-D"} {
		if err := l.EnsureAccount(ctx, Account{ID: acct, Type: AccountLiability, Currency: "840"}); err != nil {
			t.Fatal(err)
		}
	}
	for _, p := range postings {
		if err := l.Post(ctx, p.ref, p.es); err != nil {
			t.Fatalf("post %s: %v", p.ref, err)
		}
	}
	total, err := l.TotalBalance(ctx)
	if err != nil {
		t.Fatalf("total: %v", err)
	}
	if total != 0 {
		t.Fatalf("ledger total = %d, want 0", total)
	}
	if bal, _ := l.Balance(ctx, "INCOME:FEES"); bal != 150 {
		t.Fatalf("fee income = %d, want 150", bal)
	}
}

func newTestReconciler(l *Ledger) *Reconciler {
	return &Reconciler{Store: l.store, MemberAccountPrefix: "M:"}
}

func TestReconcileClean(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := newTestReconciler(l)
	statements := []Statement{
		{CycleID: "cyc", Member: "ACQ-A", Amount: 12100, Direction: Debit},
		{CycleID: "cyc", Member: "ISS-C", Amount: 16200, Direction: Credit},
	}
	for _, st := range statements {
		member := st.Member
		var es []Entry
		if st.Direction == Debit {
			es = []Entry{
				{AccountID: "M:" + member, Direction: Debit, Amount: st.Amount, Currency: "840"},
				{AccountID: "CASH", Direction: Credit, Amount: st.Amount, Currency: "840"},
			}
		} else {
			es = []Entry{
				{AccountID: "CASH", Direction: Debit, Amount: st.Amount, Currency: "840"},
				{AccountID: "M:" + member, Direction: Credit, Amount: st.Amount, Currency: "840"},
			}
		}
		if err := l.Post(ctx, "cyc:"+member, es); err != nil {
			t.Fatal(err)
		}
	}
	rep, err := r.Reconcile(ctx, statements)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if !rep.OK() {
		t.Fatalf("expected clean reconciliation, got mismatches %+v", rep.Mismatches)
	}
	if rep.Matched != 2 {
		t.Fatalf("matched = %d, want 2", rep.Matched)
	}
}

func TestReconcileAmountMismatch(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := newTestReconciler(l)
	if err := l.Post(ctx, "cyc:ISS-C", []Entry{
		{AccountID: "CASH", Direction: Debit, Amount: 16000, Currency: "840"},
		{AccountID: "M:ISS-C", Direction: Credit, Amount: 16000, Currency: "840"},
	}); err != nil {
		t.Fatal(err)
	}
	rep, err := r.Reconcile(ctx, []Statement{
		{CycleID: "cyc", Member: "ISS-C", Amount: 16200, Direction: Credit},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rep.OK() {
		t.Fatal("expected reconciliation to fail")
	}
	found := false
	for _, m := range rep.Mismatches {
		if m.Kind == AmountMismatch && m.Expected == 16200 && m.Actual == 16000 {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected amount-mismatch, got %+v", rep.Mismatches)
	}
}

func TestReconcileMissingInLedger(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := newTestReconciler(l)
	rep, err := r.Reconcile(ctx, []Statement{
		{CycleID: "cyc", Member: "ACQ-A", Amount: 12100, Direction: Debit},
	})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range rep.Mismatches {
		if m.Kind == MissingInLedger {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected missing-in-ledger, got %+v", rep.Mismatches)
	}
}

func TestReconcileDuplicate(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := newTestReconciler(l)
	rep, err := r.Reconcile(ctx, []Statement{
		{CycleID: "cyc", Member: "ACQ-A", Amount: 500, Direction: Debit},
		{CycleID: "cyc", Member: "ACQ-A", Amount: 500, Direction: Debit},
	})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range rep.Mismatches {
		if m.Kind == Duplicate {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected duplicate, got %+v", rep.Mismatches)
	}
}

func TestReconcileSkipsSchemeOperator(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := &Reconciler{Store: l.store, MemberAccountPrefix: "M:", SchemeOperatorID: "CLARA"}
	// Member statement for the cycle, so orphan detection runs over cycle refs.
	if err := l.Post(ctx, "cyc:ISS-C", []Entry{
		{AccountID: "CASH", Direction: Debit, Amount: 16200, Currency: "840"},
		{AccountID: "M:ISS-C", Direction: Credit, Amount: 16200, Currency: "840"},
	}); err != nil {
		t.Fatal(err)
	}
	// The operator's internal fee posting: no settlement statement line.
	if err := l.Post(ctx, "cyc:CLARA", []Entry{
		{AccountID: "CASH", Direction: Debit, Amount: 150, Currency: "840"},
		{AccountID: "INCOME:FEES", Direction: Credit, Amount: 150, Currency: "840"},
	}); err != nil {
		t.Fatal(err)
	}
	rep, err := r.Reconcile(ctx, []Statement{
		{CycleID: "cyc", Member: "ISS-C", Amount: 16200, Direction: Credit},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range rep.Mismatches {
		if m.Kind == OrphanInLedger {
			t.Fatalf("operator posting flagged as orphan: %+v", m)
		}
	}
	if !rep.OK() {
		t.Fatalf("expected clean reconciliation, got %+v", rep.Mismatches)
	}
}

func TestReconcileOrphanInLedger(t *testing.T) {
	ctx := context.Background()
	l := newTestLedger(t)
	r := newTestReconciler(l)
	if err := l.Post(ctx, "cyc:ACQ-A", []Entry{
		{AccountID: "M:ACQ-A", Direction: Debit, Amount: 12100, Currency: "840"},
		{AccountID: "CASH", Direction: Credit, Amount: 12100, Currency: "840"},
	}); err != nil {
		t.Fatal(err)
	}
	// The statement omits ACQ-A entirely.
	rep, err := r.Reconcile(ctx, []Statement{
		{CycleID: "cyc", Member: "ISS-C", Amount: 16200, Direction: Credit},
	})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range rep.Mismatches {
		if m.Kind == OrphanInLedger && m.Member == "ACQ-A" && m.Actual == -12100 {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected orphan-in-ledger, got %+v", rep.Mismatches)
	}
}
