# Task 8 Brief: dashboards — stat card, overview, transactions, ops summary + money unit test

## Task Description (from plan, with controller corrections)

First real data dashboards over the BFF, plus the `fmtMinor` unit test the plan schedules here.

**Files:**
- Create: `web/src/lib/adminapi.ts` (NEW — plan references it but never defines it), `web/src/components/cards/stat-card.tsx`, `web/src/app/(app)/overview/page.tsx`, `web/src/app/(app)/transactions/page.tsx`, `web/src/app/(app)/ops/page.tsx` (replace Task 7 placeholder), `web/src/lib/money/minor.test.ts`
- Tests: `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"` — expected 1 test, 3 asserts pass.

**Interfaces:**
- Consumes: Task 3 helpers (`fmtMinor`, `fmtTs`), Task 6 BFF (`/api/data/...`), real admin API shapes (below).
- Produces: `fetchAdmin<T>(path)` server helper; overview/transactions/ops pages.

## PLAN DEFECTS — controller rulings (follow these; do NOT implement raw plan code):

- **D1 (missing lib):** the plan's pages import `fetchAdmin` from `@/lib/adminapi` but no task creates that file. THIS TASK creates it. It is a server-only helper that fetches the same-origin BFF and turns non-2xx into a thrown error:

```ts
// src/lib/adminapi.ts
import "server-only";

export async function fetchAdmin<T>(path: string): Promise<T> {
  const url = `/api/data${path}`;
  const res = await fetch(url, { cache: "no-store" });
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    const msg = body && (body as { error?: string }).error ? (body as { error: string }).error : String(res.status);
    throw new Error(`adminapi ${path}: ${msg}`);
  }
  return res.json() as Promise<T>;
}
```

- **D2 (Tx shape wrong):** the plan's `Tx` interface names fields that don't exist. The admin API returns `AuditEvent` (`internal/adminapi/store_tx.go:11-18`) as: `{ stan, mti, pan, amount, responseCode, destination, createdAt }`. `amount` is a STRING (already formatted); there is NO `pan_masked`/`response_code`/`created_at`/`dest`. Use the corrected interface and columns below.
- **D3 (async client):** all pages under `(app)` run in the shell; Task 8 pages are data pages — they do NOT need another `getUser()` call (the `(app)` layout already gates). Data pages just call `fetchAdmin` server-side. Do NOT re-await `createServerClient` in these pages.
- The BFF allowlist (Task 6 fix) permits `/dashboard` and `/transactions?limit=...` for scheme_operator; `/dashboard` for viewer. Pages in this task only call those.

## Steps

- **Step 1: money test (verbatim from plan):**

```ts
// src/lib/money/minor.test.ts
import assert from "node:assert/strict";
import { test } from "node:test";
import { fmtMinor } from "./minor";

test("fmtMinor formats minor units", () => {
  assert.equal(fmtMinor(123456), "1,234.56 EUR");
  assert.equal(fmtMinor(0), "0.00 EUR");
  assert.equal(fmtMinor(-5000), "-50.00 EUR");
});
```

Run `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"`. If `tsx` isn't installed, add it as a dev dep: `cmd /c "cd web && npm.cmd i -D tsx"`. Expected 3 asserts pass.

- **Step 2: stat card (verbatim from plan):**

```tsx
// src/components/cards/stat-card.tsx
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function StatCard({ title, value, hint }: { title: string; value: string; hint?: string }) {
  return (
    <Card><CardHeader className="pb-2"><CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle></CardHeader>
      <CardContent><div className="text-2xl font-semibold">{value}</div>
        {hint && <p className="mt-1 text-xs text-muted-foreground">{hint}</p>}</CardContent></Card>
  );
}
```

