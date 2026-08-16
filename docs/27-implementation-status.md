# 27. Clara Network — Implementation Status

This document records what has actually been **built, tested, and released**
(as of `v0.1.0-beta`) against the blueprint in `25-clara-network-system-design.md`.
All ten build phases are implemented. Each section maps a blueprint phase to
the concrete packages, commands, tests, and simulators that realize it, plus
where the blueprint and the implementation diverge.

## 27.1 Repository layout (as built)

```
clara-network/
├── cmd/
│   ├── switch/           # ISO 8583 message switch (long-running server)
│   ├── issuer-sim/       # issuer host simulator (long-running server)
│   ├── acquirer-sim/     # acquirer host simulator (client driving auths)
│   ├── cardsvc/          # issuing stack REST API (long-running server)
│   ├── card-sim/         # issuing stack demo client
│   ├── clearing-sim/     # clearing + net settlement demo
│   ├── ledger-sim/       # double-entry ledger + reconciliation demo
│   ├── acquiring-sim/    # merchant boarding / funding demo
│   ├── disputes-sim/     # disputes + arbitration demo
│   ├── hsm-sim/          # HSM simulation demo
│   ├── resilience-sim/   # outage chaos drill demo
│   └── instant-sim/      # instant-payments (RTP) demo
├── internal/
│   ├── iso8583/          # ISO 8583 message model, bitmap parse/build
│   ├── framing/          # 2-byte length-prefixed TCP framing
│   ├── switchsrv/        # switch: routing, failover, idempotency, risk hook, stand-in, metrics
│   ├── binrouting/       # BIN -> issuing-institution routing table
│   ├── risk/             # rule engine: velocity counters (memory/Redis)
│   ├── acquirersim/      # acquirer host logic behind acquirer-sim
│   ├── issuersim/        # issuer host logic behind issuer-sim
│   ├── clearing/         # capture, netting, prefund caps, default fund, pacs.009
│   ├── ledger/           # append-only double-entry ledger, reconciliation
│   ├── cardsvc/          # BIN ranges, personalization, ARQC, token vault, provisioning
│   ├── acquiring/        # boarding, MATCH/OFAC, MCC tiering, fee/reserve funding
│   ├── disputes/         # reason codes, representment, arbitration, monitoring
│   ├── hsm/              # in-process HSM simulation (see 27.9)
│   ├── resilience/       # stand-in, circuit breakers, metrics, 91-burst detection
│   ├── instant/          # ISO 20022 pacs.008/002, prefunded RTP engine
│   └── env/              # CLARA_* config helpers
├── deploy/docker-compose.yml   # full 14-service stack
├── Dockerfile                   # multi-stage build for every cmd
└── Makefile                     # build/test/vet/run-<sim> targets
```

Module `github.com/0xMudit/Clara-Network`, Go `1.25`. Direct dependencies
are only `github.com/jackc/pgx/v5` (optional PostgreSQL) and
`github.com/redis/go-redis/v9` (optional Redis).

## 27.2 Phase-by-phase status

> Blueprint phase numbers refer to `docs/25` §25.4.

### Phase 1 — Switch skeleton + ISO 8583 model ✅

- `internal/iso8583` — message classes, MTI/BMP, data elements, build/parse
  with correct bitmap encoding; `internal/framing` — 2-byte big-endian
  length-prefixed TCP framing.
- `cmd/switch` — the long-running network core; `cmd/acquirer-sim` and
  `cmd/issuer-sim` drive the round trip.
- Tests: `internal/iso8583` (7), `internal/switchsrv/server_test.go` (7).
- **Success check met**: acquirer-sim → switch → issuer-sim → response
  round trip with correct MTI/bitmap.

### Phase 2 — Authorization flow ✅

- `internal/switchsrv` — routing, BIN-based routing via `internal/binrouting`,
  risk check in path (`internal/risk`), failover across a comma-separated
  route list, replay-safe idempotency (in-memory, or Redis), response-code
  mapping (`00` approve, `51` insufficient funds, `59` risk decline, `91`
  issuer unavailable).
