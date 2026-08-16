# 20. Disputes, Chargebacks & Arbitration

## 20.1 The dispute lifecycle

A dispute (chargeback) is the mechanism by which funds are returned from the
acquirer (and ultimately the merchant) to the issuer when a cardholder (or the
issuer) contests a transaction. The scheme defines the process, timeframes,
reason codes, and liability rules.

### Visa Claims Resolution (VCR) lifecycle

1. **Pre-dispute** — cardholder contacts issuer; issuer may request the
   merchant's transaction data before filing (inquiry / retrieval request).
2. **Dispute submission** — issuer files a dispute with the network, using a
   reason code and supporting evidence.
3. **Dispute response** — acquirer (via the merchant) responds with
   **representment** evidence (proof of purchase, AVS/CVV data, 3DS status,
   refunds, etc.).
4. **Pre-arbitration / arbitration** — if the response is rejected, the dispute
   escalates; an arbiter (the scheme or third party) issues a final decision
   with associated fees for the losing party.

### Timing (VCR)

- Standard disputes: **~46 days** from filing to resolution; contentious
  disputes can exceed **100 days**.
- **Expedited process** targets **≤31 days** by compressing response windows
  and streamlining evidence exchange.
- Fraud and authorization disputes follow a streamlined cycle because the
  evidence requirements are simpler.

### Associated Transactions

An "associated transaction" check runs before a dispute is validated: if the
cardholder was already credited (credit/reversal/adjustment) for the disputed
amount, the dispute is **invalid** and rejected. Always reconcile refunds
against disputes first.

## 20.2 Reason codes (high level)

| Scheme | Category | Examples |
|--------|----------|----------|
| Visa | Authorization | Services not provided, not authorized, counterfeit, transaction not recognized |
| Visa | Fraud | Cardholder-not-present fraud, lost/stolen card, counterfeit |
| Visa | Processing errors | Duplicate processing, late presentment, incorrect amount |
| Mastercard | Fraud | EMV liability shift, cardholder account compromised |
| Mastercard | Authorization | Transaction not recognized, recurring billing issues |
| Mastercard | Processing | Data not present, MCC mismatch, incorrect currency |

Reason codes drive **liability**: who carries the loss when a dispute is won
or lost (issuer, acquirer/merchant, or network).

## 20.3 Merchant defense

Acquirers defend disputes on behalf of merchants. Winning evidence typically
includes:

- Transaction receipt / invoice / delivery confirmation.
- Authorization response (approval code, AVS/CVV results, 3DS result).
- Proof of delivery (tracking, signature).
- **Chip/CVV data** — EMV chip transactions shift counterfeit liability to the
  acquirer unless chip data is present.
- Clear **refund history** (to satisfy the associated-transaction check).
- Clear **terms & conditions** (for subscription/recurring disputes).

## 20.4 Chargeback monitoring programs

Schemes penalize members/merchants with excessive chargeback ratios:

- **Visa** — Visa Merchant Purchase Monitoring / VMDP: chargeback ratio
  thresholds trigger monitoring, fees, and possible merchant termination.
- **Mastercard** — Chargeback Monitoring Program (CMP) and Excessive Chargeback
  Merchant program; high ratios lead to fines and loss of acquirer privileges.
- Ratios are typically measured as **chargebacks ÷ transactions** (and
  sometimes dollar value) over rolling periods.

## 20.5 Designing your own dispute framework

If you build a network, you need a rulebook dispute chapter (see
`16-membership-rulebook-certification.md`):

- **Reason code taxonomy** — a documented, versioned list of reason codes.
- **Timeframes** — filed, response, escalation, arbitration deadlines per
  category; track SLA adherence.
- **Evidence requirements** — what proof each side must provide.
- **Dispute engine** — system to accept dispute messages, validate associated
  transactions, route to acquirer/merchant, track deadlines, and compute fees.
- **Liability rules** — who wins/loses per scenario (fraud liability shift,
  processing errors, authorization issues).
- **Arbitration** — a final, fee-bearing escalation path with clear rules.

## 20.6 Key references

- Visa — *Claims Resolution* (VCR) documentation and reason-code list.
- Mastercard — *Chargeback Guide* and dispute/resolution rules.
- Visa — *Visa Merchant Purchase Monitoring Program* (VMDP) thresholds.
- ISO 20022 — message types for disputes (`adm`/`camt` related inquiry
  messages; see `05-iso20022.md`).