- **Step 3: overview page** (from plan; correct: it's already at `(app)/overview/page.tsx` from Task 7 — REPLACE its placeholder content):

```tsx
// src/app/(app)/overview/page.tsx
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";

export default async function OverviewPage() {
  const d = await fetchAdmin<DashboardSummary>("/dashboard");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Network overview</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Transactions" value={d.transactions.toLocaleString()} />
        <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} />
        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
        <StatCard title="Cards" value={d.cards.toLocaleString()} />
        <StatCard title="Tokens" value={d.tokens.toLocaleString()} />
      </div>
    </div>
  );
}
```

- **Step 4: transactions page — CORRECTED shape (D2):**

```tsx
// src/app/(app)/transactions/page.tsx
import { fetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import type { Page } from "@/types/admin";

interface AuditEvent {
  stan: string;
  mti: string;
  pan: string;
  amount: string;
  responseCode: string;
  destination: string;
  createdAt: string;
}

export default async function TransactionsPage() {
  const page = await fetchAdmin<Page<AuditEvent>>("/transactions?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Transactions</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th><th className="px-3 py-2">PAN</th>
            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Resp</th><th className="px-3 py-2">Dest</th>
            <th className="px-3 py-2">Time</th>
          </tr></thead>
          <tbody>
            {page.items.map(t => (
              <tr key={t.stan} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{t.stan}</td>
                <td className="px-3 py-2 font-mono">{t.mti}</td>
                <td className="px-3 py-2 font-mono">{t.pan}</td>
                <td className="px-3 py-2">{t.amount}</td>
                <td className="px-3 py-2 font-mono">{t.responseCode}</td>
                <td className="px-3 py-2 font-mono">{t.destination}</td>
                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.createdAt)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No transactions yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}
```

- **Step 5: ops summary — REPLACE Task 7 placeholder** at `src/app/(app)/ops/page.tsx`. Keep minimal this task: reuse `StatCard`s from `/dashboard` with an ops slant and link to the transactions page. Suggested content (role is already scheme_operator-only via page gate from Task 7 — KEEP that gate and file; only replace the returned JSX):

```tsx
// src/app/(app)/ops/page.tsx  (keep the existing async role gate; swap the return JSX below)
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import Link from "next/link";
import type { DashboardSummary } from "@/types/admin";

// inside the existing OpsPage, after the gate:
const d = await fetchAdmin<DashboardSummary>("/dashboard");
return (
  <div className="grid gap-4">
    <h1 className="text-2xl font-semibold">Operations</h1>
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <StatCard title="Transactions today" value={d.transactions.toLocaleString()} hint="Authorizations via the switch" />
      <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} hint="Captured this settlement window" />
      <StatCard title="Merchants onboarded" value={d.merchants.toLocaleString()} />
    </div>
    <p className="text-sm text-muted-foreground">
      <Link href="/transactions" className="underline underline-offset-4 hover:text-foreground">View transaction log →</Link>
    </p>
  </div>
);
```

Note: keep the import of `notFound` and the scheme_operator gate that Task 7 already put in this file. `d` is re-fetched (fine for this demo). The previous `await createServerClient()` in this page may be dropped ONLY if the gate no longer needs it — check how Task 7 wrote the gate; preserve its behavior exactly (gate must still run).

- **Step 6: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"
cmd /c "cd web && npx.cmd next build"
```

Expected: tests 3 passes, build pass. Commit `feat(web): overview + transactions dashboards`.

## Global Constraints (binding)

- Node v24.19.0; `npm.cmd`/`npx.cmd`. Windows PowerShell.
- `fetchAdmin` is server-only (`server-only` package already installed by Task 6). Only ever same-origin BFF — never call CLARA_API_URL from pages.
- Never log secrets. `fmtMinor`/`fmtTs` are the only formatters.
- Unrelated working-tree files (`.github/**`, AGENTS.md, docs/28-*.md) — never stage/commit.
- No branch creation; commit on `feat/web-ui`.

## Controller Rulings

- D1 creates `fetchAdmin`; D2 corrects the AuditEvent shape; D3 pages don't re-auth (layout gates). The BFF does the auth + allowlist.
- If `/api/data/dashboard` returns 401/403 in a future local run with no session, `fetchAdmin` throws → Next shows its default error page. Acceptable for now (login is required to reach these pages in practice). Do not build error-boundary UI this task.
- Ops page: preserve Task 7's role gate; replace only the returned JSX body.