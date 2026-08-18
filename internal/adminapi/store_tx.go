package adminapi

import (
	"context"
	"time"
)

// AuditEvent is a transaction audit record (mirrors switchsrv.AuditEvent without
// importing the switch package to avoid pulling in TCP/framing dependencies).
type AuditEvent struct {
	STAN         string    `json:"stan"`
	MTI          string    `json:"mti"`
	PAN          string    `json:"pan"`
	Amount       string    `json:"amount"`
	ResponseCode string    `json:"responseCode"`
	Destination  string    `json:"destination"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Transactions returns a page of authorization audit events, newest first.
func (s *Store) Transactions(ctx context.Context, limit, offset int) (Page[AuditEvent], error) {
	total, err := s.count(ctx, "switch_transactions")
	if err != nil {
		return Page[AuditEvent]{}, err
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT stan, mti, pan_masked, amount, response_code, destination, created_at
		 FROM switch_transactions ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return Page[AuditEvent]{}, err
	}
	defer rows.Close()
	var out []AuditEvent
	for rows.Next() {
		var ev AuditEvent
		if err := rows.Scan(&ev.STAN, &ev.MTI, &ev.PAN, &ev.Amount,
			&ev.ResponseCode, &ev.Destination, &ev.CreatedAt); err != nil {
			return Page[AuditEvent]{}, err
		}
		out = append(out, ev)
	}
	return Page[AuditEvent]{Items: out, Total: total}, rows.Err()
}
