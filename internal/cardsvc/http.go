package cardsvc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

// httpServer exposes the issuing stack as a REST API for members, wallets,
// and the card-sim harness. It never returns card keys or clear PANs.
type httpServer struct {
	svc   *Service
	vault *TokenVault
	log   *slog.Logger
}

// ListenAndServe starts the cardsvc HTTP server.
func ListenAndServe(ctx context.Context, addr string, svc *Service, vault *TokenVault, log *slog.Logger) error {
	s := &httpServer{svc: svc, vault: vault, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /cards", s.createCard)
	mux.HandleFunc("POST /cards/{ref}/arqc", s.computeARQC)
	mux.HandleFunc("POST /cards/{ref}/verify-arqc", s.verifyARQC)
	mux.HandleFunc("POST /tokens", s.tokenize)
	mux.HandleFunc("GET /tokens/{token}", s.getToken)
	mux.HandleFunc("POST /tokens/{token}/provision", s.provision)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cardsvc: listen: %w", err)
	}
	log.Info("cardsvc listening", "addr", ln.Addr())
	srv := &http.Server{Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.Serve(ln); err != nil && ctx.Err() == nil {
		return fmt.Errorf("cardsvc: serve: %w", err)
	}
	return nil
}

func (s *httpServer) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createCardReq struct {
	PAN     string `json:"pan"`
	Expiry  string `json:"expiry"`
	Product string `json:"product"`
}

type cardResp struct {
	Ref     string `json:"ref"`
	PANMask string `json:"panMask"`
	BIN     string `json:"bin"`
	Expiry  string `json:"expiry"`
	Status  string `json:"status"`
	Product string `json:"product"`
}

func (s *httpServer) createCard(w http.ResponseWriter, r *http.Request) {
	var req createCardReq
	if !decode(w, r, &req) {
		return
	}
	card, err := s.svc.CreateCard(r.Context(), req.PAN, req.Expiry, req.Product)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, cardResp{
		Ref: card.Ref, PANMask: card.PANMask, BIN: card.BIN,
		Expiry: card.Expiry, Status: card.Status, Product: card.Product,
	})
}

type arqcReq struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	STAN     string `json:"stan"`
	Date     string `json:"date"`
	ATC      uint16 `json:"atc"`
	UN       string `json:"un"`
	ARQC     string `json:"arqc"` // hex, for verification
}

func (s *httpServer) computeARQC(w http.ResponseWriter, r *http.Request) {
	var req arqcReq
	if !decode(w, r, &req) {
		return
	}
	arqc, err := s.svc.ComputeARQC(r.Context(), r.PathValue("ref"), req.data())
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"arqc": hex.EncodeToString(arqc)})
}

func (s *httpServer) verifyARQC(w http.ResponseWriter, r *http.Request) {
	var req arqcReq
	if !decode(w, r, &req) {
		return
	}
	raw, err := hex.DecodeString(req.ARQC)
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("arqc must be hex: %w", err))
		return
	}
	valid, err := s.svc.VerifyARQC(r.Context(), r.PathValue("ref"), req.data(), raw)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"valid": valid})
}

type tokenizeReq struct {
	PAN string `json:"pan"`
}

type tokenResp struct {
	Token    string `json:"token"`
	PAR      string `json:"par"`
	BIN      string `json:"bin"`
	Status   string `json:"status"`
	DeviceID string `json:"deviceId,omitempty"`
}

func (s *httpServer) tokenize(w http.ResponseWriter, r *http.Request) {
	var req tokenizeReq
	if !decode(w, r, &req) {
		return
	}
	tok, err := s.vault.Tokenize(r.Context(), s.svc, req.PAN)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, tokenResp{Token: tok.Number, PAR: tok.PAR, BIN: tok.BIN, Status: tok.Status})
}

func (s *httpServer) getToken(w http.ResponseWriter, r *http.Request) {
	tok, ok, err := s.vault.store.Token(r.Context(), r.PathValue("token"))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		writeErr(w, http.StatusNotFound, fmt.Errorf("unknown token"))
		return
	}
	writeJSON(w, http.StatusOK, tokenResp{Token: tok.Number, PAR: tok.PAR, BIN: tok.BIN, Status: tok.Status, DeviceID: tok.DeviceID})
}

type provisionReq struct {
	DeviceID  string `json:"deviceId"`
	Requestor string `json:"requestor"`
}

func (s *httpServer) provision(w http.ResponseWriter, r *http.Request) {
	var req provisionReq
	if !decode(w, r, &req) {
		return
	}
	p, err := s.vault.Provision(r.Context(), r.PathValue("token"), req.DeviceID, req.Requestor)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (q arqcReq) data() ARQCData {
	return ARQCData{
		Amount: q.Amount, Currency: q.Currency, STAN: q.STAN,
		Date: q.Date, ATC: q.ATC, UN: q.UN,
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("bad json: %w", err))
		return false
	}
	return true
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}
