"""
POST /updateFabricSimStatus  ->  /api/fm/updateFabricSimStatus

Update simulation status.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def update_fabric_sim_status(client: Client, *, name: str, status: str | None = None) -> str:
    path = "updateFabricSimStatus"
    body: dict[str, Any] = {}
    body["name"] = name
    if status is not None:
        body["status"] = status
    return client.call_api("POST", path, json_body=body or None)
