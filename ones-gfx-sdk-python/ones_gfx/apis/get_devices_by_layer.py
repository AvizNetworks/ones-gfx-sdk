"""
GET /getDevicesByLayer  ->  /api/fm/getDevicesByLayer

Device IPs for a layer.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_devices_by_layer(*, layer: str) -> list[str]:
    client = get_client()
    path = "getDevicesByLayer"
    params: dict[str, Any] = {}
    params["layer"] = layer
    return client.call_api("GET", path, params=params or None)
