# Task 4 Brief: Supabase client (server + middleware) + roles

## Task Description (from plan)

Auth plumbing: per-request Supabase server client, middleware session-guard, and the role model.

**Files:**
- Create: `web/src/lib/supabase/server.ts`, `web/src/lib/supabase/middleware.ts`, `web/src/lib/roles.ts`, `web/src/middleware.ts`
- Modify: `web/package.json` (deps via CLI)

**Interfaces:**
- Consumes: Task 3 `getEnv()`.
- Produces:
  - `createServerClient()` per-request with cookie methods.
  - `updateSession(request)` used by middleware; redirects unauthenticated users to `/login`.
  - `Role = "scheme_operator" | "issuer" | "acquirer" | "merchant" | "viewer"`
  - `roleFromAppMetadata(meta?): Role | null` reads `user.app_metadata.role`.
  - `ROLE_LABEL: Record<Role, string>`
  - `HOME_BY_ROLE: Record<Role, string>`
  - `dashboardAccess(role): string[]`

## Steps

- **Step 1: install + package pin**

```bash
cmd /c "cd web && npm.cmd i @supabase/supabase-js @supabase/ssr"
```

- **Step 2: server client**

```ts
// src/lib/supabase/server.ts
import { createServerClient } from "@supabase/ssr";
import { cookies } from "next/headers";

export function createServerClient() {
  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
  const cookieStore = cookies();
  return createServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
    cookies: { getAll: () => cookieStore.getAll(), setAll: () => {} },
  });
}
```

- **Step 3: middleware client**

```ts
// src/lib/supabase/middleware.ts
import { createServerClient } from "@supabase/ssr";
import { NextResponse, type NextRequest } from "next/server";

export async function updateSession(request: NextRequest) {
  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
  let response = NextResponse.next({ request });
  const supabase = createServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
    cookies: {
      getAll: () => request.cookies.getAll(),
      setAll: (cookiesToSet) => {
        cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value));
        response = NextResponse.next({ request });
        cookiesToSet.forEach(({ name, value }) => response.cookies.set(name, value));
      },
    },
  });
  const { data } = await supabase.auth.getUser();
  if (!data.user) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("next", request.nextUrl.pathname);
    return NextResponse.redirect(url);
  }
  return response;
}
```

- **Step 4: roles module**

```ts
// src/lib/roles.ts
export type Role = "scheme_operator" | "issuer" | "acquirer" | "merchant" | "viewer";
export const ROLE_LABEL: Record<Role, string> = {
  scheme_operator: "Scheme Operator", issuer: "Issuer", acquirer: "Acquirer",
  merchant: "Merchant", viewer: "Viewer (HR)",
};
export function roleFromAppMetadata(meta?: Record<string, unknown>): Role | null {
  const r = meta?.role;
  return r === "scheme_operator" || r === "issuer" || r === "acquirer" || r === "merchant" || r === "viewer" ? r : null;
}
export const HOME_BY_ROLE: Record<Role, string> = {
  scheme_operator: "/ops", issuer: "/issuer", acquirer: "/acquirer",
  merchant: "/merchant", viewer: "/overview",
};
export const DASHBOARD_ACCESS: Record<Role, string[]> = {
  scheme_operator: ["/ops", "/transactions", "/clearing", "/settlement", "/ledger", "/cards", "/merchants", "/disputes"],
  issuer: ["/issuer", "/cards", "/tokens"],
  acquirer: ["/acquirer", "/merchants", "/funding", "/disputes"],
  merchant: ["/merchant", "/funding", "/disputes"],
  viewer: ["/overview"],
};
export function dashboardAccess(role: Role): string[] { return DASHBOARD_ACCESS[role] ?? []; }
```

- **Step 5: middleware.ts**

```ts
// src/middleware.ts
import { updateSession } from "@/lib/supabase/middleware";
import { type NextRequest } from "next/server";

export async function middleware(request: NextRequest) {
  return updateSession(request);
}
export const config = { matcher: ["/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)"] };
```

- **Step 6: verify + commit**

```bash
cmd /c "cd web && npx.cmd next build"
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Expected: pass (run `next build` first to generate `.next/types` if `tsc` trips on LayoutProps). Commit `feat(web): supabase auth clients, roles, auth middleware`.

## Global Constraints (binding)

- Node v24.19.0; npm via `npm.cmd`/`npx.cmd`. Windows PowerShell.
- Role vocabulary is the single source of truth in `web/src/lib/roles.ts`: `scheme_operator` | `issuer` | `acquirer` | `merchant` | `viewer`.
- Roles are read from Supabase Auth `app_metadata.role` (user JWT) — no separate role table.
- BFF/auth wiring must never expose secrets; anon key is public by design.

## Controller Rulings for this task

- Use `@/lib/roles` path alias in imports inside `roles.ts` consumers; the roles module itself is plain exports.
- `setAll: () => {}` in server client is intentional (server components are read-only for cookies) — do NOT "fix" it to write cookies; middleware handles refreshes.
- Keep the `require("../env")` shape as shown (it works in Next server bundling); if tsc complains about `require` typing, add `const env = require` won't be needed — `@types/node` is present; use `require` as in brief.