// Package clearing implements the Clara Network clearing and net settlement
// layer: it captures clearing (financial) messages from acquirers, computes
// per-member net positions per settlement cycle, enforces prefunded caps,
// sizes and applies the default fund, and issues ISO 20022 pacs.009
// settlement instructions to the settlement agent.
//
// Money is always carried in minor units (cents) as int64 to avoid rounding.
package clearing

import (
	"fmt"
	"strings"
)

// SchemeOperatorID is the member identifier of the scheme's own account,
// which collects scheme fees as part of the same netting run.
const SchemeOperatorID = "CLARA"

// ClearingRecord is a single captured clearing (financial) message.
type ClearingRecord struct {
	CycleID     string
	STAN        string
	MTI         string // 0220/0221 clearing request, 0420 reversal
	Sender      string // acquirer receiving institution
	Receiver    string // issuer receiving institution
	AmountMinor int64
	Interchange int64 // interchange fee owed by the acquirer to the issuer
	Currency    string
	RefID       string // original authorization reference (RRN)
}

// Validate checks a clearing record for obvious data errors.
func (r ClearingRecord) Validate() error {
	switch r.MTI {
	case "0220", "0221", "0420":
	default:
		return fmt.Errorf("clearing: unsupported MTI %q", r.MTI)
	}
	if r.CycleID == "" || r.STAN == "" || r.RefID == "" {
		return fmt.Errorf("clearing: missing cycle/stan/reference")
	}
	if r.Sender == "" || r.Receiver == "" {
		return fmt.Errorf("clearing: missing sender or receiver")
	}
	if r.Sender == r.Receiver {
		return fmt.Errorf("clearing: sender %q cannot also be receiver", r.Sender)
	}
	if r.AmountMinor <= 0 {
		return fmt.Errorf("clearing: amount must be positive")
	}
	if r.Interchange < 0 || r.Interchange >= r.AmountMinor {
		return fmt.Errorf("clearing: interchange out of range [0,%d)", r.AmountMinor)
	}
	if len(r.Currency) != 3 || !isDigits(r.Currency) {
		return fmt.Errorf("clearing: invalid currency %q", r.Currency)
	}
	return nil
}

func isDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return strings.TrimLeft(s, "0") != ""
}

// FormatAmount renders minor units as major units with two decimals.
func FormatAmount(minor int64) string {
	sign := ""
	if minor < 0 {
		sign = "-"
		minor = -minor
	}
	return fmt.Sprintf("%s%d.%02d", sign, minor/100, minor%100)
}
