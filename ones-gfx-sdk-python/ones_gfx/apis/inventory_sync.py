"""
POST /fabrics/{fabricName}/inventorySync  ->  /api/fm/fabrics/{fabricName}/inventorySync

Force immediate UFM inventory sync.
"""
from __future__ import annotations


from ..client import get_client


def inventory_sync(fabricName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/inventorySync"
    return client.call_api("POST", path)
