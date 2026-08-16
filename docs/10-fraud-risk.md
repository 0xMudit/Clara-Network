# 10. Fraud Detection & Risk Management

## 10.1 The fraud landscape in card payments

Fraud vectors relevant to a card platform:

- **CNP fraud** — stolen card data used online; biggest and fastest-growing
  category (mitigated by 3DS, see `08-3ds.md`, and tokens, `07-tokenization.md`).
- **Counterfeit/skimming** — cloned magstripe/chip data (mitigated by EMV).
- **Lost/stolen cards** — mitigated by EMV CVM (PIN), velocity checks, network
  alerts.
- **Account takeover** — credential theft at merchants/issuers.
- **Friendly fraud / chargeback abuse** — legitimate cardholders disputing valid
  transactions.
- **Synthetic identity / bust-out fraud** — in issuing (credit risk).
- **Card testing** — attackers test stolen PANs with tiny transactions before
  big spends (detected via velocity + AVS/CVV patterns).
- **Enumeration** — BIN/account number enumeration at your APIs; requires rate
  limiting and anomaly detection.

## 10.2 Risk architecture: three layers

### Layer 1 — Pre-transaction / onboarding
- **KYC/KYB** for cardholders and merchants; sanctions screening (AML).
- Merchant **underwriting** — MCC, history, chargeback ratio, reserve policies.
- Cardholder **underwriting** (credit) for credit products.

### Layer 2 — Real-time transaction scoring (hot path)
Runs on every authorization with a **sub-100 ms budget**:

- **Rule engine** (deterministic): thresholds on amount, country, velocity,
  MCC, BIN range, AVS/CVV mismatch, terminal/POS data, 3DS status.
- **Model scoring** (ML): features like device fingerprint, cardholder
  behavior, merchant behavior, network risk scores (Visa TC40/SAFE, Mastercard
  SPA).
- **Velocity checks**: transaction counts/amounts per card, per merchant, per
  device, per IP, per BIN in rolling windows.
- **Neural/statistical anomaly detection** on cardholder spend patterns.
- **Decision**: approve, decline, step-up (3DS/challenge), refer, block device,
  queue for review.

### Layer 3 — Post-authorization & lifecycle
- **Network alerts**: Visa TC40 (cardholder dispute), SAFE (safe reporting of
  fraud), Mastercard SPA/AMR (alerts on compromised merchants/accounts).
- **Chargeback management** and dispute defense.
- **Account monitoring**: unusual activity, card testing, bust-out detection.
- **Card reissue workflows** (compromise → replace, token refresh).

## 10.3 Key real-time signals

| Signal | Example |
|--------|---------|
| Transaction velocity | >10 txns/card in 10 min |
| Amount anomaly | First txn 5× average basket |
| Geography/time | CP txn in country A, CNP txn 30 min later in country B |
| Device fingerprint | Emulator, new device on old account |
| Card not present + high value | Elevated risk without 3DS |
| AVS/CVV results | `N` mismatch on high-value txn |
| Terminal/POS entry mode | Magstripe fallback on EMV terminal; manual key-entry CNP |
| BIN/MCC behavior | Gaming MCC, gambling, etc. |
| 3DS status | Unauthenticated high-value purchase |
| Token/PAN status | Suspended token used |

## 10.4 Rules engine vs ML

- **Rules**: transparent, auditable, low-latency, easy to tune for regulatory
  and scheme requirements (e.g. always decline on hot-card list). Handles
  deterministic fraud.
- **ML**: catches novel patterns at scale; requires labeled data, feature
  pipelines, monitoring for drift, and human review loops.
- Production systems run **rules first, model second**, with model output
  integrated into a final risk score and override policies.

## 10.5 The "stand-in" problem

When the issuer host is unavailable, networks offer **stand-in processing**
(authorization by the network using stored rules and risk limits) so merchants
aren't left hanging. If you build issuing infrastructure, decide your
**stand-in policy**: approve within limits, decline, or route.

## 10.6 Implementing fraud in your platform

**Component design (see `12-system-design.md`):**

- **Risk API** — synchronous call within the authorization path; must be fast,
  highly available, and **never lose** decision reasons for audit.
- **Feature store** — precomputed cardholder/merchant/device features.
- **Rules engine** — versioned rules with dry-run + shadow mode.
- **Model serving** — low-latency inference; A/B and shadow scoring; drift
  monitoring.
- **Event stream** — all transactions to a stream for feature computation and
  analytics.
- **Case management** — review queues, decisions, and SAR/AML integration.
- **Blocklists/hotlists** — card, device, IP, fingerprint, BIN blacklists with
  TTLs and reason codes.
- **Rate limiting & anti-enumeration** at the API gateway.

**Testing**:
- Replay historical data to measure detection rate / false-positive rate
  before go-live.
- Predefine SLA: p99 decision latency < 100 ms; graceful degradation (fail
  open with logging) vs fail closed policy by risk tier.

## 10.7 Fraud-specific compliance

- **Scheme reporting**: report fraud to networks (TC40/SAFE, SPA); scheme
  monitoring programs track chargeback and fraud ratios.
- **Regulatory**: SAR (suspicious activity reports) for AML; data-protection
  rules for storing device/PII.
- **Dispute evidence**: retain decision data to defend chargebacks.

## 10.8 Key references

- Visa — TC40 / SAFE reporting, Visa Risk Manager (VRM), Visa Advanced
  Authorization
- Mastercard — SPA (Secure Payment Alert), Decision Intelligence / AMR
- Stripe Radar / Adyen RevenueAccelerate — PSP fraud product docs
- Industry system-design articles on fraud scoring latency budgets (see
  `12-system-design.md` references)
