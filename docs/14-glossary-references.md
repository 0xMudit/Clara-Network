# 14. Glossary & Source References

## 14.1 Glossary

| Term | Definition |
|------|------------|
| **Acquirer** | Bank/entity that contracts merchants, submits transactions to the network, and settles funds to the merchant |
| **ACS (Access Control Server)** | Issuer-side 3DS system that performs authentication |
| **AID** | Application Identifier (RID + PIX) on an EMV card |
| **AIP** | Application Interchange Profile (EMV card capabilities bitmap) |
| **ARQC / TC / AAC** | EMV application cryptograms: Authorization Request, Transaction Certificate, Application Authentication Cryptogram |
| **Authorization** | Real-time approve/decline of funds and card validity |
| **AVS** | Address Verification Service |
| **BIN / IIN** | Bank/Issuer Identification Number — leading 6–8 digits of the PAN |
| **BIN sponsorship** | A bank sub-licensing BIN range to a non-bank issuer |
| **Chargeback** | Reversal of a transaction initiated by the issuer on the cardholder's behalf |
| **Clearing** | Exchange of final transaction details between acquirer and issuer |
| **CNP** | Card-not-present (online/phone/mail) transaction |
| **CP** | Card-present (face-to-face) transaction |
| **CVV / CVC / CVV2** | Card verification value; security code |
| **DE** | Data element (ISO 8583 field) |
| **EMV** | Europay, Mastercard, Visa — chip card standard (EMVCo) |
| **F2F** | Face-to-face (card present) |
| **HSM** | Hardware Security Module — key/PIN/cryptogram security |
| **Interchange** | Fee acquirer pays to issuer via the network |
| **IPM** | Mastercard's ISO 8583-based clearing format (Integrated Product Messages) |
| **Issuer** | Bank that issues the card and holds the cardholder account |
| **Luhn** | Mod-10 checksum algorithm for PANs |
| **MAC** | Message Authentication Code (ISO 8583 DE 64/128) |
| **MCC** | Merchant Category Code |
| **MDR** | Merchant Discount Rate — total fee merchants pay to accept cards |
| **MID / TID** | Merchant ID / Terminal ID |
| **MII** | Major Industry Identifier (first digit of PAN) |
| **MTI** | Message Type Identifier (ISO 8583) |
| **mTLS** | Mutual TLS — client certificate authentication used by scheme APIs |
| **Network (scheme)** | Visa, Mastercard, Amex, Discover, etc. — routes messages, sets rules/fees |
| **PAN** | Primary Account Number |
| **PAR** | Payment Account Reference — links a PAN and its tokens |
| **PCI DSS** | Payment Card Industry Data Security Standard |
| **PDOL** | EMV Processing Data Object List (terminal→card) |
| **POS** | Point of Sale |
| **PSP** | Payment Service Provider |
| **P2PE** | Point-to-Point Encryption |
| **Representment** | Acquirer re-submitting a disputed transaction with evidence |
| **Retrieval reference number** | ISO 8583 DE 37 — transaction matching key |
| **Reversal** | Cancel an authorization (DE 90 original data) |
| **RTGS** | Real-Time Gross Settlement |
| **SAQ / ROC / QSA / AOC** | PCI validation: Self-Assessment Questionnaire / Report on Compliance / Qualified Security Assessor / Attestation of Compliance |
| **SCA** | Strong Customer Authentication (PSD2) |
| **SDA / DDA / CDA** | EMV offline data authentication methods |
| **Settlement** | Actual transfer of funds between members/merchants |
| **STAN** | Systems Trace Audit Number (ISO 8583 DE 11) |
| **Switch** | Network software/hardware routing messages between acquirers and issuers |
| **Token** | Surrogate for PAN (network or acquirer token) |
| **TSP** | Token Service Provider |
| **TVR / TSI** | EMV Terminal Verification Results / Transaction Status Information |
| **3DS** | 3-D Secure — CNP cardholder authentication protocol |
| **UETR** | Universally Unique End-to-end Transaction Reference (ISO 20022) |

## 14.2 Standards bodies & portals

