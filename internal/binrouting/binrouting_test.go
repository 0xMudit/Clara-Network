package binrouting

import "testing"

func TestBINOf(t *testing.T) {
	cases := []struct {
		pan  string
		bin  string
		want bool
	}{
		{"4000001234567890", "400000", true},
		{"4000051234567890", "400005", true},
		{"12345", "", false},
		{"40000X123456789", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := BINOf(c.pan)
		if ok != c.want || got != c.bin {
			t.Fatalf("BINOf(%q) = %q,%v want %q,%v", c.pan, got, ok, c.bin, c.want)
		}
	}
}

func TestLookup(t *testing.T) {
	tab, err := FromJSON([]byte(`{"entries":{"400000":"1000001000"},"default":"9999999999"}`))
	if err != nil {
		t.Fatalf("FromJSON: %v", err)
	}

	cases := []struct {
		pan  string
		id   string
		want bool
	}{
		{"4000001234567890", "1000001000", true},
		{"4000051234567890", "9999999999", true},
		{"12345", "", false},
	}
	for _, c := range cases {
		got, ok := tab.Lookup(c.pan)
		if ok != c.want || got != c.id {
			t.Fatalf("Lookup(%q) = %q,%v want %q,%v", c.pan, got, ok, c.id, c.want)
		}
	}
}

func TestLookupNoDefault(t *testing.T) {
	tab, err := FromJSON([]byte(`{"entries":{"400000":"1000001000"}}`))
	if err != nil {
		t.Fatalf("FromJSON: %v", err)
	}
	if _, ok := tab.Lookup("4000051234567890"); ok {
		t.Fatal("expected no route for unknown BIN without default")
	}
}

func TestFromJSONInvalid(t *testing.T) {
	if _, err := FromJSON([]byte(`{not json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
