# Clara Network

**Clara Network** is an open-source project to design and build a
Mastercard/Visa-style card payment network end-to-end: scheme (network)
operator, issuer, and acquirer infrastructure.

## What's in this repository

The repo holds the **research and specification library** for the network under
[`docs/`](./docs/), plus the Go implementation of the switch and settlement
engine described in [`docs/25-clara-network-system-design.md`](./docs/25-clara-network-system-design.md).

Start with [`docs/00-README.md`](./docs/00-README.md).

### Documentation library (all docs on `main`)

| # | Document | What it covers |
|---|----------|----------------|
| 00 | [Documentation Index](./docs/00-README.md) | How to use the library, reading paths, coverage matrix |
| 01 | [Card Payment Ecosystem & Four-Party Model](./docs/01-card-payment-ecosystem.md) | Participants, roles, open vs closed networks, fee flows |
| 02 | [Card Numbering & Identification](./docs/02-card-numbering-and-identification.md) | PAN structure, IIN/BIN, MII, Luhn check digit, BIN sponsorship |
| 03 | [Payment Flows: Authorization, Clearing, Settlement](./docs/03-payment-flows.md) | The three phases and message exchange lifecycle |
| 04 | [ISO 8583 - Card-originated Interchange Messages](./docs/04-iso8583.md) | The message standard for card transaction switching |
| 05 | [ISO 20022 - Payments Messaging](./docs/05-iso20022.md) | XML message standard for interbank clearing & settlement |
| 06 | [EMV Chip & Contactless](./docs/06-emv-chip.md) | Chip card spec, cryptograms, terminal transaction flow |
| 07 | [Tokenization & Network Tokens](./docs/07-tokenization.md) | EMV payment tokens, TSP, PAR, network token services |
| 08 | [3-D Secure (EMV 3DS)](./docs/08-3ds.md) | E-commerce cardholder authentication protocol |
| 09 | [PCI DSS Compliance](./docs/09-pci-dss.md) | Security standard for the cardholder data environment |
| 10 | [Fraud Detection & Risk Management](./docs/10-fraud-risk.md) | Real-time scoring, rules, velocity, machine learning |
| 11 | [Card Scheme Developer APIs](./docs/11-apis-integration.md) | Mastercard/Visa developer platforms and integration patterns |
| 12 | [System Design: Building the Platform](./docs/12-system-design.md) | Microservices, idempotency, ledger, reconciliation, scaling |
| 13 | [Fees, Interchange & Settlement Mechanics](./docs/13-fees-interchange-settlement.md) | How money and fees actually move |
| 14 | [Glossary & Source References](./docs/14-glossary-references.md) | Terminology and research sources |
| 15 | [Payment System Governance & Regulation](./docs/15-payment-system-governance-regulation.md) | Licensing, PFMI principles, SIPS oversight |
| 16 | [Membership, Rulebook & Certification](./docs/16-membership-rulebook-certification.md) | Member tiers, sponsorship, rulebook, host/L3 certification |
| 17 | [Key Management, HSM & Message Security](./docs/17-message-security-key-management.md) | PIN blocks (ISO 9564), HSMs, MACs, key hierarchy |
| 18 | [Settlement & Liquidity Infrastructure](./docs/18-settlement-liquidity-infrastructure.md) | RTGS accounts, prefunding, net/gross settlement, default fund |
| 19 | [Stand-In Processing & Availability](./docs/19-stand-in-processing-availability.md) | Stand-in rules, response codes 91/P, operational resilience |
| 20 | [Disputes, Chargebacks & Arbitration](./docs/20-disputes-chargeback-management.md) | Dispute lifecycle, reason codes, VCR, monitoring programs |
| 21 | [Cross-Border, FX & DCC](./docs/21-cross-border-fx-dcc.md) | FX conversion, dynamic currency conversion, cross-currency settlement |
| 22 | [Card Production & Lifecycle](./docs/22-card-production-lifecycle.md) | Personalization pipeline, EMV data, lifecycle, instant issuance |
| 23 | [Merchant Acquiring & Underwriting](./docs/23-merchant-acquiring-underwriting.md) | Boarding, MATCH/OFAC screening, MCC, reserves, monitoring |
| 24 | [Instant Payments & Real-Time Processing](./docs/24-instant-payments-rtp.md) | RTP/TIPS/Pix/UPI, prefunding, settlement models, ISO 20022 |
| 25 | [Clara Network System Design (Build Blueprint)](./docs/25-clara-network-system-design.md) | Decisive stack, module map, build phases 1–10, data flow |
| 26 | [Architecture Diagram Prompts](./docs/26-architecture-image-prompts.md) | Image-generation prompts; generated diagrams in [`docs/images/`](./docs/images/) |

