# 24. Instant Payments & Real-Time Processing

## 24.1 What "instant" means

Instant payments (fast payments) are push credit transfers that are
**settled in real time, 24/7/365**, with funds available to the beneficiary
within seconds. They are a complement to (not a replacement for) card rails.

Representative systems: **UPI** (India, launched 2016), **Pix** (Brazil,
2020), **TIPS** (Eurosystem, 2018), **RTP** (U.S., 2017), **FedNow** (U.S.,
2023), **Faster Payments** (UK, 2008).

## 24.2 Settlement models (World Bank FPS study)

Two major settlement models:

| Model | How it works | Risk profile |
|-------|--------------|--------------|
| **Real-time settlement** | Each payment settles immediately, gross, in central-bank money; sending PSP must pre-fund or hold sufficient balance | Finality at settlement; requires 24/7 central-bank (or dedicated) accounts and liquidity management |
| **Deferred settlement** | Payments cleared instantly but net positions settle periodically (e.g., at RTGS windows) | Faster to launch, but introduces settlement risk between clear and settle |

In all studied jurisdictions, **final settlement occurs in central-bank money**
(either in the RTGS or a dedicated settlement account).

### RTP (U.S.) prefunding model

- RTP uses a **fully prefunded joint account** at the Federal Reserve, owned by
  all participating entities.
- Each sending participant must **pre-fund**; non-funding institutions need a
  funding arrangement with another participant.
- Before forwarding a payment, RTP **verifies and reserves settlement capacity**
  for the sender; insufficient prefunded position → the payment is **rejected**
  (overdrafts are not permitted).
- Settlement is complete when RTP records the decrease in the sender's
  position and the increase in the receiver's position.

## 24.3 TIPS design (Eurosystem) — a reference architecture

TIPS (TARGET Instant Payment Settlement) settles SEPA Instant Credit
Transfers (SCT Inst) in central-bank money, 24/7/365, with capacity around
**2,000 payments/second** and **43M+ transactions/day** demonstrated.

**Design principles:**

1. Provides instant **settlement**, not clearing; validation includes format
   compliance, account existence, authorization, and sufficiency of funds.
2. Settles exclusively in **central bank money**; currency-agnostic design
   (non-euro currencies supported, e.g., Swedish krona).
3. Operates **24/7/365**.
4. Uses **ISO 20022** (A2A only): `pacs.008 FIToFICustomerCreditTransfer`
   and `pacs.002 FIToFIPaymentStatusReport`.
5. **20-second timeout** — the scheme must respond to the originator within a
   configurable threshold (default 20 s); timed-out payments are rejected.
6. **Liquidity management** — dedicated cash accounts (DCAs) top up from RTGS
   accounts during RTGS hours; supports 24/7 settlement without RTGS being open.

**Two settlement models:**

- **2-Instructing Parties (2-IP)** — conditional settlement: TIPS reserves the
  amount on the originator's account, forwards the payment to the beneficiary
  PSP for confirmation, then settles in a second phase.
- **Single Instructing Party (SIP)** — immediate settlement without prior
  reservation; a single instructing party validates with both PSPs and submits
  a pre-accepted payment; TIPS settles immediately.

## 24.4 System design considerations

- **Microservices settlement core** — BIS **Project FuSSE** demonstrates a
  modular settlement engine (Kubernetes, Kafka, Redis) with routing-slip
  messaging (the message carries its itinerary; each service prunes its step),
  stateless horizontal scaling, and demonstrated throughput of **10,000 TPS**.
- **Elastic scaling** — cloud-native, containerized, stateless components that
  scale independently (security-intensive ops scale separately).
- **PQC-ready** — embed quantum-resistant cryptography for long-lived systems.
- **Fraud screening** — instant systems have no time for batch fraud checks;
  real-time scoring must fit in the same 20-second budget (see
  `10-fraud-risk.md`).
- **Liquidity transfer windows** — bridge RTGS ↔ instant accounts; define
  funding and auto-refill rules.

## 24.5 Building instant payments alongside a card scheme

- Reuse the **same settlement infrastructure** as card clearing
  (see `18-settlement-liquidity-infrastructure.md`): prefunding caps, default
  fund, and net/gross decisions.
- Decide your **settlement model**: prefunded real-time (RTP-style) vs
  deferred net (Faster-Payments-style) vs central-bank settlement (TIPS-style).
- Implement **ISO 20022** for the instant layer (see `05-iso20022.md`) while
  keeping **ISO 8583** for card authorization.
- Set your **timeout SLA** (20 s industry norm), **finality rules**, and
  **recall/investigation** flows (SCT Inst defines recalls and investigations).
- Plan **cross-currency instant settlement** if multi-currency (see
  `21-cross-border-fx-dcc.md`).

## 24.6 Key references

- ECB — *TIPS User Requirements* (R2025.JUN) and *TIPS principles and
  proposals* background documents.
- BIS Innovation Hub — *Project FuSSE* (Fully Scalable Settlement Engine).
- World Bank — *Settlement models in fast payment systems and implications
  for participant access*.
- The Clearing House — RTP rules & prefunded account model.
- EPC — SEPA Instant Credit Transfer (SCT Inst) scheme rulebook.
