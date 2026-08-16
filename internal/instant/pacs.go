package instant

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// marshalDoc renders an ISO 20022 document with an XML declaration.
func marshalDoc(doc any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	buf.Write(enc)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// formatAmount renders minor units as major units with two decimals.
func formatAmount(minor int64) string {
	sign := ""
	if minor < 0 {
		sign = "-"
		minor = -minor
	}
	return fmt.Sprintf("%s%d.%02d", sign, minor/100, minor%100)
}

// parseAmount parses "major.minor" into minor units, rejecting more than two
// decimal places.
func parseAmount(s string) (int64, error) {
	sign := int64(1)
	switch s[0] {
	case '-':
		sign = -1
		s = s[1:]
	case '+':
		s = s[1:]
	}
	dot := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dot = i
			break
		}
	}
	var major, frac string
	if dot < 0 {
		major, frac = s, ""
	} else {
		major, frac = s[:dot], s[dot+1:]
	}
	if len(frac) > 2 {
		return 0, fmt.Errorf("more than two decimals in %q", s)
	}
	for len(frac) < 2 {
		frac += "0"
	}
	if !isDigits(major) || !isDigits(frac) {
		return 0, fmt.Errorf("invalid amount %q", s)
	}
	var out int64
	for _, c := range major {
		out = out*10 + int64(c-'0')
	}
	for _, c := range frac {
		out = out*10 + int64(c-'0')
	}
	return sign * out, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
