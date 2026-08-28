"""
Main client entry point for the ONES Spectrum-X SDK.

Single-token auth model: log in once, persist the token to a local
``secrets.json``, and attach it as a raw ``authorization`` header on every
subsequent request. The token is also exposed as the module-level global
``auth_token``, which is populated from ``secrets.json`` on import.
"""

from __future__ import annotations

import json
from collections.abc import Mapping
from pathlib import Path
from typing import Any

import requests

from ._types import AuthResponse

# Local file that persists the auth token between runs. Defined as a module
# constant so the location is easy to change.
SECRETS_FILE = Path("secrets.json")

# Global auth token, populated on startup from SECRETS_FILE and refreshed on
# every login/refresh. Sent as the raw ``authorization`` header.
auth_token: str | None = None


class NotAuthenticatedError(RuntimeError):
    """Raised when an authenticated API call is made before login()."""


def _load_token() -> str | None:
    """Populate the global ``auth_token`` from ``secrets.json`` if present.

    Tolerates a missing file or malformed JSON by leaving the token as None.
    """
    global auth_token
    try:
        data = json.loads(SECRETS_FILE.read_text())
        auth_token = data.get("auth_token")
    except (OSError, json.JSONDecodeError, ValueError):
        auth_token = None
    return auth_token


def _save_token(token: str | None) -> None:
    """Set the global token and persist it to ``secrets.json``.

    Merges with any existing keys in the file so unrelated data is preserved.
    """
    global auth_token
    auth_token = token
    data: dict[str, Any] = {}
    try:
        existing = json.loads(SECRETS_FILE.read_text())
        if isinstance(existing, dict):
            data = existing
    except (OSError, json.JSONDecodeError, ValueError):
        data = {}
    data["auth_token"] = token
    SECRETS_FILE.write_text(json.dumps(data, indent=2))


# Populate the global token on import.
_load_token()


class ONESClient:
    def __init__(self, base_url: str, verify_tls: bool | str = True):
        if not base_url:
            raise ValueError("base_url is required.")
        self.base_url = base_url.rstrip("/")
        self.verify_tls = verify_tls
        self._session = requests.Session()
        # Silence the InsecureRequestWarning when TLS verification is off.
        if verify_tls is False:
            requests.packages.urllib3.disable_warnings(
                requests.packages.urllib3.exceptions.InsecureRequestWarning
            )

    def _auth_headers(self) -> dict[str, str]:
        """Header dict carrying the raw token, or empty when unauthenticated."""
        return {"authorization": auth_token} if auth_token else {}

    def call_api(
        self,
        method: str,
        path: str,
        *,
        params: Mapping[str, Any] | None = None,
        json_body: Any = None,
        files: dict[str, Any] | None = None,
        data: dict[str, Any] | None = None,
        extra_headers: dict[str, str] | None = None,
    ) -> Any:
        """Call a Fabric Manager endpoint under ``/api/fm/``.

        This is the shared entry point used by every generated method in
        ``ones_gfx.apis``. It enforces the auth guard, prefixes the path with
        ``/api/fm/``, attaches the raw ``authorization`` header, and unwraps the
        standard ``{"data": ...}`` response envelope when present.

        Args:
            method: HTTP verb (``GET``, ``POST``, ``PATCH``, ``DELETE``, ...).
            path: Endpoint path as documented (e.g. ``getAllFabrics`` or
                ``fabrics/{fabricName}/tenants``); prefixed with ``/api/fm/``.
            params: Query string parameters (``None`` values are dropped by
                callers before this point).
            json_body: JSON request body.
            files: Multipart file parts, forwarded to ``requests``.
            data: Multipart/form field values, forwarded to ``requests``.
            extra_headers: Additional request headers (e.g. ``Prefer``).

        Raises:
            NotAuthenticatedError: if no auth token is set (login() not called).
        """
        if not auth_token:
            raise NotAuthenticatedError(
                "Not authenticated: call login() before calling Fabric Manager APIs."
            )
        url = f"{self.base_url}/api/fm/{path.lstrip('/')}"
        headers = self._auth_headers()
        if extra_headers:
            headers.update(extra_headers)
        response = self._session.request(
            method=method.upper(),
            url=url,
            params=params,
            json=json_body,
            files=files,
            data=data,
            headers=headers,
            verify=self.verify_tls,
        )
        response.raise_for_status()
        if not response.content:
            return None
        try:
            body = response.json()
        except ValueError:
            return response.text
        if isinstance(body, dict) and "data" in body:
            return body["data"]
        return body

    def login(self, username: str, password: str) -> AuthResponse:
        """Authenticate and persist the returned token.

        POST {base_url}/api/user/login with {"username", "password"}. Returns
        ``{"data": {"message", "token", "isPwdResetNeeded"}}``; stores
        ``data.token`` and returns the full parsed response.
        """
        response = self._session.post(
            f"{self.base_url}/api/user/login",
            json={"username": username, "password": password},
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        token = body.get("data", {}).get("token")
        if token:
            _save_token(token)
        return body

    def refresh(self) -> AuthResponse:
        """Refresh the auth token and persist it.

        POST {base_url}/api/user/refresh with the current token in the
        ``authorization`` header. Returns ``{"data": {"message", "token"}}`` —
        the same shape as login, but with message "Token refreshed" and no
        ``isPwdResetNeeded``. Stores the new ``data.token`` in the same place as
        login (global + secrets.json) and returns the full parsed response.
        """
        response = self._session.post(
            f"{self.base_url}/api/user/refresh",
            headers=self._auth_headers(),
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        token = body.get("data", {}).get("token")
        if token:
            _save_token(token)
        return body

    def logout(self) -> Any:
        """Log out server-side, then clear the local token.

        POST {base_url}/api/logout with the auth header, then blank the global
        and the ``auth_token`` field in ``secrets.json`` regardless of outcome.
        """
        try:
            response = self._session.post(
                f"{self.base_url}/api/user/logout",
                headers=self._auth_headers(),
                verify=self.verify_tls,
            )
            body = response.json() if response.content else None
        finally:
            _save_token(None)
        return body

    def getTeleDevices(self) -> Any:
        """Fetch the inventory devices list.

        GET {base_url}/api/inventory/Devices with the auth header. Returns the
        ``data`` field of the response when present, otherwise the full body.
        """
        response = self._session.get(
            f"{self.base_url}/api/inventory/Devices",
            headers=self._auth_headers(),
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        if isinstance(body, dict) and "data" in body:
            return body["data"]
        return body


# Shared client used by the ones_gfx.apis functions. Configure it once at
# startup with configure(); the api modules import get_client() and use it.
_client: ONESClient | None = None


def configure(base_url: str, verify_tls: bool | str = True) -> ONESClient:
    """Create and register the shared client used by ones_gfx.apis."""
    global _client
    _client = ONESClient(base_url, verify_tls=verify_tls)
    return _client


def get_client() -> ONESClient:
    """Return the shared client, or raise if configure() was never called."""
    if _client is None:
        raise RuntimeError(
            "Client not configured: call ones_gfx.client.configure(base_url) first."
        )
    return _client
