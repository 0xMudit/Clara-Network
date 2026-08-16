# 3. Payment Flows: Authorization, Clearing & Settlement

Card payments happen in **three distinct phases** with very different timing and
failure modes. Engineering a payment platform means modeling each phase
explicitly.

## 3.1 Phase 1 — Authorization (real-time, milliseconds)

Authorization is the **approve/decline** decision on funds, card validity, and
fraud risk, exchanged in real time.

### 3.1.1 Happy path (card present)

```
Cardholder
   │ taps / inserts / swipes
   ▼
POS Terminal ──(PAN, amount, POS data)──▶ Gateway/Processor
   ▲                                        │
   │        approve / decline               ▼
   │◀──────────────────────────────────  Acquirer
   │                                        │
   │                                        ▼
   │                                     Card Scheme
   │                                        │  (route by BIN)
   │                                        ▼
   │                                     Issuer
   │                                 (checks funds, limit,
   │                                  fraud, CVV, cryptogram,
   │                                  3DS result)
   │
   │  Response flows back in reverse order
```

Key points:

- **Message types** (ISO 8583): `0100` authorization request, `0110` response,
  plus reversals `0400/0410` and advisements `0120/0121`.
- The issuer can **approve, decline, or refer** the transaction; a response
  code (e.g. `00` approved, `05` do not honor, `51` insufficient funds) is
  returned.
- **Authorizations are not money movement.** The issuer typically places a
  **hold/reservation** on the cardholder's available balance. Holds expire
  (e.g. 7-30 days) if the transaction is not captured and cleared.
- Online card systems typically complete in **1-2 seconds** end-to-end.
- Offline (EMV) authorization may be approved locally by terminal/card under
  floor limits (see `06-emv-chip.md`).

### 3.1.2 Timeouts, retries, reversals

- If the response is lost (timeout), the terminal/processor must retry with the
  **same systems trace audit number (STAN, DE 11)** and message id so the issuer
  can dedupe — **idempotency**.
- **Reversals** (`0400`) reverse a previously approved authorization (e.g.
  terminal timed out but the auth succeeded). These are not refunds.
- Authorization can be **voided** (full reversal of an uncaptured auth).

### 3.1.3 Hold / pre-authorization & capture (two-step)

Common in hotels, car rental, gas pumps, e-commerce:

1. **Auth-only** - validate funds, place a hold; no money moves.
2. **Capture** later - when goods ship / stay completes, capture the final
   amount (same, less, or more than the auth via **incremental auth**).

Captured amounts must match within scheme tolerances; mismatches can trigger
fees or declines.

## 3.2 Phase 2 — Clearing (near-real-time to next-day)

Clearing is the **exchange of final transaction details** between acquirer and
issuer so both sides can post transactions to their books, calculate
**interchange**, and reconcile. It is separate from authorization.

### 3.2.1 Batch vs real-time clearing

- **Batch (traditional):** the acquirer collects authorized transactions into a
  **batch file** and submits at end of day (merchant "closes the batch").
- **Real-time clearing:** with instant payments and modern networks, clearing
  messages are sent per-transaction shortly after authorization.

### 3.2.2 Network clearing systems

- **Visa** uses **BASE II** (part of VisaNet): collects, validates, and
  transmits transaction data to receiving institutions for reconciliation and
  settlement. Generates reconciliation reports. Files: acquirer submits
  "Outgoing", issuer receives "Incoming".
- **Mastercard** uses **GCMS** (Global Clearing Management System), based on
  **IPM** (Integrated Product Messages, an ISO 8583 variant). Acquirers submit
  **R111** files; issuers receive **T112** confirmation reports.
- The network **validates** fields, **calculates interchange** and settlement
  totals, handles **currency conversion**, and returns edited files to members.

### 3.2.3 What clearing carries

Per transaction: PAN (or token), amount(s) in transaction/settlement/issuer
currencies, currency conversion rates, fees, MCC, POS data, retrieval
reference number, dates, and the interchange qualification data.

### 3.2.4 Clearing outcomes

