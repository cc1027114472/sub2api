#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"
STATE_DIR="${PROJECT_ROOT}/.deploy"

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
command -v git >/dev/null 2>&1 || fail "git is not installed"

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

mkdir -p -- "${STATE_DIR}"
chmod 700 "${STATE_DIR}"

log "creating a database backup before update"
bash "${SCRIPT_DIR}/backup.sh"

current_image_id="$(docker inspect --format '{{.Image}}' sub2api 2>/dev/null || true)"
if [ -n "${current_image_id}" ]; then
    rollback_tag="sub2api-fork:rollback-$(date +%Y%m%d%H%M%S)"
    docker tag "${current_image_id}" "${rollback_tag}"
    printf '%s\n' "${rollback_tag}" > "${STATE_DIR}/rollback-image"
    chmod 600 "${STATE_DIR}/rollback-image"
fi

branch="$(git -C "${PROJECT_ROOT}" symbolic-ref --short HEAD)"
log "pulling the latest ${branch} commit from origin"
git -C "${PROJECT_ROOT}" fetch --prune origin "${branch}"
git -C "${PROJECT_ROOT}" pull --ff-only origin "${branch}"

new_revision="$(git -C "${PROJECT_ROOT}" rev-parse --short=12 HEAD)"
set_env_value SUB2API_IMAGE_TAG "git-${new_revision}"

log "building and starting the updated fork image"
bash "${SCRIPT_DIR}/build-fork-image.sh"
compose up -d --remove-orphans

bash "${SCRIPT_DIR}/healthcheck.sh"
printf '%s\n' "${new_revision}" > "${STATE_DIR}/last-successful-revision"
chmod 600 "${STATE_DIR}/last-successful-revision"
log "update completed at revision ${new_revision}"
