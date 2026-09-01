#!/usr/bin/env bash
#
# smoke.sh - one-click end-to-end smoke test for Clara Network.
#
# Boots the local docker-compose stack (postgres, redis, switch, issuer-sim,
# cardsvc, adminapi), seeds demo data via the one-shot sims, then runs three
# check suites and prints a PASS/FAIL table:
#
#   1. DB suite      - schema tables exist + each domain has seeded rows
#   2. Backend suite - every adminapi endpoint returns HTTP 200 + valid JSON
#   3. Frontend suite- `next build` compiles all routes + BFF auth guard 401s
#
# Usage:
#   scripts/smoke.sh            # full run: boot, seed, all suites, report
#   scripts/smoke.sh --no-seed  # reuse an already-running, already-seeded stack
#   scripts/smoke.sh --keep     # leave the stack running when done
#   scripts/smoke.sh --skip-frontend # skip the (slow) Next.js build
#
# Exit code is non-zero if ANY check fails (so it fails CI).
#
# Note: on Windows, run this from Git Bash. Go has no host toolchain requirement:
# the Go services run in Docker.

set -u

# ---------------------------------------------------------------------------
# Config
# ---------------------------------------------------------------------------
COMPOSE_FILE="deploy/docker-compose.yml"
ADMINAPI_PORT="${CLARA_ADMINAPI_PORT:-18083}"
ADMINAPI="http://localhost:${ADMINAPI_PORT}"
WEB_DIR="web"
SEED=1
KEEP=0
SKIP_FRONTEND=0

# Global state
PASS=0
FAIL=0
RESULTS=()          # "name|PASS|detail"
START_TIME=$(date +%s)

# ---------------------------------------------------------------------------
# Colors / helpers
# ---------------------------------------------------------------------------
if [ -t 1 ]; then
  C_GREEN=$'\033[32m'; C_RED=$'\033[31m'; C_YELLOW=$'\033[33m'; C_BOLD=$'\033[1m'; C_RESET=$'\033[0m'
else
  C_GREEN=""; C_RED=""; C_YELLOW=""; C_BOLD=""; C_RESET=""
fi
green() { printf '%s%s%s' "$C_GREEN" "$1" "$C_RESET"; }
red()   { printf '%s%s%s' "$C_RED" "$1" "$C_RESET"; }
yellow(){ printf '%s%s%s' "$C_YELLOW" "$1" "$C_RESET"; }
log()   { printf '%s\n' "$1"; }
step()  { printf '\n%s==%s %s ==%s\n' "$C_BOLD" "$1" "$C_RESET"; }

# record checks a single assertion result.
#   check <name> <pass(0/1)> <detail>
check() {
  local name="$1" ok="$2" detail="${3:-}"
  if [ "$ok" -eq 0 ]; then
    PASS=$((PASS+1)); RESULTS+=("$name|PASS|$detail")
    printf '  %s  %s  %s\n' "$(green PASS)" "$name" "${detail:+($detail)}"
  else
    FAIL=$((FAIL+1)); RESULTS+=("$name|FAIL|$detail")
    printf '  %s  %s  %s\n' "$(red FAIL)" "$name" "${detail:+($detail)}"
  fi
}

# json_ok parses an endpoint response and asserts HTTP 200 + non-empty JSON.
#   json_ok <name> <url> [jq-path-expected-not-null]
json_ok() {
  local name="$1" url="$2" json_path="${3:-}"
  local code body
  code=$(curl -s -o /tmp/clara_body.$$ -w '%{http_code}' "$url")
  body=$(cat /tmp/clara_body.$$ 2>/dev/null); rm -f /tmp/clara_body.$$
  if [ "$code" != "200" ]; then
    check "$name" 1 "HTTP $code: $body"
    return
  fi
  if ! echo "$body" | grep -q '^{'; then
    check "$name" 1 "not JSON: $body"
    return
  fi
  if [ -n "$json_path" ]; then
    if command -v jq >/dev/null 2>&1; then
      local val
      val=$(echo "$body" | jq -r "$json_path" 2>/dev/null)
      case "$val" in
        ""|null|0|[]|{}) check "$name" 1 "empty at $json_path" ;;
        *) check "$name" 0 "ok" ;;
      esac
    else
      check "$name" 0 "ok (jq not present; HTTP 200 + JSON only)"
    fi
  else
    check "$name" 0 "ok"
  fi
}

