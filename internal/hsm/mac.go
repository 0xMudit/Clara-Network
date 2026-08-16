package hsm

import (
	"crypto/des"
	"crypto/subtle"
	"errors"
)

// RetailMAC computes an ISO 9797-1 MAC algorithm 3 (ANSI X9.19 retail MAC)
// over data using a 16-byte key treated as two DES keys (K1, K2). Padding is
// method 1 (zeros). The final block is double-encrypted as E_K1(D_K2(C_n)),
// which resists extension attacks that break plain CBC-MAC.
func RetailMAC(key, data []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, errors.New("hsm: retail MAC requires a 16-byte key")
	}
	k1, err := des.NewCipher(key[:8])
	if err != nil {
		return nil, err
	}
	k2, err := des.NewCipher(key[8:])
	if err != nil {
		return nil, err
	}

	padded := zeroPad(data, 8)
	if len(padded) == 0 {
		return nil, errors.New("hsm: empty message has no MAC")
	}

	prev := make([]byte, 8)
	for i := 0; i < len(padded); i += 8 {
		block := padded[i : i+8]
		for j := range block {
			block[j] ^= prev[j]
		}
		k1.Encrypt(block, block)
		copy(prev, block)
	}

	// Retail MAC final step: E_K1(D_K2(C_n)).
	final := make([]byte, 8)
	copy(final, prev)
	k2.Decrypt(final, final)
	k1.Encrypt(final, final)
	return final, nil
}

// VerifyMAC reports whether mac matches the retail MAC of data in constant
// time.
func VerifyMAC(key, data, mac []byte) bool {
	got, err := RetailMAC(key, data)
	if err != nil {
		return false
	}
	if len(got) != len(mac) {
		return false
	}
	return subtle.ConstantTimeCompare(got, mac) == 1
}

func zeroPad(data []byte, size int) []byte {
	pad := size - len(data)%size
	if pad == size {
		pad = 0
	}
	out := make([]byte, len(data)+pad)
	copy(out, data)
	return out
}
