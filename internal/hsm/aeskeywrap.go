package hsm

import (
	"crypto/aes"
	"crypto/subtle"
	"errors"
	"fmt"
)

// AES Key Wrap (RFC 3394) wraps key material for storage and transport. The
// NIST test vector below pins the implementation (RFC 3394 §4.5).
const (
	rfc3394KEK     = "000102030405060708090a0b0c0d0e0f"
	rfc3394Key     = "00112233445566778899aabbccddeeff"
	rfc3394Wrapped = "1fa68b0a8112b447aef34bd8fb5a7b829d3e862371d2cfe5"
)

// wrapAESKey wraps a 16-byte key with a 16-byte KEK using AES Key Wrap
// (RFC 3394). The returned block is len(key)+8 bytes.
func wrapAESKey(kek, key []byte) ([]byte, error) {
	if len(kek) != 16 || len(key)%8 != 0 || len(key) < 16 {
		return nil, errors.New("hsm: key wrap requires a 16-byte KEK and an 8-multiple key of at least 16 bytes")
	}
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	n := len(key) / 8
	a := make([]byte, 8)
	copy(a, []byte{0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6})
	r := make([]byte, len(key))
	copy(r, key)

	buf := make([]byte, 16)
	for j := 0; j < 6; j++ {
		for i := 1; i <= n; i++ {
			copy(buf[0:8], a)
			copy(buf[8:16], r[(i-1)*8:i*8])
			block.Encrypt(buf, buf)
			t := n*j + i
			// XOR A with the 64-bit integer t (big-endian) in place.
			buf[7] ^= byte(t)
			buf[6] ^= byte(t >> 8)
			if t >= 1<<16 {
				buf[5] ^= byte(t >> 16)
				buf[4] ^= byte(t >> 24)
			}
			copy(a, buf[0:8])
			copy(r[(i-1)*8:i*8], buf[8:16])
		}
	}
	out := make([]byte, 0, len(key)+8)
	out = append(out, a...)
	out = append(out, r...)
	return out, nil
}

// unwrapAESKey reverses wrapAESKey (RFC 3394).
func unwrapAESKey(kek, block []byte) ([]byte, error) {
	if len(kek) != 16 || len(block)%8 != 0 || len(block) < 24 {
		return nil, errors.New("hsm: key unwrap requires a 16-byte KEK and a valid wrapped block")
	}
	blockCipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	n := len(block)/8 - 1
	a := make([]byte, 8)
	copy(a, block[0:8])
	r := make([]byte, len(block)-8)
	copy(r, block[8:])

	buf := make([]byte, 16)
	for j := 5; j >= 0; j-- {
		for i := n; i >= 1; i-- {
			t := n*j + i
			// XOR the round counter into A (big-endian) before decrypting.
			a[7] ^= byte(t)
			a[6] ^= byte(t >> 8)
			if t >= 1<<16 {
				a[5] ^= byte(t >> 16)
				a[4] ^= byte(t >> 24)
			}
			copy(buf[0:8], a)
			copy(buf[8:16], r[(i-1)*8:i*8])
			blockCipher.Decrypt(buf, buf)
			copy(a, buf[0:8])
			copy(r[(i-1)*8:i*8], buf[8:16])
		}
	}
	for _, b := range a {
		if b != 0xA6 {
			return nil, errors.New("hsm: integrity check failed while unwrapping key (wrong KEK or tampered block)")
		}
	}
	return r, nil
}

// keyWrapKnownAnswer verifies the RFC 3394 §4.5 vector.
func keyWrapKnownAnswer() error {
	kek, err := hexBytes(rfc3394KEK)
	if err != nil {
		return err
	}
	key, err := hexBytes(rfc3394Key)
	if err != nil {
		return err
	}
	want, err := hexBytes(rfc3394Wrapped)
	if err != nil {
		return err
	}
	got, err := wrapAESKey(kek, key)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return fmt.Errorf("hsm: RFC 3394 wrap KAT mismatch: got %x want %x", got, want)
	}
	back, err := unwrapAESKey(kek, got)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(back, key) != 1 {
		return fmt.Errorf("hsm: RFC 3394 unwrap round-trip mismatch")
	}
	return nil
}
