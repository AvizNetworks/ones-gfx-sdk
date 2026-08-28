"""
POST /addFabricSim  ->  /api/fm/addFabricSim

Register a simulation.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def add_fabric_sim(*, id: int | None = None, fabricName: str | None = None, simulationId: str | None = None, username: str | None = None, token: str | None = None, orgUuid: str | None = None, status: str | None = None, uiLink: str | None = None) -> str:
    client = get_client()
    path = "addFabricSim"
    body: dict[str, Any] = {}
    if id is not None:
        body["id"] = id
    if fabricName is not None:
        body["fabricName"] = fabricName
    if simulationId is not None:
        body["simulationId"] = simulationId
    if username is not None:
        body["username"] = username
    if token is not None:
        body["token"] = token
    if orgUuid is not None:
        body["orgUuid"] = orgUuid
    if status is not None:
        body["status"] = status
    if uiLink is not None:
        body["uiLink"] = uiLink
    return client.call_api("POST", path, json_body=body or None)
