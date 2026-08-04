import os
import shlex
import subprocess
import sys
from pathlib import Path

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"
REPO_URL = "https://github.com/cc1027114472/sub2api.git"
REMOTE_BASE = "/root/sub2api"
REMOTE_DEPLOY = "/root/sub2api/deploy"
LOCAL_DEPLOY = Path(__file__).resolve().parent.parent / "deploy"
CACHED_IMAGE = "sub2api-fork:git-1c67cc5e1f5b"


def run(ssh, cmd, check=True):
    stdin, stdout, stderr = ssh.exec_command(cmd, get_pty=True)
    out = stdout.read().decode("utf-8", "ignore")
    err = stderr.read().decode("utf-8", "ignore")
    code = stdout.channel.recv_exit_status()
    if out:
        sys.stdout.buffer.write(out.encode("utf-8", "replace"))
        sys.stdout.buffer.flush()
    if err:
        sys.stderr.buffer.write(err.encode("utf-8", "replace"))
        sys.stderr.buffer.flush()
    if check and code != 0:
        raise RuntimeError(f"command failed: {cmd}\nexit={code}")
    return code, out, err


def upload(ssh, local_path: Path, remote_path: str) -> None:
    sftp = ssh.open_sftp()
    try:
        sftp.put(str(local_path), remote_path)
    finally:
        sftp.close()


def main():
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(
        hostname=HOST,
        username=USER,
        password=PASSWORD,
        timeout=15,
        banner_timeout=15,
        auth_timeout=15,
    )

    try:
        print("== remote state ==")
        for cmd in [
            "whoami",
            "hostnamectl | sed -n '1,8p'",
            "docker --version",
            "docker compose version",
            "docker ps --format '{{.Names}} {{.Image}} {{.Ports}}'",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=False)

        print("\n== locate or clone repo ==")
        repo_root = REMOTE_BASE
        if run(
            ssh,
            f"test -d {shlex.quote(repo_root)} || git clone {shlex.quote(REPO_URL)} {shlex.quote(repo_root)}",
            check=True,
        ):
            pass

        print(f"\nrepo_root={repo_root}")

        print("\n== update repo ==")
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(repo_root)} && git fetch --all --prune && git checkout main && git pull --ff-only'",
            check=False,
        )

        print("\n== repo layout ==")
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(repo_root)} && find . -maxdepth 2 -type f | grep -E \"deploy/(server-install|docker-deploy|install)\\.sh|docker-compose\\.local\\.yml|docker-compose\\.fork\\.yml\" | sort'",
            check=False,
        )

        workdir = REMOTE_DEPLOY if run(ssh, f"test -d {shlex.quote(REMOTE_DEPLOY)}", check=False)[0] == 0 else repo_root
        print(f"\nworkdir={workdir}")

        print("\n== current container mounts ==")
        for cmd in [
            "docker inspect sub2api --format '{{range .Mounts}}{{.Source}} -> {{.Destination}}{{println}}{{end}}'",
            "docker inspect sub2api-postgres --format '{{range .Mounts}}{{.Source}} -> {{.Destination}}{{println}}{{end}}'",
            "docker inspect sub2api-redis --format '{{range .Mounts}}{{.Source}} -> {{.Destination}}{{println}}{{end}}'",
        ]:
            run(ssh, cmd, check=False)

        print("\n== stop old named containers ==")
        run(
            ssh,
            "docker rm -f sub2api sub2api-postgres sub2api-redis 2>/dev/null || true",
            check=False,
        )

        print("\n== ensure deploy files ==")
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(workdir)} && if [ ! -f .env ] || [ ! -f docker-compose.yml ]; then bash docker-deploy.sh; fi'",
            check=True,
        )

        print("\n== normalize low-memory build setting ==")
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(workdir)} && "
            "if grep -q \"^NODE_MAX_OLD_SPACE_SIZE=\" .env; then "
            "sed -i \"s/^NODE_MAX_OLD_SPACE_SIZE=.*/NODE_MAX_OLD_SPACE_SIZE=768/\" .env; "
            "else printf \"\\nNODE_MAX_OLD_SPACE_SIZE=768\\n\" >> .env; fi'",
            check=True,
        )

        for filename in [
            "docker-compose.fork.yml",
            "build-fork-image.sh",
            "backup.sh",
            "restore.sh",
            "update.sh",
            "rollback.sh",
            "healthcheck.sh",
        ]:
            local_file = LOCAL_DEPLOY / filename
            remote_file = f"{workdir}/{filename}"
            print(f"\n== upload {filename} ==")
            upload(ssh, local_file, remote_file)
            run(ssh, f"chmod +x {shlex.quote(remote_file)}", check=False)

        print("\n== choose image ==")
        have_cached, _, _ = run(ssh, f"docker image inspect {shlex.quote(CACHED_IMAGE)} >/dev/null 2>&1", check=False)
        use_cached = have_cached == 0
        if use_cached:
            print(f"using cached image: {CACHED_IMAGE}")
        else:
            print("\n== build fork image ==")
            run(
                ssh,
                f"bash -lc 'cd {shlex.quote(workdir)} && bash build-fork-image.sh'",
                check=True,
            )

        print("\n== start stack ==")
        image_env = ""
        if use_cached:
            image_env = f"SUB2API_IMAGE={shlex.quote(CACHED_IMAGE.split(':', 1)[0])} SUB2API_IMAGE_TAG={shlex.quote(CACHED_IMAGE.split(':', 1)[1])} "
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(workdir)} && {image_env}docker compose -f docker-compose.local.yml -f docker-compose.fork.yml up -d --remove-orphans'",
            check=True,
        )

        print("\n== wait for health ==")
        run(
            ssh,
            "bash -lc 'for i in $(seq 1 120); do curl -fsS http://127.0.0.1:8080/health >/dev/null && exit 0; sleep 2; done; exit 1'",
            check=True,
        )

        print("\n== verify backup ==")
        run(
            ssh,
            f"bash -lc 'cd {shlex.quote(workdir)} && COMPOSE_PROJECT_NAME=deploy BACKUP_ROOT=/root/sub2api/deploy/backups bash backup.sh >/tmp/sub2api-backup.log && tail -n 20 /tmp/sub2api-backup.log'",
            check=True,
        )

        print("\n== final state ==")
        for cmd in [
            "docker ps --format '{{.Names}} {{.Image}} {{.Ports}}'",
            "curl -fsS http://127.0.0.1:8080/health",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=False)
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
