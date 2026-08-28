"""
POST /addHostTenantData  ->  /api/fm/addHostTenantData

Register a host-based tenant record.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def add_host_tenant_data(client: Client, *, name: str, fabricName: str, description: str | None = None, hostsAllocated: int | None = None, vniId: int | None = None, config_status: str | None = None) -> str:
    path = "addHostTenantData"
    body: dict[str, Any] = {}
    body["name"] = name
    body["fabricName"] = fabricName
    if description is not None:
        body["description"] = description
    if hostsAllocated is not None:
        body["hostsAllocated"] = hostsAllocated
    if vniId is not None:
        body["vniId"] = vniId
    if config_status is not None:
        body["config_status"] = config_status
    return client.call_api("POST", path, json_body=body or None)
