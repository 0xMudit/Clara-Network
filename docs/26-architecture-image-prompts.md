# 26. Architecture Diagram Prompts (for image generation)

Use these five prompts with your image tool (e.g., Nano). Each prompt is
self-contained: specify a **clean architectural diagram**, flat vector style,
white background, and **short labels** (image models garble long text — keep
labels to 1–3 words). Generate one image per prompt, then use them in the
system-design documentation.

## Prompt 1 — End-to-end system architecture

> Create a clean professional system architecture diagram, flat vector style,
> white background. Layered top-down layout. Layer 1 (top, "ENTRY"): three
> boxes side by side labeled "POS TERMINAL", "E-COMMERCE", "MOBILE WALLET".
> Layer 2 ("ACQUIRER"): one wide box labeled "ACQUIRER HOST" with sub-boxes
> "Merchant Mgmt" and "Clearing Capture". Layer 3 ("NETWORK"): a wide
> highlighted box labeled "CLARA SWITCH" with sub-boxes "ISO 8583 Router",
> "Risk Engine", "Stand-In". Layer 4 ("ISSUER"): one box labeled "ISSUER
> HOST" with sub-boxes "Auth Decision", "Card Data", "Limits". Layer 5
> ("CORE"): three boxes labeled "LEDGER", "HSM", "KAFKA". Draw blue arrows
> from Layer 1 down through each layer to Layer 5, and a single green return
> arrow from "ISSUER HOST" back up to Layer 1 on the right side labeled
> "AUTH RESPONSE". Consistent blue/green color coding, no gradients, readable
> dark text labels.

## Prompt 2 — ISO 8583 authorization flow (sequence diagram)

> Create a clean sequence diagram, flat vector style, white background. Four
> vertical lifelines left to right labeled "ACQUIRER", "SWITCH", "RISK",
> "ISSUER". Draw these arrows in order: arrow "AUTH REQUEST (ISO 8583)" from
> ACQUIRER to SWITCH; arrow "SCORE CHECK" from SWITCH to RISK; arrow "APPROVE"
> from RISK back to SWITCH; arrow "AUTH FORWARD" from SWITCH to ISSUER;
> arrow "DECISION" from ISSUER back to SWITCH; arrow "AUTH RESPONSE" from
> SWITCH back to ACQUIRER. Below, a small red dashed arrow from SWITCH to
> ISSUER labeled "TIMEOUT" and a second red dashed arrow from SWITCH back to
> ACQUIRER labeled "STAND-IN 91/P". Use blue for normal flow, red for
> fallback, numbered steps 1-6, dark text, thin clean lines.

## Prompt 3 — Clearing, settlement & liquidity

> Create a clean financial infrastructure architecture diagram, flat vector
> style, white background. Top-down: top row, three boxes labeled "ISSUER",
> "ACQUIRER", "OTHER MEMBERS". Arrows from each down to a wide box labeled
> "CLARING ENGINE" (with sub-boxes "Capture", "Netting", "Reconciliation").
> Below it, a box labeled "SETTLEMENT ACCOUNTS" showing three stack cards
> labeled "Prefund", "Cap", "Net Position". Below, two boxes side by side:
> "DEFAULT FUND" and "CENTRAL BANK RTGS". Arrows: from each member box down
> to CLARING ENGINE labeled "CLEARING FILE"; from CLARING ENGINE to
> SETTLEMENT ACCOUNTS labeled "NET POSITION"; from SETTLEMENT ACCOUNTS to
> CENTRAL BANK RTGS labeled "SETTLE"; from DEFAULT FUND to CENTRAL BANK RTGS
> a dashed red arrow labeled "ON DEFAULT". Color: members blue, engine
> green, settlement gold, central bank dark blue. Clean minimal text.

## Prompt 4 — Issuer, tokenization & card stack

> Create a clean architectural diagram, flat vector style, white background.
> Left side, a vertical stack of boxes labeled "CARD PRODUCTION",
> "PERSONALIZATION", "HSM KEYS", "CARD DATA". Arrows from these into a wide
> center box labeled "ISSUER HOST" containing sub-boxes "BIN Ranges",
> "Cryptogram Verify", "Card Lifecycle". Right of the center box, a second
> wide box labeled "TOKEN VAULT" containing sub-boxes "Network Token",
> "PAR", "Token Status". From TOKEN VAULT, arrows to two smaller boxes
> labeled "APPLE PAY" and "GOOGLE PAY". A dashed secure channel drawn between
> "HSM KEYS" and "Cryptogram Verify" labeled "HSM". Color code: production
> purple, issuing blue, tokens green, wallets gray. Dark readable labels,
> no gradient, thin borders.

## Prompt 5 — Security, HSM & resilience

> Create a clean cybersecurity architecture diagram, flat vector style, white
> background. Center, a large shield-shaped box labeled "SECURITY CORE"
> containing three sub-boxes "HSM", "KEY HIERARCHY", "MAC / PIN BLOCK".
> Around it, four surrounding boxes labeled "SWITCH", "ISSUER", "ACQUIRER",
> "CARD PROD" each connected to the shield by a closed padlock-style dashed
> line. Below the shield, three boxes side by side labeled "SITE A",
> "SITE B", "DR SITE" with arrows between them labeled "ACTIVE-ACTIVE",
> "FAILOVER", "RTO 1H". A small red badge on "SITE A" labeled "OUTAGE" and a
> green badge on "SITE B" labeled "TAKEOVER". Use dark blue for the shield,
> green for resilience flow, red for the outage. Clean flat style, short
> labels, no photo backgrounds.
