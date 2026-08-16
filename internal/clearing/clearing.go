package clearing

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Config configures the clearing and settlement service.
type Config struct {
	SchemeFee         int64   // per-transaction scheme fee charged to acquirers (minor units)
	Currency          string  // settlement currency code, e.g. "840" (USD)
	DefaultFundFactor float64 // default fund target = factor x largest net debit
	Log               *slog.Logger
}

// Service runs clearing cycles for the scheme.
type Service struct {
	cfg   Config
	store Store
	log   *slog.Logger
}

// NewService builds a clearing service over the given store.
func NewService(store Store, cfg Config) *Service {
	if cfg.SchemeFee <= 0 {
		cfg.SchemeFee = 25
	}
	if cfg.Currency == "" {
		cfg.Currency = "840"
	}
	if cfg.DefaultFundFactor <= 0 {
		cfg.DefaultFundFactor = 2.0
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Service{cfg: cfg, store: store, log: cfg.Log}
}

// SubmitBatch validates and captures a clearing (clearing-file) batch under
// the given cycle.
func (s *Service) SubmitBatch(ctx context.Context, cycleID string, records []ClearingRecord) error {
	for i := range records {
		records[i].CycleID = cycleID
		if err := records[i].Validate(); err != nil {
			return fmt.Errorf("clearing: record %d: %w", i, err)
		}
	}
	if err := s.store.AppendRecords(ctx, records); err != nil {
		return fmt.Errorf("clearing: persist batch: %w", err)
	}
	return nil
}

// Fund sets or updates a member's prefunded account.
func (s *Service) Fund(ctx context.Context, acct PrefundAccount) error {
	if err := acct.Validate(); err != nil {
		return err
	}
	return s.store.SetPrefund(ctx, acct)
}

// FundDefaultFund sets the default fund balance.
func (s *Service) FundDefaultFund(ctx context.Context, amount int64) error {
	if amount < 0 {
		return fmt.Errorf("clearing: negative default fund")
	}
	return s.store.SetDefaultFundBalance(ctx, amount)
}

// CycleResult is the outcome of a settlement cycle.
type CycleResult struct {
	CycleID             string
	Positions           []NetPosition
	Instructions        []SettlementInstruction
	Events              []DefaultEvent
	Accounts            map[string]PrefundAccount
	DefaultFundBalance  int64
	DefaultFundTarget   int64
	Final               bool
}

// RunCycle nets all captured records for cycleID, applies prefunded accounts
// and the default fund, and persists the result.
func (s *Service) RunCycle(ctx context.Context, cycleID string, asOf time.Time) (*CycleResult, error) {
	records, err := s.store.Records(ctx, cycleID)
	if err != nil {
		return nil, fmt.Errorf("clearing: load records: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("clearing: no records for cycle %q", cycleID)
	}

	positions := NetPositions(records, s.cfg.SchemeFee)
	if total := TotalNet(positions); total != 0 {
		return nil, fmt.Errorf("clearing: net positions do not balance (%d)", total)
	}

	accounts := map[string]PrefundAccount{}
	for _, p := range positions {
		if p.Member == SchemeOperatorID {
			continue
		}
		a, ok, err := s.store.PrefundAccount(ctx, p.Member)
		if err != nil {
			return nil, fmt.Errorf("clearing: load prefund %s: %w", p.Member, err)
		}
		if !ok {
			a = PrefundAccount{Member: p.Member}
		}
		accounts[p.Member] = a
	}

	df, err := s.store.DefaultFundBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("clearing: load default fund: %w", err)
	}

	res, err := Settle(cycleID, positions, accounts, df, s.cfg.Currency, asOf)
	if err != nil {
		return nil, fmt.Errorf("clearing: settle: %w", err)
	}

	if err := s.store.SavePositions(ctx, positions); err != nil {
		return nil, fmt.Errorf("clearing: persist positions: %w", err)
	}
	if err := s.store.SaveInstructions(ctx, res.Instructions); err != nil {
		return nil, fmt.Errorf("clearing: persist instructions: %w", err)
	}
	for _, a := range res.Accounts {
		if err := s.store.SetPrefund(ctx, a); err != nil {
			return nil, fmt.Errorf("clearing: persist prefund: %w", err)
		}
	}
	if err := s.store.SetDefaultFundBalance(ctx, res.DFBalance); err != nil {
		return nil, fmt.Errorf("clearing: persist default fund: %w", err)
	}

	return &CycleResult{
		CycleID:            cycleID,
		Positions:          positions,
		Instructions:       res.Instructions,
		Events:             res.Events,
		Accounts:           res.Accounts,
		DefaultFundBalance: res.DFBalance,
		DefaultFundTarget:  TargetDefaultFund(positions, s.cfg.DefaultFundFactor),
		Final:              res.Final,
	}, nil
}
