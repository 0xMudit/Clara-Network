// Package metrics provides Prometheus-compatible counters and histograms for
// Clara Network services. It uses only the stdlib and a simple text format
// so no external dependency is needed.
package metrics

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing metric.
type Counter struct {
	name   string
	help   string
	labels []string
	mu     sync.Mutex
	// keys -> value
	counts map[string]*atomic.Int64
	order  []string
}

// NewCounter registers a counter with the given name, help text, and label names.
func NewCounter(name, help string, labels ...string) *Counter {
	return &Counter{
		name:   name,
		help:   help,
		labels: labels,
		counts: make(map[string]*atomic.Int64),
	}
}

// Inc increments the counter by 1 with the given label values.
func (c *Counter) Inc(labelValues ...string) {
	c.Add(1, labelValues...)
}

// Add increments the counter by delta.
func (c *Counter) Add(delta int64, labelValues ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := strings.Join(labelValues, ",")
	if v, ok := c.counts[key]; ok {
		v.Add(delta)
	} else {
		v = &atomic.Int64{}
		v.Add(delta)
		c.counts[key] = v
		c.order = append(c.order, key)
	}
}

func (c *Counter) write(w io.Writer) {
	fmt.Fprintf(w, "# TYPE %s counter\n# HELP %s %s\n", c.name, c.name, c.help)
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range c.order {
		v := c.counts[key]
		if len(c.labels) > 0 && key != "" {
			pairs := make([]string, len(c.labels))
			vals := strings.Split(key, ",")
			for i, l := range c.labels {
				if i < len(vals) {
					pairs[i] = fmt.Sprintf("%s=%q", l, vals[i])
				}
			}
			fmt.Fprintf(w, "%s{%s} %d\n", c.name, strings.Join(pairs, ","), v.Load())
		} else {
			fmt.Fprintf(w, "%s %d\n", c.name, v.Load())
		}
	}
}

// Histogram tracks value distributions with fixed buckets.
type Histogram struct {
	name    string
	help    string
	buckets []float64
	mu      sync.Mutex
	// bucket boundary -> count
	counts map[float64]*atomic.Int64
	sum    atomic.Int64
	count  atomic.Int64
}

// NewHistogram registers a histogram.
func NewHistogram(name, help string, buckets ...float64) *Histogram {
	h := &Histogram{name: name, help: help, buckets: buckets, counts: make(map[float64]*atomic.Int64)}
	for _, b := range buckets {
		h.counts[b] = &atomic.Int64{}
	}
	return h
}

// Observe records a value.
func (h *Histogram) Observe(v float64) {
	h.sum.Add(int64(v * 1000)) // store as micro-units for integer atomic
	h.count.Add(1)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, b := range h.buckets {
		if v <= b {
			h.counts[b].Add(1)
		}
	}
}

func (h *Histogram) write(w io.Writer) {
	fmt.Fprintf(w, "# TYPE %s histogram\n# HELP %s %s\n", h.name, h.name, h.help)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, b := range h.buckets {
		fmt.Fprintf(w, "%s_bucket{le=\"%g\"} %d\n", h.name, b, h.counts[b].Load())
	}
	fmt.Fprintf(w, "%s_sum %d\n", h.name, h.sum.Load())
	fmt.Fprintf(w, "%s_count %d\n", h.name, h.count.Load())
}

// Registry holds all registered metrics and serves /metrics.
type Registry struct {
	mu       sync.Mutex
	counters []*Counter
	histograms []*Histogram
}

// NewRegistry returns a new empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// MustRegister adds metrics to the registry.
func (r *Registry) MustRegister(cs ...*Counter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters = append(r.counters, cs...)
}

// MustRegisterHistogram adds histograms to the registry.
func (r *Registry) MustRegisterHistogram(hs ...*Histogram) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.histograms = append(r.histograms, hs...)
}

// Handler returns an HTTP handler that serves Prometheus exposition format.
func (r *Registry) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, c := range r.counters {
			c.write(w)
		}
		for _, h := range r.histograms {
			h.write(w)
		}
	}
}
