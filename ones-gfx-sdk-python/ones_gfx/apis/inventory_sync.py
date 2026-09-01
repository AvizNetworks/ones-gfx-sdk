"""
POST /fabrics/{fabricName}/inventorySync  ->  /api/fm/fabrics/{fabricName}/inventorySync

Force immediate UFM inventory sync.
"""
from __future__ import annotations


from ..client import Client


def inventory_sync(client: Client, fabricName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/inventorySync"
    return client.call_api("POST", path)
