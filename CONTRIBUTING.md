# Contributing to Clara Network

Thank you for your interest in contributing to Clara Network! This guide covers
everything you need to get started — from setting up your local environment to
submitting your first pull request.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Local Development Setup](#local-development-setup)
- [Project Architecture](#project-architecture)
- [Build Phases and Code Map](#build-phases-and-code-map)
- [Running the Full Stack](#running-the-full-stack)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Branch Naming](#branch-naming)
- [Pull Request Process](#pull-request-process)
- [Definition of Done](#definition-of-done)
- [Finding Things to Work On](#finding-things-to-work-on)
- [Getting Help](#getting-help)

## Code of Conduct

This project follows the
[Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating,
you agree to uphold it.

## Quick Start

```sh
# Clone the repo
git clone https://github.com/0xMudit/Clara-Network.git
cd Clara-Network

# Run the full stack
docker compose -f deploy/docker-compose.yml up --build

# Watch the demo
docker compose -f deploy/docker-compose.yml logs -f acquirer-sim

# In another terminal, hit the Admin API
curl http://localhost:18083/api/v1/dashboard
```

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.26+ | `go version` to check |
| Docker Desktop | Latest | Linux containers enabled |
| Git | 2.x+ | `git version` to check |
| Make | Optional | For convenience targets |

> **Windows users:** Tests must run on Linux. Use WSL2, Docker, or
> `docker build --target test` to run the test suite. Windows hosts with
> AppLocker or Microsoft Defender Application Control block unsigned test
> binaries.

## Local Development Setup

### Option A: Docker (recommended for first-time contributors)

```sh
# Start the full stack (Postgres, Redis, switch, issuer, acquirer, all sims)
docker compose -f deploy/docker-compose.yml up --build

# Verify everything is running
curl http://localhost:18083/health          # Admin API
curl http://localhost:18083/api/v1/dashboard  # Dashboard summary
```

### Option B: Native Go

```sh
# Install dependencies
go mod download

# Run tests
go test ./...

# Run individual services in separate terminals
go run ./cmd/switch          # Terminal 1 — the network core
go run ./cmd/issuer-sim      # Terminal 2 — issuer host
go run ./cmd/acquirer-sim    # Terminal 3 — drives authorization requests
go run ./cmd/adminapi        # Terminal 4 — Admin REST API
```

### Option C: Run specific sims

Each sim is self-contained and exits when done:

```sh
go run ./cmd/clearing-sim     # Clearing + net settlement demo
go run ./cmd/ledger-sim       # Ledger + reconciliation demo
go run ./cmd/card-sim         # Card personalization + token demo
go run ./cmd/acquiring-sim    # Merchant boarding + funding demo
go run ./cmd/disputes-sim     # Disputes + arbitration demo
go run ./cmd/hsm-sim          # HSM key ceremonies demo
go run ./cmd/resilience-sim   # Circuit breaker + stand-in chaos drill
go run ./cmd/instant-sim      # Instant payments (RTP) demo
```

## Project Architecture

```
clara-network/
├── cmd/                        # Entry points (one per microservice)
│   ├── switch/                 # ISO 8583 message switch (the network core)
│   ├── issuer-sim/             # Issuer host simulator
│   ├── acquirer-sim/           # Acquirer host simulator (client)
│   ├── cardsvc/                # Card service REST API
│   ├── card-sim/               # Card service demo client
│   ├── clearing-sim/           # Clearing + net settlement demo
│   ├── ledger-sim/             # Ledger + reconciliation demo
│   ├── acquiring-sim/          # Merchant boarding / funding demo
│   ├── disputes-sim/           # Disputes + arbitration demo
│   ├── hsm-sim/                # HSM simulation demo
│   ├── resilience-sim/         # Outage chaos drill demo
│   ├── instant-sim/            # Instant payments demo
│   └── adminapi/               # Read-only Admin REST API
├── internal/                   # Core library packages
│   ├── iso8583/                # ISO 8583 message model, bitmap parse/build
│   ├── framing/                # 2-byte length-prefixed TCP framing
│   ├── switchsrv/              # Switch: routing, failover, idempotency, risk, stand-in
│   ├── binrouting/             # BIN → issuing-institution routing table
│   ├── risk/                   # Rule engine: velocity counters (memory/Redis)
│   ├── acquirersim/            # Acquirer host logic
│   ├── issuersim/              # Issuer host logic
│   ├── clearing/               # Clearing capture, netting, prefund caps, default fund, pacs.009
│   ├── ledger/                 # Append-only double-entry ledger, reconciliation
│   ├── cardsvc/                # BIN ranges, personalization, ARQC, token vault, provisioning
│   ├── acquiring/              # Merchant boarding, MATCH/OFAC screening, MCC tiering, funding
│   ├── disputes/               # Reason codes, representment, arbitration, monitoring
│   ├── hsm/                    # In-process HSM simulation
│   ├── resilience/             # Stand-in, circuit breakers, metrics, 91-burst detection
│   ├── instant/                # ISO 20022 pacs.008/002, prefunded RTP engine
│   ├── adminapi/               # Admin REST API: server, handlers, store
│   └── env/                    # CLARA_* config helpers
├── deploy/                     # Docker Compose stack (14 services)
├── docs/                       # 27-document research & specification library
├── web/                        # Next.js admin dashboard frontend
├── Dockerfile                  # Multi-stage build for every cmd
├── Makefile                    # Build/test/vet/run convenience targets
└── go.mod                      # Module: github.com/0xMudit/Clara-Network
```

## Build Phases and Code Map

Each of the 10 completed phases maps to specific packages. Use this to find
where a given feature lives:

| Phase | Feature | Core Package(s) | Demo Command |
|-------|---------|------------------|--------------|
| 1 | Switch skeleton + ISO 8583 | `internal/iso8583`, `internal/framing` | `cmd/switch` |
| 2 | Authorization flow | `internal/switchsrv`, `internal/binrouting`, `internal/risk` | `cmd/acquirer-sim` |
| 3 | Clearing + net settlement | `internal/clearing` | `cmd/clearing-sim` |
| 4 | Ledger + reconciliation | `internal/ledger` | `cmd/ledger-sim` |
| 5 | Issuing stack (cards, tokens) | `internal/cardsvc` | `cmd/cardsvc`, `cmd/card-sim` |
| 6 | Acquiring stack (merchants) | `internal/acquiring` | `cmd/acquiring-sim` |
| 7 | Disputes engine | `internal/disputes` | `cmd/disputes-sim` |
| 8 | Key management + HSM | `internal/hsm` | `cmd/hsm-sim` |
| 9 | Resilience (stand-in, circuits) | `internal/resilience` | `cmd/resilience-sim` |
| 10 | Instant payments (RTP) | `internal/instant` | `cmd/instant-sim` |

For a detailed contributor on-ramp, see
[`docs/28-contributor-architecture-guide.md`](docs/28-contributor-architecture-guide.md).

## Running the Full Stack

### What a successful demo run looks like

```sh
docker compose -f deploy/docker-compose.yml up --build
```

When you run the full stack, the one-shot sims execute in sequence and log
their progress:

1. **Acquirer-sim** sends 6 authorization requests through the switch with BIN
   routing and a velocity rule. The first 5 are approved (`response_code: 00`);
   the 6th is declined (`response_code: 59`) by the velocity rule.

2. **Clearing-sim** processes a clearing cycle, computes net positions for each
   member, and writes ISO 20022 pacs.009 settlement instructions as XML.

3. **Ledger-sim** posts every net position as a balanced double-entry journal
   and reconciles against the settlement agent's statement.

4. **Card-sim** personalizes a card, verifies an ARQC cryptogram, creates a
   network token, and provisions it to a mobile wallet.

5. **Acquiring-sim** boards merchants with MATCH/OFAC screening, assigns MCCs,
   and runs the funding engine (fees, reserves, payouts).

6. **Disputes-sim** files disputes, processes representments, runs arbitration,
   and monitors chargeback ratios.

7. **Hsm-sim** performs key ceremonies (dual-control M-of-N), PIN block
   verification, retail MAC with tamper detection, key rotation, and
   zeroize.

8. **Resilience-sim** runs a chaos drill: kills the primary issuer (circuit
   breaker trips, failover), kills the secondary (stand-in + 91-burst alert),
   then recovers (half-open probe re-closes the circuit).

9. **Instant-sim** settles pacs.008 credit transfers in real time, tests every
   rejection code, exercises SLA timeout with reservation release, and
   verifies position conservation.

All sims should exit with code 0. The switch, issuer-sim, cardsvc, Postgres,
and Redis stay running.

### Verifying the Admin API

```sh
# Dashboard summary
curl http://localhost:18083/api/v1/dashboard

# Expected output (counts vary by demo run):
# {"transactions":6,"clearingRecords":...,"merchants":...,"disputes":...,"cards":...,"tokens":...}

# Switch transaction audit log
curl http://localhost:18083/api/v1/transactions?limit=3

# Health check
curl http://localhost:18083/health
```

## Running Tests

```sh
# Run all tests (must run on Linux — use WSL2 or Docker on Windows)
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests for a specific package
go test ./internal/clearing/ -v

# Run tests in Docker (works on all platforms)
docker build --target test -t clara-network-test .

# Run vet
go vet ./...

# Build everything
go build ./...
```

**Current test count:** 128 tests across 12 packages, all green.

## Code Style

### Go conventions

- Follow standard Go conventions: `gofmt`, `go vet` pass cleanly.
- Package names are lowercase, single-word: `clearing`, `ledger`, `cardsvc`.
- Use table-driven tests where practical.
- Wrap errors with context: `fmt.Errorf("clearing net positions: %w", err)`.
- Prefer the standard library — external dependencies are limited to
  `pgx/v5` (PostgreSQL) and `go-redis/v9` (Redis).

### Configuration

All configuration uses `CLARA_*` environment variables. Use the helpers in
`internal/env/` — do not call `os.Getenv` directly in business logic.

### Security

- Never commit real keys, passwords, or DSN strings.
- The demo master key (`2b7e151628aed2a6abf7158809cf4f3c`) is for
  demonstration only.
- `CLARA_PG_DSN` in `deploy/docker-compose.yml` uses dummy credentials.

### Documentation

- Keep `docs/27-implementation-status.md` in sync with implementation changes.
- New config variables should be documented in both the README and `docs/27`.
- Architecture diagrams live in `docs/images/`.

## Commit Messages

Use clear, descriptive commit messages following Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:** `feat`, `fix`, `docs`, `test`, `refactor`, `ci`, `chore`

**Examples:**

```
feat(clearing): add multi-currency netting support

Adds support for netting positions across different currencies within
the same clearing cycle. Uses the settlement agent's FX rates to
convert to a base currency before computing net positions.

Closes #123
```

```
test(hsm): add AES key wrap round-trip test

Verifies that a key wrapped with RFC 3394 can be unwrapped to the
original value, covering both single-key and key-bar cases.

Relates to #456
```

- Keep the subject line under 72 characters.
- Wrap the body at 80 characters.
- Reference issues with `Closes #N` or `Relates to #N`.

## Branch Naming

Use descriptive branch names with a type prefix:

| Prefix | Use for |
|--------|---------|
| `feat/` | New features |
| `fix/` | Bug fixes |
| `docs/` | Documentation changes |
| `test/` | Test additions or fixes |
| `refactor/` | Code refactoring |
| `ci/` | CI/infrastructure changes |

**Examples:**
- `feat/clearing-multi-currency`
- `fix/hsm-key-rotation-audit`
- `docs/phase-11-fraud-ml`
- `test/instant-concurrent-settlement`

## Pull Request Process

1. **Fork** the repository and create a branch from `main`.
2. **Make your changes** with tests.
3. **Run the full test suite** and verify `go vet` is clean.
4. **Update documentation** if your change affects public behavior, config, or
   architecture.
5. **Open a PR** against `main` with a clear title and description.
6. **Fill out the PR template** checklist.
7. **Wait for CI** to pass and for a maintainer review.

### Review expectations

- PRs require at least one maintainer review.
- CI must pass (`go vet`, `go test`, `golangci-lint`, Docker smoke test).
- Large changes may be split into smaller PRs for easier review.

## Definition of Done

A pull request is considered complete when:

- [ ] `go test ./...` passes (including `-race` on CI)
- [ ] `go vet ./...` is clean
- [ ] `golangci-lint run` reports no new warnings
- [ ] New code has corresponding tests
- [ ] Documentation is updated (if the change affects public behavior, config,
      or architecture)
- [ ] `deploy/schema.sql` is updated (if new tables or columns were added)
- [ ] `docs/27-implementation-status.md` is updated (if a new phase or
      significant feature was added)
- [ ] No secrets, keys, or credentials are committed
- [ ] The PR description explains the "why" behind the change

## Finding Things to Work On

- **Good first issues** are tagged [`good first issue`](https://github.com/0xMudit/Clara-Network/labels/good%20first%20issue) on GitHub.
- **Help wanted** issues are tagged [`help wanted`](https://github.com/0xMudit/Clara-Network/labels/help%20wanted).
- Check the [ROADMAP.md](ROADMAP.md) for planned features.
- Browse the [open issues](https://github.com/0xMudit/Clara-Network/issues)
  for bugs and feature requests.

### High-impact contribution areas

| Area | Examples | Difficulty |
|------|----------|------------|
| New ISO 8583 fields | Add DE support, new MTI types | Medium |
| Test coverage | Edge cases, fuzz tests, integration tests | Easy–Medium |
| Documentation | Translations, tutorials, config var docs | Easy |
| Admin API endpoints | New query endpoints, filtering, pagination | Medium |
| Risk engine | New rule types, ML scoring hooks | Hard |
| Deployment | Kubernetes manifests, Terraform, Helm charts | Medium |
| Web dashboard | UI improvements, new views | Easy–Medium |

## Getting Help

- **Issues** — [GitHub Issues](https://github.com/0xMudit/Clara-Network/issues)
- **Discussions** — [GitHub Discussions](https://github.com/0xMudit/Clara-Network/discussions)
- **Security** — [Security Advisories](https://github.com/0xMudit/Clara-Network/security/advisories/new)

Thank you for helping build open payment infrastructure! 🚀
