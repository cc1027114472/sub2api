#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-sub2api}"
SWAP_SIZE_MB="${SWAP_SIZE_MB:-2048}"
ENABLE_UFW="${ENABLE_UFW:-true}"

log() {
    printf '[INFO] %s\n' "$*"
}

warn() {
    printf '[WARN] %s\n' "$*" >&2
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
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

require_root() {
    [ "${EUID}" -eq 0 ] || fail "run this script as root"
}

check_platform() {
    [ -f /etc/os-release ] || fail "cannot identify the operating system"
    # shellcheck disable=SC1091
    . /etc/os-release
    case "${ID:-}" in
        ubuntu|debian) ;;
        *) fail "unsupported operating system: ${ID:-unknown}" ;;
    esac
    command_exists apt-get || fail "apt-get is required"
}

install_host_packages() {
    log "refreshing apt metadata"
    export DEBIAN_FRONTEND=noninteractive
    apt-get update

    log "installing host prerequisites"
    apt-get install -y \
        ca-certificates \
        curl \
        fail2ban \
        git \
        openssl \
        ufw

    if ! command_exists docker; then
        apt-get install -y docker.io
    fi

    if ! docker compose version >/dev/null 2>&1; then
        if apt-cache show docker-compose-v2 >/dev/null 2>&1; then
            apt-get install -y docker-compose-v2
        elif apt-cache show docker-compose-plugin >/dev/null 2>&1; then
            apt-get install -y docker-compose-plugin
        else
            fail "Docker Compose v2 package is not available in the configured apt sources"
        fi
    fi

    if ! docker buildx version >/dev/null 2>&1; then
        if apt-cache show docker-buildx >/dev/null 2>&1; then
            apt-get install -y docker-buildx
        elif apt-cache show docker-buildx-plugin >/dev/null 2>&1; then
            apt-get install -y docker-buildx-plugin
        else
            fail "Docker Buildx package is not available in the configured apt sources"
        fi
    fi

    systemctl enable --now docker
    systemctl enable --now fail2ban
    docker version --format '{{.Server.Version}}' >/dev/null
    docker compose version >/dev/null
    docker buildx version >/dev/null
}

ensure_swap() {
    if swapon --show --noheadings 2>/dev/null | grep -q .; then
        log "swap already enabled"
        return
    fi

    if [ -e /swapfile ]; then
        fail "/swapfile exists but is not active; refusing to overwrite it"
    fi

    log "creating ${SWAP_SIZE_MB} MiB swap file"
    if command_exists fallocate; then
        fallocate -l "${SWAP_SIZE_MB}M" /swapfile
    else
        dd if=/dev/zero of=/swapfile bs=1M count="${SWAP_SIZE_MB}" status=progress
    fi
    chmod 600 /swapfile
    mkswap /swapfile >/dev/null
    swapon /swapfile

    if ! grep -qE '^/swapfile[[:space:]]' /etc/fstab; then
        printf '/swapfile none swap sw 0 0\n' >> /etc/fstab
    fi
}

configure_firewall() {
    if [ "${ENABLE_UFW}" != "true" ]; then
        warn "UFW setup skipped because ENABLE_UFW=${ENABLE_UFW}"
        return
    fi

    log "configuring UFW for SSH and web traffic"
    ufw default deny incoming >/dev/null
    ufw default allow outgoing >/dev/null
    ufw allow 22/tcp >/dev/null
    ufw allow 80/tcp >/dev/null
    ufw allow 443/tcp >/dev/null
    ufw --force enable >/dev/null
}

ensure_directories() {
    log "creating deployment data directories"
    install -d -m 700 \
        "${SCRIPT_DIR}/data" \
        "${SCRIPT_DIR}/postgres_data" \
        "${SCRIPT_DIR}/redis_data" \
        "${PROJECT_ROOT}/.deploy"
}

generate_secret() {
    openssl rand -hex 32
}