The `docs/references/` directory holds source PDFs downloaded from public
institutions (BIS/CPMI, ECB, World Bank, FDIC, OCC, Visa, Mastercard, PCI SSC,
NIST) plus a manifest tracking every needed source, including member-gated and
paywalled documents.

## Architecture

<img src="docs/images/architecture-overview.png" width="620" alt="Clara Network end-to-end architecture">

*End-to-end system architecture: entry points, acquirer host, Clara switch,
issuer host, and core services.*

<img src="docs/images/auth-flow-sequence.png" width="620" alt="ISO 8583 authorization flow">

*ISO 8583 authorization flow, including the stand-in fallback path.*

<img src="docs/images/clearing-settlement.png" width="620" alt="Clearing, settlement & liquidity">

*Clearing, settlement, and liquidity: netting, prefunded accounts, default
fund, and central-bank RTGS.*

<img src="docs/images/issuer-tokenization.png" width="620" alt="Issuer, tokenization & card stack">

*Issuer, tokenization, and card stack: card production, HSM keys, BIN ranges,
and network tokens.*

<img src="docs/images/security-hsm-resilience.png" width="620" alt="Security, HSM & resilience">

*Security, HSM, and resilience: key hierarchy, MAC/PIN blocks, and active-active
site topology.*

The five diagrams were generated from the prompts in
[`docs/26-architecture-image-prompts.md`](docs/26-architecture-image-prompts.md).

## Building & running

Phases 1–10 implement the **authorization flow, net settlement, scheme ledger,
issuing stack, acquiring stack, disputes engine, key management,
operational resilience, and an instant-payment layer**: acquirer → switch →
(risk check) → issuer authorization with BIN-based routing, failover,
idempotent replay, in-path risk scoring, stand-in processing; a clearing
engine that captures clearing files, computes per-member net positions,
enforces prefunded caps, applies the default fund, and emits ISO 20022
pacs.009 settlement instructions; an append-only double-entry ledger that
posts every net position as a balanced journal and reconciles the ledger
against the settlement agent's statement; the issuing stack — BIN ranges,
card personalization, EMV-style ARQC cryptogram verification (with ATC
anti-replay), token vault (PAN → token + PAR), and mobile-wallet provisioning;
the acquiring stack — merchant boarding with MATCH/OFAC negative-list
screening, MCC assignment with risk tiering, and a funding engine that
withholds processing fees and rolling reserves and schedules merchant payouts;
the disputes engine — a reason-code taxonomy, the file → representment →
rule → arbitration lifecycle with fees charged to the losing party, the
associated-transaction (prior-credit) check, SLA deadline tracking, and
merchant chargeback-ratio monitoring; the key-management layer — a
Hardware Security Module simulation with dual-control key ceremonies
(M-of-N), AES key wrap (RFC 3394), TR-31-style key blocks for transport to
members, ISO 9564 PIN blocks (formats 0 and 4) verified inside the HSM,
ISO 9797-1 retail MACs with tamper detection, key rotation, a full audit
trail, and dual-control zeroize; the resilience layer — issuer stand-in
processing (SIP/STIP) with per-issuer limits and negative/valid-card files,
per-route circuit breakers with half-open probing (primary → secondary →
stand-in → decline), outcome metrics with approximate p99 latency, and
burst detection of issuer-inoperative (91) responses that flags an issuer
outage; and the instant-payment layer — ISO 20022 pacs.008 customer credit
transfers settled in real time, 24/7/365, against fully prefunded member
positions (the RTP model) with a 20-second scheme SLA, verify-and-reserve
settlement capacity checks, rejection reason codes (AC04/AC01/AG01/FF01),
SLA timeout handling with reservation release (NOAS), and pacs.002 status
reports. It requires Go 1.26+ and Docker Desktop.

