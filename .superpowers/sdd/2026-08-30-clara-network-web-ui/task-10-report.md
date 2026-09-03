# Task 10 Report — seed scripts + demo users + CLI runbook

**Status:** DONE
**Commit:** `7b4413a` — `feat(web): seed scripts, cli runbook, demo users` (on `feat/web-ui`)
**Branch:** `feat/web-ui` (no branch created, per constraint)

## What I implemented

Files created / modified exactly per Steps 1–6 of the brief:

1. **`web/scripts/seed-docker.ps1`** (Windows-first, PowerShell 5.1):
   - Reads `CLARA_PG_DSN` (falls back to `DATABASE_URL`); prints usage + exits non-zero if missing.
   - Writes a temp Compose override (OS temp dir, UTF-8 **no BOM** via
     `[System.IO.File]::WriteAllText` + `UTF8Encoding($false)`), single-quote-escaping the DSN
     (`'` → `''`) defensively, for `switch`, `cardsvc`, `clearing-sim`, `ledger-sim`, `card-sim`,
     `acquiring-sim`, `disputes-sim`.
   - `docker compose -f <compose> -f <override> up -d --build redis switch cardsvc`.
   - Readiness: TCP connect loop to `localhost:18081` up to 60s; on timeout prints `cardsvc` logs
     and fails.
   - One-shot data sims run sequentially (`run --rm --no-deps`): clearing-sim, ledger-sim,
     card-sim, acquiring-sim, disputes-sim; then `acquirer-sim` to generate authorizations through
     the running switch.
   - Teardown: `docker compose stop postgres redis switch cardsvc` (no `down -v`); override file
     removed; cleanup in `finally` even on failure.
   - Prints `seed complete (target: host/db)` with the password stripped (regex removes
     `scheme://user:pass@`).
   - `$ErrorActionPreference = "Stop"`; every compose step checks `$LASTEXITCODE -eq 0` and aborts
     naming the failing service/step; PS 5.1 syntax only (no ternary / `??`).

2. **`web/scripts/seed.sh`** — bash twin (`set -euo pipefail`, `trap cleanup EXIT`): same DSN
   fallback + usage, same override services + YAML quoting, readiness via `bash` `/dev/tcp` to
   `localhost:18081` (pgrep-free), `run --rm --no-deps` one-shots, masked print via `sed`.

3. **`web/scripts/create-users.mjs`** — ESM, uses `import { createClient } from "@supabase/supabase-js"`
   (dual CJS/ESM, resolves fine). Reads `SUPABASE_URL` / `SUPABASE_SERVICE_ROLE_KEY`. Idempotent loop
   over the D1 manifest: `auth.admin.listUsers({ page, perPage: 200 })` page-scan by exact email →
   `updateUserById(id, { app_metadata: { role } })` if present, else `createUser({ email, password,
   email_confirm: true, app_metadata: { role } })`. Per-user created/updated/FAILED output; exits 1
   on any API error. Contains the required `// syntax check via: node --check scripts/create-users.mjs`.
   Client options `{ auth: { autoRefreshToken: false, persistSession: false } }` per brief.

4. **`web/sql/demo-profile.sql`** — SQL-comment-only documentation manifest (role → email → shared
   password → `app_metadata.role`), with the required "documentation ONLY / created via
   create-users.mjs with the service-role key" header.