# wait_health polls a URL until it returns HTTP 200 or the timeout elapses.
wait_health() {
  local url="$1" label="$2" tries="${3:-60}"
  log "  waiting for $label ($url) ..."
  local i=0
  while [ "$i" -lt "$tries" ]; do
    code=$(curl -s -o /dev/null -w '%{http_code}' "$url" 2>/dev/null || true)
    if [ "$code" = "200" ]; then
      return 0
    fi
    i=$((i+1)); sleep 2
  done
  return 1
}

# ---------------------------------------------------------------------------
# Args
# ---------------------------------------------------------------------------
for arg in "$@"; do
  case "$arg" in
    --no-seed) SEED=0 ;;
    --keep)    KEEP=1 ;;
    --skip-frontend) SKIP_FRONTEND=1 ;;
    -h|--help)
      grep -E '^#( |$)' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) log "$(red "unknown arg: $arg")"; grep -E '^#( |$)' "$0" | sed 's/^# \{0,1\}//'; exit 2 ;;
  esac
done

log "$(yellow "Clara Network smoke test")"
log "  adminapi: $ADMINAPI   compose: $COMPOSE_FILE"

# ---------------------------------------------------------------------------
# 0. Boot the stack
# ---------------------------------------------------------------------------
log ""
if ! docker compose version >/dev/null 2>&1; then
  log "$(red "ERROR: docker compose not available on PATH")"
  exit 2
fi

step "Boot docker-compose stack"
log "  docker compose up -d --build (this can take a while on first run) ..."
if ! docker compose -f "$COMPOSE_FILE" up -d --build; then
  log "$(red "ERROR: docker compose up failed")"
  docker compose -f "$COMPOSE_FILE" logs --tail=50
  exit 2
fi

if ! wait_health "$ADMINAPI/health" "adminapi" 60; then
  log "$(red "ERROR: adminapi did not become healthy at $ADMINAPI/health in 120s")"
  docker compose -f "$COMPOSE_FILE" logs adminapi --tail=80
  exit 2
fi
log "  adminapi is $(green healthy)"

# ---------------------------------------------------------------------------
# 1. Seed demo data (unless --no-seed)
# ---------------------------------------------------------------------------
if [ "$SEED" -eq 1 ]; then
  step "Seed demo data (one-shot sims)"
  # The one-shot sims populate postgres with clearing/ledger/cards/merchants/
  # disputes/auth data. They exit on their own, so use --no-deps + restart no.
  for svc in clearing-sim ledger-sim card-sim acquiring-sim disputes-sim acquirer-sim; do
    log "  running $svc ..."
    docker compose -f "$COMPOSE_FILE" run --rm --no-deps "$svc" >/dev/null 2>&1 \
      || docker compose -f "$COMPOSE_FILE" run --rm --no-deps "$svc"
    # Ignore exit codes; data assertions below are the real signal.
  done
fi

# ---------------------------------------------------------------------------
# 2. DB suite
# ---------------------------------------------------------------------------
step "DB suite (schema + seed data)"
psqlq() { docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U clara -d clara -tA -c "$1" 2>&1; }

for table in switch_transactions clearing_records net_positions \
             settlement_instructions prefund_accounts default_fund \
             ledger_accounts ledger_entries bin_ranges cards tokens \
             merchants funding_lines screening_lists disputes dispute_transactions; do
  n=$(psqlq "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='$table';")
  if [ "$n" = "1" ]; then check "db: table $table exists" 0; else check "db: table $table exists" 1 "$n"; fi
done

for t in switch_transactions clearing_records net_positions settlement_instructions \
         ledger_accounts ledger_entries bin_ranges cards tokens merchants \
         funding_lines screening_lists disputes; do
  n=$(psqlq "SELECT count(*) FROM $t;" 2>/dev/null)
  if [ -n "$n" ] && [ "$n" -gt 0 ] 2>/dev/null; then
    check "db: $t has $n row(s)" 0
  else
    check "db: $t has data" 1 "row count=$n"
  fi
