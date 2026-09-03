package iso8583

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Message is an ISO 8583 message: a 4-digit MTI plus the present data elements.
type Message struct {
	MTI    string
	Fields map[int]string
}

// New creates an empty message with the given MTI.
func New(mti string) *Message {
	return &Message{MTI: mti, Fields: make(map[int]string)}
}

// Set stores a value for a data element.
func (m *Message) Set(field int, value string) *Message {
	m.Fields[field] = value
	return m
}

// Get returns the value of a data element, or "" if absent.
func (m *Message) Get(field int) string { return m.Fields[field] }

// Has reports whether a data element is present.
func (m *Message) Has(field int) bool {
	_, ok := m.Fields[field]
	return ok
}

// Marshal encodes the message to its wire form: MTI, bitmap(s), then elements.
func (m *Message) Marshal() ([]byte, error) {
	if len(m.MTI) != 4 || !allDigits(m.MTI) {
		return nil, fmt.Errorf("iso8583: invalid MTI %q", m.MTI)
	}

	var primary, secondary uint64
	fields := make([]int, 0, len(m.Fields))
	for f := range m.Fields {
		if f < 2 || f > 128 {
			return nil, fmt.Errorf("iso8583: data element %d out of range", f)
		}
		fields = append(fields, f)
		if f > 64 {
			secondary |= 1 << (128 - f)
		} else {
			primary |= 1 << (64 - f)
		}
	}
	sort.Ints(fields)

	if secondary != 0 {
		primary |= 1 << 63
	}

	var b strings.Builder
	b.Grow(len(m.MTI) + 32 + len(fields)*4)
	b.WriteString(m.MTI)
	fmt.Fprintf(&b, "%016X", primary)
	if secondary != 0 {
		fmt.Fprintf(&b, "%016X", secondary)
	}
	for _, f := range fields {
		enc, err := encodeField(m.Fields[f], DataElements[f])
		if err != nil {
			return nil, fmt.Errorf("iso8583: field %d: %w", f, err)
		}
		b.WriteString(enc)
	}
	return []byte(b.String()), nil
}

// Parse decodes an ISO 8583 message from its wire form.
func Parse(data []byte) (*Message, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("iso8583: message too short (%d bytes)", len(data))
	}
	s := string(data)

	mti := s[:4]
	if !allDigits(mti) {
		return nil, fmt.Errorf("iso8583: invalid MTI %q", mti)
	}

	primary, err := strconv.ParseUint(s[4:20], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("iso8583: invalid primary bitmap: %w", err)
	}

	m := New(mti)
	offset := 20

	var secondary uint64
	if primary&(1<<63) != 0 {
		if len(s) < offset+16 {
			return nil, io.ErrUnexpectedEOF
		}
		if secondary, err = strconv.ParseUint(s[offset:offset+16], 16, 64); err != nil {
			return nil, fmt.Errorf("iso8583: invalid secondary bitmap: %w", err)
		}
		offset += 16
	}

	for f := 2; f <= 128; f++ {
		var present bool
		if f <= 64 {
			present = primary&(1<<(64-f)) != 0
		} else {
			present = secondary&(1<<(128-f)) != 0
		}
		if !present {
			continue
		}
		spec, ok := DataElements[f]
		if !ok {
			return nil, fmt.Errorf("iso8583: unsupported data element %d", f)
		}
		val, n, err := readField(s, offset, spec)
		if err != nil {
			return nil, fmt.Errorf("iso8583: field %d: %w", f, err)
		}
		m.Fields[f] = val
		offset += n
	}
	return m, nil
}

func encodeField(value string, spec FieldSpec) (string, error) {
	switch spec.Type {
	case TypeN, TypeA, TypeAN, TypeANS:
		if len(value) > spec.MaxLen {
			return "", fmt.Errorf("value %q exceeds length %d", value, spec.MaxLen)
		}
		if spec.Type == TypeN {
			return leftPad(value, spec.MaxLen, '0'), nil
		}
		return rightPad(value, spec.MaxLen, ' '), nil
	case TypeNLLVAR, TypeNLLLVAR:
		if !allDigits(value) {
			return "", fmt.Errorf("value %q is not numeric", value)
		}
		prefix := 2
		if spec.Type == TypeNLLLVAR {
			prefix = 3
		}
		return encodeVariable(value, prefix, spec.MaxLen)
	case TypeANLLVAR, TypeANLLLVAR:
		prefix := 2
		if spec.Type == TypeANLLLVAR {
			prefix = 3
		}
		return encodeVariable(value, prefix, spec.MaxLen)
	default:
		return "", fmt.Errorf("unknown field type %d", spec.Type)
	}
}

func encodeVariable(value string, prefixLen, maxLen int) (string, error) {
	if len(value) > maxLen {
		return "", fmt.Errorf("value exceeds max length %d", maxLen)
	}
	return fmt.Sprintf("%0*d%s", prefixLen, len(value), value), nil
}

func readField(s string, off int, spec FieldSpec) (string, int, error) {
	switch spec.Type {
	case TypeN, TypeA, TypeAN, TypeANS:
		if off+spec.MaxLen > len(s) {
			return "", 0, io.ErrUnexpectedEOF
		}
		v := s[off : off+spec.MaxLen]
		if spec.Type != TypeN {
			v = strings.TrimRight(v, " ")
		}
		return v, spec.MaxLen, nil
	case TypeNLLVAR, TypeANLLVAR:
		return readVariable(s, off, 2, spec.MaxLen)
	case TypeNLLLVAR, TypeANLLLVAR:
		return readVariable(s, off, 3, spec.MaxLen)
	default:
		return "", 0, fmt.Errorf("unknown field type %d", spec.Type)
	}
}

func readVariable(s string, off, prefixLen, maxLen int) (string, int, error) {
	if off+prefixLen > len(s) {
		return "", 0, io.ErrUnexpectedEOF
	}
	length, err := strconv.Atoi(s[off : off+prefixLen])
	if err != nil {
		return "", 0, fmt.Errorf("invalid length prefix: %w", err)
	}
	if length > maxLen {
		return "", 0, fmt.Errorf("length %d exceeds max %d", length, maxLen)
	}
	if off+prefixLen+length > len(s) {
		return "", 0, io.ErrUnexpectedEOF
	}
	return s[off+prefixLen : off+prefixLen+length], prefixLen + length, nil
}

func allDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func leftPad(s string, n int, c byte) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat(string(c), n-len(s)) + s
}

func rightPad(s string, n int, c byte) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(string(c), n-len(s))
}
