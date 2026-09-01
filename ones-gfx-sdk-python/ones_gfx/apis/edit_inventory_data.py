"""
PATCH /fabrics/{fabricName}/editInventoryData  ->  /api/fm/fabrics/{fabricName}/editInventoryData

Partial inventory edit scoped to a fabric.
"""
from __future__ import annotations


from .._types import InventoryItem
from ..client import Client


def edit_inventory_data(client: Client, fabricName: str, *, items: list[InventoryItem]) -> str:
    path = f"fabrics/{fabricName}/editInventoryData"
    return client.call_api("PATCH", path, json_body=items)
