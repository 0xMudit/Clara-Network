package resilience

import (
	"sync"
	"time"
)

// OpResult is the observed outcome of one authorization attempt.
type OpResult struct {
	Code    string
	Dest    string
	StandIn bool
	Latency time.Duration
}

// MetricsSnapshot is a point-in-time view of switch health.
type MetricsSnapshot struct {
	Total           int64
	ByCode          map[string]int64
	ByDest          map[string]int64
	StandInApproved int64
	StandInDeclined int64
	P99Latency      time.Duration
	Now             time.Time
}

// latencyBounds are the upper bounds of the p99 histogram buckets.
var latencyBounds = []time.Duration{
	time.Millisecond,
	2 * time.Millisecond,
	5 * time.Millisecond,
	10 * time.Millisecond,
	25 * time.Millisecond,
	50 * time.Millisecond,
	100 * time.Millisecond,
	250 * time.Millisecond,
	500 * time.Millisecond,
	time.Second,
	5 * time.Second,
}

// Metrics records authorization outcomes for monitoring dashboards (docs/19
// §19.4): response-code distributions, per-destination volume, stand-in
// usage, approximate p99 latency, and bursts of issuer-inoperative (91)
// responses that flag an issuer outage.
type Metrics struct {
	mu              sync.Mutex
	total           int64
	byCode          map[string]int64
	byDest          map[string]int64
	standInApproved int64
	standInDeclined int64
	latBuckets      []int64
	burst91         []time.Time
}

// NewMetrics returns an empty metrics accumulator.
func NewMetrics() *Metrics {
	return &Metrics{
		byCode:     map[string]int64{},
		byDest:     map[string]int64{},
		latBuckets: make([]int64, len(latencyBounds)+1),
	}
}

// Record observes a single authorization outcome.
func (m *Metrics) Record(r OpResult) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.total++
	m.byCode[r.Code]++
	m.byDest[r.Dest]++
	if r.StandIn {
		if r.Code == "00" {
			m.standInApproved++
		} else {
			m.standInDeclined++
		}
	}
	if r.Code == "91" {
		m.burst91 = append(m.burst91, time.Now())
	}
	m.latBuckets[latBucket(r.Latency)]++
}

func latBucket(d time.Duration) int {
	for i, b := range latencyBounds {
		if d < b {
			return i
		}
	}
	return len(latencyBounds)
}

// Burst91 counts issuer-inoperative responses observed inside the window,
// trimming older entries as it goes.
func (m *Metrics) Burst91(window time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := time.Now().Add(-window)
	for len(m.burst91) > 0 && m.burst91[0].Before(cutoff) {
		m.burst91 = m.burst91[1:]
	}
	return len(m.burst91)
}

// P99 estimates the 99th-percentile latency from the histogram buckets.
func (m *Metrics) P99() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.p99Locked()
}

func (m *Metrics) p99Locked() time.Duration {
	if m.total == 0 {
		return 0
	}
	target := (m.total*99 + 99) / 100
	cum := int64(0)
	for i, c := range m.latBuckets {
		cum += c
		if cum >= target {
			if i >= len(latencyBounds) {
				return 5 * time.Second
			}
			return latencyBounds[i]
		}
	}
	return 5 * time.Second
}

// Snapshot returns a point-in-time view of the metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	codes := make(map[string]int64, len(m.byCode))
	for k, v := range m.byCode {
		codes[k] = v
	}
	dests := make(map[string]int64, len(m.byDest))
	for k, v := range m.byDest {
		dests[k] = v
	}
	return MetricsSnapshot{
		Total:           m.total,
		ByCode:          codes,
		ByDest:          dests,
		StandInApproved: m.standInApproved,
		StandInDeclined: m.standInDeclined,
		P99Latency:      m.p99Locked(),
		Now:             time.Now(),
	}
}
