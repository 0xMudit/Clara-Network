package resilience

import (
	"testing"
	"time"
)

func TestStandInApprovesWithinLimit(t *testing.T) {
	s := NewStandIn(100000)
	dec := s.Decide("1000001000", "4000001234567890", 5000)
	if !dec.Approve || dec.Code != "00" {
		t.Fatalf("decide = %+v, want approve", dec)
	}
}

func TestStandInDeclinesAboveLimit(t *testing.T) {
	s := NewStandIn(100000)
	dec := s.Decide("1000001000", "4000001234567890", 999999999)
	if dec.Approve || dec.Code != "91" {
		t.Fatalf("decide = %+v, want 91 decline", dec)
	}
}

func TestStandInPerIssuerLimit(t *testing.T) {
	s := NewStandIn(100000)
	s.SetPolicy(Policy{IssuerID: "1000001000", Enabled: true, Limit: 25000})
	if dec := s.Decide("1000001000", "4000001234567890", 30000); dec.Approve {
		t.Fatalf("high-value tx should decline under per-issuer limit: %+v", dec)
	}
	if dec := s.Decide("1000001000", "4000001234567890", 10000); !dec.Approve {
		t.Fatalf("low-value tx should approve: %+v", dec)
	}
	// A second issuer with no explicit policy uses the default limit.
	if dec := s.Decide("1000002000", "4000001234567890", 50000); !dec.Approve {
		t.Fatalf("default-limit issuer should approve 50000: %+v", dec)
	}
}

func TestStandInDisabledForIssuer(t *testing.T) {
	s := NewStandIn(100000)
	s.SetPolicy(Policy{IssuerID: "1000001000", Enabled: false, Limit: 99999999})
	dec := s.Decide("1000001000", "4000001234567890", 100)
	if dec.Approve || dec.Code != "91" {
		t.Fatalf("disabled stand-in must decline with 91, got %+v", dec)
	}
}

func TestStandInHotCardDecline(t *testing.T) {
	s := NewStandIn(100000)
	s.SetPolicy(Policy{
		IssuerID:      "1000001000",
		Enabled:       true,
		Limit:         100000,
		NegativeCards: map[string]bool{"4000000000000666": true},
	})
	dec := s.Decide("1000001000", "4000000000000666", 100)
	if dec.Approve || dec.Code != "05" {
		t.Fatalf("hot card should decline with 05, got %+v", dec)
	}
}

func TestStandInRestrictedBINDecline(t *testing.T) {
	s := NewStandIn(100000)
	s.SetPolicy(Policy{
		IssuerID:       "1000001000",
		Enabled:        true,
		Limit:          100000,
		RestrictedBINs: map[string]bool{"411111": true},
	})
	dec := s.Decide("1000001000", "4111111234567890", 100)
	if dec.Approve || dec.Code != "57" {
		t.Fatalf("restricted BIN should decline with 57, got %+v", dec)
	}
}

func TestStandInPositiveFile(t *testing.T) {
	s := NewStandIn(100000)
	s.SetPolicy(Policy{
		IssuerID:   "1000001000",
		Enabled:    true,
		Limit:      100000,
		ValidCards: map[string]bool{"4000001234567890": true},
	})
	if dec := s.Decide("1000001000", "4000001234567890", 100); !dec.Approve {
		t.Fatalf("listed card should approve, got %+v", dec)
	}
	if dec := s.Decide("1000001000", "4000001234567899", 100); dec.Approve {
		t.Fatalf("unlisted card must decline, got %+v", dec)
	}
}

func TestCircuitOpensAfterThreshold(t *testing.T) {
	h := NewRouteHealth(2, time.Second)
	if !h.Healthy("127.0.0.1:1") {
		t.Fatal("fresh route must be healthy")
	}
	h.RecordFailure("127.0.0.1:1")
	if h.State("127.0.0.1:1") != RouteClosed {
		t.Fatalf("one failure must not trip, got %s", h.State("127.0.0.1:1"))
	}
	h.RecordFailure("127.0.0.1:1")
	if h.State("127.0.0.1:1") != RouteOpen {
		t.Fatalf("circuit should be open after threshold, got %s", h.State("127.0.0.1:1"))
	}
	if h.Healthy("127.0.0.1:1") {
		t.Fatal("open circuit must be unhealthy before cooldown")
	}
}

