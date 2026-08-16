# 17. Key Management, HSM & Message Security

## 17.1 What must be protected

Card messages carry three classes of secrets, each with its own protection:

| Secret | Example | Mechanism |
|--------|---------|-----------|
| **Card data at rest** | PAN, track 2, cardholder data | Encryption + PCI DSS (see `09-pci-dss.md`) |
| **PIN (at the host/switch)** | PIN block, PIN verification value (PVV) | ISO 9564 PIN blocks, encryption through an HSM/SCD |
| **Message integrity/authenticity** | Authorization, clearing, settlement messages | Retail MAC / digital signature per ISO 8583 |

## 17.2 PIN blocks — ISO 9564

ISO 9564-1 defines the standard PIN block formats:

- **Formats 0–3** (classic, 3DES-based): ISO-0 (standard, PAN-derived), ISO-1,
  ISO-2, ISO-3.
- **Format 4** (ISO 9564-1:2017): AES-based, combines the PIN and PAN; enables
  a move away from 3DES and toward modern algorithms.

**PCI PIN Security Standard Requirement 1** requires that all PIN
translations, regardless of format, occur **inside a Secure Cryptographic
Device (SCD)** — i.e., an HSM or a tamper-resistant security module. PIN blocks
must never be in clear outside the SCD.

### PCI SSC "Implementing ISO Format 4 PIN Blocks" highlights

- Format 4 uses AES and couples PIN with PAN in a way that resists dictionary
  attacks and replay.
- Migration path: keep ISO 0–3 for legacy interchange, adopt Format 4 for
  new/high-security environments.
- Applies the same SCD translation rules as older formats — every
  encrypt/decrypt/translate step happens in the HSM.

## 17.3 HSM usage across the network

An HSM (Hardware Security Module) is the tamper-resistant device that protects
cryptographic keys and performs cryptographic operations. In a card network:

- **Zone/master keys (ZMKs, TMKs)** — key-encryption keys distributed between
  network and members.
- **Transaction keys** — keys used to encrypt PINs and compute MACs for each
  acquirer/issuer link.
- **Key ceremonies** — controlled generation, split (dual control / M-of-N),
  distribution, rotation, and destruction of master keys.
- **DUKPT / derived keys** — per-terminal and per-transaction keys for POS/ATM
  (see ISO 9564-1 Annex for PIN key management under DUKPT).

## 17.4 Message authentication (MAC)

ISO 8583 messages between parties are authenticated with a **MAC** computed
over the message data (see `04-iso8583.md`):

- Algorithms: retail MAC (ISO 9797-1), CMAC, EMV-style; 3DES or AES.
- The MAC key is a transaction key derived from a zone key, managed inside the
  HSM.
- The network **verifies the acquirer MAC on inbound** and **computes the
  issuer MAC on outbound** — all MAC work in HSMs to protect the keys.
- Beyond MAC, sensitive fields (PIN, sometimes PAN) are **field-encrypted**
  end-to-end between the member and the network.

## 17.5 Key hierarchy design (reference architecture)

```
Root/Zone Master Key (ZM)          <- stored in HSM, never leaves clear
   |-- KEK (key-encrypting keys)   <- encrypts other keys in transit
        |-- PIN encryption key     <- per acquirer/issuer link
        |-- MAC key                <- per acquirer/issuer link
        |-- Tokenization keys      <- see 07-tokenization.md
        |-- Card personalization keys <- see 22-card-production-lifecycle.md
```

- Keys are identified by **Key Version Number (KVN)** and rolled on schedule.
- NIST SP 800-57 gives algorithm/lifetime guidance; card networks typically
  enforce 3DES-weakness rollovers and AES adoption.

## 17.6 Building secure key management into your platform

- Use **approved HSMs** (FIPS 140-2/140-3 validated) at every cryptographic
  boundary: network switch, issuer host, acquirer host, card production
  (see `22-card-production-lifecycle.md`).
- **Dual control**: no single person can create/use/delete a master key.
- **No cleartext keys in software** — keys live only inside HSMs or as
  encrypted key blocks (TR-31/TR-34).
- **Separate keys by function**: PIN, MAC, data encryption, card data, and
  certificates should use different keys.
- **Audit**: log every key event (create, rotate, load, destroy) and every
  SCD operation.

## 17.7 Key references

- ISO 9564-1 — *PIN management and security* (PIN block formats 0–4).
- PCI SSC — *Implementing ISO Format 4 PIN Blocks* (information supplement).
- PCI SSC — *PCI PIN Security Requirements* (Requirement 1: SCD controls).
- PCI SSC — *Point-to-Point Encryption (P2PE) Standard* and *HSM requirements*.
- NIST SP 800-57 — key management recommendations.
