package clearing

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store persists clearing records, net positions, prefund accounts, and the
// default fund.
type Store interface {
	AppendRecords(ctx context.Context, records []ClearingRecord) error
	Records(ctx context.Context, cycleID string) ([]ClearingRecord, error)
	SavePositions(ctx context.Context, positions []NetPosition) error
	SaveInstructions(ctx context.Context, instructions []SettlementInstruction) error
	PrefundAccount(ctx context.Context, member string) (PrefundAccount, bool, error)
	SetPrefund(ctx context.Context, acct PrefundAccount) error
	DefaultFundBalance(ctx context.Context) (int64, error)
	SetDefaultFundBalance(ctx context.Context, amount int64) error
}

// MemoryStore is an in-process Store for tests and single-instance runs.
type MemoryStore struct {
	mu           sync.Mutex
	records      []ClearingRecord
	positions    []NetPosition
	instructions []SettlementInstruction
	prefunds     map[string]PrefundAccount
	dfBalance    int64
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{prefunds: map[string]PrefundAccount{}}
}

// AppendRecords implements Store.
func (m *MemoryStore) AppendRecords(_ context.Context, records []ClearingRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, records...)
	return nil
}

// Records implements Store.
func (m *MemoryStore) Records(_ context.Context, cycleID string) ([]ClearingRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []ClearingRecord
	for _, r := range m.records {
		if r.CycleID == cycleID {
			out = append(out, r)
		}
	}
	return out, nil
}

// SavePositions implements Store.
func (m *MemoryStore) SavePositions(_ context.Context, positions []NetPosition) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.positions = append(m.positions, positions...)
	return nil
}

// SaveInstructions implements Store.
func (m *MemoryStore) SaveInstructions(_ context.Context, instructions []SettlementInstruction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.instructions = append(m.instructions, instructions...)
	return nil
}

// PrefundAccount implements Store.
func (m *MemoryStore) PrefundAccount(_ context.Context, member string) (PrefundAccount, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.prefunds[member]
	return a, ok, nil
}

// SetPrefund implements Store.
func (m *MemoryStore) SetPrefund(_ context.Context, acct PrefundAccount) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prefunds[acct.Member] = acct
	return nil
}

// DefaultFundBalance implements Store.
func (m *MemoryStore) DefaultFundBalance(_ context.Context) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dfBalance, nil
}

// SetDefaultFundBalance implements Store.
func (m *MemoryStore) SetDefaultFundBalance(_ context.Context, amount int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dfBalance = amount
	return nil
}

// PostgresStore persists the clearing layer in PostgreSQL.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// AppendRecords implements Store.
func (p *PostgresStore) AppendRecords(ctx context.Context, records []ClearingRecord) error {
	for _, r := range records {
		if _, err := p.Pool.Exec(ctx,
			`INSERT INTO clearing_records
			 (cycle_id, stan, mti, sender, receiver, amount_minor, interchange, currency, ref_id)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			r.CycleID, r.STAN, r.MTI, r.Sender, r.Receiver, r.AmountMinor, r.Interchange, r.Currency, r.RefID); err != nil {
			return err
		}
	}
	return nil
}

// Records implements Store.
func (p *PostgresStore) Records(ctx context.Context, cycleID string) ([]ClearingRecord, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT cycle_id, stan, mti, sender, receiver, amount_minor, interchange, currency, ref_id
		 FROM clearing_records WHERE cycle_id = $1 ORDER BY id`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ClearingRecord
	for rows.Next() {
		var r ClearingRecord
		if err := rows.Scan(&r.CycleID, &r.STAN, &r.MTI, &r.Sender, &r.Receiver,
			&r.AmountMinor, &r.Interchange, &r.Currency, &r.RefID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// SavePositions implements Store.
func (p *PostgresStore) SavePositions(ctx context.Context, positions []NetPosition) error {
	for _, pos := range positions {
		if _, err := p.Pool.Exec(ctx,
			`INSERT INTO net_positions (cycle_id, member, net_minor) VALUES ($1,$2,$3)`,
			pos.CycleID, pos.Member, pos.Net); err != nil {
			return err
		}
	}
	return nil
}

// SaveInstructions implements Store.
func (p *PostgresStore) SaveInstructions(ctx context.Context, instructions []SettlementInstruction) error {
	for _, in := range instructions {
		if _, err := p.Pool.Exec(ctx,
			`INSERT INTO settlement_instructions
			 (cycle_id, msg_id, member, amount_minor, direction, currency, instruction_time, final)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			in.CycleID, in.MsgID, in.Member, in.Amount, in.Direction, in.Currency, in.Instruction, in.Final); err != nil {
			return err
		}
	}
	return nil
}

// PrefundAccount implements Store.
func (p *PostgresStore) PrefundAccount(ctx context.Context, member string) (PrefundAccount, bool, error) {
	var a PrefundAccount
	err := p.Pool.QueryRow(ctx,
		`SELECT member, balance, cap FROM prefund_accounts WHERE member = $1`, member).
		Scan(&a.Member, &a.Balance, &a.Cap)
	if err == pgx.ErrNoRows {
		return a, false, nil
	}
	if err != nil {
		return a, false, err
	}
	return a, true, nil
}

// SetPrefund implements Store.
func (p *PostgresStore) SetPrefund(ctx context.Context, acct PrefundAccount) error {
	if _, err := p.Pool.Exec(ctx,
		`INSERT INTO prefund_accounts (member, balance, cap) VALUES ($1,$2,$3)
		 ON CONFLICT (member) DO UPDATE SET balance = EXCLUDED.balance, cap = EXCLUDED.cap`,
		acct.Member, acct.Balance, acct.Cap); err != nil {
		return err
	}
	return nil
}

// DefaultFundBalance implements Store.
func (p *PostgresStore) DefaultFundBalance(ctx context.Context) (int64, error) {
	var bal int64
	err := p.Pool.QueryRow(ctx,
		`SELECT balance FROM default_fund WHERE id = 1`).Scan(&bal)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return bal, err
}

// SetDefaultFundBalance implements Store.
func (p *PostgresStore) SetDefaultFundBalance(ctx context.Context, amount int64) error {
	if _, err := p.Pool.Exec(ctx,
		`INSERT INTO default_fund (id, balance) VALUES (1, $1)
		 ON CONFLICT (id) DO UPDATE SET balance = EXCLUDED.balance`, amount); err != nil {
		return err
	}
	return nil
}
