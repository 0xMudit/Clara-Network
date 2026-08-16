# 12. System Design: Building the Card Payment Platform

## 12.1 The core truth: correctness, not throughput

Card payment volume is large but not "web-scale": even a big processor handles
thousands of transactions per second, not millions. **The hard problem is
correctness under failure**, not raw RPS.

> "A single duplicated charge or missed transaction can cause legal liability,
> chargebacks, and user trust loss." — Md Sanwar Hossain

Non-negotiables:
- **Exactly-once execution**: one payment intent = exactly one charge, no
  matter how many retries. (Implemented as **at-least-once delivery +
  idempotent processing + reconciliation**.)
- **Consistent state**: ledger must always reflect true money movement.
- **Auditability**: every state transition is traceable (who, what, when, why).
- **CAP**: tilt decisively toward **consistency** over availability.

## 12.2 Reference architecture

```
                ┌──────────────────────────────────────────────────┐
                │                    API Gateway                    │
                │  authN (mTLS/OAuth), rate limiting, WAF           │
                └───────────────────────┬──────────────────────────┘
                                        │ Idempotency-Key (Redis + DB)
                 ┌──────────────────────▼──────────────────────────┐
                 │                Payment Service                   │
                 │  charge/refund/cancel orchestration, state machine│
                 └───────────┬──────────────────────┬───────────────┘
                             │                      │
        ┌────────────────────▼──────┐   ┌───────────▼─────────────┐
        │  Fraud / Risk Engine      │   │  PSP Adapter / Switch   │
        │  rules + ML (<100ms)      │   │  ISO 8583 / gateway API │
        └───────────────────────────┘   └───────────┬─────────────┘
                                                    │
        ┌───────────────────────────────┐   ┌───────▼──────────────┐
        │   Ledger (append-only)        │   │  Clearing &          │
        │   double-entry                │   │  Settlement Service  │
        └───────────────┬───────────────┘   │  batch files, netting│
                        │                   └───────┬──────────────┘
              ┌─────────▼─────────┐                 │
              │ Outbox → Event    │   ┌─────────────▼─────────────┐
              │ bus (Kafka)       │──▶│  Reconciliation Service   │
              └───────────────────┘   │  match reports vs ledger  │
                                       └───────────────────────────┘
```

Services (from public architecture write-ups):

| Component | Responsibility | Typical implementation |
|-----------|----------------|------------------------|
| Payment API / Gateway | Idempotent REST surface | Idempotency-Key header; validations |
| Charge orchestrator | Auth/capture/refund sagas | State machine; workflow engine |
| PSP Adapter / Switch | Card network communication | ISO 8583 stack, or PSP SDK; retries, circuit breaker |
| Fraud Engine | Real-time risk scoring | Rules + ML, feature store |
| Ledger | Double-entry bookkeeping | PostgreSQL append-only entries |
| Clearing & Settlement | Batch files, interchange, netting | File processors, RTGS instructions |
| Reconciliation | Match internal vs external | Daily batch + anomaly detection |
| Webhook Handler | Idempotent async event sink | Signed receiver + dedupe + DLQ |

## 12.3 Idempotency — the foundation

Network failures are guaranteed: clients time out and retry. Without
idempotency, every retry is a potential duplicate charge.

**Contract** (Stripe-style):
- Client generates a UUID per payment intent; sends `Idempotency-Key` header.
- Same key + same body → replay cached response.
- Same key + different body → `409 conflict`.
- Same key while in flight → `409 request_in_progress` (retry with backoff).

**Implementation**:
- Persist the key **before** doing the work (upsert), with a **unique
  constraint** in the DB (the durable guarantee) and Redis as a fast cache
  (SETNX as a lock) — key window 24 h; store response body, not just status.
- Scope keys to customer + operation (`charge:{customer_id}:{uuid}`).
- Make **webhook processing** idempotent on event ID too (PSPs deliver
  at-least-once).
- Pass your internal intent ID as the **provider-level idempotency key**.

## 12.4 Payment state machine

A payment is a small lifecycle, not a status field. Model explicitly with
allowed transitions:

```
created → requires_action (3DS) → processing → succeeded
                                              → failed
       → processing → succeeded
                    → failed
       → succeeded → refunded / partially_refunded
       → cancelled (before capture)
auth (hold) ──capture──▶ cleared ──settled ──funded
```

Rules:
- Only `created/pending` intents can be charged; once succeeded, create a new
  refund intent — **never modify the original**.
- Store the mutable `payments` current-state row (optimized for queries) plus an
  **append-only payment_events** log (audit/debug).
- Apply transitions with `SELECT ... FOR UPDATE` or optimistic locking so
  concurrent webhooks can't double-apply.

## 12.5 Double-entry ledger

Every money movement is recorded as two entries (debit + credit) summing to
zero, in a single ACID transaction.

- **Append-only**: never UPDATE/DELETE. Corrections are reversing entries.
- **Amounts** in integer minor units (cents) — never floats.
- **Balance** = SUM(credits) − SUM(debits), recomputed from the ledger;
  materialize for read latency, but the ledger is truth.
- Schema: `(entry_id, txn_id, account_id, direction, amount, currency, ts)`,
  indexed by `(account_id, ts)`, clustered by `txn_id` for atomic two-row
  writes.
- The balance invariant check inside the transaction is cheap insurance;
  roll back if unbalanced.
