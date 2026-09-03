# Task 6 Report: BFF data proxy + route handler

## Status: DONE

## What I implemented

1. **`web/src/lib/clara.ts`** — server-only typed fetch helper.
   - `claraFetch<T>(path, accessToken, init?)` reads `CLARA_API_URL` from `getEnv()` and
     fetches `${CLARA_API_URL}${path}` server-side with
     `headers: { Authorization: "Bearer <token>", Accept: "application/json" }`.
   - Corrected the brief's **intentional syntax error** (`headers: auth: {...}`) to a proper
     headers object. Also dropped the doubled `Bearer Bearer` prefix to a single `Bearer`
     (that was part of the erroneous line).
   - Throws `Error("clara <path>: <status>")` on non-2xx so the route handler can map it.
   - `import "server-only"` guards against accidental client bundling.

2. **`web/src/app/api/data/[...path]/route.ts`** — read-only BFF proxy.
   - `export const dynamic = "force-dynamic"` (runtime env lookup; build passes without
     `web/.env.local`).
   - `GET` awaits `createServerClient()` (async) and `await params` (Next 16 route-handler
     contract) for `Promise<{ path: string[] }>`.
   - Rejects unauthenticated requests with 401 `{ error: "unauthorized" }`.
   - Role gate via `app_metadata.role` against `ALLOWED_ROUTES` (scheme_operator/issuer/
     acquirer/merchant); unknown roles default to `[]` → 403 `{ error: "forbidden" }`.
     No admin API reachable beyond the allowlist.
   - Proxies via `claraFetch(target, data.user.id, { cache: "no-store" })`. Per controller
     ruling: **no Redis/DB cache** — correctness first, caching deferred.
   - Upstream 4xx → 502, everything else → 500 (`{ error: "upstream error", detail }`).
   - No secrets logged; no direct browser→adminAPI path (all traffic goes through BFF).

3. **Dependency** — added `server-only@^0.0.1` (installed successfully with
   `set npm_config_allow_scripts=` prefix; no EALLOWSCRIPTS issue encountered).

## What I tested (exact commands + output summaries)

- `cmd /c "cd web && npx.cmd tsc --noEmit"` → exit 0, no diagnostics.
- `cmd /c "cd web && npx.cmd eslint src/app/api src/lib/clara.ts"` → exit 0, 0 errors,
  1 warning: unused `getEnv` import in route.ts (verbatim from the brief's code block;
  `no-unused-vars` is warn-only in eslint-config-next).
- `cmd /c "cd web && npx.cmd next build"` → exit 0. Compiled successfully; route table shows
  `/api/data/[...path]` registered as `ƒ (Dynamic)`. Only pre-existing warning is the
  middleware→proxy deprecation (Task 5's file, untouched).

## Files changed

- `web/src/lib/clara.ts` (new)
- `web/src/app/api/data/[...path]/route.ts` (new)
- `web/package.json` (+`server-only`)
- `web/package-lock.json` (+`server-only`)

## Commit

- `42ee415` feat(web): BFF data proxy, role-gated admin API access (on `feat/web-ui`).
  Staged ONLY the 4 task files; `.github/`, `AGENTS.md`, `.golangci.yml`, etc. left untouched
  for the other process.

## Self-review findings

- `claraFetch` spreads `...init` after `headers` **is** correct per brief (caller overrides
  win); `cache: "no-store"` flows through. The brief's `as RequestInit` cast was unnecessary
  once the syntax was fixed — omitted.
- Error message format `clara ${path}: ${status}` matches the route's
  `/^clara \S+ 4\d\d/` upstream-4xx→502 regex.
- `req` first arg intentionally kept (Next route-handler signature); skipped by
  `no-unused-vars` (`args: "after-used"`).
- Role `undefined`/unknown → `[]` allowlist → 403: safe default; the `app_metadata.role`
  string is only used as an index key, never reflected in output.
- `data.user.id` reused as the admin-API token per spec — acceptable since the admin API has
  no auth and the role gate is the sole authorization layer.

## Issues / concerns

- None blocking. Minor: the brief's route.ts imports `getEnv` but never uses it → 1 eslint
  warning (kept verbatim to match spec).
- Next logs a deprecation warning for `middleware`→`proxy` (pre-existing from Task 5; out of
  scope here).