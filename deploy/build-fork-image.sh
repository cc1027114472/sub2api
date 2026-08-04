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

require_value() {
    local variable_name="$1"
    local default_value="$2"
    local shell_value="${!variable_name-}"
    local file_value=""

    if [ -f "${ENV_FILE}" ]; then
        file_value="$(env_value "${variable_name}")"
    fi

    if [ -n "${shell_value}" ]; then
        printf '%s\n' "${shell_value}"
    elif [ -n "${file_value}" ]; then
        printf '%s\n' "${file_value}"
    else
        printf '%s\n' "${default_value}"
    fi
}

build_frontend_cache() {
    local node_image
    local node_heap
    local npm_registry
    local image_name
    local cache_tag

    node_image="$(require_value NODE_IMAGE node:24-alpine)"
    node_heap="$(require_value NODE_MAX_OLD_SPACE_SIZE 2048)"
    npm_registry="$(require_value NPM_CONFIG_REGISTRY https://registry.npmjs.org)"
    image_name="$(require_value SUB2API_IMAGE sub2api-fork)"
    cache_tag="${image_name}:frontend-cache"

    log "building and caching the frontend stage (Node heap: ${node_heap} MiB)"
    DOCKER_BUILDKIT=1 docker build \
        --pull \
        --progress=plain \
        --target frontend-builder \
        --file "${PROJECT_ROOT}/Dockerfile" \
        --tag "${cache_tag}" \
        --build-arg "NODE_IMAGE=${node_image}" \
        --build-arg "NODE_MAX_OLD_SPACE_SIZE=${node_heap}" \
        --build-arg "NPM_CONFIG_REGISTRY=${npm_registry}" \
        "${PROJECT_ROOT}"
}

main() {
    [ -f "${ENV_FILE}" ] || fail "missing ${ENV_FILE}"
    [ -f "${PROJECT_ROOT}/Dockerfile" ] || fail "missing ${PROJECT_ROOT}/Dockerfile"
    command -v docker >/dev/null 2>&1 || fail "docker is not installed"
    docker buildx version >/dev/null 2>&1 || fail "Docker Buildx is required; rerun deploy/server-install.sh to install it"

    log "validating the merged Compose configuration"
    compose config >/dev/null

    build_frontend_cache

    log "building the backend and final runtime image"
    DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 compose build --pull sub2api
}

main "$@"
