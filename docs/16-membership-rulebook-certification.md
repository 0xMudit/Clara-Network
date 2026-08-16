# 16. Network Membership, Rulebook & Certification

## 16.1 Membership model

A card scheme's participants are **members** (usually financial institutions).
There are two broad tiers:

- **Principal (direct) members** — hold a settlement account (direct or via
  central bank), are assessed on the scheme's books, have full voting rights,
  and hold a BIN/ICA-style identifier.
- **Associate / sponsored members** — participate through a principal
  (sponsor). The sponsor remains financially liable to the scheme.

**Key membership mechanics:**

- **Sponsorship** — non-members (merchants, ISOs, PSPs) must be sponsored by a
  member. The member is responsible to the scheme for the sponsored party's
  behavior.
- **BIN/ICA assignment** — each member gets a unique institution identifier and
  BIN ranges (see `02-card-numbering-and-identification.md`).
- **Settlement accounts** — members hold accounts (often with the central bank)
  through which net/gross positions settle (see `18-settlement-liquidity-infrastructure.md`).
- **Member agreements** — signed contracts incorporating the rulebook by
  reference.

## 16.2 The rulebook

The rulebook is the legal contract that binds members. It typically covers:

- Obligations of issuers and acquirers (authorization, clearing, settlement,
  reporting).
- Merchant agreements and merchant obligations (acceptance, signage, MCC).
- Dispute resolution and chargeback rules (see `20-disputes-chargeback-management.md`).
- Data security requirements (PCI DSS alignment, see `09-pci-dss.md`).
- Sanctions and compliance (AML, OFAC/SDN screening).
- Fees, interchange, and scheme fees.
- Certification and testing obligations.
- Default, suspension, and termination procedures.

## 16.3 Certification & testing — the acceptance gate

Before an issuer or acquirer host can join the network, it must pass
**certification** proving it implements the network's message specs correctly.

### Host certification (acquirer + issuer processing)

Per the **U.S. Payments Forum "Payment Network Host and Level 3 Requirements"**
publication, network certification programs cover:

- **Network-level host certification** for both the acquirer-side and
  issuer-side hosts, run by the scheme (Amex, Discover, JCB, UPI, Mastercard,
  Visa all run similar programs).
- **EMV Level 3 terminal certification** — applies to chip-capable terminals
  (POS/ATM/contactless). L3 testing validates terminal application behavior
  against the EMV spec plus each network's application specifications using
  approved test tools/cards.
- Scope: authorization, clearing/settlement, reversal, chargeback, and
  exception flows; exact message formats per ISO 8583; field population rules;
  error handling.

### Typical certification stages

1. **Read the specification** (network interface specs / developer portals).
2. **Self-test** against the network's simulators / test environment.
3. **Submit test evidence** (transactions, logs, and results) for review.
4. **Live/parallel testing** in the scheme's pilot environment.
5. **Production go-live** — formal approval and member activation.

## 16.4 Running your own certification program

If you build a scheme, you need a certification function too:

- Maintain **test simulators** for ISO 8583/ISO 20022 message flows
  (see `04-iso8583.md`, `05-iso20022.md`).
- Define **test cards** (valid/invalid PANs, cryptogram vectors, PIN test
  values) — keyed to your crypto standards (see `17-message-security-key-management.md`).
- Define **negative test cases**: timeouts, malformed messages, unexpected
  response codes, offline approval scenarios.
- Run **periodic re-certification** when you change specs (major/minor version
  gating).
- Maintain a **certification register** of approved hosts, terminals, and
  test facilities.

## 16.5 Key references

- U.S. Payments Forum — *Payment Network Host and Level 3 Requirements*
  (Final, 09/2023) — summary of host/terminal certification across Amex,
  Discover, JCB, UPI, Mastercard, Visa.
- EMVCo — EMV Level 3 terminal testing framework and approved tools.
- Mastercard — certification portal and Chip/EMV specifications.
- Visa — developer portal, integration/certification guides.
