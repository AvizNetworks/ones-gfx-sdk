"""
GET /getHostsList/{fabricName}  ->  /api/fm/getHostsList/{fabricName}

Get list of host records registered in a fabric.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class HostItem(TypedDict, total=False):
    """Models/Hosts.java JSON"""

    id: int
    hostname: str
    hostStatus: str
    tenantName: str
    fabricName: str
    config_status: str
    createdAt: str
    updatedAt: str


def get_hosts_list(fabricName: str) -> list[HostItem]:
    client = get_client()
    path = f"getHostsList/{fabricName}"
    return client.call_api("GET", path)