- Purpose-built ledgers (TigerBeetle) when you exceed ~50k writes/sec or need
  active-active cross-region.

## 12.6 Sagas, not 2PC

Authorization takes 1–2 s; you can't hold locks across participants. Use
**orchestrated sagas** with explicit compensating transactions:

| Step | Forward action | Compensation |
|------|----------------|--------------|
| 1 | Authorize card | Void authorization |
| 2 | Score fraud | Release risk hold |
| 3 | Capture funds | Issue refund |
| 4 | Post to ledger | Post reversal entries |
| 5 | Notify | Send cancellation |

- Drive sagas from a durable workflow engine (Temporal, Step Functions, or a
  state machine) that survives host crashes.
- Compensations are retried until success.

## 12.7 Transactional outbox

To publish events without dual-write inconsistency: write the domain event to
an **outbox table in the same transaction** as the business change; a relay
process publishes to Kafka. Event published iff transaction committed.

## 12.8 PSP / switch integration robustness

- **Timeouts**: 2–5 s for charges, ~10 s for refunds; never wait forever.
- **Retries** with the same idempotency key, max ~3, exponential backoff +
  jitter.
- **Unknown state**: on timeout with unknown result, set `unknown`; a
  reconciliation job resolves within minutes by polling the PSP.
- **Circuit breaker**: trip after error-rate threshold; fail fast.
- **Multi-provider failover / smart routing**: rule engine (+ML for
  cost/auth-rate) picks the processor per transaction; keep reconciliation
  per-provider.

## 12.9 Reconciliation

Every morning (or per batch), compare the ledger against provider/bank/scheme
settlement files:

1. **Charges**: every charge appears in the provider's records, amounts match.
2. **Settlement**: bank statement total = sum(captures) − fees.
3. **Refunds**: every refund matches.
4. **Auth/capture/settle bridge**: authorization holds → cleared → settled →
   funded.

Mismatch classification: missing transaction, amount mismatch, orphan (in
report but not in ledger), duplicate. Alert on any unresolved discrepancy.

## 12.10 Data model & storage

- **payments**: id, merchant_id, customer_id, idempotency_key (unique), amount,
  currency, status, payment_method_id (token, not PAN), provider_id, created_at.
- **payment_events**: id, payment_id, type, payload, created_at (append-only).
- **ledger_entries**: as above.
- **customers / payment_methods**: store provider tokens, not card data.
- **settlements / reconciliation**: batches, net positions, match status.

Scaling:
- Shard payments/ledger by `merchant_id` or `customer_id` (single-shard
  locality for a payment).
- Partition ledger by month; archive old partitions to cold storage.
- Read replicas for reporting; never hit the write primary for analytics.
- Redis for idempotency cache + rate limiting; PostgreSQL for durable truth.

## 12.11 Security & compliance wiring

- **HSM**: PIN blocks, cryptograms, CVV/iCVV, MAC, key storage (see
  `09-pci-dss.md`).
- **Edge tokenization**: keep PAN out of your environment → SAQ A scope.
- **Secrets management**: HSM-backed KMS, no keys in code.
- **Audit logs** to SIEM; sign webhooks; dedupe + DLQ.
- **Rate limiting / anti-enumeration** at the gateway.

## 12.12 Performance targets (typical, from public case studies)

| Metric | Target |
|--------|--------|
| Peak TPS | ~5,000 (Black Friday ≈ 10× daily mean) |
| p99 charge latency | < 1 s end-to-end (excluding 3DS step-up) |
| Fraud decision latency | < 100 ms |
| Idempotency store | ~24 h TTL, Redis + Postgres |
| Ledger growth | ~1 TB/yr @ 100M txns/yr (4–8 entries/txn) |

## 12.13 Build vs buy decision

- **Path A — build own rails**: ISO 8583 switch, BIN sponsorship, sponsor bank,
  HSM, scheme certification. Best when volume is high enough for interchange
  savings to dominate (>$1B/yr), unique flows, or data-residency mandates.
- **Path B — third-party PSP**: tokenization/fraud/acquiring live in the vendor;
  you build an idempotent API, ledger, reconciliation, and webhook consumer.
- **Path C — direct acquirer via sponsor**: selective; keep PSP for tokenization
  and fraud, add smart routing later.

## 12.14 Production readiness checklist

- [ ] Every write endpoint checks idempotency (Redis + unique constraint).
- [ ] Payment state machine enforces valid transitions.
- [ ] Transactional outbox for all event publishing.
- [ ] PSP calls have timeout + retry + circuit breaker; provider-level idempotency
      keys.
- [ ] Webhooks verify signatures, process idempotently, DLQ on failure.
- [ ] Ledger append-only, integer minor units, balanced on every write.
- [ ] Daily reconciliation with alerting on mismatches.
- [ ] No raw PAN/CVV anywhere in your systems.
- [ ] Fraud scoring in <100 ms and non-blocking for the hot path.
- [ ] Load-tested to 2× peak TPS including concurrent retry storms.

## 12.15 Key references

- Stripe/payment-engineering articles on idempotency, ledger, outbox, sagas
- "Design a Payment System" — systemdesign.one newsletter; HLD Handbook case
  study; singhajit.com; sujeet.pro; letsbuildsolutions.com; mdsanwarhossain.me
- Visa/Mastercard network documentation for clearing/settlement files
- EMVCo / PCI SSC specs referenced in `04`, `06`, `07`, `09`