```sh
# unit + integration tests
go test ./...

# run the full stack using Docker (postgres, redis, switch, issuer-sim,
# acquirer-sim, clearing-sim, ledger-sim, cardsvc, card-sim, acquiring-sim,
# disputes-sim, hsm-sim, resilience-sim, instant-sim): 6 auth requests with BIN routing and a velocity rule that
# declines the 6th with response code 59, then a settlement cycle with a
# member default covered by the default fund, a clean ledger + reconciliation
# run, the issuing stack demo (cryptogram verify, tokenize, provision), the
# acquiring stack demo (boarding decisions, fee/reserve funding, reserve
# release), the disputes demo (representment rulings, arbitration,
# associated-transaction rejection, chargeback ratios), the HSM demo
# (key ceremonies, PIN verify, retail MAC + tamper detection, key rotation,
# audit trail, zeroize), the resilience demo (failover to a secondary
# issuer, circuit-breaker trip, stand-in approvals/declines, 91-burst alert,
# half-open probe recovery), and the instant-payments demo (pacs.008 in ->
# ACSC/RJCT out, prefunded positions, SLA timeout with reservation release,
# position conservation)
docker compose -f deploy/docker-compose.yml up --build
docker compose -f deploy/docker-compose.yml logs switch acquirer-sim clearing-sim ledger-sim card-sim acquiring-sim disputes-sim hsm-sim resilience-sim instant-sim

# or run locally: terminal 1 -> switch, terminal 2 -> issuer-sim,
# terminal 3 -> acquirer-sim, terminal 4 -> clearing-sim, terminal 5 -> ledger-sim,
# terminal 6 -> cardsvc, terminal 7 -> card-sim, terminal 8 -> acquiring-sim,
# terminal 9 -> disputes-sim, terminal 10 -> hsm-sim, terminal 11 -> resilience-sim,
# terminal 12 -> instant-sim
go run ./cmd/switch
go run ./cmd/issuer-sim
go run ./cmd/acquirer-sim
go run ./cmd/clearing-sim
go run ./cmd/ledger-sim
go run ./cmd/cardsvc
go run ./cmd/card-sim
go run ./cmd/acquiring-sim
go run ./cmd/disputes-sim
go run ./cmd/hsm-sim
go run ./cmd/resilience-sim
go run ./cmd/instant-sim
```

Key config (via env):

- `CLARA_ISSUER_ROUTES` — JSON `{receiving-institution-id:
  host:port}`; a value may be a comma-separated failover list.
- `CLARA_BIN_TABLE` — JSON `{"entries":{"400000":"1000001000"}}` routes by
  PAN BIN when the message omits DE100.
- `CLARA_RISK_RULES` — JSON rule set; velocity counters (per card / per
  merchant) are counted in Redis and can decline with a configurable code.
- `CLARA_REDIS_ADDR` — idempotency + risk counters.
- `CLARA_PG_DSN` — audit log, clearing records, net positions, prefund
  accounts, default fund.
- `CLARA_SEND_DE100=false` (acquirer-sim) — omit DE100 to exercise BIN routing.
- `CLARA_SCENARIO` (clearing-sim) — `default` (prefund covers) or `default`
  run with a member default; settlement pacs.009 XML is written to
  `CLARA_OUT` (default `out/clearing`).
