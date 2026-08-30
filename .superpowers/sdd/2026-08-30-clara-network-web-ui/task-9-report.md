# Task 9 Report — remaining dashboards + role homes

**Status:** DONE_WITH_CONCERNS
**Branch:** `feat/web-ui` (no branch created)
**Date:** 2026-08-31

## What I implemented

### Commit 1 — `cca1bfb` `fix(web): forward query strings in BFF proxy` (Step 0)

`web/src/app/api/data/[...path]/route.ts`: the allowlist check now runs against the
pathname-only `base` (`"/" + path.join("/")`), and the upstream call appends the
incoming query string via `url.search` (`base + url.search`), where `url = new URL(req.url)`.
`url.search` is `""` when the request has no query params, so the no-query case is
unchanged. This unblocks the mandatory `?cycle=` params for `/clearing/records`,
`/clearing/positions`, and `/settlement/instructions` (400 without them), and makes
`?limit=` behave explicitly instead of silently relying on the admin API's default 50.

### Commit 2 — `15397cc` `feat(web): settlement/ledger/cards/tokens/merchant/dispute dashboards` (Steps 1–9)

Seven new dashboards + three role-home digests (all server components, all gated via
`DASHBOARD_ACCESS` membership per D5, all fetching same-origin `/api/data` through
`fetchAdmin<T>`):

