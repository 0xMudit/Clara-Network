# 1. Card Payment Ecosystem & the Four-Party Model

## 1.1 What a card payment system is

A card payment system is a network of banks, merchants, and a **card scheme**
(also called a *network* or *brand*) that agree to common rules, message
formats, and fee structures so that a card issued by one bank can be accepted
anywhere the network is accepted.

> "Credit card networks provide the infrastructure for customers to make
> purchases from merchants. They allow financial institutions to talk to each
> other, approve transactions, and ultimately hand over the money."
> — Investopedia

The major international card schemes are **Visa** (largest by transaction
share), **Mastercard**, **American Express**, **Discover**, plus regional
networks such as **UnionPay** (China) and **Mir** (Russia). As of the most
recent Nilson data cited in public sources, Visa handled ~39% of global card
transactions, UnionPay ~33%, and Mastercard ~25%.

## 1.2 The four-party model

The dominant structure for Visa and Mastercard is the **four-party (open loop)
model**:

```
  Cardholder (consumer)
        │  (card agreement)
        ▼
   Issuing Bank ─────────────┐   Card Scheme / Network
        ▲                     │   (Visa / Mastercard)
        │                     │   routes messages & sets rules
   [funds move]               │
        │                     │
   Acquiring Bank ────────────┘
        │  (merchant account)
        ▼
      Merchant
```

| Party | Who they are | What they do |
|-------|--------------|--------------|
| **Cardholder** | The person/company holding the card | Initiates payment at a merchant |
| **Issuer (issuing bank)** | The cardholder's bank; issues the card | Approves/declines each authorization, funds the cardholder, bears fraud risk on authenticated transactions, manages the account and billing |
| **Merchant** | The business selling goods/services | Accepts card payments |
| **Acquirer (acquiring bank)** | The merchant's bank | Contracts merchants, submits transactions to the network, settles funds to the merchant, holds the merchant account |
| **Card scheme / network** | Visa, Mastercard, etc. | Routes messages between acquirer and issuer, writes the rulebook, sets interchange and scheme fees, referees disputes |
| **Processor** | Technology provider (e.g. Worldpay/FIS, Stripe) | Handles the technical transport of authorization/clearing/settlement messages between the acquirer and the network — and/or on the issuing side |
| **Payment gateway** | Front-end software | Captures payment data from the checkout/POS, encrypts it, forwards to the processor/acquirer |
| **Sponsor bank** | A licensed bank that sponsors a non-bank | Lets fintechs/processors access the card networks without being a bank; liable for the sponsored entity's regulatory and scheme compliance |

The issuer and acquirer are normally **not** Visa/Mastercard themselves.
Visa/Mastercard are networks — they do not hold consumer balances or credit
risk. American Express and Discover run **closed (three-party) networks** where
the network itself acts as issuer and (often) acquirer.

## 1.3 Open vs closed networks

- **Open networks** (Visa, Mastercard, UnionPay): any member bank can issue
  cards or acquire merchants under the network's rules. This creates the
  **interchange** fee paid from the acquiring side to the issuing side.
- **Closed networks** (Amex, Discover historically, private-label/store cards):
  the network both issues the card and (directly or via partners) acquires the
  merchant. No separate issuing/acquiring banks, no interchange between banks;
  the merchant pays the network directly.

## 1.4 How a transaction moves through the ecosystem

A simplified card-present (in-person) transaction:

1. Cardholder presents card at the **POS terminal** (insert chip, tap, or swipe)
   or enters details at **online checkout**.
2. The merchant's **payment gateway/terminal** sends the authorization request
   to the **acquirer** (or its processor).
3. The **acquirer** submits the request to the **card scheme**.
4. The **scheme** routes the request to the correct **issuer** (identified by
   the BIN/IIN in the PAN).
5. The **issuer** checks account status, available funds/credit limit, and
   fraud signals, then **approves or declines** (a typical response in 1-2
   seconds).
6. The decision flows back through the scheme → acquirer → merchant terminal.
7. Later, **clearing** exchanges final transaction details between acquirer and
   issuer; **settlement** moves funds: issuer → scheme → acquirer → merchant
   account (usually T+1 to T+3 business days).

### 1.4.1 The three processing stages

1. **Authorization** - real-time approve/decline decision on funds and card
   validity. Also covers cancellations (reversals) and advisements.
2. **Clearing** - exchange of final financial transaction details between
   acquirer and issuer so each can post the transaction to its books and
   reconcile. Includes fee calculation (interchange).
3. **Settlement** - the actual exchange of funds. The scheme calculates each
   member's **net position** and instructs transfers between settlement
   accounts (often via central bank RTGS rails).

> Mastercard implements this with its **Authorization Platform**, **Global
> Clearing Management System (GCMS)**, and **Settlement Account Management
> (SAM)**. Visa implements it with **VisaNet**, including the **BASE II**
> clearing system.

## 1.5 The fee stack

A merchant does not pay one fee; it pays a **merchant discount rate (MDR)** —
typically 1%–3% — which is split as:

1. **Interchange fee** - paid to the **issuer** (the largest slice, ~70% of
   total merchant card fees on average). Set by the network, varies by card
   type, transaction environment (card-present vs card-not-present), merchant
   category code (MCC), and risk. Example: 1.5%–3.5% for U.S. credit; regulated
   caps exist in some jurisdictions (e.g. EU caps, U.S. Durbin debit caps).
2. **Assessment / scheme fee** - paid to the **network** (roughly 0.1%–0.2%).
3. **Processor / acquirer margin** - the remainder, negotiated between merchant
   and processor/acquirer.

Money path in dollars (illustrative, from public sources):

```
Sale:      $100.00
Interchange (1.23% + $0.10):  -$1.33   → issuer
Scheme fee (0.15% + $0.10):   -$0.25   → network
Acquirer margin:              -$0.22   → acquirer
Merchant receives:             $98.20
```

## 1.6 Chargebacks and disputes

In a dispute, roles reverse: the cardholder asks the **issuer** to reverse a
charge. The issuer initiates a **chargeback** through the network to the
**acquirer**, which recovers funds from the merchant (or the merchant defends
via **representment**). Chargeback monitoring programs impose thresholds on
merchants; high chargeback ratios can lead to scheme fines or termination.

## 1.7 Regulatory environment

Card systems are heavily regulated. Relevant frameworks include:

- **PCI DSS** (PCI Security Standards Council) - security requirements for
  everyone handling card data (see `09-pci-dss.md`).
- **Interchange regulation** - e.g. U.S. Regulation II (Durbin) debit caps;
  EU interchange caps under the Interchange Fee Regulation; Australia's RBA
  caps.
- **Strong Customer Authentication (SCA)** - EU PSD2, implemented via EMV 3DS
  (see `08-3ds.md`).
- **Scheme rules** - the network rulebooks (Visa Core Rules, Mastercard Rules)
  that all members contractually agree to.
- **AML/KYC** - anti-money-laundering and customer due diligence obligations.

## 1.8 Key references

- Stripe - "Credit card networks explained", "Acquirer vs issuer"
- Mastercard - "Switching explained" (authorization, clearing, settlement)
- FIS/Worldpay - "Life cycle of a transaction"
- TSG Payments - "Payments 101: Credit Card Transaction Flow"
- Investopedia - "Credit Card Networks"
- RBA - "Backgrounder on Interchange and Scheme Fees"
- U.S. Federal Reserve - "Interchange Fees and Payment Card Networks"
