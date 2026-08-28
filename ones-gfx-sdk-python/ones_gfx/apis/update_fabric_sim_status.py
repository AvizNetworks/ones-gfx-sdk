"""
POST /updateFabricSimStatus  ->  /api/fm/updateFabricSimStatus

Update simulation status.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def update_fabric_sim_status(*, name: str, status: str | None = None) -> str:
    client = get_client()
    path = "updateFabricSimStatus"
    body: dict[str, Any] = {}
    body["name"] = name
    if status is not None:
        body["status"] = status
    return client.call_api("POST", path, json_body=body or None)
