// Command resilience-sim runs a chaos drill for the Clara Network resilience
// layer (docs/19): a switch fronts a primary and a secondary issuer, then the
// primary dies (circuit breaker trips, traffic fails over), the secondary dies
// too (stand-in processing approves within limits and declines against
// negative files, 91 responses burst), and finally the primary comes back
// (half-open probe re-closes the circuit).
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
	"github.com/0xMudit/Clara-Network/internal/issuersim"
	"github.com/0xMudit/Clara-Network/internal/resilience"
	"github.com/0xMudit/Clara-Network/internal/switchsrv"
)

const (
	switchAddr    = "127.0.0.1:19080"
	primaryAddr   = "127.0.0.1:19081"
	secondaryAddr = "127.0.0.1:19082"
	issuerID      = "1000001000"

	panNormal = "4000001234567890"
	panHot    = "4000000000000666"
	panRestr  = "4111111234567890"
)

var seq atomic.Int64

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Stand-in policy for the issuer.
	standIn := resilience.NewStandIn(100000)
	standIn.SetPolicy(resilience.Policy{
		IssuerID:       issuerID,
		Enabled:        true,
		Limit:          500000, // approve up to 5,000.00 in stand-in
		NegativeCards:  map[string]bool{panHot: true},
		RestrictedBINs: map[string]bool{"411111": true},
	})

	// Circuit breakers trip after 2 consecutive failures and probe again
	// after a short cooldown so the demo moves quickly.
	health := resilience.NewRouteHealth(2, 3*time.Second)
	metrics := resilience.NewMetrics()

	stopPrimary := startIssuer(ctx, logger, primaryAddr, "primary")
	defer stopPrimary()
	stopSecondary := startIssuer(ctx, logger, secondaryAddr, "secondary")
	defer stopSecondary()

	sw, err := switchsrv.New(switchsrv.Config{
		ListenAddr:   switchAddr,
		IssuerRoutes: map[string]string{issuerID: primaryAddr + "," + secondaryAddr},
		StandIn:      standIn,
		RouteHealth:  health,
		Metrics:      metrics,
		Log:          logger,
	})
	if err != nil {
		logger.Error("switch boot failed", "err", err)
		os.Exit(1)
	}
	go func() { _ = sw.ListenAndServe(ctx) }()

	step(logger, "1. baseline: primary and secondary live")
	runAuths(logger, sw.Addr().String(), []trial{
		{panNormal, 10000, "normal purchase"},
		{panNormal, 20000, "another purchase"},
	})
	snapshot(logger, sw, metrics)

	step(logger, "2. primary dies: traffic fails over to the secondary")
	stopPrimary()
	runAuths(logger, sw.Addr().String(), []trial{
		{panNormal, 10000, "after primary failure"},
		{panNormal, 15000, "still on secondary"},
		{panNormal, 12000, "primary circuit open"},
	})
	snapshot(logger, sw, metrics)

	step(logger, "3. secondary dies too: stand-in processing takes over")
	stopSecondary()
	runAuths(logger, sw.Addr().String(), []trial{
		{panNormal, 10000, "stand-in approve within limit"},
		{panHot, 10000, "hot card -> 05"},
		{panRestr, 10000, "restricted BIN -> 57"},
		{panNormal, 999999999, "above stand-in limit -> 91"},
		{panNormal, 999999999, "another 91 -> burst"},
		{panNormal, 999999999, "third 91 -> burst alert"},
	})
	snapshot(logger, sw, metrics)

	step(logger, "4. primary recovers: half-open probe re-closes the circuit")
	startIssuer(ctx, logger, primaryAddr, "primary (recovered)")
	// Let the cooldown elapse so the probe is admitted.
	logger.Info("waiting for circuit cooldown", "cooldown", 3*time.Second)
	time.Sleep(3 * time.Second)
	runAuths(logger, sw.Addr().String(), []trial{
		{panNormal, 10000, "after recovery"},
		{panNormal, 25000, "recovered primary"},
	})
	snapshot(logger, sw, metrics)

	logger.Info("resilience demo complete")
}

type trial struct {
	pan    string
	amount int
	label  string
}

func step(logger *slog.Logger, s string) {
	logger.Info("")
	logger.Info("== " + s + " ==")
}

func runAuths(logger *slog.Logger, addr string, trials []trial) {
	for _, tr := range trials {
		code, si := auth(logger, addr, tr.pan, tr.amount)
		logger.Info("auth",
			"label", tr.label,
			"pan", tr.pan,
			"amount", tr.amount,
			"code", code,
			"si_marker", si)
	}
}

func auth(logger *slog.Logger, addr string, pan string, amount int) (string, bool) {
	now := time.Now()
	req := iso8583.New("0100").
		Set(2, pan).
		Set(3, "000000").
		Set(4, fmt.Sprintf("%012d", amount)).
		Set(7, now.Format("0102150405")).
		Set(11, fmt.Sprintf("%06d", seq.Add(1))).
		Set(12, now.Format("150405")).
		Set(13, now.Format("0102")).
		Set(22, "022").
		Set(25, "00").
		Set(32, "1000001").
		Set(49, "840").
		Set(100, issuerID)

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		logger.Error("dial switch failed", "err", err)
		return "96", false
	}
	defer conn.Close()

	raw, err := req.Marshal()
	if err != nil {
		logger.Error("marshal request failed", "err", err)
		return "96", false
	}
	if err := framing.WriteFrame(conn, raw); err != nil {
		logger.Error("write request failed", "err", err)
		return "96", false
	}
	frame, err := framing.ReadFrame(conn)
	if err != nil {
		logger.Error("read response failed", "err", err)
		return "96", false
	}
	resp, err := iso8583.Parse(frame)
	if err != nil {
		logger.Error("parse response failed", "err", err)
		return "96", false
	}
	return resp.Get(39), resp.Get(62) == "SI"
}

func snapshot(logger *slog.Logger, sw *switchsrv.Server, metrics *resilience.Metrics) {
	ms, routes := sw.Resilience()
	logger.Info("metrics",
		"total", ms.Total,
		"by_code", ms.ByCode,
		"by_dest", ms.ByDest,
		"standin_approved", ms.StandInApproved,
		"standin_declined", ms.StandInDeclined,
		"p99", ms.P99Latency)
	for _, r := range routes {
		logger.Info("route",
			"addr", r.Addr,
			"state", r.State,
			"failures", r.Failures,
			"threshold", r.Threshold)
	}
	burst := metrics.Burst91(30 * time.Second)
	if burst >= 3 {
		logger.Warn("ISSUER OUTAGE ALERT: 91 burst", "count_30s", burst)
	} else if burst > 0 {
		logger.Info("91 count in window", "count_30s", burst)
	}
}

func startIssuer(ctx context.Context, logger *slog.Logger, addr, name string) func() {
	srv, err := issuersim.New(issuersim.Config{
		ListenAddr: addr,
		ID:         name,
		Decision:   issuersim.DefaultDecision,
		Log:        logger,
	})
	if err != nil {
		logger.Error("issuer boot failed", "name", name, "err", err)
		os.Exit(1)
	}
	go func() { _ = srv.ListenAndServe(ctx) }()
	logger.Info("issuer listening", "name", name, "addr", addr)
	return func() { _ = srv.Close() }
}
