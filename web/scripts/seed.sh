#!/usr/bin/env bash
# seed.sh - populate an EXTERNAL Clara Network target database (e.g.
# Supabase/Postgres) with demo data by running the compose data sims against it.
# bash (set -euo pipefail). Linux/macOS twin of scripts/seed-docker.ps1.
#
# Usage:
#   CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" ./seed.sh
# (CLARA_PG_DSN falls back to DATABASE_URL.)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COMPOSE_FILE="$REPO_ROOT/deploy/docker-compose.yml"

if [ ! -f "$COMPOSE_FILE" ]; then
  echo "ERROR: compose file not found: $COMPOSE_FILE" >&2
  exit 1
fi

DSN="${CLARA_PG_DSN:-${DATABASE_URL:-}}"
if [ -z "$DSN" ]; then
  echo "ERROR: no target database DSN provided." >&2
  echo 'Set CLARA_PG_DSN (or DATABASE_URL), e.g.:' >&2
  echo '  CLARA_PG_DSN="postgres://user:pass@host:5432/db?sslmode=require" ./seed.sh' >&2
  exit 1
fi

OVERRIDE="$(mktemp)"
STARTED=0
cleanup() {
  if [ "$STARTED" = "1" ]; then
    docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" stop postgres redis switch cardsvc >/dev/null 2>&1 || true
    echo "[seed] stopped background services (postgres, redis, switch, cardsvc)"
  fi
  rm -f "$OVERRIDE"
}
trap cleanup EXIT

DSN_YAML="'$(printf '%s' "$DSN" | sed "s/'/''/g")'"
{
  echo "services:"
  echo "  postgres:"
  echo "    ports: !reset []"
  for svc in switch cardsvc clearing-sim ledger-sim card-sim acquiring-sim disputes-sim; do
    echo "  $svc:"
    echo "    environment:"
    echo "      CLARA_PG_DSN: $DSN_YAML"
  done
} > "$OVERRIDE"

echo "[seed] wrote compose override: $OVERRIDE"

echo "[seed] starting redis, switch, cardsvc (image builds on first run)"
docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" up -d --build redis switch cardsvc
STARTED=1

echo "[seed] waiting for cardsvc to accept connections (localhost:18081, up to 60s)..."
READY=0
for _ in $(seq 1 30); do
  if (echo > /dev/tcp/localhost/18081) >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 2
done
if [ "$READY" != "1" ]; then
  echo "ERROR: cardsvc did not become ready on localhost:18081 within 60s" >&2
  docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" logs --tail 100 cardsvc || true
  exit 1
fi
echo "[seed] cardsvc is ready"

for svc in clearing-sim ledger-sim card-sim acquiring-sim disputes-sim; do
  echo "[seed] running one-shot $svc (clearing/ledger/cards/merchants/disputes)"
  docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" run --rm --no-deps "$svc"
done

echo "[seed] running acquirer-sim to generate authorizations through the running switch"
docker compose -f "$COMPOSE_FILE" -f "$OVERRIDE" run --rm --no-deps acquirer-sim

MASKED="$(printf '%s' "$DSN" | sed -E 's#^[a-zA-Z][a-zA-Z0-9+.-]*://[^@/]*@##' | sed -E 's#\?.*$##')"
echo "seed complete (target: $MASKED)"