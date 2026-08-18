package adminapi

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5"
)

// CardRecord is a card record as the issuer holds it (masked, no clear PAN).
type CardRecord struct {
	Ref     string `json:"ref"`
	PANHash string `json:"panHash"`
	PANMask string `json:"panMask"`
	BIN     string `json:"bin"`
	Expiry  string `json:"expiry"`
	Status  string `json:"status"`
	Product string `json:"product"`
	LastATC uint16 `json:"lastAtc"`
}

// Cards returns a page of issued cards (masked PANs, no clear data).
func (s *Store) Cards(ctx context.Context, limit, offset int) (Page[CardRecord], error) {
	total, err := s.count(ctx, "cards")
	if err != nil {
		return Page[CardRecord]{}, err
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT ref, pan_hash, pan_masked, bin, expiry, status, product, last_atc
		 FROM cards ORDER BY ref LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return Page[CardRecord]{}, err
	}
	defer rows.Close()
	var out []CardRecord
	for rows.Next() {
		var c CardRecord
		var hash []byte
		if err := rows.Scan(&c.Ref, &hash, &c.PANMask, &c.BIN, &c.Expiry,
			&c.Status, &c.Product, &c.LastATC); err != nil {
			return Page[CardRecord]{}, err
		}
		c.PANHash = hex.EncodeToString(hash)
		out = append(out, c)
	}
	return Page[CardRecord]{Items: out, Total: total}, rows.Err()
}

// Card returns a single card by reference.
func (s *Store) Card(ctx context.Context, ref string) (CardRecord, bool, error) {
	var c CardRecord
	var hash []byte
	err := s.Pool.QueryRow(ctx,
		`SELECT ref, pan_hash, pan_masked, bin, expiry, status, product, last_atc
		 FROM cards WHERE ref = $1`, ref).
		Scan(&c.Ref, &hash, &c.PANMask, &c.BIN, &c.Expiry, &c.Status, &c.Product, &c.LastATC)
	if err == pgx.ErrNoRows {
		return c, false, nil
	}
	if err != nil {
		return c, false, err
	}
	c.PANHash = hex.EncodeToString(hash)
	return c, true, nil
}

// BinRangeRecord is a BIN range the issuer owns.
type BinRangeRecord struct {
	BIN      string `json:"bin"`
	Low      int64  `json:"low"`
	High     int64  `json:"high"`
	Currency string `json:"currency"`
	Product  string `json:"product"`
}

// BinRanges returns all registered BIN ranges.
func (s *Store) BinRanges(ctx context.Context) ([]BinRangeRecord, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT bin, low, high, currency, product FROM bin_ranges ORDER BY bin`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BinRangeRecord
	for rows.Next() {
		var r BinRangeRecord
		if err := rows.Scan(&r.BIN, &r.Low, &r.High, &r.Currency, &r.Product); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// TokenRecord is a payment token (no PAN hash exposed).
type TokenRecord struct {
	Number    string    `json:"token"`
	PAR       string    `json:"par"`
	Status    string    `json:"status"`
	BIN       string    `json:"bin"`
	Requestor string    `json:"requestor"`
	DeviceID  string    `json:"deviceId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Tokens returns a page of payment tokens.
func (s *Store) Tokens(ctx context.Context, limit, offset int) (Page[TokenRecord], error) {
	total, err := s.count(ctx, "tokens")
	if err != nil {
		return Page[TokenRecord]{}, err
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT token, par, status, bin, COALESCE(trid,''), COALESCE(device_id,''), created_at
		 FROM tokens ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return Page[TokenRecord]{}, err
	}
	defer rows.Close()
	var out []TokenRecord
	for rows.Next() {
		var t TokenRecord
		if err := rows.Scan(&t.Number, &t.PAR, &t.Status, &t.BIN,
			&t.Requestor, &t.DeviceID, &t.CreatedAt); err != nil {
			return Page[TokenRecord]{}, err
		}
		out = append(out, t)
	}
	return Page[TokenRecord]{Items: out, Total: total}, rows.Err()
}
