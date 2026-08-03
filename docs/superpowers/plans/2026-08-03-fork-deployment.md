# Fork Deployment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task with review checkpoints.

**Goal:** Add a safe, fork-aware Docker deployment workflow and use it to install the checked-out fork on the target Ubuntu server.

**Architecture:** Keep the repository's existing Compose files unchanged. Add a Compose override that builds the application from the fork's root `Dockerfile`, plus idempotent host, health, backup, update, and rollback scripts. The server runs the application on localhost with PostgreSQL and Redis isolated inside the Compose network.

**Tech Stack:** Bash 4+, Docker Engine, Docker Compose v2, PostgreSQL 18 Alpine, Redis 8 Alpine, Caddy template, Ubuntu 22.04.

---

### Task 1: Add the fork-specific Compose override

**Files:**
- Create: `deploy/docker-compose.fork.yml`

- [ ] Define a `sub2api` service override with `build.context: ..`, `dockerfile: Dockerfile`, reproducible Go/npm registry build arguments, and an image name derived from `SUB2API_IMAGE`.
- [ ] Leave PostgreSQL, Redis, health checks, persistence, and network behavior inherited from `docker-compose.local.yml`.
- [ ] Validate the merged model with:

```bash
docker compose -f deploy/docker-compose.local.yml -f deploy/docker-compose.fork.yml config
```

### Task 2: Add idempotent host installation

**Files:**
- Create: `deploy/server-install.sh`

- [ ] Require root and Ubuntu/Debian-compatible apt.
- [ ] Install `ca-certificates`, `curl`, `git`, `openssl`, `ufw`, `fail2ban`, `docker.io`, and a Compose v2 package when missing.
- [ ] Enable Docker and create a 2 GiB swap file only when no swap is active.
- [ ] Generate `.env` once with localhost binding, fixed secrets, low-memory connection-pool defaults, and secure URL validation defaults.
- [ ] Create `data`, `postgres_data`, and `redis_data` with restrictive permissions.
- [ ] Run Compose config validation, build the fork image, start the stack, and poll the local health endpoint.
- [ ] Never format or mount `/dev/vdb`, disable the current root login, or print secret values.

### Task 3: Add health, backup, update, and rollback operations

**Files:**
- Create: `deploy/healthcheck.sh`
- Create: `deploy/backup.sh`
- Create: `deploy/update.sh`
- Create: `deploy/rollback.sh`

- [ ] Keep health output secret-free and return non-zero on failed application, PostgreSQL, or Redis readiness.
- [ ] Use `pg_dump` inside the PostgreSQL container and archive application data plus `.env` into a timestamped backup directory with mode `700`.
- [ ] Make updates build the checked-out fork revision, record the current image ID, start the new stack, and verify health.
- [ ] Make rollback restore the recorded image tag without deleting persistent data.

### Task 4: Validate, commit, publish, and deploy

**Files:**
- Modify: none beyond the files above

- [ ] Run `bash -n` against every new shell script and the repository's deployment tests.
- [ ] Commit the deployment changes on `codex/server-deploy`.
- [ ] Push the branch to the user's fork.
- [ ] On the server, clone the pushed branch into `/opt/sub2api`, run `deploy/server-install.sh`, and collect health output.
- [ ] Verify containers, localhost health, database/Redis readiness, and host listeners.
