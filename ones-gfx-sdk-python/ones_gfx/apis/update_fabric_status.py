"""
POST /updateFabricStatus  ->  /api/fm/updateFabricStatus

Update fabric status/intent/deployment type.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def update_fabric_status(client: Client, *, name: str, status: str | None = None, intent: str | None = None, description: str | None = None, deploymentType: str | None = None) -> str:
    path = "updateFabricStatus"
    body: dict[str, Any] = {}
    body["name"] = name
    if status is not None:
        body["status"] = status
    if intent is not None:
        body["intent"] = intent
    if description is not None:
        body["description"] = description
    if deploymentType is not None:
        body["deploymentType"] = deploymentType
    return client.call_api("POST", path, json_body=body or None)
