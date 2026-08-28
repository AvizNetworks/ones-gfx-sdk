"""
GET /fabrics/{fabricName}/inventoryHosts  ->  /api/fm/fabrics/{fabricName}/inventoryHosts

Live host list from UFM.
"""
from __future__ import annotations


from ..client import get_client


def get_inventory_hosts(fabricName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/inventoryHosts"
    return client.call_api("GET", path)
