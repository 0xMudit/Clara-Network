package clearing

import (
	"context"
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

const (
	memberA = "A00001"
	memberB = "B00002"
	cycleID = "20260816"
)

func record(sender, receiver string, amount, interchange int64) ClearingRecord {
	return ClearingRecord{
		STAN:        "100000",
		MTI:         "0221",
		Sender:      sender,
		Receiver:    receiver,
		AmountMinor: amount,
		Interchange: interchange,
		Currency:    "840",
		RefID:       "REF-1",
	}
}

func TestNetPositions(t *testing.T) {
	records := []ClearingRecord{
		record(memberA, memberB, 1000, 200), // B +800, A -825 (incl 25 fee), CLARA +25
		record(memberB, memberA, 5000, 500), // A +4500, B -4525, CLARA +25
	}
	positions := NetPositions(records, 25)

	got := map[string]int64{}
	for _, p := range positions {
		got[p.Member] = p.Net
	}
	if got[memberA] != 3675 {
		t.Fatalf("A net = %d want 3675", got[memberA])
	}
	if got[memberB] != -3725 {
		t.Fatalf("B net = %d want -3725", got[memberB])
	}
	if got[SchemeOperatorID] != 50 {
		t.Fatalf("CLARA net = %d want 50", got[SchemeOperatorID])
	}
	if total := TotalNet(positions); total != 0 {
		t.Fatalf("positions do not balance: %d", total)
	}
}

func TestNetPositionsEmpty(t *testing.T) {
	if positions := NetPositions(nil, 25); positions != nil {
		t.Fatalf("expected nil positions, got %v", positions)
	}
}

func TestSettleNoDefaults(t *testing.T) {
	positions := []NetPosition{
		{CycleID: cycleID, Member: memberA, Net: -4450},
		{CycleID: cycleID, Member: memberB, Net: 4400},
		{CycleID: cycleID, Member: SchemeOperatorID, Net: 50},
	}
	accounts := map[string]PrefundAccount{
		memberA: {Member: memberA, Balance: 5000, Cap: 5000},
		memberB: {Member: memberB, Balance: 1000, Cap: 1000},
	}
	res, err := Settle(cycleID, positions, accounts, 100000, "840", time.Now())
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	if !res.Final {
		t.Fatal("cycle should settle with no defaults")
	}
	if len(res.Events) != 0 {
		t.Fatalf("unexpected default events: %v", res.Events)
	}
	if got := res.Accounts[memberA].Balance; got != 550 {
		t.Fatalf("A balance = %d want 550", got)
	}
	if got := res.Accounts[memberB].Balance; got != 5400 {
		t.Fatalf("B balance = %d want 5400", got)
	}
	if got := res.DFBalance; got != 100000 {
		t.Fatalf("DF balance changed without defaults: %d", got)
	}
	if len(res.Instructions) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(res.Instructions))
	}
}

func TestSettleDefaultFundCoversShortfall(t *testing.T) {
	positions := []NetPosition{
		{CycleID: cycleID, Member: memberA, Net: -30000},
		{CycleID: cycleID, Member: memberB, Net: 29500},
		{CycleID: cycleID, Member: SchemeOperatorID, Net: 500},
	}
	accounts := map[string]PrefundAccount{
		memberA: {Member: memberA, Balance: 10000, Cap: 10000},
		memberB: {Member: memberB, Balance: 1000, Cap: 1000},
	}
	res, err := Settle(cycleID, positions, accounts, 50000, "840", time.Now())
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	if !res.Final {
		t.Fatal("cycle should settle because DF covers the shortfall")
	}
	if len(res.Events) != 1 {
		t.Fatalf("expected one default event, got %v", res.Events)
	}
	ev := res.Events[0]
	if ev.Member != memberA || ev.Shortfall != 20000 || ev.Covered != 20000 || ev.Uncovered != 0 {
		t.Fatalf("unexpected event: %+v", ev)
	}
	if got := res.DFBalance; got != 30000 {
		t.Fatalf("DF balance = %d want 30000", got)
	}
	if got := res.Accounts[memberA].Balance; got != 0 {
		t.Fatalf("A balance should be depleted, got %d", got)
	}
}

func TestSettleUncovered(t *testing.T) {
	positions := []NetPosition{
		{CycleID: cycleID, Member: memberA, Net: -30000},
		{CycleID: cycleID, Member: memberB, Net: 29500},
		{CycleID: cycleID, Member: SchemeOperatorID, Net: 500},
	}
	accounts := map[string]PrefundAccount{
		memberA: {Member: memberA, Balance: 10000, Cap: 10000},
	}
	res, err := Settle(cycleID, positions, accounts, 5000, "840", time.Now())
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	if res.Final {
		t.Fatal("cycle must not be final with uncovered shortfall")
	}
	ev := res.Events[0]
	if ev.Covered != 5000 || ev.Uncovered != 15000 {
		t.Fatalf("expected covered=5000 uncovered=15000, got %+v", ev)
	}
	if res.DFBalance != 0 {
		t.Fatalf("DF balance should be exhausted, got %d", res.DFBalance)
	}
}

