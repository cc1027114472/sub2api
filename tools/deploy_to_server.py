import os
import sys
import time
import paramiko

# Ensure UTF-8 output even in Windows cmd/pwsh
if sys.stdout.encoding != 'utf-8':
    try:
        sys.stdout.reconfigure(encoding='utf-8', errors='replace')
        sys.stderr.reconfigure(encoding='utf-8', errors='replace')
    except Exception:
        pass

HOST = "154.23.162.32"
USER = "root"
PASSWORD = "2Q99EAjlnb"
WORK_DIR = "/root/sub2api"
DEPLOY_DIR = "/root/sub2api/deploy"

def main():
    print(f"Connecting to {HOST} as {USER}...")
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(
        HOST,
        username=USER,
        password=PASSWORD,
        allow_agent=False,
        look_for_keys=False,
        timeout=20,
    )

    def run_cmd(cmd, check=True):
        print(f"\n>>> Running: {cmd}")
        stdin, stdout, stderr = ssh.exec_command(cmd, get_pty=True)
        out_chunks = []
        while not stdout.channel.exit_status_ready():
            if stdout.channel.recv_ready():
                chunk = stdout.channel.recv(4096).decode("utf-8", "ignore")
                try:
                    sys.stdout.write(chunk)
                    sys.stdout.flush()
                except Exception:
                    sys.stdout.buffer.write(chunk.encode("utf-8", "replace"))
                    sys.stdout.buffer.flush()
                out_chunks.append(chunk)
            time.sleep(0.05)
        
        # Read remaining
        rest = stdout.read().decode("utf-8", "ignore")
        if rest:
            try:
                sys.stdout.write(rest)
                sys.stdout.flush()
            except Exception:
                sys.stdout.buffer.write(rest.encode("utf-8", "replace"))
                sys.stdout.buffer.flush()
            out_chunks.append(rest)

        code = stdout.channel.recv_exit_status()
        if check and code != 0:
            raise RuntimeError(f"Command failed with exit code {code}: {cmd}")
        return code, "".join(out_chunks)

    try:
        print("=== Step 1: Ensure Adequate Swap Space ===")
        run_cmd("""
            if [ ! -f /swapfile2 ]; then
                echo "Creating 2GB extra swapfile..."
                fallocate -l 2G /swapfile2 || dd if=/dev/zero of=/swapfile2 bs=1M count=2048
                chmod 600 /swapfile2
                mkswap /swapfile2
                swapon /swapfile2
            else
                swapon /swapfile2 2>/dev/null || true
            fi
            swapon --show
            free -h
        """)

        print("\n=== Step 2: Ensure Environment File & Settings ===")
        # Restore .env if needed
        run_cmd(f"if [ ! -f {DEPLOY_DIR}/.env ] && [ -f /root/sub2api_env_backup/.env.bak ]; then cp /root/sub2api_env_backup/.env.bak {DEPLOY_DIR}/.env; fi")
        
        # Set Node heap limit to 2560MB so vue-tsc completes cleanly
        run_cmd(f"""
            cd {DEPLOY_DIR} && \
            if grep -q "^NODE_MAX_OLD_SPACE_SIZE=" .env; then \
                sed -i "s/^NODE_MAX_OLD_SPACE_SIZE=.*/NODE_MAX_OLD_SPACE_SIZE=2560/" .env; \
            else \
                echo "NODE_MAX_OLD_SPACE_SIZE=2560" >> .env; \
            fi
        """)

        # Pull latest code from GitHub origin/main
        run_cmd(f"cd {WORK_DIR} && git fetch origin main && git checkout main && git reset --hard origin/main")

        # Get latest short revision
        _, rev_out = run_cmd(f"cd {WORK_DIR} && git rev-parse --short=12 HEAD")
        revision = rev_out.strip()
        print(f"Target Git Revision: {revision}")

        run_cmd(f"""
            cd {DEPLOY_DIR} && \
            if grep -q "^SUB2API_IMAGE_TAG=" .env; then \
                sed -i "s/^SUB2API_IMAGE_TAG=.*/SUB2API_IMAGE_TAG=git-{revision}/" .env; \
            else \
                echo "SUB2API_IMAGE_TAG=git-{revision}" >> .env; \
            fi
        """)

        # Make all scripts executable
        run_cmd(f"cd {DEPLOY_DIR} && chmod +x *.sh")

        print("\n=== Step 3: Build Fork Image ===")
        run_cmd(f"cd {DEPLOY_DIR} && COMPOSE_PROJECT_NAME=deploy bash build-fork-image.sh")

        print("\n=== Step 4: Stop Existing Containers and Start Stack ===")
        run_cmd("docker rm -f sub2api sub2api-postgres sub2api-redis 2>/dev/null || true")
        run_cmd(f"cd {DEPLOY_DIR} && COMPOSE_PROJECT_NAME=deploy docker compose -f docker-compose.local.yml -f docker-compose.fork.yml up -d --remove-orphans")

        print("\n=== Step 5: Healthcheck Polling ===")
        run_cmd("bash -c 'for i in $(seq 1 60); do if curl -fsS http://127.0.0.1:8080/health >/dev/null 2>&1; then echo \"Health check passed on attempt $i!\"; exit 0; fi; echo \"Waiting for service to be healthy... ($i/60)\"; sleep 2; done; echo \"Health check timed out\"; exit 1'")

        print("\n=== Step 6: Verify Remote Services & Endpoint ===")
        run_cmd("docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'")
        run_cmd("curl -sI http://127.0.0.1:8080/health")
        run_cmd("curl -sI -k --resolve ukapi.cc:443:127.0.0.1 https://ukapi.cc/health || true")
        run_cmd("docker logs --tail 30 sub2api")

        print("\n==========================================")
        print("  Deployment Completed Successfully!  ")
        print("==========================================")

    finally:
        ssh.close()

if __name__ == "__main__":
    main()
