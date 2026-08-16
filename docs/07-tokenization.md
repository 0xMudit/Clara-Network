# 7. Tokenization & Network Tokens

## 7.1 What payment tokenization is

**Payment tokenization** replaces the **PAN** with a unique surrogate value — an
**EMV Payment Token** — that is constrained in where and how it can be used (a
specific merchant, device, or transaction type). It removes the most valuable
fraud target from merchant and wallet environments.

> "Reducing the value and limiting the use of stolen or compromised payments
> information is critical... EMV Payment Tokenisation replaces valuable card
> data with payment tokens." — EMVCo

Standard: **EMV Payment Tokenisation Specification – Technical Framework**
(EMVCo), v2.x. Related: the **Network Token** model operated by Visa, Mastercard
et al., and **acquirer/security tokens** (Stripe, gateway tokens) which never
cross the network.

## 7.2 Token ecosystem roles

| Role | Function |
|------|----------|
| **Token Requestor** | Entity authorized to request tokens from a TSP (merchant, wallet, gateway, PSP) |
| **Token Service Provider (TSP)** | Generates/stores tokens, maps token↔PAN, validates usage; registers via EMVCo TSP Registration Programme (TSP Code) |
| **Token Requestor ID (TRID)** | Unique ID per requestor, issued by the TSP/scheme |
| **Payment Account Reference (PAR)** | 29-character value linking a PAN and its tokens; enables fraud/AML/loyalty analytics without the PAN |
| **Token Reference ID** | Opaque reference for the token, used in APIs instead of the token number |
| **BIN Controller / BCID** | Issuer-side registration to link tokens to PAN via PAR |
| **Issuer** | Approves token requests (with **ID&V**), controls token lifecycle |
| **Merchant Identifier** | Constrains a token to one merchant |

## 7.3 Token characteristics

- Same **format/length** as a PAN (e.g. 16 digits) so back-end systems need no
  changes; starts with the issuer's BIN range (or a token BIN).
- **Domain-limited** — e.g. only usable at one merchant, on one device, for
  e-commerce, or for a limited time.
- **Token cryptogram** — each transaction adds a fresh **cryptogram** (CVC3,
  dynamic). "The token is only useful with the correct cryptogram, which changes
  per transaction."
- **Lifecycle** — active/suspended/resumed/deleted; can be re-provisioned
  without re-entering card data; token can be replaced without card replacement.

## 7.4 Tokenization vs other techniques

| Technique | What it does | Crosses network? |
|-----------|--------------|------------------|
| **EMV Payment Token (network token)** | Replaces PAN end-to-end from merchant to issuer | Yes |
| **Acquirer/gateway token** (Stripe, Adyen) | Replaces PAN between merchant and processor | No — processor converts to PAN |
| **Format-preserving encryption / vault** | Encrypts/hashes PAN at rest | No |
| **Point-to-point encryption (P2PE)** | Encrypts at the terminal, decrypts at the gateway | No |

Merchant/acquirer tokens reduce PCI scope at the merchant, but a data breach of
the processor still exposes PANs. **Network tokens** remove the PAN from the
entire flow.

## 7.5 Network token benefits (from Visa/Mastercard public data)

- **Fraud reduction**: Visa reports token transactions drive ~30% reduction in
  online fraud vs PAN, and a 4%+ authorization uplift.
- **Auth-rate uplift**: CNP token transactions see ~3–4.6% higher authorization
  rates (fewer false declines, because the network confirms the underlying card
  is still valid and can auto-update the token on card reissue).
- **Card-on-file updates**: TSP can refresh expiry/credentials automatically
  without merchant re-entry.
- **Limited impact** of a compromise: tokens for one merchant are useless
  elsewhere.

## 7.6 Token provisioning flow (Visa Token Service example)

```
Cardholder enrolls PAN + security code with a digital wallet / merchant
            │
            ▼
Token Requestor ── request token ──▶ TSP (e.g. Visa Token Service)
            ▲                             │  (with ID&V attributes)
            │                             ▼
            │                          Issuer (ID&V: approve/decline/step-up)
            │                             │  approve
            │                             ▼
            │                        TSP issues Token + cryptogram material
            │  token returned            │
            └────────────────────────────┘
Then:
Token + cryptogram ──▶ merchant ──▶ acquirer ──▶ network ──▶ issuer (detokenize via TSP)
```

ID&V (Identification & Verification) can be:
- **On-behalf-of** (network validates)
- **Issuer-interactive** (issuer receives the ID&V call, approves/declines,
  delivers card art, terms, token reference)

## 7.7 Token lifecycle management APIs

Issuer-side capabilities (Visa Token Service, Mastercard MDES):

- **Token provision / notification**
- **Token inquiry** (details, device/risk info, list tokens for a PAN)
- **Token lifecycle controls** (activate, suspend, resume, delete)
- **Update PAN / expiry (VAU — Visa Account Updater / Account Updater)**
- **Token reference ID** usage instead of the token number in APIs

Mastercard's equivalent is **MDES (Mastercard Digital Enablement Service)** —
provisioning, lifecycle, and **CoF (card-on-file)** APIs; token expiry
synchronization via account updater.

## 7.8 Token acceptance & routing rules

- **PAN tokenization does not change BIN routing requirements** — token BINs
  must be routable by the acquirer/network.
- For **card-present** (tap-to-pay in wallets), tokens are provisioned to the
  SE/secure element or TEE and used with a device cryptogram.
- For **e-commerce**, tokens are used with CVC2 equivalent (dynamic CVC/CVC3)
  and **3DS** (see `08-3ds.md`).
- **PAR** lets you link token and PAN transactions for analytics without
  storing PANs.

## 7.9 Implementing tokenization in your platform

**If you are a TSP:**
- Build a **token vault**: token↔PAN mapping, tokenization/format-preserving
  crypto, domain controls, lifecycle state machine.
- Register with **EMVCo TSP Registration Programme** for a global TSP Code.
- Integrate ID&V flows with issuers.
- Maintain **cryptogram verification** (HSM) for token-based transactions.

**If you are an issuer:**
- Integrate TSP APIs (Visa VTS / Mastercard MDES) for ID&V, lifecycle, and
  notifications.
- Maintain PAR linking and account-updater feeds.
- Verify cryptograms and honor domain controls at authorization.

**If you are a merchant/PSP:**
- Use a PCI-validated tokenization provider or gateway tokens; prefer **network
  tokens** for auth-rate and security benefits.
- Store only the token + PAR + last-4; never PAN/CVV (see `09-pci-dss.md`).

## 7.10 Key references

- EMVCo — "EMV Payment Tokenisation" (Technical Framework, Guide to Use Cases,
  TSP/BCID registration programmes)
- Visa — "Token Service Provisioning and Credential Management" developer docs;
  "A Deep Dive into Tokenized Transactions"
- Mastercard Developers — "Tokenization" (Gateway); MDES docs
- W3C — "Tokenized Card Payment" (PaymentRequest integration)
