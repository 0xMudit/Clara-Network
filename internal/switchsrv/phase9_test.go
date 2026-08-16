package switchsrv

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xMudit/Clara-Network/internal/binrouting"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
	"github.com/0xMudit/Clara-Network/internal/issuersim"
	"github.com/0xMudit/Clara-Network/internal/resilience"
)

// authReqSTAN builds an authorization request with a caller-chosen STAN so
// tests are not deduplicated by the idempotency store.
func authReqSTAN(pan string, amount int, stan string) *iso8583.Message {
	req := authReq(pan, amount)
	req.Set(11, stan)
	return req
}

// startSwitchServerCfg starts the switch and returns the server so tests can
// inspect its Resilience() snapshot; the listener address is via Addr().
func startSwitchServerCfg(t testing.TB, cfg Config) (*Server, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cfg.ListenAddr = "127.0.0.1:0"
	cfg.Log = testLogger()
	if cfg.Idempotency == nil {
		cfg.Idempotency = NewMemoryIdempotency()
	}
	if cfg.IdempotencyTTL == 0 {
		cfg.IdempotencyTTL = 5 * time.Second
	}
	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("start switch: %v", err)
	}
	go func() { _ = srv.ListenAndServe(ctx) }()
	return srv, cancel
}

func TestStandInPolicyViaSwitch(t *testing.T) {
	standIn := resilience.NewStandIn(100000)
	standIn.SetPolicy(resilience.Policy{
		IssuerID:       "1000001000",
		Enabled:        true,
		Limit:          25000,
		NegativeCards:  map[string]bool{"4000000000000666": true},
		RestrictedBINs: map[string]bool{"411111": true},
	})
	sw, stopSwitch := startSwitchServerCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1"},
		StandIn:      standIn,
	})
	defer stopSwitch()

	cases := []struct {
		name, pan string
		amount    int
		stan      string
		wantCode  string
		wantSI    bool
	}{
		{"low-value approves", "4000001234567890", 5000, "100001", "00", true},
		{"hot card declines", "4000000000000666", 1000, "100002", "05", false},
		{"restricted BIN declines", "4111111234567890", 1000, "100003", "57", false},
		{"above limit declines", "4000001234567890", 999999999, "100004", "91", false},
	}
	for _, tc := range cases {
		resp := send(t, sw.Addr().String(), authReqSTAN(tc.pan, tc.amount, tc.stan))
		if resp.Get(39) != tc.wantCode {
			t.Errorf("%s: code = %q, want %q", tc.name, resp.Get(39), tc.wantCode)
		}
		gotSI := resp.Get(62) == standInMarker
		if gotSI != tc.wantSI {
			t.Errorf("%s: SI marker present = %t, want %t", tc.name, gotSI, tc.wantSI)
		}
	}
}

func TestCircuitBreakerSkipsDownPrimary(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	health := resilience.NewRouteHealth(2, time.Hour)
	sw, stopSwitch := startSwitchServerCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1, " + issuer},
		RouteHealth:  health,
	})
	defer stopSwitch()

	for i := 0; i < 3; i++ {
		if resp := send(t, sw.Addr().String(), authReqSTAN("4000001234567890", 2500, fmt.Sprintf("20000%d", i))); resp.Get(39) != "00" {
			t.Fatalf("attempt %d should reach the live secondary, got %q", i, resp.Get(39))
		}
	}

	_, routes := sw.Resilience()
	var primary, secondary *resilience.Route
	for i := range routes {
		if routes[i].Addr == "127.0.0.1:1" {
			primary = &routes[i]
		}
		if routes[i].Addr == issuer {
			secondary = &routes[i]
		}
	}
	if primary == nil || primary.State != resilience.RouteOpen {
		t.Fatalf("primary should be open after two failures, got %+v", primary)
	}
	if primary.Failures != 2 {
		t.Fatalf("primary failures = %d, want 2 (third attempt skipped)", primary.Failures)
	}
	if secondary == nil || secondary.State != resilience.RouteClosed || secondary.Failures != 0 {
		t.Fatalf("secondary should stay closed, got %+v", secondary)
	}
}

func TestMetricsTrackedAcrossOutcomes(t *testing.T) {
	standIn := resilience.NewStandIn(100000)
	metrics := resilience.NewMetrics()
	sw, stopSwitch := startSwitchServerCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1"},
		StandIn:      standIn,
		Metrics:      metrics,
	})
	defer stopSwitch()

	if resp := send(t, sw.Addr().String(), authReqSTAN("4000001234567890", 5000, "300001")); resp.Get(39) != "00" {
		t.Fatalf("stand-in approval expected, got %q", resp.Get(39))
	}
	if resp := send(t, sw.Addr().String(), authReqSTAN("4000001234567890", 999999999, "300002")); resp.Get(39) != "91" {
		t.Fatalf("stand-in decline expected, got %q", resp.Get(39))
	}

	snap, _ := sw.Resilience()
	if snap.Total != 2 {
		t.Fatalf("total = %d, want 2", snap.Total)
	}
	if snap.ByCode["00"] != 1 || snap.ByCode["91"] != 1 {
		t.Fatalf("byCode = %v, want 00:1 91:1", snap.ByCode)
	}
	if snap.StandInApproved != 1 || snap.StandInDeclined != 1 {
		t.Fatalf("standin approved=%d declined=%d, want 1/1", snap.StandInApproved, snap.StandInDeclined)
	}
	if n := metrics.Burst91(time.Minute); n != 1 {
		t.Fatalf("91 burst count = %d, want 1", n)
	}
	if snap.ByDest["1000001000"] != 2 {
		t.Fatalf("byDest = %v, want dest 1000001000:2", snap.ByDest)
	}
}

func TestStandInDisabledFallsBackTo91(t *testing.T) {
	standIn := resilience.NewStandIn(100000)
	standIn.SetPolicy(resilience.Policy{IssuerID: "1000001000", Enabled: false})
	sw, stopSwitch := startSwitchServerCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1"},
		StandIn:      standIn,
	})
	defer stopSwitch()

	resp := send(t, sw.Addr().String(), authReqSTAN("4000001234567890", 5000, "400001"))
	if resp.Get(39) != "91" {
		t.Fatalf("disabled stand-in must decline with 91, got %q", resp.Get(39))
	}
}

func TestBINRouteThenStandInWithPositiveFile(t *testing.T) {
	tab := binrouting.New(&binrouting.Config{Entries: map[string]string{"400000": "1000001000"}})
	standIn := resilience.NewStandIn(100000)
	standIn.SetPolicy(resilience.Policy{
		IssuerID:   "1000001000",
		Enabled:    true,
		Limit:      100000,
		ValidCards: map[string]bool{"4000001234567890": true},
	})
	sw, stopSwitch := startSwitchServerCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1"},
		BINTable:     tab,
		StandIn:      standIn,
	})
	defer stopSwitch()

	listed := authReqSTAN("4000001234567890", 5000, "500001")
	delete(listed.Fields, 100)
	if resp := send(t, sw.Addr().String(), listed); resp.Get(39) != "00" {
		t.Fatalf("listed card should be stand-in approved, got %q", resp.Get(39))
	}
	unlisted := authReqSTAN("4000001234567899", 5000, "500002")
	delete(unlisted.Fields, 100)
	if resp := send(t, sw.Addr().String(), unlisted); resp.Get(39) != "91" {
		t.Fatalf("unlisted card must decline with 91, got %q", resp.Get(39))
	}
}
