# Task 5 Report: login page + user avatar/logout

## What I implemented

- `web/src/lib/supabase/client.ts` — client-side Supabase factory using `createBrowserClient` from `@supabase/ssr`, reading statically-referenced `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` (not `@/lib/env`).
- `web/src/app/login/page.tsx` — server-rendered login route. Awaits async `createServerClient()`, redirects already-authenticated users to their role home (or `/`), renders `<LoginForm>` inside `<Suspense>`. Added `export const dynamic = "force-dynamic"` so the page is NOT statically prerendered (see concerns).
- `web/src/app/login/login-form.tsx` — client form using Card/Label/Input/Button, `signInWithPassword`, role-based redirect honoring `?next=`.
- `web/src/components/logout-button.tsx` — client button calling `signOut()`, redirecting to `/login` + `router.refresh()`.
- `web/.env.local.example` — added `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` placeholders.
- Added shadcn UI components `card`, `input`, `label` (only `button` existed). Note `@base-ui/react` was already a dependency, so no new deps.

## What I tested

- `cmd /c "cd web && npx.cmd tsc --noEmit"` → passed (no output).
- `cmd /c "cd web && npx.cmd next build"` → passed. `/login` shows as `ƒ` (dynamic server-rendered). (First build failed — see concerns — fixed with `force-dynamic`.)
- `cmd /c "cd web && npx.cmd eslint src/app/login src/components/logout-button.tsx src/lib/supabase/client.ts"` → exit 0, no errors.

## Files changed

- `web/src/lib/supabase/client.ts` (new)
- `web/src/app/login/page.tsx` (new)
- `web/src/app/login/login-form.tsx` (new)
- `web/src/components/logout-button.tsx` (new)
- `web/src/components/ui/card.tsx` (new, shadcn)
- `web/src/components/ui/input.tsx` (new, shadcn)
- `web/src/components/ui/label.tsx` (new, shadcn)
- `web/.env.local.example` (modified)

Commit: `d70c688` feat(web): login page, logout, demo role flows (8 files, +241)

## Self-review findings

- **Deviated from verbatim code (fix):** The brief's `client.ts` defined `export function createBrowserClient()` that imported `createBrowserClient` from `@supabase/ssr` under the same name — the local function shadows the import, causing infinite self-recursion at runtime. I aliased the import (`createSupabaseBrowserClient`) to fix, matching the existing `server.ts` convention.
- **Added `export const dynamic = "force-dynamic"`** to the login page (beyond the brief's verbatim code) — required for `next build` to pass without env vars.
- Verified `--color-destructive` is defined in globals.css so `text-destructive` renders.
- Confirmed only the task's `web/` files were committed; unrelated `.github/`, docs, AGENTS.md, etc. were left unstaged.

## Issues / concerns

1. **Potential middleware redirect loop on `/login` (out of scope, not modified):** `src/lib/supabase/middleware.ts` (Task 4) redirects unauthenticated requests to `/login?next=<path>`. Since it matches `/login` too, a logged-out user hitting `/login` would be redirected to `/login?next=/login` → same URL → possible redirect loop. I did not modify Task 4's middleware (out of my task's declared file scope). Recommend the controller decide whether `/login` (and static/public routes) should be excluded from the auth matcher.
2. **No real Supabase creds** — runtime behavior (sign in / sign out / session routing) is untested; only build/typecheck/lint verified.
3. Pre-existing `Next.js middleware deprecated` build warning (from Task 4's `middleware.ts`) is unrelated to this change.
