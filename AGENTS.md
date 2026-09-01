# AGENTS.md

Instructions for AI coding agents (Cursor, Copilot, Codebuff, etc.) working on
the Clara Network codebase.

## Project overview

Clara Network is an open-source, Mastercard/Visa-style card payment network
built in Go. It implements a full end-to-end payment stack: ISO 8583 message
switch, authorization with BIN routing, clearing & net settlement, a
double-entry ledger, issuing (cards, tokens, EMV cryptograms), acquiring
(merchant boarding, funding), disputes, key management (HSM simulation),
resilience (circuit breakers, stand-in), and instant payments (pacs.008).

The project is designed so that **county and bank engineers can fork it and
build their own payment system with a few tweaks**.

## Repository layout

```
clara-network/
├── cmd/                    # Entry points (one per microservice)
│   ├── switch/             # ISO 8583 message switch
│   ├── issuer-sim/         # Issuer host simulator
│   ├── acquirer-sim/       # Acquirer host simulator (drives auths)
│   ├── cardsvc/            # Card service REST API
│   ├── card-sim/           # Card service demo client
│   ├── clearing-sim/       # Clearing + net settlement demo
│   ├── ledger-sim/         # Ledger + reconciliation demo
│   ├── acquiring-sim/      # Merchant boarding / funding demo
│   ├── disputes-sim/       # Disputes + arbitration demo
│   ├── hsm-sim/            # HSM simulation demo
│   ├── resilience-sim/     # Outage chaos drill demo
│   ├── instant-sim/        # Instant payments demo
│   └── adminapi/           # Read-only Admin REST API
├── internal/               # Core library packages
│   ├── iso8583/            # ISO 8583 message model, bitmap parse/build
│   ├── framing/            # 2-byte length-prefixed TCP framing
│   ├── switchsrv/          # Switch: routing, failover, idempotency, risk, stand-in
│   ├── binrouting/         # BIN → issuing-institution routing table
│   ├── risk/               # Rule engine: velocity counters (memory/Redis)
│   ├── acquirersim/        # Acquirer host logic
│   ├── issuersim/          # Issuer host logic
│   ├── clearing/           # Clearing capture, netting, prefund caps, default fund, pacs.009
│   ├── ledger/             # Append-only double-entry ledger, reconciliation
│   ├── cardsvc/            # BIN ranges, personalization, ARQC, token vault, provisioning
│   ├── acquiring/          # Merchant boarding, MATCH/OFAC screening, MCC tiering, funding
│   ├── disputes/           # Reason codes, representment, arbitration, monitoring
│   ├── hsm/                # In-process HSM simulation
│   ├── resilience/         # Stand-in, circuit breakers, metrics, 91-burst detection
│   ├── instant/            # ISO 20022 pacs.008/002, prefunded RTP engine
│   ├── adminapi/           # Admin REST API: server, handlers, store
│   └── env/                # CLARA_* config helpers
├── deploy/                 # Docker Compose stack (14 services)
├── docs/                   # 27-document research & specification library
└── web/                    # Next.js admin dashboard frontend
```

## Tech stack

- **Language:** Go 1.26+
- **Database:** PostgreSQL 16 (optional — all services fall back to in-memory)
- **Cache:** Redis 7 (optional — used for idempotency + velocity counters)
- **Container runtime:** Docker with Docker Compose
- **Standards:** ISO 8583 (card switching), ISO 20022 (clearing/settlement)

## Build & test commands

```sh
# Build everything
go build ./...

# Run all tests
go test ./...

# Run vet
go vet ./...

# Run tests in a Linux container (required on Windows — AppLocker blocks unsigned binaries)
docker build --target test -t clara-network-test .

# Run the full stack
docker compose -f deploy/docker-compose.yml up --build

# Individual services
go run ./cmd/switch
go run ./cmd/issuer-sim
go run ./cmd/acquirer-sim
go run ./cmd/adminapi
```

## Coding conventions

### Go style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Package names are lowercase, single-word: `clearing`, `ledger`, `cardsvc`.
- Exported types and functions use PascalCase; unexported use camelCase.
- Error handling: always check `err != nil`; wrap with `fmt.Errorf("context: %w", err)`.
- Prefer table-driven tests.
- No external frameworks — use only `net/http` stdlib for HTTP handlers.
- Direct dependencies are minimal: `pgx/v5` (Postgres), `go-redis/v9` (Redis).

### Configuration

All config is via `CLARA_*` environment variables. The `internal/env` package
provides helpers. Do not use `os.Getenv` directly in business logic — go
through the env helpers.

### Testing

- Tests live alongside the code they test (`foo.go` / `foo_test.go`).
- In-memory stores are preferred for unit tests.
- Integration tests use the `CLARA_PG_DSN` env var (skip gracefully if unset).
- Tests must pass on Linux. Windows requires WSL2 or Docker.

### Documentation

- The `docs/` directory contains the research and specification library.
- Keep documentation in sync with implementation changes.
- Architecture diagrams live in `docs/images/`.

### Security

- Never commit real keys, passwords, or DSN strings.
- The demo master key (`2b7e151628aed2a6abf7158809cf4f3c`) is for
  demonstration only — production deployments must use real HSM-backed keys.
- `CLARA_PG_DSN` in `deploy/docker-compose.yml` uses dummy credentials.

## Adding a new phase or feature

1. Create the core logic in `internal/<package>/`.
2. Create the CLI entry point in `cmd/<name>/`.
3. Write tests in `internal/<package>/_test.go`.
4. Add the service to `deploy/docker-compose.yml` if it's a long-running server.
5. Update `docs/27-implementation-status.md` with the new phase.
6. Add a Makefile target if useful.

## Common tasks

### Adding a new ISO 8583 message type

1. Add the MTI constant and field definitions in `internal/iso8583/`.
2. Add bitmap parsing/building support.
3. Add a handler in `internal/switchsrv/`.
4. Write round-trip tests.

### Adding a new API endpoint

1. Add the handler in `internal/adminapi/handlers.go`.
2. Add the store query in `internal/adminapi/store.go`.
3. Register the route in the server setup.
4. Add a test.
5. Update the README endpoint table.

### Adding a new clearing/settlement feature

1. Add the logic in `internal/clearing/`.
2. Update `internal/clearing/clearing_test.go`.
3. Add a new scenario to `cmd/clearing-sim/` if it's a demo-able flow.
4. Update the schema in `deploy/schema.sql` if new tables are needed.

## What NOT to do

- Do not add heavy dependencies (frameworks, ORM, etc.) — keep the stack minimal.
- Do not change the `CLARA_*` env var naming convention.
- Do not break backward compatibility of the ISO 8583 message format.
- Do not commit `.env` files or secrets.
- Do not modify `deploy/schema.sql` without a corresponding migration strategy.
- Do not use `go generate` or code generation tools without documenting them.
