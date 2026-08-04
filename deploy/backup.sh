#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"
BACKUP_ROOT="${BACKUP_ROOT:-${SCRIPT_DIR}/backups}"

log() {
    printf '[INFO] %s\n' "$*"
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

[ "${EUID}" -eq 0 ] || fail "run this script as root"
[ -f "${ENV_FILE}" ] || fail "missing ${ENV_FILE}"
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

timestamp="$(date +%Y%m%d-%H%M%S)"
backup_dir="${BACKUP_ROOT}/${timestamp}"
mkdir -p -- "${backup_dir}"
chmod 700 "${BACKUP_ROOT}" "${backup_dir}"

POSTGRES_USER="$(env_value POSTGRES_USER)"
POSTGRES_DB="$(env_value POSTGRES_DB)"
POSTGRES_USER="${POSTGRES_USER:-sub2api}"
POSTGRES_DB="${POSTGRES_DB:-sub2api}"

cleanup() {
    compose start redis sub2api >/dev/null 2>&1 || true
}
trap cleanup EXIT

log "writing PostgreSQL logical backup"
compose exec -T postgres pg_dump \
    --username "${POSTGRES_USER}" \
    --dbname "${POSTGRES_DB}" \
    --format=custom \
    | gzip -9 > "${backup_dir}/postgres.dump.gz"

log "pausing application and Redis while copying runtime data"
compose stop sub2api redis >/dev/null
tar -czf "${backup_dir}/runtime-data.tgz" \
    -C "${SCRIPT_DIR}" \
    data \
    redis_data

cp -- "${ENV_FILE}" "${backup_dir}/.env"
chmod 600 "${backup_dir}/.env"
git -C "${PROJECT_ROOT}" rev-parse HEAD > "${backup_dir}/git-revision"
docker inspect --format '{{.Image}}' sub2api 2>/dev/null > "${backup_dir}/application-image-id" || true
printf '%s\n' "${timestamp}" > "${backup_dir}/created-at"

log "backup created at ${backup_dir}"
