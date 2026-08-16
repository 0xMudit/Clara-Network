# 11. Card Scheme Developer APIs & Integration Patterns

## 11.1 Two integration worlds

1. **Scheme APIs** — Mastercard/Visa developer platforms for issuers,
   processors, and approved third parties (issuing, tokenization, transaction
   services, virtual cards, click-to-pay, etc.).
2. **Gateway/processor APIs** — REST/JSON APIs that let merchants/PSPs capture
   and process payments without touching the scheme directly (Mastercard
   Gateway, Stripe, Adyen, Worldpay, Checkout.com).

For "build your own card system," you typically need **1** (to interconnect as
an issuer/acquirer) or **2** (to run a merchant-facing PSP on existing rails).

## 11.2 Mastercard Developers platform

Developer portal: `developer.mastercard.com`. Authentication via **mTLS
certificates** and OAuth2 (client credentials) for newer APIs. Sandbox access
is instant; production requires on-boarding, certificate approval, and often a
mastercard approval process.

Key product families:

| Family | What it provides |
|--------|------------------|
| **Processing / MI Issuing** | Full issuing processing APIs: Card Issuance, Card Management, Card Controls, Authorization Management, Transaction Management (APMEA Digital First reference app on GitHub) |
| **Mastercard Digital Enablement Service (MDES)** | Network token provisioning, lifecycle, card-on-file (see `07-tokenization.md`) |
| **Mastercard Gateway (Open Payment API)** | REST/JSON payment gateway: payments, refunds, 3DS, tokenization, Click to Pay, "simple API", SDKs (JS, mobile) |
| **Mastercard Send** | Push payment (disbursement) API |
| **Settlement / Clearing** | IPM files (R111), Clearing Optimizer, settlement reporting (see `03-payment-flows.md`) |
| **Risk / Fraud** | Decision Intelligence, SPA alerts, chargeback APIs |

API mechanics:
- **mTLS** client certificates; field-level encryption for payloads
  (symmetric key wrap via public key).
- **OpenAPI specs** available for Issuing APIs; use code generators.
- Sandbox uses dummy data and instantly accessible credentials.

## 11.3 Visa Developer platform

Portal: `developer.visa.com`. Authentication via mutual TLS **with
certificates** and signature headers (SHA-256 with your private key; VDP
certificates). Sandbox + production projects.

Key product families:

| Family | What it provides |
|--------|------------------|
| **VisaNet Connect – Issuing** | Direct REST connectivity to VisaNet for issuers: Authorizations API (approve/decline), Card Services (card creation/update), with **HSM-as-a-Service** (PIN & cryptogram verification done by Visa) |
| **Token Service (VTS)** | Provisioning, ID&V, credential management, lifecycle, Token Reference ID (see `07-tokenization.md`) |
| **Visa Card Program Management (VCPE/VCDI)** | Instant/digital issuance, card enrollment, credential inquiry |
| **DPS Forward** | Full issuer processing (digital issuance, tokenization, wallet provisioning, ISO 20022 + REST APIs, stand-in processing) |
| **Visa Direct** | Push payments (fund disbursement / P2P) |
| **Visa B2B Connect** | Cross-border B2B payments (ISO 20022) |
| **Visa Risk Manager / VRM** | Issuer-side authorization controls, TC40/SAFE |

VisaNet Connect notes:
- REST APIs use **ISO 20022-style (ATICA)** naming conventions.
- Visa can verify CVV/CVV2/iCVV/CAAV or the issuer can (via HSM or
  HSM-as-a-Service).
- Clearing requests submitted by acquirers are processed by the same interface.

## 11.4 Payment gateway APIs (merchant-side)

The common merchant contract (used by Stripe, Mastercard Gateway, Adyen):

- **Idempotency-Key** header on every POST (see `12-system-design.md`).
- Resources: PaymentIntent/Charge, PaymentMethod/token, Capture, Refund,
  Dispute/Chargeback, Payout/Settlement, 3DS sessions.
- Webhooks for asynchronous events (at-least-once, signed).
- Hosted fields / JS SDK for **edge tokenization** (PCI scope reduction).
- Risk products (e.g. Stripe Radar) and network-token features.

## 11.5 Direct acquiring vs sponsor vs PSP

| Path | Description | Best when |
|------|-------------|-----------|
| **Direct scheme member** | Become an acquirer/issuer member of Visa/Mastercard (bank or regulated entity with settlement accounts) | You are (or partner with) a bank; full control, full liability |
| **Sponsor bank** | Non-bank uses a bank's scheme membership (BIN sponsorship for issuing; acquiring sponsorship for merchants) | Fast to market without banking license |
| **PSP/PayFac** | Aggregate under a PSP's acquiring (Stripe, Adyen); you get merchant APIs, no scheme membership | Building merchant products quickly |
| **In-house switch** | Build an ISO 8583 switch for your own network (new scheme) | Building a *new* network (not interoperating with Visa/MC) |

## 11.6 Integration checklist

- [ ] Determine role: issuer, acquirer, PSP, TSP, switch, or new scheme.
- [ ] Select build-vs-buy: scheme APIs (build own stack), processor (use their
      platform), or PSP (product on rails).
- [ ] Onboard developer portals: mTLS certs, sandbox keys, production approval.
- [ ] Certify with the scheme (Level 3 tests, BIN certification).
- [ ] Implement auth flow (REST/ISO 20022 for VisaNet Connect; ISO 8583/IPM for
      direct switch; gateway API for PSP path).
- [ ] Implement clearing/settlement feeds + reconciliation (BASE II/IPM files,
      net settlement reports).
- [ ] Implement tokenization + 3DS integrations.
- [ ] Fraud, chargebacks, and disputes API integrations.
- [ ] PCI compliance and scheme rule compliance programs.

## 11.7 Key references

- Mastercard Developers — `developer.mastercard.com` (products, OpenAPI specs,
  issuing reference app on GitHub)
- Visa Developer — `developer.visa.com` (VisaNet Connect – Issuing, VTS, DPS
  Forward, Card Program Management)
- Mastercard Gateway — Open Payment API reference (REST/JSON)
- W3C — Payment Request API / Tokenized Card Payment
