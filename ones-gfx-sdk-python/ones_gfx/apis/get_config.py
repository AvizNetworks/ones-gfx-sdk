"""
GET /getconfig  ->  /api/fm/getconfig

Current device config.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_config(client: Client, *, deviceip: str) -> object:
    path = "getconfig"
    params: dict[str, Any] = {}
    params["deviceip"] = deviceip
    return client.call_api("GET", path, params=params or None)
