package adminapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// ClearingRecord mirrors a row from clearing_records.
type ClearingRecord struct {
	CycleID     string `json:"cycleId"`
	STAN        string `json:"stan"`
	MTI         string `json:"mti"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	AmountMinor int64  `json:"amountMinor"`
	Interchange int64  `json:"interchange"`
	Currency    string `json:"currency"`
	RefID       string `json:"refId"`
}

// ClearingRecords returns all clearing records for a settlement cycle.
func (s *Store) ClearingRecords(ctx context.Context, cycleID string) ([]ClearingRecord, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT cycle_id, stan, mti, sender, receiver, amount_minor, interchange, currency, ref_id
		 FROM clearing_records WHERE cycle_id = $1 ORDER BY id`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ClearingRecord
	for rows.Next() {
		var r ClearingRecord
		if err := rows.Scan(&r.CycleID, &r.STAN, &r.MTI, &r.Sender, &r.Receiver,
			&r.AmountMinor, &r.Interchange, &r.Currency, &r.RefID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ClearingCycles returns the distinct cycle IDs in the clearing records.
func (s *Store) ClearingCycles(ctx context.Context) ([]string, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT DISTINCT cycle_id FROM clearing_records ORDER BY cycle_id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// NetPosition mirrors a row from net_positions.
type NetPosition struct {
	CycleID string `json:"cycleId"`
	Member  string `json:"member"`
	Net     int64  `json:"net"`
}

// NetPositions returns the net settlement positions for a cycle.
func (s *Store) NetPositions(ctx context.Context, cycleID string) ([]NetPosition, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT cycle_id, member, net_minor FROM net_positions WHERE cycle_id = $1 ORDER BY member`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []NetPosition
	for rows.Next() {
		var p NetPosition
		if err := rows.Scan(&p.CycleID, &p.Member, &p.Net); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// SettlementInstruction mirrors a row from settlement_instructions.
type SettlementInstruction struct {
	CycleID     string    `json:"cycleId"`
	MsgID       string    `json:"msgId"`
	Member      string    `json:"member"`
	Amount      int64     `json:"amount"`
	Direction   string    `json:"direction"`
	Currency    string    `json:"currency"`
	Instruction time.Time `json:"instruction"`
	Final       bool      `json:"final"`
}

// SettlementInstructions returns settlement instructions for a cycle.
func (s *Store) SettlementInstructions(ctx context.Context, cycleID string) ([]SettlementInstruction, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT cycle_id, msg_id, member, amount_minor, direction, currency, instruction_time, final
		 FROM settlement_instructions WHERE cycle_id = $1 ORDER BY id`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SettlementInstruction
	for rows.Next() {
		var i SettlementInstruction
		if err := rows.Scan(&i.CycleID, &i.MsgID, &i.Member, &i.Amount, &i.Direction,
			&i.Currency, &i.Instruction, &i.Final); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

// PrefundAccount mirrors a row from prefund_accounts.
type PrefundAccount struct {
	Member  string `json:"member"`
	Balance int64  `json:"balance"`
	Cap     int64  `json:"cap"`
}

// PrefundAccounts returns all member prefund accounts.
func (s *Store) PrefundAccounts(ctx context.Context) ([]PrefundAccount, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT member, balance, cap FROM prefund_accounts ORDER BY member`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PrefundAccount
	for rows.Next() {
		var a PrefundAccount
		if err := rows.Scan(&a.Member, &a.Balance, &a.Cap); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// DefaultFundBalance returns the default fund balance.
func (s *Store) DefaultFundBalance(ctx context.Context) (int64, error) {
	var bal int64
	err := s.Pool.QueryRow(ctx, `SELECT balance FROM default_fund WHERE id = 1`).Scan(&bal)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return bal, err
}
