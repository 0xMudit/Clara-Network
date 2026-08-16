// Command ledger-sim exercises the Clara Network ledger and reconciliation
// layer: it runs a settlement cycle, posts the net positions as balanced
// double-entry journals (docs/12 §12.5), then reconciles the ledger against
// the settlement agent's statement (docs/12 §12.9).
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/clearing"
	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/0xMudit/Clara-Network/internal/ledger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const currency = "840"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cycleID := env.Get("CLARA_CYCLE", time.Now().UTC().Format("20060102"))
	scenario := env.Get("CLARA_SCENARIO", "default")

	var clr clearing.Store = clearing.NewMemoryStore()
	var ldgr ledger.Store = ledger.NewMemoryStore()
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		if pool, err := pgxpool.New(ctx, dsn); err == nil {
			clr = &clearing.PostgresStore{Pool: pool}
			ldgr = &ledger.PostgresStore{Pool: pool}
		} else {
			logger.Warn("postgres unavailable, using in-memory stores", "err", err)
		}
	}
	defer ldgr.Close()

	svc := clearing.NewService(clr, clearing.Config{})
	book := ledger.NewLedger(ldgr)

	if err := seedClearing(ctx, svc, cycleID, scenario); err != nil {
		logger.Error("seeding clearing failed", "err", err)
		os.Exit(1)
	}

	res, err := svc.RunCycle(ctx, cycleID, time.Now())
	if err != nil {
		logger.Error("settlement cycle failed", "err", err)
		os.Exit(1)
	}

	if err := seedAccounts(ctx, book, res); err != nil {
		logger.Error("seeding ledger accounts failed", "err", err)
		os.Exit(1)
	}
	if err := postPositions(ctx, book, res); err != nil {
		logger.Error("posting net positions failed", "err", err)
		os.Exit(1)
	}

	statements := statementsFrom(res)
	if os.Getenv("CLARA_MISMATCH") != "" {
		statements = injectMismatch(statements)
		logger.Warn("injected settlement statement mismatch for demo")
	}

	report, err := (&ledger.Reconciler{
		Store:               ldgr,
		MemberAccountPrefix: "M:",
		SchemeOperatorID:    clearing.SchemeOperatorID,
	}).Reconcile(ctx, statements)
	if err != nil {
		logger.Error("reconciliation failed", "err", err)
		os.Exit(1)
	}

	printSettlement(logger, res)
	printLedger(ctx, logger, book)
	printReport(logger, report)

	if report.OK() {
		logger.Info("reconciliation: OK")
	} else {
		logger.Warn("reconciliation: MISMATCHES FOUND")
	}
}

func seedClearing(ctx context.Context, svc *clearing.Service, cycleID, scenario string) error {
	records := []clearing.ClearingRecord{
		clearingRecord("ACQ-A", "ISS-C", 2500, 300),
		clearingRecord("ACQ-A", "ISS-C", 1800, 200),
		clearingRecord("ACQ-A", "ISS-C", 5000, 600),
		clearingRecord("ACQ-A", "ISS-D", 4200, 400),
		clearingRecord("ACQ-B", "ISS-C", 9000, 1000),
		clearingRecord("ACQ-B", "ISS-D", 1500, 150),
	}
	if err := svc.SubmitBatch(ctx, cycleID, records); err != nil {
		return err
	}

	acqAPrefund := int64(20000)
	if scenario == "default" {
		acqAPrefund = 3000
	}
	for _, a := range []clearing.PrefundAccount{
		{Member: "ACQ-A", Balance: acqAPrefund, Cap: 20000},
		{Member: "ACQ-B", Balance: 15000, Cap: 15000},
		{Member: "ISS-C", Balance: 1000, Cap: 1000},
		{Member: "ISS-D", Balance: 1000, Cap: 1000},
	} {
		if err := svc.Fund(ctx, a); err != nil {
			return err
		}
	}
	return svc.FundDefaultFund(ctx, 50000)
}

func clearingRecord(sender, receiver string, amount, interchange int64) clearing.ClearingRecord {
	return clearing.ClearingRecord{
		STAN:        "100000",
		MTI:         "0221",
		Sender:      sender,
		Receiver:    receiver,
		AmountMinor: amount,
		Interchange: interchange,
		Currency:    currency,
		RefID:       fmt.Sprintf("%s-%d", sender, amount),
	}
}

// seedAccounts creates the scheme's chart of accounts: a cash account, member
// liability (deposit) accounts, and the fee income account.
func seedAccounts(ctx context.Context, book *ledger.Ledger, res *clearing.CycleResult) error {
	if err := book.EnsureAccount(ctx, ledger.Account{ID: "CASH", Type: ledger.AccountAsset, Currency: currency}); err != nil {
		return err
	}
	if err := book.EnsureAccount(ctx, ledger.Account{ID: "INCOME:FEES", Type: ledger.AccountIncome, Currency: currency}); err != nil {
		return err
	}
	for _, p := range res.Positions {
		if p.Member == clearing.SchemeOperatorID {
			continue
		}
		if err := book.EnsureAccount(ctx, ledger.Account{
			ID: "M:" + p.Member, Type: ledger.AccountLiability, Currency: currency,
		}); err != nil {
			return err
		}
	}
	return nil
}