done

# ---------------------------------------------------------------------------
# 3. Backend API suite (adminapi)
# ---------------------------------------------------------------------------
step "Backend API suite (adminapi endpoints)"
check "api: /health" 0 "via boot wait"
json_ok "api: /dashboard"        "$ADMINAPI/api/v1/dashboard"
json_ok "api: /transactions"     "$ADMINAPI/api/v1/transactions?limit=5"
json_ok "api: clearing/cycles"   "$ADMINAPI/api/v1/clearing/cycles"
# clearing/records, clearing/positions and settlement/instructions require a
# ?cycle= param (matching how the frontend drives them), so resolve a real one
# from /clearing/cycles first.
CYCLE=$(curl -s "$ADMINAPI/api/v1/clearing/cycles" | grep -o '"items":\[[^]]*\]' | grep -o '"[^"]*"' | tail -n +2 | tr -d '"' | head -n 1)
if [ -z "$CYCLE" ]; then
  check "api: resolve clearing cycle" 1 "no cycles returned by /clearing/cycles"
else
  check "api: resolve clearing cycle" 0 "cycle=$CYCLE"
  json_ok "api: clearing/records?cycle"   "$ADMINAPI/api/v1/clearing/records?cycle=$CYCLE"
  json_ok "api: clearing/positions?cycle" "$ADMINAPI/api/v1/clearing/positions?cycle=$CYCLE"
  json_ok "api: settlement/instructions?cycle" "$ADMINAPI/api/v1/settlement/instructions?cycle=$CYCLE"
fi
json_ok "api: settlement/prefunds"     "$ADMINAPI/api/v1/settlement/prefunds"
json_ok "api: settlement/default-fund" "$ADMINAPI/api/v1/settlement/default-fund"
json_ok "api: ledger/accounts"   "$ADMINAPI/api/v1/ledger/accounts"
json_ok "api: ledger/entries"    "$ADMINAPI/api/v1/ledger/entries?limit=5"
json_ok "api: cards"             "$ADMINAPI/api/v1/cards?limit=5"
json_ok "api: bin-ranges"        "$ADMINAPI/api/v1/bin-ranges"
json_ok "api: tokens"            "$ADMINAPI/api/v1/tokens?limit=5"
json_ok "api: merchants"         "$ADMINAPI/api/v1/merchants?limit=5"
json_ok "api: disputes"          "$ADMINAPI/api/v1/disputes?limit=5"
json_ok "api: disputes/overdue"  "$ADMINAPI/api/v1/disputes/overdue"

# ---------------------------------------------------------------------------
# 4. Frontend suite (build + runtime auth/BFF guard)
# ---------------------------------------------------------------------------
if [ "$SKIP_FRONTEND" -eq 1 ]; then
  log "\n$(yellow "frontend suite skipped (--skip-frontend)")"
