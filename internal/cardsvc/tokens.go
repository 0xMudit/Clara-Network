package cardsvc

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

// TokenStatus values.
const (
	TokenActive   = "active"
	TokenSuspended = "suspended"
)

// Token is an EMV payment token (docs/07): a PAN-like surrogate bound to a
// Payment Account Reference (PAR) that links the token to the real PAN. The
// vault never stores the PAN alongside the token; only its hash.
type Token struct {
	Number    string // PAN-like token number
	PANHash   []byte
	PAR       string
	Status    string
	BIN       string
	Requestor string // token requestor ID (TRID)
	DeviceID  string // provisioned wallet/device, if any
	CreatedAt time.Time
}

// Provisioned returns the token's wallet provisioning payload (docs/07 §7.5,
// docs/25 §25.4 phase 5): the wallet receives the token plus the PAR.
type Provisioned struct {
	Token     string
	PAR       string
	BIN       string
	Requestor string
	DeviceID  string
}

const (
	parLength    = 29
	parAlphabet  = "ABCDEFGHJKMNPQRSTUVWXYZ23456789" // excludes I, O, L, 0, 1
	tokenLength  = 16
)

// TokenVault maps PANs to payment tokens.
type TokenVault struct {
	store Store
}

// NewTokenVault builds a token vault over the store.
func NewTokenVault(store Store) *TokenVault {
	return &TokenVault{store: store}
}

// Tokenize creates (or returns the existing) token and PAR for a PAN. The
// PAN must be on an issued card.
func (v *TokenVault) Tokenize(ctx context.Context, svc *Service, pan string) (*Token, error) {
	hash := panHash(pan)
	if tok, ok, err := v.store.TokenByPANHash(ctx, hash); err != nil {
		return nil, err
	} else if ok && tok.Status == TokenActive {
		return &tok, nil
	}
	if _, err := svc.GetCard(ctx, cardRef(pan)); err != nil {
		return nil, fmt.Errorf("cardsvc: cannot tokenize unknown card: %w", err)
	}
	r, ok, err := svc.FindBinRange(ctx, pan)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cardsvc: pan is not in an issued BIN range")
	}
	token, err := v.generateToken(ctx, r.BIN)
	if err != nil {
		return nil, err
	}
	t := Token{
		Number:    token,
		PANHash:   hash,
		PAR:       v.generatePAR(),
		Status:    TokenActive,
		BIN:       r.BIN,
		CreatedAt: time.Now().UTC(),
	}
	if err := v.store.SaveToken(ctx, t); err != nil {
		return nil, fmt.Errorf("cardsvc: persist token: %w", err)
	}
	return &t, nil
}

// Detokenize returns the PAN hash for a token if the token is active.
func (v *TokenVault) Detokenize(ctx context.Context, token string) ([]byte, error) {
	t, ok, err := v.store.Token(ctx, token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cardsvc: unknown token %s", token)
	}
	if t.Status != TokenActive {
		return nil, fmt.Errorf("cardsvc: token %s is %s", token, t.Status)
	}
	return t.PANHash, nil
}

// Provision simulates mobile-wallet provisioning: it binds the token to a
// device and returns the payload delivered to the wallet.
func (v *TokenVault) Provision(ctx context.Context, token, deviceID, requestor string) (*Provisioned, error) {
	t, ok, err := v.store.Token(ctx, token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cardsvc: unknown token %s", token)
	}
	if t.Status != TokenActive {
		return nil, fmt.Errorf("cardsvc: token %s is %s", token, t.Status)
	}
	t.DeviceID = deviceID
	t.Requestor = requestor
	if err := v.store.SaveToken(ctx, t); err != nil {
		return nil, fmt.Errorf("cardsvc: persist provision: %w", err)
	}
	return &Provisioned{
		Token:     t.Number,
		PAR:       t.PAR,
		BIN:       t.BIN,
		Requestor: requestor,
		DeviceID:  deviceID,
	}, nil
}

func (v *TokenVault) generateToken(ctx context.Context, bin string) (string, error) {
	for i := 0; i < 5; i++ {
		digits := make([]byte, tokenLength-1)
		copy(digits, []byte(bin))
		if _, err := rand.Read(digits[len(bin):]); err != nil {
			return "", fmt.Errorf("cardsvc: random token: %w", err)
		}
		for j := len(bin); j < len(digits); j++ {
			digits[j] = '0' + digits[j]%10
		}
		token := string(digits) + string(luhnCheckDigit(digits))
		if _, ok, err := v.store.Token(ctx, token); err != nil {
			return "", err
		} else if !ok {
			return token, nil
		}
	}
	return "", fmt.Errorf("cardsvc: could not allocate a unique token")
}

func (v *TokenVault) generatePAR() string {
	buf := make([]byte, parLength)
	rand.Read(buf)
	out := make([]byte, parLength)
	for i := range buf {
		out[i] = parAlphabet[int(buf[i])%len(parAlphabet)]
	}
	return string(out)
}

func luhnCheckDigit(digits []byte) byte {
	sum := 0
	double := true
	for i := len(digits) - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return byte('0' + (10 - sum%10) % 10)
}
