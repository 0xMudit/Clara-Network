// Package risk implements the in-path authorization risk engine: velocity
// counters and amount thresholds that can decline a transaction before it is
// forwarded to the issuer. It must stay well under the 100 ms latency budget.
package risk

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store counts events per key within a window. RedisStore backs the hot path;
// MemoryStore is used in tests and single-instance deployments.
type Store interface {
	// Incr increments the counter for key (applying ttl on each call) and
	// returns the new count.
	Incr(ctx context.Context, key string, ttl time.Duration) (int64, error)
	Close() error
}

// Built-in rule kinds.
const (
	KindCardVelocity     = "velocity-card"
	KindMerchantVelocity = "velocity-merchant"
	KindAmountLimit      = "amount"
)

// Rule is a single risk rule.
type Rule struct {
	Name    string        `json:"name"`
	Kind    string        `json:"kind"`
	Limit   int64         `json:"limit"`
	Window  time.Duration `json:"window"` // seconds for velocity rules
	Code    string        `json:"code"`   // decline response code
	Enabled bool          `json:"enabled"`
}

// Decision is the result of risk evaluation.
type Decision struct {
	Allow  bool
	Code   string
	Reason string
}

// Engine evaluates transactions against a rule set.
type Engine struct {
	store Store
	rules []Rule
}

// New builds a risk engine.
func New(store Store, rules []Rule) *Engine {
	return &Engine{store: store, rules: rules}
}

// FromConfig builds an engine from a JSON config document and a store.
func FromConfig(data []byte, store Store) (*Engine, error) {
	var cfg struct {
		Rules []Rule `json:"rules"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return New(store, cfg.Rules), nil
}

// Evaluate runs the rules for a transaction. It always allows when the engine
// has no store (not configured).
func (e *Engine) Evaluate(ctx context.Context, pan, merchantID string, amountMinor int64) (*Decision, error) {
	if e == nil || e.store == nil {
		return &Decision{Allow: true}, nil
	}
	for _, r := range e.rules {
		if !r.Enabled {
			continue
		}
		var key string
		switch r.Kind {
		case KindCardVelocity:
			key = "risk:vel:card:" + pan
		case KindMerchantVelocity:
			key = "risk:vel:mch:" + merchantID
		case KindAmountLimit:
			if amountMinor > r.Limit {
				return &Decision{Allow: false, Code: r.Code, Reason: r.Name}, nil
			}
			continue
		default:
			continue
		}
		n, err := e.store.Incr(ctx, key, r.Window*time.Second)
		if err != nil {
			return nil, err
		}
		if n > r.Limit {
			return &Decision{Allow: false, Code: r.Code, Reason: r.Name}, nil
		}
	}
	return &Decision{Allow: true}, nil
}

// MemoryStore is an in-process Store for tests.
type MemoryStore struct {
	mu   sync.Mutex
	vals map[string]int64
	ttls map[string]time.Time
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{vals: map[string]int64{}, ttls: map[string]time.Time{}}
}

// Incr increments key, resetting its window when expired.
func (m *MemoryStore) Incr(_ context.Context, key string, ttl time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if exp, ok := m.ttls[key]; !ok || now.After(exp) {
		m.vals[key] = 0
		m.ttls[key] = now.Add(ttl)
	}
	m.vals[key]++
	return m.vals[key], nil
}

// Close implements Store.
func (m *MemoryStore) Close() error { return nil }

// RedisStore implements Store on top of Redis INCR + EXPIRE.
type RedisStore struct {
	rdb *redis.Client
}

// NewRedisStore connects to the Redis instance at addr.
func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{rdb: redis.NewClient(&redis.Options{Addr: addr})}
}

// Incr increments key and (re)applies the TTL in a single pipeline round-trip.
func (r *RedisStore) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := r.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// Close implements Store.
func (r *RedisStore) Close() error { return r.rdb.Close() }
