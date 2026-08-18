package metrics

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCounterInc(t *testing.T) {
	c := NewCounter("test_counter", "test help")
	c.Inc()
	c.Inc("label1")
	c.Inc("label1")
	if v := c.counts[""].Load(); v != 1 {
		t.Fatalf("bare count = %d, want 1", v)
	}
	if v := c.counts["label1"].Load(); v != 2 {
		t.Fatalf("label1 count = %d, want 2", v)
	}
}

func TestHistogramObserve(t *testing.T) {
	h := NewHistogram("test_hist", "test", 10, 100, 1000)
	h.Observe(5)
	h.Observe(50)
	h.Observe(500)
	h.Observe(5000)
	if c := h.count.Load(); c != 4 {
		t.Fatalf("count = %d, want 4", c)
	}
	if v := h.counts[10].Load(); v != 1 {
		t.Fatalf("le=10 count = %d, want 1", v)
	}
	if v := h.counts[100].Load(); v != 2 {
		t.Fatalf("le=100 count = %d, want 2", v)
	}
	if v := h.counts[1000].Load(); v != 3 {
		t.Fatalf("le=1000 count = %d, want 3", v)
	}
}

func TestRegistryHandler(t *testing.T) {
	r := NewRegistry()
	c := NewCounter("http_requests_total", "Total requests", "method", "status")
	r.MustRegister(c)
	h := NewHistogram("http_duration_seconds", "Request duration", 0.1, 1, 10)
	r.MustRegisterHistogram(h)
	c.Inc("GET", "200")
	c.Inc("POST", "201")
	h.Observe(0.05)
	h.Observe(0.5)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	r.Handler()(w, req)
	body := w.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Fatal("response missing counter")
	}
	if !strings.Contains(body, "http_duration_seconds") {
		t.Fatal("response missing histogram")
	}
	if w.Header().Get("Content-Type") != "text/plain; version=0.0.4; charset=utf-8" {
		t.Fatalf("content-type = %q", w.Header().Get("Content-Type"))
	}
}

func TestCounterLabelValues(t *testing.T) {
	c := NewCounter("labeled", "help", "a", "b")
	c.Inc("x", "y")
	c.Inc("x", "y")
	c.Inc("x", "z")
	if v := c.counts["x,y"].Load(); v != 2 {
		t.Fatalf("x,y = %d, want 2", v)
	}
	if v := c.counts["x,z"].Load(); v != 1 {
		t.Fatalf("x,z = %d, want 1", v)
	}
}
