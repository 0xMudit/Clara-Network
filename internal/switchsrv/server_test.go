package switchsrv

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
	"github.com/0xMudit/Clara-Network/internal/issuersim"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func startIssuer(t *testing.T, decision issuersim.DecisionFunc) (string, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	srv, err := issuersim.New(issuersim.Config{ListenAddr: "127.0.0.1:0", Decision: decision, Log: testLogger()})
	if err != nil {
		t.Fatalf("start issuer: %v", err)
	}
	go func() { _ = srv.ListenAndServe(ctx) }()
	return srv.Addr().String(), cancel
}

func startSwitch(t *testing.T, routes map[string]string) (string, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	srv, err := New(Config{
		ListenAddr:     "127.0.0.1:0",
		IssuerRoutes:   routes,
		Idempotency:    NewMemoryIdempotency(),
		IdempotencyTTL: 5 * time.Second,
		Log:            testLogger(),
	})
	if err != nil {
		t.Fatalf("start switch: %v", err)
	}
	go func() { _ = srv.ListenAndServe(ctx) }()
	return srv.Addr().String(), cancel
}

func authReq(pan string, amount int) *iso8583.Message {
	now := time.Now()
	return iso8583.New("0100").
		Set(2, pan).
		Set(3, "000000").
		Set(4, fmt.Sprintf("%012d", amount)).
		Set(7, now.Format("0102150405")).
		Set(11, "123456").
		Set(12, now.Format("150405")).
		Set(13, now.Format("0102")).
		Set(22, "022").
		Set(25, "00").
		Set(32, "1000001").
		Set(49, "840").
		Set(100, "1000001000")
}

func send(t *testing.T, addr string, req *iso8583.Message) *iso8583.Message {
	t.Helper()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		t.Fatalf("dial switch: %v", err)
	}
	defer conn.Close()
	raw, err := req.Marshal()
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	if err := framing.WriteFrame(conn, raw); err != nil {
		t.Fatalf("write request: %v", err)
	}
	frame, err := framing.ReadFrame(conn)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	resp, err := iso8583.Parse(frame)
	if err != nil {
		t.Fatalf("parse response: %v", err)
	}
	return resp
}

func TestAuthorizationApproved(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": issuer})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234567890", 2500))
	if resp.Get(39) != "00" {
		t.Fatalf("expected approval, got response code %q", resp.Get(39))
	}
	if resp.MTI != "0110" {
		t.Fatalf("expected MTI 0110, got %q", resp.MTI)
	}
}

func TestAuthorizationDeclinedDoNotHonor(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": issuer})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234560000", 2500))
	if resp.Get(39) != "05" {
		t.Fatalf("expected do-not-honor, got response code %q", resp.Get(39))
	}
}

func TestAuthorizationInsufficientFunds(t *testing.T) {
	issuer, stopIssuer := startIssuer(t, issuersim.DefaultDecision)
	defer stopIssuer()
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": issuer})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234567890", 600000))
	if resp.Get(39) != "51" {
		t.Fatalf("expected insufficient funds, got response code %q", resp.Get(39))
	}
}

func TestIdempotentReplayReturnsCachedResponse(t *testing.T) {
	var decision atomic.Value
	decision.Store("00")
	issuer, stopIssuer := startIssuer(t, func(*iso8583.Message) string {
		return decision.Load().(string)
	})
	defer stopIssuer()
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": issuer})
	defer stopSwitch()

	req := authReq("4000001234567890", 2500)
	first := send(t, sw, req)
	if first.Get(39) != "00" {
		t.Fatalf("first request should approve, got %q", first.Get(39))
	}

	// The issuer would now decline, but the replay must return the cached
	// original response without reaching it.
	decision.Store("51")
	second := send(t, sw, req)
	if second.Get(39) != "00" {
		t.Fatalf("replay should return cached approval, got %q", second.Get(39))
	}
}

func TestStandInApprovesWithinLimit(t *testing.T) {
	// No issuer: route points at a closed port.
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": "127.0.0.1:1"})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234567890", 5000))
	if resp.Get(39) != "00" {
		t.Fatalf("stand-in should approve within limit, got %q", resp.Get(39))
	}
	if resp.Get(62) != "SI" {
		t.Fatalf("stand-in response should carry SI marker, got %q", resp.Get(62))
	}
}

func TestStandInDeclinesAboveLimit(t *testing.T) {
	sw, stopSwitch := startSwitch(t, map[string]string{"1000001000": "127.0.0.1:1"})
	defer stopSwitch()

	resp := send(t, sw, authReq("4000001234567890", 999999999))
	if resp.Get(39) != "91" {
		t.Fatalf("stand-in should decline above limit with 91, got %q", resp.Get(39))
	}
}

func TestFormatErrorWhenNoDestination(t *testing.T) {
	sw, stopSwitch := startSwitch(t, map[string]string{})
	defer stopSwitch()

	req := authReq("4000001234567890", 1000)
	delete(req.Fields, 100)
	resp := send(t, sw, req)
	if resp.Get(39) != "30" {
		t.Fatalf("expected format error 30, got %q", resp.Get(39))
	}
}
