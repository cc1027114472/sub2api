import sys

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"
NEW_EMAIL = "admin@ukapi.cc"


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
        sql = (
            "update users set email = 'admin@ukapi.cc', updated_at = now() where id = 1; "
            "update auth_identities set provider_subject = 'admin@ukapi.cc', metadata = "
            "jsonb_set(metadata, '{email}', to_jsonb('admin@ukapi.cc'::text)) "
            "where user_id = 1 and provider_type = 'email';"
        )
        command = f'docker exec sub2api-postgres psql -U sub2api -d sub2api -v ON_ERROR_STOP=1 -c "{sql}"'
        stdin, stdout, stderr = ssh.exec_command(command, get_pty=True)
        out = stdout.read().decode("utf-8", "ignore")
        err = stderr.read().decode("utf-8", "ignore")
        code = stdout.channel.recv_exit_status()
        sys.stdout.buffer.write(out.encode("utf-8", "replace"))
        sys.stderr.buffer.write(err.encode("utf-8", "replace"))
        if code:
            raise SystemExit(code)
        print(f"admin email normalized to {NEW_EMAIL}")
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
