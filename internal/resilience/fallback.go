package resilience

import (
	"sort"
	"sync"
	"time"
)

// RouteState is the circuit-breaker state of an issuer route.
type RouteState string

// Circuit-breaker states.
const (
	RouteClosed   RouteState = "closed"
	RouteOpen     RouteState = "open"
	RouteHalfOpen RouteState = "half-open"
)

// Route is the observed health of a single issuer address.
type Route struct {
	Addr        string
	State       RouteState
	Failures    int
	Threshold   int
	Cooldown    time.Duration
	OpenedAt    time.Time
	LastAttempt time.Time
	LastSuccess time.Time
}

// RouteHealth tracks consecutive failures per issuer address and trips a
// circuit breaker so traffic skips a dead primary and fails over to the
// secondary (docs/19 §19.3 fallback ordering). An open circuit admits a
// half-open probe once the cooldown has elapsed; a probe success re-closes
// the circuit, a probe failure re-opens it.
type RouteHealth struct {
	mu               sync.Mutex
	routes           map[string]*Route
	defaultThreshold int
	defaultCooldown  time.Duration
}

// NewRouteHealth builds a circuit-breaker table. A threshold <= 0 defaults to
// 3 consecutive failures; a cooldown <= 0 defaults to 10s.
func NewRouteHealth(threshold int, cooldown time.Duration) *RouteHealth {
	if threshold <= 0 {
		threshold = 3
	}
	if cooldown <= 0 {
		cooldown = 10 * time.Second
	}
	return &RouteHealth{
		routes:           map[string]*Route{},
		defaultThreshold: threshold,
		defaultCooldown:  cooldown,
	}
}

func (h *RouteHealth) routeLocked(addr string) *Route {
	r, ok := h.routes[addr]
	if !ok {
		r = &Route{
			Addr:      addr,
			State:     RouteClosed,
			Threshold: h.defaultThreshold,
			Cooldown:  h.defaultCooldown,
		}
		h.routes[addr] = r
	}
	return r
}

// Healthy reports whether traffic may be attempted against the address now.
// An open circuit becomes a half-open probe once the cooldown has elapsed.
func (h *RouteHealth) Healthy(addr string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := h.routeLocked(addr)
	r.LastAttempt = time.Now()
	switch r.State {
	case RouteOpen:
		if time.Since(r.OpenedAt) >= r.Cooldown {
			r.State = RouteHalfOpen
			return true
		}
		return false
	default:
		return true
	}
}

// RecordSuccess marks a successful attempt, re-closing the circuit.
func (h *RouteHealth) RecordSuccess(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := h.routeLocked(addr)
	r.Failures = 0
	r.State = RouteClosed
	r.LastSuccess = time.Now()
}

// RecordFailure counts a failed attempt and opens the circuit at threshold.
func (h *RouteHealth) RecordFailure(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := h.routeLocked(addr)
	r.Failures++
	if r.State == RouteHalfOpen {
		// A failed half-open probe re-opens the circuit and restarts the
		// cooldown instead of accumulating toward the original threshold.
		r.State = RouteOpen
		r.OpenedAt = time.Now()
		return
	}
	if r.Failures >= r.Threshold {
		r.State = RouteOpen
		r.OpenedAt = time.Now()
	}
}

// State returns the current circuit state for an address.
func (h *RouteHealth) State(addr string) RouteState {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := h.routeLocked(addr)
	return r.State
}

// Snapshot returns a copy of the circuit table sorted by address.
func (h *RouteHealth) Snapshot() []Route {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Route, 0, len(h.routes))
	for _, r := range h.routes {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Addr < out[j].Addr })
	return out
}

// Order reorders the configured addresses for an attempt, putting healthy
// routes first (primary preferred), half-open probes next, and dropping open
// circuits that are still cooling down (they fail fast without a dial).
func (h *RouteHealth) Order(addrs []string) []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	var healthy, probing []string
	for _, a := range addrs {
		r := h.routeLocked(a)
		r.LastAttempt = time.Now()
		switch r.State {
		case RouteOpen:
			if time.Since(r.OpenedAt) >= r.Cooldown {
				r.State = RouteHalfOpen
				probing = append(probing, a)
			}
		default:
			healthy = append(healthy, a)
		}
	}
	return append(healthy, probing...)
}
