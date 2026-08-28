"""
GET /getHostTenantsList/{fabricName}  ->  /api/fm/getHostTenantsList/{fabricName}

List all host-tenant records for a fabric.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class HosttenantsRecord(TypedDict, total=False):
    """Models/Hosttenants.java JSON (response)"""

    id: int
    name: str
    description: str
    hostsAllocated: int
    vniId: int
    fabricName: str
    config_status: str
    createdAt: str
    updatedAt: str


def get_host_tenants_list(client: Client, fabricName: str) -> list[HosttenantsRecord]:
    path = f"getHostTenantsList/{fabricName}"
    return client.call_api("GET", path)
