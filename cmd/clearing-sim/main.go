// Command clearing-sim exercises the Clara Network clearing and net
// settlement layer: it seeds sample clearing records and prefunded accounts,
// runs one settlement cycle, prints the net positions and instructions, and
// writes the ISO 20022 pacs.009 settlement instructions as XML files.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/clearing"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	acqA = "ACQ-A"
	acqB = "ACQ-B"
	issC = "ISS-C"
	issD = "ISS-D"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cycleID := env.Get("CLARA_CYCLE", time.Now().UTC().Format("20060102"))
	scenario := env.Get("CLARA_SCENARIO", "default")
	outDir := env.Get("CLARA_OUT", "out/clearing")

	var store clearing.Store = clearing.NewMemoryStore()
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logger.Error("postgres unavailable, using in-memory store", "err", err)
		} else {
			store = &clearing.PostgresStore{Pool: pool}
		}
	}

	svc := clearing.NewService(store, clearing.Config{})

	if err := seed(ctx, svc, cycleID, scenario); err != nil {
		logger.Error("seeding failed", "err", err)
		os.Exit(1)
	}

	res, err := svc.RunCycle(ctx, cycleID, time.Now())
	if err != nil {
		logger.Error("settlement cycle failed", "err", err)
		os.Exit(1)
	}

	printResult(logger, res)
	if err := writeInstructions(outDir, res); err != nil {
		logger.Error("writing pacs.009 failed", "err", err)
		os.Exit(1)
	}
	logger.Info("settlement instructions written", "dir", outDir)
}

// seed loads clearing records, prefunded accounts, and the default fund.
func seed(ctx context.Context, svc *clearing.Service, cycleID, scenario string) error {
	records := []clearing.ClearingRecord{
		rec(acqA, issC, 2500, 300),
		rec(acqA, issC, 1800, 200),
		rec(acqA, issC, 5000, 600),
		rec(acqA, issD, 4200, 400),
		rec(acqB, issC, 9000, 1000),
		rec(acqB, issD, 1500, 150),
	}
	if err := svc.SubmitBatch(ctx, cycleID, records); err != nil {
		return fmt.Errorf("submit batch: %w", err)
	}

	// In the "default" scenario the largest acquirer pre-funds far below its
	// net obligation, forcing the default fund to step in.
	acqAPrefund := int64(20000)
	if scenario == "default" {
		acqAPrefund = 3000
	}
	prefunds := []clearing.PrefundAccount{
		{Member: acqA, Balance: acqAPrefund, Cap: 20000},
		{Member: acqB, Balance: 15000, Cap: 15000},
		{Member: issC, Balance: 1000, Cap: 1000},
		{Member: issD, Balance: 1000, Cap: 1000},
	}
	for _, a := range prefunds {
		if err := svc.Fund(ctx, a); err != nil {
			return fmt.Errorf("fund %s: %w", a.Member, err)
		}
	}
	if err := svc.FundDefaultFund(ctx, 50000); err != nil {
		return fmt.Errorf("fund default fund: %w", err)
	}
	return nil
}

func rec(sender, receiver string, amount, interchange int64) clearing.ClearingRecord {
	return clearing.ClearingRecord{
		STAN:        "100000",
		MTI:         "0221",
		Sender:      sender,
		Receiver:    receiver,
		AmountMinor: amount,
		Interchange: interchange,
		Currency:    "840",
		RefID:       fmt.Sprintf("%s-%d", sender, amount),
	}
}

func printResult(logger *slog.Logger, res *clearing.CycleResult) {
	logger.Info("settlement cycle", "cycle", res.CycleID, "final", res.Final,
		"default_fund_balance", res.DefaultFundBalance, "default_fund_target", res.DefaultFundTarget)
	for _, p := range res.Positions {
		logger.Info("net position", "member", p.Member, "net", clearing.FormatAmount(p.Net))
	}
	for _, in := range res.Instructions {
		logger.Info("settlement instruction", "member", in.Member, "direction", in.Direction,
			"amount", clearing.FormatAmount(in.Amount), "msg_id", in.MsgID)
	}
	for _, ev := range res.Events {
		logger.Warn("default event", "member", ev.Member, "shortfall", clearing.FormatAmount(ev.Shortfall),
			"covered", clearing.FormatAmount(ev.Covered), "uncovered", clearing.FormatAmount(ev.Uncovered))
	}
	for m, a := range res.Accounts {
		logger.Info("prefund account", "member", m, "balance", clearing.FormatAmount(a.Balance), "cap", clearing.FormatAmount(a.Cap))
	}
}

func writeInstructions(outDir string, res *clearing.CycleResult) error {
	for _, in := range res.Instructions {
		xml, err := clearing.Pacs009XML(in)
		if err != nil {
			return err
		}
		dir := filepath.Join(outDir, in.CycleID)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		path := filepath.Join(dir, in.Member+".pacs.009.xml")
		if err := os.WriteFile(path, xml, 0o644); err != nil {
			return err
		}
	}
	return nil
}
