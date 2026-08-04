#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"

log() {
    printf '[INFO] %s\n' "$*"
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

command -v docker >/dev/null 2>&1 || fail "docker is not installed"
command -v curl >/dev/null 2>&1 || fail "curl is not installed"
[ -f "${ENV_FILE}" ] || fail "missing ${ENV_FILE}"

compose() {
    (
        cd -- "${SCRIPT_DIR}"
        docker compose \
            --project-name "${COMPOSE_PROJECT_NAME}" \
            -f docker-compose.local.yml \
            -f docker-compose.fork.yml \
            "$@"
    )
}

env_value() {
    local key="$1"
    awk -F= -v wanted="${key}" '
        $1 == wanted {
            sub(/^[^=]*=/, "")
            print
            exit
        }
    ' "${ENV_FILE}"
}

log "Compose service status:"
compose ps

if ! curl -fsS --max-time 10 http://127.0.0.1:8080/health >/dev/null; then
    fail "application health endpoint is not ready"
fi
log "application health: OK"

POSTGRES_USER="$(env_value POSTGRES_USER)"
POSTGRES_DB="$(env_value POSTGRES_DB)"
POSTGRES_USER="${POSTGRES_USER:-sub2api}"
POSTGRES_DB="${POSTGRES_DB:-sub2api}"

if ! compose exec -T postgres pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null; then
    fail "PostgreSQL is not ready"
fi
log "PostgreSQL readiness: OK"

if ! compose exec -T redis redis-cli ping 2>/dev/null | grep -qx 'PONG'; then
    fail "Redis is not ready"
fi
log "Redis readiness: OK"

if [ -n "$(compose port postgres 5432 2>/dev/null || true)" ]; then
    fail "PostgreSQL is unexpectedly published on the host"
fi
if [ -n "$(compose port redis 6379 2>/dev/null || true)" ]; then
    fail "Redis is unexpectedly published on the host"
fi
log "database and Redis host exposure: none"

if command -v ss >/dev/null 2>&1; then
    log "Relevant host listeners:"
    ss -lnt 2>/dev/null | grep -E ':(22|80|443|8080)[[:space:]]' || true
fi

log "health check passed"