- **ISO** — iso.org: ISO 8583 (card messages), ISO 20022 (payments messaging),
  ISO/IEC 7812 (PAN/IIN), ISO/IEC 7816 (contact cards), ISO/IEC 14443
  (contactless), ISO 9564 (PIN), ISO 18245 (MCC).
  - ISO 8583 free annexes: `standards.iso.org/iso/8583/ed-3/en/`
  - X9 (Accredited Standards Committee) ISO 8583 maintenance: `x9.org/iso-8583-mas`
  - ISO 20022: `iso20022.org`
- **EMVCo** — emvco.com: chip specs (Books 1–4, contactless), payment
  tokenisation, 3-D Secure, TSP/BCID registration.
- **PCI SSC** — pcisecuritystandards.org: PCI DSS v4.0.1, SAQs, P2PE, SAD
  storage rules.
- **BIS CPMI** — bis.org: harmonised ISO 20022 cross-border data requirements
  (d230), PFMI principles.
- **SWIFT** — swift.com: MT↔MX migration (CBPR+).
- **European Payments Council** — europeanpaymentscouncil.eu: SEPA rulebooks &
  ISO 20022 implementation guidelines.

## 14.3 Scheme developer portals

- Mastercard Developers — `developer.mastercard.com`
- Visa Developer — `developer.visa.com`
- Mastercard Gateway API reference — `na.gateway.mastercard.com/api/documentation`
- GitHub: Mastercard issuing-api-reference-app

## 14.4 Primary research sources used in this library

**Ecosystem & flows**
- Stripe — "Credit card networks explained"; "Acquirer vs issuer"
- Mastercard — "Switching explained" (Authorization/GCMS/SAM)
- FIS/Worldpay — "Life cycle of a transaction"
- TSG Payments — "Payments 101: Credit Card Transaction Flow"
- Investopedia — "Credit Card Networks"
- Intelica — "Mastercard & Visa clearing process"

**Numbering**
- ISO/IEC 7812-1:2017; Visa "Eight Digit BIN FAQ"; Wikipedia "Payment card
  number"; Scientific American "What Is the Luhn Algorithm?"; Barclaycard BIN
  rules

**ISO 8583 / ISO 20022**
- ISO 8583:2023 + annexes; Worldpay ISO 8583 Reference Guide; X9 ISO 8583 MAS
- ISO 20022 portal; BIS CPMI d230 + technical annex; J.P. Morgan "Migration to
  ISO 20022"; EPC SEPA IGs; CPG ISO 20022 explainers

**EMV / 3DS / Tokenization**
- EMVCo spec portals & knowledge hub (contact chip, contactless, 3DS,
  tokenisation)
- EMV '96 public specifications (cardspec/applspec PDFs)
- mstcompany.net EMV transaction flow series; openscdp.org EMV tutorial
- Visa — "A Deep Dive into Tokenized Transactions"; VTS docs
- Stripe "3D Secure 2"; Cardinal EMV 3DS docs; 3DSecure.io protocol docs

**PCI DSS**
- PCI SSC v4.0.1 standard (public copy via Middlebury), "In a Nutshell", QRG,
  Summary of Changes
- Venn / CrowdStrike / pcicompliancehub 12-requirements explainers

**Fees & settlement**
- CRS IF11893; RBA "Backgrounder on Interchange and Scheme Fees"; Federal
  Reserve FEDS 2009-23; Checkout.com interchange explainer; Stripe Interchange
  101; Synctera interchange docs; PXP settlement glossary

**System design**
- systemdesign.one "Design a Payment System"; HLD Handbook case study;
  singhajit.com; sujeet.pro; letsbuildsolutions.com; mdsanwarhossain.me
- Visa/Mastercard network operating specs (authorization, clearing files,
  settlement)

## 14.5 Disclaimer

This documentation set is an **educational research summary** compiled from
public sources. It is not legal, compliance, or financial advice. Building or
operating any part of a payment scheme, issuing, acquiring, or processing
business requires licensing, scheme membership or sponsorship, regulatory
approval, and professional legal/compliance review in the relevant
jurisdictions.
