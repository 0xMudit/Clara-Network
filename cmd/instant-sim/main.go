// Command instant-sim exercises the Clara Network instant-payment layer
// (docs/24): ISO 20022 pacs.008 customer credit transfers settled in real
// time, 24/7/365, against fully prefunded member positions (the RTP model)
// with a 20-second scheme SLA. It demonstrates the happy path, insufficient
// funds (AC04), unknown beneficiary (AC01), unknown sender / self transfer
// (AG01), format errors (FF01), and a beneficiary that misses the SLA (NOAS,
// with the reservation released), then verifies position conservation.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/0xMudit/Clara-Network/internal/env"
	"github.com/0xMudit/Clara-Network/internal/instant"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Beneficiary PSPs confirm instantly, except SLEEPY which never answers
	// within the SLA (used to demonstrate the NOAS reject and release).
	forward := func(ctx context.Context, p instant.Payment) error {
		if p.Beneficiary == "SLEEPY" {
			<-ctx.Done()
			return ctx.Err()
		}
		logger.Info("beneficiary PSP confirmed", "beneficiary", p.Beneficiary, "txid", p.TxID)
		return nil
	}

	engine := instant.NewEngine(instant.Config{Currency: "USD", Forward: forward})
	fund(logger, engine, "BANK-A", "1000.00")
	fund(logger, engine, "BANK-B", "500.00")

	initial := sumPositions(engine.Positions())
	logger.Info("instant-payment engine booted",
		"currency", "USD", "sla", "20s", "prefunded_total", format(initial))

	send(logger, engine, ctx, "PAY-001", "BANK-A", "BANK-B", "25.00",
		"DE89370400440532013000", "GB29NWBK60161331926819", "invoice 42") // ACSC
	send(logger, engine, ctx, "PAY-002", "BANK-A", "BANK-B", "5000.00",
		"DE89370400440532013000", "GB29NWBK60161331926819", "overdraft") // AC04
	send(logger, engine, ctx, "PAY-003", "BANK-A", "NO-SUCH", "10.00",
		"DE89370400440532013000", "NL91ABNA0417164300", "wrong bank") // AC01
	send(logger, engine, ctx, "PAY-004", "NO-SUCH", "BANK-B", "10.00",
		"DE89370400440532013000", "GB29NWBK60161331926819", "unknown") // AG01

	// A EUR payment into a USD-only engine is a format/rule error (FF01).
	p := instant.Payment{
		MsgID: "PAY-005", InstrID: "PAY-005", EndToEndID: "PAY-005-E2E", TxID: "TX-PAY-005",
		Sender: "BANK-A", Beneficiary: "BANK-B",
		SenderIBAN: "DE89370400440532013000", BeneficiaryIBAN: "GB29NWBK60161331926819",
		AmountMinor: 1000, Currency: "EUR", Remittance: "currency mismatch",
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
	res := engine.Transfer(ctx, p)
	logger.Info("instant payment",
		"msg_id", "PAY-005", "sender", "BANK-A", "beneficiary", "BANK-B",
		"amount", "10.00", "status", "RJCT ("+res.Reason+")", "final", res.Final)

	// A beneficiary that misses the SLA: the engine rejects with NOAS and
	// releases the reservation. A separate engine with a short SLA shows the
	// configurable timeout without waiting 20 seconds.
	sla := time.Duration(env.Int("CLARA_INSTANT_SLA", 3)) * time.Second
	sleepy := instant.NewEngine(instant.Config{Currency: "USD", Timeout: sla, Forward: forward})
	fund(logger, sleepy, "BANK-A", "1000.00")
	fund(logger, sleepy, "SLEEPY", "500.00")
	logger.Info("timeout drill", "engine_sla", sla)
	send(logger, sleepy, ctx, "PAY-006", "BANK-A", "SLEEPY", "10.00",
		"DE89370400440532013000", "FR1420041010050500013M02606", "slow bank") // NOAS

	printPositions(logger, engine.Positions())
	printHistory(logger, engine.History())

	final := sumPositions(engine.Positions())
	logger.Info("conservation check",
		"prefunded_total", format(initial), "final_total", format(final),
		"conserved", initial == final)
	if initial != final {
		logger.Error("position conservation violated", "initial", initial, "final", final)
		os.Exit(1)
	}
	logger.Info("instant payments demo complete")
}

