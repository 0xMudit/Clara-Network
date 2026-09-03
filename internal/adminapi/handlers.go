package adminapi

import (
	"fmt"
	"net/http"
)

// ── Dashboard ────────────────────────────────────────────────────────────────

func (s *server) getDashboard(w http.ResponseWriter, r *http.Request) {
	summary, err := s.store.DashboardSummary(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// ── Dashboard time series ────────────────────────────────────────────────────

func (s *server) getDashboardSeries(w http.ResponseWriter, r *http.Request) {
	days := queryInt(r, "days", 14)
	series, err := s.store.TransactionSeries(r.Context(), days)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": series})
}

// ── Transactions ─────────────────────────────────────────────────────────────

func (s *server) getTransactions(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.Transactions(r.Context(), limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

// ── Clearing ─────────────────────────────────────────────────────────────────

func (s *server) getClearingCycles(w http.ResponseWriter, r *http.Request) {
	cycles, err := s.store.ClearingCycles(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string][]string{"items": cycles})
}

func (s *server) getClearingRecords(w http.ResponseWriter, r *http.Request) {
	cycleID := queryString(r, "cycle")
	if cycleID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("missing ?cycle= parameter"))
		return
	}
	records, err := s.store.ClearingRecords(r.Context(), cycleID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": records})
}

func (s *server) getNetPositions(w http.ResponseWriter, r *http.Request) {
	cycleID := queryString(r, "cycle")
	if cycleID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("missing ?cycle= parameter"))
		return
	}
	positions, err := s.store.NetPositions(r.Context(), cycleID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": positions})
}

// ── Settlement ───────────────────────────────────────────────────────────────

func (s *server) getSettlementInstructions(w http.ResponseWriter, r *http.Request) {
	cycleID := queryString(r, "cycle")
	if cycleID == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("missing ?cycle= parameter"))
		return
	}
	instructions, err := s.store.SettlementInstructions(r.Context(), cycleID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": instructions})
}

func (s *server) getPrefundAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.PrefundAccounts(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": accounts})
}

func (s *server) getDefaultFund(w http.ResponseWriter, r *http.Request) {
	balance, err := s.store.DefaultFundBalance(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"balance": balance})
}

// ── Ledger ───────────────────────────────────────────────────────────────────

func (s *server) getLedgerAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.LedgerAccountsWithBalances(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": accounts})
}

func (s *server) getLedgerEntries(w http.ResponseWriter, r *http.Request) {
	accountID := queryString(r, "account")
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.LedgerEntries(r.Context(), accountID, limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

// ── Cards ────────────────────────────────────────────────────────────────────

func (s *server) getCards(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.Cards(r.Context(), limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *server) getCard(w http.ResponseWriter, r *http.Request) {
	ref := r.PathValue("ref")
	card, ok, err := s.store.Card(r.Context(), ref)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		writeErr(w, http.StatusNotFound, fmt.Errorf("unknown card %s", ref))
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (s *server) getBinRanges(w http.ResponseWriter, r *http.Request) {
	ranges, err := s.store.BinRanges(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": ranges})
}

func (s *server) getTokens(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.Tokens(r.Context(), limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

// ── Merchants ────────────────────────────────────────────────────────────────

func (s *server) getMerchants(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.Merchants(r.Context(), limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *server) getMerchant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	merchant, ok, err := s.store.Merchant(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		writeErr(w, http.StatusNotFound, fmt.Errorf("unknown merchant %s", id))
		return
	}
	writeJSON(w, http.StatusOK, merchant)
}

func (s *server) getMerchantFunding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	lines, err := s.store.MerchantFunding(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": lines})
}

// ── Disputes ─────────────────────────────────────────────────────────────────

func (s *server) getDisputes(w http.ResponseWriter, r *http.Request) {
	stage := queryString(r, "stage")
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	page, err := s.store.Disputes(r.Context(), stage, limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *server) getOverdueDisputes(w http.ResponseWriter, r *http.Request) {
	disputes, err := s.store.OverdueDisputes(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": disputes})
}

func (s *server) getDispute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	dispute, ok, err := s.store.Dispute(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		writeErr(w, http.StatusNotFound, fmt.Errorf("unknown dispute %s", id))
		return
	}
	writeJSON(w, http.StatusOK, dispute)
}

func (s *server) getDisputeRatio(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ratio, status, err := s.store.DisputeRatio(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"merchantId": id,
		"ratio":      ratio,
		"status":     status,
	})
}
