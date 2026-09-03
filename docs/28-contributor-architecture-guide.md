# 28. Clara Network — Contributor Architecture Guide

This document is a **fast contributor on-ramp** — distinct from the research
library (docs 00–24) and the implementation status (doc 27). If you're a new
contributor looking to understand the codebase and find where to make changes,
start here.

## 28.1 How the codebase is organized

Every feature in Clara Network follows the same pattern:

```
internal/<package>/     # Core logic (library code)
cmd/<name>/             # Entry point (main.go that wires the package)
internal/<package>/_test.go  # Tests
```

The `internal/` packages contain the actual implementation. The `cmd/`
directories are thin wrappers that parse config, wire dependencies, and start
the service. Most contributions happen in `internal/`.

## 28.2 Phase-to-code map

Use this table to find the code for any feature:

| Phase | What it does | Core package(s) | Entry point(s) | Tests |
|-------|-------------|-----------------|-----------------|-------|
| **1** | Switch + ISO 8583 | `internal/iso8583`, `internal/framing` | `cmd/switch` | `internal/iso8583/` (7) |
| **2** | Authorization flow | `internal/switchsrv`, `internal/binrouting`, `internal/risk` | `cmd/acquirer-sim`, `cmd/issuer-sim` | `internal/switchsrv/` (12), `internal/binrouting/` (4), `internal/risk/` (5) |
| **3** | Clearing + settlement | `internal/clearing` | `cmd/clearing-sim` | `internal/clearing/` (11) |
| **4** | Ledger + reconciliation | `internal/ledger` | `cmd/ledger-sim` | `internal/ledger/` (14) |
| **5** | Issuing (cards, tokens) | `internal/cardsvc` | `cmd/cardsvc`, `cmd/card-sim` | `internal/cardsvc/` (14) |
| **6** | Acquiring (merchants) | `internal/acquiring` | `cmd/acquiring-sim` | `internal/acquiring/` (10) |
| **7** | Disputes | `internal/disputes` | `cmd/disputes-sim` | `internal/disputes/` (9) |
| **8** | Key management / HSM | `internal/hsm` | `cmd/hsm-sim` | `internal/hsm/` (11) |
| **9** | Resilience | `internal/resilience` | `cmd/resilience-sim` | `internal/resilience/` (13), `internal/switchsrv/phase9_test.go` (5) |
| **10** | Instant payments | `internal/instant` | `cmd/instant-sim` | `internal/instant/` (13) |

Cross-cutting packages (not tied to a single phase):

| Package | Purpose |
|---------|---------|
| `internal/iso8583` | ISO 8583 message model — used by switch, issuer, acquirer |
| `internal/framing` | TCP framing — used by switch and issuer |
| `internal/switchsrv` | Switch server — phases 1, 2, 9 |
| `internal/binrouting` | BIN routing table — used by switch in phase 2 |
| `internal/risk` | Risk/velocity engine — used by switch in phase 2 |
| `internal/env` | `CLARA_*` config helpers — used everywhere |
| `internal/metrics` | Latency histograms — used by switch and resilience |
| `internal/adminapi` | Admin REST API — reads from all tables |

## 28.3 "Where do I add X?" guide

### "I want to add a new ISO 8583 data element (field)"

1. **Define the field** in `internal/iso8583/` — add the field number
   constant, its length/type, and bitmap position.
2. **Update the parser/builder** to handle the new field in
   `internal/iso8583/parse.go` and `internal/iso8583/build.go`.
3. **Add a test** in `internal/iso8583/iso8583_test.go` — round-trip
   parse → build → parse the new field.
4. **Wire it through the switch** if the switch needs to read or forward it
   (`internal/switchsrv/`).

### "I want to add a new ISO 8583 message type (MTI)"

1. **Add the MTI constant** in `internal/iso8583/` (e.g., `MTIAuthReq = "0100"`).
2. **Add the message class** and its required/optional fields.
3. **Add a handler** in `internal/switchsrv/` for routing the new MTI.
4. **Add a simulator** in `cmd/` if it needs a demo client.

### "I want to add a new dispute reason code"

1. **Add the code constant** in `internal/disputes/` (e.g.,
   `ReasonCodeFraud = "4837"`).
2. **Add the code to the taxonomy** — update the reason code lookup table.
3. **Add test coverage** in `internal/disputes/disputes_test.go` for any
   lifecycle behavior specific to this code.
4. **Update docs/20-disputes-chargeback-management.md** with the new code.

### "I want to add a new risk rule type"

1. **Define the rule kind** in `internal/risk/` — add a new `Kind` constant.
2. **Implement the counter/scorer** — the existing velocity rules are a
   good template.
3. **Add tests** in `internal/risk/risk_test.go`.
4. **Update the `CLARA_RISK_RULES` JSON format** documentation in the
   README and docs/27.

### "I want to add a new Admin API endpoint"

1. **Add the store query** in `internal/adminapi/store.go`.
2. **Add the handler** in `internal/adminapi/handlers.go`.
3. **Register the route** in the server setup (look for `mux.HandleFunc`).
4. **Add a test** in `internal/adminapi/` (use the existing test patterns).
5. **Update the README** endpoint table and `docs/27-implementation-status.md`.

### "I want to add a new settlement or clearing feature"

1. **Add the logic** in `internal/clearing/`.
2. **Add tests** in `internal/clearing/clearing_test.go`.
3. **Add a new scenario** to `cmd/clearing-sim/` if it's demo-able.
4. **Update `deploy/schema.sql`** if new tables or columns are needed.
5. **Update docs/27** with the new feature.

### "I want to add a new card or token feature"

