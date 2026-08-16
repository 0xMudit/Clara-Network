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

## Status

Research & specification phase. Contributions are welcome.

## License

[MIT](./LICENSE) — see the LICENSE file for details.
