"""开启公开注册，并要求新用户完成邮箱验证。"""

import json
import os
import sys
import urllib.error
import urllib.request


BASE_URL = os.environ.get("SUB2API_URL", "https://ukapi.cc").strip().rstrip("/")
ADMIN_EMAIL = os.environ.get("SUB2API_ADMIN_EMAIL", "").strip()
ADMIN_PASSWORD = os.environ.get("SUB2API_ADMIN_PASSWORD", "").strip()


def request_json(path: str, method: str, payload: dict | None = None, token: str = "") -> dict:
    data = None
    headers = {"Accept": "application/json", "User-Agent": "sub2api-deploy-script"}
    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"
    if token:
        headers["Authorization"] = f"Bearer {token}"

    request = urllib.request.Request(
        f"{BASE_URL}/api/v1{path}",
        data=data,
        headers=headers,
        method=method,
    )
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            return json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", "replace")
        raise RuntimeError(f"{method} {path} failed with HTTP {exc.code}: {body}") from exc


def main() -> None:
    if not ADMIN_EMAIL or not ADMIN_PASSWORD:
        raise RuntimeError(
            "请设置 SUB2API_ADMIN_EMAIL 和 SUB2API_ADMIN_PASSWORD 环境变量"
        )

    login = request_json(
        "/auth/login",
        "POST",
        {"email": ADMIN_EMAIL, "password": ADMIN_PASSWORD},
    )
    login_data = login.get("data", login)
    access_token = login_data.get("access_token")
    if not access_token:
        raise RuntimeError("管理员登录未返回 access token")

    updated_response = request_json(
        "/admin/settings",
        "PUT",
        {
            "registration_enabled": True,
            "email_verify_enabled": True,
        },
        access_token,
    )
    updated = updated_response.get("data", updated_response)
    print(
        json.dumps(
            {
                "registration_enabled": updated.get("registration_enabled"),
                "email_verify_enabled": updated.get("email_verify_enabled"),
            },
            ensure_ascii=False,
        )
    )


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        raise
