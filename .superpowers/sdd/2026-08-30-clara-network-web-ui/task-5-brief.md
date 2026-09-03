# Task 5 Brief: login page + user avatar/logout

## Task Description (from plan)

Login/authentication UI: a server-rendered login page with a client form, plus a logout button.

**Files:**
- Create: `web/src/app/login/page.tsx`, `web/src/app/login/login-form.tsx`, `web/src/components/logout-button.tsx`

**Interfaces:**
- Consumes: Task 4 `createServerClient()` (NOW ASYNC — see Ruling), `roleFromAppMetadata`, `HOME_BY_ROLE`.
- Produces: route `/login`; login form signs in via `signInWithPassword`, redirects to `next` or role home; `LogoutButton` component.

## Steps

- **Step 1: login page** (server component — note `createServerClient` is async per Ruling):

```tsx
// src/app/login/page.tsx
import { Suspense } from "react";
import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { LoginForm } from "./login-form";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export default async function LoginPage() {
  const supabase = await createServerClient();
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

- **Step 2: login form** (client component). IMPORTANT — RULING: this is a CLIENT component; it CANNOT use `getEnv()` from `@/lib/env` (browser cannot read arbitrary `process.env`). Instead use the public Supabase constants read through statically-referenced `NEXT_PUBLIC_*` env vars. Create `web/src/lib/supabase/client.ts` first:

```ts
// src/lib/supabase/client.ts
"use client";
import { createBrowserClient } from "@supabase/ssr";

export function createBrowserClient() {
  return createBrowserClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
  );
}
```

Then the form:

```tsx
// src/app/login/login-form.tsx
"use client";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
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
    const supabase = createBrowserClient();
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

- **Step 3: logout button**

```tsx
// src/components/logout-button.tsx
"use client";
import { useRouter } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";

export function LogoutButton() {
  const router = useRouter();
  async function logout() {
    const supabase = createBrowserClient();
    await supabase.auth.signOut();
    router.replace("/login");
    router.refresh();
  }
  return <Button variant="outline" onClick={logout}>Sign out</Button>;
}
```

- **Step 4: NEXT_PUBLIC env wiring**

Add to `web/.env.local.example` (and getEnv is server-only — NEXT_PUBLIC vars are read directly in the client module, NOT through getEnv):
```
# Client-side public Supabase vars (next.js requires NEXT_PUBLIC_ prefix, values match SUPABASE_URL/SUPABASE_ANON_KEY)
NEXT_PUBLIC_SUPABASE_URL=https://<ref>.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=<anon key>
```

- **Step 5: demo users** — created in Task 10 (script). Skip here.

- **Step 6: verify + commit**

```bash
cmd /c "cd web && npx.cmd next build"
cmd /c "cd web && npx.cmd tsc --noEmit"
```

Expected: pass. Commit `feat(web): login page, logout, demo role flows`.

## Global Constraints (binding)

- Node v24.19.0; npm via `npm.cmd`. Windows PowerShell.
- Roles from `app_metadata.role`; anon key is public by design.
- Never log/secrets; `.env.local.example` placeholders only.
- Client components must NOT import `@/lib/env` (server-only). Use `@/lib/supabase/client`.

## Controller Rulings for this task

- **RULING (carried from Task 4):** `createServerClient()` in `@/lib/supabase/server` is async (Next 16 async cookies). **AWAIT it** in the login page. Also the supabase-js client module is `createBrowserClient()` (new file `src/lib/supabase/client.ts`) — Do NOT reuse server client in client components.
- RULING (plan defect fix): `getEnv()` cannot be used in client components. The plan's `login-form`/`logout-button` used `await import("@/lib/env")...` — replaced by statically-referenced `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` via a dedicated client factory. Cost if wrong: two extra placeholder env vars — correct-by-default.
- Use existing UI components (Card/Label/Input/Button) that shadcn created — if `label.tsx` / `input.tsx` are missing, add them via `npx.cmd shadcn@latest add label input` (v4 CLI: `-b base -p nova` flags as needed). Add only what's required.