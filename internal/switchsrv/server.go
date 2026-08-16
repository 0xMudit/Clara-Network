// Package switchsrv implements the Clara Network ISO 8583 message switch:
// it accepts authorization requests from acquirers, routes them to the
// correct issuer, enforces idempotent replay, and applies stand-in
// processing when an issuer is unreachable.
package switchsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/0xMudit/Clara-Network/internal/binrouting"
	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
	"github.com/0xMudit/Clara-Network/internal/risk"
)

const (
	respCodeFormatError    = "30" // Format error
	respCodeStandInDecline = "91" // Issuer or switch inoperative
	respCodeSystemError    = "96" // System malfunction
	respCodeRiskFraud      = "59" // Suspected fraud
	standInMarker          = "SI"
)

// Config configures the switch.
type Config struct {
	ListenAddr     string
	IssuerRoutes   map[string]string // receiving institution ID (DE100) -> comma-separated issuer host:port list
	BINTable       *binrouting.Table // derives DE100 from the PAN when absent
	Risk           *risk.Engine      // in-path risk evaluation; nil disables
	Idempotency    IdempotencyStore
	Audit          AuditSink
	IdempotencyTTL time.Duration
	StandInLimit   int64 // max amount (minor units) authorized under stand-in
	Log            *slog.Logger
}

// Server is the Clara Network ISO 8583 switch.
type Server struct {
	cfg Config
	ln  net.Listener
	log *slog.Logger
}

// New binds the switch listener and returns a ready server.
func New(cfg Config) (*Server, error) {
	if cfg.IdempotencyTTL == 0 {
		cfg.IdempotencyTTL = 60 * time.Second
	}
	if cfg.StandInLimit == 0 {
		cfg.StandInLimit = 100000
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	ln, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("switch: listen: %w", err)
	}
	return &Server{cfg: cfg, ln: ln, log: cfg.Log}, nil
}

// Addr returns the bound listener address (useful with port 0 in tests).
func (s *Server) Addr() net.Addr { return s.ln.Addr() }

// ListenAndServe accepts connections until ctx is cancelled.
func (s *Server) ListenAndServe(ctx context.Context) error {
	s.log.Info("switch listening", "addr", s.ln.Addr())
	go func() {
		<-ctx.Done()
		s.ln.Close()
	}()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("switch: accept: %w", err)
		}
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		frame, err := framing.ReadFrame(conn)
		if err != nil {
			return
		}
		resp, err := s.process(ctx, frame)
		if err != nil {
			s.log.Warn("processing failed", "err", err)
			return
		}
		if err := framing.WriteFrame(conn, resp); err != nil {
			return
		}
	}
}

func (s *Server) process(ctx context.Context, frame []byte) ([]byte, error) {
	req, err := iso8583.Parse(frame)
	if err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}
	switch req.MTI {
	case "0100", "0200", "0420":
		return s.handleAuth(ctx, req)
	default:
		return s.marshal(s.buildResponse(req, respCodeSystemError)), nil
	}
}

func (s *Server) handleAuth(ctx context.Context, req *iso8583.Message) ([]byte, error) {
	destID := req.Get(100)
	if destID == "" {
		destID = s.routeByBIN(req)
	}
	if destID == "" {
		return s.marshal(s.buildResponse(req, respCodeFormatError)), nil
	}

	key := idemKey(req)
	if s.cfg.Idempotency != nil {
		if cached, ok, err := s.cfg.Idempotency.Get(ctx, key); err == nil && ok {
			if resp, perr := iso8583.Parse([]byte(cached)); perr == nil {
				s.log.Debug("idempotent replay", "stan", req.Get(11), "dest", destID)
				s.audit(req, resp)
				return []byte(cached), nil
			}
		}
	}

	if respBytes, declined, err := s.checkRisk(ctx, req); declined || err != nil {
		if declined {
			return respBytes, nil
		}
		return nil, err
	}

	respBytes, err := s.forward(ctx, req, destID)
	if err != nil {
		s.log.Warn("issuer unavailable", "dest", destID, "err", err)
		resp := s.standIn(req)
		if respBytes, err = resp.Marshal(); err != nil {
			return nil, err
		}
	} else if s.cfg.Idempotency != nil {
		if err := s.cfg.Idempotency.Set(ctx, key, string(respBytes), s.cfg.IdempotencyTTL); err != nil {
			s.log.Debug("idempotency store set failed", "err", err)
		}
	}

	resp, err := iso8583.Parse(respBytes)
	if err != nil {
		return nil, err
	}
	s.audit(req, resp)
	return respBytes, nil
}

// routeByBIN derives the destination institution from the PAN using the BIN
// table and records it as DE100 on the request so downstream parties agree.
func (s *Server) routeByBIN(req *iso8583.Message) string {
	if s.cfg.BINTable == nil {
		return ""
	}
	id, ok := s.cfg.BINTable.Lookup(req.Get(2))
	if !ok {
		s.log.Warn("no BIN route", "bin", binOf(req.Get(2)))
		return ""
	}
	req.Set(100, id)
	return id
}