ensure_env() {
    if [ -f "${ENV_FILE}" ] && [ "${FORCE_REGENERATE_ENV:-false}" != "true" ]; then
        chmod 600 "${ENV_FILE}"
        log "preserving existing ${ENV_FILE}"
        return
    fi

    if [ -f "${ENV_FILE}" ]; then
        local backup_file
        backup_file="${ENV_FILE}.backup.$(date +%Y%m%d%H%M%S)"
        cp -- "${ENV_FILE}" "${backup_file}"
        chmod 600 "${backup_file}"
        warn "existing .env backed up to ${backup_file}"
    fi

    local commit
    local image_tag
    local postgres_password
    local redis_password
    local jwt_secret
    local totp_key
    local admin_password

    commit="$(git -C "${PROJECT_ROOT}" rev-parse --short=12 HEAD 2>/dev/null || printf 'local')"
    image_tag="git-${commit}"
    postgres_password="${POSTGRES_PASSWORD:-$(generate_secret)}"
    redis_password="${REDIS_PASSWORD:-$(generate_secret)}"
    jwt_secret="${JWT_SECRET:-$(generate_secret)}"
    totp_key="${TOTP_ENCRYPTION_KEY:-$(generate_secret)}"
    admin_password="${ADMIN_PASSWORD:-$(generate_secret)}"

    log "creating ${ENV_FILE} with generated secrets"
    (
        umask 077
        cat > "${ENV_FILE}" <<EOF
# Generated by deploy/server-install.sh. Keep this file private.
BIND_HOST=127.0.0.1
SERVER_PORT=8080
SERVER_MODE=release
RUN_MODE=standard
TZ=Asia/Shanghai

SUB2API_IMAGE=sub2api-fork
SUB2API_IMAGE_TAG=${image_tag}
NPM_CONFIG_REGISTRY=https://registry.npmjs.org
GOPROXY=https://proxy.golang.org,direct
GOSUMDB=sum.golang.org

POSTGRES_USER=sub2api
POSTGRES_PASSWORD=${postgres_password}
POSTGRES_DB=sub2api
POSTGRES_MAX_CONNECTIONS=50
POSTGRES_SHARED_BUFFERS=64MB
POSTGRES_EFFECTIVE_CACHE_SIZE=512MB
POSTGRES_MAINTENANCE_WORK_MEM=32MB

REDIS_PASSWORD=${redis_password}
REDIS_POOL_SIZE=32
REDIS_MIN_IDLE_CONNS=4

DATABASE_MAX_OPEN_CONNS=20
DATABASE_MAX_IDLE_CONNS=5
DATABASE_CONN_MAX_LIFETIME_MINUTES=30
DATABASE_CONN_MAX_IDLE_TIME_MINUTES=5

JWT_SECRET=${jwt_secret}
TOTP_ENCRYPTION_KEY=${totp_key}
ADMIN_EMAIL=${ADMIN_EMAIL:-admin@sub2api.local}
ADMIN_PASSWORD=${admin_password}

SERVER_MAX_REQUEST_BODY_SIZE=67108864
GATEWAY_MAX_BODY_SIZE=67108864
SECURITY_URL_ALLOWLIST_ENABLED=false
SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP=false
SECURITY_URL_ALLOWLIST_ALLOW_PRIVATE_HOSTS=false
EOF
        chmod 600 "${ENV_FILE}"
    )
}

build_and_start() {
    log "validating the merged Compose configuration"
    compose config >/dev/null

    log "building the application image from the checked-out fork"
    bash "${SCRIPT_DIR}/build-fork-image.sh"

    log "starting the application stack"
    compose up -d --remove-orphans

    log "waiting for service health"
    local attempt
    for attempt in $(seq 1 180); do
        if curl -fsS --max-time 5 http://127.0.0.1:8080/health >/dev/null 2>&1; then
            bash "${SCRIPT_DIR}/healthcheck.sh"
            return
        fi
        sleep 2
    done

    compose ps || true
    fail "application did not become healthy; inspect logs with deploy/healthcheck.sh and docker compose logs"
}

main() {
    require_root
    check_platform
    [ -f "${SCRIPT_DIR}/docker-compose.local.yml" ] || fail "repository deployment files are missing"
    [ -f "${SCRIPT_DIR}/docker-compose.fork.yml" ] || fail "fork Compose override is missing"

    install_host_packages
    ensure_swap
    configure_firewall
    ensure_directories
    ensure_env
    build_and_start

    warn "root SSH password login was not changed by this script; harden SSH after validating your public-key login"
    warn "the 50 GiB unallocated disk was not formatted or mounted"
    log "Sub2API fork deployment completed"
    log "application is bound to 127.0.0.1:8080; configure a real domain and Caddy before public exposure"
}

main "$@"
