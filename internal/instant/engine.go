package instant

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Forwarder delivers a settled payment to the beneficiary PSP and must
// confirm within the scheme SLA; returning an error (or taking longer than
// the SLA) causes the reservation to be released and the payment rejected.
type Forwarder func(ctx context.Context, p Payment) error

// Config configures the instant settlement engine.
type Config struct {
	Currency string        // settlement currency (ISO 4217 alpha), e.g. "USD"
	Timeout  time.Duration // scheme SLA for the beneficiary to confirm; default 20s
	Forward  Forwarder     // delivers the payment to the beneficiary PSP; nil settles locally
	Log      *slog.Logger
}

// Settlement is one recorded outcome in the engine's history.
type Settlement struct {
	MsgID       string
	TxID        string
	Sender      string
	Beneficiary string
	AmountMinor int64
	Currency    string
	Status      string // ACSC / RJCT
	Reason      string // reason code when RJCT
	Final       bool
	SettledAt   time.Time
}

// Result is the outcome of a Transfer and the positions after it.
type Result struct {
	MsgID     string
	TxID      string
	Status    string
	Reason    string
	Final     bool
	Positions map[string]int64
	SettledAt time.Time
}

// Engine settles instant payments in real time against fully prefunded
// member positions (docs/24 §24.2 RTP model): it verifies and reserves the
// sender's settlement capacity before forwarding, and settlement is complete
// when the sender's position is debited and the beneficiary's is credited.
type Engine struct {
	cfg      Config
	mu       sync.Mutex
	balances map[string]int64 // available prefunded position per PSP
	history  []Settlement
	log      *slog.Logger
}

// NewEngine builds an instant settlement engine.
func NewEngine(cfg Config) *Engine {
	if cfg.Currency == "" {
		cfg.Currency = "USD"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 20 * time.Second
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Engine{
		cfg:      cfg,
		balances: map[string]int64{},
		log:      cfg.Log,
	}
}

// SetPosition pre-funds a PSP's settlement account (a dedicated cash account
// topped up from its RTGS account) and returns the new balance.
func (e *Engine) SetPosition(psp string, balance int64) (int64, error) {
	if psp == "" {
		return 0, fmt.Errorf("instant: empty PSP")
	}
	if balance < 0 {
		return 0, fmt.Errorf("instant: negative position for %q", psp)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.balances[psp] = balance
	return balance, nil
}

// Position returns the available prefunded position of a PSP.
func (e *Engine) Position(psp string) (int64, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	b, ok := e.balances[psp]
	return b, ok
}

// Positions returns a snapshot of every prefunded position.
func (e *Engine) Positions() map[string]int64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make(map[string]int64, len(e.balances))
	for k, v := range e.balances {
		out[k] = v
	}
	return out
}

// History returns a copy of the settlement history in arrival order.
func (e *Engine) History() []Settlement {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]Settlement(nil), e.history...)
}

// Transfer settles an instant payment. It validates the message, verifies and
// reserves the sender's settlement capacity, forwards to the beneficiary PSP
// within the SLA, then settles (debit sender, credit beneficiary) with
// finality on confirmation. Rejections never move funds.
func (e *Engine) Transfer(ctx context.Context, p Payment) Result {
	if err := p.Validate(); err != nil {
		return e.reject(p, ReasonFormat, err.Error())
	}
	if p.Currency != e.cfg.Currency {
		return e.reject(p, ReasonFormat, fmt.Sprintf("currency %s not settled in %s", p.Currency, e.cfg.Currency))
	}

	if ok, reason := e.reserve(p); !ok {
		return e.reject(p, reason, reasonLabel(reason))
	}

	fwdCtx, cancel := context.WithTimeout(ctx, e.cfg.Timeout)
	defer cancel()
	if e.cfg.Forward != nil {
		if err := e.cfg.Forward(fwdCtx, p); err != nil {
			reason := ReasonForbidden
			if fwdCtx.Err() == context.DeadlineExceeded {
				reason = ReasonNoAnswer
			}
			e.release(p)
			e.log.Warn("instant payment rejected",
				"txid", p.TxID, "sender", p.Sender, "beneficiary", p.Beneficiary,
				"amount", p.AmountMinor, "reason", reason, "err", err)
			return e.reject(p, reason, reasonLabel(reason))
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.balances[p.Beneficiary] += p.AmountMinor
	st := Settlement{
		MsgID: p.MsgID, TxID: p.TxID, Sender: p.Sender, Beneficiary: p.Beneficiary,
		AmountMinor: p.AmountMinor, Currency: p.Currency, Status: StatusACSC, Final: true,
		SettledAt: time.Now().UTC(),
	}
	e.history = append(e.history, st)
	return e.resultLocked(p, st)
}

// reserve verifies the sender and beneficiary exist and that the sender's
// prefunded position covers the amount, then debits it as a reservation.
func (e *Engine) reserve(p Payment) (bool, string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.balances[p.Sender]; !ok {
		return false, ReasonForbidden // unknown sender PSP
	}
	if _, ok := e.balances[p.Beneficiary]; !ok {
		return false, ReasonAccount // unknown beneficiary PSP
	}
	if e.balances[p.Sender] < p.AmountMinor {
		return false, ReasonInsufficientFunds
	}
	e.balances[p.Sender] -= p.AmountMinor
	return true, ""
}

// release returns an unused reservation to the sender's position.
func (e *Engine) release(p Payment) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.balances[p.Sender] += p.AmountMinor
}

func (e *Engine) reject(p Payment, reason, detail string) Result {
	e.mu.Lock()
	defer e.mu.Unlock()
	st := Settlement{
		MsgID: p.MsgID, TxID: p.TxID, Sender: p.Sender, Beneficiary: p.Beneficiary,
		AmountMinor: p.AmountMinor, Currency: p.Currency, Status: StatusRJCT, Reason: reason,
		SettledAt: time.Now().UTC(),
	}
	e.history = append(e.history, st)
	e.log.Warn("instant payment rejected",
		"txid", p.TxID, "sender", p.Sender, "beneficiary", p.Beneficiary,
		"amount", p.AmountMinor, "reason", reason, "detail", detail)
	return e.resultLocked(p, st)
}

func (e *Engine) resultLocked(p Payment, st Settlement) Result {
	positions := make(map[string]int64, len(e.balances))
	for k, v := range e.balances {
		positions[k] = v
	}
	return Result{
		MsgID:     p.MsgID,
		TxID:      p.TxID,
		Status:    st.Status,
		Reason:    st.Reason,
		Final:     st.Final,
		Positions: positions,
		SettledAt: st.SettledAt,
	}
}

func reasonLabel(reason string) string {
	switch reason {
	case ReasonInsufficientFunds:
		return "insufficient funds in sender position"
	case ReasonAccount:
		return "unknown or incorrect beneficiary"
	case ReasonForbidden:
		return "transaction forbidden"
	case ReasonFormat:
		return "invalid message format"
	case ReasonNoAnswer:
		return "no answer from beneficiary within SLA"
	default:
		return reason
	}
}
