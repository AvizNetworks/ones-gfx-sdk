"""
GET /getDevicesByLayer  ->  /api/fm/getDevicesByLayer

Device IPs for a layer.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_devices_by_layer(client: Client, *, layer: str) -> list[str]:
    path = "getDevicesByLayer"
    params: dict[str, Any] = {}
    params["layer"] = layer
    return client.call_api("GET", path, params=params or None)
