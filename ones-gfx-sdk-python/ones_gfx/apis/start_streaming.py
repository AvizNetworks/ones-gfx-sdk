"""
GET /start  ->  /api/fm/start

Mark a file for WebSocket streaming.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def start_streaming(client: Client, *, filename: str) -> str:
    path = "start"
    params: dict[str, Any] = {}
    params["filename"] = filename
    return client.call_api("GET", path, params=params or None)
