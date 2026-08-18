package adminapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// LedgerAccount is an account with its recomputed balance.
type LedgerAccount struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Balance  int64  `json:"balance"`
	Currency string `json:"-"`
}

// LedgerAccountsWithBalances returns every ledger account with its
// recomputed balance (credits minus debits) in a single query.
func (s *Store) LedgerAccountsWithBalances(ctx context.Context) ([]LedgerAccount, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT a.id, a.type, a.currency,
		        COALESCE(SUM(CASE WHEN e.direction='credit' THEN e.amount ELSE 0 END)
		               - SUM(CASE WHEN e.direction='debit'  THEN e.amount ELSE 0 END), 0)
		 FROM ledger_accounts a
		 LEFT JOIN ledger_entries e ON e.account_id = a.id
		 GROUP BY a.id, a.type, a.currency
		 ORDER BY a.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LedgerAccount
	for rows.Next() {
		var a LedgerAccount
		if err := rows.Scan(&a.ID, &a.Type, &a.Currency, &a.Balance); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// LedgerEntry is a single journal line.
type LedgerEntry struct {
	JournalID string    `json:"journalId"`
	AccountID string    `json:"accountId"`
	Direction string    `json:"direction"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Reference string    `json:"reference"`
	PostedAt  time.Time `json:"postedAt"`
}

// LedgerEntries returns journal entries, optionally filtered by account, with
// pagination. When accountID is empty all entries are returned.
func (s *Store) LedgerEntries(ctx context.Context, accountID string, limit, offset int) (Page[LedgerEntry], error) {
	var total int64
	var err error

	if accountID != "" {
		err = s.Pool.QueryRow(ctx,
			`SELECT count(*) FROM ledger_entries WHERE account_id = $1`, accountID).Scan(&total)
	} else {
		err = s.Pool.QueryRow(ctx,
			`SELECT count(*) FROM ledger_entries`).Scan(&total)
	}
	if err != nil {
		return Page[LedgerEntry]{}, err
	}

	var rows pgx.Rows
	if accountID != "" {
		rows, err = s.Pool.Query(ctx,
			`SELECT journal_id, account_id, direction, amount, currency, reference, posted_at
			 FROM ledger_entries WHERE account_id = $1
			 ORDER BY posted_at DESC LIMIT $2 OFFSET $3`, accountID, limit, offset)
	} else {
		rows, err = s.Pool.Query(ctx,
			`SELECT journal_id, account_id, direction, amount, currency, reference, posted_at
			 FROM ledger_entries ORDER BY posted_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return Page[LedgerEntry]{}, err
	}
	defer rows.Close()
	var out []LedgerEntry
	for rows.Next() {
		var e LedgerEntry
		if err := rows.Scan(&e.JournalID, &e.AccountID, &e.Direction, &e.Amount,
			&e.Currency, &e.Reference, &e.PostedAt); err != nil {
			return Page[LedgerEntry]{}, err
		}
		out = append(out, e)
	}
	return Page[LedgerEntry]{Items: out, Total: total}, rows.Err()
}
