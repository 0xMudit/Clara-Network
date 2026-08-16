// Command switch runs the Clara Network ISO 8583 message switch.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/binrouting"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/0xMudit/Clara-Network/internal/risk"
	"github.com/0xMudit/Clara-Network/internal/switchsrv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := switchsrv.Config{
		ListenAddr:     env.Get("CLARA_LISTEN", ":8080"),
		IssuerRoutes:   map[string]string{},
		IdempotencyTTL: 60 * time.Second,
		Log:            logger,
	}

	if raw := os.Getenv("CLARA_ISSUER_ROUTES"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &cfg.IssuerRoutes); err != nil {
			logger.Error("invalid CLARA_ISSUER_ROUTES", "err", err)
			os.Exit(1)
		}
	}

	cfg.Idempotency = switchsrv.NewMemoryIdempotency()
	if addr := os.Getenv("CLARA_REDIS_ADDR"); addr != "" {
		client := redis.NewClient(&redis.Options{Addr: addr})
		if err := client.Ping(ctx).Err(); err != nil {
			logger.Warn("redis unavailable, using in-memory idempotency", "err", err)
		} else {
			cfg.Idempotency = &switchsrv.RedisIdempotency{Client: client}
		}
	}

	if raw := os.Getenv("CLARA_BIN_TABLE"); raw != "" {
		tab, err := binrouting.FromJSON([]byte(raw))
		if err != nil {
			logger.Error("invalid CLARA_BIN_TABLE", "err", err)
			os.Exit(1)
		}
		cfg.BINTable = tab
	}

	if raw := os.Getenv("CLARA_RISK_RULES"); raw != "" {
		var store risk.Store = risk.NewMemoryStore()
		if addr := os.Getenv("CLARA_REDIS_ADDR"); addr != "" {
			store = risk.NewRedisStore(addr)
		}
		engine, err := risk.FromConfig([]byte(raw), store)
		if err != nil {
			logger.Error("invalid CLARA_RISK_RULES", "err", err)
			os.Exit(1)
		}
		cfg.Risk = engine
	}

	cfg.Audit = switchsrv.NoopAudit{}
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logger.Warn("postgres unavailable, audit disabled", "err", err)
		} else {
			cfg.Audit = &switchsrv.PostgresAudit{Pool: pool}
		}
	}

	srv, err := switchsrv.New(cfg)
	if err != nil {
		logger.Error("failed to start switch", "err", err)
		os.Exit(1)
	}
	if err := srv.ListenAndServe(ctx); err != nil {
		logger.Error("switch exited", "err", err)
		os.Exit(1)
	}
}
