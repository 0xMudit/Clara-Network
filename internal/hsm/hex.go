package hsm

import "encoding/hex"

func hexBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
