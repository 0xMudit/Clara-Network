package switchsrv

import (
	"net"
	"testing"
	"time"

	"github.com/0xMudit/Clara-Network/internal/binrouting"
	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/issuersim"
	"github.com/0xMudit/Clara-Network/internal/risk"
)

func velocityEngine(limit int64) *risk.Engine {
	return risk.New(risk.NewMemoryStore(), []risk.Rule{
		{Name: "card-velocity", Kind: risk.KindCardVelocity, Limit: limit, Window: 60, Code: "59", Enabled: true},
	})
}

func TestRouteByBIN(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	tab := binrouting.New(&binrouting.Config{Entries: map[string]string{"400000": "1000001000"}})
	sw, stopSwitch := startSwitchCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": issuer},
		BINTable:     tab,
	})
	defer stopSwitch()

	req := authReq("4000001234567890", 2500)
	delete(req.Fields, 100)
	resp := send(t, sw, req)
	if resp.Get(39) != "00" {
		t.Fatalf("expected approval via BIN routing, got %q", resp.Get(39))
	}
	if resp.Get(100) != "1000001000" {
		t.Fatalf("response should echo derived DE100, got %q", resp.Get(100))
	}
}

func TestFormatErrorWhenUnknownBIN(t *testing.T) {
	tab := binrouting.New(&binrouting.Config{Entries: map[string]string{"400000": "1000001000"}})
	sw, stopSwitch := startSwitchCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": "127.0.0.1:1"},
		BINTable:     tab,
	})
	defer stopSwitch()

	req := authReq("4111111234567890", 1000)
	delete(req.Fields, 100)
	resp := send(t, sw, req)
	if resp.Get(39) != "30" {
		t.Fatalf("expected format error 30 for unknown BIN, got %q", resp.Get(39))
	}
}

func TestIssuerFailover(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": "127.0.0.1:1, " + issuer})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234567890", 2500))
	if resp.Get(39) != "00" {
		t.Fatalf("failover should reach the live issuer, got %q", resp.Get(39))
	}
}

func TestRiskDeclinesAfterVelocityLimit(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitchCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": issuer},
		Risk:         velocityEngine(1),
	})
	defer stopSwitch()

	first := authReq("4000001234567890", 2500)
	first.Set(11, "100000")
	if resp := send(t, sw, first); resp.Get(39) != "00" {
		t.Fatalf("first request should approve, got %q", resp.Get(39))
	}

	second := authReq("4000001234567890", 2500)
	second.Set(11, "100001")
	if resp := send(t, sw, second); resp.Get(39) != "59" {
		t.Fatalf("second request should be risk-declined with 59, got %q", resp.Get(39))
	}
}

func TestRiskReplayDoesNotDoubleCount(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	// Limit 2: without double counting the replay, two fresh transactions
	// are allowed and the third is declined.
	sw, stopSwitch := startSwitchCfg(t, Config{
		IssuerRoutes: map[string]string{"1000001000": issuer},
		Risk:         velocityEngine(2),
	})
	defer stopSwitch()

	req := authReq("4000001234567890", 2500)
	if resp := send(t, sw, req); resp.Get(39) != "00" {
		t.Fatalf("first request should approve, got %q", resp.Get(39))
	}
	if resp := send(t, sw, req); resp.Get(39) != "00" {
		t.Fatalf("idempotent replay must return cached approval, got %q", resp.Get(39))
	}

	// Second fresh transaction still within limit: the replay must not have
	// incremented the counter a second time.
	next := authReq("4000001234567890", 2500)
	next.Set(11, "200000")
	if resp := send(t, sw, next); resp.Get(39) != "00" {
		t.Fatalf("velocity counter must not double on replay, got %q", resp.Get(39))
	}

	// Third fresh transaction exceeds the limit.
	last := authReq("4000001234567890", 2500)
	last.Set(11, "300000")
	if resp := send(t, sw, last); resp.Get(39) != "59" {
		t.Fatalf("third transaction should be risk-declined, got %q", resp.Get(39))
	}
}

func BenchmarkAuthorizationRoundTrip(b *testing.B) {
	issuer, stopIssuer := startIssuer(b, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitch(b, map[string]string{"1000001000": issuer})
	defer stopSwitch()

	conn, err := net.DialTimeout("tcp", sw, 5*time.Second)
	if err != nil {
		b.Fatalf("dial switch: %v", err)
	}
	defer conn.Close()

	req := authReq("4000001234567890", 2500)
	raw, err := req.Marshal()
	if err != nil {
		b.Fatalf("marshal request: %v", err)
	}
	if err := framing.WriteFrame(conn, raw); err != nil {
		b.Fatalf("warmup write: %v", err)
	}
	if _, err := framing.ReadFrame(conn); err != nil {
		b.Fatalf("warmup read: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := framing.WriteFrame(conn, raw); err != nil {
			b.Fatal(err)
		}
		if _, err := framing.ReadFrame(conn); err != nil {
			b.Fatal(err)
		}
	}
}
