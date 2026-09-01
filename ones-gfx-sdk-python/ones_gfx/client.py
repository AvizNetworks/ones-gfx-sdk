"""
Main client entry point for the ONES Spectrum-X SDK.

Single-token auth model, held per client instance: log in (or supply a token),
and it is attached as a raw ``authorization`` header on every subsequent
request. Nothing is stored globally or on disk — construct as many clients as
you need, for different hosts or different users.

Pass the client explicitly into the ``ones_gfx.apis`` functions::

    from ones_gfx import Client
    from ones_gfx.apis import get_all_fabrics

    client = Client.initialize_with_creds("https://host:3002", "admin", "secret")
    fabrics = get_all_fabrics(client)
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any

import requests

from ._types import AuthResponse


class NotAuthenticatedError(RuntimeError):
    """Raised when an authenticated API call is made without a token."""


class Client:
    """A configured connection to one ONES instance.

    Holds the base URL, the auth token, and (optionally) the credentials used to
    obtain it. Prefer the :meth:`initialize_with_creds` / :meth:`initialize_with_token`
    factories over the constructor.

    Args:
        verify_tls: strictness of TLS certificate checking. ``False`` (the
            default) accepts self-signed certificates, which is how ONES is
            normally deployed. ``True`` enforces strict verification. A string
            is treated as a path to a CA bundle to verify against.
    """

    def __init__(
        self,
        baseUrl: str,
        authToken: str | None = None,
        username: str | None = None,
        password: str | None = None,
        verify_tls: bool | str = False,
    ) -> None:
        if not baseUrl:
            raise ValueError("baseUrl is required.")
        self.baseUrl = baseUrl.rstrip("/")
        self.authToken = authToken
        self.username = username
        self.password = password
        self.verify_tls = verify_tls
        self._session = requests.Session()
        # Silence the InsecureRequestWarning when TLS verification is off.
        if verify_tls is False:
            requests.packages.urllib3.disable_warnings(
                requests.packages.urllib3.exceptions.InsecureRequestWarning
            )

    # ------------------------------------------------------------------
    # Factories
    # ------------------------------------------------------------------

    @classmethod
    def initialize_with_creds(
        cls,
        baseUrl: str,
        username: str,
        password: str,
        verify_tls: bool | str = False,
    ) -> "Client":
        """Build a client, log in, and store the returned token.

        Raises whatever ``login`` raises if authentication fails, so a returned
        client is always usable.

        ``verify_tls`` defaults to ``False`` (self-signed certificates
        accepted); pass ``True`` for strict checking or a CA-bundle path.
        """
        client = cls(
            baseUrl=baseUrl,
            username=username,
            password=password,
            verify_tls=verify_tls,
        )
        client.login()
        return client

    @classmethod
    def initialize_with_token(
        cls,
        baseUrl: str,
        token: str,
        verify_tls: bool | str = False,
    ) -> "Client":
        """Build a client from an existing token.

        ``verify_tls`` defaults to ``False`` (self-signed certificates
        accepted); pass ``True`` for strict checking or a CA-bundle path.

        Raises:
            ValueError: if ``token`` is missing or empty.
        """
        if not token:
            raise ValueError("token is required when initializing with a token.")
        return cls(baseUrl=baseUrl, authToken=token, verify_tls=verify_tls)

    # ------------------------------------------------------------------
    # Getters / setters
    # ------------------------------------------------------------------

    def get_base_url(self) -> str:
        return self.baseUrl

    def set_base_url(self, baseUrl: str) -> None:
        if not baseUrl:
            raise ValueError("baseUrl is required.")
        self.baseUrl = baseUrl.rstrip("/")

    def get_auth_token(self) -> str | None:
        return self.authToken

    def set_auth_token(self, authToken: str | None) -> None:
        self.authToken = authToken

    def get_username(self) -> str | None:
        return self.username

    def set_username(self, username: str | None) -> None:
        self.username = username

    def get_password(self) -> str | None:
        return self.password

    def set_password(self, password: str | None) -> None:
        self.password = password

    # ------------------------------------------------------------------
    # Token management
    # ------------------------------------------------------------------

    def update_token(self, token: str) -> None:
        """Replace the stored token.

        Unlike :meth:`set_auth_token` this rejects an empty value, so it is the
        safer entry point when rotating a token.
        """
        if not token:
            raise ValueError("token is required.")
        self.authToken = token

    def is_authenticated(self) -> bool:
        """True when a token is held."""
        return bool(self.authToken)

    def _auth_headers(self) -> dict[str, str]:
        """Header dict carrying the raw token, or empty when unauthenticated."""
        return {"authorization": self.authToken} if self.authToken else {}

    # ------------------------------------------------------------------
    # Session endpoints (/api/user/...)
    # ------------------------------------------------------------------

    def login(self, username: str | None = None, password: str | None = None) -> AuthResponse:
        """Authenticate and store the returned token.

        POST {baseUrl}/api/user/login with {"username", "password"}, falling back
        to the credentials already on the client. Returns
        ``{"data": {"message", "token", "isPwdResetNeeded"}}``.
        """
        user = username if username is not None else self.username
        pwd = password if password is not None else self.password
        if not user or not pwd:
            raise ValueError("username and password are required to log in.")
        # Remember them so refresh/re-login works without re-supplying.
        self.username = user
        self.password = pwd

        response = self._session.post(
            f"{self.baseUrl}/api/user/login",
            json={"username": user, "password": pwd},
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        token = body.get("data", {}).get("token")
        if token:
            self.authToken = token
        return body

    def refresh_auth(self) -> AuthResponse:
        """Exchange the current token for a new one and store it.

        POST {baseUrl}/api/user/refresh with the current token in the
        ``authorization`` header. Returns ``{"data": {"message", "token"}}`` —
        the same shape as login, with message "Token refreshed".

        Raises:
            NotAuthenticatedError: if there is no current token to exchange.
        """
        if not self.authToken:
            raise NotAuthenticatedError(
                "Not authenticated: no token to refresh. Call login() first."
            )
        response = self._session.post(
            f"{self.baseUrl}/api/user/refresh",
            headers=self._auth_headers(),
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        token = body.get("data", {}).get("token")
        if token:
            self.update_token(token)
        return body

    def logout(self) -> Any:
        """Log out server-side, then clear the stored token.

        POST {baseUrl}/api/user/logout with the auth header. The local token is
        cleared regardless of the server's response.
        """
        try:
            response = self._session.post(
                f"{self.baseUrl}/api/user/logout",
                headers=self._auth_headers(),
                verify=self.verify_tls,
            )
            body = response.json() if response.content else None
        finally:
            self.authToken = None
        return body

    # ------------------------------------------------------------------
    # Request plumbing
    # ------------------------------------------------------------------

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

        This is the shared entry point used by every function in
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
            NotAuthenticatedError: if no auth token is set.
        """
        if not self.authToken:
            raise NotAuthenticatedError(
                "Not authenticated: call login() before calling Fabric Manager APIs."
            )
        url = f"{self.baseUrl}/api/fm/{path.lstrip('/')}"
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

    def get_tele_devices(self) -> Any:
        """Fetch the inventory devices list.

        GET {baseUrl}/api/inventory/Devices with the auth header. Returns the
        ``data`` field of the response when present, otherwise the full body.
        """
        response = self._session.get(
            f"{self.baseUrl}/api/inventory/Devices",
            headers=self._auth_headers(),
            verify=self.verify_tls,
        )
        response.raise_for_status()
        body = response.json()
        if isinstance(body, dict) and "data" in body:
            return body["data"]
        return body

    def close(self) -> None:
        """Release the underlying HTTP session."""
        self._session.close()

    def __enter__(self) -> "Client":
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.close()


# Backwards-compatible alias for the previous class name.
ONESClient = Client
