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

Phases 1–2 implement the **ISO 8583 authorization flow** end-to-end:
acquirer → switch → (risk check) → issuer authorization, BIN-based routing,
issuer failover, idempotent replay, in-path risk scoring (velocity), stand-in
processing, and audit logging. It requires Go 1.26+ and Docker Desktop.

```sh
# unit + integration tests
go test ./...

# run a full round-trip using Docker (postgres, redis, switch, issuer-sim,
# acquirer-sim), then watch the logs: 6 auth requests with BIN routing and
# a velocity rule that declines the 6th with response code 59
docker compose -f deploy/docker-compose.yml up --build
docker compose -f deploy/docker-compose.yml logs switch acquirer-sim

# or run locally: terminal 1 -> switch, terminal 2 -> issuer-sim,
# terminal 3 -> acquirer-sim
go run ./cmd/switch
go run ./cmd/issuer-sim
go run ./cmd/acquirer-sim
```

Key config (via env):

- `CLARA_ISSUER_ROUTES` — JSON `{receiving-institution-id:
  host:port}`; a value may be a comma-separated failover list.
- `CLARA_BIN_TABLE` — JSON `{"entries":{"400000":"1000001000"}}` routes by
  PAN BIN when the message omits DE100.
- `CLARA_RISK_RULES` — JSON rule set; velocity counters (per card / per
  merchant) are counted in Redis and can decline with a configurable code.
- `CLARA_REDIS_ADDR` — idempotency + risk counters.
- `CLARA_PG_DSN` — audit log.
- `CLARA_SEND_DE100=false` (acquirer-sim) — omit DE100 to exercise BIN routing.

## Status

Research & specification library (docs 00–24), phase 1 (ISO 8583 switch),
and phase 2 (authorization flow with BIN routing, risk, failover)
implemented. Contributions are welcome.

## License

[MIT](./LICENSE) — see the LICENSE file for details.
