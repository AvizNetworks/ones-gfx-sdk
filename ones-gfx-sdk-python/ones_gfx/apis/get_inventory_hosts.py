"""
GET /fabrics/{fabricName}/inventoryHosts  ->  /api/fm/fabrics/{fabricName}/inventoryHosts

Live host list from UFM.
"""
from __future__ import annotations


from ..client import Client


def get_inventory_hosts(client: Client, fabricName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/inventoryHosts"
    return client.call_api("GET", path)
