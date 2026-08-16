package iso8583

import (
	"reflect"
	"strings"
	"testing"
)

func TestRoundTripAllFields(t *testing.T) {
	m := New("0100").
		Set(2, "4000001234567890").
		Set(3, "000000").
		Set(4, "000000100000").
		Set(7, "0816142030").
		Set(11, "123456").
		Set(12, "142030").
		Set(13, "0816").
		Set(22, "022").
		Set(24, "100").
		Set(25, "00").
		Set(32, "1000001").
		Set(37, "ABCDEFGHIJKL").
		Set(39, "00").
		Set(41, "TST00001").
		Set(42, "CLARA000000001").
		Set(49, "840").
		Set(62, "SI").
		Set(70, "301").
		Set(100, "1000001000")

	raw, err := m.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got, err := Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.MTI != m.MTI {
		t.Fatalf("MTI mismatch: want %q got %q", m.MTI, got.MTI)
	}
	if !reflect.DeepEqual(got.Fields, m.Fields) {
		t.Fatalf("fields mismatch:\nwant %#v\ngot  %#v", m.Fields, got.Fields)
	}
}

func TestSecondaryBitmap(t *testing.T) {
	m := New("0100").Set(100, "1000001000")
	raw, err := m.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// 4 (MTI) + 16 (primary with secondary flag) + 16 (secondary) + 2 + 10 (DE100).
	if len(raw) != 48 {
		t.Fatalf("unexpected wire length %d, want 48", len(raw))
	}
	if raw[4] != '8' {
		t.Fatalf("expected secondary bitmap flag (leading 8), got %q", raw[:16])
	}
}

func TestPadding(t *testing.T) {
	m := New("0100").Set(4, "123").Set(41, "AB")
	raw, err := m.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got, err := Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.Get(4) != "000000000123" {
		t.Fatalf("DE4 not zero padded: %q", got.Get(4))
	}
	// Alpha fields are space padded on the wire but trimmed on read.
	if got.Get(41) != "AB" {
		t.Fatalf("DE41 not trimmed on read: %q", got.Get(41))
	}
	if !strings.Contains(string(raw), "AB      ") {
		t.Fatalf("wire form should space-pad DE41, got %q", raw)
	}
}

func TestRejectsInvalidMTI(t *testing.T) {
	if _, err := New("12AB").Marshal(); err == nil {
		t.Fatal("expected error for non-numeric MTI")
	}
}

func TestRejectsOverlongValue(t *testing.T) {
	if _, err := New("0100").Set(3, "1234567").Marshal(); err == nil {
		t.Fatal("expected error for overlong DE3")
	}
}

func TestParseTooShort(t *testing.T) {
	if _, err := Parse([]byte("0100")); err == nil {
		t.Fatal("expected error for short message")
	}
}

func TestParseEmptyBitmapFieldError(t *testing.T) {
	// DE100 declared in bitmap but truncated body.
	raw := "0100800000000000000000000000000001"
	if _, err := Parse([]byte(raw)); err == nil {
		t.Fatal("expected error for truncated variable field")
	}
}
