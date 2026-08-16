# 13. Fees, Interchange & Settlement Mechanics

## 13.1 The merchant discount rate (MDR)

Merchants pay their acquirer/PSP a **merchant discount rate** — typically
**1%–3%** of transaction value — which bundles three components:

```
MDR  =  Interchange (to issuer)  +  Scheme/assessment (to network)
        +  Processor/acquirer margin
```

Example split on $100 (from a public PSP explainer):

| Component | Amount | Recipient |
|-----------|--------|-----------|
| Interchange (1.23% + $0.10) | $1.33 | Issuer |
| Scheme fee (0.15% + $0.10) | $0.25 | Network |
| Acquirer margin | $0.22 | Acquirer/PSP |
| **Merchant receives** | **$98.20** | Merchant |

## 13.2 Interchange fees

- **Interchange** is the fee the acquirer (via the network) pays to the
  **issuer**. It is **set by the network**, published in rate tables, and
  adjusted semi-annually (typically April/October).
- It compensates the issuer for cardholder credit losses, funds rewards, fraud
  management, and processing costs. Interchange drives ~70% of total merchant
  card acceptance costs on average.
- **Interchange categories** determine the rate: card type (credit vs debit vs
  prepaid, rewards tiers), transaction environment (card-present vs
  card-not-present), **MCC**, transaction value, and qualification data
  submitted by the acquirer (the **IRD** — interchange rate designator).
- At clearing, the network **validates** the acquirer's claimed interchange
  category and applies the correct rate; mismatches produce
  **interchange-compliance adjustments**.

### 13.2.1 Regulatory caps (examples)

- **U.S. (Durbin / Reg II)**: debit interchange for large issuers capped at
  ~$0.05 + $0.21 per transaction (later adjusted). Credit not capped.
- **EU (IFR)**: credit 0.3%, debit 0.2% caps.
- **Australia (RBA)**: credit capped at a weighted average 0.50% (individual
  ceiling 0.80%); debit capped at a weighted average 8¢ (ceiling 10¢ or 0.20%).

## 13.3 Scheme / assessment fees

- Paid by **both acquirers and issuers** to the network (e.g. Visa, Mastercard)
  for the switching/clearing/settlement service and brand.
- Typical magnitude ~0.1%–0.2% of a transaction, plus per-transaction charges,
  cross-border surcharges, and optional product fees.
- Australian RBA data: net scheme fees ~A$1.8 billion in 2023/24; issuers get
  significant rebates (net scheme fees paid by issuers much lower than by
  acquirers) because networks compete for issuing volume.

## 13.4 Three-party (closed-loop) economics

- **Amex/Discover-style**: the network is issuer+acquirer; merchants pay a
  merchant service charge directly to the network; **no interchange** between
  banks. Merchant fees typically higher (Amex ~2%+), historically because the
  network keeps both interchange and assessment and covers cardholder credit
  risk.

## 13.5 Fee-relevant transaction attributes (qualification)

The acquirer must send the right data at **clearing** for the best interchange
qualification:

- MCC (merchant category)
- Card type & product tier
- CP vs CNP, POS entry mode (DE 22), condition code (DE 25)
- Level 2/3 data (purchase order, line-item detail, tax, ship-to zip) — lowers
  rates for corporate cards
- AVS/CVV results, 3DS authentication status (see `08-3ds.md`)
- Surcharge/visa electronic indicator, authorization vs sale type (auth vs
  auth-and-capture)
- Transaction date/time vs clearing date/time (delayed clearing can
  **downgrade** interchange)

## 13.6 Settlement mechanics in detail

```
Cardholder's account (issuer)  ──debit──▶  Issuer's settlement account
                                                 │
                                                 ▼  net position to network
                                    Card Network (Visa/Mastercard)
                                                 │
                    ┌────────────────────────────┤
                    ▼                            ▼
            Acquirer's settlement    Fees paid: interchange → issuers,
            account (net position)   scheme fees → network
                    │
                    ▼
            Merchant's account (T+1–T+3), minus fees
```

1. **Clearing** determines amounts, interchange, and net settlement totals.
2. The network computes each member's **net position** (all credits minus all
   debits across all members and currencies, with fees and FX).
3. **Settlement instructions** are sent to settlement banks; funds move
   typically via **RTGS** (Fedwire, TARGET2, CHIPS-style) or the scheme's
   designated settlement banks.
4. Mastercard does this via **SAM**; Visa via VisaNet settlement. Members must
   maintain **settlement accounts** with approved settlement banks and
   pre-fund/debit lines (settlement risk management — credit limits, collateral,
   real-time monitoring).
5. Merchant funding on T+1/T+2/T+3, or faster funding products.

## 13.7 Settlement risk & failure handling

- **Settlement risk** (Herstatt risk in FX): a member fails before settlement.
  Schemes enforce: membership capital/guarantee funds, settlement credit
  limits, daily settlement cycles, and recovery processes.
- **Unsuccessful settlement** → scheme funding/guarantee arrangements (e.g.
  Visa's Settlement Service Guarantee, Mastercard's settlement guarantee fund).
- **Payment system risk frameworks** apply to systemically important payment
  systems (PFMI principles from BIS CPMI).

## 13.8 Pricing models (merchant side)

| Model | Description |
|-------|-------------|
| **Interchange-plus (cost-plus)** | Exact interchange + acquirer markup; most transparent |
| **Blended / flat-rate** | Single rate covers interchange+scheme+markup |
| **Tiered** | Qualified/mid-qualified/non-qualified buckets |
| **Subscription / monthly** | Fixed fee with lower per-txn rates |

## 13.9 Fee flows to reconcile in your platform

- Issuer: earns interchange, pays scheme fees, absorbs credit/fraud losses.
- Acquirer: earns merchant margin, pays interchange + scheme fees, manages
  chargeback/reserve exposure.
- Network: earns scheme fees, sets interchange, runs monitoring programs.

Reconciliation must track: transaction amount, interchange line item, scheme
fee line item, acquirer markup, FX/conversion fees, surcharge, and
chargeback/refund adjustments (see `12-system-design.md`).

## 13.10 Key references

- CRS (Congressional Research Service) — "Merchant Discount, Interchange, and
  Other Transaction Fees"
- RBA — "Backgrounder on Interchange and Scheme Fees"
- U.S. Federal Reserve — "Interchange Fees and Payment Card Networks" (FEDS 2009-23)
- Checkout.com — "What are interchange fees?"
- Stripe — "Interchange Fees 101"
- Synctera — "Card Interchange" docs (IRD, qualification, compliance)
- PXP — "Settlement" glossary