func send(logger *slog.Logger, engine *instant.Engine, ctx context.Context,
	msgID, sender, beneficiary, amount, senderIBAN, beneficiaryIBAN, remittance string) {
	minor, err := parse(amount)
	if err != nil {
		logger.Error("bad amount", "amount", amount, "err", err)
		os.Exit(1)
	}
	p := instant.Payment{
		MsgID: msgID, InstrID: msgID, EndToEndID: msgID + "-E2E", TxID: "TX-" + msgID,
		Sender: sender, Beneficiary: beneficiary,
		SenderIBAN: senderIBAN, BeneficiaryIBAN: beneficiaryIBAN,
		AmountMinor: minor, Currency: "USD", Remittance: remittance,
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
	inbound, err := instant.BuildPacs008(p)
	if err != nil {
		logger.Error("build pacs.008 failed", "err", err)
		os.Exit(1)
	}
	res := engine.Transfer(ctx, p)
	outbound, err := instant.BuildPacs002(instant.StatusReport{
		MsgID: "STS-" + msgID, OriginalMsgID: p.MsgID, OriginalMsgName: "pacs.008.001.09",
		EndToEndID: p.EndToEndID, TxID: p.TxID, Status: res.Status, Reason: res.Reason,
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})
	if err != nil {
		logger.Error("build pacs.002 failed", "err", err)
		os.Exit(1)
	}
	status := res.Status
	if res.Reason != "" {
		status += " (" + res.Reason + ")"
	}
	logger.Info("instant payment",
		"msg_id", msgID, "sender", sender, "beneficiary", beneficiary,
		"amount", amount, "status", status, "final", res.Final)
	logger.Info("pacs.008 inbound", "xml", string(inbound))
	logger.Info("pacs.002 outbound", "xml", string(outbound))
}

func fund(logger *slog.Logger, engine *instant.Engine, psp, amount string) {
	minor, err := parse(amount)
	if err != nil {
		logger.Error("bad amount", "amount", amount, "err", err)
		os.Exit(1)
	}
	if _, err := engine.SetPosition(psp, minor); err != nil {
		logger.Error("prefund failed", "psp", psp, "err", err)
		os.Exit(1)
	}
	logger.Info("position funded", "psp", psp, "balance", amount)
}

func printPositions(logger *slog.Logger, positions map[string]int64) {
	logger.Info("final positions")
	for _, psp := range sortedKeys(positions) {
		logger.Info("position", "psp", psp, "balance", format(positions[psp]))
	}
}

func printHistory(logger *slog.Logger, history []instant.Settlement) {
	logger.Info("settlement history")
	for _, st := range history {
		status := st.Status
		if st.Reason != "" {
			status += " (" + st.Reason + ")"
		}
		logger.Info("settlement", "txid", st.TxID, "sender", st.Sender,
			"beneficiary", st.Beneficiary, "amount", format(st.AmountMinor),
			"status", status, "final", st.Final)
	}
}

func sortedKeys(m map[string]int64) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sumPositions(m map[string]int64) int64 {
	var total int64
	for _, v := range m {
		total += v
	}
	return total
}

// parse converts "major.minor" to minor units.
func parse(s string) (int64, error) {
	var major, frac int64
	dot := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dot = i
			break
		}
	}
	if dot < 0 {
		for _, c := range s {
			if c < '0' || c > '9' {
				return 0, fmt.Errorf("invalid amount %q", s)
			}
			major = major*10 + int64(c-'0')
		}
		return major * 100, nil
	}
	for _, c := range s[:dot] {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		major = major*10 + int64(c-'0')
	}
	fracStr := s[dot+1:]
	for len(fracStr) < 2 {
		fracStr += "0"
	}
	for _, c := range fracStr[:2] {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		frac = frac*10 + int64(c-'0')
	}
	return major*100 + frac, nil
}

func format(minor int64) string {
	return fmt.Sprintf("%d.%02d", minor/100, minor%100)
}
