#!/usr/bin/env bash
#
# End-to-end smoke check for the dockerized Couch Potatoes stack (ticket #13).
#
# Builds every image from this repository, starts the stack with the documented
# example environment, waits (bounded) for the database to finish seeding, then
# asserts the observable contract:
#
#   - the graph contains the seed marker and at least one node,
#   - GET /healthcheck on the API succeeds,
#   - GET /ready succeeds and actually depends on the database: stopping the
#     database makes /ready fail with 503 while the API process stays up, and
#     restarting the database brings /ready back,
#   - the client serves its page on http://localhost:3000.
#
# Any failure or timeout dumps `docker compose ps` and the container logs and
# exits non-zero. The stack is always torn down, and db/data is left with only
# .gitkeep.
#
# Runnable by hand from anywhere inside the repository, and by CI:
#
#   scripts/stack-smoke.sh
#
# The timeouts below can be overridden through the environment (seconds).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

ENV_FILE="${ENV_FILE:-.env.example}"
DB_SERVICE="cp-db"
API_SERVICE="cp-api"
API_BASE_URL="${API_BASE_URL:-http://localhost:4001}"
CLIENT_URL="${CLIENT_URL:-http://localhost:3000/}"
SEED_MARKER_ID="couch-potatoes-seed-v1"

SEED_TIMEOUT_SECONDS="${SEED_TIMEOUT_SECONDS:-300}"
HTTP_TIMEOUT_SECONDS="${HTTP_TIMEOUT_SECONDS:-180}"
DB_DOWN_TIMEOUT_SECONDS="${DB_DOWN_TIMEOUT_SECONDS:-120}"
DB_RESTORE_TIMEOUT_SECONDS="${DB_RESTORE_TIMEOUT_SECONDS:-180}"
POLL_INTERVAL_SECONDS="${POLL_INTERVAL_SECONDS:-3}"

COMPOSE=(docker compose --env-file "$ENV_FILE")

log() {
  printf '==> %s\n' "$*"
}

fail() {
  printf 'smoke check failed: %s\n' "$*" >&2
  exit 1
}

compose() {
  "${COMPOSE[@]}" "$@"
}

# Generous cap: with the database stopped, the driver needs a few seconds to
# give up and turn /ready into a 503, and a short cap would hide that status.
http_status() {
  curl --silent --show-error --max-time 15 --output /dev/null \
    --write-out '%{http_code}' "$1" 2>/dev/null || true
}

# Reads one scalar from the database using the credentials the container
# already has in NEO4J_AUTH, so the script never duplicates them. The query is
# passed as a positional argument so the container shell does not strip the
# Cypher string quotes.
database_value() {
  local query="$1"
  # shellcheck disable=SC2016 # these variables must expand inside the container
  compose exec -T "$DB_SERVICE" sh -c '
    cypher-shell --format plain \
      -u "${NEO4J_AUTH%%/*}" -p "${NEO4J_AUTH#*/}" "$1"
  ' sh "$query" 2>/dev/null | tail -n 1 | tr -d '[:space:]'
}

seed_marker_count() {
  database_value "MATCH (s:SeedMarker {id: '$SEED_MARKER_ID'}) RETURN count(s) AS count" || true
}

node_count() {
  database_value "MATCH (n) RETURN count(n) AS count" || true
}

seeded_graph() {
  local count
  count="$(seed_marker_count)"
  [ "$count" = "1" ]
}

api_health_ok() {
  [ "$(http_status "$API_BASE_URL/healthcheck")" = "200" ]
}

api_ready_ok() {
  [ "$(http_status "$API_BASE_URL/ready")" = "200" ]
}

api_ready_failing() {
  [ "$(http_status "$API_BASE_URL/ready")" = "503" ]
}

client_serving() {
  local body
  body="$(curl --silent --show-error --max-time 5 "$CLIENT_URL" 2>/dev/null)" || return 1
  grep -q 'id="root"' <<<"$body"
}

