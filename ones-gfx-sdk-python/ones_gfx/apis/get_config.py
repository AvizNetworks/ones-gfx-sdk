"""
GET /getconfig  ->  /api/fm/getconfig

Current device config.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_config(*, deviceip: str) -> object:
    client = get_client()
    path = "getconfig"
    params: dict[str, Any] = {}
    params["deviceip"] = deviceip
    return client.call_api("GET", path, params=params or None)
