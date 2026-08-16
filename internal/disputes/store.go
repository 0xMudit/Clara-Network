package disputes

import (
	"context"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store persists disputes and monitored transactions.
type Store interface {
	SaveDispute(ctx context.Context, d Dispute) error
	Dispute(ctx context.Context, id string) (Dispute, bool, error)
	Disputes(ctx context.Context, merchantID string) ([]Dispute, error)
	OpenDisputes(ctx context.Context) ([]Dispute, error)
	SaveTransaction(ctx context.Context, tx MonitoredTransaction) error
	Transaction(ctx context.Context, refID string) (MonitoredTransaction, bool, error)
	Transactions(ctx context.Context, merchantID string) ([]MonitoredTransaction, error)
}

// MemoryStore is an in-process Store for tests and single-instance runs.
type MemoryStore struct {
	mu       sync.Mutex
	disputes map[string]Dispute
	txs      map[string]MonitoredTransaction
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{disputes: map[string]Dispute{}, txs: map[string]MonitoredTransaction{}}
}

// SaveDispute implements Store.
func (m *MemoryStore) SaveDispute(_ context.Context, d Dispute) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.disputes[d.ID] = d
	return nil
}

// Dispute implements Store.
func (m *MemoryStore) Dispute(_ context.Context, id string) (Dispute, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.disputes[id]
	return d, ok, nil
}

// Disputes implements Store.
func (m *MemoryStore) Disputes(_ context.Context, merchantID string) ([]Dispute, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Dispute
	for _, d := range m.disputes {
		if d.MerchantID == merchantID {
			out = append(out, d)
		}
	}
	return out, nil
}

// OpenDisputes implements Store.
func (m *MemoryStore) OpenDisputes(_ context.Context) ([]Dispute, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Dispute
	for _, d := range m.disputes {
		if d.Stage != StageResolved && d.Stage != StageInvalid {
			out = append(out, d)
		}
	}
	return out, nil
}

// SaveTransaction implements Store.
func (m *MemoryStore) SaveTransaction(_ context.Context, tx MonitoredTransaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txs[tx.RefID] = tx
	return nil
}

// Transaction implements Store.
func (m *MemoryStore) Transaction(_ context.Context, refID string) (MonitoredTransaction, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tx, ok := m.txs[refID]
	return tx, ok, nil
}

// Transactions implements Store.
func (m *MemoryStore) Transactions(_ context.Context, merchantID string) ([]MonitoredTransaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []MonitoredTransaction
	for _, tx := range m.txs {
		if tx.MerchantID == merchantID {
			out = append(out, tx)
		}
	}
	return out, nil
}

// PostgresStore persists the disputes engine in PostgreSQL.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// SaveDispute implements Store.
func (p *PostgresStore) SaveDispute(ctx context.Context, d Dispute) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO disputes
		 (id, ref_id, merchant_id, cardholder, amount_minor, currency, reason_code, category,
		  stage, status, filed_at, response_due, responded_at, escalated_at, evidence,
		  decision, winner, decision_at, dispute_fee, arbitration_fee, note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
		 ON CONFLICT (id) DO UPDATE SET stage = EXCLUDED.stage, status = EXCLUDED.status,
		   responded_at = EXCLUDED.responded_at, escalated_at = EXCLUDED.escalated_at,
		   evidence = EXCLUDED.evidence, decision = EXCLUDED.decision, winner = EXCLUDED.winner,
		   decision_at = EXCLUDED.decision_at, dispute_fee = EXCLUDED.dispute_fee,
		   arbitration_fee = EXCLUDED.arbitration_fee, note = EXCLUDED.note`,
		d.ID, d.RefID, d.MerchantID, d.Cardholder, d.AmountMinor, d.Currency, d.ReasonCode, d.Category,
		d.Stage, d.Status, d.FiledAt, d.ResponseDue, d.RespondedAt, d.EscalatedAt,
		join(d.Evidence), d.Decision, d.Winner, d.DecisionAt, d.DisputeFee, d.ArbitrationFee, d.Note)
	return err
}

// Dispute implements Store.
func (p *PostgresStore) Dispute(ctx context.Context, id string) (Dispute, bool, error) {
	var d Dispute
	var evidence string
	err := p.Pool.QueryRow(ctx,
		`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency, reason_code, category,
		        stage, status, filed_at, response_due, responded_at, escalated_at, evidence,
		        decision, winner, decision_at, dispute_fee, arbitration_fee, note
		 FROM disputes WHERE id = $1`, id).
		Scan(&d.ID, &d.RefID, &d.MerchantID, &d.Cardholder, &d.AmountMinor, &d.Currency, &d.ReasonCode,
			&d.Category, &d.Stage, &d.Status, &d.FiledAt, &d.ResponseDue, &d.RespondedAt, &d.EscalatedAt,
			&evidence, &d.Decision, &d.Winner, &d.DecisionAt, &d.DisputeFee, &d.ArbitrationFee, &d.Note)
	if err == pgx.ErrNoRows {
		return d, false, nil
	}
	if err != nil {
		return d, false, err
	}
	d.Evidence = split(evidence)
	return d, true, nil
}

// Disputes implements Store.
func (p *PostgresStore) Disputes(ctx context.Context, merchantID string) ([]Dispute, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency, reason_code, category,
		        stage, status, filed_at, response_due, responded_at, escalated_at, evidence,
		        decision, winner, decision_at, dispute_fee, arbitration_fee, note
		 FROM disputes WHERE merchant_id = $1`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDisputes(rows)
}

