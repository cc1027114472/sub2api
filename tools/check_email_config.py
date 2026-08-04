import sys

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"


def run(ssh: paramiko.SSHClient, cmd: str) -> tuple[int, str, str]:
    stdin, stdout, stderr = ssh.exec_command(cmd, get_pty=True)
    out = stdout.read().decode("utf-8", "ignore")
    err = stderr.read().decode("utf-8", "ignore")
    code = stdout.channel.recv_exit_status()
    return code, out, err


def write_block(title: str, code: int, out: str, err: str) -> None:
    sys.stdout.buffer.write(f"\n== {title} [exit {code}] ==\n".encode("utf-8"))
    if out:
        sys.stdout.buffer.write(out.encode("utf-8", "replace"))
        if not out.endswith("\n"):
            sys.stdout.buffer.write(b"\n")
    if err:
        sys.stdout.buffer.write(err.encode("utf-8", "replace"))
        if not err.endswith("\n"):
            sys.stdout.buffer.write(b"\n")
    sys.stdout.buffer.flush()


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
        checks = [
            (
                "deploy env mail vars",
                r"grep -n -i -E 'resend|smtp|mail|email' /root/sub2api/deploy/.env 2>/dev/null || true",
            ),
            (
                "container env mail vars",
                r"docker exec sub2api /bin/sh -lc ""env | grep -i -E 'resend|smtp|mail|email' || true""",
            ),
            (
                "config yaml mail vars",
                r"grep -n -i -E 'smtp_|mail|email' /root/sub2api/deploy/config.yaml 2>/dev/null | sed -n '1,120p' || true",
            ),
        ]
        for title, cmd in checks:
            code, out, err = run(ssh, cmd)
            write_block(title, code, out, err)
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
