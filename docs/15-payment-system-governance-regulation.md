# 15. Payment System Governance, Regulation & PFMI Compliance

## 15.1 Why a card network is regulated

Running a card scheme is closer to running a **financial market
infrastructure (FMI)** than running a software product. A scheme operator:

- **Owns the rulebook** that binds issuers, acquirers, merchants and
  cardholders.
- **Settles inter-member positions** (directly or via a central bank), which
  makes it a payments system operator.
- May be designated as a **systemically important payment system (SIPS)**,
  triggering formal oversight by the central bank and PFMI assessment.

Before you can operate, you typically need: a **scheme/network license** (or
regulatory designation), **e-money / payment-institution licenses** (PSD2 in
Europe, money transmitter licenses in the U.S.), **members**, and **a central
bank settlement relationship**.

## 15.2 PFMI — the international standard for FMIs

The **CPMI-IOSCO Principles for Financial Market Infrastructures (PFMI)** are
the de-facto global benchmark. Central banks and regulators use them to
supervise payment systems, CSDs, CCPs and trade repositories.

### The 24 principles (highlights)

| # | Principle | What it means for a card scheme |
|---|-----------|--------------------------------|
| 1 | **Legal basis** | The system has a sound, clear, enforceable legal basis (rulebook, contracts, laws). |
| 2 | **Governance** | Clear governance with clear accountability; board oversight; risk committee. |
| 3 | **Framework for comprehensive management of risks** | Risk management across legal, credit, liquidity, operational, and general business risk. |
| 4 | **Credit risk** | Limit and manage credit exposure between participants (incl. settlement banks). |
| 8 | **Settlement finality** | Clear and final settlement at least by end of value date; ideally real-time. |
| 9 | **Money settlements** | Settle in central bank money where practical; otherwise in commercial bank money with controls. |
| 11 | **Central securities depositories** | (N/A for card networks; applies to securities FMIs.) |
| 13 | **Default procedures** | Rules and procedures for participant default, incl. loss/allocation of losses. |
| 17 | **Operational risk** | High availability, scalability, and resilience (see `19-stand-in-processing-availability.md`). |
| 19 | **Tiered participation** | Manage risks arising from tiered (indirect/sponsor) participation. |
| 20 | **FMI links** | Safe and efficient links with other FMIs (e.g. RTGS, ACH, instant systems). |
| 21 | **Efficiency & effectiveness** | Meet needs of participants and markets while being efficient. |
| 22 | **Communication procedures & standards** | Practical, safe, efficient procedures and standards (ISO 8583 / ISO 20022). |

### The 5 Responsibilities of authorities (A–E)

Authorities (central banks, regulators) are responsible for:

- **A. Regulation, supervision & oversight** of FMIs.
- **B. Powers and resources** — adequate powers and resources to carry out A.
- **C. Policy disclosure** — clear, transparent disclosure of policies.
- **D. Application of principles** — consistent application, incl. assessment.
- **E. Cooperation** across authorities and jurisdictions (esp. for
  cross-border systems).

## 15.3 Licensing paths by jurisdiction

| Jurisdiction | License / approval | Notes |
|--------------|--------------------|-------|
| **U.S.** | Money transmitter license (state-by-state, 50 states + RTP/ACH access), or bank charter | Scheme operators get special treatment via FED/Risk-focused supervision; card networks historically operate under the Bank Holding Company Act framework. |
| **EU/EEA** | Payment Institution (PSD2), E-Money Institution (EMD2), or bank license | Needed to provide payment services; interoperability obligations under PSD2. |
| **UK** | FCA payment institution / e-money license | Post-Brexit equivalence + Bank of England RTGS access. |
| **India** | RBI authorization as Payment System Operator | UPI participation; PPI licenses for issuing. |
| **Brazil** | BCB licensing (Pix participation via direct/sponsored access) | Pix is central-bank operated; private IPS need BCB approval. |

## 15.4 Governance design for a scheme operator

- **Board & committees** — board, risk committee, audit committee, technical
  standards committee, disputes/arbitration panel.
- **Rulebook governance** — a controlled change process for scheme rules
  (versioned rulebooks, member consultation windows, published effective
  dates).
- **Membership governance** — admission, suspension, and exit criteria
  (see `16-membership-rulebook-certification.md`).
- **Oversight relationship** — maintain a formal engagement with the central
  bank / regulator; publish a **self-assessment against PFMI** periodically.
- **Transparency** — disclose fees, interchange, operating procedures, and
  risk management practices.

## 15.5 Key references

- BIS CPMI-IOSCO — *Principles for Financial Market Infrastructures* (PFMI),
  CPMI paper d101a (April 2012).
- BIS CPMI — *Guidance on cyber resilience for FMIs* (d146).
- ECB — TIPS regulatory framework (TARGET Instant Payment Settlement).
- Federal Reserve — policy on payment system risk (PSR), intraday credit.
- RBA / MAS / RBI / BCB — payment system oversight statements and PFMI
  assessment reports.
