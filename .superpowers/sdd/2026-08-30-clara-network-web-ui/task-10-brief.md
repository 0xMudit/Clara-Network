# Task 10 Brief: seed scripts + demo users + CLI runbook

## Task Description (from plan, with controller corrections)

Tooling + runbook to populate the target DB (Supabase/Postgres) with demo data and demo users, and to drive auth via terminal-only CLIs. THIS TASK ONLY CREATES FILES — live Supabase/Vercel/Railway provisioning happens in Task 11 (needs the user's accounts). Do not attempt real logins here.

**Files:**
- Create: `web/scripts/seed-docker.ps1`, `web/scripts/seed.sh`, `web/scripts/create-users.mjs`, `web/sql/demo-profile.sql`, `web/README.md`
- Modify: `web/package.json` (add `db:seed`, `db:users` scripts)

## PLAN DEFECTS / CONTROLLER RULINGS (binding)

- **D1 (demo users via admin API, not raw SQL):** `demo-profile.sql` cannot insert into Supabase's `auth.users` (service-managed). Create it as a documented MANIFEST (role / email / password / app_metadata) in SQL-comment form for reference — the ACTUAL creation is `create-users.mjs` using `supabase.auth.admin.createUser(...)` with the service-role key. Manifest (single shared demo password, all roles, `email_confirm: true`):
  - scheme_operator → `scheme-operator@clara.demo`, password `ClaraDemo!2026`
  - issuer → `issuer@clara.demo`, `ClaraDemo!2026`
  - acquirer → `acquirer@clara.demo`, `ClaraDemo!2026`
  - merchant → `merchant@clara.demo`, `ClaraDemo!2026`
  - viewer → `viewer@clara.demo`, `ClaraDemo!2026`
  All get `app_metadata: { role: "<roleslug>" }`. The script is IDEMPOTENT: look up by email; if present, update app_metadata to ensure the role is set; else create.
- **D2 (seed data plan):** `deploy/docker-compose.yml` (shown in repo) has the data sims. `switch_transactions` come from `switch` + `acquirer-sim` (needs `redis`); cards/tokens from `cardsvc`+`card-sim`; merchants/funding from `acquiring-sim`; clearing/settlement/prefund/default-fund from `clearing-sim`; ledger from `ledger-sim`; disputes from `disputes-sim`. `adminapi` is NOT required for seeding data. The script must point the sims' `CLARA_PG_DSN` at the EXTERNAL target DB (an env-supplied DSN), NOT the compose-local `postgres`.
- **D3 (Windows-first):** `seed-docker.ps1` is the primary script (this is a Windows dev machine). `seed.sh` is a thin Linux/macOS equivalent for CI/other developers. Both must override `CLARA_PG_DSN` for every data service.

## Step 1 — seed-docker.ps1 (Windows PowerShell 5.1)

Behavior:
1. Read `CLARA_PG_DSN` (or fall back to `DATABASE_URL`) from the environment; if missing, print usage (`$env:CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" .\seed-docker.ps1`) and exit non-zero.
2. Generate a temporary Compose OVERRIDE file in the OS temp dir with the DSN baked in, providing `environment: { CLARA_PG_DSN: "<dsn>" }` for services: `switch`, `cardsvc`, `clearing-sim`, `ledger-sim`, `card-sim`, `acquiring-sim`, `disputes-sim`. Write it with UTF-8 (no BOM), quoting the DSN defensively.
3. Start the servers needed for data generation in the background against the external DB:
   - `docker compose -f <repo>/deploy/docker-compose.yml -f <override> up -d --build redis switch cardsvc`
   - (switch needs redis+issuer-sim services from the stack; the compose `depends_on` handles ordering. issuer-sim doesn't touch the DB, fine.)
4. Wait for `cardsvc` to be ready: loop `docker compose ... ps cardsvc --format json` / try a TCP connect to `localhost:18081` up to ~60s; if never ready, print an error and exit.
5. Run the one-shot data sims sequentially (each `docker compose -f ... -f <override> run --rm <svc>`), building once on first use:
   - `clearing-sim` (clearing_records, net_positions, settlement_instructions, prefund_accounts, default_fund)
   - `ledger-sim` (ledger_accounts, ledger_entries)
   - `card-sim` (cards + tokens via cardsvc)
   - `acquiring-sim` (merchants, funding lines)
   - `disputes-sim` (disputes)
6. Generate authorizations: `docker compose -f ... -f <override> run --rm acquirer-sim` (writes `switch_transactions` through the running switch).
7. Tear down the background servers: `docker compose -f <override-file> -f <compose> stop redis switch cardsvc` (leave one-shot containers exited; do NOT `down -v`). Remove the override file.
8. Print `seed complete` and the DSN host (masked: show only host/db, never the password).

Notes: use `$ErrorActionPreference = "Stop"`; log each step; verify each one-shot exits `$LASTEXITCODE -eq 0` else stop with that service's name in the error. Use only PowerShell 5.1-compatible syntax (no ternary, no `??`). Validate the file with `[System.Management.Automation.Language.Parser]::ParseFile()` — report parse errors.

## Step 2 — seed.sh (bash, Linux)

Same pipeline, bash (set -euo pipefail). `docker compose -f ... up -d --build redis switch cardsvc`, readiness via `docker compose exec -T cardsvc wget -q --spider http://localhost:8081/health` style or `pgrep`-free; one-shot `docker compose run --rm` runs; teardown; masked-print. Use `DATABASE_URL`/`CLARA_PG_DSN` the same way. (You won't be able to execute it on this Windows box — syntax-check only via `bash -n`.)

## Step 3 — create-users.mjs (Node ≥ 20, ESM)

- Reads `SUPABASE_URL` and `SUPABASE_SERVICE_ROLE_KEY` from process.env (runner passes `--env-file=.env.local`).
- Uses `@supabase/supabase-js` `createClient(SUPABASE_URL, SERVICE_ROLE_KEY, { auth: { autoRefreshToken: false, persistSession: false } })`.
- Manifest table (D1 emails/passwords/roles). Loop: `supabase.auth.admin.listUsers()` paginated OR `getUserById` search by exact `email`; if user exists → `auth.admin.updateUserById(id, { app_metadata: { role } })`; else → `auth.admin.createUser({ email, password, email_confirm: true, app_metadata: { role } })`.
- Prints per-user result (created / updated) and exits non-zero on any API error. Idempotent by design.
- Add `// syntax check via: node --check scripts/create-users.mjs`.
- If `@supabase/supabase-js` default import shape causes issues under ESM, use `import { createClient } from "@supabase/supabase-js"` (works — it ships dual CJS/ESM). Validate `node --check` passes.

## Step 4 — demo-profile.sql

SQL-comment manifest documenting the five demo users (role → email → shared password) and the exact `app_metadata.role` each maps to, plus a header: "This file is documentation ONLY. Demo users are created via `scripts/create-users.mjs` with the service-role key (managed `auth.users` cannot be inserted via SQL)."

## Step 5 — web/README.md runbook

Sections:
1. **Overview** — one paragraph: read-only admin console over the Clara admin API via a role-gated Next.js BFF (Vercel + Supabase + Railway).
2. **Prereqs** — Node 24, Docker Desktop, Go image build (via Docker), accounts: Vercel/Supabase/Railway.
3. **Local setup** — copy `.env.local.example` → `.env.local`, fill SUPABASE_URL/ANON_KEY/CLARA_API_URL; `npm install`; `npm run dev`.
4. **Database + seed** — push `deploy/schema.sql` (see Step 6 snippet), then run `db:seed` (seed-docker.ps1 / seed.sh) against the target DSN, then `db:users` (create-users.mjs) for demo logins.
5. **Demo matrix** — table: role | email | password | landing page | what's visible.
6. **Deploy (CLI-only)** — `npx vercel --prod` push from `web/`; Railway `railway up --deploy` from repo root; Supabase `npx supabase db push --project-ref <ref> --include-all`; env var list (SUPABASE_URL, SUPABASE_ANON_KEY, CLARA_API_URL, NEXT_PUBLIC_* pair).
7. **Tear-down / reset** — `docker compose -f deploy/docker-compose.yml down` (with `-v` to wipe local volumes), Supabase project delete, Vercel project remove.

## Step 6 — package.json scripts

Add:
```json
"scripts": {
  "db:seed": "powershell -NoProfile -ExecutionPolicy Bypass -File scripts/seed-docker.ps1",
  "db:users": "node --env-file=.env.local scripts/create-users.mjs"
}
```
(Add these two keys to the EXISTING scripts object — do not replace existing scripts.)

## Step 7 — verify + commit

Non-requiring any live services:
- `cmd /c "cd web && npx.cmd tsc --noEmit"`
- `cmd /c "cd web && npx.cmd next build"`
- `node --check web/scripts/create-users.mjs`
- PS parse check on `seed-docker.ps1` (see Step 1) and `bash -n` on `seed.sh` if bash exists on PATH (else note it)
- `cmd /c "cd web && npx.cmd tsx --test src/lib/money/minor.test.ts"` (regression)

Commit `feat(web): seed scripts, cli runbook, demo users`.

## Global Constraints (binding)

- PowerShell 5.1-compatible script only; ALWAYS `npm.cmd`/`npx.cmd`; Windows-first.
- Never hardcode secrets: DSN comes from env; demo password documented (it's a demo, not a secret).
- Do NOT execute the scripts against any live Supabase/Railway/Vercel (Task 11 does that) — create + syntax-check only.
- Unrelated working-tree files (`.github/**`, AGENTS.md, docs/28-*.md) — never stage/commit.
- No branch creation; commit on `feat/web-ui`.

## Controller Rulings

- D1 manifest emails/password are THE canonical demo identities (README + demo-profile.sql + create-users.mjs must all agree).
- The compose override approach (temp file with DSN baked in) is REQUIRED — do not `sed` the original compose, do not edit `deploy/docker-compose.yml`.