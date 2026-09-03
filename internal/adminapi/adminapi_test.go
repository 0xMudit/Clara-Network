package adminapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("CLARA_PG_DSN")
	if dsn == "" {
		t.Skip("CLARA_PG_DSN not set; skipping integration test")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool: %v", err)
	}
	t.Cleanup(func() { pool.Close() })
	return pool
}

func newTestServer(t *testing.T, pool *pgxpool.Pool) *httptest.Server {
	t.Helper()
	store := &Store{Pool: pool}
	s := &server{store: store, log: testLogger()}
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

	srv := httptest.NewServer(cors(mux))
	t.Cleanup(func() { srv.Close() })
	return srv
}

func getJSON(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return resp
}

func decodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	return v
}

func TestHealth(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/health")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[map[string]string](t, resp)
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestDashboard(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/dashboard")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[DashboardSummary](t, resp)
	if body.Cards < 0 {
		t.Fatal("cards count should be non-negative")
	}
}

func TestDashboardSeries(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/dashboard/series?days=30")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[struct {
		Items []SeriesPoint `json:"items"`
	}](t, resp)
	if len(body.Items) != 30 {
		t.Fatalf("expected 30 bucket points, got %d", len(body.Items))
	}
	for _, p := range body.Items {
		if p.Count < 0 {
			t.Fatalf("count should be non-negative, got %d (%s)", p.Count, p.Date)
		}
	}
}

func TestTransactions(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/transactions?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[Page[AuditEvent]](t, resp)
	if body.Total < 0 {
		t.Fatal("total should be non-negative")
	}
}

func TestClearingCycles(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/clearing/cycles")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestClearingRecordsRequiresCycle(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/clearing/records")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 without ?cycle=, got %d", resp.StatusCode)
	}
}

func TestClearingRecords(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/clearing/records?cycle=test-cycle")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNetPositions(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/clearing/positions?cycle=test-cycle")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSettlementInstructionsRequiresCycle(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/settlement/instructions")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 without ?cycle=, got %d", resp.StatusCode)
	}
}

func TestPrefundAccounts(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/settlement/prefunds")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDefaultFund(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/settlement/default-fund")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[map[string]int64](t, resp)
	if body["balance"] < 0 {
		t.Fatal("balance should be non-negative")
	}
}

func TestLedgerAccounts(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/ledger/accounts")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLedgerEntries(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/ledger/entries?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[Page[LedgerEntry]](t, resp)
	if body.Total < 0 {
		t.Fatal("total should be non-negative")
	}
}

func TestCards(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/cards?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[Page[CardRecord]](t, resp)
	if body.Total < 0 {
		t.Fatal("total should be non-negative")
	}
}

func TestCardNotFound(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/cards/nonexistent")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestBinRanges(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/bin-ranges")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestTokens(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/tokens?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestMerchants(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/merchants?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[Page[MerchantRecord]](t, resp)
	if body.Total < 0 {
		t.Fatal("total should be non-negative")
	}
}

func TestMerchantNotFound(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/merchants/nonexistent")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestMerchantFunding(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/merchants/M-test/funding")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDisputes(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/disputes?limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[Page[DisputeRecord]](t, resp)
	if body.Total < 0 {
		t.Fatal("total should be non-negative")
	}
}

func TestDisputesByStage(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/disputes?stage=filed&limit=5")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestOverdueDisputes(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/disputes/overdue")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDisputeNotFound(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/disputes/nonexistent")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDisputeRatio(t *testing.T) {
	pool := testDB(t)
	srv := newTestServer(t, pool)

	resp := getJSON(t, srv.URL+"/api/v1/disputes/M-test/ratio")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
