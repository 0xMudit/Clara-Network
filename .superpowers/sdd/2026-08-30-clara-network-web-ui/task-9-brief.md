# Task 9 Brief: remaining dashboards — clearing, settlement, ledger, cards, tokens, merchants, funding, disputes + role homes

## Task Description (from plan, with controller corrections)

All remaining read-only dashboards over the BFF, plus the issuer/acquirer/merchant role-home digests.

**Files:**
- Fix: `web/src/app/api/data/[...path]/route.ts` (Step 0 — query-string forwarding)
- Create: `web/src/app/(app)/{clearing,settlement,ledger,cards,tokens,merchants,funding,disputes}/page.tsx`
- Create/replace: `web/src/app/(app)/{issuer,acquirer,merchant}/page.tsx` (digests, keep Task 7 gates)

**Interfaces:**
- Consumes: `fetchAdmin` (Task 8), `fmtMinor`, `fmtTs`, BFF allowlist (Task 6 fix), real admin-API JSON shapes (below).
- Produces: dashboards wired to `/clearing/cycles`, `/clearing/records?cycle=`, `/clearing/positions?cycle=`, `/settlement/prefunds`, `/settlement/default-fund`, `/settlement/instructions?cycle=`, `/ledger/accounts`, `/cards`, `/tokens`, `/merchants`, `/disputes`.

## PLAN DEFECTS / CONTROLLER RULINGS (binding)

- **D0 (BFF drops query strings):** `route.ts` builds `target = "/" + path.join("/")` — NO query string. `?limit=`, `?cycle=` are silently dropped (Task 8's `/transactions?limit=50` only worked because the admin API's default limit is 50). `?cycle=` is REQUIRED by clearing records/positions and settlement instructions (400 without it). Step 0 fixes this BEFORE the pages land.
- **D1 (missing /clearing page):** the plan's nav (`DASHBOARD_ACCESS`/`TITLES`) references `/clearing` but no task creates it. THIS task creates `(app)/clearing/page.tsx`.
- **D2 (no /funding-lines endpoint):** the plan's funding page targets `/funding-lines`, which does NOT exist (admin API serves `/merchants/{id}/funding`). Build the funding page from REAL data: `Page<Merchant>` via `/merchants`, showing funding-relevant columns (name, risk tier, reserve balance, volume, funding delay, transaction limit).
- **D3 (ledger currency hidden):** `LedgerAccount.Currency` is `json:"-"` — the ledger table must NOT show a currency column. Balance is minor units → `fmtMinor`.
- **D4 (envelope vs Page):** clearing/ledger/prefund etc. return `{items: [...]}` (NOT `Page<T>`). Cards/tokens/merchants/disputes/transactions return `Page<T>`. Read the shapes below and fetch accordingly.
- **D5 (gates):** every new dashboard gates via `DASHBOARD_ACCESS` membership (single source of truth) — `const role = roleFromAppMetadata(data.user?.app_metadata); if (!role || !DASHBOARD_ACCESS[role].includes("/<route>")) notFound();` (with `await createServerClient()`). The three role-home pages KEEP their strict single-role gates from Task 7 (they are HOME_BY_ROLE targets).

## Step 0 — BFF query-string forwarding (`route.ts`), own commit

