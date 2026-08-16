// Command issuer-sim runs a simulated issuer host for development.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/0xMudit/Clara-Network/internal/issuersim"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := issuersim.Config{
		ListenAddr: env.Get("CLARA_LISTEN", ":8082"),
		ID:         env.Get("CLARA_ID", ""),
		Log:        logger,
	}
	srv, err := issuersim.New(cfg)
	if err != nil {
		logger.Error("failed to start issuer-sim", "err", err)
		os.Exit(1)
	}
	if err := srv.ListenAndServe(ctx); err != nil {
		logger.Error("issuer-sim exited", "err", err)
		os.Exit(1)
	}
}
