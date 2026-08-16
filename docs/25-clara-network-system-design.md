# 25. Clara Network — Build Blueprint (System Design)

## 25.1 The goal

Clara Network is a Mastercard/Visa-style **card payment scheme** plus the
issuing and acquiring infrastructure to run it. Realistically this is built
as a **monorepo of services** centered on an **ISO 8583 message switch**,
extended with clearing/settlement, issuing, acquiring, risk, key management,
disputes, and an instant-payment layer. The research library (docs 00–24)
contains the domain knowledge each component encodes.

## 25.2 Decisive stack

| Concern | Choice | Why |
|---------|--------|-----|
| Language | **Go** | Low-latency concurrency for the switch, single static binaries, strong stdlib networking |
| Message switch | **ISO 8583 engine** (custom, bitmapped) | The network core; see `04-iso8583.md` |
| Interbank messaging | **ISO 20022** (pacs/pain/camt) | Clearing, settlement, instant layer; see `05-iso20022.md` |
| Service communication | gRPC internal + **Kafka events** | Idempotent, replayable message bus; see `12-system-design.md` |
| Data | **PostgreSQL** (ledger, accounts, merchant/card data), **Redis** (hot cache, velocity, session) | Correctness (ledger) + hot-path speed |
| Crypto | **HSM** (Thales/Safenet) for all key operations; never keys in software | See `17-message-security-key-management.md` |
| Ops | Kubernetes, 2+ AZs, active-active | See `19-stand-in-processing-availability.md` |

## 25.3 Module map

```
clara-network/
├── switch/          # ISO 8583 message switch (the network core)
│   ├── parser/      # bitmap ISO 8583 parse/build
│   ├── router/      # route acquirer→issuer by BIN, failover, stand-in
│   ├── idempotency/ # replay-safe request/response handling
│   └── sim/         # acquirer & issuer simulators for testing
├── issuer/          # issuer host: auth decisions, accounts, card data, stand-in params
├── acquirer/        # acquirer host: merchants, terminals, clearing capture, funding
├── risk/            # rules engine + ML scoring, velocity (sub-100ms)
├── clearing/        # clearing capture, netting, settlement instructions, prefunding caps
├── ledger/          # double-entry ledger, reconciliation
├── tokenvault/      # network tokens, PAR (see 07-tokenization.md)
├── keymgmt/         # HSM integration, PIN blocks, MACs, key hierarchy (see 17)
├── disputes/        # chargeback engine, reason codes, deadlines (see 20)
├── instant/         # ISO 20022 pacs.008, 24/7/365, 20s SLA (see 24)
├── cardsvc/         # card lifecycle, BIN ranges, CVV/cryptogram data (see 22)
├── adminapi/        # member boarding, merchant boarding, dashboards (see 16, 23)
├── internal/        # shared libs: message models, errors, tracing, config
└── deploy/          # k8s manifests, terraform, helm
```

## 25.4 Build phases (ordered by dependency)

1. **Switch skeleton + ISO 8583 model** — message spec, parser, simulator
   harness. Success check: acquirer-sim → switch → issuer-sim → response
   round-trip with correct MTI/bitmap.
2. **Authorization flow** — route by BIN, risk scoring in path, stand-in when
   issuer is down (response codes `91`/`P`). Success check: end-to-end auth
   under 100 ms p99 with replay-safe idempotency.
3. **Clearing + net settlement** — clearing capture, netting per member,
   prefunded settlement accounts, settlement instruction to the settlement
   agent, default fund sizing. Success check: batch settles with correct
   net positions and finality (see `18-settlement-liquidity-infrastructure.md`).
4. **Ledger + reconciliation** — double-entry posting on authorization and
   clearing; reconcile to central-bank/settlement-agent statements.
5. **Issuing stack** — BIN ranges, card data, cryptogram verification via
   HSM, token vault, mobile-wallet provisioning.
6. **Acquiring stack** — merchant boarding, MCC assignment, MATCH/OFAC
   screening, reserves, funding (see `23-merchant-acquiring-underwriting.md`).
7. **Disputes engine** — reason codes, evidence capture, deadlines, fees.
8. **Key management + security** — HSM integration, PIN blocks (ISO 9564),
   MACs, dual-control ceremonies (see `17`).
9. **Resilience** — active-active multi-AZ, RTO/RPO targets, chaos testing,
   monitoring (see `19`).
10. **Instant payments layer** — ISO 20022, 24/7/365, prefunding, 20 s
    timeout (see `24`).

## 25.5 Data flow (authorization, happy path)

```
POS/terminal ──ISO 8583──> acquirer ──> switch ──risk check──> issuer host
                                                            (approve/decline)
POS <────────────── ISO 8583 response (auth code) <────────── switch
acquirer ──clearing file──> clearing engine ──net positions──> settlement agent
members ──prefunded funds──> settlement account
```

## 25.6 Critical design details

- **Idempotency**: every request carries a unique trace/STAN; re-deliveries
  return the stored response (see `12-system-design.md`).
- **Finality**: net positions become final per scheme rules; prefunding caps
  guarantee funds even on member default (see `18`).
- **Latency budget**: risk and switch logic must fit a p99 < 100 ms budget.
- **Security**: all PIN/MAC/key operations inside the HSM; no keys in
  application memory; PCI DSS scope documented (see `09`, `17`).
- **Testing**: simulator-driven conformance testing before any member/cert
  certification (see `16-membership-rulebook-certification.md`).
