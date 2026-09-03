# Security Policy

## Important Disclaimer

Clara Network is a **research and simulation project**. It is not a certified
production payment system and should not be deployed to handle real money,
cardholder data, or live transactions without significant additional hardening,
certification, and compliance work.

The HSM, key management, PIN block, and MAC code in `internal/hsm` are
**in-process simulations** for educational and architectural purposes. They
implement the published algorithms (AES key wrap, TR-31 key blocks, ISO 9564
PIN blocks, ISO 9797-1 MACs) but operate entirely in software memory. For
production use, these must be backed by a certified hardware security module
(e.g., Thales, Utimaco, Entrust).

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in Clara Network, please report it
responsibly. **Do not open a public GitHub issue for security vulnerabilities.**

### How to report

1. **Email** — Send a description to the maintainer at
   <https://github.com/0xMudit> via GitHub's private contact feature (click
   "Security" on the profile page).
2. **GitHub Security Advisories** — Use the
   [Security Advisories](https://github.com/0xMudit/Clara-Network/security/advisories/new)
   feature for private disclosure.

### What to include

- A description of the vulnerability and its potential impact.
- Steps to reproduce or a proof of concept.
- Any suggested fix, if you have one.

### What to expect

- **Acknowledgment** within 48 hours of your report.
- **Assessment** within 7 days, with a status update.
- **Resolution timeline** communicated once the issue is confirmed.
- Credit in the release notes (unless you prefer to remain anonymous).

### Scope

The following are in scope for security reports:

- Cryptographic implementation bugs in `internal/hsm` (key wrap, PIN blocks,
  MACs, key derivation).
- Authentication or authorization bypass in the switch or Admin API.
- SQL injection or data leakage through the Admin API or card service.
- Remote code execution or denial-of-service vectors.
- Information disclosure (e.g., PAN leakage, key material in logs).

The following are **out of scope**:

- The demo master key (`2b7e151628aed2a6abf7158809cf4f3c`) — this is
  explicitly documented as for demonstration only.
- The demo Postgres credentials in `deploy/docker-compose.yml` — these are
  for local development only.
- The in-memory fallback stores (by design, these are for simulation).

## Security Best Practices for Forks

If you fork Clara Network for a production or pre-production deployment:

1. **Replace all demo keys and credentials** before deployment.
2. **Enable TLS** on all network connections (the beta uses plaintext TCP).
3. **Back the HSM simulation with a real HSM** for any key material that
   protects cardholder data or transaction integrity.
4. **Enable PostgreSQL authentication** and restrict network access.
5. **Run `go vet ./...`** and address all warnings before deployment.
6. **Restrict Admin API access** — it currently has no authentication and
   exposes operational data.
7. **Review PCI DSS requirements** (`docs/09-pci-dss.md`) if you handle
   real cardholder data.

## Acknowledgments

We thank security researchers who report vulnerabilities responsibly. Your
efforts help make payment infrastructure safer for everyone.
