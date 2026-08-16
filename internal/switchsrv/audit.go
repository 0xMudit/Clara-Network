package switchsrv

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditEvent is a record of a processed transaction.
type AuditEvent struct {
	MTI          string
	STAN         string
	PAN          string
	Amount       string
	ResponseCode string
	Destination  string
	CreatedAt    time.Time
}

// AuditSink records processed transactions for compliance and reconciliation.
type AuditSink interface {
	Record(ctx context.Context, ev AuditEvent) error
}

// NoopAudit discards events.
type NoopAudit struct{}

func (NoopAudit) Record(context.Context, AuditEvent) error { return nil }

// PostgresAudit persists events to the switch_transactions table.
type PostgresAudit struct {
	Pool *pgxpool.Pool
}

func (p *PostgresAudit) Record(ctx context.Context, ev AuditEvent) error {
	_, err := p.Pool.Exec(ctx,
		`INSERT INTO switch_transactions (stan, mti, pan_masked, amount, response_code, destination, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ev.STAN, ev.MTI, ev.PAN, ev.Amount, ev.ResponseCode, ev.Destination, ev.CreatedAt)
	return err
}
