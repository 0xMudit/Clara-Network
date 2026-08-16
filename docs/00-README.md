# Building Your Own Mastercard/Visa-Style Card System - Documentation Index

This documentation set is a research library for designing and building a **card
payment scheme** (network) and the surrounding **issuing** and **acquiring**
infrastructure — an end-to-end system similar to Visa or Mastercard.

It was compiled from public sources including ISO standards documentation,
EMVCo specifications, PCI SSC publications, card scheme developer portals,
central-bank explainers (RBA, BIS CPMI, U.S. Federal Reserve), and engineering
articles on production payment system design.

## Document Map

| # | Document | What it covers |
|---|----------|----------------|
| 01 | [Card Payment Ecosystem & Four-Party Model](./01-card-payment-ecosystem.md) | Participants, roles, open vs closed networks, fee flows |
| 02 | [Card Numbering & Identification](./02-card-numbering-and-identification.md) | PAN structure, IIN/BIN, MII, Luhn check digit, BIN sponsorship |
| 03 | [Payment Flows: Authorization, Clearing, Settlement](./03-payment-flows.md) | The three phases and message exchange lifecycle |
| 04 | [ISO 8583 - Card-originated Interchange Messages](./04-iso8583.md) | The message standard for card transaction switching |
| 05 | [ISO 20022 - Payments Messaging](./05-iso20022.md) | XML message standard for interbank clearing & settlement |
| 06 | [EMV Chip & Contactless](./06-emv-chip.md) | Chip card spec, cryptograms, terminal transaction flow |
| 07 | [Tokenization & Network Tokens](./07-tokenization.md) | EMV payment tokens, TSP, PAR, network token services |
| 08 | [3-D Secure (EMV 3DS)](./08-3ds.md) | E-commerce cardholder authentication protocol |
| 09 | [PCI DSS Compliance](./09-pci-dss.md) | Security standard for the cardholder data environment |
| 10 | [Fraud Detection & Risk Management](./10-fraud-risk.md) | Real-time scoring, rules, velocity, machine learning |
| 11 | [Card Scheme Developer APIs](./11-apis-integration.md) | Mastercard/Visa developer platforms and integration patterns |
| 12 | [System Design: Building the Platform](./12-system-design.md) | Microservices, idempotency, ledger, reconciliation, scaling |
| 13 | [Fees, Interchange & Settlement Mechanics](./13-fees-interchange-settlement.md) | How money and fees actually move |
| 14 | [Glossary & Source References](./14-glossary-references.md) | Terminology and research sources |
| 15 | [Payment System Governance & Regulation](./15-payment-system-governance-regulation.md) | Licensing, PFMI principles, SIPS oversight |
| 16 | [Membership, Rulebook & Certification](./16-membership-rulebook-certification.md) | Member tiers, sponsorship, rulebook, host/L3 certification |
| 17 | [Key Management, HSM & Message Security](./17-message-security-key-management.md) | PIN blocks (ISO 9564), HSMs, MACs, key hierarchy |
| 18 | [Settlement & Liquidity Infrastructure](./18-settlement-liquidity-infrastructure.md) | RTGS accounts, prefunding, net/gross settlement, default fund |
| 19 | [Stand-In Processing & Availability](./19-stand-in-processing-availability.md) | Stand-in rules, response codes 91/P, operational resilience |
| 20 | [Disputes, Chargebacks & Arbitration](./20-disputes-chargeback-management.md) | Dispute lifecycle, reason codes, VCR, monitoring programs |
| 21 | [Cross-Border, FX & DCC](./21-cross-border-fx-dcc.md) | FX conversion, dynamic currency conversion, cross-currency settlement |
| 22 | [Card Production & Lifecycle](./22-card-production-lifecycle.md) | Personalization pipeline, EMV data, lifecycle, instant issuance |
| 23 | [Merchant Acquiring & Underwriting](./23-merchant-acquiring-underwriting.md) | Boarding, MATCH/OFAC screening, MCC, reserves, monitoring |
| 24 | [Instant Payments & Real-Time Processing](./24-instant-payments-rtp.md) | RTP/TIPS/Pix/UPI, prefunding, settlement models, ISO 20022 |

## How to use this library

- **Executives / product** → start with `01`, `03`, `13`, `15`.
- **Architects / engineers** → start with `12`, `04`, `03`, `06`, `18`, `24`.
- **Compliance / security** → start with `09`, `10`, `07`, `17`, `15`.
- **Integrators building on Visa/Mastercard rails** → `11`, `05`, `23`.
- **Operations / resilience** → `19`, `18`, `20`.

## End-to-end coverage matrix

Build-your-own-network checklist mapped to the docs:

| Capability needed to run a network | Where it's covered |
|------------------------------------|--------------------|
| Participants & four-party model | 01 |
| Card numbering, BIN assignment | 02 |
| Authorization / clearing / settlement flows | 03, 04, 13 |
| Messaging standards | 04 (ISO 8583), 05 (ISO 20022) |
| Chip / contactless / tokenization / 3DS | 06, 07, 08 |
| Data security (PCI DSS) & fraud | 09, 10 |
| Scheme APIs & platform architecture | 11, 12 |
| **Licensing & PFMI compliance** | 15 |
| **Membership, rulebook & certification** | 16 |
| **Crypto security: PIN, keys, HSMs, MACs** | 17 |
| **Settlement accounts, liquidity, default fund** | 18 |
| **Stand-in processing & resilience** | 19 |
| **Disputes, chargebacks & arbitration** | 20 |
| **Cross-border, FX & DCC** | 21 |
| **Card production & personalization** | 22 |
| **Merchant acquiring, underwriting & boarding** | 23 |
| **Instant / real-time payments** | 24 |

## Important disclaimer

Visa and Mastercard are **card schemes** (networks), not payment software you
can download. "Building your own" means one of these realistic paths:

1. **Network / scheme operator** - requires licensing, regulatory approval,
   central-bank sponsorship, members, and enormous scale. Not feasible for a
   startup; you effectively create a new ISO 8583 switch.
2. **Issuer / processor** - use a BIN sponsor, a sponsor bank, and a card
   processor, or build the issuing stack yourself.
3. **Acquirer / payment gateway / PSP** - the most common path; connect to
   Visa/Mastercard through a sponsor bank or directly.
4. **In-house card program** on existing rails via scheme APIs.

The realistic deliverable for most teams is a **payment platform** (option 2/3)
that interoperates with existing schemes. This library documents both the
full "from scratch" model and the build-vs-buy decision.
