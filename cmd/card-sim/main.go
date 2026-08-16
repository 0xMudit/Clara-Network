// Command card-sim exercises the Clara Network issuing stack end to end: it
// personalizes a card, verifies an EMV-style online cryptogram (valid,
// tampered, replayed), tokenizes the PAN, detokenizes it, and provisions a
// mobile wallet (docs/25 §25.4 phase 5).
package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/cardsvc"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const masterKeyHex = "2b7e151628aed2a6abf7158809cf4f3c"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	masterKey, err := hex.DecodeString(env.Get("CLARA_ISSUER_MASTER_KEY", masterKeyHex))
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

	if err := svc.AddBinRange(ctx, cardsvc.BinRange{BIN: "400000", Low: 0, High: 9999999999, Currency: "840", Product: "classic"}); err != nil {
		logger.Error("seeding bin range failed", "err", err)
		os.Exit(1)
	}

	pan := env.Get("CLARA_PAN", "4000001234567899")
	card, err := svc.CreateCard(ctx, pan, "3012", env.Get("CLARA_PRODUCT", "platinum"))
	if err != nil {
		logger.Error("card issuance failed", "err", err)
		os.Exit(1)
	}
	logger.Info("card issued", "mask", card.PANMask, "bin", card.BIN, "product", card.Product, "status", card.Status)

	// --- EMV cryptogram verification ---
	txn := cardsvc.ARQCData{Amount: 250000, Currency: "840", STAN: "100000", Date: "260816", ATC: 7, UN: "1234"}
	arqc, err := svc.ComputeARQC(ctx, card.Ref, txn)
	if err != nil {
		logger.Error("arqc computation failed", "err", err)
		os.Exit(1)
	}
	logger.Info("arqc computed", "hex", hex.EncodeToString(arqc), "data", fmt.Sprintf("%x", txn.Bytes()))

	valid, err := svc.VerifyARQC(ctx, card.Ref, txn, arqc)
	if err != nil || !valid {
		logger.Error("arqc verification failed", "err", err, "valid", valid)
		os.Exit(1)
	}
	logger.Info("arqc verified", "valid", true)

	tampered := txn
	tampered.Amount = 1
	if valid, _ := svc.VerifyARQC(ctx, card.Ref, tampered, arqc); valid {
		logger.Error("tampered arqc must not verify")
		os.Exit(1)
	}
	logger.Info("tampered arqc rejected", "valid", false)

	if valid, _ := svc.VerifyARQC(ctx, card.Ref, txn, arqc); valid {
		logger.Error("replayed arqc must be rejected")
		os.Exit(1)
	}
	logger.Info("replayed arqc rejected", "valid", false)

	// --- Token vault ---
	tok, err := vault.Tokenize(ctx, svc, pan)
	if err != nil {
		logger.Error("tokenization failed", "err", err)
		os.Exit(1)
	}
	logger.Info("tokenized", "token", tok.Number, "par", tok.PAR, "bin", tok.BIN, "luhn", cardsvc.ValidLuhn(tok.Number))

	hash, err := vault.Detokenize(ctx, tok.Number)
	if err != nil {
		logger.Error("detokenization failed", "err", err)
		os.Exit(1)
	}
	if !bytes.Equal(hash, card.PANHash) {
		logger.Error("detokenized PAN hash does not match the card")
		os.Exit(1)
	}
	logger.Info("detokenized", "pan_hash_matches", true)

	// --- Mobile wallet provisioning ---
	p, err := vault.Provision(ctx, tok.Number, env.Get("CLARA_DEVICE_ID", "device-42"), env.Get("CLARA_TRID", "TRID001"))
	if err != nil {
		logger.Error("provisioning failed", "err", err)
		os.Exit(1)
	}
	logger.Info("provisioned wallet", "token", p.Token, "par", p.PAR, "device", p.DeviceID, "requestor", p.Requestor)
	logger.Info("card-sim done")
}