func TestCircuitHalfOpenProbeRecovers(t *testing.T) {
	h := NewRouteHealth(2, 50*time.Millisecond)
	h.RecordFailure("a:1")
	h.RecordFailure("a:1")
	if h.State("a:1") != RouteOpen {
		t.Fatalf("circuit should open, got %s", h.State("a:1"))
	}
	// After cooldown the probe is admitted (half-open).
	time.Sleep(60 * time.Millisecond)
	if !h.Healthy("a:1") {
		t.Fatal("half-open probe should be admitted after cooldown")
	}
	if h.State("a:1") != RouteHalfOpen {
		t.Fatalf("probe state should be half-open, got %s", h.State("a:1"))
	}
	h.RecordSuccess("a:1")
	if h.State("a:1") != RouteClosed {
		t.Fatalf("success should re-close the circuit, got %s", h.State("a:1"))
	}
}

func TestCircuitProbeFailureReopens(t *testing.T) {
	h := NewRouteHealth(2, 50*time.Millisecond)
	h.RecordFailure("a:1")
	h.RecordFailure("a:1")
	time.Sleep(60 * time.Millisecond)
	h.Healthy("a:1") // admit probe
	h.RecordFailure("a:1")
	if h.State("a:1") != RouteOpen {
		t.Fatalf("failed probe must re-open the circuit, got %s", h.State("a:1"))
	}
}

func TestOrderDropsCoolingOpenCircuit(t *testing.T) {
	h := NewRouteHealth(2, time.Hour)
	addrs := []string{"primary:1", "secondary:2"}
	h.RecordFailure("primary:1")
	h.RecordFailure("primary:1") // primary open
	ordered := h.Order(addrs)
	if len(ordered) != 1 || ordered[0] != "secondary:2" {
		t.Fatalf("order = %v, want [secondary:2] only", ordered)
	}
	// Recovery via success on the secondary; primary is still cooling down.
	h.RecordSuccess("secondary:2")
	ordered = h.Order(addrs)
	if len(ordered) != 1 || ordered[0] != "secondary:2" {
		t.Fatalf("order = %v, want [secondary:2] while primary cools down", ordered)
	}
	// Once the primary also succeeds it is preferred first again.
	h.RecordSuccess("primary:1")
	ordered = h.Order(addrs)
	if len(ordered) != 2 || ordered[0] != "primary:1" || ordered[1] != "secondary:2" {
		t.Fatalf("order = %v, want [primary:1 secondary:2]", ordered)
	}
}

func TestMetricsCountsAndStandIn(t *testing.T) {
	m := NewMetrics()
	m.Record(OpResult{Code: "00", Dest: "1000001000"})
	m.Record(OpResult{Code: "00", Dest: "1000001000", StandIn: true})
	m.Record(OpResult{Code: "91", Dest: "1000001000", StandIn: true})
	m.Record(OpResult{Code: "05", Dest: "1000001000"})
	snap := m.Snapshot()
	if snap.Total != 4 {
		t.Fatalf("total = %d, want 4", snap.Total)
	}
	if snap.ByCode["00"] != 2 || snap.ByCode["91"] != 1 || snap.ByCode["05"] != 1 {
		t.Fatalf("byCode = %v", snap.ByCode)
	}
	if snap.StandInApproved != 1 || snap.StandInDeclined != 1 {
		t.Fatalf("standin approved=%d declined=%d, want 1/1", snap.StandInApproved, snap.StandInDeclined)
	}
}

func TestMetricsP99AndBurst(t *testing.T) {
	m := NewMetrics()
	for i := 0; i < 95; i++ {
		m.Record(OpResult{Code: "00", Latency: 500 * time.Microsecond})
	}
	for i := 0; i < 5; i++ {
		m.Record(OpResult{Code: "00", Latency: 400 * time.Millisecond})
	}
	if m.P99() != latencyBounds[8] { // 500ms bucket: the 99th percentile is in the tail
		t.Fatalf("p99 = %v, want 500ms", m.P99())
	}
	// Three 91s in a short window trip a burst alert.
	m.Record(OpResult{Code: "91", Dest: "1000001000", StandIn: true})
	m.Record(OpResult{Code: "91", Dest: "1000001000", StandIn: true})
	m.Record(OpResult{Code: "91", Dest: "1000001000", StandIn: true})
	if n := m.Burst91(time.Minute); n != 3 {
		t.Fatalf("burst count = %d, want 3", n)
	}
}
