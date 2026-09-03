# Task 6 Brief: BFF data proxy + route handler

## Task Description (from plan)

Read-only BFF proxy over admin API. Guards with role + cookies; fetches from CLARA_API_URL; caches with Supabase Redis.

**Files:**
- Create: `web/src/app/api/data/[...path]/route.ts`, `web/src/lib/clara.ts`

**Interfaces:**
- Consumes: `createServerClient()` (ASYNC), `CLARA_API_URL` from `getEnv()`.
- Produces:
  - `RouteContext` (`params: Promise<{ path: string[] }>`).
  - `GET /api/data/<endpoint>` streaming proxy with `fetch`.
  - `claraFetch<T>(path, accessToken, {cache, ttl})` — typed GET with optional caching.
  - `UNAUTHORIZED` / `FORBIDDEN` `NextResponse` JSON errors.

## Steps

- **Step 1: clara.ts** (server-only fetch helper)

```ts
// src/lib/clara.ts
import "server-only";
import { getEnv } from "./env";

export function claraFetch<T>(path: string, accessToken: string, init?: RequestInit): Promise<T> {
  const { CLARA_API_URL } = getEnv();
  return fetch(`${CLARA_API_URL}${path}`, {
    headers: auth: { Authorization: `Bearer Bearer ${accessToken}`, Accept: "application/json" },
    ...init,
  } as RequestInit).then(r => {
    if (!r.ok) throw new Error(`clara ${path}: ${r.status}`);
    return r.json() as Promise<T>;
  });
}
```

NOTE: `"server-only"` package may need install: `cmd /c "cd web && npm.cmd i server-only"`. The code above has a syntax error (`headers: auth:` is wrong) — correct it as `headers: { Authorization: ..., Accept: ... }`.

- **Step 2: route handler** (server, auth-guarded)

```ts
// src/app/api/data/[...path]/route.ts
import { NextResponse } from "next/server";
import { createServerClient } from "@/lib/supabase/server";
import { getEnv } from "@/lib/env";
import { claraFetch } from "@/lib/clara";

export const dynamic = "force-dynamic";

type Res<T> = NextResponse<T | { error: string }>;

const ALLOWED_ROUTES: Record<string, string[]> = {
  scheme_operator: [
    "/dashboard", "/transactions", "/clearing/cycles", "/clearing/records",
    "/clearing/positions", "/settlement/instructions", "/settlement/prefunds",
    "/settlement/default-fund", "/ledger/accounts", "/ledger/entries",
  ],
  issuer: ["/dashboard", "/cards", "/tokens"],
  acquirer: ["/dashboard", "/acquirer", "/merchants", "/funding", "/disputes"],
  merchant: ["/dashboard", "/merchant", "/funding", "/disputes"],
};

export async function GET(
  req: Request,
  { params }: { params: Promise<{ path: string[] }> },
): Promise<Res<unknown>> {
  const supabase = await createServerClient();
  const { data, error } = await supabase.auth.getUser();
  if (error || !data.user) return NextResponse.json({ error: "unauthorized" }, { status: 401 });
  const role = data.user.app_metadata.role as string;
  const { path } = await params;
  const target = "/" + path.join("/");
  const allowed = ALLOWED_ROUTES[role] ?? [];
  if (!allowed.includes(target)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
  try {
    const body = await claraFetch(target, data.user.id, { cache: "no-store" });
    return NextResponse.json(body);
  } catch (e) {
    const status = e instanceof Error && /^clara \S+ 4\d\d/.test(e.message) ? 502 : 500;
    return NextResponse.json({ error: "upstream error", detail: e instanceof Error ? e.message : String(e) }, { status });
  }
}
```

- **Step 3: verify + commit**

```bash
cmd /c "cd web && npx.cmd next build"
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd eslint src/app/api src/lib/clara.ts"
```

Expected: pass. Commit `feat(web): BFF data proxy, role-gated admin API access`.

## Global Constraints (binding)

- Read-only proxy: method GET only. Admin API has no auth — the role gate in the BFF is THE authorization layer.
- Never log secrets. Never expose admin API directly (no direct client fetch to CLARA_API_URL).
- Roles from `app_metadata.role`.
- Node v24.19.0; `npm.cmd`. Windows PowerShell.

## Controller Rulings for this task

- `createServerClient()` is ASYNC — await it.
- The plan mentioned caching (Redis/Supabase). SIMPLIFY: `cache: "no-store"` for now (correctness first); caching re-examined in a later task if the demo needs it. Do NOT add a Redis client this task.
- The demo data is small; the BFF should remain readable and dependency-free. If `server-only` install fails with EALLOWSCRIPTS, prefix `set npm_config_allow_scripts=` as in Task 4.