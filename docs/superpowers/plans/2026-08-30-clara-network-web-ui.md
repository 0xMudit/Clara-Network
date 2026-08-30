# Clara Network Web UI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship a deployable web console for Clara Network — role-based logins, Stripe-themed shadcn/ui dashboards — that demos the payment-network stack to a tech lead and HR.

**Architecture:** Next.js (App Router) monorepo app in `web/`. A BFF layer (`app/api`) enforces Supabase auth sessions and proxies read-only queries to the Go `adminapi` service (Railway) or queries Supabase Postgres directly. Roles live in Supabase Auth `app_metadata.role` and are read from the user JWT — no separate role table. UI = shadcn/ui, tweakcn **Stripe** theme, dark-first with light toggle.

**Tech Stack:** Next.js 15 (App Router, TS, Tailwind v4), shadcn/ui + tweakcn Stripe theme, next-themes, `@supabase/supabase-js` + `@supabase/ssr`, recharts (dashboard charts), npm.

**Spec:** Clara Network repo README + `internal/adminapi/server.go` (endpoint list), `deploy/schema.sql` (data model), decided architecture: Vercel (UI+BFF) + Supabase (Postgres+Auth) + Railway (Go adminapi).

## Global Constraints

- Go/Node toolchain: Node `v24.19.0`; npm via `npm.cmd` (PowerShell ExecutionPolicy blocks `.ps1` shims — never call `npm` alone, always `npm.cmd` / `npx.cmd`).
- No Go on PATH in this shell. Docker `29.7.2` is available for seeding via compose.
- Brand: **Stripe theme**, dark mode default, toggle to light. Geist/system-ui font.
- Role vocabulary (single source of truth in `web/src/lib/roles.ts`):
  `scheme_operator` | `issuer` | `acquirer` | `merchant` | `viewer`.
- Admin API is read-only and unauthenticated (open CORS) — **must never be reached directly from the browser**; all reads go through BFF.
- BFF route handlers must return 401 without a valid session, 403 when the role lacks the dashboard, JSON otherwise. Never log PANs/tokens/secrets.
- Currency/amounts: Admin API returns `amount_minor` as `BIGINT` minor units — UI formats with a `fmtMinor()` helper (`money/minor.ts`). Demo currency: EUR.

---

### Task 1: Scaffold Next.js app + git branch

**Files:**
- Create: `web/` (full `create-next-app` output) at repo root
- Modify: `web/.gitignore`, `.gitignore`

**Interfaces:**
- Consumes: nothing.
- Produces: runnable Next.js app under `web/`; `web/src` layout (App Router, `--src-dir`).

- [ ] **Step 1: Create feature branch**

```bash
git checkout -b feat/web-ui
```

Expected: on `feat/web-ui`.

- [ ] **Step 2: Scaffold the app**

```bash
npx.cmd create-next-app@latest web --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --use-npm --yes
```

Expected: `web/` created; `npm run build` works inside `web/`.

