package clearing

import "sort"

// NetPosition is a member's net settlement amount for a cycle. A positive
// value means the member is a net creditor; negative means a net debtor.
type NetPosition struct {
	CycleID string
	Member  string
	Net     int64 // minor units
}

// NetAmountDue returns the amount the member must pay (as a positive value)
// when they are a net debtor.
func (p NetPosition) NetAmountDue() int64 {
	if p.Net < 0 {
		return -p.Net
	}
	return 0
}

// NetPositions aggregates clearing records into per-member net positions.
// Settlement math (four-party model): for each transaction of amount A and
// interchange I, the issuer reimburses the acquirer A-I; the acquirer also
// pays a scheme fee. The scheme operator therefore nets a positive position.
func NetPositions(records []ClearingRecord, schemeFee int64) []NetPosition {
	if len(records) == 0 {
		return nil
	}
	net := map[string]int64{}
	for _, r := range records {
		settle := r.AmountMinor - r.Interchange
		net[r.Receiver] += settle
		net[r.Sender] -= settle
		net[r.Sender] -= schemeFee
		net[SchemeOperatorID] += schemeFee
	}

	members := make([]string, 0, len(net))
	for m := range net {
		members = append(members, m)
	}
	sort.Strings(members)

	positions := make([]NetPosition, 0, len(members))
	for _, m := range members {
		positions = append(positions, NetPosition{CycleID: records[0].CycleID, Member: m, Net: net[m]})
	}
	return positions
}

// TotalNet is the sum of all net positions; it must be zero for a valid
// cycle (a ledger invariant).
func TotalNet(positions []NetPosition) int64 {
	var total int64
	for _, p := range positions {
		total += p.Net
	}
	return total
}
