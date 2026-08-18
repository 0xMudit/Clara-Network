package adminapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// DisputeRecord is a dispute case through its lifecycle.
type DisputeRecord struct {
	ID             string    `json:"id"`
	RefID          string    `json:"refId"`
	MerchantID     string    `json:"merchantId"`
	Cardholder     string    `json:"cardholder"`
	AmountMinor    int64     `json:"amountMinor"`
	Currency       string    `json:"currency"`
	ReasonCode     string    `json:"reasonCode"`
	Category       string    `json:"category"`
	Stage          string    `json:"stage"`
	Status         string    `json:"status"`
	FiledAt        time.Time `json:"filedAt"`
	ResponseDue    time.Time `json:"responseDue"`
	RespondedAt    time.Time `json:"respondedAt,omitempty"`
	EscalatedAt    time.Time `json:"escalatedAt,omitempty"`
	Evidence       []string  `json:"evidence,omitempty"`
	Decision       string    `json:"decision,omitempty"`
	Winner         string    `json:"winner,omitempty"`
	DecisionAt     time.Time `json:"decisionAt,omitempty"`
	DisputeFee     int64     `json:"disputeFee"`
	ArbitrationFee int64     `json:"arbitrationFee"`
	Note           string    `json:"note,omitempty"`
}

// Disputes returns a page of disputes, optionally filtered by stage.
func (s *Store) Disputes(ctx context.Context, stage string, limit, offset int) (Page[DisputeRecord], error) {
	var total int64
	var err error
	var rows pgx.Rows

	if stage != "" {
		err = s.Pool.QueryRow(ctx,
			`SELECT count(*) FROM disputes WHERE stage = $1`, stage).Scan(&total)
		if err != nil {
			return Page[DisputeRecord]{}, err
		}
		rows, err = s.Pool.Query(ctx,
			`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency,
			        reason_code, category, stage, status, filed_at, response_due,
			        responded_at, escalated_at, evidence, decision, winner,
			        decision_at, dispute_fee, arbitration_fee, note
			 FROM disputes WHERE stage = $1 ORDER BY filed_at DESC LIMIT $2 OFFSET $3`,
			stage, limit, offset)
	} else {
		err = s.Pool.QueryRow(ctx,
			`SELECT count(*) FROM disputes`).Scan(&total)
		if err != nil {
			return Page[DisputeRecord]{}, err
		}
		rows, err = s.Pool.Query(ctx,
			`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency,
			        reason_code, category, stage, status, filed_at, response_due,
			        responded_at, escalated_at, evidence, decision, winner,
			        decision_at, dispute_fee, arbitration_fee, note
			 FROM disputes ORDER BY filed_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return Page[DisputeRecord]{}, err
	}
	defer rows.Close()
	out, err := scanDisputes(rows)
	return Page[DisputeRecord]{Items: out, Total: total}, err
}

// OverdueDisputes returns open disputes past their response deadline.
func (s *Store) OverdueDisputes(ctx context.Context) ([]DisputeRecord, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency,
		        reason_code, category, stage, status, filed_at, response_due,
		        responded_at, escalated_at, evidence, decision, winner,
		        decision_at, dispute_fee, arbitration_fee, note
		 FROM disputes
		 WHERE stage <> 'resolved' AND stage <> 'invalid'
		   AND response_due < now()
		 ORDER BY response_due`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDisputes(rows)
}

// Dispute returns a single dispute by ID.
func (s *Store) Dispute(ctx context.Context, id string) (DisputeRecord, bool, error) {
	var d DisputeRecord
	var evidence string
	err := s.Pool.QueryRow(ctx,
		`SELECT id, ref_id, merchant_id, cardholder, amount_minor, currency,
		        reason_code, category, stage, status, filed_at, response_due,
		        responded_at, escalated_at, evidence, decision, winner,
		        decision_at, dispute_fee, arbitration_fee, note
		 FROM disputes WHERE id = $1`, id).
		Scan(&d.ID, &d.RefID, &d.MerchantID, &d.Cardholder, &d.AmountMinor,
			&d.Currency, &d.ReasonCode, &d.Category, &d.Stage, &d.Status,
			&d.FiledAt, &d.ResponseDue, &d.RespondedAt, &d.EscalatedAt,
			&evidence, &d.Decision, &d.Winner, &d.DecisionAt,
			&d.DisputeFee, &d.ArbitrationFee, &d.Note)
	if err == pgx.ErrNoRows {
		return d, false, nil
	}
	if err != nil {
		return d, false, err
	}
	d.Evidence = splitCSV(evidence)
	return d, true, nil
}

// DisputeRatio computes a merchant's chargeback ratio from stored data.
func (s *Store) DisputeRatio(ctx context.Context, merchantID string) (ratio float64, status string, err error) {
	var txCount, cbCount int64
	err = s.Pool.QueryRow(ctx,
		`SELECT count(*) FROM dispute_transactions WHERE merchant_id = $1 AND is_credit = false`, merchantID).Scan(&txCount)
	if err != nil {
		return 0, "", err
	}
	err = s.Pool.QueryRow(ctx,
		`SELECT count(*) FROM disputes WHERE merchant_id = $1 AND decision = 'accepted'`, merchantID).Scan(&cbCount)
	if err != nil {
		return 0, "", err
	}
	if txCount > 0 {
		ratio = float64(cbCount) * 100 / float64(txCount)
	}
	status = "normal"
	if ratio >= 1.0 {
		status = "excessive"
	} else if ratio >= 0.5 {
		status = "watched"
	}
	return ratio, status, nil
}

func scanDisputes(rows pgx.Rows) ([]DisputeRecord, error) {
	var out []DisputeRecord
	for rows.Next() {
		var d DisputeRecord
		var evidence string
		if err := rows.Scan(&d.ID, &d.RefID, &d.MerchantID, &d.Cardholder,
			&d.AmountMinor, &d.Currency, &d.ReasonCode, &d.Category,
			&d.Stage, &d.Status, &d.FiledAt, &d.ResponseDue,
			&d.RespondedAt, &d.EscalatedAt, &evidence,
			&d.Decision, &d.Winner, &d.DecisionAt,
			&d.DisputeFee, &d.ArbitrationFee, &d.Note); err != nil {
			return nil, err
		}
		d.Evidence = splitCSV(evidence)
		out = append(out, d)
	}
	return out, rows.Err()
}
