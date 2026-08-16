package switchsrv

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// IdempotencyStore caches response frames keyed by a request identity so that
// replays return the original response instead of re-processing.
type IdempotencyStore interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type memEntry struct {
	value   string
	expires time.Time
}

// MemoryIdempotency is an in-process store used in tests and when Redis is
// not configured.
type MemoryIdempotency struct {
	mu sync.Mutex
	m  map[string]memEntry
}

// NewMemoryIdempotency returns an empty in-memory store.
func NewMemoryIdempotency() *MemoryIdempotency {
	return &MemoryIdempotency{m: make(map[string]memEntry)}
}

func (s *MemoryIdempotency) Get(_ context.Context, key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[key]
	if !ok {
		return "", false, nil
	}
	if time.Now().After(e.expires) {
		delete(s.m, key)
		return "", false, nil
	}
	return e.value, true, nil
}

func (s *MemoryIdempotency) Set(_ context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = memEntry{value: value, expires: time.Now().Add(ttl)}
	return nil
}

// RedisIdempotency uses Redis for a shared, durable idempotency cache.
type RedisIdempotency struct {
	Client *redis.Client
}

func (r *RedisIdempotency) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (r *RedisIdempotency) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}
