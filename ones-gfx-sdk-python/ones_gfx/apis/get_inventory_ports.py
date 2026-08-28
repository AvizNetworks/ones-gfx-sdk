"""
GET /fabrics/{fabricName}/inventoryPorts  ->  /api/fm/fabrics/{fabricName}/inventoryPorts

IB ports grouped by host.
"""
from __future__ import annotations


from ..client import get_client


def get_inventory_ports(fabricName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/inventoryPorts"
    return client.call_api("GET", path)
