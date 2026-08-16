package clearing

import "time"

// DefaultEvent records that a member could not meet its net obligation from
// its prefunded balance and the default fund (and, failing that, the scheme)
// had to cover the shortfall.
type DefaultEvent struct {
	CycleID   string
	Member    string
	Shortfall int64 // obligation the member's prefund could not cover
	Covered   int64 // amount the default fund absorbed
	Uncovered int64 // amount neither prefund nor default fund covered
}

// TargetDefaultFund sizes the default fund to cover the default of the
// largest member's net debit position scaled by factor (docs/18 §18.5).
func TargetDefaultFund(positions []NetPosition, factor float64) int64 {
	var maxDebit int64
	for _, p := range positions {
		if due := p.NetAmountDue(); due > maxDebit {
			maxDebit = due
		}
	}
	if maxDebit == 0 {
		return 0
	}
	target := int64(float64(maxDebit)*factor/100) * 100
	if target < 100 {
		target = 100
	}
	return target
}

// SettleResult is the outcome of applying net positions to the prefunded
// accounts.
type SettleResult struct {
	Accounts     map[string]PrefundAccount
	Instructions []SettlementInstruction
	Events       []DefaultEvent
	DFBalance    int64
	// Final is true when every member's obligation was fully met.
	Final bool
}

// Settle applies the net positions to the prefunded accounts using the
// default-fund waterfall: a member's prefunded balance first, then the
// default fund, then uncovered (scheme capital / pro-rata allocation — not
// applied here, surfaced as DefaultEvent.Uncovered).
func Settle(cycleID string, positions []NetPosition, accounts map[string]PrefundAccount, dfBalance int64, currency string, asOf time.Time) (*SettleResult, error) {
	next := make(map[string]PrefundAccount, len(accounts))
	for m, a := range accounts {
		if err := a.Validate(); err != nil {
			return nil, err
		}
		next[m] = a
	}

	res := &SettleResult{
		Accounts:  next,
		DFBalance: dfBalance,
		Final:     true,
	}

	for _, p := range positions {
		if p.Net == 0 {
			continue
		}
		// The scheme operator's fee income nets into its own settlement
		// account; it does not produce an RTGS instruction.
		if p.Member == SchemeOperatorID {
			continue
		}

		acct := next[p.Member]

		inst := SettlementInstruction{
			CycleID:     cycleID,
			MsgID:       cycleID + "-" + p.Member,
			Member:      p.Member,
			Amount:      p.Net,
			Currency:    currency,
			Instruction: asOf,
			Final:       true,
		}
		if inst.Amount < 0 {
			inst.Amount = -inst.Amount
		}

		if p.Net > 0 {
			// Net creditor: the scheme pays out.
			inst.Direction = DirCredit
			acct.Balance += p.Net
			next[p.Member] = acct
			res.Instructions = append(res.Instructions, inst)
			continue
		}

		// Net debtor: prefund balance first, then the default fund.
		inst.Direction = DirDebit
		required := -p.Net
		fromPrefund := acct.Balance
		if fromPrefund > required {
			fromPrefund = required
		}
		acct.Balance -= fromPrefund
		shortfall := required - fromPrefund

		covered := int64(0)
		uncovered := int64(0)
		if shortfall > 0 {
			covered = res.DFBalance
			if covered > shortfall {
				covered = shortfall
			}
			res.DFBalance -= covered
			uncovered = shortfall - covered
			res.Final = uncovered == 0
			res.Events = append(res.Events, DefaultEvent{
				CycleID:   cycleID,
				Member:    p.Member,
				Shortfall: shortfall,
				Covered:   covered,
				Uncovered: uncovered,
			})
		}
		next[p.Member] = acct
		res.Instructions = append(res.Instructions, inst)
	}

	return res, nil
}
