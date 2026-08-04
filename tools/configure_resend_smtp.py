import json
import os
import sys
import urllib.error
import urllib.request


BASE_URL = os.environ.get("SUB2API_URL", "https://ukapi.cc").strip().rstrip("/")
ADMIN_EMAIL = os.environ.get("SUB2API_ADMIN_EMAIL", "").strip()
ADMIN_PASSWORD = os.environ.get("SUB2API_ADMIN_PASSWORD", "").strip()
RESEND_API_KEY = os.environ.get("RESEND_API_KEY", "").strip()
TEST_RECIPIENT = os.environ.get("RESEND_TEST_RECIPIENT", "delivered@resend.dev").strip()
FROM_EMAIL = os.environ.get("RESEND_FROM_EMAIL", "openapi@ukapi.cc").strip()
FROM_NAME = os.environ.get("RESEND_FROM_NAME", "ukapi.cc").strip()

# Resend's SMTP credentials: username is the literal "resend", password is the API key.
SMTP_CONFIG = {
    "smtp_host": "smtp.resend.com",
    "smtp_port": 587,
    "smtp_username": "resend",
    "smtp_password": RESEND_API_KEY,
    "smtp_from_email": FROM_EMAIL,
    "smtp_from_name": FROM_NAME,
    "smtp_use_tls": True,
}


def request_json(path: str, method: str, payload: dict | None = None, token: str = "") -> dict:
    data = None
    headers = {
        "Accept": "application/json",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/131.0 Safari/537.36",
    }
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
    missing = [
        name
        for name, value in {
            "SUB2API_ADMIN_EMAIL": ADMIN_EMAIL,
            "SUB2API_ADMIN_PASSWORD": ADMIN_PASSWORD,
            "RESEND_API_KEY": RESEND_API_KEY,
        }.items()
        if not value
    ]
    if missing:
        raise RuntimeError("missing environment variables: " + ", ".join(missing))

    print("== login ==")
    print(f"login email: {ADMIN_EMAIL!r}")
    login = request_json(
        "/auth/login",
        "POST",
        {"email": ADMIN_EMAIL, "password": ADMIN_PASSWORD},
    )
    login_data = login.get("data", login)
    access_token = login_data.get("access_token")
    if not access_token:
        raise RuntimeError("login did not return an access token")
    print("admin login successful")

    print("\n== update smtp settings ==")
    update_payload = {
        **SMTP_CONFIG,
        "frontend_url": BASE_URL,
    }
    updated_response = request_json("/admin/settings", "PUT", update_payload, access_token)
    updated = updated_response.get("data", updated_response)
    print(
        "saved:",
        {
            "smtp_host": updated.get("smtp_host"),
            "smtp_port": updated.get("smtp_port"),
            "smtp_username": updated.get("smtp_username"),
            "smtp_password_configured": updated.get("smtp_password_configured"),
            "smtp_from_email": updated.get("smtp_from_email"),
            "smtp_use_tls": updated.get("smtp_use_tls"),
            "frontend_url": updated.get("frontend_url"),
        },
    )

    print("\n== test smtp connection ==")
    connection_response = request_json(
        "/admin/settings/test-smtp",
        "POST",
        {
            "smtp_host": SMTP_CONFIG["smtp_host"],
            "smtp_port": SMTP_CONFIG["smtp_port"],
            "smtp_username": SMTP_CONFIG["smtp_username"],
            "smtp_password": SMTP_CONFIG["smtp_password"],
            "smtp_use_tls": SMTP_CONFIG["smtp_use_tls"],
        },
        access_token,
    )
    connection = connection_response.get("data", connection_response)
    print(connection)

    print("\n== send test email ==")
    sent_response = request_json(
        "/admin/settings/send-test-email",
        "POST",
        {
            **SMTP_CONFIG,
            "email": TEST_RECIPIENT,
        },
        access_token,
    )
    sent = sent_response.get("data", sent_response)
    print(sent)
    print(f"test recipient: {TEST_RECIPIENT}")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        raise
