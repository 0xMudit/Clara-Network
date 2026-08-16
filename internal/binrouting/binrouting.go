// Package binrouting maps card BINs (the first 6 digits of a PAN) to the
// receiving institution (issuer) that owns them. The switch uses it to derive
// the destination when a message does not carry DE100.
package binrouting

import "encoding/json"

// Config is the serialisable BIN routing table.
type Config struct {
	// Entries maps a BIN prefix (usually 6 digits) to an issuer ID (DE100).
	Entries map[string]string `json:"entries"`
	// Default is applied when no entry matches, if non-empty.
	Default string `json:"default,omitempty"`
}

// Table resolves a PAN BIN to an issuer ID.
type Table struct {
	entries map[string]string
	def     string
}

// New builds a table from config. A nil config produces an empty table that
// never routes.
func New(cfg *Config) *Table {
	t := &Table{entries: map[string]string{}}
	if cfg == nil {
		return t
	}
	for k, v := range cfg.Entries {
		t.entries[k] = v
	}
	t.def = cfg.Default
	return t
}

// FromJSON parses a JSON config document.
func FromJSON(data []byte) (*Table, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return New(&cfg), nil
}

// BINOf returns the first 6 digits of a PAN. ok is false when the PAN is
// shorter than 6 digits or contains non-numeric characters.
func BINOf(pan string) (string, bool) {
	if len(pan) < 6 {
		return "", false
	}
	bin := pan[:6]
	for i := 0; i < len(bin); i++ {
		if bin[i] < '0' || bin[i] > '9' {
			return "", false
		}
	}
	return bin, true
}

// Lookup resolves a PAN to an issuer ID.
func (t *Table) Lookup(pan string) (string, bool) {
	bin, ok := BINOf(pan)
	if !ok {
		return "", false
	}
	if id, ok := t.entries[bin]; ok {
		return id, true
	}
	if t.def != "" {
		return t.def, true
	}
	return "", false
}
