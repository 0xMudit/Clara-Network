# Clara Network Roadmap

This document outlines proposed future work beyond the 10 completed build
phases. Items are **drafts pending owner confirmation** — priorities and
timelines may change.

## Current Status

All 10 blueprint phases are implemented and released as **v0.1.0-beta**.
The system simulates a complete card payment network: ISO 8583 switch,
authorization with BIN routing, clearing & net settlement, double-entry
ledger, issuing (cards, tokens, EMV cryptograms), acquiring (merchants,
funding), disputes, HSM key management, resilience (circuit breakers,
stand-in), and instant payments (pacs.008 RTP).

---

## Phase 11 — Fraud Detection & ML Scoring

> Status: **Proposed** | Docs: [`docs/10-fraud-risk.md`](docs/10-fraud-risk.md)

The current risk engine (`internal/risk`) uses velocity-based rules. Phase 11
would add machine learning scoring.

**Planned work:**

- [ ] Pluggable scoring interface (`RiskScorer`) that wraps the existing
  velocity rules and adds an ML hook
- [ ] Feature extraction pipeline (transaction amount, time-of-day, merchant
  MCC, geo-velocity, device fingerprint)
- [ ] Integration with a model serving layer (ONNX Runtime or gRPC to an
  external model server)
- [ ] Real-time scoring with sub-10ms latency budget
- [ ] Model versioning and A/B scoring (shadow mode before live enforcement)
- [ ] Fraud case management UI in the web dashboard

---

## Phase 12 — 3-D Secure (EMV 3DS) Authentication

> Status: **Proposed** | Docs: [`docs/08-3ds.md`](docs/08-3ds.md)

E-commerce cardholder authentication protocol.

**Planned work:**

- [ ] 3DS Server implementation ( merchant integration endpoint)
- [ ] Directory Server simulation (protocol version negotiation)
- [ ] ACS simulation (cardholder authentication challenge/response)
- [ ] Frictionless flow vs. challenge flow decision engine
- [ ] Integration with the authorization flow (switch reads ARes data)

---

## Phase 13 — Cross-Border & FX Settlement

> Status: **Proposed** | Docs: [`docs/21-cross-border-fx-dcc.md`](docs/21-cross-border-fx-dcc.md)

Multi-currency and cross-border transaction support.

**Planned work:**

- [ ] FX rate provider interface (configurable rates or external API hook)
- [ ] Multi-currency clearing — net positions per currency pair
- [ ] Dynamic Currency Conversion (DCC) at the acquirer
- [ ] Cross-border interchange fee rules
- [ ] Settlement in multiple currencies with FX conversion

---

## Phase 14 — Expanded Card Production & Instant Issuance

> Status: **Proposed** | Docs: [`docs/22-card-production-lifecycle.md`](docs/22-card-production-lifecycle.md)

Full card lifecycle management beyond the current personalization.

**Planned work:**

- [ ] Card lifecycle states (active, suspended, lost, stolen, expired,
  cancelled)
- [ ] Instant issuance flow (branch/ATM card printing simulation)
- [ ] Card renewal and replacement logic
- [ ] EMV applet updates ( remote configuration)
- [ ] CVV/CVC2 generation and verification

---

## Phase 15 — Additional Regulatory Jurisdictions

> Status: **Proposed** | Docs: [`docs/15-payment-system-governance-regulation.md`](docs/15-payment-system-governance-regulation.md)

Configurable compliance rules for different countries/regions.

**Planned work:**

- [ ] Jurisdiction configuration (country code → regulatory rules)
- [ ] SIPS (Systemically Important Payment System) compliance hooks
- [ ] PFMI principle checks (operational risk, settlement finality)
- [ ] Transaction limits and sanctions screening per jurisdiction
- [ ] Audit trail export in jurisdiction-specific formats

---

## Phase 16 — Kubernetes & Production Deployment

> Status: **Proposed** | Docs: [`docs/12-system-design.md`](docs/12-system-design.md)

Move from Docker Compose to production-grade deployment.

**Planned work:**

- [ ] Helm chart for the full stack
- [ ] Kubernetes manifests (Deployments, Services, ConfigMaps, Secrets)
- [ ] Horizontal autoscaling for the switch
- [ ] PostgreSQL operator integration (CloudNativePG or Zalando)
- [ ] Redis Sentinel / Cluster configuration
- [ ] Observability: Prometheus metrics, OpenTelemetry tracing, structured
  logging
- [ ] Load testing framework (k6 or vegeta scripts)

---

## Phase 17 — Web Dashboard Expansion

> Status: **Proposed** | Current: `web/` (Next.js)

The admin dashboard currently shows raw API data. Phase 17 would make it
a full operational tool.

**Planned work:**

- [ ] Real-time transaction feed (WebSocket from switch)
- [ ] Settlement cycle visualization (net positions, fund flows)
- [ ] Disputes workflow UI (file, respond, escalate, evidence upload)
- [ ] Merchant onboarding portal
- [ ] Card management dashboard (issue, suspend, view tokens)
- [ ] HSM key ceremony UI (dual-control with audit trail)
- [ ] Role-based access control (admin, operator, auditor)

---

## Ongoing — Community & Documentation

- [ ] Translations: Spanish, French, Arabic, Hindi, Portuguese
- [ ] Video walkthrough of the authorization flow
- [ ] Interactive tutorial (browser-based demo)
- [ ] API reference docs (OpenAPI/Swagger for the Admin API)
- [ ] Compliance checklist guides per jurisdiction

---

## How to suggest features

Open a [Feature Request](https://github.com/0xMudit/Clara-Network/issues/new?template=feature_request.md)
on GitHub with the details of what you'd like to see.

## How to contribute to roadmap items

1. Comment on the relevant issue (or create one if none exists).
2. Read [`CONTRIBUTING.md`](CONTRIBUTING.md) for the development workflow.
3. Read [`docs/28-contributor-architecture-guide.md`](docs/28-contributor-architecture-guide.md)
   for where the code lives.
4. Submit a PR with tests and documentation.
