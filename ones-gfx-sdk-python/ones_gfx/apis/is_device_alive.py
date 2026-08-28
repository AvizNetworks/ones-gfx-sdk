"""
GET /isDeviceAlive  ->  /api/fm/isDeviceAlive

Reachability probe.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def is_device_alive(*, hostIP: str) -> bool:
    client = get_client()
    path = "isDeviceAlive"
    params: dict[str, Any] = {}
    params["hostIP"] = hostIP
    return client.call_api("GET", path, params=params or None)
