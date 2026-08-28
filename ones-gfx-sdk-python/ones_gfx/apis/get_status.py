"""
GET /status  ->  /api/fm/status

Read status lines from a file on disk.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_status(*, fileName: str | None = None) -> list[str] | None:
    client = get_client()
    path = "status"
    params: dict[str, Any] = {}
    if fileName is not None:
        params["fileName"] = fileName
    return client.call_api("GET", path, params=params or None)
