BASE 15397cc

7b4413a feat(web): seed scripts, cli runbook, demo users

=== STAT ===
 web/README.md                | 143 ++++++++++++++++++++++++++++++++++++-------
 web/package.json             |   4 +-
 web/scripts/create-users.mjs |  78 +++++++++++++++++++++++
 web/scripts/seed-docker.ps1  | 135 ++++++++++++++++++++++++++++++++++++++++
 web/scripts/seed.sh          |  83 +++++++++++++++++++++++++
 web/sql/demo-profile.sql     |  24 ++++++++
 6 files changed, 444 insertions(+), 23 deletions(-)
diff --git a/web/README.md b/web/README.md
index e215bc4..6928b6c 100644
--- a/web/README.md
+++ b/web/README.md
@@ -1,36 +1,135 @@
-This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).
+# Clara Network Web Console
 
-## Getting Started
+## 1. Overview
 
-First, run the development server:
+The Clara Network web console is a read-only admin dashboard over the Clara
+`adminapi`. A Next.js (App Router) BFF proxies the Clara admin API to the
+browser behind a Supabase Auth login, and gates every route by the signed-in
+user's `app_metadata.role` (`scheme_operator | issuer | acquirer | merchant |
+viewer`). Hosting is CLI-only: Next.js on Vercel, PostgreSQL + Auth on Supabase,
+and the Go `adminapi` on Railway.
 
-```bash
-npm run dev
-# or
-yarn dev
-# or
-pnpm dev
-# or
-bun dev
+## 2. Prereqs
+
+- Node.js 24 (requires 20.6+ for `--env-file`, 24 is tested)
+- Docker Desktop (Windows) or Docker Engine (Linux/macOS) ÔÇö the Go sims run in
+  containers (a Go image build avoids Windows AppLocker issues with local
+  binaries)
+- Accounts: Vercel (frontend), Supabase (DB + Auth), Railway (Go `adminapi`)
+- A terminal. PowerShell 5.1 on Windows; bash on Linux/macOS
+
+## 3. Local setup
+
+```sh
+cd web
+copy .env.local.example .env.local     # Windows
+# cp .env.local.example .env.local     # on Linux/macOS
+```
+
+Fill in `.env.local`:
+- `SUPABASE_URL` / `SUPABASE_ANON_KEY` ÔÇö Supabase project ÔåÆ Settings ÔåÆ API
+- `CLARA_API_URL` ÔÇö http://localhost:18083 during local dev, else the Railway URL
+- `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` ÔÇö same values as
+  the non-public pair (Next requires the `NEXT_PUBLIC_` prefix)
+
+Then:
+
+```sh
+npm install
+npm run dev        # http://localhost:3000
 ```
 
-Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.
+For demo-user provisioning also add `SUPABASE_SERVICE_ROLE_KEY`
+(Supabase project ÔåÆ Settings ÔåÆ API, `service_role` key) to `.env.local`.
+
+## 4. Database + seed
+
+1. **Push the schema** to the target Supabase/Postgres DB:
+
+   ```sh
+   npx supabase db push --project-ref <ref> --include-all
+   # or apply deploy/schema.sql directly to any Postgres with:
+   #   psql "<SUPABASE_POSTGRES_DSN>" -f deploy/schema.sql
+   ```
 
-You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.
+2. **Seed demo data** with the one-shot data sims, pointed at the external DSN
+   (`CLARA_PG_DSN`, falling back to `DATABASE_URL`; the compose override never
+   touches the compose-local `postgres`):
 
-This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.
+   ```sh
+   # Windows (primary):
+   $env:CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require"
+   npm run db:seed          # -> powershell scripts/seed-docker.ps1
+   # Linux/macOS:
+   CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" \
+     npm run db:seed
+   # or directly:  ./scripts/seed.sh
+   ```
 
