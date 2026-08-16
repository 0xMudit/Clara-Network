// Package cardsvc implements the issuing stack (docs/25 §25.4 phase 5): BIN
// ranges, card data, EMV-style cryptogram verification, and the token vault.
package cardsvc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
)

// CardStatus values.
const (
	StatusActive   = "active"
	StatusBlocked  = "blocked"
	StatusExpired  = "expired"
)

// BinRange is a range of PANs the issuer owns under one BIN (docs/02 §2).
type BinRange struct {
	BIN      string // first 6 PAN digits
	Low      int64  // lower bound of the range (numeric PAN tail)
	High     int64  // upper bound of the range
	Currency string
	Product  string // product / program name
}

// Card is a card record as the issuer holds it: the PAN is never stored in
// clear; a masked form and a SHA-256 hash are kept, plus the per-card key
// used to verify online cryptograms.
type Card struct {
	Ref     string // stable identifier derived from the PAN hash
	PANHash []byte
	PANMask string
	BIN     string
	Expiry  string // YYMM
	Status  string
	Product string
	UDK     []byte // unique derived key
	LastATC uint16 // last accepted ATC (anti-replay)
}

// Config configures the card service.
type Config struct {
	IssuerMasterKey []byte // 16-byte AES key; the "HSM master key"
	Log             *slog.Logger
}

// Service issues cards and verifies cryptograms.
type Service struct {
	cfg   Config
	store Store
	log   *slog.Logger
}

// NewService builds a card service over the given store.
func NewService(store Store, cfg Config) (*Service, error) {
	if len(cfg.IssuerMasterKey) != cmacSize {
		return nil, fmt.Errorf("cardsvc: issuer master key must be %d bytes", cmacSize)
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Service{cfg: cfg, store: store, log: cfg.Log}, nil
}

// AddBinRange registers a BIN range for issuance.
func (s *Service) AddBinRange(ctx context.Context, r BinRange) error {
	if len(r.BIN) != 6 {
		return fmt.Errorf("cardsvc: BIN must be 6 digits, got %q", r.BIN)
	}
	if r.High < r.Low {
		return fmt.Errorf("cardsvc: bin %s high < low", r.BIN)
	}
	return s.store.SaveRange(ctx, r)
}

// FindBinRange returns the range the PAN belongs to, if issued.
func (s *Service) FindBinRange(ctx context.Context, pan string) (BinRange, bool, error) {
	if len(pan) < 7 {
		return BinRange{}, false, nil
	}
	ranges, err := s.store.Ranges(ctx)
	if err != nil {
		return BinRange{}, false, err
	}
	bin := pan[:6]
	tail, err := parseInt64(pan[6:])
	if err != nil {
		return BinRange{}, false, nil
	}
	for _, r := range ranges {
		if r.BIN == bin && tail >= r.Low && tail <= r.High {
			return r, true, nil
		}
	}
	return BinRange{}, false, nil
}

// CreateCard personalizes a new card in an issued BIN range and derives its
// per-card application key (docs/22 §22.2).
func (s *Service) CreateCard(ctx context.Context, pan, expiry, product string) (*Card, error) {
	if !ValidLuhn(pan) {
		return nil, fmt.Errorf("cardsvc: pan %s fails Luhn check", pan)
	}
	if len(expiry) != 4 {
		return nil, fmt.Errorf("cardsvc: expiry must be YYMM")
	}
	r, ok, err := s.FindBinRange(ctx, pan)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cardsvc: pan %s is not in an issued BIN range", pan)
	}

	udk, err := DeriveCardKey(s.cfg.IssuerMasterKey, pan, 1)
	if err != nil {
		return nil, err
	}
	card := Card{
		Ref:     cardRef(pan),
		PANHash: panHash(pan),
		PANMask: MaskPAN(pan),
		BIN:     r.BIN,
		Expiry:  expiry,
		Status:  StatusActive,
		Product: productOr(product, r.Product),
		UDK:     udk,
	}
	if err := s.store.SaveCard(ctx, card); err != nil {
		return nil, fmt.Errorf("cardsvc: persist card: %w", err)
	}
	s.log.Info("card issued", "bin", card.BIN, "mask", card.PANMask, "product", card.Product)
	return &card, nil
}

// GetCard loads a card by its reference.
func (s *Service) GetCard(ctx context.Context, ref string) (*Card, error) {
	card, ok, err := s.store.Card(ctx, ref)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cardsvc: unknown card %s", ref)
	}
	return &card, nil
}

// ComputeARQC simulates the chip: it derives a one-time cryptogram for the
// transaction data using the card's application key.
func (s *Service) ComputeARQC(ctx context.Context, ref string, data ARQCData) ([]byte, error) {
	card, err := s.GetCard(ctx, ref)
	if err != nil {
		return nil, err
	}
	return ComputeARQC(card.UDK, data)
}

// VerifyARQC validates an online cryptogram against the card's key and
// enforces ATC anti-replay, updating the stored ATC on success (docs/06 §6.5).
func (s *Service) VerifyARQC(ctx context.Context, ref string, data ARQCData, arqc []byte) (bool, error) {
	card, err := s.GetCard(ctx, ref)
	if err != nil {
		return false, err
	}
	valid, replay, err := VerifyARQC(card.UDK, data, arqc, card.LastATC)
	if err != nil {
		return false, err
	}
	if valid {
		card.LastATC = data.ATC
		if err := s.store.SaveCard(ctx, *card); err != nil {
			return false, fmt.Errorf("cardsvc: persist atc: %w", err)
		}
	}
	return valid && !replay, nil
}

func panHash(pan string) []byte {
	h := sha256.Sum256([]byte(pan))
	return h[:]
}

func cardRef(pan string) string {
	return hex.EncodeToString(panHash(pan))[:16]
}

// MaskPAN renders the card-holder data protection form, e.g. 400000******7890.
func MaskPAN(pan string) string {
	if len(pan) < 12 {
		return strings.Repeat("*", len(pan)-4) + pan[len(pan)-4:]
	}
	return pan[:6] + strings.Repeat("*", len(pan)-10) + pan[len(pan)-4:]
}

// ValidLuhn checks the ISO/IEC 7812 check digit (docs/02).
func ValidLuhn(pan string) bool {
	if len(pan) < 13 {
		return false
	}
	var sum int
	double := false
	for i := len(pan) - 1; i >= 0; i-- {
		d := int(pan[i] - '0')
		if d < 0 || d > 9 {
			return false
		}
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

func parseInt64(s string) (int64, error) {
	var n int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a digit: %c", c)
		}
		n = n*10 + int64(c-'0')
	}
	return n, nil
}

func productOr(given, fallback string) string {
	if given != "" {
		return given
	}
	if fallback != "" {
		return fallback
	}
	return "classic"
}
