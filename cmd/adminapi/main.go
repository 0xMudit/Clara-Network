// Command adminapi runs the Clara Network admin REST API: a read-only
// aggregator that serves dashboard data from the shared PostgreSQL schema.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/adminapi"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/0xMudit/Clara-Network/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("CLARA_PG_DSN")
	if dsn == "" {
		logger.Error("CLARA_PG_DSN is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("postgres connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	store := &adminapi.Store{Pool: pool}
	addr := env.Get("CLARA_LISTEN", ":8083")

	reg := metrics.NewRegistry()

	if err := adminapi.ListenAndServeWithMetrics(ctx, addr, store, logger, reg); err != nil {
		logger.Error("adminapi exited", "err", err)
		os.Exit(1)
	}
}
