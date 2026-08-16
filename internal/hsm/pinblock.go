package hsm

import (
	"crypto/aes"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
)

// PIN block formats (ISO 9564-1).
const (
	PinFormat0 = 0 // ISO-0, PAN-derived, 3DES
	PinFormat4 = 4 // ISO 9564-1:2017, AES, combines PIN and PAN
)

// format0Block builds the 8-byte ISO-0 block: format nibble 0x0, PIN length,
// PIN digits, F fill, XORed with the PAN-derived data.
func format0Block(pin, pan string) ([]byte, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return nil, errors.New("hsm: PIN must be 4-12 digits")
	}
	p := make([]byte, 8)
	p[0] = byte(0x0<<4 | len(pin)) // format 0, PIN length
	digits := pin + strings.Repeat("F", 14-len(pin))
	for i, d := range digits {
		if i%2 == 0 {
			p[1+i/2] |= digitValue(byte(d)) << 4
		} else {
			p[1+i/2] |= digitValue(byte(d))
		}
	}
	panData, err := panBytes(pan)
	if err != nil {
		return nil, err
	}
	for i := range p {
		p[i] ^= panData[i]
	}
	return p, nil
}

// format4Block builds the 16-byte ISO 9564-1:2017-style block: format nibble
// 0x4, PIN length, PIN digits, F fill, XORed with the PAN-derived data, then
// AES-encrypted with the PIN key.
func format4Block(pinKey, pin, pan string) ([]byte, error) {
	if len(pin) < 4 || len(pin) > 15 {
		return nil, errors.New("hsm: PIN must be 4-15 digits for format 4")
	}
	key, err := hexBytes(pinKey)
	if err != nil {
		return nil, err
	}
	if len(key) != 16 {
		return nil, errors.New("hsm: format 4 requires a 16-byte AES PIN key")
	}
	p := make([]byte, 16)
	p[0] = byte(0x4<<4 | len(pin))
	digits := pin + strings.Repeat("F", 30-len(pin))
	for i, d := range digits {
		if i%2 == 0 {
			p[1+i/2] |= digitValue(byte(d)) << 4
		} else {
			p[1+i/2] |= digitValue(byte(d))
		}
	}
	panData, err := panBytes16(pan)
	if err != nil {
		return nil, err
	}
	for i := range p {
		p[i] ^= panData[i]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ct := make([]byte, 16)
	block.Encrypt(ct, p)
	return ct, nil
}

// decryptFormat4 reverses format4Block and returns the plaintext block.
func decryptFormat4(pinKey string, ct []byte) ([]byte, error) {
	key, err := hexBytes(pinKey)
	if err != nil {
		return nil, err
	}
	if len(ct) != 16 {
		return nil, errors.New("hsm: format 4 block must be 16 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	pt := make([]byte, 16)
	block.Decrypt(pt, ct)
	return pt, nil
}

// panBytes derives the ISO-0 PAN block: the rightmost 12 PAN digits excluding
// the check digit, left-padded with zeros to 16 nibbles.
func panBytes(pan string) ([]byte, error) {
	return packDigits(padding12(pan))
}

// panBytes16 derives the format-4 PAN block: the rightmost 12 PAN digits
// excluding the check digit, left-padded with zeros to 32 nibbles.
func panBytes16(pan string) ([]byte, error) {
	return packDigits(padding32(pan))
}

func padding12(pan string) string {
	core := panCore(pan)
	return strings.Repeat("0", 16-len(core)) + core
}

func padding32(pan string) string {
	core := panCore(pan)
	return strings.Repeat("0", 32-len(core)) + core
}

// panCore returns the rightmost 12 PAN digits excluding the check digit.
func panCore(pan string) string {
	d := ""
	for _, r := range pan {
		if r >= '0' && r <= '9' {
			d += string(r)
		}
	}
	if len(d) < 13 {
		return d
	}
	// Rightmost 12 digits before the Luhn check digit.
	return d[len(d)-13 : len(d)-1]
}

func packDigits(digits string) ([]byte, error) {
	if len(digits)%2 != 0 {
		return nil, fmt.Errorf("hsm: odd digit count %d", len(digits))
	}
	out := make([]byte, len(digits)/2)
	for i := 0; i < len(digits); i++ {
		if i%2 == 0 {
			out[i/2] |= digitValue(digits[i]) << 4
		} else {
			out[i/2] |= digitValue(digits[i])
		}
	}
	return out, nil
}

func digitValue(d byte) byte {
	switch {
	case d >= '0' && d <= '9':
		return d - '0'
	case d == 'F' || d == 'f':
		return 0xF
	default:
		return 0xF
	}
}

// pinFromFormat0 extracts the PIN from an ISO-0 block given the PAN.
func pinFromFormat0(block []byte, pan string) (string, error) {
	panData, err := panBytes(pan)
	if err != nil {
		return "", err
	}
	clear := make([]byte, 8)
	for i := range block {
		clear[i] = block[i] ^ panData[i]
	}
	if clear[0]>>4 != 0 {
		return "", fmt.Errorf("hsm: not an ISO-0 block (format nibble %x)", clear[0]>>4)
	}
	n := int(clear[0] & 0x0F)
	if n < 4 || n > 12 {
		return "", fmt.Errorf("hsm: bad PIN length %d", n)
	}
	var b strings.Builder
	for i := 1; i < 1+n/2+1 && b.Len() < n; i++ {
		hi, lo := nibble(clear[i]>>4), nibble(clear[i]&0x0F)
		b.WriteByte(hi)
		b.WriteByte(lo)
	}
	return b.String()[:n], nil
}

func nibble(v byte) byte {
	if v <= 9 {
		return '0' + v
	}
	return 'F'
}

// pinFromFormat4 decrypts a format 4 block and extracts the PIN.
func pinFromFormat4(key string, ct []byte) (string, error) {
	pt, err := decryptFormat4(key, ct)
	if err != nil {
		return "", err
	}
	if pt[0]>>4 != 4 {
		return "", fmt.Errorf("hsm: not a format 4 block (format nibble %x)", pt[0]>>4)
	}
	n := int(pt[0] & 0x0F)
	var b strings.Builder
	for i := 1; i < 16 && b.Len() < n; i++ {
		hi, lo := nibble(pt[i]>>4), nibble(pt[i]&0x0F)
		b.WriteByte(hi)
		b.WriteByte(lo)
	}
	return b.String()[:n], nil
}

// verifyPIN compares a decrypted PIN against the expected value in constant
// time.
func verifyPIN(got, want string) bool {
	if len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}