// OpenDisputes implements Store.
func (p *PostgresStore) OpenDisputes(ctx context.Context) ([]Dispute, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency, reason_code, category,
		        stage, status, filed_at, response_due, responded_at, escalated_at, evidence,
		        decision, winner, decision_at, dispute_fee, arbitration_fee, note
		 FROM disputes WHERE stage <> 'resolved' AND stage <> 'invalid'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDisputes(rows)
}

// SaveTransaction implements Store.
func (p *PostgresStore) SaveTransaction(ctx context.Context, tx MonitoredTransaction) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO dispute_transactions (ref_id, merchant_id, amount_minor, currency, is_credit, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (ref_id) DO NOTHING`,
		tx.RefID, tx.MerchantID, tx.AmountMinor, tx.Currency, tx.IsCredit, tx.CreatedAt)
	return err
}

// Transaction implements Store.
func (p *PostgresStore) Transaction(ctx context.Context, refID string) (MonitoredTransaction, bool, error) {
	var tx MonitoredTransaction
	err := p.Pool.QueryRow(ctx,
		`SELECT ref_id, merchant_id, amount_minor, currency, is_credit, created_at
		 FROM dispute_transactions WHERE ref_id = $1`, refID).
		Scan(&tx.RefID, &tx.MerchantID, &tx.AmountMinor, &tx.Currency, &tx.IsCredit, &tx.CreatedAt)
	if err == pgx.ErrNoRows {
		return tx, false, nil
	}
	if err != nil {
		return tx, false, err
	}
	return tx, true, nil
}

// Transactions implements Store.
func (p *PostgresStore) Transactions(ctx context.Context, merchantID string) ([]MonitoredTransaction, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT ref_id, merchant_id, amount_minor, currency, is_credit, created_at
		 FROM dispute_transactions WHERE merchant_id = $1`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []MonitoredTransaction
	for rows.Next() {
		var tx MonitoredTransaction
		if err := rows.Scan(&tx.RefID, &tx.MerchantID, &tx.AmountMinor, &tx.Currency, &tx.IsCredit, &tx.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, tx)
	}
	return out, rows.Err()
}

func scanDisputes(rows pgx.Rows) ([]Dispute, error) {
	var out []Dispute
	for rows.Next() {
		var d Dispute
		var evidence string
		if err := rows.Scan(&d.ID, &d.RefID, &d.MerchantID, &d.Cardholder, &d.AmountMinor, &d.Currency,
			&d.ReasonCode, &d.Category, &d.Stage, &d.Status, &d.FiledAt, &d.ResponseDue, &d.RespondedAt,
			&d.EscalatedAt, &evidence, &d.Decision, &d.Winner, &d.DecisionAt, &d.DisputeFee,
			&d.ArbitrationFee, &d.Note); err != nil {
			return nil, err
		}
		d.Evidence = split(evidence)
		out = append(out, d)
	}
	return out, rows.Err()
}

func join(parts []string) string {
	return strings.Join(parts, ",")
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
