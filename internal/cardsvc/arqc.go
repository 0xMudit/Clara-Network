package cardsvc

import (
	"crypto/subtle"
	"encoding/binary"
	"fmt"
)

// ARQCData is the transaction data the chip authenticates with an EMV-style
// Authorization Request Cryptogram (docs/06 §6.4–6.5). The fields are
// concatenated into a canonical byte string; a per-card derived key turns it
// into a one-time cryptogram that proves card authenticity per transaction.
type ARQCData struct {
	Amount   int64  // minor units, encoded as 12 ASCII digits
	Currency string // 3-digit numeric currency code, e.g. "840"
	STAN     string // 6-digit system trace audit number
	Date     string // YYMMDD transaction date
	ATC      uint16 // application transaction counter, anti-replay
	UN       string // unpredictable number from the terminal, 4 digits
}

// Bytes returns the canonical concatenation used for the cryptogram.
func (d ARQCData) Bytes() []byte {
	out := make([]byte, 0, 33)
	out = append(out, []byte(fmt.Sprintf("%012d", d.Amount))...)
	out = append(out, []byte(fmt.Sprintf("%03s", d.Currency))...)
	out = append(out, []byte(fmt.Sprintf("%06s", d.STAN))...)
	out = append(out, []byte(fmt.Sprintf("%06s", d.Date))...)
	out = binary.BigEndian.AppendUint16(out, d.ATC)
	out = append(out, []byte(fmt.Sprintf("%04s", d.UN))...)
	return out
}

// DeriveCardKey derives the per-card application key from the issuer master
// key and the card's PAN plus an issuance sequence number, mirroring how a
// personalization system loads card keys at production time (docs/22 §22.2).
func DeriveCardKey(issuerMasterKey []byte, pan string, sequence byte) ([]byte, error) {
	if len(issuerMasterKey) != cmacSize {
		return nil, fmt.Errorf("cardsvc: issuer master key must be %d bytes", cmacSize)
	}
	if len(pan) == 0 {
		return nil, fmt.Errorf("cardsvc: empty pan")
	}
	in := make([]byte, 0, len(pan)+1)
	in = append(in, []byte(pan)...)
	in = append(in, sequence)
	return cmac(issuerMasterKey, in)
}

// ComputeARQC produces the one-time cryptogram for the transaction data. In a
// real deployment this runs inside the chip; the issuer verifies it with an
// HSM. The sim exposes both sides so the full flow is testable.
func ComputeARQC(cardKey []byte, data ARQCData) ([]byte, error) {
	tag, err := cmac(cardKey, data.Bytes())
	if err != nil {
		return nil, fmt.Errorf("cardsvc: compute arqc: %w", err)
	}
	return tag, nil
}

// VerifyARQC recomputes the cryptogram over the transaction data and compares
// it in constant time (EMV online cryptogram validation). replay indicates
// whether the cryptogram was authentic but its ATC had already been used,
// which is a replay attempt.
func VerifyARQC(cardKey []byte, data ARQCData, arqc []byte, lastATC uint16) (valid, replay bool, err error) {
	expected, err := ComputeARQC(cardKey, data)
	if err != nil {
		return false, false, err
	}
	if len(arqc) != len(expected) {
		return false, false, nil
	}
	if subtle.ConstantTimeCompare(arqc, expected) != 1 {
		return false, false, nil
	}
	if data.ATC <= lastATC {
		return false, true, nil
	}
	return true, false, nil
}
