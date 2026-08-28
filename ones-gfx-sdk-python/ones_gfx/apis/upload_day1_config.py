"""
POST /uploadDay1Config  ->  /api/fm/uploadDay1Config

Upload a Day-1 intent file.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def upload_day1_config(*, file: str, x_request_origin: str | None = None) -> str:
    client = get_client()
    path = "uploadDay1Config"
    headers: dict[str, str] = {}
    if x_request_origin is not None:
        headers["x-request-origin"] = x_request_origin
    files: dict[str, Any] = {}
    files["file"] = open(file, "rb")
    return client.call_api("POST", path, files=files or None, extra_headers=headers or None)
