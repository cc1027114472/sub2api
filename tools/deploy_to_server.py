import os
import subprocess
import sys
import time
from pathlib import Path
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
LOCAL_ROOT = Path(__file__).resolve().parent.parent

def main():
    print("=== Step 1: Create Local Git Bundle of Latest Commits ===")
    bundle_path = LOCAL_ROOT / "sub2api_diff.bundle"
    cmd = ["git", "bundle", "create", str(bundle_path), "HEAD"]
    print("Running:", " ".join(cmd))
    subprocess.check_call(cmd, cwd=LOCAL_ROOT)
    print(f"Bundle created at {bundle_path} ({bundle_path.stat().st_size} bytes)")

    print(f"\nConnecting to {HOST} as {USER}...")
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
        print("\n=== Step 2: Upload Bundle to Server ===")
        sftp = ssh.open_sftp()
        remote_bundle = "/tmp/sub2api_diff.bundle"
        print(f"Uploading {bundle_path} -> {remote_bundle}")
        sftp.put(str(bundle_path), remote_bundle)
        sftp.close()
        print("Upload complete.")

        print("\n=== Step 3: Apply Bundle to Server Git Repo ===")
        run_cmd(f"cd {WORK_DIR} && git bundle verify {remote_bundle}")
        run_cmd(f"cd {WORK_DIR} && git fetch {remote_bundle} HEAD && git checkout -B main FETCH_HEAD")
        run_cmd(f"cd {WORK_DIR} && git log -n 3 --oneline")
        run_cmd(f"rm -f {remote_bundle}")

        print("\n=== Step 4: Ensure Environment Settings & Tags ===")
        # Get target revision
        _, rev_out = run_cmd(f"cd {WORK_DIR} && git rev-parse --short=12 HEAD")
        revision = rev_out.strip()
        print(f"Target Git Revision: {revision}")

        # Update SUB2API_IMAGE_TAG and ensure NODE_MAX_OLD_SPACE_SIZE
        run_cmd(f"""
            cd {DEPLOY_DIR} && \
            if grep -q "^SUB2API_IMAGE_TAG=" .env; then \
                sed -i "s/^SUB2API_IMAGE_TAG=.*/SUB2API_IMAGE_TAG=git-{revision}/" .env; \
            else \
                echo "SUB2API_IMAGE_TAG=git-{revision}" >> .env; \
            fi && \
            if grep -q "^NODE_MAX_OLD_SPACE_SIZE=" .env; then \
                sed -i "s/^NODE_MAX_OLD_SPACE_SIZE=.*/NODE_MAX_OLD_SPACE_SIZE=2560/" .env; \
            else \
                echo "NODE_MAX_OLD_SPACE_SIZE=2560" >> .env; \
            fi
        """)

        # Make all scripts executable
        run_cmd(f"cd {DEPLOY_DIR} && chmod +x *.sh")

        print("\n=== Step 5: Build Fork Image ===")
        run_cmd(f"cd {DEPLOY_DIR} && COMPOSE_PROJECT_NAME=deploy bash build-fork-image.sh")

        print("\n=== Step 6: Stop Existing Containers and Start Stack ===")
        run_cmd("docker rm -f sub2api sub2api-postgres sub2api-redis 2>/dev/null || true")
        run_cmd(f"cd {DEPLOY_DIR} && COMPOSE_PROJECT_NAME=deploy docker compose -f docker-compose.local.yml -f docker-compose.fork.yml up -d --remove-orphans")

        print("\n=== Step 7: Healthcheck Polling ===")
        run_cmd("bash -c 'for i in $(seq 1 60); do if curl -fsS http://127.0.0.1:8080/health >/dev/null 2>&1; then echo \"Health check passed on attempt $i!\"; exit 0; fi; echo \"Waiting for service to be healthy... ($i/60)\"; sleep 2; done; echo \"Health check timed out\"; exit 1'")

        print("\n=== Step 8: Verify Remote Services & Endpoint ===")
        run_cmd("docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'")
        run_cmd("curl -s http://127.0.0.1:8080/health")
        run_cmd("curl -s -k --resolve ukapi.cc:443:127.0.0.1 https://ukapi.cc/health")
        run_cmd("docker logs --tail 25 sub2api")

        print("\n==========================================")
        print("  Redeployment Completed Successfully!  ")
        print("==========================================")

    finally:
        ssh.close()
        if bundle_path.exists():
            bundle_path.unlink()

if __name__ == "__main__":
    main()
