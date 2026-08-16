// Package acquirersim simulates an acquirer host that sends batches of
// ISO 8583 authorization requests to the Clara Network switch.
package acquirersim

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/0xMudit/Clara-Network/internal/framing"
	"github.com/0xMudit/Clara-Network/internal/iso8583"
)

// Config describes a batch of simulated authorization requests.
type Config struct {
	SwitchAddr string
	AcquirerID string
	IssuerID   string
	PAN        string
	Count      int
	Amount     int // base amount in minor units
	Step       int // amount increment per request
	Log        *slog.Logger
}

// Run sends Count authorization requests to the switch and logs each response.
func Run(ctx context.Context, cfg Config) error {
	if cfg.Count <= 0 {
		cfg.Count = 1
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", cfg.SwitchAddr)
	if err != nil {
		return fmt.Errorf("acquirersim: connect to switch: %w", err)
	}
	defer conn.Close()

	now := time.Now()
	approved, declined := 0, 0
	for i := 0; i < cfg.Count; i++ {
		stan := fmt.Sprintf("%06d", 100000+i)
		amount := cfg.Amount + i*cfg.Step
		req := iso8583.New("0100").
			Set(2, cfg.PAN).
			Set(3, "000000").
			Set(4, fmt.Sprintf("%012d", amount)).
			Set(7, now.Format("0102150405")).
			Set(11, stan).
			Set(12, now.Format("150405")).
			Set(13, now.Format("0102")).
			Set(22, "022").
			Set(25, "00").
			Set(32, cfg.AcquirerID).
			Set(49, "840").
			Set(100, cfg.IssuerID)

		raw, err := req.Marshal()
		if err != nil {
			return fmt.Errorf("acquirersim: marshal %s: %w", stan, err)
		}
		if err := framing.WriteFrame(conn, raw); err != nil {
			return fmt.Errorf("acquirersim: write %s: %w", stan, err)
		}
		respFrame, err := framing.ReadFrame(conn)
		if err != nil {
			return fmt.Errorf("acquirersim: read %s: %w", stan, err)
		}
		resp, err := iso8583.Parse(respFrame)
		if err != nil {
			return fmt.Errorf("acquirersim: parse response %s: %w", stan, err)
		}
		code := resp.Get(39)
		if code == "00" {
			approved++
		} else {
			declined++
		}
		cfg.Log.Info("auth response",
			"stan", stan, "mti", resp.MTI, "pan", cfg.PAN,
			"amount", amount, "code", code)
	}
	cfg.Log.Info("summary", "sent", cfg.Count, "approved", approved, "declined", declined)
	return nil
}