else
  step "Frontend suite (next build + runtime auth/BFF guard)"
  if [ -d "$WEB_DIR/node_modules" ]; then
    log "  next build ..."
    ( cd "$WEB_DIR"
      # Build-time env is not required because getEnv() reads at request time,
      # but NEXT_PUBLIC_* vars must be present so the client bundle inlines them.
      NEXT_PUBLIC_SUPABASE_URL="${NEXT_PUBLIC_SUPABASE_URL:-http://localhost:54321}" \
      NEXT_PUBLIC_SUPABASE_ANON_KEY="${NEXT_PUBLIC_SUPABASE_ANON_KEY:-smoke-placeholder-anon}" \
      npm run build >/tmp/clara_nextbuild.$$ 2>&1
      if [ "$?" -eq 0 ]; then
        check "fe: next build" 0
      else
        check "fe: next build" 1 "see /tmp/clara_nextbuild.$$"
      fi
    )

    # Boot the production server with smoke env and probe route/auth behaviour.
    log "  next start (port 3111) ..."
    FE_PORT=3111
    ( cd "$WEB_DIR"
      PORT="$FE_PORT" \
      NEXT_PUBLIC_SUPABASE_URL="${NEXT_PUBLIC_SUPABASE_URL:-http://localhost:54321}" \
      NEXT_PUBLIC_SUPABASE_ANON_KEY="${NEXT_PUBLIC_SUPABASE_ANON_KEY:-smoke-placeholder-anon}" \
      SUPABASE_URL="${NEXT_PUBLIC_SUPABASE_URL:-http://localhost:54321}" \
      SUPABASE_ANON_KEY="${NEXT_PUBLIC_SUPABASE_ANON_KEY:-smoke-placeholder-anon}" \
      CLARA_API_URL="${CLARA_API_URL:-$ADMINAPI}" \
      npx next start -p "$FE_PORT" >/tmp/clara_nextstart.$$ 2>&1 &
      FE_PID=$!
    )
    # wait for the server to come up
    fe_ready=""
    for i in $(seq 1 30); do
      code=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/login" 2>/dev/null || true)
      if [ "$code" != "000" ] && [ -n "$code" ]; then fe_ready=1; break; fi
      sleep 1
    done
    if [ -z "$fe_ready" ]; then
      check "fe: server startup" 1 "next start did not come up; see /tmp/clara_nextstart.$$"
    else
      check "fe: server startup" 0
      code=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/login")
      [ "$code" = "200" ] && check "fe: /login returns 200" 0 || check "fe: /login returns 200" 1 "HTTP $code"
      code=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/")
      [ "$code" = "307" ] && check "fe: / redirects to login (307)" 0 || check "fe: / redirects to login (307)" 1 "HTTP $code"
      code=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/ops")
      [ "$code" = "307" ] && check "fe: protected page /ops redirects (307)" 0 || check "fe: protected page /ops redirects (307)" 1 "HTTP $code"
      # BFF proxy must return 401 JSON for unauthenticated API calls, NOT a 307
      # redirect (see web/src/middleware.ts matcher). This catches the bug where
      # the auth middleware hijacked /api/data/* away from the route's 401.
      no_redirect=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/api/data/dashboard")
      [ "$no_redirect" = "401" ] && check "fe: BFF /api/data/dashboard unauthenticated -> 401" 0 || check "fe: BFF /api/data/dashboard unauthenticated -> 401 (expected)" 1 "HTTP $no_redirect (307 = middleware hijack bug)"
      no_redirect2=$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$FE_PORT/api/data/clearing/records")
      [ "$no_redirect2" = "401" ] && check "fe: BFF /api/data/clearing/records -> 401" 0 || check "fe: BFF /api/data/clearing/records -> 401" 1 "HTTP $no_redirect2"
      # stop the dev server (it was launched from a subshell, so kill by port)
      if command -v fuser >/dev/null 2>&1; then
        fuser -k "${FE_PORT}/tcp" 2>/dev/null
      else
        pkill -f "next start.*-p $FE_PORT" 2>/dev/null
      fi
    fi
  else
    check "fe: next build" 1 "web/node_modules missing; run npm install in web/"
  fi
fi

# ---------------------------------------------------------------------------
# 5. Report
# ---------------------------------------------------------------------------
step "Results"
elapsed=$(( $(date +%s) - START_TIME ))
printf '  %s\n' "$(yellow "summary: $PASS passed, $FAIL failed (${elapsed}s)")"
for r in "${RESULTS[@]}"; do
  IFS='|' read -r name ok detail <<< "$r"
  if [ "$ok" = "PASS" ]; then printf '  %s  %s\n' "$(green PASS)" "$name"; else printf '  %s  %s\n' "$(red FAIL)" "$name"; fi
done

if [ "$KEEP" -eq 1 ]; then
  log "\n$(yellow "stack left running (--keep). Tear down with:") docker compose -f $COMPOSE_FILE down"
else
  log "\n  tearing down stack ..."
  docker compose -f "$COMPOSE_FILE" down >/dev/null 2>&1
fi

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
log "\n$(green "ALL CHECKS PASSED")"
exit 0
