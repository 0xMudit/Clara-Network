# Clara Network Web Console

## 1. Overview

The Clara Network web console is a read-only admin dashboard over the Clara
`adminapi`. A Next.js (App Router) BFF proxies the Clara admin API to the
browser behind a Supabase Auth login, and gates every route by the signed-in
user's `app_metadata.role` (`scheme_operator | issuer | acquirer | merchant |
viewer`). Hosting is CLI-only: Next.js on Vercel, PostgreSQL + Auth on Supabase,
and the Go `adminapi` on Railway.

## 2. Prereqs

- Node.js 24 (requires 20.6+ for `--env-file`, 24 is tested)
- Docker Desktop (Windows) or Docker Engine (Linux/macOS) — the Go sims run in
  containers (a Go image build avoids Windows AppLocker issues with local
  binaries)
- Accounts: Vercel (frontend), Supabase (DB + Auth), Railway (Go `adminapi`)
- A terminal. PowerShell 5.1 on Windows; bash on Linux/macOS

## 3. Local setup

```sh
cd web
copy .env.local.example .env.local     # Windows
# cp .env.local.example .env.local     # on Linux/macOS
```

Fill in `.env.local`:
- `SUPABASE_URL` / `SUPABASE_ANON_KEY` — Supabase project → Settings → API
- `CLARA_API_URL` — http://localhost:18083 during local dev, else the Railway URL
- `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` — same values as
  the non-public pair (Next requires the `NEXT_PUBLIC_` prefix)

Then:

```sh
npm install
npm run dev        # http://localhost:3000
```

For demo-user provisioning also add `SUPABASE_SERVICE_ROLE_KEY`
(Supabase project → Settings → API, `service_role` key) to `.env.local`.

## 4. Database + seed

1. **Push the schema** to the target Supabase/Postgres DB:

   ```sh
   npx supabase db push --project-ref <ref> --include-all
   # or apply deploy/schema.sql directly to any Postgres with:
   #   psql "<SUPABASE_POSTGRES_DSN>" -f deploy/schema.sql
   ```

2. **Seed demo data** with the one-shot data sims, pointed at the external DSN
   (`CLARA_PG_DSN`, falling back to `DATABASE_URL`; the compose override never
   touches the compose-local `postgres`):

   ```sh
   # Windows (primary):
   $env:CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require"
   npm run db:seed          # -> powershell scripts/seed-docker.ps1
   # Linux/macOS:
   CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" \
     npm run db:seed
   # or directly:  ./scripts/seed.sh
   ```

   What it runs: `switch` + `redis` + `cardsvc` in the background, then
   `clearing-sim` (clearing/settlement/prefund/default fund), `ledger-sim`
   (ledger accounts + entries), `card-sim` (cards + tokens), `acquiring-sim`
   (merchants + funding), `disputes-sim` (disputes), and finally `acquirer-sim`
   to drive authorizations that write `switch_transactions` through the running
   switch. The `adminapi` is not needed to seed.

3. **Create the demo users** in Supabase Auth (service-role admin API,
   idempotent — updates roles, creates what's missing):

   ```sh
   npm run db:users         # node --env-file=.env.local scripts/create-users.mjs
   ```

## 5. Demo matrix

| role            | email                     | password        | landing page | what's visible                                                    |
| --------------- | ------------------------- | --------------- | ------------ | ----------------------------------------------------------------- |
| scheme_operator | scheme-operator@clara.demo| `ClaraDemo!2026`| `/ops`       | ops, transactions, clearing, settlement, ledger, cards, merchants, disputes |
| issuer          | issuer@clara.demo         | `ClaraDemo!2026`| `/issuer`    | issuer, cards, tokens                                             |
| acquirer        | acquirer@clara.demo       | `ClaraDemo!2026`| `/acquirer`  | acquirer, merchants, funding, disputes                             |
| merchant        | merchant@clara.demo       | `ClaraDemo!2026`| `/merchant`  | merchant, funding, disputes                                        |
| viewer          | viewer@clara.demo         | `ClaraDemo!2026`| `/overview`  | overview                                                           |

The canonical identities live in `scripts/create-users.mjs` and
`sql/demo-profile.sql` (documentation manifest) — keep all three in sync.

## 6. Deploy (CLI-only)

```sh
# Vercel (frontend), from web/:
cd web
npx vercel --prod

# Railway (Go adminapi), from the repo root:
railway up --deploy

# Supabase (schema), from the repo root:
npx supabase db push --project-ref <ref> --include-all
```

Environment variables to set on the deployed services:

| var                           | service  |
| ----------------------------- | -------- |
| `SUPABASE_URL`                | Vercel   |
| `SUPABASE_ANON_KEY`           | Vercel   |
| `CLARA_API_URL`               | Vercel   |
| `NEXT_PUBLIC_SUPABASE_URL`    | Vercel   |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY`| Vercel  |
| `SUPABASE_SERVICE_ROLE_KEY`   | Vercel (for `db:users`) |
| `CLARA_PG_DSN`                | Railway  |
| `CLARA_LISTEN` (=`:8083`)     | Railway  |

## 7. Tear-down / reset

```sh
# Local sims / data containers (add -v to wipe local volumes):
docker compose -f deploy/docker-compose.yml down -v

# Reset a local dev database's seed data: re-run section 4 steps 1-2.

# Cloud:
#   Supabase: dashboard -> Project Settings -> Danger Zone -> Delete project
#   Vercel:   npx vercel remove <project> --yes
#   Railway:  railway down (or delete the project in the dashboard)
```