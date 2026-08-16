// Command disputes-sim exercises the Clara Network disputes engine: it seeds
// monitored transactions, files disputes across reason-code categories, runs
// the representment -> scheme rule -> arbitration lifecycle with fees charged
// to the losing party, checks the associated-transaction rule and SLA
// overdues, and reports merchant chargeback ratios.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/disputes"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var store disputes.Store = disputes.NewMemoryStore()
	if dsn := os.Getenv("CLARA_PG_DSN"); dsn != "" {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logger.Error("postgres unavailable, using in-memory store", "err", err)
		} else {
			store = &disputes.PostgresStore{Pool: pool}
		}
	}

	svc := disputes.NewService(store)

	seed(ctx, logger, svc)

	// Case 1: fraud dispute, acquirer represents with full chip/3DS evidence -> acquirer wins.
	runCase(ctx, logger, svc, "T1", "M-online", "Mona Reyes", 8000, "4837",
		[]string{"chip", "3ds", "receipt"}, false)

	// Case 2: non-receipt dispute, missing delivery proof -> chargeback stands.
	runCase(ctx, logger, svc, "T2", "M-online", "John Kim", 12000, "4870",
		[]string{"receipt"}, false)

	// Case 3: cancelled-recurring dispute, acquirer wins, then issuer escalates to arbitration.
	runCase(ctx, logger, svc, "T3", "M-online", "Aria Shah", 20000, "4841",
		[]string{"terms", "receipt"}, true)

	// Case 4: the cardholder was already refunded -> associated-transaction rule rejects the filing.
	d, err := svc.File(ctx, disputes.FileRequest{RefID: "T4", Cardholder: "Kai Mehta", Amount: 3000, Currency: "840", ReasonCode: "13"})
	if err != nil {
		logger.Warn("filing rejected by associated-transaction rule", "err", err)
	} else {
		logger.Info("dispute filed", "id", d.ID)
	}

	// SLA: show overdue open disputes.
	overdue, err := svc.Overdue(ctx)
	if err != nil {
		logger.Error("overdue check failed", "err", err)
	}
	logger.Info("SLA check", "overdue_count", len(overdue))

	// Monitoring: chargeback ratios per merchant.
	for _, m := range []string{"M-online", "M-grocery"} {
		ratio, status, err := svc.MonitorRatio(ctx, m)
		if err != nil {
			logger.Error("monitoring failed", "merchant", m, "err", err)
			continue
		}
		logger.Info("chargeback monitoring", "merchant", m,
			"ratio_pct", fmt.Sprintf("%.2f%%", ratio), "status", status)
	}
}

func seed(ctx context.Context, logger *slog.Logger, svc *disputes.Service) {
	txs := []disputes.MonitoredTransaction{
		{RefID: "T1", MerchantID: "M-online", AmountMinor: 8000, Currency: "840"},
		{RefID: "T2", MerchantID: "M-online", AmountMinor: 12000, Currency: "840"},
		{RefID: "T3", MerchantID: "M-online", AmountMinor: 20000, Currency: "840"},
		{RefID: "T4", MerchantID: "M-online", AmountMinor: 3000, Currency: "840", IsCredit: true},
		{RefID: "G1", MerchantID: "M-grocery", AmountMinor: 1000, Currency: "840"},
	}
	for i := 0; i < 200; i++ {
		txs = append(txs, disputes.MonitoredTransaction{
			RefID: "G" + fmt.Sprint(1000+i), MerchantID: "M-grocery", AmountMinor: 1000, Currency: "840",
		})
	}
	for _, tx := range txs {
		if err := svc.RecordTransaction(ctx, tx); err != nil {
			logger.Error("seed transaction failed", "ref", tx.RefID, "err", err)
			os.Exit(1)
		}
	}
}

func runCase(ctx context.Context, logger *slog.Logger, svc *disputes.Service,
	ref, merchant, holder string, amount int64, code string, evidence []string, arbitrate bool) {
	d, err := svc.File(ctx, disputes.FileRequest{
		RefID: ref, MerchantID: merchant, Cardholder: holder, Amount: amount, Currency: "840", ReasonCode: code,
	})
	if err != nil {
		logger.Error("filing failed", "ref", ref, "err", err)
		return
	}
	if _, err := svc.Represent(ctx, d.ID, evidence); err != nil {
		logger.Error("representment failed", "id", d.ID, "err", err)
		return
	}
	ruled, err := svc.Rule(ctx, d.ID)
	if err != nil {
		logger.Error("ruling failed", "id", d.ID, "err", err)
		return
	}
	logOutcome(logger, ruled)
	if arbitrate && ruled.Decision == disputes.DecisionRejected {
		if _, err := svc.Escalate(ctx, d.ID); err != nil {
			logger.Error("escalation failed", "id", d.ID, "err", err)
			return
		}
		final, err := svc.Arbitrate(ctx, d.ID, true)
		if err != nil {
			logger.Error("arbitration failed", "id", d.ID, "err", err)
			return
		}
		logger.Warn("arbitration", "id", final.ID, "winner", final.Winner,
			"arbitration_fee", disputes.Amount(final.ArbitrationFee))
	}
}

func logOutcome(logger *slog.Logger, d *disputes.Dispute) {
	logger.Info("dispute ruled", "id", d.ID, "reason", d.ReasonCode+" "+d.Category,
		"decision", d.Decision, "winner", d.Winner,
		"dispute_fee", disputes.Amount(d.DisputeFee), "note", d.Note,
		"response_due", d.ResponseDue.Format("2006-01-02"), "days", int(d.ResponseDue.Sub(d.FiledAt)/(24*time.Hour)))
}
