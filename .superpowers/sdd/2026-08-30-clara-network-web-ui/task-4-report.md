# Task 4 Report: Supabase client (server + middleware) + roles

**Status:** DONE_WITH_CONCERNS
**Commit:** `7149a29` feat(web): supabase auth clients, roles, auth middleware (on `feat/web-ui`)

## What I implemented

Per the brief: Supabase server client, middleware session guard, role model, root `middleware.ts`, and the two `@supabase/*` deps.

- `web/src/lib/supabase/server.ts` — per-request `createServerClient()` using `getEnv()` + `@supabase/ssr`.
- `web/src/lib/supabase/middleware.ts` — `updateSession(request)`; refreshes cookies, redirects unauthenticated users to `/login?next=...`.
- `web/src/lib/roles.ts` — `Role` union, `ROLE_LABEL`, `roleFromAppMetadata`, `HOME_BY_ROLE`, `DASHBOARD_ACCESS`, `dashboardAccess`.
- `web/src/middleware.ts` — Next middleware wrapper with the brief's matcher.
- `web/package.json` / `web/package-lock.json` — `@supabase/ssr@^0.12.5`, `@supabase/supabase-js@^2.112.4`.

Note: all four source files and both dep entries were already present on disk uncommitted from a prior partial run. I verified them against the brief, fixed the build failures (below), and committed.

## Deviations from the brief (required to build)

1. **`server.ts` name collision (TS2440/TS2554).** The verbatim brief exports `createServerClient()` while importing `createServerClient` from `@supabase/ssr`; the call inside the function resolves to itself (0 args) instead of the Supabase factory. Fixed by aliasing the import: `createServerClient as createSupabaseServerClient`. Exported interface (`createServerClient()`) is unchanged.
2. **`server.ts` async cookies (TS2339).** Next 16's `cookies()` is async (`Promise<ReadonlyRequestCookies>`), so `cookieStore.getAll()` fails. Fixed with `const cookieStore = await cookies();` and made the function `async`. Consumers in later tasks must `await createServerClient()`.
3. **eslint `@typescript-eslint/no-require-imports`.** eslint-config-next/typescript forbids `require()`. The controller ruling requires keeping the `require("../env")` shape, so I added `// eslint-disable-next-line @typescript-eslint/no-require-imports` to the two `require` lines rather than converting to a static import or disabling the rule globally.

Everything else (cookie handling in `middleware.ts`, `setAll: () => {}` no-op, `/login?next=...` redirect, matcher, roles module verbatim) matches the brief exactly.

## What I tested

- `cmd /c "cd web && npx.cmd next build"` — **passes.** Compiled + TypeScript + static generation OK. Warning: `middleware` convention is deprecated in Next 16 in favor of `proxy` (kept `middleware.ts` per brief; not blocking).
- `cmd /c "cd web && npx.cmd tsc --noEmit"` — **passes** (clean, no output).
- `cmd /c "cd web && npx.cmd eslint src/lib/supabase src/lib/roles.ts src/middleware.ts"` — **passes** after the two eslint-disable directives.
- Full `next build` run before the server.ts fix **failed** with TS2440/TS2554/TS2339, confirming the fix was necessary.

## Files changed (committed)

- `web/package.json`, `web/package-lock.json`
- `web/src/lib/roles.ts` (new)
- `web/src/lib/supabase/server.ts` (new)
- `web/src/lib/supabase/middleware.ts` (new)
- `web/src/middleware.ts` (new)

Unrelated working-tree changes (`.github/ci.yml` modified, other untracked files) were intentionally **not** committed.

## Self-review findings

- Interface contract produced: `createServerClient()` per-request with cookie methods (now async — see above), `updateSession(request)`, `Role`, `roleFromAppMetadata`, `ROLE_LABEL`, `HOME_BY_ROLE`, `dashboardAccess`. All present.
- Roles read from `app_metadata.role` (no role table); anon key only, no secrets exposed. Matches global constraints.
- `setAll: () => {}` in server client kept intentionally (server components are read-only; middleware owns refreshes).
- Path alias `@/lib/supabase/middleware` used in `src/middleware.ts`.

## Issues / concerns

1. **Brief's verbatim `server.ts` does not compile** under this stack (Next 16 async `cookies()` + import/local name collision). The brief's "Expected: pass" only holds with the two fixes above. If the controller prefers the exact verbatim shape, the alternative is a Next 15 pin (`next@15`) — not recommended.
2. **`createServerClient()` is now async** — Task 5/6 consumers must `await` it. This matches the official Supabase Next.js pattern.
3. **Next 16 deprecation warning**: `src/middleware.ts` should ideally be `src/proxy.ts` at some point (Next 16 prefers the `proxy` convention). Kept `middleware` per brief; flagging so a later task can decide.
4. `web/.env.local` is absent — `getEnv()` will throw at runtime for any page that calls it. Expectation confirmed: build/typecheck unaffected.