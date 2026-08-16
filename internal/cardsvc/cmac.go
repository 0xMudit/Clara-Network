package cardsvc

import (
	"crypto/aes"
	"fmt"
)

// cmac is AES-CMAC (NIST SP 800-38B) with a 128-bit key and a 128-bit tag.
// It is used to derive per-card keys from the issuer master key and to compute
// the EMV application cryptograms.

const cmacSize = aes.BlockSize

func dbl(b []byte) []byte {
	out := make([]byte, cmacSize)
	var carry byte
	for i := cmacSize - 1; i >= 0; i-- {
		next := b[i] >> 7
		out[i] = b[i]<<1 | carry
		carry = next
	}
	if carry != 0 {
		out[cmacSize-1] ^= 0x87
	}
	return out
}

// cmac returns the AES-CMAC tag of msg under a 128-bit key.
func cmac(key, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cmac: %w", err)
	}
	L := make([]byte, cmacSize)
	block.Encrypt(L, make([]byte, cmacSize))
	k1 := dbl(L)
	k2 := dbl(k1)

	n := (len(msg) + cmacSize - 1) / cmacSize
	if n == 0 {
		n = 1
	}
	x := make([]byte, cmacSize)
	rest := msg
	for i := 0; i < n; i++ {
		blockIn := make([]byte, cmacSize)
		switch {
		case i == n-1 && len(rest) == cmacSize:
			copy(blockIn, rest)
			xorBytes(blockIn, k1)
		case i == n-1:
			copy(blockIn, rest)
			blockIn[len(rest)] = 0x80
			xorBytes(blockIn, k2)
		default:
			copy(blockIn, rest[:cmacSize])
			rest = rest[cmacSize:]
		}
		xorBytes(x, blockIn)
		block.Encrypt(x, x)
	}
	return x, nil
}

func xorBytes(dst, src []byte) {
	for i := range src {
		dst[i] ^= src[i]
	}
}
