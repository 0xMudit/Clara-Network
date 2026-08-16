package clearing

import (
	"fmt"
	"time"
)

// Direction of a settlement instruction from the scheme to the settlement
// agent.
type Direction string

const (
	// DirDebit instructs the settlement agent to debit the member (it must
	// pay its net obligation into the scheme account).
	DirDebit Direction = "DEBIT"
	// DirCredit instructs the settlement agent to credit the member (it is a
	// net creditor and receives funds).
	DirCredit Direction = "CREDIT"
)

// SettlementInstruction is a single net settlement movement for a member.
type SettlementInstruction struct {
	CycleID     string
	MsgID       string
	Member      string
	Amount      int64 // minor units, always positive
	Direction   Direction
	Currency    string
	Instruction time.Time
	Final       bool
}

// PrefundAccount tracks a member's prefunded balance and its cap, which
// bounds the member's maximum net debit position (see docs/18 §18.4).
type PrefundAccount struct {
	Member  string
	Balance int64
	Cap     int64
}

// Validate checks a prefund account for obvious data errors. The cap bounds
// the maximum net debit a member may take on; the balance is what the member
// has funded and may legitimately exceed the cap after credits.
func (a PrefundAccount) Validate() error {
	if a.Member == "" {
		return fmt.Errorf("prefund: empty member")
	}
	if a.Balance < 0 || a.Cap < 0 {
		return fmt.Errorf("prefund: negative balance or cap")
	}
	return nil
}