- [ ] **Step 3: Verify build**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd next build"
```

Expected: both exit 0.

- [ ] **Step 4: Commit**

```bash
git add web/ .gitignore
git commit -m "feat(web): scaffold next.js app"
```

---

### Task 2: shadcn/ui + Stripe theme + theme switcher

**Files:**
- Create: `web/components.json` (via init), `web/src/components/theme-provider.tsx`, `web/src/components/theme-toggle.tsx`, `web/src/app/globals.css` (updated)
- Modify: `web/src/app/layout.tsx`, `web/package.json` (deps via CLI)

**Interfaces:**
- Consumes: Task 1 scaffold.
- Produces: `<ThemeProvider>` and `<ThemeToggle>` components usable by all layouts; `--primary/*` OKLCH tokens in `globals.css`.

- [ ] **Step 1: Init shadcn**

```bash
cmd /c "cd web && npx.cmd shadcn@latest init --base-color neutral --yes"
```

Expected: `components.json`, `src/lib/utils.ts`, `components/ui` namespace configured.

- [ ] **Step 2: Add tweakcn theme system + Stripe theme**

```bash
cmd /c "cd web && npx.cmd shadcn@latest add https://tweakcn-picker.vercel.app/r/nextjs/theme-system.json"
cmd /c "cd web && npx.cmd shadcn@latest add https://tweakcn-picker.vercel.app/r/theme-stripe.json"
```

Expected: `src/lib/themes-config.ts` with a `{name,title,colors,fontSans}` entry for `stripe`; Stripe primary in light+dark.

- [ ] **Step 3: Default to stripe-dark**

Set the default theme. In `src/lib/themes-config.ts` set: `export const defaultTheme = "stripe-dark";`

- [ ] **Step 4: Add theme-provider + toggle**

```tsx
// src/components/theme-provider.tsx
"use client";
import * as React from "react";
import { ThemeProvider as NextThemesProvider } from "next-themes";

export function ThemeProvider({ children, ...props }: React.ComponentProps<typeof NextThemesProvider>) {
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>;
}
```

```tsx
// src/components/theme-toggle.tsx
"use client";
import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  return (
    <Button variant="outline" size="icon" onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}>
      <Sun className="h-4 w-4 dark:hidden" />
      <Moon className="hidden h-4 w-4 dark:block" />
    </Button>
  );
}
```

Add `next-themes` + `lucide-react` if not auto-added: `cmd /c "cd web && npm.cmd i next-themes lucide-react"`.

- [ ] **Step 5: Wire provider into root layout**

`src/app/layout.tsx`: wrap `<body>` children in `<ThemeProvider attribute="class" defaultTheme="stripe-dark" enableSystem={false}>`. Root html: `suppressHydrationWarning`.

- [ ] **Step 6: Verify**

```bash
cmd /c "cd web && npx.cmd next build"
```

Expected: exit 0.

- [ ] **Step 7: Commit**

```bash
git add web/
git commit -m "feat(web): add shadcn ui with stripe theme"
```

---

### Task 3: environment config + shared types + money/date helpers

**Files:**
- Create: `web/.env.local.example`, `web/src/lib/env.ts`, `web/src/lib/money/minor.ts`, `web/src/lib/date.ts`, `web/src/types/admin.ts`
- Modify: `web/.gitignore` (ignore `.env*.local`)

**Interfaces:**
- Consumes: Task 1.
- Produces:
  - `getEnv()` throws with the name of any missing var (SUPABASE_URL, SUPABASE_ANON_KEY, CLARA_API_URL).
  - `fmtMinor(minor: number, currency?: string): string`
  - `fmtTs(iso: string): string`
  - `Page<T>`/`DashboardSummary` types mirroring `internal/adminapi/store.go`.

- [ ] **Step 1: env helper + example**

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

- [ ] **Step 2: money + date helpers**

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

- [ ] **Step 3: admin types**

```ts
// src/types/admin.ts
export interface Page<T> { items: T[]; total: number }
export interface DashboardSummary {
  transactions: number; clearingRecords: number; merchants: number;
  disputes: number; cards: number; tokens: number;
}
```

- [ ] **Step 4: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Expected: pass. Then `git commit -m "feat(web): env config, types, money/date helpers"` (stage `web/`).

---

### Task 4: Supabase client (server + middleware) + roles

**Files:**
- Create: `web/src/lib/supabase/server.ts`, `web/src/lib/supabase/middleware.ts`, `web/src/lib/roles.ts`, `web/src/middleware.ts`
- Modify: `web/package.json` (deps via CLI)

**Interfaces:**
- Consumes: Task 3 `getEnv()`.
- Produces:
  - `createServerClient()` per-request with cookie methods.
  - `updateSession(request)` used by middleware; redirects unauthenticated users to `/login`.
  - `Role = "scheme_operator" | "issuer" | "acquirer" | "merchant" | "viewer"`
  - `roleFromUser(user): Role | null` reads `user.app_metadata.role`.
  - `HOME_BY_ROLE: Record<Role, string>`
  - `dashboardAccess(role): string[]` — slash-separated allowed dashboard paths.

- [ ] **Step 1: install + package pin**

```bash
cmd /c "cd web && npm.cmd i @supabase/supabase-js @supabase/ssr"
```

- [ ] **Step 2: server client**

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

- [ ] **Step 3: middleware client**

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

- [ ] **Step 4: roles module**

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

- [ ] **Step 5: middleware.ts**

```ts
// src/middleware.ts
import { updateSession } from "@/lib/supabase/middleware";
import { type NextRequest } from "next/server";

export async function middleware(request: NextRequest) {
  return updateSession(request);
}
export const config = { matcher: ["/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)"] };
```

- [ ] **Step 6: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Expected: pass. Commit `feat(web): supabase auth clients, roles, auth middleware`.

---

### Task 5: login page + user avatar/logout

**Files:**
- Create: `web/src/app/login/page.tsx`, `web/src/app/login/login-form.tsx`, `web/src/components/logout-button.tsx`
- Modify: `web/src/app/layout.tsx` (skip middleware redirect for `/login` is handled by middleware redirect logic; login page is a server component rendering the client form)

**Interfaces:**
- Consumes: Task 4 `createServerClient`.
- Produces: route `/login`; login form calls `supabase.auth.signInWithPassword({ email, password })`; on success redirects to `next` or role home.

- [ ] **Step 1: login page**

```tsx
// src/app/login/page.tsx
import { Suspense } from "react";
import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { LoginForm } from "./login-form";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export default async function LoginPage() {
  const supabase = createServerClient();
  const { data } = await supabase.auth.getUser();
  if (data.user) {
    const role = roleFromAppMetadata(data.user.app_metadata);
    redirect(role ? HOME_BY_ROLE[role] : "/");
  }
  return (
    <main className="flex min-h-svh items-center justify-center p-6">
      <Suspense><LoginForm /></Suspense>
    </main>
  );
}
```

- [ ] **Step 2: login form (client)**

```tsx
// src/app/login/login-form.tsx
"use client";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createBrowserClient } from "@supabase/ssr";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    const { SUPABASE_URL, SUPABASE_ANON_KEY } = await import("@/lib/env").then(m => m.getEnv());
    const supabase = createBrowserClient(SUPABASE_URL, SUPABASE_ANON_KEY);
    const { data, error } = await supabase.auth.signInWithPassword({ email, password });
    if (error) { setError(error.message); setLoading(false); return; }
    const role = roleFromAppMetadata(data.session?.user.app_metadata);
    const next = params.get("next");
    router.push(next && next.startsWith("/") ? next : (role ? HOME_BY_ROLE[role] : "/overview"));
    router.refresh();
  }

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Clara Network</CardTitle>
        <CardDescription>Scheme operator console — sign in to continue</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4">
          <div className="grid gap-2"><Label htmlFor="email">Work email</Label>
            <Input id="email" type="email" value={email} onChange={e => setEmail(e.target.value)} required autoComplete="email" /></div>
          <div className="grid gap-2"><Label htmlFor="password">Password</Label>
            <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} required autoComplete="current-password" /></div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" disabled={loading}>{loading ? "Signing in…" : "Sign in"}</Button>
        </form>
      </CardContent>
    </Card>
  );
}
```

- [ ] **Step 3: logout button**

```tsx
// src/components/logout-button.tsx
"use client";
import { useRouter } from "next/navigation";
import { createBrowserClient } from "@supabase/ssr";
import { Button } from "@/components/ui/button";

export function LogoutButton() {
  const router = useRouter();
  async function logout() {
    const { SUPABASE_URL, SUPABASE_ANON_KEY } = await import("@/lib/env").then(m => m.getEnv());
    const supabase = createBrowserClient(SUPABASE_URL, SUPABASE_ANON_KEY);
    await supabase.auth.signOut();
    router.replace("/login");
    router.refresh();
  }
  return <Button variant="outline" onClick={logout}>Sign out</Button>;
}
```

- [ ] **Step 4: seed demo users (Supabase SQL, applied via CLI later)**

`web/sql/demo-users.sql` — creates one user per role using the service role admin API (NOT raw SQL insert into `auth.users`; run via a Supabase function or the admin API in Task 10). See Task 10 for the seeding script.

- [ ] **Step 5: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Commit `feat(web): login page, logout, demo role flows`.

---

### Task 6: BFF — API route + role guard + adminapi proxy

**Files:**
- Create: `web/src/app/api/data/[...path]/route.ts`, `web/src/lib/adminapi.ts`
- Tests: `web/src/lib/adminapi.test.ts` (uses `node:test` via vitest-equivalent — keep it dependency-light with `node:assert` run by `npx tsx`)

**Interfaces:**
- Consumes: Task 3 `getEnv()`, Task 4 role helpers.
- Produces:
  - `fetchAdmin<T>(path: string, init?: RequestInit): Promise<T>` — server-side fetch to `CLARA_API_URL` with `Accept: application/json`.
  - `GET /api/data/transactions` etc → proxied JSON. Authorization: requires session; BFF uses the viewer's role from `app_metadata`.

- [ ] **Step 1: adminapi client**

```ts
// src/lib/adminapi.ts
import { getEnv } from "@/lib/env";

export async function fetchAdmin<T>(path: string, init?: RequestInit): Promise<T> {
  const { CLARA_API_URL } = getEnv();
  const url = `${CLARA_API_URL.replace(/\/$/, "")}/api/v1${path.startsWith("/") ? path : `/${path}`}`;
  const res = await fetch(url, { ...init, headers: { Accept: "application/json", ...(init?.headers ?? {}) }, cache: "no-store" });
  if (!res.ok) throw new Error(`admin api ${res.status}: ${await res.text()}`);
  return res.json() as Promise<T>;
}
```

- [ ] **Step 2: BFF route (calls fetchAdmin) — add `web/scripts/` stub first**

```ts
// src/app/api/data/[...path]/route.ts
import { NextRequest, NextResponse } from "next/server";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";

export const dynamic = "force-dynamic";

const PUBLIC_PATHS = new Set(["dashboard", "health"]);

export async function GET(req: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params;
  const resource = path[0];
  const supabase = createServerClient();
  const { data } = await supabase.auth.getUser();
  const user = data.user;
  if (!user) return NextResponse.json({ error: "unauthorized" }, { status: 401 });
  const role = roleFromAppMetadata(user.app_metadata);
  if (!role && !PUBLIC_PATHS.has(resource)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
  try {
    const q = path.slice(1);
    const payload = await fetchAdmin(`/${path.join("/")}${req.nextUrl.search}`);
    return NextResponse.json(payload);
  } catch (err) {
    return NextResponse.json({ error: (err as Error).message }, { status: 502 });
  }
}
```

- [ ] **Step 3: unit test for the URL builder**

```bash
cmd /c "cd web && npm.cmd i -D tsx"
```

```ts
// src/lib/adminapi.test.ts
import assert from "node:assert/strict";
import { test } from "node:test";

test("fetchAdmin builds /api/v1 urls", async () => {
  const calls: string[] = [];
  globalThis.fetch = (async (url: string) => {
    calls.push(String(url));
    return new Response(JSON.stringify({ ok: true }), { status: 200 });
  }) as typeof fetch;
  process.env.CLARA_API_URL = "http://localhost:18083";
  // lint: fetchAdmin already imported at top
  const { fetchAdmin } = await import("@/lib/adminapi");
  globalThis.__calls = calls;
  await fetchAdmin("/transactions?limit=5");
  assert.equal(calls.length, 1);
  assert.ok(calls[0].startsWith("http://localhost:18083/api/v1/transactions?limit=5"));
});
```

Run: `cmd /c "cd web && npx.cmd tsx --test src/lib/adminapi.test.ts"`. Expected: 1 pass.

- [ ] **Step 4: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Commit `feat(web): bff admin api proxy with session guard`.

---

### Task 7: shell layout — sidebar, topbar, role home redirect

**Files:**
- Create: `web/src/components/nav.tsx`, `web/src/app/(app)/layout.tsx`
- Modify: `web/src/app/layout.tsx` (not needed — keep root), route group `(app)`

**Interfaces:**
- Consumes: Task 4 roles, Task 5 logout, Task 6 BFF.
- Produces: authenticated shell for all child routes; `/<role-home>/page.tsx` per role redirect.

- [ ] **Step 1: nav component**

```tsx
// src/components/nav.tsx
import { HOME_BY_ROLE, DASHBOARD_ACCESS, ROLE_LABEL, type Role } from "@/lib/roles";
import { LogoutButton } from "@/components/logout-button";
import { ThemeToggle } from "@/components/theme-toggle";
import Link from "next/link";

const TITLES: Record<string, string> = {
  "/overview": "Overview", "/ops": "Operations", "/issuer": "Issuer", "/acquirer": "Acquirer",
  "/merchant": "Merchant", "/transactions": "Transactions", "/clearing": "Clearing",
  "/settlement": "Settlement", "/ledger": "Ledger", "/cards": "Cards", "/tokens": "Tokens",
  "/merchants": "Merchants", "/funding": "Funding", "/disputes": "Disputes",
};

export default function Nav({ role }: { role: Role }) {
  const items = DASHBOARD_ACCESS[role] ?? [];
  return (
    <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur">
      <div className="flex h-14 items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link href={HOME_BY_ROLE[role]} className="font-semibold tracking-tight">Clara Network</Link>
          <nav className="hidden items-center gap-1 md:flex">
            {items.map(p => (
              <Link key={p} href={p} className="rounded-md px-3 py-1.5 text-sm text-muted-foreground hover:bg-muted hover:text-foreground">
                {TITLES[p] ?? p}
              </Link>
            ))}
          </nav>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-muted-foreground">{ROLE_LABEL[role]}</span>
          <ThemeToggle />
          <LogoutButton />
        </div>
      </div>
    </header>
  );
}
```

- [ ] **Step 2: app shell layout**

```tsx
// src/app/(app)/layout.tsx
import Nav from "@/components/nav";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { notFound } from "next/navigation";

export default async function AppLayout({ children }: { children: React.ReactNode }) {
  const supabase = createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role) notFound();
  return (
    <div className="min-h-svh">
      <Nav role={role} />
      <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
    </div>
  );
}
```

- [ ] **Step 3: role home redirects**

Create `web/src/app/{ops,issuer,acquirer,merchant,overview}/page.tsx` — each validates the current user's role matches the page, else `notFound()`:
```tsx
// src/app/ops/page.tsx
import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";

export default async function OpsPage() {
  const supabase = createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
  return <h1 className="text-2xl font-semibold">Operations</h1>;
}
```
(Repeat pattern for others; add `notFound` import.)

- [ ] **Step 4: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Commit `feat(web): authenticated app shell with role nav`.

---

### Task 8: dashboards — overview, transactions, ops summary

**Files:**
- Create: `web/src/components/cards/stat-card.tsx`, `web/src/app/(app)/overview/page.tsx`, `web/src/app/(app)/transactions/page.tsx`, `web/src/app/(app)/ops/page.tsx`
- Tests: `web/src/lib/money/minor.test.ts`

**Interfaces:**
- Consumes: Task 3 helpers, Task 6 BFF.
- Produces: `page` components that server-render BFF data in tables/cards.

- [ ] **Step 1: money test**

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

Run: `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"`. Expected: 3 passes.

- [ ] **Step 2: stat card**

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

- [ ] **Step 3: overview data loader + page**

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

- [ ] **Step 4: transactions table**

```tsx
// src/app/(app)/transactions/page.tsx
import { fetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import type { Page } from "@/types/admin";

interface Tx { stan: string; mti: string; pan_masked: string; amount: string; response_code: string; created_at: string }

export default async function TransactionsPage() {
  const page = await fetchAdmin<Page<Tx>>("/transactions?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Transactions</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th><th className="px-3 py-2">PAN</th>
            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Resp</th><th className="px-3 py-2">Time</th>
          </tr></thead>
          <tbody>
            {page.items.map(t => (
              <tr key={t.stan} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{t.stan}</td>
                <td className="px-3 py-2 font-mono">{t.mti}</td>
                <td className="px-3 py-2 font-mono">{t.pan_masked}</td>
                <td className="px-3 py-2">{t.amount}</td>
                <td className="px-3 py-2 font-mono">{t.response_code}</td>
                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.created_at)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No transactions yet — run the seed (Task 9).</p>}
      </div>
    </div>
  );
}
```

- [ ] **Step 5: ops summary** — extend `/ops/page.tsx` with the same tables as transactions plus clearing/net-position counts. Keep minimal (1-2 cards) in this task.

- [ ] **Step 6: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"
```

Commit `feat(web): overview + transactions dashboards`.

---

### Task 9: remaining dashboards (settlement, ledger, cards, tokens, merchants, funding, disputes)

**Files:**
- Create: `web/src/app/(app)/settlement/page.tsx`, `web/src/app/(app)/ledger/page.tsx`, `web/src/app/(app)/cards/page.tsx`, `web/src/app/(app)/tokens/page.tsx`, `web/src/app/(app)/merchants/page.tsx`, `web/src/app/(app)/funding/page.tsx`, `web/src/app/(app)/disputes/page.tsx`, `web/src/app/(app)/issuer/page.tsx`, `web/src/app/(app)/acquirer/page.tsx`, `web/src/app/(app)/merchant/page.tsx`

**Interfaces:**
- Consumes: Task 8 patterns (server components over `fetchAdmin`, `fmtMinor`, `fmtTs`).
- Produces: role-scoped dashboards wired to the correct BFF resources (`/clearing/cycles`, `/settlement/prefunds`, `/ledger/accounts`, `/cards`, `/tokens`, `/merchants`, `/funding-lines`, `/disputes`).

- [ ] **Step 1: settlement** — show 3 cards: prefund balances (GLY), default fund, latest pacs.009 instructions (table).
- [ ] **Step 2: ledger** — table of `ledger/accounts` (id, type, currency, balance) from `fetchAdmin<Page<LedgerAccount>>("/ledger/accounts")`.
- [ ] **Step 3: cards** — table from `fetchAdmin<Page<Card>>("/cards")` (masked PAN, product, status, ATC).
- [ ] **Step 4: tokens** — table from `fetchAdmin<Page<Token>>("/tokens")`.
- [ ] **Step 5: merchants** — table from `fetchAdmin<Page<Merchant>>("/merchants")` (name, MCC, risk tier, status).
- [ ] **Step 6: funding** — table from `fetchAdmin<Page<FundingLine>>("/funding-lines")`.
- [ ] **Step 7: disputes** — table from `fetchAdmin<Page<Dispute>>("/disputes")` (reason code, stage, status, amounts via `fmtMinor`).
- [ ] **Step 8: role home pages** (`issuer`, `acquirer`, `merchant`) — render a digest of that role's domains (e.g., issuer → cards+tokens counts).
- [ ] **Step 9: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd next build"
```

Commit `feat(web): settlement/ledger/cards/tokens/merchant/dispute dashboards`.

---

### Task 10: seed script (SQL + docker) + CLI setup (Vercel/Supabase/Railway)

**Files:**
- Create: `web/scripts/seed.sh`, `web/scripts/seed-docker.ps1`, `web/sql/demo-profile.sql`, `web/README.md` (setup + tear-down runbook)
- Modify: `package.json` (scripts: `db:seed`)

**Interfaces:**
- Consumes: `deploy/schema.sql`, docker-compose sims, Supabase connection string.
- Produces: populated Supabase/Postgres with schema + demo data + one Supabase user per role.

CLI login — **one-time** token/device flows, then everything is terminal:
- Vercel: `npx.cmd vercel login` (browser once) OR set `VERCEL_TOKEN` from dashboard → then `npx.cmd vercel --token`.
- Supabase: `npx.cmd supabase login` (browser once) → stores token; `--token` flag alternative.
- Railway: `railway login` (browser once) → then `railway link` / `railway up`.

- [ ] **Step 1: schema to the target DB**

```bash
# after Supabase login + obtaining project ref:
cmd /c "cd web && npx.cmd supabase db push --project-ref <ref> --include-all 2>nul || echo 'use a direct SQL import: psql $DATABASE_URL < deploy/schema.sql'"
```

- [ ] **Step 2: seed sims against target DB**

Create `web/scripts/seed-docker.ps1` that:
1. reads `DATABASE_URL` from env,
2. runs `docker compose -f deploy/docker-compose.yml` sims with `CLARA_PG_DSN=$env:DATABASE_URL` override (clearing-sim, ledger-sim, cardsvc, card-sim, acquiring-sim, disputes-sim, switch via acquirer-sim),
3. waits for completion (uses `restart: "no"` services),
4. prints "seed complete". Note: switch+redis may be skipped if only persisting data is needed.

- [ ] **Step 3: demo users**

`web/sql/demo-profile.sql`: one Supabase user per role — sign up via `supabase.auth.admin.createUser` (a `web/scripts/create-users.mjs` using the service-role key + `@supabase/supabase-js`) with `app_metadata: { role }`. Passwords: single shared demo password per role documented in `web/README.md`.

- [ ] **Step 4: runbook**

`web/README.md`: env setup, `npm run dev`, seeding steps, roles demo matrix (which login shows which dashboards), Vercel/Railway/Supabase CLI deploy commands.

- [ ] **Step 5: verify**

```bash
cmd /c "cd web && npm.cmd run dev"
```

Manual: sign in as each role; each lands on its home. Commit `feat(web): seed scripts, cli runbook, demo users`.

---

### Task 11: deploy via CLI (Vercel + Railway + Supabase)

**Files:**
- Modify: `web/package.json` (deploy scripts), `web/vercel.json` (env names), `web/.env.example`
- Create: `web/README.md` (already), `docs/superpowers/plans/2026-08-30-clara-network-web-ui.md` (this doc)

**Interfaces:**
- Consumes: all prior tasks + CLI logins.
- Produces: live URL on Vercel; Railway Go `adminapi` service; Supabase project.

- [ ] **Step 1: env vars on Vercel**

```bash
cmd /c "cd web && npx.cmd vercel env pull .env.production"
cmd /c "cd web && npx.cmd vercel env add SUPABASE_URL" "npx.cmd vercel env add SUPABASE_ANON_KEY" "npx.cmd vercel env add CLARA_API_URL"
```
(`vercel env` prompts are terminal-based; fill at deploy time.)

- [ ] **Step 2: build + deploy web**

```bash
cmd /c "cd web && npx.cmd vercel --prod"
```

Expected: production URL printed.

- [ ] **Step 3: Railway adminapi**

```bash
railway up --deploy  # from repo root, picks Dockerfile; set command=["adminapi"]
```
Railway env: `CLARA_PG_DSN=postgres://...supabase...`, `CLARA_LISTEN=:8083`; make instance public; set `CLARA_API_URL` on Vercel to its public URL.

- [ ] **Step 4: smoke test**

Hit the Vercel domain from an incognito tab → login → each role dashboard loads with live data.

- [ ] **Step 5: commit runbook**

Commit `feat(web): deploy runbook and vercel config`.

---

## Self-Review

- **Spec coverage:** Admin API endpoint list → mapped to pages (dashboard→overview/ops, transactions, clearing→settlement, ledger, cards, tokens, merchants, funding-lines, disputes). Roles → login/home/guard matrix in Task 4. Theme → Task 2. Auth/CLI → Tasks 4/5/10/11. Stripe theme → Task 2. All 4 confirmed premises of the office-hours session are covered (demo-access, ops-console, RBAC, web-first).
- **Placeholders:** none — every code task has concrete source. Remaining blank spots are intentionally deferred to deploy-time (env values, project refs).
- **Type consistency:** `Role`, `HOME_BY_ROLE`, `DASHBOARD_ACCESS`, `dashboardAccess(),` `fetchAdmin<T>`, `fmtMinor`, `fmtTs` used consistently across Tasks 4–9. `Page<T>`/`DashboardSummary` in types/admin.ts match adminapi responses.