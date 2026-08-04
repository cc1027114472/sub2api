#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"
BACKUP_DIR="${1:-}"
RESTORE_CONFIRM="${RESTORE_CONFIRM:-}"

log() {
    printf '[INFO] %s\n' "$*"
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

[ "${EUID}" -eq 0 ] || fail "run this script as root"
[ -n "${BACKUP_DIR}" ] || fail "usage: RESTORE_CONFIRM=YES $0 /path/to/backup"
[ "${RESTORE_CONFIRM}" = "YES" ] || fail "set RESTORE_CONFIRM=YES to acknowledge the destructive restore operation"
[ -f "${ENV_FILE}" ] || fail "missing ${ENV_FILE}"
[ -d "${BACKUP_DIR}" ] || fail "backup directory does not exist: ${BACKUP_DIR}"
[ -f "${BACKUP_DIR}/postgres.dump.gz" ] || fail "backup is missing postgres.dump.gz"
[ -f "${BACKUP_DIR}/runtime-data.tgz" ] || fail "backup is missing runtime-data.tgz"
command -v docker >/dev/null 2>&1 || fail "docker is not installed"

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

POSTGRES_USER="$(env_value POSTGRES_USER)"
POSTGRES_DB="$(env_value POSTGRES_DB)"
POSTGRES_USER="${POSTGRES_USER:-sub2api}"
POSTGRES_DB="${POSTGRES_DB:-sub2api}"

restore_staging="${PROJECT_ROOT}/.deploy/restore-$(date +%Y%m%d-%H%M%S)"
mkdir -p -- "${restore_staging}"
chmod 700 "${PROJECT_ROOT}/.deploy" "${restore_staging}"

log "stopping application and Redis"
compose stop sub2api redis >/dev/null

log "moving current runtime data aside"
for directory in data redis_data; do
    if [ -e "${SCRIPT_DIR}/${directory}" ]; then
        mv -- "${SCRIPT_DIR}/${directory}" "${restore_staging}/${directory}.previous"
    fi
done

log "restoring runtime data"
tar -xzf "${BACKUP_DIR}/runtime-data.tgz" -C "${SCRIPT_DIR}"
chmod 700 "${SCRIPT_DIR}/data" "${SCRIPT_DIR}/redis_data"

log "restoring PostgreSQL logical backup"
gzip -dc "${BACKUP_DIR}/postgres.dump.gz" \
    | compose exec -T postgres pg_restore \
        --username "${POSTGRES_USER}" \
        --dbname "${POSTGRES_DB}" \
        --clean \
        --if-exists \
        --no-owner \
        --no-privileges

log "starting restored services"
compose start redis sub2api
bash "${SCRIPT_DIR}/healthcheck.sh"

log "restore completed"
log "previous runtime data is recoverable at ${restore_staging}"