- `(app)/clearing/page.tsx` (Step 1, scheme_operator): fetches `/clearing/cycles`; on
  empty → `<h1>Clearing</h1>` + standard empty-state; else latest cycle, then
  `/clearing/records?cycle=` + `/clearing/positions?cycle=` (both `encodeURIComponent`'d).
  Caption `Cycle {cycle}`; records table (STAN/MTI mono, sender, receiver, amount +
  interchange via `fmtMinor(amount, currency)`); positions table (member, net via `fmtMinor`).
- `(app)/settlement/page.tsx` (Step 2, scheme_operator): three sections — Default fund
  `StatCard` (`fmtMinor(balance, "EUR")`), Prefund balances table (`fmtMinor(..., "EUR")`),
  Latest instructions (latest cycle → `/settlement/instructions?cycle=`, columns msgId mono,
  member, dir, amount `fmtMinor(amount, currency)`, final ✅/—, instruction `fmtTs`).
  No cycles → instructions section shows the standard empty-state.
- `(app)/ledger/page.tsx` (Step 3, scheme_operator): `/ledger/accounts`, columns id (mono),
  type, balance `fmtMinor(balance)` — NO currency column (D3, `LedgerAccount.Currency` is `json:"-"`).
- `(app)/cards/page.tsx` (Step 4, `/cards` in access): `Page<Card>` via `/cards?limit=50`;
  ref/panMask/bin mono, product, status, expiry, lastAtc.
- `(app)/tokens/page.tsx` (Step 5, `/tokens` in access): `Page<Token>` via `/tokens?limit=50`;
  token/par mono, bin, requestor, status, createdAt `fmtTs`.
- `(app)/merchants/page.tsx` (Step 6, `/merchants` in access): `Page<Merchant>` via
  `/merchants?limit=50`; name, dba, mccs joined, riskTier, status, reserveBalance/volume `fmtMinor`.
- `(app)/funding/page.tsx` (Step 7, `/funding` in access): real-data funding view per D2
  (no `/funding-lines` endpoint) — `Page<Merchant>` via `/merchants?limit=50`, titled
  "Acquirer funding", columns name, fundingDelayDays, transactionLimit/reserveBalance/volume `fmtMinor`.
- `(app)/disputes/page.tsx` (Step 8, `/disputes` in access): `Page<Dispute>` via
  `/disputes?limit=50`; id mono, cardholder, reasonCode, category, stage, status,
  amount `fmtMinor(amountMinor, currency)`, filedAt `fmtTs`, responseDue `fmtTs` (empty → "—").
- `(app)/issuer/page.tsx`, `(app)/acquirer/page.tsx`, `(app)/merchant/page.tsx` (Step 9):
  kept each strict single-role gate (`!== role` → `notFound()`, with awaited
  `createServerClient()`), replaced placeholder JSX with `/dashboard` digests:
  issuer → StatCards Cards + Tokens, links `/cards`, `/tokens`; acquirer → StatCards
  Merchants + Disputes, links `/merchants`, `/disputes`; merchant → StatCards Merchants +
  Disputes, links `/funding`, `/disputes` (see Concern 1).

## What I tested

All on Windows PowerShell with `cmd /c "cd web && ..."`:

1. `npx.cmd tsc --noEmit` — pass (both commits).
2. `npx.cmd next build` — pass both times. Route list after commit 2 includes all new pages:
   `/cards /clearing /disputes /funding /ledger /merchants /settlement /tokens` (all `ƒ` dynamic,
   inheriting `(app)` layout `force-dynamic`). Expected middleware deprecation warning only.
3. `npx.cmd tsx --test src/lib/money/minor.test.ts` — 1 test, 1 pass, 0 fail.
   (Project tests are `node:test`; a throwaway `npx vitest` run behaved badly and was discarded —
   nothing from it landed in the repo.)
4. `npx.cmd eslint src` — 0 errors, 1 pre-existing warning (unused `getEnv` import in
   `route.ts`, a Task-6 deferred minor already tracked in progress.md; not introduced here).

The Step-0 query-forward is also verified by inspection against the real admin API:
`getClearingRecords`/`getNetPositions`/`getSettlementInstructions` in
`internal/adminapi/handlers.go` return 400 on a missing `?cycle=`, and the pages only call
those endpoints with a concrete cycle from `/clearing/cycles`.

## Files changed

**Commit `cca1bfb`** (1 file):
- `M web/src/app/api/data/[...path]/route.ts`

**Commit `15397cc`** (11 files, +561 −3):
- `A web/src/app/(app)/clearing/page.tsx`
- `A web/src/app/(app)/settlement/page.tsx`
- `A web/src/app/(app)/ledger/page.tsx`
- `A web/src/app/(app)/cards/page.tsx`
- `A web/src/app/(app)/tokens/page.tsx`
- `A web/src/app/(app)/merchants/page.tsx`
- `A web/src/app/(app)/funding/page.tsx`
- `A web/src/app/(app)/disputes/page.tsx`
- `M web/src/app/(app)/issuer/page.tsx`
- `M web/src/app/(app)/acquirer/page.tsx`
- `M web/src/app/(app)/merchant/page.tsx`

No BFF allowlist changes were needed (Task 6 fix already covers every route called). No
unrelated working-tree files (`.github/**`, AGENTS.md, docs/28-*, etc.) were staged.

## Self-review findings

- **No dangling nav links — ASSERTED:** every path in `DASHBOARD_ACCESS` (all roles)
  resolves to a page file under `web/src/app/(app)/`:
  - scheme_operator: `/ops /transactions /clearing /settlement /ledger /cards /merchants /disputes` — all exist.
  - issuer: `/issuer /cards /tokens` — all exist.
  - acquirer: `/acquirer /merchants /funding /disputes` — all exist.
  - merchant: `/merchant /funding /disputes` — all exist.
  - viewer: `/overview` — exists.
  Confirmed both against the filesystem (14 page.tsx under `(app)`) and the `next build`
  route listing. `Nav` `TITLES` already covers every one of these paths.
- **D4 envelopes respected:** clearing/ledger/prefund/default-fund/cycles return `{items}` /
  `{balance}` and are fetched as such; cards/tokens/merchants/disputes are `Page<T>`.
- **D3 respected:** ledger has no currency column.
- **D0 respected:** BFF forwards query strings; pages use `encodeURIComponent` on cycle values.
- **D5 respected:** every new dashboard gates via `DASHBOARD_ACCESS[role].includes("/<route>")`:
  `if (!role || !DASHBOARD_ACCESS[role].includes("/clearing")) notFound();` — release 1.
- **D2 respected:** funding page built on real `/merchants` data, no `/funding-lines` reference.
- Interfaces checked against the real Go JSON tags (`internal/adminapi/store_{cards,clearing,ledger,acquiring,disputes}.go`).

## Issues / concerns

1. **Merchant-home link deviates from the Step 9 text (DONE_WITH_CONCERNS).** Step 9 says the
   merchant home links to `/merchants` + `/disputes`, but `DASHBOARD_ACCESS.merchant` =
   `["/merchant", "/funding", "/disputes"]` — no `/merchants`. Since D5 mandates a
   `DASHBOARD_ACCESS` gate on the merchants page, a merchant user clicking `/merchants`
   would hit a 404. I chose the zero-dead-link option: the merchant home links to
   `/funding` + `/disputes` (both in merchant's nav/access) instead of `/merchants`. This
   also avoids exposing competitor merchant data (reserve/volume) to merchant-role users,
   which excluding `/merchants` from merchant's access appears intentional for. If the
   reviewer prefers literal Step 9 compliance, the fix is to add `"/merchants"` to
   `DASHBOARD_ACCESS.merchant` (merchant nav would then also show a Merchants item).
2. Pre-existing, not introduced here: unused `getEnv` import in `route.ts` (Task-6 deferred minor).

## Report path

`C:\Users\mudit\Videos\openSourceProjects\Clara-Network\.superpowers\sdd\2026-08-30-clara-network-web-ui\task-9-report.md`