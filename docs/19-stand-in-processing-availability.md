# 19. Stand-In Processing & Operational Resilience

## 19.1 The stand-in problem

When an **issuer host is unavailable** (outage, timeout, switch inoperative),
merchants must not be left hanging at the point of sale. Card networks respond
on the issuer's behalf with **stand-in processing**: the network authorizes
(or declines) the transaction using stored rules and risk limits instead of
the issuer.

## 19.2 How the major schemes do it

### Mastercard — Stand-In Processing (SIP)

- Documented in the *Mastercard Switch Rules* manual (§5.4.2 Stand-In
  Processing Service).
- The network approves transactions within the issuer's defined stand-in
  limits (amount thresholds, risk rules, valid/invalid card lists) when the
  issuer cannot be reached.
- Transactions processed in stand-in are flagged so the issuer can reconcile
  and, if needed, mark them for review.

### Visa — STIP (Stand-In Processing)

- Visa's STIP authorizes transactions for issuers that are offline, subject to
  issuer-configured **STIP limits** and negative files.
- Response code `91` — **Issuer or switch is inoperative**; also used when STIP
  is **not applicable/available** for the transaction.
- Response code `P` — **CVV2 not performed** because STIP (stand-in) responded
  to an issuer-unavailable condition.
- Key risk signal: a burst of `91` responses typically indicates an issuer
  outage and triggers monitoring.

### Stand-in risk model

- **Approve-within-limits** — approve only low-risk, low-value transactions;
  decline high-value or high-risk ones.
- **Decline by default** — safer but poor customer experience and merchant
  friction.
- **Hybrid** — approve within limits based on card/merchant risk profile,
  decline everything above thresholds.

## 19.3 Designing your own stand-in

If you operate the network:

- **Rules engine** (see `10-fraud-risk.md`) evaluates transactions when the
  issuer is unreachable, using stored card-level limits and velocity.
- **Negative files** — hot cards, stolen cards, restricted BINs enforced during
  stand-in.
- **Stand-in markers** — add a flag/code so downstream clearing, disputes, and
  reporting know the transaction was stand-in approved.
- **Fallback ordering** — primary → secondary issuer routing → stand-in →
  decline. Track timeouts per route.

## 19.4 Operational resilience for the network itself

PFMI Principle 17 (operational risk) requires high availability, scalability,
and business-continuity planning.

### Availability targets & design

- Target **99.99%+** availability for authorization; design for graceful
  degradation (never a hard crash).
- **Active-active multi-site** — two (or more) independently operating data
  centers, geographically dispersed; traffic fails over automatically.
- **No single points of failure** — redundant network links, HSMs (see
  `17-message-security-key-management.md`), switches, databases.
- **Graceful degradation** — offload non-critical load (fraud scoring,
  analytics) during stress; keep the core authorization path up.

### Business continuity planning (BCP/DR)

- **RTO / RPO** targets (e.g., RTO ≤ 1 h, RPO ≈ 0 for the ledger).
- **Cyber resilience** — per BIS CPMI guidance (d146): defend, detect,
  respond, recover; test response playbooks.
- **Failover testing** — regular chaos/recovery drills; load tests at peak and
  above-peak volumes (see `12-system-design.md`).
- **Monitoring & alerting** — heartbeat, latency, error-rate and response-code
  dashboards; burst detection (e.g., spike in `91`s).

## 19.5 Key references

- Mastercard — *Switch Rules* (Stand-In Processing Service, §5.4.2).
- Visa — *International Operating Regulations* (STIP, response codes 91 / P).
- ISO 8583 — response/action code semantics for issuer-unavailable conditions
  (see `04-iso8583.md`).
- BIS CPMI — *Guidance on cyber resilience for FMIs* (d146).
