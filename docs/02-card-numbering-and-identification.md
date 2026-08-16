# 2. Card Numbering & Identification (PAN, IIN/BIN, Luhn)

## 2.1 Primary Account Number (PAN) structure

A **PAN** (payment card number) is composed of 8 to 19 digits. It is allocated
according to **ISO/IEC 7812**. Its three parts:

```
┌─────────────┬───────────────────────────────┬──────────┐
│ IIN (BIN)   │  Individual account number    │ Check    │
│ 6 or 8 digs │  (up to 12 digits)            │ digit    │
└─────────────┴───────────────────────────────┴──────────┘
```

- **Issuer Identification Number (IIN)** — also called **Bank Identification
  Number (BIN)** by Visa/Mastercard. The leading 6 digits (historically) or 8
  digits (per the 2017 5th edition of ISO/IEC 7812-1) identify the issuing
  institution. The **first digit** is the **Major Industry Identifier (MII)**.
- **Individual account number** — the rest of the PAN (minus check digit),
  assigned by the issuer to identify the cardholder account.
- **Check digit** — a single digit computed with the **Luhn (mod-10)** algorithm
  (see 2.3).

Note: an **8-digit IIN** was standardized in ISO/IEC 7812-1:2017 (IIN length
went 6 → 8, minimum PAN length went 8 → 10). Visa moved to assigning only
8-digit BINs after **April 2022**, while continuing to support legacy 6-digit
BINs.

## 2.2 Major Industry Identifier (MII)

The first digit of the PAN identifies the broad industry of the issuer:

| MII | Industry |
|-----|----------|
| 0 | ISO/TC 68 and other future industry assignments |
| 1 | Airlines |
| 2 | Airlines and banking/financial |
| 3 | Travel & entertainment and banking/financial |
| 4 | Banking/financial |
| 5 | Banking/financial |
| 6 | Merchandising and banking/financial |
| 7 | Petroleum and other future assignments |
| 8 | Healthcare, telecommunications |
| 9 | National assignment (national standards bodies) |

## 2.3 Luhn (mod-10) check digit algorithm

The Luhn algorithm is a checksum that catches any single-digit error and most
adjacent transposition errors. It is **not** security — it only catches typos.

Procedure (from ISO/IEC 7812-1 Annex B and public references):

1. Write all digits except the last (the check digit).
2. Starting from the right, **double every other digit** (the 2nd, 4th, 6th…
   from the right, i.e. alternating positions).
3. For any digit that doubles to ≥ 10, **sum its digits** (e.g. 7 → 14 → 1+4 = 5).
4. Sum all resulting single digits.
5. Add the check digit. If the total is a **multiple of 10**, the PAN is valid.

Example (pseudo-code):

```text
input: digits d[0..n-1]  (check digit = d[n-1])
sum = 0
parity = n-1 mod 2
for i in 0..n-2:
    d = d[i]
    if i mod 2 == parity:
        d *= 2
        if d > 9: d -= 9
    sum += d
check_digit = (10 - (sum mod 10)) mod 10
```

Implementation in most payment SDKs and card-validation libraries; it is also
used for PAN validation in ISO 8583 systems (in addition to scheme-based BIN
checks).

## 2.4 Network BIN/IIN ranges (from public sources)

| Issuing network | IIN/BIN ranges | Typical PAN length | Validation |
|-----------------|----------------|--------------------|------------|
| Visa | `4` (e.g. 4xxx) | 13, 16, 19 | Luhn |
| Visa Electron | 4026, 417500, 4844, 4913, 4917 | 16 | Luhn |
| Mastercard | `51–55`, `2221–2720` (2-series since 2016) | 16 | Luhn |
| American Express | `34`, `37` | 15 | Luhn |
| Discover | 6011, 644–649, 65 | 16–19 | Luhn |
| UnionPay | `62` | 16–19 | Luhn (per scheme rules) |
| Maestro | 5018, 5020, 5038, 5893, 6304, 6759, 6761–6763 | 12–19 | Luhn |
| Diners Club | 30, 36, 38, 39 | 14–19 | Luhn |
| JCB | 3528–3589 | 16–19 | Luhn |
| Mir (Russia) | 2200–2204 | 16–19 | Luhn |

> BIN ranges above are for **routing/identification** during development and
> testing. In production you can only use ranges licensed to you by a scheme
> and/or a BIN sponsor.

## 2.5 How a network routes on BIN

When a transaction arrives at the scheme, the network:

1. Reads the **BIN (first 6–8 digits)** from the PAN.
2. Looks up the BIN in its **BIN routing table** to find the responsible issuer
   (and its node address / institution ID).
3. Routes the authorization/clearing message to that issuer.

When building your own scheme or switch, you must maintain a **BIN table**
mapping BIN → issuer node, merchant/acquirer IDs → acquirer node, plus a
**scheme routing plan** (primary + alternate routes).

## 2.6 BIN sponsorship

Getting your own BIN from a scheme/ISO is costly and demands regulatory and
compliance overhead. A **BIN sponsor** (a regulated bank) sub-licenses part of
its BIN range to you. The sponsor is contractually liable to the network for
your card program. BIN sponsorship is the standard fast path for fintechs to
issue cards without being a bank.

- Benefits: fast to market, no direct scheme/ISO application.
- Costs: sponsor fees, scheme fees, and the sponsor controls (or heavily
  monitors) your compliance, KYC/AML, and program rules.
- Visa/Mastercard treat sponsor arrangements via their **issuing rules**;
  the PAN and tokens remain 16 digits regardless of 6- or 8-digit BIN.

## 2.7 Related identifiers

- **Acquiring Identifier** - 6-digit identifiers Visa uses for acquirers (no
  longer called BINs; only *issuing* numbers are ISO BINs).
- **Merchant Category Code (MCC)** - 4-digit code (ISO 18245) describing the
  merchant's business type; drives interchange tiering and risk (e.g. 5411
  supermarkets, 5812 restaurants, 5999 misc.).
- **Merchant ID (MID)** - the acquirer's unique ID for a merchant.
- **Terminal ID (TID)** - identifies a POS terminal.
- **Token Reference ID / Payment Account Reference (PAR)** - see
  `07-tokenization.md`.

## 2.8 Best practices for building an in-house card numbering service

- Use a **sequence allocator** (e.g. HSM-backed or distributed ID service) for
  individual account numbers, never sequential guessable numbers in production.
- Compute and validate the **Luhn digit** at issuance and at entry.
- Store the **BIN table** in a highly available service with sub-ms lookup;
  it is on the hot path for routing.
- For testing, use **official test card numbers** from Visa/Mastercard
  (e.g. `4111 1111 1111 1111` Visa, `5555 5555 5555 4444` Mastercard) and never
  real cardholder data.
- Consider **tokenized issuance** (see `07-tokenization.md`) so PANs are never
  exposed to merchant environments.

## 2.9 Key references

- ISO/IEC 7812-1:2017 (Identification cards — Identification of issuers —
  Part 1: Numbering system)
- Wikipedia - "Payment card number"
- Visa - "Eight Digit BIN Frequently Asked Questions"
- Scientific American - "What Is the Luhn Algorithm?"
- Barclaycard UK - "BIN Rules" document (BIN range tables)