api_process_running() {
  compose ps --status running --services 2>/dev/null | grep -qx "$API_SERVICE"
}

# Polls a predicate until it succeeds or the bounded timeout expires. Never
# sleeps a fixed interval and hopes.
wait_until() {
  local description="$1"
  local timeout_seconds="$2"
  shift 2
  local deadline=$((SECONDS + timeout_seconds))

  log "waiting up to ${timeout_seconds}s for ${description}"
  until "$@"; do
    if (( SECONDS >= deadline )); then
      printf 'timed out after %ss waiting for %s\n' "$timeout_seconds" "$description" >&2
      return 1
    fi
    sleep "$POLL_INTERVAL_SECONDS"
  done
  log "confirmed: ${description}"
}

dump_diagnostics() {
  {
    printf '\n===== docker compose ps =====\n'
    compose ps || true
    printf '\n===== docker compose logs =====\n'
    compose logs --no-color || true
  } >&2
}

clean_db_data() {
  if [ ! -d "$ROOT_DIR/db/data" ]; then
    return 0
  fi

  # Neo4j writes into this bind mount as its own user, so the files here are
  # not always deletable by the caller. When that happens, clean up from a
  # throwaway container as root and hand the directory back to the host user.
  if ! {
    find "$ROOT_DIR/db/data" -mindepth 1 ! -name .gitkeep -exec rm -rf {} + 2>/dev/null &&
      touch "$ROOT_DIR/db/data/.gitkeep" 2>/dev/null
  }; then
    log "db/data is owned by the database user; resetting it from a root container"
    docker run --rm --volume "$ROOT_DIR/db/data:/data" alpine:3.24 \
      sh -c "find /data -mindepth 1 ! -name .gitkeep -exec rm -rf {} +; touch /data/.gitkeep; chown -R $(id -u):$(id -g) /data" ||
      log "could not fully reset db/data; clear it by hand if the next start looks stale"
  fi

  touch "$ROOT_DIR/db/data/.gitkeep" 2>/dev/null || true
}

teardown() {
  log "tearing the stack down"
  compose down --volumes --remove-orphans --timeout 30 || true
  clean_db_data
}

on_exit() {
  local status=$?
  if (( status != 0 )); then
    dump_diagnostics
  fi
  teardown
  exit "$status"
}

trap on_exit EXIT
trap 'exit 130' INT TERM

[ -f "$ENV_FILE" ] || fail "environment file ${ENV_FILE} does not exist"

log "using environment file ${ENV_FILE}"
log "building every image from source"
compose build

log "starting the stack"
compose up --detach

wait_until "the database to finish seeding" "$SEED_TIMEOUT_SECONDS" seeded_graph

count="$(node_count)"
[[ "$count" =~ ^[0-9]+$ ]] || fail "could not read a node count from the database (got '${count}')"
(( count > 0 )) || fail "the database contains no nodes after seeding"
log "database seeded with ${count} nodes"

wait_until "the API health endpoint" "$HTTP_TIMEOUT_SECONDS" api_health_ok
wait_until "the API readiness endpoint" "$HTTP_TIMEOUT_SECONDS" api_ready_ok
wait_until "the client page on ${CLIENT_URL}" "$HTTP_TIMEOUT_SECONDS" client_serving

log "stopping ${DB_SERVICE} to prove /ready checks database connectivity"
compose stop --timeout 30 "$DB_SERVICE" >/dev/null

wait_until "readiness to fail with 503 while the database is stopped" \
  "$DB_DOWN_TIMEOUT_SECONDS" api_ready_failing
api_process_running || fail "the API process exited when the database stopped"
log "API process is still running while /ready reports 503"

log "restarting ${DB_SERVICE} to prove readiness recovers"
compose start "$DB_SERVICE" >/dev/null
wait_until "readiness to recover" "$DB_RESTORE_TIMEOUT_SECONDS" api_ready_ok

log "stack smoke check passed in ${SECONDS}s"
