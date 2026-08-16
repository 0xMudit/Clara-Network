# Clara Network

**Clara Network** is an open-source project to design and build a
Mastercard/Visa-style card payment network end-to-end: scheme (network)
operator, issuer, and acquirer infrastructure.

## What's in this repository

The repo currently holds the **research and specification library** for the
network under [`docs/`](./docs/):

- **25 numbered technical documents** (`docs/00-README.md` → `docs/24-...`)
  covering the full stack: four-party model, BIN/PAN numbering, authorization /
  clearing / settlement flows, ISO 8583 & ISO 20022, EMV chip & contactless,
  tokenization, 3-D Secure, PCI DSS, fraud & risk, scheme APIs, system design,
  fees/interchange, governance & PFMI, membership & certification, key
  management & HSM, settlement & liquidity, stand-in processing, disputes &
  chargebacks, cross-border & DCC, card production, merchant acquiring, and
  instant payments / RTP.
- **Reference library** (`docs/references/`) — source PDFs downloaded from
  public institutions (BIS/CPMI, ECB, World Bank, FDIC, OCC, Visa, Mastercard,
  PCI SSC, NIST) plus a manifest tracking every needed source, including
  member-gated and paywalled documents.

Start with [`docs/00-README.md`](./docs/00-README.md).

## Building & running

Phase 1 implements the **ISO 8583 message switch** end-to-end: acquirer →
switch → issuer authorization, idempotent replay, stand-in processing, and
audit logging. It requires Go 1.26+ and Docker Desktop.

```sh
# unit + integration tests
go test ./...

# run a full round-trip using Docker (postgres, redis, switch, issuer-sim,
# acquirer-sim), then watch the switch logs for auth responses
docker compose -f deploy/docker-compose.yml up --build
docker compose -f deploy/docker-compose.yml logs switch

# or run locally: terminal 1 -> switch, terminal 2 -> issuer-sim,
# terminal 3 -> acquirer-sim
go run ./cmd/switch
go run ./cmd/issuer-sim
go run ./cmd/acquirer-sim
```

Key config (via env): `CLARA_ISSUER_ROUTES` (JSON `{receiving-institution-id:
host:port}`), `CLARA_REDIS_ADDR` (idempotency), `CLARA_PG_DSN` (audit log).

## Status

Research & specification library (docs 00–24) plus phase 1 (ISO 8583 switch)
implemented. Contributions are welcome.

## License

[MIT](./LICENSE) — see the LICENSE file for details.
