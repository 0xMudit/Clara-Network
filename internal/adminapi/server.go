package adminapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/0xMudit/Clara-Network/internal/metrics"
)

// server is the admin API HTTP handler.
type server struct {
	store   *Store
	log     *slog.Logger
	metrics *metrics.Registry
}

// ListenAndServe starts the admin REST API.
func ListenAndServe(ctx context.Context, addr string, store *Store, log *slog.Logger) error {
	return ListenAndServeWithMetrics(ctx, addr, store, log, nil)
}

// ListenAndServeWithMetrics starts the admin REST API with an optional
// Prometheus-compatible /metrics endpoint.
func ListenAndServeWithMetrics(ctx context.Context, addr string, store *Store, log *slog.Logger, reg *metrics.Registry) error {
	s := &server{store: store, log: log, metrics: reg}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/v1/dashboard", s.getDashboard)
	mux.HandleFunc("GET /api/v1/dashboard/series", s.getDashboardSeries)
	mux.HandleFunc("GET /api/v1/transactions", s.getTransactions)
	mux.HandleFunc("GET /api/v1/clearing/cycles", s.getClearingCycles)
	mux.HandleFunc("GET /api/v1/clearing/records", s.getClearingRecords)
	mux.HandleFunc("GET /api/v1/clearing/positions", s.getNetPositions)
	mux.HandleFunc("GET /api/v1/settlement/instructions", s.getSettlementInstructions)
	mux.HandleFunc("GET /api/v1/settlement/prefunds", s.getPrefundAccounts)
	mux.HandleFunc("GET /api/v1/settlement/default-fund", s.getDefaultFund)
	mux.HandleFunc("GET /api/v1/ledger/accounts", s.getLedgerAccounts)
	mux.HandleFunc("GET /api/v1/ledger/entries", s.getLedgerEntries)
	mux.HandleFunc("GET /api/v1/cards", s.getCards)
	mux.HandleFunc("GET /api/v1/cards/{ref}", s.getCard)
	mux.HandleFunc("GET /api/v1/bin-ranges", s.getBinRanges)
	mux.HandleFunc("GET /api/v1/tokens", s.getTokens)
	mux.HandleFunc("GET /api/v1/merchants", s.getMerchants)
	mux.HandleFunc("GET /api/v1/merchants/{id}", s.getMerchant)
	mux.HandleFunc("GET /api/v1/merchants/{id}/funding", s.getMerchantFunding)
	mux.HandleFunc("GET /api/v1/disputes", s.getDisputes)
	mux.HandleFunc("GET /api/v1/disputes/overdue", s.getOverdueDisputes)
	mux.HandleFunc("GET /api/v1/disputes/{id}", s.getDispute)
	mux.HandleFunc("GET /api/v1/disputes/{id}/ratio", s.getDisputeRatio)

	if reg != nil {
		mux.Handle("GET /metrics", reg.Handler())
	}

	handler := cors(mux)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("adminapi: listen: %w", err)
	}
	log.Info("adminapi listening", "addr", ln.Addr())
	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.Serve(ln); err != nil && ctx.Err() == nil {
		return fmt.Errorf("adminapi: serve: %w", err)
	}
	return nil
}

// cors wraps an handler with permissive CORS headers for development.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ── query helpers ────────────────────────────────────────────────────────────

func queryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func queryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// ── JSON helpers (matches cardsvc convention) ────────────────────────────────

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
