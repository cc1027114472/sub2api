import sys

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"
DEPLOY_DIR = "/root/sub2api/deploy"
CACHED_IMAGE = "sub2api-fork:git-1c67cc5e1f5b"


def safe_write(text: str, is_err: bool = False) -> None:
    stream = sys.stderr if is_err else sys.stdout
    data = text.encode("utf-8", "replace")
    stream.buffer.write(data)
    stream.buffer.flush()


def run(ssh: paramiko.SSHClient, cmd: str, check: bool = True) -> tuple[int, str, str]:
    stdin, stdout, stderr = ssh.exec_command(cmd, get_pty=True)
    out = stdout.read().decode("utf-8", "ignore")
    err = stderr.read().decode("utf-8", "ignore")
    code = stdout.channel.recv_exit_status()
    if out:
        safe_write(out)
    if err:
        safe_write(err, is_err=True)
    if check and code != 0:
        raise RuntimeError(f"command failed ({code}): {cmd}")
    return code, out, err


def main() -> None:
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(
        HOST,
        username=USER,
        password=PASSWORD,
        timeout=20,
        banner_timeout=20,
        auth_timeout=20,
    )

    try:
        print("== force sub2api bind to loopback ==")
        run(
            ssh,
            (
                f"bash -lc \"cd {DEPLOY_DIR} && "
                "if grep -q '^BIND_HOST=' .env; then "
                "sed -i 's/^BIND_HOST=.*/BIND_HOST=127.0.0.1/' .env; "
                "else printf '\\nBIND_HOST=127.0.0.1\\n' >> .env; fi && "
                f"SUB2API_IMAGE={CACHED_IMAGE.split(':', 1)[0]} "
                f"SUB2API_IMAGE_TAG={CACHED_IMAGE.split(':', 1)[1]} "
                "docker compose -f docker-compose.local.yml -f docker-compose.fork.yml up -d --remove-orphans --no-build\""
            ),
            check=True,
        )

        print("\n== verify bind and health ==")
        for cmd in [
            "docker ps --format '{{.Names}} {{.Ports}}'",
            "ss -lntp | grep -E ':(80|443|8080)\\b' || true",
            "curl -fsS http://127.0.0.1:8080/health",
            "curl -fsS https://ukapi.cc/health",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=True)
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
