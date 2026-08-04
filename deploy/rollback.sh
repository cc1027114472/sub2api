#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"
STATE_DIR="${PROJECT_ROOT}/.deploy"
ROLLBACK_FILE="${STATE_DIR}/rollback-image"

log() {
    printf '[INFO] %s\n' "$*"
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

[ "${EUID}" -eq 0 ] || fail "run this script as root"
[ -f "${ENV_FILE}" ] || fail "missing ${ENV_FILE}"
[ -s "${ROLLBACK_FILE}" ] || fail "no rollback image has been recorded"
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

set_env_value() {
    local key="$1"
    local value="$2"
    local tmp_file

    tmp_file="$(mktemp "${ENV_FILE}.XXXXXX")"
    awk -F= -v wanted="${key}" -v replacement="${value}" '
        BEGIN { replaced = 0 }
        $1 == wanted {
            print wanted "=" replacement
            replaced = 1
            next
        }
        { print }
        END {
            if (!replaced) print wanted "=" replacement
        }
    ' "${ENV_FILE}" > "${tmp_file}"
    chmod 600 "${tmp_file}"
    mv -- "${tmp_file}" "${ENV_FILE}"
}

rollback_tag="$(head -n 1 "${ROLLBACK_FILE}")"
[ -n "${rollback_tag}" ] || fail "rollback image tag is empty"
docker image inspect "${rollback_tag}" >/dev/null 2>&1 || fail "rollback image ${rollback_tag} is unavailable"

log "switching application image to ${rollback_tag}"
set_env_value SUB2API_IMAGE_TAG "${rollback_tag#sub2api-fork:}"
compose up -d --no-build --force-recreate sub2api

bash "${SCRIPT_DIR}/healthcheck.sh"
printf '%s\n' "${rollback_tag}" > "${STATE_DIR}/last-rollback-image"
chmod 600 "${STATE_DIR}/last-rollback-image"
log "application rollback completed"
