package disputes

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Dispute stages and statuses.
const (
	StageFiled         = "filed"
	StageRepresentment = "representment"
	StageArbitration   = "arbitration"
	StageResolved      = "resolved"
	StageInvalid       = "invalid"
)

// Decisions.
const (
	DecisionAccepted = "accepted" // chargeback stands; issuer wins
	DecisionRejected = "rejected" // representment successful; acquirer wins
	DecisionArbitred = "arbitred"
)

// Winner labels.
const (
	WinnerIssuer   = "issuer"
	WinnerAcquirer = "acquirer"
)

// Dispute is one dispute case through its lifecycle: filed by the issuer,
// represented by the acquirer, ruled by the scheme, with optional escalation
// to arbitration (docs/20 §20.1).
type Dispute struct {
	ID             string
	RefID          string
	MerchantID     string
	Cardholder     string
	AmountMinor    int64
	Currency       string
	ReasonCode     string
	Category       string
	Stage          string
	Status         string
	FiledAt        time.Time
	ResponseDue    time.Time
	RespondedAt    time.Time
	EscalatedAt    time.Time
	Evidence       []string
	Decision       string
	Winner         string
	DecisionAt     time.Time
	DisputeFee     int64
	ArbitrationFee int64
	Note           string
}

// MonitoredTransaction is a transaction counted toward a merchant's chargeback
// ratio and checked for prior credits (associated-transaction rule).
type MonitoredTransaction struct {
	RefID       string
	MerchantID  string
	AmountMinor int64
	Currency    string
	IsCredit    bool
	CreatedAt   time.Time
}

// Config tunes the dispute framework.
type Config struct {
	DisputeFee           int64  // charged to the losing party
	ArbitrationFee       int64  // charged to the losing party
	MonitorWindowDays    int    // rolling window for chargeback ratios
	MonitorWatchedRatio  int    // percent at which monitoring starts
	MonitorExcessiveRatio int   // percent at which a merchant is excessive
}

// DefaultConfig returns the scheme's default dispute configuration.
func DefaultConfig() Config {
	return Config{
		DisputeFee:           2500,
		ArbitrationFee:       10000,
		MonitorWindowDays:    90,
		MonitorWatchedRatio:  50, // 0.50%
		MonitorExcessiveRatio: 100, // 1.00%
	}
}

// FileRequest opens a dispute.
type FileRequest struct {
	RefID      string
	MerchantID string
	Cardholder string
	Amount     int64
	Currency   string
	ReasonCode string
}

// Service drives the dispute lifecycle.
type Service struct {
	store  Store
	cfg    Config
	log    *slog.Logger
}

// NewService builds a disputes service over the store.
func NewService(store Store) *Service {
	return &Service{store: store, cfg: DefaultConfig(), log: slog.Default()}
}

// WithConfig overrides the framework configuration.
func (s *Service) WithConfig(c Config) *Service {
	s.cfg = c
	return s
}

// Config returns the current configuration.
func (s *Service) Config() Config {
	return s.cfg
}

// RecordTransaction feeds the monitoring/associated-transaction data.
func (s *Service) RecordTransaction(ctx context.Context, tx MonitoredTransaction) error {
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = time.Now().UTC()
	}
	return s.store.SaveTransaction(ctx, tx)
}

// File validates and opens a dispute (docs/20 §20.1): the referenced
// transaction must exist and must not already have been credited, or the
// dispute is rejected as invalid (associated-transaction check).
func (s *Service) File(ctx context.Context, req FileRequest) (*Dispute, error) {
	rc, ok := Lookup(req.ReasonCode)
	if !ok {
		return nil, fmt.Errorf("disputes: unknown reason code %s", req.ReasonCode)
	}

	tx, found, err := s.store.Transaction(ctx, req.RefID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("disputes: unknown transaction %s", req.RefID)
	}
	if tx.IsCredit {
		return nil, fmt.Errorf("disputes: transaction %s already credited (associated-transaction rule)", req.RefID)
	}
	if tx.AmountMinor != req.Amount {
		return nil, fmt.Errorf("disputes: amount %d differs from transaction %d", req.Amount, tx.AmountMinor)
	}

	now := time.Now().UTC()
	d := &Dispute{
		ID:          disputeID(req),
		RefID:       req.RefID,
		MerchantID:  tx.MerchantID,
		Cardholder:  req.Cardholder,
		AmountMinor: tx.AmountMinor,
		Currency:    tx.Currency,
		ReasonCode:  rc.Code,
		Category:    rc.Category,
		Stage:       StageFiled,
		Status:      StageFiled,
		FiledAt:     now,
		ResponseDue: now.AddDate(0, 0, ResponseDays(rc.Category)),
	}
	if err := s.store.SaveDispute(ctx, *d); err != nil {
		return nil, fmt.Errorf("disputes: persist dispute: %w", err)
	}
	s.log.Info("dispute filed", "id", d.ID, "ref", d.RefID, "merchant", d.MerchantID,
		"reason", rc.Code, "due", d.ResponseDue.Format("2006-01-02"))
	return d, nil
}

// Represent submits the acquirer's representment evidence (docs/20 §20.3).
func (s *Service) Represent(ctx context.Context, id string, evidence []string) (*Dispute, error) {
	d, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if d.Stage != StageFiled {
		return nil, fmt.Errorf("disputes: %s is %s, cannot represent", id, d.Stage)
	}
	d.Stage = StageRepresentment
	d.Status = StageRepresentment
	d.RespondedAt = time.Now().UTC()
	d.Evidence = evidence
	if err := s.store.SaveDispute(ctx, *d); err != nil {
		return nil, err
	}
	s.log.Info("representment submitted", "id", id, "evidence", evidence)
	return d, nil
}

