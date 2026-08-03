# Fork Deployment Design

## Goal

Deploy the current `main` branch of `cc1027114472/sub2api` to the target Ubuntu server without accidentally replacing fork changes with the upstream prebuilt image.

## Architecture

The server will clone the fork into `/opt/sub2api`, build the application image locally from the repository `Dockerfile`, and run the application, PostgreSQL, and Redis with Docker Compose. The application will bind to `127.0.0.1:8080` so an HTTPS reverse proxy can be added after a domain is provided; PostgreSQL and Redis remain on the Compose network and are not published to the host.

The deployment uses local data directories under `deploy/` so the application, database, Redis data, and `.env` can be backed up together. The initial script will prepare Docker, a 2 GiB swap file when no swap exists, secure environment defaults, and the data directories. It will not format or mount the unallocated 50 GiB disk.

## Files

- `deploy/docker-compose.fork.yml`: overrides the upstream image with a local build from the fork and keeps the upstream local-directory persistence layout.
- `deploy/server-install.sh`: idempotent host preparation, secret generation, image build, startup, and local health verification.
- `deploy/healthcheck.sh`: reports Compose service state, application health, database readiness, Redis readiness, and listening ports without printing secrets.
- `deploy/backup.sh`: stops application writes briefly, creates a PostgreSQL dump, and archives application/configuration data without deleting the live deployment.
- `deploy/update.sh`: records the current image, builds and starts the checked-out fork revision, then runs health verification.
- `deploy/rollback.sh`: restores the previous Docker image tag recorded by `update.sh` and verifies health.
- `deploy/Caddyfile`: remains the reverse-proxy template; public HTTPS is deferred until a real domain is supplied.

## Safety and operational rules

- Never print `.env`, generated passwords, OAuth secrets, or database dumps.
- Refuse to run as a non-root user for host installation.
- Refuse to overwrite an existing `.env` unless `FORCE_REGENERATE_ENV=true` is explicitly supplied.
- Refuse destructive data operations; no `down -v`, recursive deletion, disk formatting, or repartitioning.
- Bind the application to localhost by default.
- Use `docker compose config` before building or starting services.
- Treat a failed health check as a failed deployment and leave logs available for manual inspection.

## Validation

Local validation will run shell syntax checks, Compose model validation when Docker is available, and repository deployment tests. Server validation will confirm Docker/Compose versions, container health, `http://127.0.0.1:8080/health`, PostgreSQL readiness, Redis readiness, and that PostgreSQL/Redis have no host-published ports.
