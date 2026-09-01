"""
GET /fabrics/{fabricName}/inventoryPorts  ->  /api/fm/fabrics/{fabricName}/inventoryPorts

IB ports grouped by host.
"""
from __future__ import annotations


from ..client import Client


def get_inventory_ports(client: Client, fabricName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/inventoryPorts"
    return client.call_api("GET", path)
