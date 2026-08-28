"""
POST /updateroleinfo  ->  /api/fm/updateroleinfo

Update role layer name.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def update_role_info(*, layer: int, current_name: str) -> str:
    client = get_client()
    path = "updateroleinfo"
    body: dict[str, Any] = {}
    body["layer"] = layer
    body["current_name"] = current_name
    return client.call_api("POST", path, json_body=body or None)
