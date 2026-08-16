# 22. Card Production, Personalization & Lifecycle

## 22.1 The card production flow (G+D Convego model)

Card production is a secure, ceremony-controlled manufacturing process. The
**G+D Convego Issuance Partner** flow is a representative 10-step pipeline:

1. **Digital display** — the product (design, embossing, chip, contactless) is
   defined and previewed digitally.
2. **Request entry** — the issuing bank submits the card issuance request and
   cardholder data (via secure file/API).
3. **Data generation (production office)** — the production system generates
   the data needed for personalization: **EMV profile**, **keys**, **CVV/CVC**,
   **PVV**, expiry, embossing data — creating a **job ticket**.
4. **Production setup** — secure data loading into the personalization
   machines (data is encrypted and keyed into HSMs).
5. **Personalization** — printing, embossing, chip personalization (keys and
   EMV applications written), magnetic stripe encoding, contactless encoding.
6. **Quality & verification** — automated verification of every personalized
   element.
7. **Card mailer / packaging** — cards inserted into envelopes with PIN mailers
   (PIN is generated and mailed separately).
8. **Logistics** — secure dispatch to the cardholder or branch.
9. **Secure delivery** — delivery with track-and-trace and secure PIN delivery
   channels.
10. **Activation** — card activated by the issuer on first use, by phone, or
    via app (see lifecycle below).

Context: G+D reports ~**400 million cards/year** across **100+ banking
clients**, highlighting that production is an industrial operation with strict
security.

## 22.2 What gets personalized

| Element | Purpose | Security note |
|---------|---------|----------------|
| Magnetic stripe (track 1/2) | Fallback payment | Card verification data (CVC) encoded; sensitive |
| EMV chip applications | Chip payment (see `06-emv-chip.md`) | Chip keys, CVV/CVC, application data loaded via HSM |
| Contactless (NFC) | Tap payments | Same chip app data |
| CVV/CVC, PVV | Verification / PIN offset | Generated in an HSM, never stored in clear |
| Embossing / printing | Visual identification | PCI production standards apply |

## 22.3 Security & compliance for card production

- **PCI DSS scope** — card production facilities that store/process cardholder
  data must comply with PCI DSS and additional PCI card production guidelines.
- **PCI PIN Security** — PIN generation and PIN mailers are governed by the PCI
  PIN Security standard; PINs must be produced and mailed under dual control.
- **HSM-controlled keys** — all cryptographic material (chip keys, CVV keys,
  PVV keys) is generated, stored, and used inside HSMs
  (see `17-message-security-key-management.md`).
- **Key ceremonies** — card production keys are loaded under split-key dual
  control at the production site.
- **Physical security** — access-controlled facilities, escorted production,
  tamper-evident packaging, destroyed rejects.

## 22.4 Card lifecycle

| State | Description |
|-------|-------------|
| **Produced / personalized** | Data written, awaiting activation |
| **Activated** | First use or issuer-activated; token provisioning possible |
| **Active** | In normal use (chip, contactless, CNP) |
| **Suspended / blocked** | Lost/stolen, fraud-hold, or court order |
| **Expired** | Past expiry date; typically reissued automatically |
| **Cancelled / closed** | Account closed or card replaced |

Lifecycle automation matters operationally:

- **Reissue on expiry** — bulk re-personalization ahead of expiry.
- **Lost/stolen reissue** — block old card, reissue new card, refresh tokens
  and wallet cards (see `07-tokenization.md`).
- **Fraud-driven reissue** — block and reissue on compromise (see
  `10-fraud-risk.md`).

## 22.5 Instant & digital issuance

Modern schemes move beyond physical plastic:

- **Instant issuance** — a blank card is personalized at the branch from an
  HSM-backed service, or a virtual card is issued in-app immediately.
- **Digital card provisioning** — cards are tokenized into mobile wallets
  (Apple Pay, Google Pay) via secure provisioning APIs (see
  `07-tokenization.md`).
- **Virtual cards / spend controls** — issuer-defined limits, merchant
  restrictions, and expiry per virtual card.

## 22.6 Key references

- G+D — *Convego Issuance Partner* (card production & personalization).
- PCI SSC — *PCI DSS Requirements and Security Assessment Procedures*
  (CDE scope) and *PCI Card Production* guidance.
- PCI SSC — *PCI PIN Security Requirements* (PIN production & mailers).
- ISO 7816 — physical card and chip interface standards.
- Visa/Mastercard — chip personalization specs (CVN/cryptogram settings) for
  EMV applets.