// Rule is the scheme's decision: the acquirer wins if its evidence covers the
// reason code's requirements, otherwise the chargeback stands (docs/20 §20.3).
// The losing party pays the dispute fee (docs/20 §20.5).
func (s *Service) Rule(ctx context.Context, id string) (*Dispute, error) {
	d, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if d.Stage != StageRepresentment {
		return nil, fmt.Errorf("disputes: %s is %s, cannot rule yet", id, d.Stage)
	}
	rc := ReasonCodes[d.ReasonCode]
	if covers(d.Evidence, rc.RequiredEvidence) {
		d.Decision = DecisionRejected
		d.Winner = WinnerAcquirer
		d.DisputeFee = s.cfg.DisputeFee // issuer pays
		d.Note = "representment evidence satisfied " + d.ReasonCode + " requirements"
	} else {
		d.Decision = DecisionAccepted
		d.Winner = WinnerIssuer
		d.DisputeFee = s.cfg.DisputeFee // acquirer/merchant pays
		d.Note = "chargeback stands; evidence missing " + strings.Join(rc.RequiredEvidence, ", ")
	}
	d.Stage = StageResolved
	d.Status = StageResolved
	d.DecisionAt = time.Now().UTC()
	if err := s.store.SaveDispute(ctx, *d); err != nil {
		return nil, err
	}
	s.log.Info("dispute ruled", "id", id, "decision", d.Decision, "winner", d.Winner,
		"fee", s.cfg.DisputeFee)
	return d, nil
}

// Escalate moves a rejected dispute to arbitration (docs/20 §20.4).
func (s *Service) Escalate(ctx context.Context, id string) (*Dispute, error) {
	d, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if d.Stage != StageResolved || d.Decision != DecisionRejected {
		return nil, fmt.Errorf("disputes: %s cannot be escalated (stage=%s decision=%s)",
			id, d.Stage, d.Decision)
	}
	d.Stage = StageArbitration
	d.Status = StageArbitration
	d.EscalatedAt = time.Now().UTC()
	d.ResponseDue = d.EscalatedAt.AddDate(0, 0, 14)
	if err := s.store.SaveDispute(ctx, *d); err != nil {
		return nil, err
	}
	s.log.Warn("dispute escalated to arbitration", "id", id)
	return d, nil
}

// Arbitrate issues the final arbitration decision; the losing party pays the
// arbitration fee (docs/20 §20.5).
func (s *Service) Arbitrate(ctx context.Context, id string, forIssuer bool) (*Dispute, error) {
	d, err := s.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if d.Stage != StageArbitration {
		return nil, fmt.Errorf("disputes: %s is not in arbitration", id)
	}
	if forIssuer {
		d.Decision = DecisionAccepted
		d.Winner = WinnerIssuer
	} else {
		d.Decision = DecisionRejected
		d.Winner = WinnerAcquirer
	}
	d.ArbitrationFee = s.cfg.ArbitrationFee
	d.Stage = StageResolved
	d.Status = StageResolved
	d.DecisionAt = time.Now().UTC()
	d.Note = "final arbitration decision"
	if err := s.store.SaveDispute(ctx, *d); err != nil {
		return nil, err
	}
	s.log.Warn("arbitration decided", "id", id, "winner", d.Winner,
		"arbitration_fee", s.cfg.ArbitrationFee)
	return d, nil
}

// Overdue lists open disputes past their response deadline (SLA adherence).
func (s *Service) Overdue(ctx context.Context) ([]Dispute, error) {
	open, err := s.store.OpenDisputes(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	var out []Dispute
	for _, d := range open {
		if d.ResponseDue.Before(now) {
			out = append(out, d)
		}
	}
	return out, nil
}

// MonitorRatio computes a merchant's chargeback ratio: standing chargebacks
// divided by transactions over the rolling window (docs/20 §20.4).
func (s *Service) MonitorRatio(ctx context.Context, merchantID string) (float64, string, error) {
	txs, err := s.store.Transactions(ctx, merchantID)
	if err != nil {
		return 0, "", err
	}
	ds, err := s.store.Disputes(ctx, merchantID)
	if err != nil {
		return 0, "", err
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -s.cfg.MonitorWindowDays)
	var txCount, cb int64
	for _, tx := range txs {
		if tx.CreatedAt.After(cutoff) && !tx.IsCredit {
			txCount++
		}
	}
	for _, d := range ds {
		if d.Decision == DecisionAccepted && d.DecisionAt.After(cutoff) {
			cb++
		}
	}
	ratio := 0.0
	if txCount > 0 {
		ratio = float64(cb) * 100 / float64(txCount)
	}
	status := "normal"
	if ratio >= float64(s.cfg.MonitorExcessiveRatio)/100 {
		status = "excessive"
	} else if ratio >= float64(s.cfg.MonitorWatchedRatio)/100 {
		status = "watched"
	}
	return ratio, status, nil
}

func (s *Service) get(ctx context.Context, id string) (*Dispute, error) {
	d, ok, err := s.store.Dispute(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("disputes: unknown dispute %s", id)
	}
	return &d, nil
}

func covers(have, need []string) bool {
	for _, n := range need {
		found := false
		for _, h := range have {
			if strings.EqualFold(strings.TrimSpace(h), strings.TrimSpace(n)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func disputeID(req FileRequest) string {
	name := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Cardholder), " ", "-"))
	return "D-" + name + "-" + req.RefID + "-" + req.ReasonCode
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