func TestTargetDefaultFund(t *testing.T) {
	positions := []NetPosition{
		{CycleID: cycleID, Member: memberA, Net: -4450},
		{CycleID: cycleID, Member: memberB, Net: 4400},
	}
	if got := TargetDefaultFund(positions, 2.0); got != 8900 {
		t.Fatalf("target = %d want 8900", got)
	}
	if got := TargetDefaultFund(nil, 2.0); got != 0 {
		t.Fatalf("empty target = %d want 0", got)
	}
}

func TestPacs009Credit(t *testing.T) {
	inst := SettlementInstruction{
		CycleID: cycleID, MsgID: cycleID + "-" + memberB,
		Member: memberB, Amount: 4400, Direction: DirCredit,
		Currency: "840", Instruction: time.Now(),
	}
	raw, err := Pacs009XML(inst)
	if err != nil {
		t.Fatalf("pacs.009: %v", err)
	}
	s := string(raw)
	for _, want := range []string{"pacs.009.001.08", "CLRG", "44.00", memberB} {
		if !strings.Contains(s, want) {
			t.Fatalf("pacs.009 missing %q:\n%s", want, s)
		}
	}
	var doc struct {
		XMLName xml.Name `xml:"Document"`
	}
	if err := xml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("pacs.009 is not well-formed XML: %v", err)
	}
}

func TestPacs009Debit(t *testing.T) {
	inst := SettlementInstruction{
		CycleID: cycleID, MsgID: cycleID + "-" + memberA,
		Member: memberA, Amount: 4450, Direction: DirDebit,
		Currency: "840", Instruction: time.Now(),
	}
	raw, err := Pacs009XML(inst)
	if err != nil {
		t.Fatalf("pacs.009: %v", err)
	}
	s := string(raw)
	if !strings.Contains(s, SchemeOperatorID) {
		t.Fatalf("debit instruction should pay into the scheme omnibus:\n%s", s)
	}
}

func TestSubmitBatchValidation(t *testing.T) {
	svc := NewService(NewMemoryStore(), Config{})
	bad := record(memberA, memberB, 1000, 200)
	bad.MTI = "9999"
	if err := svc.SubmitBatch(context.Background(), cycleID, []ClearingRecord{bad}); err == nil {
		t.Fatal("expected validation error")
	}
	if err := svc.SubmitBatch(context.Background(), cycleID, []ClearingRecord{record(memberA, memberA, 1000, 200)}); err == nil {
		t.Fatal("expected sender==receiver validation error")
	}
}

func TestRunCycleEndToEnd(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	svc := NewService(store, Config{})

	if err := svc.SubmitBatch(ctx, cycleID, []ClearingRecord{
		record(memberA, memberB, 1000, 200),
		record(memberA, memberB, 4000, 400),
	}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if err := svc.Fund(ctx, PrefundAccount{Member: memberA, Balance: 5000, Cap: 5000}); err != nil {
		t.Fatalf("fund A: %v", err)
	}
	if err := svc.Fund(ctx, PrefundAccount{Member: memberB, Balance: 1000, Cap: 1000}); err != nil {
		t.Fatalf("fund B: %v", err)
	}
	if err := svc.FundDefaultFund(ctx, 100000); err != nil {
		t.Fatalf("fund DF: %v", err)
	}

	res, err := svc.RunCycle(ctx, cycleID, time.Now())
	if err != nil {
		t.Fatalf("run cycle: %v", err)
	}
	if !res.Final {
		t.Fatalf("cycle should be final: %+v", res.Events)
	}
	// A owes 4400 + 2*25 fee = 4450; B is owed 4400.
	if res.Accounts[memberA].Balance != 550 {
		t.Fatalf("A balance = %d want 550", res.Accounts[memberA].Balance)
	}
	if res.Accounts[memberB].Balance != 5400 {
		t.Fatalf("B balance = %d want 5400", res.Accounts[memberB].Balance)
	}
	if len(res.Instructions) != 2 {
		t.Fatalf("expected 2 settlement instructions, got %d", len(res.Instructions))
	}
	if res.DefaultFundTarget != 8900 {
		t.Fatalf("DF target = %d want 8900", res.DefaultFundTarget)
	}

	// Positions persisted in the store.
	positions, err := store.Records(ctx, cycleID)
	if err != nil || len(positions) != 2 {
		t.Fatalf("records persisted: %d %v", len(positions), err)
	}
	// A's updated prefund balance persisted.
	a, ok, err := store.PrefundAccount(ctx, memberA)
	if err != nil || !ok || a.Balance != 550 {
		t.Fatalf("prefund A persisted: %+v %v %v", a, ok, err)
	}
}

func TestRunCycleNoRecords(t *testing.T) {
	svc := NewService(NewMemoryStore(), Config{})
	if _, err := svc.RunCycle(context.Background(), "empty", time.Now()); err == nil {
		t.Fatal("expected error for empty cycle")
	}
}