// postPositions realizes each net position as a balanced journal (docs/18):
// net debtors' deposits are reduced against cash, net creditors' deposits are
// increased, and the scheme operator's net fee position is booked as income.
func postPositions(ctx context.Context, book *ledger.Ledger, res *clearing.CycleResult) error {
	for _, p := range res.Positions {
		if p.Net == 0 {
			continue
		}
		ref := res.CycleID + ":" + p.Member
		amount := p.Net
		if amount < 0 {
			amount = -amount
		}
		if p.Member == clearing.SchemeOperatorID {
			// The operator never moves funds through the settlement agent;
			// its net fee position is booked as income against cash.
			if err := book.Post(ctx, ref, []ledger.Entry{
				{AccountID: "CASH", Direction: ledger.Debit, Amount: amount, Currency: currency, Reference: ref},
				{AccountID: "INCOME:FEES", Direction: ledger.Credit, Amount: amount, Currency: currency, Reference: ref},
			}); err != nil {
				return fmt.Errorf("post operator position: %w", err)
			}
			continue
		}
		if p.Net < 0 {
			if err := book.Post(ctx, ref, []ledger.Entry{
				{AccountID: "M:" + p.Member, Direction: ledger.Debit, Amount: amount, Currency: currency, Reference: ref},
				{AccountID: "CASH", Direction: ledger.Credit, Amount: amount, Currency: currency, Reference: ref},
			}); err != nil {
				return fmt.Errorf("post debit position %s: %w", p.Member, err)
			}
			continue
		}
		if err := book.Post(ctx, ref, []ledger.Entry{
			{AccountID: "CASH", Direction: ledger.Debit, Amount: amount, Currency: currency, Reference: ref},
			{AccountID: "M:" + p.Member, Direction: ledger.Credit, Amount: amount, Currency: currency, Reference: ref},
		}); err != nil {
			return fmt.Errorf("post credit position %s: %w", p.Member, err)
		}
	}
	return nil
}

// statementsFrom derives the settlement agent's statement from the scheme's
// settlement instructions.
func statementsFrom(res *clearing.CycleResult) []ledger.Statement {
	out := make([]ledger.Statement, 0, len(res.Instructions))
	for _, in := range res.Instructions {
		dir := ledger.Credit
		if in.Direction == clearing.DirDebit {
			dir = ledger.Debit
		}
		out = append(out, ledger.Statement{
			CycleID:   in.CycleID,
			Member:    in.Member,
			Amount:    in.Amount,
			Direction: dir,
		})
	}
	return out
}

// injectMismatch corrupts the statement to demonstrate the classification of
// discrepancies: ISS-C's line differs from the ledger (amount mismatch) and
// ACQ-B's line is dropped (orphan in the ledger).
func injectMismatch(st []ledger.Statement) []ledger.Statement {
	var out []ledger.Statement
	for i, s := range st {
		switch {
		case s.Member == "ISS-C" && i < len(st):
			s.Amount += 100
			out = append(out, s)
		case s.Member == "ACQ-B":
			// drop
		default:
			out = append(out, s)
		}
	}
	return out
}

func printSettlement(logger *slog.Logger, res *clearing.CycleResult) {
	logger.Info("settlement cycle", "cycle", res.CycleID, "final", res.Final)
	for _, p := range res.Positions {
		logger.Info("net position", "member", p.Member, "net", clearing.FormatAmount(p.Net))
	}
}

func printLedger(ctx context.Context, logger *slog.Logger, book *ledger.Ledger) {
	balances, err := book.Balances(ctx)
	if err != nil {
		logger.Warn("balances unavailable", "err", err)
		return
	}
	for _, b := range balances {
		logger.Info("ledger balance", "account", b.AccountID, "balance", clearing.FormatAmount(b.Balance))
	}
}

func printReport(logger *slog.Logger, report *ledger.Report) {
	logger.Info("reconciliation", "cycle", report.CycleID,
		"statements", report.Statements, "ledger_postings", report.LedgerPostings,
		"matched", report.Matched, "statement_total", clearing.FormatAmount(report.StatementTotal),
		"ledger_total", clearing.FormatAmount(report.LedgerTotal))
	for _, m := range report.Mismatches {
		logger.Warn("mismatch", "member", m.Member, "kind", m.Kind,
			"expected", clearing.FormatAmount(m.Expected), "actual", clearing.FormatAmount(m.Actual), "detail", m.Detail)
	}
}