1. **Add the logic** in `internal/cardsvc/`.
2. **Add tests** in `internal/cardsvc/cardsvc_test.go`.
3. **Update the REST API** in `cmd/cardsvc/` if new endpoints are needed.
4. **Update `cmd/card-sim/`** to demo the new feature.

### "I want to add a new merchant acquiring feature"

1. **Add the logic** in `internal/acquiring/`.
2. **Add tests** in `internal/acquiring/acquiring_test.go`.
3. **Update `cmd/acquiring-sim/`** to demo the new feature.

### "I want to add a new HSM or key management feature"

1. **Add the logic** in `internal/hsm/`.
2. **Add tests** in `internal/hsm/hsm_test.go`.
3. **Update `cmd/hsm-sim/`** to demo the new feature.

> **Important:** The HSM module is an in-process simulation. For production
> use, the `hsm.Store` interface would be backed by a real appliance.

### "I want to add a new resilience feature"

1. **Add the logic** in `internal/resilience/`.
2. **Add tests** in `internal/resilience/resilience_test.go`.
3. **Update `cmd/resilience-sim/`** to demo the new feature.

### "I want to add a new instant payment feature"

1. **Add the logic** in `internal/instant/`.
2. **Add tests** in `internal/instant/instant_test.go`.
3. **Update `cmd/instant-sim/`** to demo the new feature.

### "I want to add a new config variable (`CLARA_*`)"

1. **Add the helper** in `internal/env/` (e.g., `env.Duration("CLARA_FOO", 5*time.Second)`).
2. **Use the helper** in the service that needs it — never call
   `os.Getenv` directly.
3. **Document it** in the README config section and `docs/27-implementation-status.md`.
4. **Add it to `deploy/docker-compose.yml`** if it needs a default in Docker.

### "I want to add a new database table"

1. **Add the `CREATE TABLE`** in `deploy/schema.sql`.
2. **Add a store method** in the relevant `internal/*/store.go` or
   `internal/adminapi/store.go`.
3. **Add tests** for the new store method.
4. **Update docs/27** with the new table and its purpose.

### "I want to improve the web dashboard"

1. The dashboard lives in `web/` and is a Next.js app.
2. It talks to the Admin API at `:8083`.
3. See `web/README.md` for local development instructions.

## 28.4 Data flow: how a card authorization works end-to-end

Understanding this flow helps you see where each package fits:

```
POS terminal
    │
    ▼
acquirer-sim ──── creates ISO 8583 auth request (MTI 0100)
    │
    ▼
cmd/switch ──── receives TCP frame
    │
    ├── internal/framing ── parses 2-byte length prefix
    ├── internal/iso8583 ── parses bitmap and data elements
    ├── internal/binrouting ── looks up BIN → issuer ID
    ├── internal/risk ── checks velocity rules (Redis or memory)
    ├── internal/switchsrv ── routes to issuer, handles failover + stand-in
    │
    ▼
cmd/issuer-sim ──── receives auth request, approves/declines
    │
    ▼
cmd/switch ──── returns ISO 8583 response (MTI 0110)
    │
    ▼
acquirer-sim ──── receives response, logs result
    │
    ▼ (later, in batch)
cmd/clearing-sim ──── captures transactions, nets positions, emits pacs.009
    │
    ▼
cmd/ledger-sim ──── posts balanced journals, reconciles
```

## 28.5 Key design decisions

These are the architectural choices that shape how you contribute:

1. **Single binary per service** — each `cmd/<name>/` builds to one binary.
   Don't create multi-binary packages.

2. **Optional persistence** — PostgreSQL and Redis are optional. Every service
   falls back to in-memory stores. New features should work without a database.

3. **Environment-driven config** — all config is `CLARA_*` env vars, parsed
   through `internal/env/`. No config files, no flags.

4. **Simulator-driven testing** — the one-shot sims (`cmd/*-sim`) serve as
   both demos and integration tests. When adding a feature, add a sim that
   exercises it.

5. **In-process HSM** — the HSM module is a simulation. The interface
   (`hsm.Store`) is designed so a real HSM can be swapped in. Don't couple
   business logic to the software simulation.

6. **Minimal dependencies** — only `pgx/v5` and `go-redis/v9` as external
   deps. Prefer the stdlib for everything else.

7. **ISO standards compliance** — when implementing a feature from the docs,
   follow the referenced ISO standard. Cite the standard number in comments.

## 28.6 Common pitfalls

- **Tests must run on Linux.** Windows hosts with AppLocker block unsigned
  test binaries. Use WSL2 or Docker.
- **Don't change `deploy/schema.sql` without a migration story.** The schema
  is loaded on first Postgres startup. Changes need to be backward-compatible
  or accompanied by a migration script.
- **Don't break the ISO 8583 message format.** Other tools may depend on
  wire-level compatibility. Add new fields; don't reorder existing ones.
- **Don't commit secrets.** The demo keys and credentials in the repo are
  explicitly for demonstration only. Production deployments must replace them.
- **Don't add heavy dependencies.** If you think you need a framework, check
  if the stdlib can do it. The project's minimal dependency footprint is
  intentional.

## 28.7 Related documents

| Document | What it covers |
|----------|----------------|
| [`docs/00-README.md`](00-README.md) | Documentation index and reading paths |
| [`docs/25-clara-network-system-design.md`](25-clara-network-system-design.md) | The build blueprint |
| [`docs/27-implementation-status.md`](27-implementation-status.md) | What's built, tested, and released |
| [`CONTRIBUTING.md`](../CONTRIBUTING.md) | Dev setup, commit style, PR process |
| [`ROADMAP.md`](../ROADMAP.md) | Planned future work |
