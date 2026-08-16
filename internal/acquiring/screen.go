package acquiring

import (
	"context"
	"strings"
)

// List names for screening lists.
const (
	ListMATCH = "MATCH"
	ListOFAC  = "OFAC"
)

// MatchEntry is a record in the Member Alert to Control High-Risk Merchants
// list (docs/23 §23.2): merchants terminated for cause by another acquirer.
type MatchEntry struct {
	MerchantName string
	TaxID        string
	Reason       string
}

// OfacEntry is a record in the OFAC SDN / sanctions list.
type OfacEntry struct {
	Name    string
	Program string
}

// Hit is a screening match that blocks boarding.
type Hit struct {
	List   string
	Detail string
}

// Screener checks merchant applications against negative lists.
type Screener struct {
	Store Store
}

// NewScreener returns a screener over the store that holds the lists.
func NewScreener(store Store) *Screener {
	return &Screener{Store: store}
}

// Screen checks the application principals and identity against the OFAC SDN
// and MATCH lists. Any hit blocks boarding.
func (s *Screener) Screen(ctx context.Context, app Application) ([]Hit, error) {
	var hits []Hit

	matches, err := s.Store.MatchEntries(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		if eqFold(m.MerchantName, app.MerchantName) || (m.TaxID != "" && m.TaxID == app.TaxID) {
			hits = append(hits, Hit{List: ListMATCH, Detail: m.MerchantName + ": " + m.Reason})
		}
	}

	ofac, err := s.Store.OfacEntries(ctx)
	if err != nil {
		return nil, err
	}
	names := append([]string{app.MerchantName}, app.Principals...)
	for _, e := range ofac {
		for _, n := range names {
			if eqFold(e.Name, n) {
				hits = append(hits, Hit{List: ListOFAC, Detail: e.Name + " (" + e.Program + ")"})
			}
		}
	}
	return hits, nil
}

func eqFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
