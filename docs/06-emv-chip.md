# 6. EMV Chip & Contactless Payments

## 6.1 What EMV is

**EMV** (Europay, Mastercard, Visa) is the global chip-card specification owned
and maintained by **EMVCo** (owned by Visa, Mastercard, Amex, Discover, JCB,
UnionPay). It defines how payment cards and terminals interoperate for
**card-present** payments — contact (insert/dip) and contactless (tap/NFC).

Core security property: the chip generates a **one-time cryptogram** per
transaction, which authenticates the card and defeats skimming/counterfeiting.
> "Chip technology validates the authenticity of a card and generates a
> one-time use security code for every transaction... The security code is
> unique to each transaction and cannot be reused." — EMVCo

Underlying standards: **ISO/IEC 7816** (contact cards), **ISO/IEC 14443**
(contactless/NFC), **ISO 8583** (host messages).

## 6.2 The four EMV Book sets

- **EMV Integrated Circuit Card (ICC) Specifications for Payment Systems**
  (the "EMV Book"):
  - Part I — Electromechanical characteristics, logical interface, transmission
    protocols
  - Part II — Data elements & commands
  - Part III — Application selection
  - Part IV — Security aspects / secure messaging
- **Book 2: Security and Key Management**
- **Book 3: Application Specification** — data elements, commands, debit/credit
  application, transaction flow, AC generation
- **Book 4: Cardholder, Attendant & Acquirer Interface**
- Contactless: **Book A/B (entry/selection), C (kernel), D (contactless
  communications), E (contactless payment)**.
- Plus: **EMV Payment Tokenisation** and **EMV 3-D Secure** (see 07, 08).

## 6.3 Card communication basics

- **Contact** — card inserted into reader; communication over ISO 7816 with
  T=0 (character-oriented) or T=1 (block-oriented) protocol; session starts
  with **ATR (Answer to Reset)**.
- **Contactless** — card tapped; ISO 14443 NFC; reader sends **RATS (Request
  for Answer To Select)**, card responds with **ATS (Answer To Select)**.
- Data objects are exchanged in **TLV (Tag-Length-Value)** format.

### 6.3.1 Application selection

The terminal builds a **candidate list** of card applications:

1. **Indirect/implicit selection via PSE** (Payment System Environment):
   terminal sends `SELECT` for `1PAY.SYS.DDF01` (contact) or `2PAY.SYS.DDF01`
   (contactless); if present, reads the directory and app candidates.
2. **Direct selection via PSA** (Payment System Application): terminal sends
   `SELECT` with the **AID** (Application Identifier).

**AID** = **RID** (5-byte Registered Application Provider Identifier, e.g.
`A000000004` = Mastercard, `A000000003` = Visa) + **PIX** (Proprietary
Application Identifier Extension, e.g. `1010` Mastercard Credit, `2010`
Mastercard Debit, `1010` Visa classic). AIDs are 5–16 bytes.

Selection rules: if multiple apps, either **cardholder selection** (manual) or
**automatic selection** using the **Application Priority Indicator (tag 87)** —
lowest value wins. Partial-name matching supports extended AIDs.

### 6.3.2 Key EMV commands (APDUs)

| Command | Purpose |
|---------|---------|
| `SELECT` | Select PSE or ADF (application) |
| `GET PROCESSING OPTIONS` (GPO, 80 A8 00 00) | Start transaction; terminal sends **PDOL** data; card returns **AIP** + **AFL** |
| `READ RECORD` | Read files pointed to by the AFL |
| `GENERATE AC` (80 AE 80/00) | Card produces **Application Cryptogram (AC)** |
| `EXTERNAL AUTHENTICATE` | Card authenticates issuer (online) |
| `GET DATA` | Retrieve issuer data |
| `PUT DATA` | Store data (e.g. issuer scripts) |
| `VERIFY` | Offline PIN verification |
| `GET CHALLENGE` | Card supplies random for CDAs/DDA |

## 6.4 The EMV transaction flow (terminal side)

Standard flow (Book 3, application spec):

1. **Initiate Application Processing** — terminal sends GPO with PDOL data
   (e.g. terminal capabilities, country code, type, amount, currency); card
   returns **AIP** (Application Interchange Profile) and **AFL** (Application
   File Locator) telling the terminal which files to read.
2. **Read Application Data** — terminal reads the files in the AFL (RECORD
   data: PAN, expiry, cardholder name, risk parameters, etc.).
3. **Card Data Authentication** — verify the card is genuine:
   - **SDA** (Static Data Authentication) — RSA signature over static data
   - **DDA** (Dynamic Data Authentication) — card signs a random challenge
   - **CDA** (Combined DDA/AC) — signature combined into the AC (strongest)
