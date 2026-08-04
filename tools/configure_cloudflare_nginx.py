import sys
from pathlib import Path

import paramiko


HOST = "154.19.229.193"
USER = "root"
PASSWORD = "2Q99EAjlnb"
DOMAIN = "ukapi.cc"
LOCAL_CERT = Path(r"C:\Users\Administrator\Desktop\sub2\原证书.txt")
LOCAL_KEY = Path(r"C:\Users\Administrator\Desktop\sub2\私钥.txt")
REMOTE_CERT_DIR = "/etc/nginx/ssl/sub2api"
REMOTE_CERT = f"{REMOTE_CERT_DIR}/{DOMAIN}.crt"
REMOTE_KEY = f"{REMOTE_CERT_DIR}/{DOMAIN}.key"
NGINX_CONF = f"/etc/nginx/sites-available/{DOMAIN}.conf"
NGINX_LINK = f"/etc/nginx/sites-enabled/{DOMAIN}.conf"


CF_REAL_IP_SNIPPET = """\
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2a06:98c0::/29;
set_real_ip_from 2c0f:f248::/32;
real_ip_header CF-Connecting-IP;
"""


NGINX_SITE = f"""\
server {{
    listen 80;
    listen [::]:80;
    server_name {DOMAIN} www.{DOMAIN};

    return 301 https://{DOMAIN}$request_uri;
}}

server {{
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name {DOMAIN};

    ssl_certificate {REMOTE_CERT};
    ssl_certificate_key {REMOTE_KEY};
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;
    client_max_body_size 64m;

    {CF_REAL_IP_SNIPPET.rstrip()}

    location / {{
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
    }}
}}

server {{
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name www.{DOMAIN};

    ssl_certificate {REMOTE_CERT};
    ssl_certificate_key {REMOTE_KEY};

    return 301 https://{DOMAIN}$request_uri;
}}
"""


def run(ssh: paramiko.SSHClient, cmd: str, check: bool = True) -> tuple[int, str, str]:
    stdin, stdout, stderr = ssh.exec_command(cmd, get_pty=True)
    out = stdout.read().decode("utf-8", "ignore")
    err = stderr.read().decode("utf-8", "ignore")
    code = stdout.channel.recv_exit_status()
    if out:
        print(out, end="")
    if err:
        print(err, end="", file=sys.stderr)
    if check and code != 0:
        raise RuntimeError(f"command failed ({code}): {cmd}")
    return code, out, err


def upload_text(ssh: paramiko.SSHClient, text: str, remote_path: str, mode: int) -> None:
    sftp = ssh.open_sftp()
    try:
        with sftp.file(remote_path, "w") as fh:
            fh.write(text)
        sftp.chmod(remote_path, mode)
    finally:
        sftp.close()


def upload_file(ssh: paramiko.SSHClient, local_path: Path, remote_path: str, mode: int) -> None:
    if not local_path.exists():
        raise FileNotFoundError(local_path)
    sftp = ssh.open_sftp()
    try:
        sftp.put(str(local_path), remote_path)
        sftp.chmod(remote_path, mode)
    finally:
        sftp.close()


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
        print("== inspect current state ==")
        for cmd in [
            "hostnamectl | sed -n '1,8p'",
            "docker ps --format '{{.Names}} {{.Image}} {{.Ports}}'",
            "curl -fsS http://127.0.0.1:8080/health",
            "ss -lntp | sed -n '1,40p'",
            "nginx -v 2>&1 || true",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=False)

        print("\n== install nginx and prerequisites ==")
        run(ssh, "apt-get update", check=True)
        run(ssh, "DEBIAN_FRONTEND=noninteractive apt-get install -y nginx", check=True)
        run(ssh, "systemctl enable nginx", check=True)

        print("\n== upload certificate and key ==")
        run(ssh, f"install -d -m 700 {REMOTE_CERT_DIR}", check=True)
        upload_file(ssh, LOCAL_CERT, REMOTE_CERT, 0o600)
        upload_file(ssh, LOCAL_KEY, REMOTE_KEY, 0o600)

        print("\n== write nginx site config ==")
        upload_text(ssh, NGINX_SITE, NGINX_CONF, 0o644)
        run(ssh, f"ln -sfn {NGINX_CONF} {NGINX_LINK}", check=True)
        run(ssh, "rm -f /etc/nginx/sites-enabled/default", check=False)

        print("\n== validate and reload nginx ==")
        run(ssh, "nginx -t", check=True)
        run(ssh, "systemctl restart nginx", check=True)

        print("\n== verify ports and local https ==")
        for cmd in [
            "ss -lntp | grep -E ':(80|443|8080)\\b' || true",
            f"curl -I -H 'Host: {DOMAIN}' http://127.0.0.1",
            f"curl -kI --resolve {DOMAIN}:443:127.0.0.1 https://{DOMAIN}",
            f"curl -kfsS --resolve {DOMAIN}:443:127.0.0.1 https://{DOMAIN}/health",
            f"openssl x509 -in {REMOTE_CERT} -noout -subject -issuer -dates",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=True)

        print("\n== external verification ==")
        for cmd in [
            f"curl -I https://{DOMAIN}",
            f"curl -fsS https://{DOMAIN}/health",
        ]:
            print(f"\n$ {cmd}")
            run(ssh, cmd, check=True)
    finally:
        ssh.close()


if __name__ == "__main__":
    main()
