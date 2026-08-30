# Task 8 Report: dashboards — stat card, overview, transactions, ops summary + money unit test

## Status: DONE

## What I implemented

1. **`web/src/lib/adminapi.ts` (NEW — ruling D1).** `fetchAdmin<T>(path)` server-only helper (`import "server-only"`), fetches same-origin `/api/data${path}` with `cache: "no-store"`, throws `adminapi ${path}: ${msg}` on non-2xx (body `.error` preferred, else status code). Verbatim from brief.

2. **`web/src/components/cards/stat-card.tsx` (NEW).** `StatCard({ title, value, hint })` over the existing shadcn `Card` primitives. Verbatim from brief.

3. **`web/src/app/(app)/overview/page.tsx` (REPLACED Task 7 placeholder).** Server component calling `fetchAdmin<DashboardSummary>("/dashboard")` and rendering 6 `StatCard`s (transactions, clearingRecords, merchants, disputes, cards, tokens). Matches `web/src/types/admin.ts` `DashboardSummary`. No re-auth (ruling D3 — the `(app)` layout gates).

4. **`web/src/app/(app)/transactions/page.tsx` (NEW).** Server component calling `fetchAdmin<Page<AuditEvent>>("/transactions?limit=50")` with the **corrected D2 shape** — `{ stan, mti, pan, amount, responseCode, destination, createdAt }`, `amount` a string — verified against `internal/adminapi/store_tx.go:10-18` (JSON tags `stan/mti/pan/amount/responseCode/destination/createdAt`) and `Page{items,total}` against the Go store's `Page` struct. Renders a table using `fmtTs(t.createdAt)`; empty state line for pre-seed.

5. **`web/src/app/(app)/ops/page.tsx` (REPLACED Task 7 JSX only).** KEPT the Task 7 scheme_operator gate verbatim (`createServerClient()` → `getUser()` → `roleFromAppMetadata(...) !== "scheme_operator" → notFound()`) — the `notFound` import and gate behavior are unchanged. Added `const d = await fetchAdmin<DashboardSummary>("/dashboard")` after the gate and replaced the returned JSX with the ops stat cards + "View transaction log →" link to `/transactions`.

6. **`web/src/lib/money/minor.test.ts` (NEW, verbatim from brief).** `fmtMinor` unit test — 3 asserts (positive with comma group, zero, negative).

7. **BUG FIX — `web/src/lib/money/minor.ts` (surgical).** The brief-mandated test asserts `fmtMinor(-5000) === "-50.00 EUR"`, but Task 3's implementation returned `"- 50.00 EUR"` (sign + literal space, then `.trim()` — trim doesn't remove the interior space). First run **failed** on that assert. Fixed `minor.ts` to `` `${sign}${intl}.${frac} ${currency}` `` (dropped the spurious space and `.trim()`). Positive/zero output unchanged (`"1,234.56 EUR"`, `"0.00 EUR"`). This is the plan's scheduled money unit test catching a real latent bug — the fix makes `fmtMinor` semantics match the plan's explicit "→ 1,234.56 EUR" + sign-handling spec. `minor.ts` now matches the formatting style already used by the test's other two assertions.

8. **Dev dep:** installed `tsx@^4.23.13` (`web/package.json` + lockfile).

## What I tested (exact commands + results)

| Command | Result |
|---|---|
| `cmd /c "cd web && npm.cmd i -D tsx"` | added 3 packages, 0 vulnerabilities. Note: 2 install scripts (esbuild, unrs-resolver) gated by npm `allow-scripts` config — esbuild still ran the test fine. |
| `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"` | **1st run: FAIL** — `fmtMinor(-5000)` returned `'- 50.00 EUR'` vs expected `'-50.00 EUR'` (1 fail). After the `minor.ts` fix: **1 test, 1 pass, 0 fail, 3 asserts pass.** |
| `cmd /c "cd web && npx.cmd tsc --noEmit"` | passed, no output |
| `cmd /c "cd web && npx.cmd eslint "`. (7 changed files) | passed, no output |
| `cmd /c "cd web && npx.cmd next build"` | Compiled successfully; routes `/overview`, `/ops`, `/transactions` all build as ƒ (dynamic) — no static prerender failure, so no pages needed their own `force-dynamic` (the `(app)` layout's covers them). Only warning: pre-existing `middleware` deprecation from Task 7, not from this task. |

## Files changed

- Created: `web/src/lib/adminapi.ts`, `web/src/components/cards/stat-card.tsx`, `web/src/app/(app)/transactions/page.tsx`, `web/src/lib/money/minor.test.ts`
- Modified: `web/src/app/(app)/overview/page.tsx`, `web/src/app/(app)/ops/page.tsx`, `web/src/lib/money/minor.ts`, `web/package.json`, `web/package-lock.json`
- Commit: `0f711eb` `feat(web): overview + transactions dashboards` (9 files, +638/−4). Staged `web/` only — no `.github/**`, AGENTS.md, docs/28-*.md, etc. Branch `feat/web-ui` (no new branch).

## Self-review findings

- **Ruling D1** satisfied: `fetchAdmin` created server-only; pages never touch `CLARA_API_URL` (the BFF route / `claraFetch` does the proxying).
- **Ruling D2** satisfied: `AuditEvent` interface matches `internal/adminapi/store_tx.go` JSON tags exactly; `Page<T>{items,total}` matches the Go `Page` struct.
- **Ruling D3** satisfied: only ops keeps its page-level `getUser()` (its Task 7 gate — required to preserve; brief says keep it); overview/transactions rely on the `(app)` layout gate + BFF auth.
- **Allowlist check:** pages only call `/dashboard` (overview, ops) and `/transactions?limit=50` (transactions) — all allowed for the calling roles (scheme_operator; viewer gets `/dashboard` only, and the transactions page isn't in viewer nav).
- ops gate behavior preserved bit-for-bit (same auth client call, same role check, same `notFound`), only the returned JSX + data fetch changed; `d` re-fetch noted as acceptable in the brief.
- All brief code blocks used verbatim except the intentional `minor.ts` fix (documented above).

## Issues / concerns

- **`minor.ts` deviation (required).** The brief's own test would fail against Task 3's `fmtMinor` (space bug). The minimal fix changes only negative-value output; positive/zero untouched. Flagging loudly so the controller can veto if the Task 3 formatter was intentionally `"- 50.00 EUR"` (unlikely — the test is this task's spec).
- New files from the `write` tool lack a trailing newline (repo convention is mixed; nav.tsx and the Task 7 `(app)/layout.tsx` also lack one). Cosmetic; build/lint/test unaffected.
- npm `allow-scripts` gating means esbuild's postinstall didn't run under this npm config; `tsx` still executed (test passed), so no action needed. Could bite on a future `npm ci` on a fresh Windows box if `set npm_config_allow_scripts=` isn't used — pre-existing project npm config quirk, not task-specific.