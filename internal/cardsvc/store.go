package cardsvc

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store persists BIN ranges, cards, and tokens.
type Store interface {
	SaveRange(ctx context.Context, r BinRange) error
	Ranges(ctx context.Context) ([]BinRange, error)
	SaveCard(ctx context.Context, c Card) error
	Card(ctx context.Context, ref string) (Card, bool, error)
	SaveToken(ctx context.Context, t Token) error
	Token(ctx context.Context, token string) (Token, bool, error)
	TokenByPANHash(ctx context.Context, panHash []byte) (Token, bool, error)
}

// MemoryStore is an in-process Store for tests and single-instance runs.
type MemoryStore struct {
	mu     sync.Mutex
	ranges []BinRange
	cards  map[string]Card
	tokens map[string]Token
	byPan  map[string]Token
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		cards:  map[string]Card{},
		tokens: map[string]Token{},
		byPan:  map[string]Token{},
	}
}

// SaveRange implements Store.
func (m *MemoryStore) SaveRange(_ context.Context, r BinRange) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ranges = append(m.ranges, r)
	return nil
}

// Ranges implements Store.
func (m *MemoryStore) Ranges(_ context.Context) ([]BinRange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]BinRange, len(m.ranges))
	copy(out, m.ranges)
	return out, nil
}

// SaveCard implements Store.
func (m *MemoryStore) SaveCard(_ context.Context, c Card) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cards[c.Ref] = c
	return nil
}

// Card implements Store.
func (m *MemoryStore) Card(_ context.Context, ref string) (Card, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cards[ref]
	return c, ok, nil
}

// SaveToken implements Store.
func (m *MemoryStore) SaveToken(_ context.Context, t Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[t.Number] = t
	m.byPan[string(t.PANHash)] = t
	return nil
}

// Token implements Store.
func (m *MemoryStore) Token(_ context.Context, token string) (Token, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tokens[token]
	return t, ok, nil
}

// TokenByPANHash implements Store.
func (m *MemoryStore) TokenByPANHash(_ context.Context, panHash []byte) (Token, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.byPan[string(panHash)]
	return t, ok, nil
}

// PostgresStore persists the issuing stack in PostgreSQL.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// SaveRange implements Store.
func (p *PostgresStore) SaveRange(ctx context.Context, r BinRange) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO bin_ranges (bin, low, high, currency, product) VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (bin) DO UPDATE SET low = EXCLUDED.low, high = EXCLUDED.high,
		 currency = EXCLUDED.currency, product = EXCLUDED.product`,
		r.BIN, r.Low, r.High, r.Currency, r.Product)
	return err
}

// Ranges implements Store.
func (p *PostgresStore) Ranges(ctx context.Context) ([]BinRange, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT bin, low, high, currency, product FROM bin_ranges ORDER BY bin`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BinRange
	for rows.Next() {
		var r BinRange
		if err := rows.Scan(&r.BIN, &r.Low, &r.High, &r.Currency, &r.Product); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// SaveCard implements Store.
func (p *PostgresStore) SaveCard(ctx context.Context, c Card) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO cards (ref, pan_hash, pan_masked, bin, expiry, status, product, udk, last_atc)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (ref) DO UPDATE SET status = EXCLUDED.status, last_atc = EXCLUDED.last_atc`,
		c.Ref, c.PANHash, c.PANMask, c.BIN, c.Expiry, c.Status, c.Product, c.UDK, c.LastATC)
	return err
}

// Card implements Store.
func (p *PostgresStore) Card(ctx context.Context, ref string) (Card, bool, error) {
	var c Card
	err := p.Pool.QueryRow(ctx,
		`SELECT ref, pan_hash, pan_masked, bin, expiry, status, product, udk, last_atc
		 FROM cards WHERE ref = $1`, ref).
		Scan(&c.Ref, &c.PANHash, &c.PANMask, &c.BIN, &c.Expiry, &c.Status, &c.Product, &c.UDK, &c.LastATC)
	if err == pgx.ErrNoRows {
		return c, false, nil
	}
	if err != nil {
		return c, false, err
	}
	return c, true, nil
}

// SaveToken implements Store.
func (p *PostgresStore) SaveToken(ctx context.Context, t Token) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO tokens (token, pan_hash, par, status, bin, trid, device_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT (token) DO UPDATE SET status = EXCLUDED.status,
		 trid = EXCLUDED.trid, device_id = EXCLUDED.device_id`,
		t.Number, t.PANHash, t.PAR, t.Status, t.BIN, t.Requestor, t.DeviceID, t.CreatedAt)
	return err
}

// Token implements Store.
func (p *PostgresStore) Token(ctx context.Context, token string) (Token, bool, error) {
	var t Token
	err := p.Pool.QueryRow(ctx,
		`SELECT token, pan_hash, par, status, bin, trid, device_id, created_at
		 FROM tokens WHERE token = $1`, token).
		Scan(&t.Number, &t.PANHash, &t.PAR, &t.Status, &t.BIN, &t.Requestor, &t.DeviceID, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return t, false, nil
	}
	if err != nil {
		return t, false, err
	}
	return t, true, nil
}

// TokenByPANHash implements Store.
func (p *PostgresStore) TokenByPANHash(ctx context.Context, panHash []byte) (Token, bool, error) {
	var t Token
	err := p.Pool.QueryRow(ctx,
		`SELECT token, pan_hash, par, status, bin, trid, device_id, created_at
		 FROM tokens WHERE pan_hash = $1 ORDER BY created_at DESC LIMIT 1`, panHash).
		Scan(&t.Number, &t.PANHash, &t.PAR, &t.Status, &t.BIN, &t.Requestor, &t.DeviceID, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return t, false, nil
	}
	if err != nil {
		return t, false, err
	}
	return t, true, nil
}
