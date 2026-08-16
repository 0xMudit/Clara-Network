# 9. PCI DSS Compliance

## 9.1 What PCI DSS is

The **Payment Card Industry Data Security Standard (PCI DSS)** is a set of
security requirements from the **PCI Security Standards Council (PCI SSC)** —
the body created in 2006 by Visa, Mastercard, Amex, Discover, and JCB. It
protects **account data** (cardholder data and sensitive authentication data)
for every entity that **stores, processes, or transmits** card data, or whose
systems could impact the **Cardholder Data Environment (CDE)**.

> "PCI DSS provides a baseline of technical and operational requirements
> designed to protect payment account data." — PCI SSC

Enforcement is **contractual**, via payment brands and acquirers (not a
government agency). Compliance is validated with the **SAQ** (Self-Assessment
Questionnaire) or an **ROC** (Report on Compliance) by a **QSA** (Qualified
Security Assessor), depending on volume level.

## 9.2 Account data taxonomy

| Account data | Includes | Storage after authorization |
|--------------|----------|------------------------------|
| **Cardholder data (CHD)** | PAN, cardholder name, expiration date, service code | Only PAN may be stored; must be rendered **unreadable** |
| **Sensitive authentication data (SAD)** | Full track data (magstripe/chip equivalent), card verification code (CVV/CVC/CVV2/CID), PINs/PIN blocks | **Must NOT be stored**, even encrypted |

The **PAN** is the defining element of cardholder data. Name/expiry/service
code are protected when present with PAN. SAD is never retained after
authorization.

## 9.3 PCI DSS v4.0 — the 12 requirements

v4.0 (published 2022, mandatory transition; **v4.0.1** June 2024) keeps the
12-requirement structure under six goals, and offers **Defined Approach** and
**Customized Approach** assessment paths. Requirements:

### Build and Maintain a Secure Network and Systems
1. **Install and maintain network security controls** — firewall/NSC policy,
   control of traffic into and within the CDE, segmentation.
2. **Apply secure configurations to all system components** — hardening
   standards, change default passwords, remove unnecessary services/software.

### Protect Account Data
3. **Protect stored account data** — minimize storage; render PAN unreadable
   (encryption/hashing/truncation/tokenization); secure keys and the key
   lifecycle; **SAD never stored after authorization**.
4. **Protect cardholder data with strong cryptography during transmission over
   open/public networks** — TLS 1.2+ for PAN in transit; encrypt or secure
   sessions.

### Maintain a Vulnerability Management Program
5. **Protect all systems and networks from malicious software** — anti-malware
   on all relevant systems.
6. **Develop and maintain secure systems and applications** — secure SDLC,
   patch management, vulnerability remediation (priorities/timeframes per
   severity).

### Implement Strong Access Control Measures
7. **Restrict access by business need-to-know** — least privilege, access
   review.
8. **Identify users and authenticate access** — unique IDs, strong
   authentication (MFA for admin/remote access per v4.0).
9. **Restrict physical access** — facility controls, media protection, device
   tamper protection.

### Regularly Monitor and Test Networks
10. **Log and monitor all access** — audit logging of all access to system
    components and CHD; protect and review logs.
11. **Test security systems and networks regularly** — vulnerability scanning,
    penetration testing (internal/external), segmentation testing, and for
    SAQs the quarterly ASV scans (external scanning).

### Maintain an Information Security Policy
12. **Support information security with organizational policies and programs** —
    policy, risk assessments (targeted risk analysis), **annual scope
    confirmation (12.5.2)**, third-party service provider management (12.8),
    incident response (12.10), security awareness (12.6).

## 9.4 Scope: the Cardholder Data Environment (CDE)

The CDE includes:
- System components (people, process, tech) that **store, process, or transmit**
  CHD/SAD.
- Systems with **unrestricted connectivity** to those components.
- Systems that could **impact the security** of CHD/SAD (auth servers, remote
  access, logging, backup, failover).

**Scope reduction** is the single highest-leverage compliance decision:
- **Tokenization** (network or gateway) — replace PANs; remove PAN from your
  environment (see `07-tokenization.md`).
- **Edge tokenization / hosted fields** (Stripe Elements, Adyen Web Components)
  — PAN never touches your servers; qualifies for **SAQ A** (short) instead of
  **SAQ D** (full questionnaire, ~300+ controls).
- **P2PE** (point-to-point encryption) — encrypted at the terminal, decrypted
  at a validated P2PE solution provider; can reduce scope.
- **Network segmentation** — isolate the CDE from the rest of the network.

## 9.5 Validation levels (U.S. example, Visa/Mastercard)

| Level | Criteria (Visa) | Validation |
|-------|------------------|------------|
| 1 | >6M Visa/MC transactions/year (or any level per acquirer/brand) | Annual **ROC by QSA** + quarterly ASV scans + Attestation of Compliance (AOC) |
| 2 | 1M–6M | Annual **SAQ-D** or ROC + quarterly ASV |
| 3 | 20k–1M e-comm | Annual SAQ + quarterly ASV |
| 4 | <20k e-comm | Annual SAQ |

Service providers (processors, gateways, acquirers) generally require **ROC +
AOC** and listing in the **Visa Global Registry of Service Providers** and
Mastercard's equivalent.

## 9.6 Key controls to implement (engineering view)

- **Encryption**: AES-256 for PAN at rest; TLS 1.2+ (ideally 1.3) in transit;
  key management with an **HSM** and defined key lifecycle (generation, use,
  rotation, retirement). PCI forbids using the same key across production and
  test, and restricts key reuse.
- **HSM**: mandatory-usage patterns for PIN handling, cryptogram verification,
  and key storage (see `12-system-design.md`).
- **Tokenization/vault**: PCI-validated (or compliant) tokenization removes PAN
  from most systems.
- **Logging**: time-stamped audit trails of access to PAN, tamper protection,
  centralized SIEM.
- **MFA**: v4.0 makes MFA required for all access into the CDE (admin + remote).
- **Inventory & CDE mapping**: document every location/flows of account data
  (12.5).
- **Third-party (TPSP) management**: v4.0 requires you to manage TPSP
  relationships (12.8) and TPSPs to support customer compliance (12.9).
- **Customized approach**: document **targeted risk analyses** for any
  requirement met via alternative controls.

## 9.7 Compliance checklist for a card platform

- [ ] Define and document the CDE; confirm scope at least annually (12.5.2).
- [ ] Eliminate SAD storage (no CVV, track data, PINs after authorization).
- [ ] Encrypt PAN at rest; enforce TLS in transit; strong key management.
- [ ] Harden servers/containers; disable default accounts; patch promptly.
- [ ] Network segmentation + firewall rules governing CDE ingress/egress.
- [ ] Unique accounts + MFA for all CDE access.
- [ ] Centralized logging + daily review + SIEM alerting.
- [ ] Quarterly ASV scans; annual internal/external pentests; segmentation
    testing.
- [ ] Incident response plan (12.10) and security awareness training (12.6).
- [ ] TPSP due diligence, monitoring, and evidence collection (12.8/12.9).
- [ ] File the correct **SAQ/ROC + AOC** with your acquirer/brands.

## 9.8 Key references

- PCI SSC — PCI DSS v4.0 / v4.0.1 "Requirements and Testing Procedures"
- PCI SSC — "PCI DSS v4.0 in a Nutshell"; "Quick Reference Guide"; "Summary of
  Changes 3.2.1 → 4.0"
- Middlebury (public copy of the standard) — PCI-DSS-v4.0.1.pdf
- Venn, CrowdStrike, pcicompliancehub — "12 requirements of PCI DSS v4.0"
