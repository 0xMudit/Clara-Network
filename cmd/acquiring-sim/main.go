// Command acquiring-sim exercises the Clara Network acquiring stack: it seeds
// the MATCH/OFAC negative lists, underwrites a set of merchant boarding
// applications (approve/decline/conditional), runs merchant settlement batches
// through the funding engine with fee and reserve withholding, and releases
// reserve balances back to merchants.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xMudit/Clara-Network/internal/acquiring"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var store acquiring.Store = acquiring.NewMemoryStore()
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logger.Error("postgres unavailable, using in-memory store", "err", err)
		} else {
			store = &acquiring.PostgresStore{Pool: pool}
		}
	}

	svc := acquiring.NewService(store)
	eng := acquiring.NewFundingEngine(store, svc.Policy())

	if err := seedLists(ctx, store); err != nil {
		logger.Error("seeding negative lists failed", "err", err)
		os.Exit(1)
	}

	applications := []acquiring.Application{
		{
			MerchantName: "Bob's Grocery", DBA: "Bob's Grocery", TaxID: "98-7654321",
			Principals: []string{"Robert Smith"}, MCCs: []string{"5411"},
			CreditScore: 720, Volume: 10_000_000,
		},
		{
			MerchantName: "Quantum Trading LLC", DBA: "Quantum", TaxID: "12-3456789",
			Principals: []string{"Jane Doe"}, MCCs: []string{"5999"},
			CreditScore: 780, Volume: 5_000_000,
		},
		{
			MerchantName: "Petrov Trading", DBA: "Petrov", TaxID: "11-2233445",
			Principals: []string{"Vladimir Petrov"}, MCCs: []string{"5812"},
			CreditScore: 710, Volume: 2_000_000,
		},
		{
			MerchantName: "Golden Spins Casino", DBA: "Golden Spins", TaxID: "55-1122334",
			Principals: []string{"Mona Reyes"}, MCCs: []string{"7995"},
			CreditScore: 700, Volume: 1_000_000, EnhancedDD: true,
		},
		{
			MerchantName: "Sketchy Casino", DBA: "Sketchy", TaxID: "55-9988776",
			Principals: []string{"Alex Kim"}, MCCs: []string{"7995"},
			CreditScore: 690, Volume: 8_000_000,
		},
	}

	boarded := map[string]*acquiring.Merchant{}
	for _, app := range applications {
		m, d, err := svc.Board(ctx, app)
		if err != nil {
			logger.Error("boarding failed", "merchant", app.MerchantName, "err", err)
			continue
		}
		logger.Info("decision", "merchant", app.MerchantName, "status", d.Status,
			"tier", d.RiskTier, "reserve_bps", d.ReserveRateBPS,
			"delay_days", d.FundingDelayDays, "limit", acquiring.Amount(d.TransactionLimit),
			"reasons", d.Reasons)
		if m != nil && m.Status == acquiring.StatusActive {
			boarded[m.ID] = m
		}
	}

	// Settlement: low-risk grocery settles next day with no reserve; the
	// high-risk casino has 10% withheld into a rolling reserve.
	batch := []struct {
		id    string
		gross int64
	}{
		{"M-bob's-grocery", 1_000_000},
		{"M-bob's-grocery", 750_000},
		{"M-golden-spins-casino", 500_000},
		{"M-golden-spins-casino", 500_000},
		{"M-golden-spins-casino", 500_000},
	}
	for _, b := range batch {
		line, err := eng.SettleBatch(ctx, b.id, b.gross)
		if err != nil {
			logger.Error("funding failed", "merchant", b.id, "err", err)
			continue
		}
		logger.Info("funding", "merchant", b.id, "batch", line.BatchID,
			"gross", acquiring.Amount(line.Gross), "fees", acquiring.Amount(line.Fees),
			"reserve_hold", acquiring.Amount(line.ReserveHold), "net", acquiring.Amount(line.Net),
			"funding_date", line.Date.Format("2006-01-02"))
	}

	for id := range boarded {
		m, _ := svc.GetMerchant(ctx, id)
		logger.Info("merchant state", "id", id, "status", m.Status, "tier", m.RiskTier,
			"reserve_balance", acquiring.Amount(m.ReserveBalance),
			"reserve_cap", acquiring.Amount(m.ReserveCap()))
	}

	// The casino reserve exceeds its cap -> release the excess to funding.
	line, err := eng.ReleaseReserve(ctx, "M-golden-spins-casino")
	if err != nil {
		logger.Warn("no reserve release", "err", err)
	} else {
		logger.Info("reserve released", "merchant", line.MerchantID,
			"amount", acquiring.Amount(line.Net), "batch", line.BatchID)
	}

	m, err := svc.GetMerchant(ctx, "M-bob's-grocery")
	if err != nil {
		logger.Error("load merchant failed", "err", err)
		os.Exit(1)
	}
	fmt.Printf("\nFunding lines for %s:\n", m.Name)
	lines, err := store.Funding(ctx, m.ID)
	if err != nil {
		logger.Error("load funding failed", "err", err)
		os.Exit(1)
	}
	for _, l := range lines {
		fmt.Printf("  %-38s gross=%s fees=%s reserve=%s net=%s date=%s\n",
			l.BatchID, acquiring.Amount(l.Gross), acquiring.Amount(l.Fees),
			acquiring.Amount(l.ReserveHold), acquiring.Amount(l.Net), l.Date.Format("2006-01-02"))
	}
}

func seedLists(ctx context.Context, store acquiring.Store) error {
	if err := store.SaveMatchEntries(ctx, []acquiring.MatchEntry{
		{MerchantName: "Quantum Trading LLC", TaxID: "12-3456789", Reason: "excessive chargebacks"},
	}); err != nil {
		return err
	}
	return store.SaveOfacEntries(ctx, []acquiring.OfacEntry{
		{Name: "Vladimir Petrov", Program: "OFAC SDN"},
	})
}
