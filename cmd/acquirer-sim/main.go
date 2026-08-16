// Command acquirer-sim sends a batch of ISO 8583 authorization requests to
// the Clara Network switch and reports the responses.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/acquirersim"
	"github.com/0xMudit/Clara-Network/internal/env"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := acquirersim.Config{
		SwitchAddr: env.Get("CLARA_SWITCH", "localhost:8080"),
		AcquirerID: env.Get("CLARA_ACQUIRER_ID", "1000001"),
		IssuerID:   env.Get("CLARA_ISSUER_ID", "1000001000"),
		PAN:        env.Get("CLARA_PAN", "4000001234567890"),
		Count:      env.Int("CLARA_COUNT", 3),
		Amount:     env.Int("CLARA_AMOUNT", 1000),
		Step:       env.Int("CLARA_STEP", 2500),
		SendDE100:  env.Get("CLARA_SEND_DE100", "true") != "false",
		Log:        logger,
	}
	if err := acquirersim.Run(ctx, cfg); err != nil {
		logger.Error("acquirer-sim failed", "err", err)
		os.Exit(1)
	}
}
