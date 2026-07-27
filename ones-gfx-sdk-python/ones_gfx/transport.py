from __future__ import annotations

import json
from typing import Any
from urllib.parse import urljoin

import requests

from .auth import AuthBase
from .enums import OperationMode
from .exceptions import (
    APIError,
    AuthenticationError,
    BadRequestError,
    ConflictError,
    NotFoundError,
    ServerError,
    TransportError,
)


class Transport:
    def __init__(self, base_url: str, auth: AuthBase, verify_tls: bool | str = True, timeout: int = 30):
        if not base_url:
            raise ValueError("base_url is required.")
        self._base_url = base_url if base_url.endswith("/") else base_url + "/"
        self._auth = auth
        self._verify_tls = verify_tls
        self._timeout = timeout
        self._session = requests.Session()

    def get(self, path: str, timeout: int | None = None) -> Any:
        return self._request("GET", path, timeout=timeout)

    def post(self, path: str, json_body: dict | None = None, mode: OperationMode = OperationMode.SYNCHRONOUS, timeout: int | None = None) -> Any:
        return self._request("POST", path, json_body=json_body, mode=mode, timeout=timeout)

    def patch(self, path: str, json_body: dict | None = None, mode: OperationMode = OperationMode.SYNCHRONOUS, timeout: int | None = None) -> Any:
        return self._request("PATCH", path, json_body=json_body, mode=mode, timeout=timeout)

    def delete(self, path: str, mode: OperationMode = OperationMode.SYNCHRONOUS, timeout: int | None = None) -> Any:
        return self._request("DELETE", path, mode=mode, timeout=timeout)

    def close(self) -> None:
        self._session.close()
        self._auth.close()

    def _request(self, method: str, path: str, json_body: dict | None = None, mode: OperationMode = OperationMode.SYNCHRONOUS, timeout: int | None = None) -> Any:
        url = urljoin(self._base_url, path.lstrip("/"))
        if self._auth.needs_proactive_refresh():
            self._auth.refresh()
        response = self._send(method, url, json_body, mode, timeout)
        if response.status_code == 401:
            self._auth.refresh()
            response = self._send(method, url, json_body, mode, timeout)
        return self._handle_response(response)

    def _send(self, method: str, url: str, json_body: dict | None, mode: OperationMode, timeout: int | None) -> requests.Response:
        headers = {"Content-Type": "application/json", "Accept": "application/json"}
        self._auth.apply(headers)
        if mode is OperationMode.SYNCHRONOUS:
            headers["Prefer"] = "respond-sync"
        elif mode in (OperationMode.ASYNC_POLL, OperationMode.ASYNC_WEBHOOK):
            headers["Prefer"] = "respond-async"
        try:
            return self._session.request(
                method=method,
                url=url,
                headers=headers,
                json=json_body,
                timeout=self._timeout if timeout is None else timeout,
                verify=self._verify_tls,
            )
        except requests.exceptions.RequestException as e:
            raise TransportError(f"{method} {url} failed: {e}") from e

    def _handle_response(self, response: requests.Response) -> Any:
        body = self._parse_body(response)
        if 200 <= response.status_code < 300:
            return self._unwrap_envelope(body)
        message = self._extract_error_message(body) or f"HTTP {response.status_code}"
        exc_cls = self._classify_status(response.status_code)
        raise exc_cls(message, status_code=response.status_code, response_body=body)

    @staticmethod
    def _parse_body(response: requests.Response) -> Any:
        if not response.content:
            return None
        try:
            return response.json()
        except ValueError:
            return response.text

    @staticmethod
    def _unwrap_envelope(body: Any) -> Any:
        if not isinstance(body, dict):
            return body
        if "status" not in body or ("data" not in body and "message" not in body):
            return body
        if body.get("status") != "success":
            return body
        data = body.get("data")
        if isinstance(data, str):
            if data != "":
                try:
                    data = json.loads(data)
                except (json.JSONDecodeError, ValueError):
                    pass
            else:
                data = None  # empty string — fall through to message
        if data is None:
            return body.get("message")
        return data

    @staticmethod
    def _extract_error_message(body: Any) -> str | None:
        if not isinstance(body, dict):
            return str(body) if body else None
        for key in ("error", "message", "detail"):
            value = body.get(key)
            if value:
                return str(value)
        return None

    @staticmethod
    def _classify_status(status_code: int) -> type[APIError]:
        if status_code == 400:
            return BadRequestError
        if status_code in (401, 403):
            return AuthenticationError
        if status_code == 404:
            return NotFoundError
        if status_code == 409:
            return ConflictError
        if 500 <= status_code < 600:
            return ServerError
        return APIError
