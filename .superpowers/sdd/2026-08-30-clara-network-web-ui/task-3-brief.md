# Task 3 Brief: environment config + shared types + money/date helpers

## Task Description (from plan)

Core lib primitives the dashboards depend on: env helper, admin types mirroring `internal/adminapi/store.go`, and money/date formatters.

**Files:**
- Create: `web/.env.local.example`, `web/src/lib/env.ts`, `web/src/lib/money/minor.ts`, `web/src/lib/date.ts`, `web/src/types/admin.ts`
- Modify: `web/.gitignore` (ignore `.env*.local`)

**Interfaces:**
- Consumes: Task 1 scaffold.
- Produces:
  - `getEnv()` throws with the name of any missing var (SUPABASE_URL, SUPABASE_ANON_KEY, CLARA_API_URL).
  - `fmtMinor(minor: number, currency?: string): string` (minor units → "1,234.56 EUR", default currency "EUR")
  - `fmtTs(iso: string): string` (ISO ts → regional date+time)
  - `Page<T>`/`DashboardSummary` types mirroring `internal/adminapi/store.go`.

## Steps

- **Step 1: env helper + example**

```ts
// src/lib/env.ts
const required = ["SUPABASE_URL", "SUPABASE_ANON_KEY", "CLARA_API_URL"] as const;
export function getEnv(): Record<typeof required[number], string> {
  const out = {} as Record<typeof required[number], string>;
  for (const k of required) {
    const v = process.env[k];
    if (!v) throw new Error(`missing env ${k}`);
    out[k] = v;
  }
  return out;
}
```

`.env.local.example`:
```
# Supabase project (project settings -> API)
SUPABASE_URL=https://<ref>.supabase.co
SUPABASE_ANON_KEY=<anon key>
# Railway Go adminapi URL (or http://localhost:18083 during local dev)
CLARA_API_URL=http://localhost:18083
```

- **Step 2: money + date helpers**

```ts
// src/lib/money/minor.ts
export function fmtMinor(minor: number, currency = "EUR"): string {
  const sign = minor < 0 ? "-" : "";
  const abs = Math.abs(minor);
  const major = Math.floor(abs / 100);
  const frac = String(abs % 100).padStart(2, "0");
  const intl = new Intl.NumberFormat("en-US").format(major);
  return `${sign} ${intl}.${frac} ${currency}`.trim();
}
```

```ts
// src/lib/date.ts
export function fmtTs(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleString("en-GB", { dateStyle: "medium", timeStyle: "short" });
}
```

- **Step 3: admin types**

```ts
// src/types/admin.ts
export interface Page<T> { items: T[]; total: number }
export interface DashboardSummary {
  transactions: number; clearingRecords: number; merchants: number;
  disputes: number; cards: number; tokens: number;
}
```

- **Step 4: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

NOTE: per the Task 1 reviewer note, a standalone `tsc --noEmit` on a pristine checkout can raise `Cannot find name 'LayoutProps'` from the generated `.next/types` until a `next build` has run. If `tsc --noEmit` trips only on that pre-existing generated-types issue (not your new files), run `cmd /c "cd web && npx.cmd next build"` FIRST to generate types, then rerun `tsc --noEmit` — expected pass. Report which sequence you used.

Expected: pass. Then `git commit -m "feat(web): env config, types, money/date helpers"` (stage `web/`).

## Global Constraints (binding)

- Node v24.19.0; npm via `npm.cmd`/`npx.cmd`. Windows PowerShell.
- Amounts: Admin API returns `amount_minor` as `BIGINT` minor units — `fmtMinor()` formats them. Demo currency: EUR (default).
- Never log/secrets: `.env.local.example` holds only placeholder values; `.gitignore` must exclude `.env*.local`.
- Single source of truth for money/time formatting: everything else imports these helpers.

## Controller Rulings for this task

- Client/components (Task 5/6 later) may NOT use `getEnv()` (browser cannot read arbitrary process.env) — this `env.ts` is server-side only. That wiring is handled in later tasks; here just create the module.
- Keep `fmtMinor` semantics exactly as shown (sign handling for negatives); it will be unit-tested in Task 8.