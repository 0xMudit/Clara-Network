// Package issuersim implements a simulated issuer host used for development
// and integration testing. It applies a configurable decision function to
// authorization requests and responds in ISO 8583.
package issuersim

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
)

// DecisionFunc maps an authorization request to a response code.
type DecisionFunc func(req *iso8583.Message) string

// Config configures the simulated issuer host.
type Config struct {
	ListenAddr string
	ID         string
	Decision   DecisionFunc
	Log        *slog.Logger
}

// Server is a simulated issuer host.
type Server struct {
	cfg Config
	ln  net.Listener
	log *slog.Logger
	seq atomic.Uint64
}

// New binds the issuer listener and returns a ready server.
func New(cfg Config) (*Server, error) {
	if cfg.Decision == nil {
		cfg.Decision = DefaultDecision
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	ln, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("issuersim: listen: %w", err)
	}
	return &Server{cfg: cfg, ln: ln, log: cfg.Log}, nil
}

// Addr returns the bound listener address (useful with port 0 in tests).
func (s *Server) Addr() net.Addr { return s.ln.Addr() }

// ListenAndServe accepts connections until ctx is cancelled.
func (s *Server) ListenAndServe(ctx context.Context) error {
	s.log.Info("issuer-sim listening", "addr", s.ln.Addr(), "id", s.cfg.ID)
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
			return fmt.Errorf("issuersim: accept: %w", err)
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
		req, err := iso8583.Parse(frame)
		if err != nil {
			s.log.Warn("bad request", "err", err)
			continue
		}
		resp := s.respond(req)
		if err := framing.WriteFrame(conn, mustMarshal(resp)); err != nil {
			return
		}
	}
}

func (s *Server) respond(req *iso8583.Message) *iso8583.Message {
	resp := iso8583.New(responseMTI(req.MTI))
	resp.Set(7, req.Get(7))
	resp.Set(11, req.Get(11))
	resp.Set(12, req.Get(12))
	resp.Set(13, req.Get(13))
	resp.Set(32, req.Get(32))
	resp.Set(37, s.rrn())
	resp.Set(39, s.cfg.Decision(req))
	resp.Set(100, req.Get(100))
	return resp
}

func (s *Server) rrn() string {
	n := s.seq.Add(1)
	return fmt.Sprintf("%06d%06d", time.Now().UTC().Unix()%1000000, n%1000000)
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

func mustMarshal(m *iso8583.Message) []byte {
	raw, err := m.Marshal()
	if err != nil {
		panic(err)
	}
	return raw
}

// DefaultDecision approves most requests: PANs ending in "0000" are declined
// and amounts above 5000.00 are declined for insufficient funds.
func DefaultDecision(req *iso8583.Message) string {
	if strings.HasSuffix(req.Get(2), "0000") {
		return "05" // do not honor
	}
	amount, err := strconv.ParseInt(req.Get(4), 10, 64)
	if err == nil && amount > 500000 {
		return "51" // insufficient funds
	}
	return "00"
}
