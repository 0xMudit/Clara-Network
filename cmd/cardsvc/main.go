// Command cardsvc runs the issuing-stack REST service: card issuance, EMV
// cryptogram verification, tokenization, and wallet provisioning.
package main

import (
	"context"
	"encoding/hex"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/cardsvc"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	masterKey, err := hex.DecodeString(env.Get("CLARA_ISSUER_MASTER_KEY", "2b7e151628aed2a6abf7158809cf4f3c"))
	if err != nil {
		logger.Error("invalid issuer master key", "err", err)
		os.Exit(1)
	}

	var store cardsvc.Store = cardsvc.NewMemoryStore()
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		if pool, err := pgxpool.New(ctx, dsn); err == nil {
			store = &cardsvc.PostgresStore{Pool: pool}
		} else {
			logger.Warn("postgres unavailable, using in-memory store", "err", err)
		}
	}

	svc, err := cardsvc.NewService(store, cardsvc.Config{IssuerMasterKey: masterKey, Log: logger})
	if err != nil {
		logger.Error("cardsvc init failed", "err", err)
		os.Exit(1)
	}
	vault := cardsvc.NewTokenVault(store)

	// Seed the issued BIN range unless one already exists.
	if err := svc.AddBinRange(ctx, cardsvc.BinRange{
		BIN: env.Get("CLARA_BIN", "400000"), Low: 0, High: 9999999999,
		Currency: "840", Product: "classic",
	}); err != nil {
		logger.Error("seeding bin range failed", "err", err)
		os.Exit(1)
	}

	addr := env.Get("CLARA_LISTEN", ":8081")
	if err := cardsvc.ListenAndServe(ctx, addr, svc, vault, logger); err != nil {
		logger.Error("cardsvc exited", "err", err)
		os.Exit(1)
	}
}
