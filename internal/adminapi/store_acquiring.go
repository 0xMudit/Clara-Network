package adminapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// MerchantRecord is a boarded merchant.
type MerchantRecord struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	DBA              string    `json:"dba"`
	TaxID            string    `json:"taxId"`
	Principals       []string  `json:"principals"`
	MCCs             []string  `json:"mccs"`
	Status           string    `json:"status"`
	RiskTier         string    `json:"riskTier"`
	ReserveRateBPS   int64     `json:"reserveRateBps"`
	FundingDelayDays int       `json:"fundingDelayDays"`
	TransactionLimit int64     `json:"transactionLimit"`
	ReserveBalance   int64     `json:"reserveBalance"`
	Volume           int64     `json:"volume"`
	DeclineReason    string    `json:"declineReason,omitempty"`
	ApprovedAt       time.Time `json:"approvedAt"`
}

// Merchants returns a page of boarded merchants.
func (s *Store) Merchants(ctx context.Context, limit, offset int) (Page[MerchantRecord], error) {
	total, err := s.count(ctx, "merchants")
	if err != nil {
		return Page[MerchantRecord]{}, err
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT id, name, dba, tax_id, principals, mccs, status, risk_tier,
		        reserve_rate_bps, funding_delay_days, transaction_limit,
		        reserve_balance, volume, COALESCE(decline_reason,''),
		        COALESCE(approved_at, now())
		 FROM merchants ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return Page[MerchantRecord]{}, err
	}
	defer rows.Close()
	var out []MerchantRecord
	for rows.Next() {
		var m MerchantRecord
		var principals, mccs string
		if err := rows.Scan(&m.ID, &m.Name, &m.DBA, &m.TaxID, &principals, &mccs,
			&m.Status, &m.RiskTier, &m.ReserveRateBPS, &m.FundingDelayDays,
			&m.TransactionLimit, &m.ReserveBalance, &m.Volume, &m.DeclineReason,
			&m.ApprovedAt); err != nil {
			return Page[MerchantRecord]{}, err
		}
		m.Principals = splitCSV(principals)
		m.MCCs = splitCSV(mccs)
		out = append(out, m)
	}
	return Page[MerchantRecord]{Items: out, Total: total}, rows.Err()
}

// Merchant returns a single merchant by ID.
func (s *Store) Merchant(ctx context.Context, id string) (MerchantRecord, bool, error) {
	var m MerchantRecord
	var principals, mccs string
	err := s.Pool.QueryRow(ctx,
		`SELECT id, name, dba, tax_id, principals, mccs, status, risk_tier,
		        reserve_rate_bps, funding_delay_days, transaction_limit,
		        reserve_balance, volume, COALESCE(decline_reason,''),
		        COALESCE(approved_at, now())
		 FROM merchants WHERE id = $1`, id).
		Scan(&m.ID, &m.Name, &m.DBA, &m.TaxID, &principals, &mccs,
			&m.Status, &m.RiskTier, &m.ReserveRateBPS, &m.FundingDelayDays,
			&m.TransactionLimit, &m.ReserveBalance, &m.Volume, &m.DeclineReason,
			&m.ApprovedAt)
	if err == pgx.ErrNoRows {
		return m, false, nil
	}
	if err != nil {
		return m, false, err
	}
	m.Principals = splitCSV(principals)
	m.MCCs = splitCSV(mccs)
	return m, true, nil
}

// FundingLineRecord is one merchant settlement funding batch.
type FundingLineRecord struct {
	BatchID     string    `json:"batchId"`
	MerchantID  string    `json:"merchantId"`
	Gross       int64     `json:"gross"`
	Fees        int64     `json:"fees"`
	ReserveHold int64     `json:"reserveHold"`
	Net         int64     `json:"net"`
	Date        time.Time `json:"date"`
}

// MerchantFunding returns funding lines for a merchant.
func (s *Store) MerchantFunding(ctx context.Context, merchantID string) ([]FundingLineRecord, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT batch_id, merchant_id, gross, fees, reserve_hold, net, date
		 FROM funding_lines WHERE merchant_id = $1 ORDER BY date`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FundingLineRecord
	for rows.Next() {
		var l FundingLineRecord
		if err := rows.Scan(&l.BatchID, &l.MerchantID, &l.Gross, &l.Fees,
			&l.ReserveHold, &l.Net, &l.Date); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