// checkRisk evaluates the request against the risk engine. When declined it
// returns a marshalled response and declined=true; the caller short-circuits.
func (s *Server) checkRisk(ctx context.Context, req *iso8583.Message) ([]byte, bool, error) {
	if s.cfg.Risk == nil {
		return nil, false, nil
	}
	amount, err := strconv.ParseInt(req.Get(4), 10, 64)
	if err != nil {
		return s.marshal(s.buildResponse(req, respCodeFormatError)), true, nil
	}
	dec, err := s.cfg.Risk.Evaluate(ctx, req.Get(2), req.Get(41), amount)
	if err != nil {
		s.log.Warn("risk evaluation failed", "err", err)
		return nil, false, nil
	}
	if dec.Allow {
		return nil, false, nil
	}
	code := dec.Code
	if code == "" {
		code = respCodeRiskFraud
	}
	s.log.Info("risk decline", "reason", dec.Reason, "code", code,
		"stan", req.Get(11), "pan", maskPAN(req.Get(2)))
	resp := s.buildResponse(req, code)
	raw, err := resp.Marshal()
	if err != nil {
		return nil, false, err
	}
	s.audit(req, resp)
	return raw, true, nil
}

// forward relays the request to an issuer, trying each configured address in
// order until one succeeds.
func (s *Server) forward(ctx context.Context, req *iso8583.Message, destID string) ([]byte, error) {
	addrs, ok := s.cfg.IssuerRoutes[destID]
	if !ok {
		return nil, fmt.Errorf("no route for receiving institution %q", destID)
	}
	var lastErr error
	for _, addr := range strings.Split(addrs, ",") {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		raw, err := s.dialForward(ctx, req, addr)
		if err == nil {
			return raw, nil
		}
		lastErr = err
		s.log.Warn("issuer attempt failed", "dest", destID, "addr", addr, "err", err)
	}
	return nil, lastErr
}

func (s *Server) dialForward(ctx context.Context, req *iso8583.Message, addr string) ([]byte, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	raw, err := req.Marshal()
	if err != nil {
		return nil, err
	}
	if err := framing.WriteFrame(conn, raw); err != nil {
		return nil, err
	}
	return framing.ReadFrame(conn)
}

// standIn authorizes on behalf of an unreachable issuer: it approves within
// the configured limit (marking the response) and declines otherwise.
func (s *Server) standIn(req *iso8583.Message) *iso8583.Message {
	resp := s.buildResponse(req, respCodeStandInDecline)
	amount, err := strconv.ParseInt(req.Get(4), 10, 64)
	if err == nil && amount <= s.cfg.StandInLimit {
		resp.Set(39, "00")
		resp.Set(62, standInMarker)
	}
	return resp
}

func (s *Server) buildResponse(req *iso8583.Message, code string) *iso8583.Message {
	resp := iso8583.New(responseMTI(req.MTI))
	for _, f := range []int{7, 11, 12, 13, 32, 100} {
		if req.Has(f) {
			resp.Set(f, req.Get(f))
		}
	}
	resp.Set(39, code)
	return resp
}

func responseMTI(mti string) string {
	switch mti {
	case "0100":
		return "0110"
	case "0200":
		return "0210"
	case "0420":
		return "0430"
	default:
		return "0110"
	}
}

// marshal encodes a message or returns nil if encoding fails.
func (s *Server) marshal(m *iso8583.Message) []byte {
	raw, err := m.Marshal()
	if err != nil {
		s.log.Error("failed to marshal response", "err", err)
		return nil
	}
	return raw
}

func idemKey(req *iso8583.Message) string {
	return req.Get(7) + "|" + req.Get(11) + "|" + req.Get(100)
}

func (s *Server) audit(req, resp *iso8583.Message) {
	if s.cfg.Audit == nil {
		return
	}
	ev := AuditEvent{
		MTI:          resp.MTI,
		STAN:         req.Get(11),
		PAN:          maskPAN(req.Get(2)),
		Amount:       req.Get(4),
		ResponseCode: resp.Get(39),
		Destination:  req.Get(100),
		CreatedAt:    time.Now().UTC(),
	}
	go func() {
		if err := s.cfg.Audit.Record(context.Background(), ev); err != nil {
			s.log.Warn("audit write failed", "err", err)
		}
	}()
}

func maskPAN(pan string) string {
	if len(pan) < 10 {
		return pan
	}
	return pan[:6] + "******" + pan[len(pan)-4:]
}

// binOf returns the first 6 digits of a PAN, or "invalid" for diagnostics.
func binOf(pan string) string {
	if bin, ok := binrouting.BINOf(pan); ok {
		return bin
	}
	return "invalid"
}
