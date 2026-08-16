# 21. Cross-Border, FX & Dynamic Currency Conversion

## 21.1 Cross-border card transactions

A transaction is **cross-border** when the issuer and acquirer are in different
countries (or the card is used in a different currency than the issuing
country). Cross-border adds:

- **FX conversion** — the cardholder's billing currency differs from the
  transaction currency.
- **Scheme cross-border fees** — additional interchange/assessments for
  cross-border processing.
- **Regulatory differences** — data protection, sanctions, consumer protection,
  and scheme licensing per jurisdiction (see `15-payment-system-governance-regulation.md`).
- **Routing & resilience** — international processing requires robust network
  links, and cross-currency settlement between members.

## 21.2 Where FX happens

FX conversion can occur at one of three places:

| Point of conversion | Who sets the rate | Result |
|---------------------|-------------------|--------|
| **Issuer (default)** | Issuer's foreign-currency rate + markup | Cardholder billed in billing currency at issuer rate |
| **Acquirer / POS (DCC)** | Acquirer or merchant rate + markup | Cardholder sees the amount in their home currency at the terminal |
| **Network** | Scheme or settlement rate | Used when clearing/settlement crosses currencies between members |

## 21.3 Dynamic Currency Conversion (DCC) / POI currency conversion

Per the **Mastercard DCC Guide (2025, Merchant Version)**:

- **DCC = point-of-interaction (POI) currency conversion**: at the terminal or
  e-commerce checkout, the cardholder is offered their **own (billing)
  currency** instead of the local currency.
- **Cardholder choice is mandatory** — the cardholder must choose between
  paying in the local currency or in their billing currency; the choice is
  communicated clearly and the exchange rate displayed.
- **Issuer rate vs acquirer rate** — under DCC the acquirer/merchant provides
  the rate (and earns the margin); without DCC the issuer provides the rate.
- **Restrictions** — do **not** perform DCC on **prepaid travel /
  multi-currency cards**, because converting again at the issuer would cause a
  **double conversion**. The DCC provider must be able to identify such cards
  (via BIN ranges) and skip the offer.
- **Full disclosure** — the cardholder must see the amount in both currencies,
  the conversion rate, and any fees before agreeing.

## 21.4 Cross-currency settlement between members

When acquirer and issuer settle in different currencies, the scheme needs a
mechanism to move value across currencies:

- **Correspondent banking** — legacy model; slow and costly.
- **Centralized currency trading (CCT)** — a central exchange hub stores
  cross-currency PSP rates and returns the best available rate to the
  originator PSP (ECB/TIPS multi-currency model).
- **Multi-currency CSM** — a platform that debits an account in one currency
  and credits an account in another (end-to-end central-bank settlement), or
  settles the final leg in the beneficiary's currency.
- **FX risk** — the party acting as the conversion intermediary bears FX risk
  during settlement; manage with limits and rate governance.

## 21.5 Building cross-border into your network

- Define **cross-border fees** and interchange tiers per corridor.
- Support **multiple currencies** in authorization, clearing, and settlement;
  store transaction currency + billing currency.
- Publish an **FX rate policy**: issuer-default conversion, DCC at the POI,
  and settlement conversion rules.
- Enforce **DCC compliance** in the rulebook: mandatory cardholder choice,
  rate disclosure, and prepaid/multi-currency card exclusions.
- Monitor **corridors** for fraud and interchange-arbitrage (see
  `10-fraud-risk.md`).

## 21.6 Key references

- Mastercard — *DCC Guide* (Merchant Version, 2025).
- Visa — *International Operating Regulations* (currency conversion rules).
- BIS — *Cross-currency settlement of instant payments in a multi-currency
  CSM* (ECB/Banca d'Italia/Riksbank paper).
- World Bank / CPMI — cross-border payments cost and speed initiatives.