- Tests: `internal/switchsrv/phase2_test.go` (5), `internal/binrouting` (4),
  `internal/risk` (5).
- **Success check met**: end-to-end auth under the p99 budget with
  idempotent replay returning the stored response.

### Phase 3 — Clearing + net settlement ✅

- `internal/clearing` — clearing capture, per-member netting, prefunded
  settlement accounts (balance + cap), default fund applied on member
  default, ISO 20022 **pacs.009** settlement instructions written as XML.
- `cmd/clearing-sim` — one-shot batch demo; `CLARA_SCENARIO=default` covers
  prefunding, and a member-default scenario drains the default fund.
- Tests: `internal/clearing` (11).
- **Success check met**: batch settles with correct net positions and
  finality (`docs/18`).

### Phase 4 — Ledger + reconciliation ✅

- `internal/ledger` — append-only double-entry ledger; every net position
  posts as a balanced journal; reconciliation against the settlement agent's
  statement with mismatch / orphan-in-ledger classification.
- `cmd/ledger-sim` — consumes the same cycle as clearing-sim.
- Tests: `internal/ledger` (14).

### Phase 5 — Issuing stack ✅

- `internal/cardsvc` — BIN range registry, card personalization, EMV-style
  **ARQC cryptogram verification** (with ATC anti-replay), token vault
  (PAN → token + PAR), mobile-wallet provisioning.
- `cmd/cardsvc` — REST API (`:8081`): `POST /cards`, `POST /cards/{ref}/arqc`,
  `POST /cards/{ref}/verify-arqc`, `POST /tokens`, `GET /tokens/{token}`,
  `POST /tokens/{token}/provision`; `cmd/card-sim` — demo client.
- Tests: `internal/cardsvc` (14).
- **Blueprint divergence**: cryptograms are verified in-process with a
  software master key (`CLARA_ISSUER_MASTER_KEY`) rather than inside the HSM
  module (see 27.9); the HSM layer handles PIN/MAC/key material.

### Phase 6 — Acquiring stack ✅

- `internal/acquiring` — merchant boarding with **MATCH/OFAC** negative-list
  screening, MCC assignment with risk tiering, funding engine that withholds
  processing fees and rolling reserves and schedules merchant payouts.
- `cmd/acquiring-sim` — one-shot demo.
- Tests: `internal/acquiring` (10).

### Phase 7 — Disputes engine ✅

- `internal/disputes` — reason-code taxonomy, the file → representment →
  rule → arbitration lifecycle with fees charged to the losing party, the
  associated-transaction (prior-credit) check, SLA deadline tracking, and
  merchant chargeback-ratio monitoring.
- `cmd/disputes-sim` — one-shot demo.
- Tests: `internal/disputes` (9).

### Phase 8 — Key management + security ✅

- `internal/hsm` — in-process **HSM simulation** (see 27.9 for the scope
  note): dual-control key ceremonies (M-of-N), AES key wrap (RFC 3394),
  TR-31-style key blocks for transport to members, ISO 9564 PIN blocks
  (formats 0 and 4), ISO 9797-1 retail MACs with tamper detection, key
  rotation, a full audit trail, and dual-control zeroize.
- `cmd/hsm-sim` — one-shot demo.
- Tests: `internal/hsm` (11).

### Phase 9 — Resilience ✅

- `internal/resilience` — issuer **stand-in processing** (SIP/STIP) with
  per-issuer limits and negative/valid-card files; response-code semantics
  `00` (stand-in approve), `05` (hot card), `57` (restricted BIN), `91`
  (stand-in declined / issuer inoperative); per-route **circuit breakers**
  with half-open probing and a primary → secondary → stand-in → decline
  ordering; outcome metrics with approximate p99 latency; and burst detection
  of `91` responses that flags an issuer outage.