-## Learn More
+   What it runs: `switch` + `redis` + `cardsvc` in the background, then
+   `clearing-sim` (clearing/settlement/prefund/default fund), `ledger-sim`
+   (ledger accounts + entries), `card-sim` (cards + tokens), `acquiring-sim`
+   (merchants + funding), `disputes-sim` (disputes), and finally `acquirer-sim`
+   to drive authorizations that write `switch_transactions` through the running
+   switch. The `adminapi` is not needed to seed.
+
+3. **Create the demo users** in Supabase Auth (service-role admin API,
+   idempotent ÔÇö updates roles, creates what's missing):
+
+   ```sh
+   npm run db:users         # node --env-file=.env.local scripts/create-users.mjs
+   ```
+
+## 5. Demo matrix
+
+| role            | email                     | password        | landing page | what's visible                                                    |
+| --------------- | ------------------------- | --------------- | ------------ | ----------------------------------------------------------------- |
+| scheme_operator | scheme-operator@clara.demo| `ClaraDemo!2026`| `/ops`       | ops, transactions, clearing, settlement, ledger, cards, merchants, disputes |
+| issuer          | issuer@clara.demo         | `ClaraDemo!2026`| `/issuer`    | issuer, cards, tokens                                             |
+| acquirer        | acquirer@clara.demo       | `ClaraDemo!2026`| `/acquirer`  | acquirer, merchants, funding, disputes                             |
+| merchant        | merchant@clara.demo       | `ClaraDemo!2026`| `/merchant`  | merchant, funding, disputes                                        |
+| viewer          | viewer@clara.demo         | `ClaraDemo!2026`| `/overview`  | overview                                                           |
+
+The canonical identities live in `scripts/create-users.mjs` and
+`sql/demo-profile.sql` (documentation manifest) ÔÇö keep all three in sync.
+
+## 6. Deploy (CLI-only)
+
+```sh
+# Vercel (frontend), from web/:
+cd web
+npx vercel --prod
+
+# Railway (Go adminapi), from the repo root:
+railway up --deploy
+
+# Supabase (schema), from the repo root:
+npx supabase db push --project-ref <ref> --include-all
+```
 
-To learn more about Next.js, take a look at the following resources:
+Environment variables to set on the deployed services:
 
-- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
-- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.
+| var                           | service  |
+| ----------------------------- | -------- |
+| `SUPABASE_URL`                | Vercel   |
+| `SUPABASE_ANON_KEY`           | Vercel   |
+| `CLARA_API_URL`               | Vercel   |
+| `NEXT_PUBLIC_SUPABASE_URL`    | Vercel   |
+| `NEXT_PUBLIC_SUPABASE_ANON_KEY`| Vercel  |
+| `SUPABASE_SERVICE_ROLE_KEY`   | Vercel (for `db:users`) |
+| `CLARA_PG_DSN`                | Railway  |
+| `CLARA_LISTEN` (=`:8083`)     | Railway  |
 
-You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!
+## 7. Tear-down / reset
 
-## Deploy on Vercel
+```sh
+# Local sims / data containers (add -v to wipe local volumes):
+docker compose -f deploy/docker-compose.yml down -v
 
-The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.
+# Reset a local dev database's seed data: re-run section 4 steps 1-2.
 
-Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
+# Cloud:
+#   Supabase: dashboard -> Project Settings -> Danger Zone -> Delete project
+#   Vercel:   npx vercel remove <project> --yes
+#   Railway:  railway down (or delete the project in the dashboard)
+```
\ No newline at end of file
diff --git a/web/package.json b/web/package.json
index 86e1609..73d8605 100644
--- a/web/package.json
+++ b/web/package.json
@@ -1,19 +1,21 @@
 {
   "name": "web",
   "version": "0.1.0",
   "private": true,
   "scripts": {
     "dev": "next dev",
     "build": "next build",
     "start": "next start",
-    "lint": "eslint"
+    "lint": "eslint",
+    "db:seed": "powershell -NoProfile -ExecutionPolicy Bypass -File scripts/seed-docker.ps1",
+    "db:users": "node --env-file=.env.local scripts/create-users.mjs"
   },
   "dependencies": {
     "@base-ui/react": "^1.7.0",
     "@supabase/ssr": "^0.12.5",
     "@supabase/supabase-js": "^2.112.4",
     "class-variance-authority": "^0.7.1",
     "clsx": "^2.1.1",
     "lucide-react": "^1.37.0",
     "next": "16.3.3",
     "next-themes": "^0.4.6",
diff --git a/web/scripts/create-users.mjs b/web/scripts/create-users.mjs
new file mode 100644
index 0000000..3d5c3ce
--- /dev/null
+++ b/web/scripts/create-users.mjs
@@ -0,0 +1,78 @@
+// syntax check via: node --check scripts/create-users.mjs
+//
+// Creates (or updates) the five Clara Network demo users in Supabase Auth using
+// the service-role key and the admin API. Idempotent: existing users are
+// matched by exact email and only have app_metadata.role (re)set; missing users
+// are created with email_confirm: true.
+//
+// Run via `npm run db:users` from web/ (passes --env-file=.env.local, which must
+// contain SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY).
+import { createClient } from "@supabase/supabase-js";
+
+const DEMO_PASSWORD = "ClaraDemo!2026";
+
+const MANIFEST = [
+  { role: "scheme_operator", email: "scheme-operator@clara.demo", password: DEMO_PASSWORD },
+  { role: "issuer", email: "issuer@clara.demo", password: DEMO_PASSWORD },
+  { role: "acquirer", email: "acquirer@clara.demo", password: DEMO_PASSWORD },
+  { role: "merchant", email: "merchant@clara.demo", password: DEMO_PASSWORD },
+  { role: "viewer", email: "viewer@clara.demo", password: DEMO_PASSWORD },
+];
+
+const supabaseUrl = process.env.SUPABASE_URL;
+const serviceRoleKey = process.env.SUPABASE_SERVICE_ROLE_KEY;
+if (!supabaseUrl || !serviceRoleKey) {
+  console.error(
+    "SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set. " +
+      'Run via `npm run db:users` (passes --env-file=.env.local) with both keys present in web/.env.local.'
+  );
+  process.exit(1);
+}
+
+const supabase = createClient(supabaseUrl, serviceRoleKey, {
+  auth: { autoRefreshToken: false, persistSession: false },
+});
+
+async function findByEmail(email) {
+  let page = 1;
+  for (;;) {
+    const { data, error } = await supabase.auth.admin.listUsers({ page, perPage: 200 });
+    if (error) throw error;
+    const match = data.users.find((u) => u.email === email);
+    if (match) return match;
+    if (data.users.length === 0 || page >= data.lastPage) return null;
+    page += 1;
+  }
+}
+
+let failed = false;
+for (const demo of MANIFEST) {
+  try {
+    const existing = await findByEmail(demo.email);
+    if (existing) {
+      const { error } = await supabase.auth.admin.updateUserById(existing.id, {
+        app_metadata: { role: demo.role },
+      });
+      if (error) throw error;
+      console.log(`updated ${demo.email} -> role ${demo.role} (id ${existing.id})`);
+    } else {
+      const { data, error } = await supabase.auth.admin.createUser({
+        email: demo.email,
+        password: demo.password,
+        email_confirm: true,
+        app_metadata: { role: demo.role },
+      });
+      if (error) throw error;
+      console.log(`created ${demo.email} -> role ${demo.role} (id ${data.user.id})`);
+    }
+  } catch (err) {
+    failed = true;
+    console.error(`FAILED ${demo.email}: ${err?.message ?? err}`);
+  }
+}
+
+if (failed) {
+  console.error("one or more demo users failed (see above)");
+  process.exit(1);
+}
+console.log("all demo users ready");
\ No newline at end of file
diff --git a/web/scripts/seed-docker.ps1 b/web/scripts/seed-docker.ps1
new file mode 100644
index 0000000..c89159d
--- /dev/null
+++ b/web/scripts/seed-docker.ps1
@@ -0,0 +1,135 @@
+# seed-docker.ps1 - populate an EXTERNAL Clara Network target database (e.g.
+# Supabase/Postgres) with demo data by running the compose data sims against it.
+# PowerShell 5.1 compatible. Windows-first; seed.sh is the Linux/macOS twin.
+#
+# Usage:
+#   $env:CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require"
+#   .\seed-docker.ps1
+# (CLARA_PG_DSN falls back to DATABASE_URL.)
+
+$ErrorActionPreference = "Stop"
+
+$scriptDir = $PSScriptRoot
+if (-not $scriptDir) {
+    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
+}
+$repoRoot = [System.IO.Path]::GetFullPath((Join-Path $scriptDir "..\.."))
+$composeFile = Join-Path $repoRoot "deploy\docker-compose.yml"
+
+if (-not (Test-Path -LiteralPath $composeFile)) {
+    Write-Host "[seed] ERROR: compose file not found: $composeFile" -ForegroundColor Red
+    exit 1
+}
+if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
+    Write-Host "[seed] ERROR: docker not found on PATH" -ForegroundColor Red
+    exit 1
+}
+
+$dsn = $env:CLARA_PG_DSN
+if (-not $dsn) {
+    $dsn = $env:DATABASE_URL
+}
+if (-not $dsn) {
+    Write-Host @"
+[seed] ERROR: no target database DSN provided.
+Set CLARA_PG_DSN (or DATABASE_URL) before running, e.g.:
+  `$env:CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require"
+  .\seed-docker.ps1
+"@ -ForegroundColor Red
+    exit 1
+}
+
+$dsnYaml = "'" + ($dsn.Replace("'", "''")) + "'"
+$overridePath = Join-Path ([System.IO.Path]::GetTempPath()) ("clara-seed-override-" + $PID + ".yml")
+
+$overrideContent = @"
+# Temporary override generated by seed-docker.ps1 - points the data services at
+# the external target DB. Deleted on exit. Do not edit.
+services:
+  postgres:
+    ports: !reset []
+"@
+foreach ($svc in @("switch", "cardsvc", "clearing-sim", "ledger-sim", "card-sim", "acquiring-sim", "disputes-sim")) {
+    $overrideContent += @"
+
+  ${svc}:
+    environment:
+      CLARA_PG_DSN: $dsnYaml
+"@
+}
+
+function Invoke-Compose {
+    param(
+        [string[]]$Arguments,
+        [string]$Label
+    )
+    Write-Host "[seed] $Label"
+    & docker @Arguments
+    if ($LASTEXITCODE -ne 0) {
+        throw "[seed] docker compose failed ($Label): exit code $LASTEXITCODE"
+    }
+}
+
+$startedBackground = $false
+try {
+    [System.IO.File]::WriteAllText($overridePath, $overrideContent, (New-Object System.Text.UTF8Encoding($false)))
+    Write-Host "[seed] wrote compose override: $overridePath"
+
+    $baseArgs = @("compose", "-f", $composeFile, "-f", $overridePath)
+
+    Invoke-Compose -Arguments ($baseArgs + @("up", "-d", "--build", "redis", "switch", "cardsvc")) -Label "starting redis, switch, cardsvc (image builds on first run)"
+    $startedBackground = $true
+
+    Write-Host "[seed] waiting for cardsvc to accept connections (localhost:18081, up to 60s)..."
+    $ready = $false
+    $deadline = (Get-Date).AddSeconds(60)
+    while ((Get-Date) -lt $deadline) {
+        $tcp = New-Object System.Net.Sockets.TcpClient
+        try {
+            $tcp.Connect("localhost", 18081)
+            $ready = $true
+        } catch {
+            Start-Sleep -Seconds 2
+        } finally {
+            $tcp.Dispose()
+        }
+        if ($ready) {
+            break
+        }
+    }
+    if (-not $ready) {
+        Write-Host "[seed] ERROR: cardsvc did not become ready on localhost:18081 within 60s" -ForegroundColor Red
+        & docker @($baseArgs + @("logs", "--tail", "100", "cardsvc"))
+        throw "cardsvc readiness timeout"
+    }
+    Write-Host "[seed] cardsvc is ready"
+
+    foreach ($svc in @("clearing-sim", "ledger-sim", "card-sim", "acquiring-sim", "disputes-sim")) {
+        Invoke-Compose -Arguments ($baseArgs + @("run", "--rm", "--no-deps", $svc)) -Label "running one-shot $svc (clearing/ledger/cards/merchants/disputes)"
+    }
+
+    Invoke-Compose -Arguments ($baseArgs + @("run", "--rm", "--no-deps", "acquirer-sim")) -Label "running acquirer-sim to generate authorizations through the running switch"
+} catch {
+    Write-Host "[seed] ERROR: $($_.Exception.Message)" -ForegroundColor Red
+    exit 1
+} finally {
+    if ($startedBackground) {
+        try {
+            & docker @(@("compose", "-f", $composeFile, "-f", $overridePath) + @("stop", "postgres", "redis", "switch", "cardsvc")) | Out-Null
+            Write-Host "[seed] stopped background services (postgres, redis, switch, cardsvc)"
+        } catch {
+            Write-Host "[seed] WARN: background teardown failed: $($_.Exception.Message)"
+        }
+    }
+    if (Test-Path -LiteralPath $overridePath) {
+        Remove-Item -LiteralPath $overridePath -Force
+    }
+}
+
+$credentialRegex = '^[a-zA-Z][a-zA-Z0-9+.\-]*://[^@/]*@'
+if ($dsn -match $credentialRegex) {
+    $masked = $dsn -replace $credentialRegex, '' -replace '\?.*$', ''
+} else {
+    $masked = $dsn -replace '\?.*$', ''
+}
+Write-Host "seed complete (target: $masked)"
\ No newline at end of file
diff --git a/web/scripts/seed.sh b/web/scripts/seed.sh
new file mode 100644
index 0000000..59cae6a
--- /dev/null
+++ b/web/scripts/seed.sh
@@ -0,0 +1,83 @@
+#!/usr/bin/env bash
+# seed.sh - populate an EXTERNAL Clara Network target database (e.g.
+# Supabase/Postgres) with demo data by running the compose data sims against it.
+# bash (set -euo pipefail). Linux/macOS twin of scripts/seed-docker.ps1.
+#
+# Usage:
+#   CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" ./seed.sh
+# (CLARA_PG_DSN falls back to DATABASE_URL.)
+
+set -euo pipefail
+
+SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
+REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
+COMPOSE_FILE="$REPO_ROOT/deploy/docker-compose.yml"
+
+if [ ! -f "$COMPOSE_FILE" ]; then
+  echo "ERROR: compose file not found: $COMPOSE_FILE" >&2
+  exit 1
+fi
+
+DSN="${CLARA_PG_DSN:-${DATABASE_URL:-}}"
+if [ -z "$DSN" ]; then
+  echo "ERROR: no target database DSN provided." >&2
+  echo 'Set CLARA_PG_DSN (or DATABASE_URL), e.g.:' >&2
+  echo '  CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" ./seed.sh' >&2
+  exit 1
+fi
+
+OVERRIDE="$(mktemp)"
+STARTED=0
+cleanup() {
+  if [ "$STARTED" = "1" ]; then
+    docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" stop postgres redis switch cardsvc >/dev/null 2>&1 || true
+    echo "[seed] stopped background services (postgres, redis, switch, cardsvc)"
+  fi
+  rm -f "$OVERRIDE"
+}
+trap cleanup EXIT
+
+DSN_YAML="'$(printf '%s' "$DSN" | sed "s/'/''/g")'"
+{
+  echo "services:"
+  echo "  postgres:"
+  echo "    ports: !reset []"
+  for svc in switch cardsvc clearing-sim ledger-sim card-sim acquiring-sim disputes-sim; do
+    echo "  $svc:"
+    echo "    environment:"
+    echo "      CLARA_PG_DSN: $DSN_YAML"
+  done
+} > "$OVERRIDE"
+
+echo "[seed] wrote compose override: $OVERRIDE"
+
+echo "[seed] starting redis, switch, cardsvc (image builds on first run)"
+docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" up -d --build redis switch cardsvc
+STARTED=1
+
+echo "[seed] waiting for cardsvc to accept connections (localhost:18081, up to 60s)..."
+READY=0
+for _ in $(seq 1 30); do
+  if (echo > /dev/tcp/localhost/18081) >/dev/null 2>&1; then
+    READY=1
+    break
+  fi
+  sleep 2
+done
+if [ "$READY" != "1" ]; then
+  echo "ERROR: cardsvc did not become ready on localhost:18081 within 60s" >&2
+  docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" logs --tail 100 cardsvc || true
+  exit 1
+fi
+echo "[seed] cardsvc is ready"
+
+for svc in clearing-sim ledger-sim card-sim acquiring-sim disputes-sim; do
+  echo "[seed] running one-shot $svc (clearing/ledger/cards/merchants/disputes)"
+  docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" run --rm --no-deps "$svc"
+done
+
+echo "[seed] running acquirer-sim to generate authorizations through the running switch"
+docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" run --rm --no-deps acquirer-sim
+
+MASKED="$(printf '%s' "$DSN" | sed -E 's#^[a-zA-Z][a-zA-Z0-9+.-]*://[^@/]*@##' | sed -E 's#\?.*$##')"
+echo "seed complete (target: $MASKED)"
\ No newline at end of file
diff --git a/web/sql/demo-profile.sql b/web/sql/demo-profile.sql
new file mode 100644
index 0000000..080d76e
--- /dev/null
+++ b/web/sql/demo-profile.sql
@@ -0,0 +1,24 @@
+-- ============================================================================
+-- Clara Network - demo users manifest
+-- ----------------------------------------------------------------------------
+-- This file is documentation ONLY. Demo users are created via
+-- `scripts/create-users.mjs` with the service-role key (managed `auth.users`
+-- cannot be inserted via SQL). Run `npm run db:users` from web/ instead.
+--
+-- Five demo roles map to Supabase `app_metadata.role`, which gates the gallery
+-- of read-only dashboards via the Next.js BFF. One shared demo password.
+--
+--   role            | email                    | app_metadata.role
+--   ----------------+--------------------------+-------------------
+--   scheme_operator | scheme-operator@clara.demo| scheme_operator
+--   issuer          | issuer@clara.demo         | issuer
+--   acquirer        | acquirer@clara.demo       | acquirer
+--   merchant        | merchant@clara.demo       | merchant
+--   viewer          | viewer@clara.demo         | viewer
+--
+--   shared demo password (demo-only, not a secret): ClaraDemo!2026
+--
+-- create-users.mjs is idempotent: it looks each email up via
+-- supabase.auth.admin.listUsers(), updates app_metadata.role if the user
+-- already exists, and otherwise creates it with email_confirm: true.
+-- ============================================================================
\ No newline at end of file
