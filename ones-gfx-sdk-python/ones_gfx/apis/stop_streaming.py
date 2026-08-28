"""
GET /stop  ->  /api/fm/stop

Stop streaming that file.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def stop_streaming(client: Client, *, filename: str) -> str:
    path = "stop"
    params: dict[str, Any] = {}
    params["filename"] = filename
    return client.call_api("GET", path, params=params or None)
