"""
POST /updateHosts  ->  /api/fm/updateHosts

Bulk add or remove hosts for a tenant.
"""
from __future__ import annotations

from typing import Any, Literal

from ..client import Client

HostAction = Literal["ADD", "DELETE"]


def update_hosts(client: Client, *, hostnames: list[str], hostAction: HostAction, tenantName: str, fabricName: str) -> bool:
    path = "updateHosts"
    body: dict[str, Any] = {}
    body["hostnames"] = hostnames
    body["hostAction"] = hostAction
    body["tenantName"] = tenantName
    body["fabricName"] = fabricName
    return client.call_api("POST", path, json_body=body or None)
