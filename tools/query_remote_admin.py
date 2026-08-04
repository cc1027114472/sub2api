import sys

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"


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
        command = (
            "docker exec sub2api-postgres "
            "psql -U sub2api -d sub2api -Atc "
            "\"select id, email, role from users where role = chr(97)||chr(100)||chr(109)||chr(105)||chr(110) order by id; "
            "select provider_type, provider_key, provider_subject from auth_identities where user_id = 1\""
        )
        stdin, stdout, stderr = ssh.exec_command(command, get_pty=True)
        out = stdout.read().decode("utf-8", "ignore")
        err = stderr.read().decode("utf-8", "ignore")
        code = stdout.channel.recv_exit_status()
        sys.stdout.buffer.write(out.encode("utf-8", "replace"))
        sys.stderr.buffer.write(err.encode("utf-8", "replace"))
        if code:
            raise SystemExit(code)
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