In `web/src/app/api/data/[...path]/route.ts`:
1. Build the allowed-path check from the PATHNAME only, then append the search string for the upstream call. Search sits in `request.nextUrl.search` (use the existing `req`'s URL). Implement:

```ts
// after: const { path } = await params;
const url = new URL(req.url);
const base = "/" + path.join("/");
const allowed = ALLOWED_ROUTES[role] ?? [];
if (!allowed.includes(base)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
try {
  const body = await claraFetch(base + url.search, data.user.id, { cache: "no-store" });
```

(`url.search` is `""` when there are no params — safe.) Verify: `cmd /c "cd web && npx.cmd tsc --noEmit"` + `next build`. Commit `fix(web): forward query strings in BFF proxy`.

## Real admin-API JSON shapes (source: internal/adminapi/*)

- `/clearing/cycles` → `{ items: string[] }` (newest first per `ORDER BY cycle_id DESC`)
- `/clearing/records?cycle=C` → `{ items: ClearingRecord[] }`, `ClearingRecord = {cycleId, stan, mti, sender, receiver, amountMinor, interchange, currency, refId}`
- `/clearing/positions?cycle=C` → `{ items: NetPosition[] }`, `NetPosition = {cycleId, member, net}` (minor)
- `/settlement/prefunds` → `{ items: PrefundAccount[] }`, `PrefundAccount = {member, balance, cap}` (minor)
- `/settlement/default-fund` → `{ balance: number }` (minor)
- `/settlement/instructions?cycle=C` → `{ items: SettlementInstruction[] }`, `SettlementInstruction = {cycleId, msgId, member, amount, direction, currency, instruction, final}` (instruction is an ISO time; amount minor)
- `/ledger/accounts` → `{ items: LedgerAccount[] }`, `LedgerAccount = {id, type, balance}` (balance minor; NO currency field — D3)
- `/cards?limit=50` → `Page<Card>`, `Card = {ref, panHash, panMask, bin, expiry, status, product, lastAtc}`
- `/tokens?limit=50` → `Page<Token>`, `Token = {token, par, status, bin, requestor, deviceId, createdAt}`
- `/merchants?limit=50` → `Page<Merchant>`, `Merchant = {id, name, dba, taxId, mccs, status, riskTier, reserveRateBps, fundingDelayDays, transactionLimit, reserveBalance, volume, declineReason?, approvedAt}` (reserveBalance/volume/transactionLimit minor)
- `/disputes?limit=50` → `Page<Dispute>`, `Dispute = {id, refId, merchantId, cardholder, amountMinor, currency, reasonCode, category, stage, status, filedAt, responseDue, respondedAt?, escalatedAt?, decision?, winner?, decisionAt?, disputeFee, arbitrationFee, note?}`

## Step 1 — clearing page `(app)/clearing/page.tsx` (D5 gate scheme_operator only)

Fetch `/clearing/cycles`; if `items.length === 0` render `<h1>Clearing</h1>` + empty-state `"No clearing cycles yet — run the seed (Task 10)."`; else pick `const cycle = items[0]`, fetch `/clearing/records?cycle=${encodeURIComponent(cycle)}` and `/clearing/positions?cycle=${encodeURIComponent(cycle)}`. Render: `<h1>Clearing</h1>`, a `text-sm text-muted-foreground` caption `Cycle {cycle}`, then two bordered tables (same table style as Task 8): records (STAN, MTI, sender, receiver, amount via `fmtMinor(r.amountMinor, r.currency)`, interchange via `fmtMinor(r.interchange, r.currency)`) and positions (member, net via `fmtMinor(p.net)`). Table headers per Task 8 conventions (`px-3 py-2`, mono where apt).

## Step 2 — settlement page `(app)/settlement/page.tsx` (scheme_operator gate)

Three sections:
1. Prefund balances: `fetchAdmin<{items: PrefundAccount[]}>("/settlement/prefunds")` → table (member, balance `fmtMinor(balance, "EUR")`, cap `fmtMinor(cap, "EUR")`).
2. Default fund: `fetchAdmin<{balance: number}>("/settlement/default-fund")` → a `StatCard`-style readout (reuse `StatCard` from Task 8) `Default fund` / `fmtMinor(balance, "EUR")`.
3. Latest instructions: reuse `/clearing/cycles` → latest cycle → `/settlement/instructions?cycle=...` → table (msgId, member, direction, amount `fmtMinor(amount, currency)`, final ✅/—, instruction `fmtTs`).
Empty state for no cycles: "No settlement instructions yet — run the seed (Task 10)."

## Step 3 — ledger `(app)/ledger/page.tsx` (scheme_operator gate)

`fetchAdmin<{items: LedgerAccount[]}>("/ledger/accounts")`. Table: id (mono), type, balance `fmtMinor(balance)` (NOT `fmtMinor(balance, currency)` — no currency field, D3). Empty state standard.

## Step 4 — cards `(app)/cards/page.tsx` (gate: DASHBOARD_ACCESS has `/cards` → scheme_operator|issuer)

`fetchAdmin<Page<Card>>("/cards?limit=50")`. Table: ref (mono), panMask (mono), bin (mono), product, status, expiry, lastAtc.

## Step 5 — tokens `(app)/tokens/page.tsx` (gate: `/tokens` → scheme_operator|issuer)

`fetchAdmin<Page<Token>>("/tokens?limit=50")`. Table: token (mono), par (mono), bin, requestor, status, createdAt `fmtTs`.

## Step 6 — merchants `(app)/merchants/page.tsx` (gate: `/merchants` → scheme_operator|acquirer|merchant)

`fetchAdmin<Page<Merchant>>("/merchants?limit=50")`. Table: name, dba, mccs (join ", "), riskTier, status, reserveBalance `fmtMinor`, volume `fmtMinor`.

## Step 7 — funding `(app)/funding/page.tsx` (gate: `/funding` → acquirer|merchant) — D2

`fetchAdmin<Page<Merchant>>("/merchants?limit=50")`. Table titled "Acquirer funding" (columns: name, fundingDelayDays, transactionLimit `fmtMinor`, reserveBalance `fmtMinor`, volume `fmtMinor`). This is the real-data funding view (no `/funding-lines` endpoint exists).

## Step 8 — disputes `(app)/disputes/page.tsx` (gate: `/disputes` → scheme_operator|acquirer|merchant)

`fetchAdmin<Page<Dispute>>("/disputes?limit=50")`. Table: id (mono), cardholder, reasonCode, category, stage, status, amount `fmtMinor(amountMinor, currency)`, filedAt `fmtTs`, responseDue `fmtTs` (empty → "—").

## Step 9 — role-home digests (replace Task 7 placeholder JSX; KEEP each page's strict gate)

- `(app)/issuer`: after gate, fetch `/dashboard` → StatCard grid (Cards, Tokens) + hint linking `/cards` and `/tokens`. `const [,cards,tokens]`… simplest: cards/tokens counts from DashboardSummary.
- `(app)/acquirer`: after gate, fetch `/dashboard` → StatCards (Merchants, Disputes) + links to `/merchants`, `/disputes`.
- `(app)/merchant`: after gate, fetch `/dashboard` → StatCards (Merchants, Disputes) + link to `/merchants`, `/disputes`.

Each keeps `notFound()` + `await createServerClient()` from Task 7; only the `return` JSX changes (like Task 8 did for ops). If a page's gate shared pattern makes the fetch unreachable for a gate-only home, keep gate first then fetch.

## Step 10 — verify + commit

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd next build"
```

Expected: pass (all pages under `(app)` inherit layout `force-dynamic`; the BFF fix commit is Step 0). Final commit: `feat(web): settlement/ledger/cards/tokens/merchant/dispute dashboards`.

## Global Constraints (binding)

- Node v24.19.0; `npm.cmd`/`npx.cmd`; Windows PowerShell.
- Server components only; `fetchAdmin` server-only; NEVER call CLARA_API_URL from pages.
- `fmtMinor`/`fmtTs` are the only formatters. Standard empty-state copy: `<p className="p-4 text-sm text-muted-foreground">No <thing> yet — run the seed (Task 10).</p>`.
- Finish with NO dangling nav links: at end of this task every path in `DASHBOARD_ACCESS` (all roles) must resolve to a page (overview/ops/transactions/clearing/settlement/ledger/cards/tokens/merchants/funding/disputes/issuer/acquirer/merchant) — assert in your self-review.
- Unrelated working-tree files (`.github/**`, AGENTS.md, docs/28-*.md) — never stage/commit.
- No branch creation; commit on `feat/web-ui`.

## Controller Rulings

- Step 0 is required and its own commit; the rest is ONE commit (per plan).
- Do not add client-side interactivity (no filtering UI, no pagination controls, no cycle dropdown). Static server-rendered dashboards; latest cycle only.
- If any page shows a real runtime error in the demo because the DB is empty (e.g. cycle fetch returns empty), the empty-state handles it — no try/catch scaffolding.