- `CLARA_MISMATCH` (ledger-sim) — when set, corrupts the settlement agent's
  statement to demonstrate reconciliation classification (amount mismatch and
  orphan-in-ledger).
- `CLARA_ISSUER_MASTER_KEY` (cardsvc/card-sim) — 16-byte AES master key for
  per-card key derivation and cryptogram verification.
- `CLARA_BIN`, `CLARA_PAN`, `CLARA_PRODUCT` (cardsvc/card-sim) — issued BIN
  range and the PAN to personalize; `CLARA_DEVICE_ID`, `CLARA_TRID` for
  wallet provisioning.
- `CLARA_LISTEN` (cardsvc) — HTTP listen address (default `:8081`); the REST
  API exposes `POST /cards`, `POST /cards/{ref}/arqc`,
  `POST /cards/{ref}/verify-arqc`, `POST /tokens`, `GET /tokens/{token}`,
  `POST /tokens/{token}/provision`.
- `CLARA_PG_DSN` (acquiring-sim) — persists merchants, funding lines, and the
  MATCH/OFAC screening lists; without it the demo uses an in-memory store.
- `CLARA_PG_DSN` (disputes-sim) — persists dispute cases and monitored
  transactions; without it the demo uses an in-memory store.
- `CLARA_PG_DSN` (hsm-sim) — not used; the HSM simulation is fully in-process
  (keys, ceremonies, and the audit trail live inside the HSM and are wiped on
  exit or via a dual-control `Zeroize`).
- `CLARA_PG_DSN` (resilience-sim) — not used; the chaos drill runs fully
  in-process: a switch fronts a primary and a secondary issuer on localhost,
  then the primary dies (circuit breaker trips, traffic fails over), the
  secondary dies too (stand-in approves within limits, declines hot cards and
  restricted BINs, and issues 91s that trip a burst alert), and finally the
  primary recovers (a half-open probe re-closes the circuit).
- `CLARA_PG_DSN` (instant-sim) — not used; the instant-payment demo runs fully
  in-process: pacs.008 credit transfers settle in real time against
  prefunded positions with a 20-second SLA, and rejections (AC04/AC01/AG01/
  FF01/NOAS) never move funds. The SLA shown for the timeout drill is
  configurable (`CLARA_INSTANT_SLA`, default `3s`) so the drill does not wait
  twenty seconds.

## Status

Research & specification library (docs 00–26), phase 1 (ISO 8583 switch),
phase 2 (authorization flow with BIN routing, risk, failover), phase 3
(clearing + net settlement with prefunding, default fund, pacs.009), phase 4
(append-only double-entry ledger + reconciliation against the settlement
statement), phase 5 (issuing stack: BIN ranges, card personalization, EMV ARQC
verification, token vault, wallet provisioning), phase 6 (acquiring stack:
merchant boarding with MATCH/OFAC screening, MCC risk tiering, fee/reserve
funding), phase 7 (disputes engine: reason codes, representment,
arbitration, associated-transaction check, chargeback monitoring), phase 8
(key management & security: HSM simulation, dual-control key ceremonies, AES
key wrap, PIN blocks, retail MACs, key rotation, audit, zeroize), phase 9
(operational resilience: stand-in processing with per-issuer limits and
negative/valid-card files, per-route circuit breakers with half-open probing,
outcome metrics and p99 latency, 91-burst outage detection, and a chaos
drill), and phase 10 (instant payments: ISO 20022 pacs.008 customer credit
transfers settled in real time, 24/7/365, against fully prefunded member
positions with a 20-second SLA, verify-and-reserve settlement capacity
checks, rejection reason codes, SLA timeout handling with reservation
release, and pacs.002 status reports) implemented. All ten blueprint phases
are complete — the network is feature-complete for a beta release.
Contributions are welcome.

## License

[MIT](./LICENSE) — see the LICENSE file for details.
