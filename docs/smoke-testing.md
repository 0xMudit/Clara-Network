# Smoke Testing

`make smoke` (or `scripts/smoke.sh`; `scripts/smoke.ps1` is the Windows twin)
is the one-click end-to-end smoke test for Clara Network. It boots the full
local docker-compose stack, seeds demo data, and verifies all three tiers —
database, backend API, and frontend — against **live services**, then prints a
PASS / FAIL table. Exit code is non-zero if any check fails (so it can gate CI).

## What it verifies

1. **DB suite** — connects to the compose `postgres` via `psql` and asserts
   every schema table exists and each domain has seeded rows
   (`switch_transactions`, `clearing_records`, `net_positions`,
   `settlement_instructions`, `ledger_entries`, `cards`, `tokens`, `merchants`,
   `funding_lines`, `screening_lists`, `disputes`, …).

2. **Backend API suite** — hits every `adminapi` endpoint and asserts
   `HTTP 200` + valid JSON:
   `/health`, `/dashboard`, `/transactions`, `/clearing/*`,
   `/settlement/*`, `/ledger/*`, `/cards`, `/bin-ranges`, `/tokens`,
   `/merchants`, `/disputes/*`. Cycle-scoped endpoints
   (`clearing/records`, `clearing/positions`, `settlement/instructions`)
   resolve a real `?cycle=` from `/clearing/cycles` first, mirroring the frontend.

3. **Frontend suite** — `next build` (proves every route compiles), then boots
   the production server and probes runtime behaviour for an unauthenticated
   user: `/login` → 200, `/` and protected pages → 307 redirect to `/login`,
   and the BFF proxy `/api/data/*` → **401 JSON** (not a 307 redirect).

## Usage

```sh
make smoke            # full run: boot, seed, all suites, report
scripts/smoke.sh --no-seed        # reuse an already-running, seeded stack
scripts/smoke.sh --keep           # leave the stack running when done
scripts/smoke.sh --skip-frontend  # skip the (slow) Next.js build + probes
```

On Windows, run `scripts/smoke.sh` from **Git Bash** (no host Go toolchain is
needed — the Go services run in Docker). `scripts/smoke.ps1` is the native
PowerShell equivalent.

## Gotchas it catches

- **Middleware hijacking `/api/*`**: the auth middleware (`web/src/middleware.ts`)
  must exclude `/api` so the BFF route returns `401/403 JSON` instead of
  redirecting API calls to `/login`. The frontend suite fails loudly if this
  regresses (expects `401`, flags a `307` as a bug).
- **Cycle-scoped endpoints** returning `400` when the frontend would send a
  `?cycle=` — the suite passes a real cycle so a `400` here means a real defect.
- **Missing seed data** after a schema/seed change.