4. **Processing Restrictions** — check application version, usage control,
   effective/expiration dates.
5. **Cardholder Verification Method (CVM)** — verify the cardholder:
   - Offline PIN (plaintext or encrypted)
   - Online PIN (sent to issuer via DE 52/DE 55)
   - Signature
   - No CVM (e.g. low-value contactless)
   Results recorded in **CVR** (Cardholder Verification Results).
6. **Terminal Risk Management** — floor limits, random transaction selection,
   velocity checking (accumulated off-line limit). Results in **TVR**
   (Terminal Verification Results).
7. **Terminal Action Analysis** — terminal compares TVR against terminal/issuer
   action codes and decides: offline approve, offline decline, or go online.
8. **Card Action Analysis** — card compares TVR against **Issuer Action Codes**
   and its own risk management, decides:
   - **TC (Transaction Certificate)** — approve offline
   - **AAC (Application Authentication Cryptogram)** — decline
   - **ARQC (Authorization Request Cryptogram)** — go online
9. **Online Processing** — if ARQC: terminal sends authorization request to
   issuer with cryptogram (DE 55); issuer validates cryptogram (via HSM) and
   responds approve/decline; **issuer authentication** via `EXTERNAL
   AUTHENTICATE` if the card supports issuer scripts.
10. **Issuer Script Processing** — issuer may send scripts to update card
    parameters/keys.
11. **Completion** — card issues **TC** (or **AAC**); terminal stores data for
    clearing.

The **TVR**, **TSI** (Transaction Status Information), **AIP**, **CVR** and the
**AC (ARQC/TC/AAC)** travel in DE 55 (or DE 48 for some schemes) to the issuer,
which validates the cryptogram and checks risk parameters.

## 6.5 Cryptogram generation & keys

- The card and issuer share keys derived from the card's **Master Keys**
  (issuer/master card keys, master session keys). Key management uses
  **CBC**/3DES or AES with **session key diversification** (UDK, unique derived
  key per card).
- **ARQC** proves card authenticity per transaction; the **TC** provides
  non-repudiation of an offline-approved transaction.
- The issuer validates cryptograms with an **HSM** (see `09-pci-dss.md` and
  `12-system-design.md`). EMVCo also supports **online cryptogram**
  verification by the network (e.g. Visa/Mastercard HSM-as-a-Service).

## 6.6 Contactless specifics

- **Tap to pay**: cards, phones, watches using NFC; fast tap (<500 ms typical).
- **Visa payWave / Mastercard Contactless / Amex Expresspay / Discover D-PAS /
  JCB / UnionPay** are brand wrappers over EMV contactless kernels.
- Small-value transactions (below the contactless floor limit) can be approved
  **offline with no CVM**; limits are raised for authenticated mobile wallets.
- **Data elements**: same TLVs plus proximity-specific tags (e.g. 9F35
  terminal type, 9F1A country code, 9F40 additional capabilities).
- Contactless supports **fast cryptogram generation**, and may skip some
  processing for low-value taps ("quick tap").

## 6.7 EMV level testing

- **EMV Level 1** — terminal/card reader compliance: mechanical, electrical,
  protocol (transfer of data between card/device and reader).
- **EMV Level 2** — application-level (kernel) compliance.
- **EMV Level 3** — network-side (issuer host) testing (e.g. Visa BIN/Base,
  Mastercard).

## 6.8 Implementing EMV in your platform

**Card side (issuing)**
- Personalization: write PAN, expiry, keys, risk parameters to chip; use a
  **personalization bureau** or an HSM-backed personalization service.
- Cryptogram verification: online ARQC/TC verification via HSM; maintain
  card risk parameters (action codes, floor limits, CVM list).
- Support **script processing** to update cards in the field.
- EMVCo registration: RID is allocated via ISO/EMVCo; AIDs/PIX are scheme-
  approved for brand applications.

**Terminal side (acquiring)**
- Use a certified kernel (e.g. Mastercard M/Chip, Visa VIS, or third-party
  kernels) and certified terminal software (Level 1/2).
- Ensure **DE 55** passthrough and correct **processing code**, POS entry mode
  (DE 22), and condition code (DE 25).
- Handle **fallback** (chip read failure → magstripe) carefully — fallback
  signals risk and may attract fees.

## 6.9 Key references

- EMVCo — "EMV Contact Chip", "EMV Contactless Chip", specifications portal
  (emvco.com/specifications)
- EMV '96 ICC & Application Specifications (public PDFs)
- ACS Technologies — "EMV Specification" training deck
- mstcompany.net — "EMV Transaction Flow" series (PSE, AID, candidate list)
- openscdp.org — EMV tutorial (GPO, PDOL, AIP/AFL)
- Implementing Electronic Card Payment Systems (Artech House) ch. 6
