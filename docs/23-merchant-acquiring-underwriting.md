# 23. Merchant Acquiring, Underwriting & Boarding

## 23.1 The acquiring stack

- **Acquiring bank (acquirer)** — a scheme member that sponsors merchants,
  holds merchant accounts, settles funds to merchants, and bears ultimate
  liability to the scheme for its merchants.
- **Agent banks** — community banks that contract merchants on the acquirer's
  behalf; **with liability** (indemnify) vs **without liability** (referral
  only — acquirer underwrites and takes risk).
- **ISOs / MSPs (third-party organizations)** — non-members that solicit and
  sign merchants, provide processing, fraud monitoring, terminals, and
  support. They must be **registered with the scheme** (initial + annual fees),
  and the acquirer remains liable for their actions.
- **Processors** — the technical engine moving authorization/clearing messages
  between the acquirer and the network (often absorbed into the acquirer).
- **PayFac / PSP / MoR** — packaging layers that aggregate sub-merchants or
  act as merchant of record (see `01-card-payment-ecosystem.md`).

## 23.2 Merchant boarding

"Merchant boarding" is the onboarding of a new merchant by the acquirer or a
third party (ISO/MSP). The acquirer is **responsible for due diligence no
matter who performs the boarding**. Screening typically includes:

- Physical inspection / verification of business location.
- Credit history and background check on principals.
- Business plan review: projected sales volume, chargeback activity, sales
  type (card-present vs card-not-present).
- For online merchants: website content, functionality, and return policy.
- **Negative-list screening**: the merchant must **not** appear on
  - **MATCH** (Member Alert to Control High-Risk Merchants — merchants
    terminated for cause or with multiple simultaneous applications),
  - **OFAC SDN** (U.S. Specially Designated Nationals), and other sanctions
    lists.
- Signed merchant application, processing agreement, corporate resolution
  (if applicable), credit reports, and financial statements.

## 23.3 Underwriting (credit + fraud risk)

Acquiring extends effective credit to merchants because the acquirer must pay
issuers for chargebacks if the merchant fails.

- **Risk-based tiers** — segment merchants into low/medium/high risk (per Visa
  Acceptance Risk Standards: small merchants, enterprise, high-integrity-risk
  categories, pay-by-link, Future Sales).
- **Creditworthiness** — verify the business model aligns with the acquirer's
  risk tolerance; review financials/credit bureaus.
- **MCC assignment** — assign the Merchant Category Code(s) that most
  accurately describe the business; assign 2+ MCCs when separate lines of
  business exist at one outlet (each with its own agreement). Accurate MCCs
  drive interchange, monitoring, and compliance.
- **Enhanced due diligence** — for high-risk MCCs (gaming, gambling,
  adult content, money services, virtual currency, high-ticket electronics).
- **Decision** — approve, decline, or conditional approval with exposure
  mitigation: **reserves** (hold % of receipts), **delays in funding**,
  **transaction limits**, or **guarantees/chargeback insurance**.
- **Auto-boarding** — automated underwriting via models/AI is allowed but
  subject to model-risk management and regulatory requirements; significant
  risk events must trigger review of whether onboarding flaws contributed.

## 23.4 Monitoring & portfolio management

Post-boarding controls (per FFIEC / FDIC / OCC guidance):

- **Chargeback monitoring** — chargeback level and frequency per merchant;
  chargeback ratio programs (see `20-disputes-chargeback-management.md`).
- **Fraud monitoring** — average ticket size, % keyed transactions, velocity,
  same-dollar batches, odd even-amount patterns, rising decline/refusal rates,
  zero-balance DDA accounts (money-laundering red flags).
- **Risk-triggered actions** — delay funding, install front-end fraud
  monitoring, request bank statements, visit the business, or place reserves.
- **Termination & MATCH reporting** — terminate high-risk merchants for cause
  and report to the industry MATCH database.
- **PCI DSS compliance validation** — acquirers must ensure merchants and
  TPAs validate PCI DSS (QSAs for large merchants/TPAs, SAQs for smaller
  merchants) per scheme rules (see `09-pci-dss.md`).

## 23.5 Acquirer obligations to the scheme

Per **Visa Acceptance Risk Standards** (mandatory controls):

1. Define a **risk appetite/tolerance** and approval procedures by activity
   segment (permissible, conditionally restricted, prohibited) and by country.
2. Maintain **risk policies**: underwriting, monitoring, termination,
   incident investigation, settlement management, complaint handling, data
   security/retention, business continuity, disaster recovery.
3. **Merchant agreements** — written contracts specifying rights, duties, and
   obligations (fees, settlement, MCC, compliance, notification of ownership
   changes).
4. **TPA agreements** — written terms of use between acquirer and third-party
   agents/processors, including compliance and data-security clauses.
5. **Registration of ISOs/MSPs** with the scheme and ongoing reporting.

## 23.6 Key references

- Visa — *Acceptance Risk Standards* (AACQ / ATPA controls).
- FDIC — *Supervisory Guidance* (Merchant Processing, ch.19, 2024).
- OCC — *Comptroller's Handbook: Merchant Processing*.
- FFIEC — *IT Examination Handbook* (Retail Payment Systems, Merchant
  Acquiring).
- Visa/Mastercard — MATCH and high-risk merchant registration requirements.
