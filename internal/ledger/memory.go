package ledger

import (
	"context"
	"sort"
	"sync"
)

// MemoryStore is an in-process Store for tests and single-instance runs.
type MemoryStore struct {
	mu       sync.Mutex
	entries  []Entry
	accounts map[string]Account
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{accounts: map[string]Account{}}
}

// AppendEntries implements Store.
func (m *MemoryStore) AppendEntries(_ context.Context, entries []Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entries...)
	return nil
}

// EntriesByAccount implements Store.
func (m *MemoryStore) EntriesByAccount(_ context.Context, accountID string) ([]Entry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Entry
	for _, e := range m.entries {
		if e.AccountID == accountID {
			out = append(out, e)
		}
	}
	return out, nil
}

// EntriesByAccountAndReference implements Store.
func (m *MemoryStore) EntriesByAccountAndReference(_ context.Context, accountID, reference string) ([]Entry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Entry
	for _, e := range m.entries {
		if e.AccountID == accountID && e.Reference == reference {
			out = append(out, e)
		}
	}
	return out, nil
}

// ReferencesWithPrefix implements Store.
func (m *MemoryStore) ReferencesWithPrefix(_ context.Context, prefix string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	seen := map[string]bool{}
	for _, e := range m.entries {
		if len(e.Reference) >= len(prefix) && e.Reference[:len(prefix)] == prefix {
			seen[e.Reference] = true
		}
	}
	out := make([]string, 0, len(seen))
	for ref := range seen {
		out = append(out, ref)
	}
	sort.Strings(out)
	return out, nil
}

// Account implements Store.
func (m *MemoryStore) Account(_ context.Context, id string) (Account, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accounts[id]
	return a, ok, nil
}

// SaveAccount implements Store.
func (m *MemoryStore) SaveAccount(_ context.Context, acct Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts[acct.ID] = acct
	return nil
}

// Accounts implements Store.
func (m *MemoryStore) Accounts(_ context.Context) ([]Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Account, 0, len(m.accounts))
	for _, a := range m.accounts {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Close implements Store.
func (m *MemoryStore) Close() error { return nil }
