package acquiring

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// MerchantStatus values.
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusDeclined  = "declined"
	StatusTerminated = "terminated"
)

// Application is a merchant boarding application (docs/23 §23.2).
type Application struct {
	MerchantName string
	DBA          string
	TaxID        string
	Principals   []string
	MCCs         []string
	CreditScore  int
	EnhancedDD   bool
	Volume       int64 // projected monthly volume (minor units)
}

// Policy holds the underwriting decision parameters.
type Policy struct {
	CreditThreshold   int
	ReserveRateBPS    map[string]int64
	FundingDelayDays  map[string]int
	TransactionLimit  map[string]int64
	FixedFeeMinor     int64
}

// DefaultPolicy returns the scheme's default underwriting policy.
func DefaultPolicy() Policy {
	return Policy{
		CreditThreshold: 650,
		ReserveRateBPS: map[string]int64{
			TierLow: 0, TierMedium: 500, TierHigh: 1000,
		},
		FundingDelayDays: map[string]int{
			TierLow: 0, TierMedium: 1, TierHigh: 2,
		},
		TransactionLimit: map[string]int64{
			TierLow: 5_000_000, TierMedium: 2_000_000, TierHigh: 500_000,
		},
		FixedFeeMinor: 100,
	}
}

// Decision is the underwriting outcome (docs/23 §23.3): approve, decline, or
// conditional approval with exposure mitigations.
type Decision struct {
	Status            string
	RiskTier          string
	ReserveRateBPS    int64
	FundingDelayDays  int
	TransactionLimit  int64
	Reasons           []string
}

// Merchant is a boarded merchant.
type Merchant struct {
	ID               string
	Name             string
	DBA              string
	TaxID            string
	Principals       []string
	MCCs             []string
	Status           string
	RiskTier         string
	ReserveRateBPS   int64
	FundingDelayDays int
	TransactionLimit int64
	ReserveBalance   int64
	Volume           int64
	DeclineReason    string
	ApprovedAt       time.Time
}

// Service boards merchants and manages their lifecycle.
type Service struct {
	store  Store
	policy Policy
	log    *slog.Logger
}

// NewService builds an acquiring service over the store.
func NewService(store Store) *Service {
	return &Service{store: store, policy: DefaultPolicy(), log: slog.Default()}
}

// WithPolicy overrides the underwriting policy.
func (s *Service) WithPolicy(p Policy) *Service {
	s.policy = p
	return s
}

// Policy returns the service's underwriting policy.
func (s *Service) Policy() Policy {
	return s.policy
}

// Board underwrites an application: negative-list screening, MCC assignment,
// credit assessment, and conditional-approval mitigations (docs/23 §23.2-23.3).
func (s *Service) Board(ctx context.Context, app Application) (*Merchant, *Decision, error) {
	if len(app.MCCs) == 0 {
		return nil, nil, fmt.Errorf("acquiring: application has no MCC")
	}
	primary, ok := LookupMCC(app.MCCs[0])
	if !ok {
		return nil, nil, fmt.Errorf("acquiring: unknown MCC %s", app.MCCs[0])
	}

	decision := &Decision{Status: StatusActive, RiskTier: primary.Tier}

	screener := NewScreener(s.store)
	hits, err := screener.Screen(ctx, app)
	if err != nil {
		return nil, nil, err
	}
	for _, h := range hits {
		decision.Reasons = append(decision.Reasons, "screening hit on "+h.List+": "+h.Detail)
	}
	if len(hits) > 0 {
		return nil, decline(decision, "negative-list screening"), nil
	}

	// Risk tier is the highest across all assigned MCCs.
	tier := primary.Tier
	eddRequired := false
	for _, code := range app.MCCs {
		m, ok := LookupMCC(code)
		if !ok {
			return nil, nil, fmt.Errorf("acquiring: unknown MCC %s", code)
		}
		if rank(m.Tier) > rank(tier) {
			tier = m.Tier
		}
		eddRequired = eddRequired || m.EnhancedDD
	}
	decision.RiskTier = tier

	// Enhanced due diligence is mandatory for high-integrity-risk MCCs.
	if eddRequired && !app.EnhancedDD {
		decision.Reasons = append(decision.Reasons, "enhanced due diligence required for high-risk MCC")
		return nil, decline(decision, "enhanced due diligence required"), nil
	}
	if app.CreditScore < s.policy.CreditThreshold {
		decision.Reasons = append(decision.Reasons, fmt.Sprintf("credit score %d below threshold %d", app.CreditScore, s.policy.CreditThreshold))
		return nil, decline(decision, "credit assessment failed"), nil
	}

	decision.ReserveRateBPS = s.policy.ReserveRateBPS[tier]
	decision.FundingDelayDays = s.policy.FundingDelayDays[tier]
	decision.TransactionLimit = s.policy.TransactionLimit[tier]
	if tier != TierLow {
		decision.Reasons = append(decision.Reasons,
			fmt.Sprintf("conditional approval: reserve %d bps, funding delay %d day(s), transaction limit %s",
				decision.ReserveRateBPS, decision.FundingDelayDays, amount(decision.TransactionLimit)))
	}

	merchant := &Merchant{
		ID:               merchantID(app),
		Name:             app.MerchantName,
		DBA:              orDefault(app.DBA, app.MerchantName),
		TaxID:            app.TaxID,
		Principals:       app.Principals,
		MCCs:             app.MCCs,
		Status:           StatusActive,
		RiskTier:         tier,
		ReserveRateBPS:   decision.ReserveRateBPS,
		FundingDelayDays: decision.FundingDelayDays,
		TransactionLimit: decision.TransactionLimit,
		Volume:           app.Volume,
		ApprovedAt:       time.Now().UTC(),
	}
	if err := s.store.SaveMerchant(ctx, *merchant); err != nil {
		return nil, nil, fmt.Errorf("acquiring: persist merchant: %w", err)
	}
	s.log.Info("merchant boarded", "id", merchant.ID, "tier", tier,
		"reserve_bps", decision.ReserveRateBPS, "delay_days", decision.FundingDelayDays)
	return merchant, decision, nil
}

// GetMerchant loads a merchant by ID.
func (s *Service) GetMerchant(ctx context.Context, id string) (*Merchant, error) {
	m, ok, err := s.store.Merchant(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("acquiring: unknown merchant %s", id)
	}
	return &m, nil
}

// Terminate ends a merchant relationship and, in a real deployment, reports
// it to the MATCH database (docs/23 §23.4).
func (s *Service) Terminate(ctx context.Context, id, reason string) error {
	m, err := s.GetMerchant(ctx, id)
	if err != nil {
		return err
	}
	m.Status = StatusTerminated
	m.DeclineReason = reason
	if err := s.store.SaveMerchant(ctx, *m); err != nil {
		return err
	}
	s.log.Warn("merchant terminated", "id", id, "reason", reason)
	return nil
}

func decline(d *Decision, reason string) *Decision {
	d.Status = StatusDeclined
	d.Reasons = append(d.Reasons, reason)
	return d
}

func rank(tier string) int {
	switch tier {
	case TierHigh:
		return 3
	case TierMedium:
		return 2
	default:
		return 1
	}
}

func merchantID(app Application) string {
	name := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(app.MerchantName), " ", "-"))
	return "M-" + name
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