- **First presentment** — the acquirer presents the transaction to the issuer.
- **Disputes** — if the cardholder disputes, the issuer can **chargeback**; the
  acquirer can **represent** with evidence; the issuer can **second presentment**;
  arbitration follows (see 3.4).

## 3.3 Phase 3 — Settlement (T+1 to T+3, funds movement)

Settlement is the **actual transfer of funds** between members.

### 3.3.1 Net settlement

The scheme:

1. Aggregates all cleared transactions for the period.
2. Calculates each member's **net position** (gross credits minus gross debits
   plus/minus fees and interchange adjustments).
3. Sends **settlement advisements** to members.
4. Instructs funds transfers through settlement accounts — typically at a
   central bank or a designated settlement bank (**RTGS**, e.g. Fedwire, TARGET2)
   or a commercial settlement bank.

> Mastercard's **SAM** (Settlement Account Management) system calculates net
> positions and performs both advisements and funds transfer. Visa's settlement
> uses VisaNet settlement banks.

### 3.3.2 Gross vs net settlement

- **Net settlement** — members exchange net positions; fewer transfers, standard
  for card schemes.
- **Gross settlement** — each transaction settled individually; used in RTGS
  systems (real-time gross settlement).

### 3.3.3 Merchant funding

The acquirer funds the merchant's account (minus fees) on the settlement
schedule: T+1, T+2, or T+3 depending on scheme, jurisdiction, and merchant
agreement. Modern schemes offer faster/instant funding as a paid feature.

### 3.3.4 Value dates and float

Value dates determine when funds become available. Float = time between the
cardholder's debit and the merchant's credit; issuers earn float on credit
cards.

## 3.4 Disputes, chargebacks and retrievals

Dispute lifecycle (scheme-regulated, with strict timelines):

1. **Retrieval request** — issuer asks the acquirer for transaction evidence.
2. **Chargeback** — issuer reverses the transaction, debiting the acquirer.
3. **Representment** — acquirer re-presents the transaction with evidence.
4. **Pre-arbitration / arbitration** — escalation to the scheme, which rules.

Reasons: unauthorized/fraud, goods not received, services not as described,
credit not processed, etc. Chargeback codes are standardized (Visa Reason
Codes, Mastercard Chargeback Reason Codes). Chargeback monitoring programs
track ratio thresholds (e.g. >0.9% Visa, >0.75% Mastercard for some programs)
and can impose fines.

## 3.5 Refunds

- **Full/partial refunds** — acquirer-initiated credits to the cardholder.
- **Void** — reversal of an uncaptured auth (no clearing).
- Refunds flow through clearing/settlement as credits and reduce the merchant's
  funding.

## 3.6 End-to-end timeline (typical card payment)

| Time | Step |
|------|------|
| T+0 (seconds) | Authorization (hold placed on cardholder funds) |
| T+0 (EOD) | Merchant closes batch; acquirer submits clearing files |
| T+0→T+1 | Network validates, computes interchange, sends clearing reports |
| T+1→T+3 | Settlement: net positions exchanged; merchant funded (minus fees) |
| Statement date | Cardholder billed on issuer statement |

## 3.7 Design implications for your platform

- Model **three separate state machines**: authorization, clearing, settlement —
  they have different timing, idempotency windows, and retry semantics.
- **Auth → clearing amounts must reconcile**; track auth/capture mismatch.
- Clearing files are **batch**, settlement is **net**; reconciliation needs to
  bridge all three (see `12-system-design.md`).
- Preserve the **retrieval reference number (DE 37)** and **STAN (DE 11)**
  end-to-end — they are the keys the scheme uses to match messages.
- Store **original data elements (DE 90)** for reversals/chargebacks.

## 3.8 Key references

- Mastercard — "Switching explained" (Authorization Platform, GCMS, SAM)
- FIS/Worldpay — "Life cycle of a transaction"
- Intelica — "Mastercard & Visa clearing process for issuers & acquirers"
- PXP — "Settlement" glossary
- TSG Payments — "Payments 101: Credit Card Transaction Flow"
- Stripe — "Interchange Fees 101"
