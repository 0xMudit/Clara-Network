package ledger

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore persists the ledger in PostgreSQL. AppendEntries writes the
// whole journal in a single transaction, so a posting is either fully present
// or absent.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

// AppendEntries implements Store.
func (p *PostgresStore) AppendEntries(ctx context.Context, entries []Entry) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, e := range entries {
		if _, err := tx.Exec(ctx,
			`INSERT INTO ledger_entries
			 (journal_id, account_id, direction, amount, currency, reference, posted_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			e.JournalID, e.AccountID, string(e.Direction), e.Amount, e.Currency, e.Reference, e.PostedAt); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// EntriesByAccount implements Store.
func (p *PostgresStore) EntriesByAccount(ctx context.Context, accountID string) ([]Entry, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT journal_id, account_id, direction, amount, currency, reference, posted_at
		 FROM ledger_entries WHERE account_id = $1 ORDER BY id`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEntries(rows)
}

// EntriesByAccountAndReference implements Store.
func (p *PostgresStore) EntriesByAccountAndReference(ctx context.Context, accountID, reference string) ([]Entry, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT journal_id, account_id, direction, amount, currency, reference, posted_at
		 FROM ledger_entries WHERE account_id = $1 AND reference = $2 ORDER BY id`,
		accountID, reference)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEntries(rows)
}

// ReferencesWithPrefix implements Store.
func (p *PostgresStore) ReferencesWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT DISTINCT reference FROM ledger_entries
		 WHERE reference LIKE $1 ORDER BY reference`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var ref string
		if err := rows.Scan(&ref); err != nil {
			return nil, err
		}
		out = append(out, ref)
	}
	return out, rows.Err()
}

// Account implements Store.
func (p *PostgresStore) Account(ctx context.Context, id string) (Account, bool, error) {
	var a Account
	err := p.Pool.QueryRow(ctx,
		`SELECT id, type, currency FROM ledger_accounts WHERE id = $1`, id).
		Scan(&a.ID, &a.Type, &a.Currency)
	if err == pgx.ErrNoRows {
		return a, false, nil
	}
	if err != nil {
		return a, false, err
	}
	return a, true, nil
}

// SaveAccount implements Store.
func (p *PostgresStore) SaveAccount(ctx context.Context, acct Account) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO ledger_accounts (id, type, currency) VALUES ($1,$2,$3)
		 ON CONFLICT (id) DO NOTHING`,
		acct.ID, string(acct.Type), acct.Currency)
	return err
}

// Accounts implements Store.
func (p *PostgresStore) Accounts(ctx context.Context) ([]Account, error) {
	rows, err := p.Pool.Query(ctx, `SELECT id, type, currency FROM ledger_accounts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Type, &a.Currency); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Close implements Store.
func (p *PostgresStore) Close() error {
	p.Pool.Close()
	return nil
}

func scanEntries(rows pgx.Rows) ([]Entry, error) {
	var out []Entry
	for rows.Next() {
		var e Entry
		var dir string
		var posted time.Time
		if err := rows.Scan(&e.JournalID, &e.AccountID, &dir, &e.Amount,
			&e.Currency, &e.Reference, &posted); err != nil {
			return nil, err
		}
		e.Direction = Direction(dir)
		e.PostedAt = posted
		out = append(out, e)
	}
	return out, rows.Err()
}
