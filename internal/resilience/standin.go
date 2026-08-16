// Package resilience implements the Clara Network operational-resilience
// layer (docs/19): issuer stand-in processing (SIP/STIP) with negative and
// valid-card files, per-route circuit breakers with half-open probing,
// outcome metrics with approximate p99 latency, and burst detection for
// issuer-inoperative (91) responses that flag an issuer outage.
package resilience

import (
	"log/slog"
	"sync"
)

// Policy configures stand-in processing for a single issuer (docs/19 §19.3).
// The network approves transactions within the issuer's stored limits when
// the issuer cannot be reached and declines against negative files.
type Policy struct {
	IssuerID       string
	Enabled        bool
	Limit          int64           // max amount (minor units) approved in stand-in
	NegativeCards  map[string]bool // hot / lost / stolen PANs: always decline
	RestrictedBINs map[string]bool // BINs not served in stand-in
	ValidCards     map[string]bool // positive file; when non-nil only listed PANs are approved
}

// Decision is the stand-in outcome for a transaction.
type Decision struct {
	Approve bool
	Code    string
}

// StandIn evaluates transactions on behalf of unreachable issuers. Fallback
// ordering is primary -> secondary route -> stand-in -> decline.
type StandIn struct {
	mu           sync.Mutex
	defaultLimit int64
	policies     map[string]*Policy
	log          *slog.Logger
}

// NewStandIn builds a stand-in engine with a default approval limit (minor
// units) applied to issuers without an explicit policy.
func NewStandIn(defaultLimit int64) *StandIn {
	return &StandIn{
		defaultLimit: defaultLimit,
		policies:     map[string]*Policy{},
		log:          slog.Default(),
	}
}

// SetPolicy registers or replaces an issuer's stand-in policy.
func (s *StandIn) SetPolicy(p Policy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.Limit == 0 {
		p.Limit = s.defaultLimit
	}
	s.policies[p.IssuerID] = &p
}

// Policy returns a copy of the current policy for an issuer.
func (s *StandIn) Policy(issuerID string) (Policy, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.policies[issuerID]
	if !ok {
		return Policy{}, false
	}
	cp := *p
	cp.NegativeCards = cloneSet(p.NegativeCards)
	cp.RestrictedBINs = cloneSet(p.RestrictedBINs)
	cp.ValidCards = cloneSet(p.ValidCards)
	return cp, true
}

// Decide applies stand-in rules for a transaction destined for an unreachable
// issuer. Code 91 means issuer/switch inoperative (no stand-in applicable),
// 05 means do not honor (negative file), 57 transaction not permitted
// (restricted BIN). Approvals carry the SI marker set by the caller.
func (s *StandIn) Decide(issuerID, pan string, amountMinor int64) Decision {
	s.mu.Lock()
	defer s.mu.Unlock()
	pol, ok := s.policies[issuerID]
	if !ok {
		pol = &Policy{IssuerID: issuerID, Enabled: true, Limit: s.defaultLimit}
	}
	if !pol.Enabled {
		return Decision{Code: "91"} // stand-in not applicable for this issuer
	}
	if pol.ValidCards != nil && !pol.ValidCards[pan] {
		return Decision{Code: "91"} // positive file: card not eligible for stand-in
	}
	if pol.NegativeCards[pan] {
		return Decision{Code: "05"} // do not honor
	}
	if pol.RestrictedBINs[binOf(pan)] {
		return Decision{Code: "57"} // transaction not permitted
	}
	if amountMinor > pol.Limit {
		return Decision{Code: "91"} // above stand-in limit: decline
	}
	return Decision{Approve: true, Code: "00"}
}

// Issuers returns the issuer IDs that have an explicit stand-in policy.
func (s *StandIn) Issuers() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.policies))
	for id := range s.policies {
		out = append(out, id)
	}
	return out
}

func cloneSet(m map[string]bool) map[string]bool {
	if m == nil {
		return nil
	}
	out := make(map[string]bool, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func binOf(pan string) string {
	for i := 0; i < len(pan); i++ {
		if pan[i] < '0' || pan[i] > '9' {
			return ""
		}
	}
	if len(pan) < 6 {
		return pan
	}
	return pan[:6]
}