- `cmd/resilience-sim` — one-shot chaos drill: kill primary → failover, kill
  secondary → stand-in + 91 burst alert, recover primary → half-open probe
  re-closes the circuit.
- Tests: `internal/resilience` (13), `internal/switchsrv/phase9_test.go` (5).
- **Blueprint divergence**: the blueprint's active-active multi-AZ is
  documented but not provisioned (single-node compose); RTO/RPO and chaos
  practices are demonstrated by the drill rather than operated in production.

### Phase 10 — Instant payments layer ✅

- `internal/instant` — ISO 20022 **pacs.008** customer credit transfers
  settled **in real time, 24/7/365** against fully prefunded member positions
  with a **20-second scheme SLA** (configurable); verify-and-reserve
  settlement-capacity checks; rejection reason codes `AC04` (insufficient
  funds), `AC01` (unknown beneficiary), `AG01` (forbidden), `FF01`
  (format/currency/self-transfer); SLA timeout → `NOAS` with the reservation
  released; **pacs.002** status reports out.
- `cmd/instant-sim` — one-shot demo proving ACSC settlement, every rejection,
  reservation release on timeout, and position conservation.
- Tests: `internal/instant` (13), including a concurrent 10-goroutine race
  that reserves atomically (exactly 5 of 10 settle).

## 27.3 Realized architecture vs. blueprint

| Blueprint (`docs/25` §25.2) | As built | Notes |
|---|---|---|
| gRPC internal + Kafka events | Length-prefixed TCP + ISO 8583 | Single binary per service; no message bus |
| Kubernetes, 2+ AZs | Docker Compose, single node | Deployment topology is a deliberate simplification |
| HSM (Thales/Safenet) | In-process `internal/hsm` simulation | Algorithms match ISO/AES specs; see 27.9 |
| PostgreSQL everywhere | Optional PostgreSQL (audit, clearing, ledger, acquiring, disputes) | In-memory stores are the default in the sims |
| Redis hot cache | Redis for idempotency + velocity (optional) | In-memory fallbacks built in |

Everything else in the blueprint (ISO 8583 engine, ISO 20022 messages,
prefunding/default-fund settlement, double-entry ledger, stand-in, token
vault) is implemented as specified.

## 27.4 Services and ports

| Service | Command | Listen | Notes |
|---|---|---|---|
| Switch | `go run ./cmd/switch` | `:8080` (`CLARA_LISTEN`) | ISO 8583 over TCP |
| Issuer | `go run ./cmd/issuer-sim` | `:8082` (`CLARA_LISTEN`) | Issuer host |
| Cardsvc | `go run ./cmd/cardsvc` | `:8081` (`CLARA_LISTEN`) | REST API |
| Acquirer | `go run ./cmd/acquirer-sim` | — (client) | Connects to `CLARA_SWITCH` |
| Resilience drill | `cmd/resilience-sim` | `19080/19081/19082` | Switch + primary + secondary on localhost |

`docker compose -f deploy/docker-compose.yml up --build` runs the whole stack
(postgres, redis, switch, issuer-sim, acquirer-sim, clearing-sim, ledger-sim,
cardsvc, card-sim, acquiring-sim, disputes-sim, hsm-sim, resilience-sim,
instant-sim). The one-shot sims (acquirer-sim, card-sim, clearing-sim,
ledger-sim, acquiring-sim, disputes-sim, hsm-sim, resilience-sim,
instant-sim) run to completion and exit 0; switch, issuer-sim, cardsvc,
postgres, and redis stay up.

## 27.5 Configuration (`CLARA_*`)

