package ledger

import (
	"context"
	"fmt"
	"strings"
)

// Statement is one line of the settlement agent's statement for a cycle, e.g.
// the value the agent reports it actually moved for a member (mirrors the
// scheme's pacs.009 instructions; a real agent's statement may differ).
type Statement struct {
	CycleID   string
	Member    string
	Amount    int64 // minor units, always positive
	Direction Direction
}

// MismatchKind classifies a reconciliation discrepancy (docs/12 §12.9).
type MismatchKind string

const (
	// MissingInLedger: the statement shows a movement the ledger never posted.
	MissingInLedger MismatchKind = "missing-in-ledger"
	// AmountMismatch: statement and ledger agree on the member but not the amount.
	AmountMismatch MismatchKind = "amount-mismatch"
	// OrphanInLedger: the ledger posted a movement with no statement line.
	OrphanInLedger MismatchKind = "orphan-in-ledger"
	// Duplicate: more than one statement line for the same member.
	Duplicate MismatchKind = "duplicate"
)

// Mismatch is a single reconciliation discrepancy.
type Mismatch struct {
	Member   string
	Kind     MismatchKind
	Expected int64 // statement value (minor units, signed)
	Actual   int64 // ledger value (minor units, signed)
	Detail   string
}

// Report summarises a reconciliation of one cycle.
type Report struct {
	CycleID        string
	Statements     int
	LedgerPostings int
	Matched        int
	Mismatches     []Mismatch
	// StatementTotal is the signed sum of statement movements.
	StatementTotal int64
	// LedgerTotal is the signed sum of member movements posted for the cycle.
	LedgerTotal int64
}

// OK reports whether every statement line matched and the cycle totals agree.
func (r *Report) OK() bool {
	return len(r.Mismatches) == 0 && r.StatementTotal == r.LedgerTotal
}

func signed(amount int64, d Direction) int64 {
	if d == Debit {
		return -amount
	}
	return amount
}

func memberOf(reference, cycleID string) string {
	return strings.TrimPrefix(reference, cycleID+":")
}

// Reconciler matches a settlement agent's statement against the scheme's
// ledger for a cycle (docs/12 §12.9). It compares each member's net movement
// on its own ledger account against the statement line; a balanced journal's
// legs cancel out when summed, so only the member leg is counted.
type Reconciler struct {
	Store Store
	// MemberAccountPrefix maps a member ID to its ledger account, e.g. "M:".
	MemberAccountPrefix string
	// SchemeOperatorID is the operator's own member ID, whose internal fee
	// postings are settled within the scheme and never appear on a member
	// settlement statement. Empty disables the exclusion.
	SchemeOperatorID string
}

// Reconcile runs the reconciliation for the given statements. Discrepancies
// are classified per docs/12 §12.9: missing transaction, amount mismatch,
// orphan (in the ledger but not the statement), duplicate.
func (r *Reconciler) Reconcile(ctx context.Context, statements []Statement) (*Report, error) {
	report := &Report{Statements: len(statements)}
	perMember := map[string][]Statement{}
	for _, st := range statements {
		report.CycleID = st.CycleID
		report.StatementTotal += signed(st.Amount, st.Direction)
		perMember[st.Member] = append(perMember[st.Member], st)
	}

	for member, lines := range perMember {
		if len(lines) > 1 {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Member: member,
				Kind:   Duplicate,
				Detail: fmt.Sprintf("statement has %d lines for one member", len(lines)),
			})
		}

		expected := signed(lines[0].Amount, lines[0].Direction)
		actual, n, err := r.memberMovement(ctx, lines[0])
		if err != nil {
			return nil, err
		}
		report.LedgerPostings += n
		report.LedgerTotal += actual

		if n == 0 {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Member:   member,
				Kind:     MissingInLedger,
				Expected: expected,
				Detail:   "statement reports a movement with no ledger posting",
			})
			continue
		}
		if actual != expected {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Member:   member,
				Kind:     AmountMismatch,
				Expected: expected,
				Actual:   actual,
				Detail:   fmt.Sprintf("statement %d but ledger %d", expected, actual),
			})
			continue
		}
		report.Matched++
	}

	if err := r.detectOrphans(ctx, report, perMember); err != nil {
		return nil, err
	}
	return report, nil
}

// memberMovement returns the member's net ledger movement for a statement
// line: the signed sum of entries on the member's own account for the cycle.
func (r *Reconciler) memberMovement(ctx context.Context, st Statement) (int64, int, error) {
	accountID := r.MemberAccountPrefix + st.Member
	entries, err := r.Store.EntriesByAccountAndReference(ctx, accountID, reference(st.CycleID, st.Member))
	if err != nil {
		return 0, 0, fmt.Errorf("ledger: reconcile %s: %w", st.Member, err)
	}
	var total int64
	for _, e := range entries {
		total += signed(e.Amount, e.Direction)
	}
	return total, len(entries), nil
}

// detectOrphans flags ledger postings for the cycle whose member has no
// statement line.
func (r *Reconciler) detectOrphans(ctx context.Context, report *Report, perMember map[string][]Statement) error {
	refs, err := r.Store.ReferencesWithPrefix(ctx, report.CycleID+":")
	if err != nil {
		return fmt.Errorf("ledger: enumerate cycle journals: %w", err)
	}
	for _, ref := range refs {
		member := memberOf(ref, report.CycleID)
		if _, ok := perMember[member]; ok {
			continue
		}
		if r.SchemeOperatorID != "" && member == r.SchemeOperatorID {
			continue
		}
		accountID := r.MemberAccountPrefix + member
		entries, err := r.Store.EntriesByAccountAndReference(ctx, accountID, ref)
		if err != nil {
			return fmt.Errorf("ledger: load orphan %s: %w", ref, err)
		}
		var actual int64
		for _, e := range entries {
			actual += signed(e.Amount, e.Direction)
		}
		report.LedgerPostings += len(entries)
		report.LedgerTotal += actual
		report.Mismatches = append(report.Mismatches, Mismatch{
			Member: member,
			Kind:   OrphanInLedger,
			Actual: actual,
			Detail: fmt.Sprintf("ledger posted %s but the statement has no line for it", ref),
		})
	}
	return nil
}

func reference(cycleID, member string) string {
	return cycleID + ":" + member
}
