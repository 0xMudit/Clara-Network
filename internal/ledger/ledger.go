// Package ledger implements an append-only, double-entry general ledger with
// a balance invariant enforced on every write (see docs/12 §12.5). Amounts
// are integer minor units; balance is always recomputed from the journal —
// the ledger is the source of truth.
package ledger

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Direction of a single posting line.
type Direction string

const (
	// Debit increases asset and expense accounts, decreases liability and
	// income accounts.
	Debit Direction = "debit"
	// Credit increases liability and income accounts, decreases asset and
	// expense accounts.
	Credit Direction = "credit"
)

// AccountType classifies an account for reporting.
type AccountType string

const (
	AccountAsset      AccountType = "asset"
	AccountLiability  AccountType = "liability"
	AccountIncome     AccountType = "income"
	AccountExpense    AccountType = "expense"
)

// Account is a named container of value within the scheme's books.
type Account struct {
	ID       string
	Type     AccountType
	Currency string
}

// Entry is one posted line in the journal. The whole journal posts in a
// single atomic write; corrections are made with reversing entries, never by
// editing or deleting an entry.
type Entry struct {
	JournalID string // grouping key for the posting (e.g. cycle:member)
	AccountID string
	Direction Direction
	Amount    int64 // minor units, always positive
	Currency  string
	Reference string // external reference (e.g. "cycle:member")
	PostedAt  time.Time
}

// Store persists accounts and the append-only journal.
type Store interface {
	AppendEntries(ctx context.Context, entries []Entry) error
	EntriesByAccount(ctx context.Context, accountID string) ([]Entry, error)
	EntriesByAccountAndReference(ctx context.Context, accountID, reference string) ([]Entry, error)
	ReferencesWithPrefix(ctx context.Context, prefix string) ([]string, error)
	Account(ctx context.Context, id string) (Account, bool, error)
	SaveAccount(ctx context.Context, acct Account) error
	Accounts(ctx context.Context) ([]Account, error)
	Close() error
}

// Ledger posts balanced journal entries to a Store.
type Ledger struct {
	store Store
	log   *slog.Logger
}

// NewLedger returns a ledger over the given store.
func NewLedger(store Store) *Ledger {
	return &Ledger{store: store, log: slog.Default()}
}

// EnsureAccount creates the account if it does not already exist.
func (l *Ledger) EnsureAccount(ctx context.Context, acct Account) error {
	if acct.ID == "" {
		return fmt.Errorf("ledger: empty account id")
	}
	if acct.Currency == "" {
		return fmt.Errorf("ledger: account %s has no currency", acct.ID)
	}
	if _, ok, err := l.store.Account(ctx, acct.ID); err != nil {
		return fmt.Errorf("ledger: load account %s: %w", acct.ID, err)
	} else if ok {
		return nil
	}
	return l.store.SaveAccount(ctx, acct)
}

// Post validates a balanced double-entry posting and writes it atomically.
// The debit and credit totals must be equal within each currency group, every
// amount must be positive, and every referenced account must exist. Returns
// an error and writes nothing if any check fails.
func (l *Ledger) Post(ctx context.Context, journalID string, entries []Entry) error {
	if journalID == "" {
		return fmt.Errorf("ledger: empty journal id")
	}
	if len(entries) < 2 {
		return fmt.Errorf("ledger: journal %s needs at least two entries", journalID)
	}

	// Per-currency balance: sum(debits) == sum(credits).
	balances := map[string]int64{}
	for i, e := range entries {
		if e.AccountID == "" {
			return fmt.Errorf("ledger: journal %s entry %d has no account", journalID, i)
		}
		if e.Amount <= 0 {
			return fmt.Errorf("ledger: journal %s entry %d has non-positive amount %d", journalID, i, e.Amount)
		}
		if e.Currency == "" {
			return fmt.Errorf("ledger: journal %s entry %d has no currency", journalID, i)
		}
		if _, ok, err := l.store.Account(ctx, e.AccountID); err != nil {
			return fmt.Errorf("ledger: journal %s: load account %s: %w", journalID, e.AccountID, err)
		} else if !ok {
			return fmt.Errorf("ledger: journal %s references unknown account %s", journalID, e.AccountID)
		}
		switch e.Direction {
		case Debit:
			balances[e.Currency] += e.Amount
		case Credit:
			balances[e.Currency] -= e.Amount
		default:
			return fmt.Errorf("ledger: journal %s entry %d has invalid direction %q", journalID, i, e.Direction)
		}
	}
	for currency, net := range balances {
		if net != 0 {
			return fmt.Errorf("ledger: journal %s does not balance in %s (debit-credit net %d)", journalID, currency, net)
		}
	}

	now := time.Now().UTC()
	for i := range entries {
		entries[i].JournalID = journalID
		entries[i].PostedAt = now
		if entries[i].Reference == "" {
			entries[i].Reference = journalID
		}
	}
	if err := l.store.AppendEntries(ctx, entries); err != nil {
		return fmt.Errorf("ledger: persist journal %s: %w", journalID, err)
	}
	l.log.Debug("journal posted", "journal_id", journalID, "entries", len(entries))
	return nil
}

// Balance recomputes an account's balance from the journal: credits minus
// debits. The ledger is the truth; no balance column is trusted.
func (l *Ledger) Balance(ctx context.Context, accountID string) (int64, error) {
	entries, err := l.store.EntriesByAccount(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("ledger: load entries for %s: %w", accountID, err)
	}
	var balance int64
	for _, e := range entries {
		switch e.Direction {
		case Credit:
			balance += e.Amount
		case Debit:
			balance -= e.Amount
		}
	}
	return balance, nil
}

// TotalBalance returns the signed sum of all accounts; a balanced ledger
// always sums to zero.
func (l *Ledger) TotalBalance(ctx context.Context) (int64, error) {
	accounts, err := l.store.Accounts(ctx)
	if err != nil {
		return 0, fmt.Errorf("ledger: list accounts: %w", err)
	}
	var total int64
	for _, a := range accounts {
		b, err := l.Balance(ctx, a.ID)
		if err != nil {
			return 0, err
		}
		total += b
	}
	return total, nil
}

// BalanceLine is an account with its recomputed balance.
type BalanceLine struct {
	AccountID string
	Type      AccountType
	Balance   int64
}

// Balances returns every account with its recomputed balance, sorted by ID.
func (l *Ledger) Balances(ctx context.Context) ([]BalanceLine, error) {
	accounts, err := l.store.Accounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ledger: list accounts: %w", err)
	}
	out := make([]BalanceLine, 0, len(accounts))
	for _, a := range accounts {
		b, err := l.Balance(ctx, a.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, BalanceLine{AccountID: a.ID, Type: a.Type, Balance: b})
	}
	return out, nil
}