| Variable | Used by | Default | Purpose |
|---|---|---|---|
| `CLARA_LISTEN` | switch, issuer-sim, cardsvc | `:8080` / `:8082` / `:8081` | Listen address |
| `CLARA_SWITCH` | acquirer-sim | `localhost:8080` | Switch address |
| `CLARA_ISSUER_ROUTES` | switch | — | JSON `{rid: host:port}`; comma-separated failover lists |
| `CLARA_BIN_TABLE` | switch | — | JSON BIN table when DE100 is omitted |
| `CLARA_RISK_RULES` | switch | — | JSON velocity rules; Redis-backed if available |
| `CLARA_REDIS_ADDR` | switch | — | Idempotency + risk counters (memory fallback) |
| `CLARA_PG_DSN` | switch / sims | — | Optional persistence (audit, clearing, ledger, acquiring, disputes) |
| `CLARA_SEND_DE100` | acquirer-sim | `true` | Set `false` to exercise BIN routing |
| `CLARA_SCENARIO`, `CLARA_CYCLE`, `CLARA_OUT`, `CLARA_MISMATCH` | clearing/ledger sims | `default` / today / `out/clearing` / off | Batch scenario, cycle id, output dir, reconcile mismatch drill |
| `CLARA_ISSUER_MASTER_KEY` | cardsvc, card-sim | fixed demo key | AES master key for card crypto |
| `CLARA_BIN`, `CLARA_PAN`, `CLARA_PRODUCT`, `CLARA_DEVICE_ID`, `CLARA_TRID` | cardsvc, card-sim | — | Issue parameters and wallet provisioning |
| `CLARA_INSTANT_SLA` | instant-sim | `3s` | SLA for the timeout drill (scheme default is 20s) |

## 27.6 Test coverage

`go test ./...` (fresh) is green. 128 tests across the modules:

- switchsrv (17: server, phase 2, phase 9), iso8583 (7), binrouting (4),
  risk (5), clearing (11), ledger (14), cardsvc (14), acquiring (10),
  disputes (9), hsm (11), resilience (13), instant (13).

`go vet ./...` is clean. Every sim is verified to exit 0 both locally and as a
Docker image.

## 27.7 What is not implemented (deliberate)

- Real network switch/hop, kafka/gRPC internal bus, Kubernetes, TLS on the
  wire, and a production HSM are out of scope for this beta.
- The `adminapi` module (member boarding, dashboards) and `deploy/` k8s
  manifests from the blueprint module map do not exist yet.
- Risk scoring is rule/velocity-based; ML scoring is documented (`docs/10`)
  but not implemented.

## 27.8 Running the stack

- One-shot sims (no dependencies): `go run ./cmd/<x>-sim` for clearing,
  ledger, card, acquiring, disputes, hsm, resilience, instant.
- Servers: `go run ./cmd/switch` + `go run ./cmd/issuer-sim` +
  `go run ./cmd/acquirer-sim` (+ optional `go run ./cmd/cardsvc` + card-sim).
- Everything at once: `docker compose -f deploy/docker-compose.yml up --build`
  (14 services; the one-shot sims exit 0, servers stay up), then
  `docker compose ... down`.
- `make` targets: `build`, `test`, `vet`, `run-switch`, `run-issuer`,
  `run-acquirer`, `run-clearing`, `run-ledger`, `run-cardsvc`, `run-card-sim`,
  `run-acquiring`, `run-disputes`, `run-hsm`, `run-resilience`,
  `run-instant`.

## 27.9 HSM simulation scope

`internal/hsm` is an in-process simulation, not a physical Thales/Safenet
appliance. Its cryptographic behavior follows the published specifications —
AES key wrap (RFC 3394), TR-31-style key blocks, ISO 9564 PIN blocks
(formats 0 and 4), ISO 9797-1 retail MACs — and the dual-control ceremony,
audit trail, rotation, and zeroize workflows mirror HSM administration. It
uses Go's stdlib `crypto` packages and keeps keys in process memory. For a
production deployment the same `hsm.Store`/ceremony interface would be backed
by a real appliance; the message-level crypto (ARQC, PIN, MAC) is
deliberately exercised against software keys in this beta.
