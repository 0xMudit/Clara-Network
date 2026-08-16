package acquiring

import (
	"context"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store persists merchants, funding lines, and the negative screening lists.
type Store interface {
	SaveMerchant(ctx context.Context, m Merchant) error
	Merchant(ctx context.Context, id string) (Merchant, bool, error)
	AppendFunding(ctx context.Context, line FundingLine) error
	Funding(ctx context.Context, merchantID string) ([]FundingLine, error)
	SaveMatchEntries(ctx context.Context, entries []MatchEntry) error
	MatchEntries(ctx context.Context) ([]MatchEntry, error)
	SaveOfacEntries(ctx context.Context, entries []OfacEntry) error
	OfacEntries(ctx context.Context) ([]OfacEntry, error)
}

// MemoryStore is an in-process Store for tests and single-instance runs.
type MemoryStore struct {
	mu        sync.Mutex
	merchants map[string]Merchant
	funding   []FundingLine
	matches   []MatchEntry
	ofac      []OfacEntry
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{merchants: map[string]Merchant{}}
}

// SaveMerchant implements Store.
func (m *MemoryStore) SaveMerchant(_ context.Context, merchant Merchant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.merchants[merchant.ID] = merchant
	return nil
}

// Merchant implements Store.
func (m *MemoryStore) Merchant(_ context.Context, id string) (Merchant, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	merchant, ok := m.merchants[id]
	return merchant, ok, nil
}

// AppendFunding implements Store.
func (m *MemoryStore) AppendFunding(_ context.Context, line FundingLine) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.funding = append(m.funding, line)
	return nil
}

// Funding implements Store.
func (m *MemoryStore) Funding(_ context.Context, merchantID string) ([]FundingLine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []FundingLine
	for _, l := range m.funding {
		if l.MerchantID == merchantID {
			out = append(out, l)
		}
	}
	return out, nil
}

// SaveMatchEntries implements Store.
func (m *MemoryStore) SaveMatchEntries(_ context.Context, entries []MatchEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.matches = append(m.matches, entries...)
	return nil
}

// MatchEntries implements Store.
func (m *MemoryStore) MatchEntries(_ context.Context) ([]MatchEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]MatchEntry, len(m.matches))
	copy(out, m.matches)
	return out, nil
}

// SaveOfacEntries implements Store.
func (m *MemoryStore) SaveOfacEntries(_ context.Context, entries []OfacEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ofac = append(m.ofac, entries...)
	return nil
}

// OfacEntries implements Store.
func (m *MemoryStore) OfacEntries(_ context.Context) ([]OfacEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]OfacEntry, len(m.ofac))
	copy(out, m.ofac)
	return out, nil
}

// PostgresStore persists the acquiring stack in PostgreSQL.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// SaveMerchant implements Store.
func (p *PostgresStore) SaveMerchant(ctx context.Context, m Merchant) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO merchants
		 (id, name, dba, tax_id, principals, mccs, status, risk_tier, reserve_rate_bps,
		  funding_delay_days, transaction_limit, reserve_balance, volume, decline_reason, approved_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status,
		   reserve_balance = EXCLUDED.reserve_balance, decline_reason = EXCLUDED.decline_reason`,
		m.ID, m.Name, m.DBA, m.TaxID, join(m.Principals), join(m.MCCs), m.Status,
		m.RiskTier, m.ReserveRateBPS, m.FundingDelayDays, m.TransactionLimit,
		m.ReserveBalance, m.Volume, m.DeclineReason, m.ApprovedAt)
	return err
}

// Merchant implements Store.
func (p *PostgresStore) Merchant(ctx context.Context, id string) (Merchant, bool, error) {
	var m Merchant
	var principals, mccs string
	err := p.Pool.QueryRow(ctx,
		`SELECT id, name, dba, tax_id, principals, mccs, status, risk_tier, reserve_rate_bps,
		        funding_delay_days, transaction_limit, reserve_balance, volume, decline_reason, approved_at
		 FROM merchants WHERE id = $1`, id).
		Scan(&m.ID, &m.Name, &m.DBA, &m.TaxID, &principals, &mccs, &m.Status,
			&m.RiskTier, &m.ReserveRateBPS, &m.FundingDelayDays, &m.TransactionLimit,
			&m.ReserveBalance, &m.Volume, &m.DeclineReason, &m.ApprovedAt)
	if err == pgx.ErrNoRows {
		return m, false, nil
	}
	if err != nil {
		return m, false, err
	}
	m.Principals = split(principals)
	m.MCCs = split(mccs)
	return m, true, nil
}

// AppendFunding implements Store.
func (p *PostgresStore) AppendFunding(ctx context.Context, line FundingLine) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO funding_lines (batch_id, merchant_id, gross, fees, reserve_hold, net, date)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		line.BatchID, line.MerchantID, line.Gross, line.Fees, line.ReserveHold, line.Net, line.Date)
	return err
}

// Funding implements Store.
func (p *PostgresStore) Funding(ctx context.Context, merchantID string) ([]FundingLine, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT batch_id, merchant_id, gross, fees, reserve_hold, net, date
		 FROM funding_lines WHERE merchant_id = $1 ORDER BY date`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FundingLine
	for rows.Next() {
		var l FundingLine
		if err := rows.Scan(&l.BatchID, &l.MerchantID, &l.Gross, &l.Fees, &l.ReserveHold, &l.Net, &l.Date); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// SaveMatchEntries implements Store.
func (p *PostgresStore) SaveMatchEntries(ctx context.Context, entries []MatchEntry) error {
	for _, e := range entries {
		if _, err := p.Pool.Exec(ctx,
			`INSERT INTO screening_lists (list, name, tax_id, detail)
			 VALUES ('MATCH', $1, $2, $3)
			 ON CONFLICT (list, name) DO NOTHING`,
			e.MerchantName, e.TaxID, e.Reason); err != nil {
			return err
		}
	}
	return nil
}

// MatchEntries implements Store.
func (p *PostgresStore) MatchEntries(ctx context.Context) ([]MatchEntry, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT name, tax_id, detail FROM screening_lists WHERE list = 'MATCH'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []MatchEntry
	for rows.Next() {
		var e MatchEntry
		if err := rows.Scan(&e.MerchantName, &e.TaxID, &e.Reason); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// SaveOfacEntries implements Store.
func (p *PostgresStore) SaveOfacEntries(ctx context.Context, entries []OfacEntry) error {
	for _, e := range entries {
		if _, err := p.Pool.Exec(ctx,
			`INSERT INTO screening_lists (list, name, detail)
			 VALUES ('OFAC', $1, $2)
			 ON CONFLICT (list, name) DO NOTHING`,
			e.Name, e.Program); err != nil {
			return err
		}
	}
	return nil
}

// OfacEntries implements Store.
func (p *PostgresStore) OfacEntries(ctx context.Context) ([]OfacEntry, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT name, detail FROM screening_lists WHERE list = 'OFAC'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []OfacEntry
	for rows.Next() {
		var e OfacEntry
		if err := rows.Scan(&e.Name, &e.Program); err != nil {
			return nil, err
		}
		out = append(out, e)
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
