# 18. Settlement, Liquidity & Guarantee Infrastructure

## 18.1 Where settlement actually happens

Card networks clear transactions but **final settlement happens in the books
of a central bank** (or, rarely, a commercial settlement bank). For a scheme
design this is the single most important infrastructure decision:

- **Gross vs net** — settle each instruction individually (RTGS-style) or
  accumulate net positions settled at fixed times (deferred net settlement).
- **Real-time vs deferred** — card clearing is typically batched daily into a
  net settlement cycle; instant payments (see `24-instant-payments-rtp.md`)
  settle gross and in real time.

## 18.2 Settlement models (World Bank FPS study)

Three institutional setups are common; each can use gross or net settlement:

1. **Hub approach** — the scheme operator acts as clearinghouse, "mirrors"
   participant central-bank balances in its own books, clears in real time,
   and sends a single settlement instruction to the central bank.
2. **RTGS-based** — the central bank RTGS itself carries the clearing and
   settlement of the transactions.
3. **Distributed clearing** — participants validate/clear bilaterally; the
   payer's PSP then sends a settlement instruction to the central bank.

For a card scheme the **hub approach** is the norm: the network computes net
positions between members and instructs the central bank (or its settlement
bank) to move funds.

## 18.3 Account structures (Bank of England RTGS access policy)

Central banks offer different account types to payment systems:

- **Omnibus accounts** — a single account for the payment system operator
  (PSO) holding funds on behalf of its members. The PSO manages sub-positions
  in its own books; settlement finality is governed by the system's rules.
- **Prefunding accounts** — each participant pre-funds an amount up to a cap
  covering its maximum possible net obligation. Because obligations are capped
  by funded balances, **settlement risk is eliminated even if a participant
  defaults mid-cycle**.
- **Direct RTGS accounts** — participants settle their own net positions
  directly in central bank money.

## 18.4 Liquidity management

- **Intraday liquidity** — participants need funds intraday for net positions
  and for instant payment obligations.
- **Prefunding caps** — set each participant's maximum net debit position
  equal to its prefunded balance; reject or queue instructions that would
  exceed it.
- **Liquidity transfers** — allow moving funds between the scheme account and
  an RTGS account (crucial for 24/7 instant systems that settle when RTGS is
  closed, see `24-instant-payments-rtp.md`).
- **Queues & priorities** — for net systems, queue high-value vs low-value
  items to avoid gridlock.

## 18.5 Managing default & guarantee (PFMI principles 4, 9, 13)

If a member cannot pay its net obligation, the scheme must have a funded
response:

- **Default fund / guarantee fund** — pre-funded by all members to cover the
  default of the largest participant(s).
- **Loss allocation** — waterfall: defaulting member's prefunded balance →
  default fund → surviving members pro-rata → (rarely) scheme capital.
- **Net debit caps** — limit any single participant's exposure.
- **Settlement finality rules** — define exactly when obligations become final
  and irrevocable (PFMI Principle 8).

## 18.6 Reference design for a card scheme settlement layer

1. **Clearing** — network captures/aggregates authorizations and clearing
   (financial) messages; computes per-member net positions per currency.
2. **Settlement instruction** — network sends net positions to the central
   bank / settlement agent (ISO 20022 `pacs` messages, see `05-iso20022.md`).
3. **Funding** — members pre-fund (prefunding model) or hold RTGS accounts.
4. **Reconciliation** — network reconciles its own records to central-bank
   statements and notifies members (see `12-system-design.md` ledger
   reconciliation).
5. **Guarantee** — default fund sized by stress scenarios (largest participant
   default at peak volume).

## 18.7 Key references

- BIS CPMI — *Principles for Financial Market Infrastructures* (Principles 4,
  8, 9, 13, 20).
- Bank of England — RTGS *access policy* (omnibus and prefunding accounts).
- World Bank — *Settlement models in fast payment systems and implications
  for participant access*.
- Federal Reserve — *Policy on Payment System Risk* (daylight overdraft,
  net debit caps).
- ISO 20022 — pacs message types for clearing & settlement.
