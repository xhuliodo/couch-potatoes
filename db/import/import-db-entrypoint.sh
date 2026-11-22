#!/bin/bash
set -m

# Log the info with the same format as NEO4J outputs
log_info() {
  printf '%s INFO  Wrapper: %s\n' "$(date -u +"%Y-%m-%d %H:%M:%S:%3N%z")" "$1"
}

# Credentials come from the environment as NEO4J_AUTH (user/password) or as
# NEO4J_USER and NEO4J_PASSWORD; otherwise use the development defaults.
if [ -z "${NEO4J_AUTH:-}" ]; then
  NEO4J_USER="${NEO4J_USER:-neo4j}"
  NEO4J_PASSWORD="${NEO4J_PASSWORD:-letmein}"
  NEO4J_AUTH="${NEO4J_USER}/${NEO4J_PASSWORD}"
else
  NEO4J_USER="${NEO4J_USER:-${NEO4J_AUTH%%/*}}"
  NEO4J_PASSWORD="${NEO4J_PASSWORD:-${NEO4J_AUTH#*/}}"
fi
export NEO4J_AUTH NEO4J_USER NEO4J_PASSWORD

IMPORT_FILE="${IMPORT_FILE:-/import/import.cypher}"
NEO4J_IMPORT_WAIT_SECONDS="${NEO4J_IMPORT_WAIT_SECONDS:-300}"
SEED_MARKER_ID="couch-potatoes-seed-v1"

cypher_shell() {
  cypher-shell -a bolt://localhost:7687 -u "$NEO4J_USER" -p "$NEO4J_PASSWORD" "$@"
}

# Plain format prints a header and then the value; return the last line.
last_plain_value() {
  tail -n 1 | tr -d '[:space:]'
}

node_count() {
  cypher_shell --format plain "MATCH (n) RETURN count(n) AS count" 2>/dev/null | last_plain_value
}

already_seeded() {
  local count
  count=$(cypher_shell --format plain \
    "MATCH (s:SeedMarker {id: '$SEED_MARKER_ID'}) RETURN count(s) AS count" 2>/dev/null | last_plain_value)
  [ "$count" = "1" ]
}

wait_for_neo4j() {
  local deadline=$((SECONDS + NEO4J_IMPORT_WAIT_SECONDS))
  while [ "$SECONDS" -lt "$deadline" ]; do
    if ! kill -0 "$NEO4J_PID" 2>/dev/null; then
      log_info "IMPORT-SCRIPT: Neo4j exited before becoming ready"
      return 1
    fi
    if cypher_shell --format plain "RETURN 1 AS ok" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done
  log_info "IMPORT-SCRIPT: Timed out after ${NEO4J_IMPORT_WAIT_SECONDS}s waiting for Neo4j"
  return 1
}

seed_if_needed() {
  if [ ! -f "$IMPORT_FILE" ]; then
    log_info "IMPORT-SCRIPT: No import script found at ${IMPORT_FILE}"
    return
  fi

  if already_seeded; then
    log_info "IMPORT-SCRIPT: Graph already seeded, skipping import"
    return
  fi

  log_info "IMPORT-SCRIPT: Starting import, you might have to wait a while depending on vm specs"
  if cypher_shell -f "$IMPORT_FILE"; then
    cypher_shell --format plain \
      "MERGE (s:SeedMarker {id: '$SEED_MARKER_ID'}) SET s.seededAt = timestamp()" >/dev/null
    log_info "IMPORT-SCRIPT: Done importing, graph has $(node_count) nodes"
  else
    log_info "IMPORT-SCRIPT: Import failed, leaving the seed marker unset so the next start retries"
  fi
}

# Start the primary process in the background so the import can talk to it.
/docker-entrypoint.sh neo4j &
NEO4J_PID=$!

shutdown() {
  trap - TERM INT
  if kill -0 "$NEO4J_PID" 2>/dev/null; then
    log_info "IMPORT-SCRIPT: Forwarding shutdown signal to Neo4j"
    kill -TERM "$NEO4J_PID" 2>/dev/null || true
  fi
}
trap shutdown TERM INT

log_info "IMPORT-SCRIPT: Waiting until neo4j stats at :7474 ..."
if wait_for_neo4j; then
  seed_if_needed
fi

# Bring the primary process back into the foreground and propagate its exit code.
wait "$NEO4J_PID"
exit $?
