# Task 7 Brief: shell layout — nav, topbar, role home pages

## Task Description (from plan, with controller corrections)

Authenticated app shell for all role-scoped pages: a topbar nav (role-aware links, role label, theme toggle, logout) and a `(app)` route-group layout that gates on the logged-in user's role. Plus the five role-home pages.

**Files:**
- Create: `web/src/components/nav.tsx`, `web/src/app/(app)/layout.tsx`
- Create: `web/src/app/(app)/{ops,issuer,acquirer,merchant,overview}/page.tsx` (role homes)
- Modify: none (root layout untouched)

**Interfaces:**
- Consumes: Task 4 roles module, Task 5 logout button, Task 2 theme toggle.
- Produces: gated shell layout; role-home routes.

## PLAN DEFECTS — controller rulings (follow these; do NOT implement the raw plan code):

- **D1 (async client):** `createServerClient()` is ASYNC (Next 16 async cookies — recorded ruling). The plan's raw code calls it without `await`. You MUST `const supabase = await createServerClient();` in the layout AND in every role-home page.
- **D2 (route group):** the plan mixes `src/app/ops/page.tsx` (Task 7) with `src/app/(app)/overview/page.tsx` (Task 8). For the `(app)` layout to apply, ALL authenticated shell pages must live INSIDE the `(app)` group. Create the role homes as `src/app/(app)/ops/page.tsx`, `src/app/(app)/issuer/page.tsx`, `src/app/(app)/acquirer/page.tsx`, `src/app/(app)/merchant/page.tsx`, `src/app/(app)/overview/page.tsx`. Cost if wrong: pages render without the nav shell.

## Steps

- **Step 1: nav component** — verbatim from plan (reference below). Uses `HOME_BY_ROLE`, `DASHBOARD_ACCESS`, `ROLE_LABEL`, `LogoutButton`, `ThemeToggle`, `Link`. TITLES map as-is.

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

- **Step 2: app shell layout** — CORRECTED (await + `(app)` group):

```tsx
// src/app/(app)/layout.tsx
import Nav from "@/components/nav";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { notFound } from "next/navigation";

export default async function AppLayout({ children }: { children: React.ReactNode }) {
  const supabase = await createServerClient();
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

- **Step 3: role home pages** — CORRECTED (await + role gate + placeholder content). Each page lives at `web/src/app/(app)/<home>/page.tsx`:

```tsx
// src/app/(app)/ops/page.tsx
import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";

export default async function OpsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
  return <h1 className="text-2xl font-semibold">Operations</h1>;
}
```

The same pattern ×5, each checking its own role:
- `(app)/ops` → `scheme_operator`
- `(app)/issuer` → `issuer`
- `(app)/acquirer` → `acquirer`
- `(app)/merchant` → `merchant`
- `(app)/overview` → no role check (viewers + any authenticated user may land here) — render `<h1>Overview</h1>` placeholder. (Task 8 replaces its content.)

- **Step 4: verify + commit**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd next build"
```

Expected: pass (login page may need its force-dynamic — already present from Task 5). Commit `feat(web): authenticated app shell with role nav`.

## Global Constraints (binding)

- Node v24.19.0; `npm.cmd`/`npx.cmd`. Windows PowerShell.
- Roles module is the single source of truth (`Role`, `ROLE_LABEL`, `HOME_BY_ROLE`, `DASHBOARD_ACCESS`, `roleFromAppMetadata`).
- `createServerClient()` is async — ALWAYS await it in server components.
- Unrelated working-tree files (`.github/**`, AGENTS.md, docs/28-*.md) belong to another process — NEVER stage/commit them.
- Do not create/switch branches; commit on `feat/web-ui`.
- nav.tsx uses `DASHBOARD_ACCESS[role] ?? []` per plan (roles module already exports it).

## Controller Rulings for this task

- D1, D2 above are THE corrections. Do not "simplify" further.
- The five role-home pages are placeholders this task; Tasks 8-9 fill their content. Do NOT build dashboards here.
- `href` values come from `DASHBOARD_ACCESS` (e.g. `/funding`, `/tokens` will 404 until Tasks 9 creates them) — acceptable on this branch; final review checks no dangling 404s remain.