5. **`web/README.md`** — replaced create-next-app boilerplate with the 7-section runbook:
   Overview, Prereqs, Local setup, Database + seed (`npx supabase db push` / apply `schema.sql`,
   `db:seed`, `db:users`), Demo matrix (role/email/password/landing/what's-visible), Deploy
   (CLI-only Vercel/Railway/Supabase + env-var table), Tear-down/reset.

6. **`web/package.json`** — added to the existing `scripts` (nothing removed):
   `"db:seed": "powershell -NoProfile -ExecutionPolicy Bypass -File scripts/seed-docker.ps1"`,
   `"db:users": "node --env-file=.env.local scripts/create-users.mjs"`.

## What I tested (exact commands + output)

- **PS parse:** `[System.Management.Automation.Language.Parser]::ParseFile("web\scripts\seed-docker.ps1")`
  → first pass found one error (`$svc:` in the here-string — `:` needs `${svc}`); **fixed** (now
  `  ${svc}:`) → **OK (637 tokens, 0 errors)**.
- **node --check:** `node --check web\scripts\create-users.mjs` → **OK**.
- **bash -n:** `& "C:\Program Files\Git\bin\bash.exe" -n web/scripts/seed.sh` → **OK**
  (WSL `C:\WINDOWS\system32\bash.exe` is broken on this box — "Class not registered" — so Git Bash
  was used; noted in concerns).
- **Compose override merge (live, safe):** generated the override with the script's exact logic,
  then `docker compose -f deploy\docker-compose.yml -f <override> config --services` → all services
  parse; `config` output confirms all 7 data services receive the external DSN and `postgres` has no
  host ports (`!reset []`); `adminapi` left untouched.
- **tsc:** `cmd /c "cd web && npx.cmd tsc --noEmit"` → **OK**.
- **next build:** `cmd /c "cd web && npx.cmd next build"` → **compiled successfully** (17 routes +
  Proxy; only pre-existing middleware-deprecation warning).
- **money regression:** `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"` → 1
  test, **pass**.
- **create-users.mjs guard + ESM resolution:** `node scripts/create-users.mjs` (no env) → prints
  the missing-keys message, **exit 1**, no network. Confirms the module import resolves under ESM.
- **Masking logic:** verified both PS and bash strip `scheme://user:pass@` correctly
  (e.g. `...secret%40pass@aws-0....amazonaws.com:5432/postgres?sslmode=require` →
  `aws-0....amazonaws.com:5432/postgres`); handles missing-`@` DSNs.
- **Manifest consistency:** grepped README / demo-profile.sql / create-users.mjs — all five
  `*@clara.demo` emails, shared password `ClaraDemo!2026`, and the 5 role slugs match exactly.

## Files changed (all committed)

- `web/scripts/seed-docker.ps1` (new)
- `web/scripts/seed.sh` (new)
- `web/scripts/create-users.mjs` (new)
- `web/sql/demo-profile.sql` (new)
- `web/README.md` (rewritten)
- `web/package.json` (2 scripts added)

No other working-tree files staged/committed (`.github/**`, AGENTS.md, docs/28-*.md, etc. untouched).

## Self-review findings

- **D1:** three files agree on the canonical identities; `create-users.mjs` creates via admin API
  (never raw SQL), is idempotent, and `demo-profile.sql` is pure documentation. ✓
- **D2:** all data sims + switch get the external DSN via the swap-free compose override; `adminapi`
  excluded. ✓
- **D3:** PowerShell 5.1 is primary, `seed.sh` is the twin. ✓
- PS 5.1-only syntax; no ternaries/`??`; `finally` cleanup removes the temp override and stops
  background services on both success and failure; one-shot failures abort with service name.
- `powershell -ExecutionPolicy Bypass` used in the npm script so the script runs under `-File`.

## Deviations / concerns

1. **`postgres: ports: !reset []` added to the override.** `up -d redis switch cardsvc` will pull in
   the compose-local `postgres` as a dependency (switch/cardsvc have `depends_on: condition:
   service_healthy`). Without the reset it would bind host **5432** and hang/conflict on any dev
   machine already running Postgres. Compose v2.10+ supports `!reset`; confirmed working with the
   installed Docker Compose **v5.4.0**. Local postgres thus runs idle (healthy) and is stopped at
   teardown.
2. **`--no-deps` added to the `run --rm` one-shots.** Without it, `docker compose run` starts
   dependencies — i.e. the local `postgres` — which is never needed since every sim points at the
   external DSN (and would re-risk the 5432 conflict). Real deps (switch/cardsvc/redis) are already
   up in the background under the same project/network names.
3. **Compose file ordering standardized to `-f <compose> -f <override>` everywhere** (override last
   so it wins the merge). The brief's Step 7 example listed override-first; that order only matters
   for config load, but I kept override-last consistently for correctness.
4. **Teardown also stops `postgres`** (it only came up as a dependency; leaving it running would
   otherwise be a surprise). One-shot containers are `--rm`, nothing left exited.
5. **WSL bash is broken on this machine** (`bash.exe -c "echo ok"` → "Class not registered"); used
   Git Bash for `bash -n`. On other machines WSL bash would also pass.
6. **README demo matrix** uses the real role → landing/visibility mappings from
   `web/src/lib/roles.ts` (HOME_BY_ROLE / DASHBOARD_ACCESS), not invented ones.
7. Scope kept: scripts were created + syntax/merge-checked only; **no live DB, no Supabase users,
   no Vercel/Railway provisioning** (Task 11). `seed-docker.ps1`/`seed.sh` were not executed against
   any environment (they require a target DB + a multi-minute docker build).

## Report path

`.superpowers/sdd/2026-08-30-clara-network-web-ui/task-10-report.md`