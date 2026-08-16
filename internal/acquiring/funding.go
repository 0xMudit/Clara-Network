package acquiring

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// FundingLine is one merchant settlement funding batch (docs/23 §23.3): the
// gross sales amount less processing fees and the reserve hold, on the funding
// date set by the merchant's boarding tier.
type FundingLine struct {
	BatchID     string
	MerchantID  string
	Date        time.Time
	Gross       int64
	Fees        int64
	ReserveHold int64
	Net         int64
}

// ReserveCap returns the merchant's reserve release target: the reserve rate
// applied to its projected monthly volume. The reserve balance may exceed it;
// the excess is released back to the merchant.
func (m *Merchant) ReserveCap() int64 {
	if m.ReserveRateBPS <= 0 {
		return 0
	}
	return m.ReserveRateBPS * m.Volume / 10_000
}

// FundingEngine computes merchant settlements and manages reserves.
type FundingEngine struct {
	store   Store
	policy  Policy
	log     *slog.Logger
}

// NewFundingEngine builds a funding engine.
func NewFundingEngine(store Store, policy Policy) *FundingEngine {
	return &FundingEngine{store: store, policy: policy, log: slog.Default()}
}

// SettleBatch funds a merchant's gross receipts for a batch: gross − MDR fees
// − reserve hold, on the tier's funding delay.
func (f *FundingEngine) SettleBatch(ctx context.Context, merchantID string, gross int64) (*FundingLine, error) {
	if gross <= 0 {
		return nil, fmt.Errorf("acquiring: batch gross must be positive")
	}
	m, err := f.merchant(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	if m.Status != StatusActive {
		return nil, fmt.Errorf("acquiring: merchant %s is %s", merchantID, m.Status)
	}
	if gross > m.TransactionLimit {
		return nil, fmt.Errorf("acquiring: batch gross %s exceeds transaction limit %s for %s",
			amount(gross), amount(m.TransactionLimit), merchantID)
	}

	mcc, _ := LookupMCC(m.MCCs[0])
	feeRate := mcc.RateBPS
	fees := gross*feeRate/10_000 + f.policy.FixedFeeMinor
	hold := gross * m.ReserveRateBPS / 10_000
	net := gross - fees - hold

	line := FundingLine{
		BatchID:     fmt.Sprintf("B-%s-%d", merchantID, time.Now().UTC().UnixNano()),
		MerchantID:  merchantID,
		Date:        time.Now().UTC().AddDate(0, 0, m.FundingDelayDays),
		Gross:       gross,
		Fees:        fees,
		ReserveHold: hold,
		Net:         net,
	}

	m.ReserveBalance += hold
	if err := f.store.SaveMerchant(ctx, *m); err != nil {
		return nil, fmt.Errorf("acquiring: persist reserve: %w", err)
	}
	if err := f.store.AppendFunding(ctx, line); err != nil {
		return nil, fmt.Errorf("acquiring: persist funding line: %w", err)
	}
	f.log.Info("merchant funded", "merchant", merchantID, "batch", line.BatchID,
		"gross", amount(gross), "fees", amount(fees), "reserve_hold", amount(hold), "net", amount(net))
	return &line, nil
}

// ReleaseReserve returns the reserve balance in excess of the target to the
// merchant as a funding line (docs/23 §23.3 reserves).
func (f *FundingEngine) ReleaseReserve(ctx context.Context, merchantID string) (*FundingLine, error) {
	m, err := f.merchant(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	cap := m.ReserveCap()
	if m.ReserveBalance <= cap {
		return nil, fmt.Errorf("acquiring: merchant %s reserve %s not above target %s",
			merchantID, amount(m.ReserveBalance), amount(cap))
	}
	excess := m.ReserveBalance - cap
	line := FundingLine{
		BatchID:     fmt.Sprintf("R-%s-%d", merchantID, time.Now().UTC().UnixNano()),
		MerchantID:  merchantID,
		Date:        time.Now().UTC().AddDate(0, 0, m.FundingDelayDays),
		Gross:       excess,
		Fees:        0,
		ReserveHold: 0,
		Net:         excess,
	}
	m.ReserveBalance = cap
	if err := f.store.SaveMerchant(ctx, *m); err != nil {
		return nil, fmt.Errorf("acquiring: persist reserve release: %w", err)
	}
	if err := f.store.AppendFunding(ctx, line); err != nil {
		return nil, fmt.Errorf("acquiring: persist release line: %w", err)
	}
	f.log.Info("reserve released", "merchant", merchantID, "amount", amount(excess), "reserve_after", amount(m.ReserveBalance))
	return &line, nil
}

func (f *FundingEngine) merchant(ctx context.Context, id string) (*Merchant, error) {
	m, ok, err := f.store.Merchant(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("acquiring: unknown merchant %s", id)
	}
	return &m, nil
}

// Amount renders minor units as major units with two decimals.
func Amount(minor int64) string {
	sign := ""
	if minor < 0 {
		sign = "-"
		minor = -minor
	}
	return fmt.Sprintf("%s%d.%02d", sign, minor/100, minor%100)
}

func amount(minor int64) string {
	return Amount(minor)
}
