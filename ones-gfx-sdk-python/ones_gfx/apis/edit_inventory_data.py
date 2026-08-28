"""
PATCH /fabrics/{fabricName}/editInventoryData  ->  /api/fm/fabrics/{fabricName}/editInventoryData

Partial inventory edit scoped to a fabric.
"""
from __future__ import annotations


from .._types import InventoryItem
from ..client import get_client


def edit_inventory_data(fabricName: str, *, items: list[InventoryItem]) -> str:
    client = get_client()
    path = f"fabrics/{fabricName}/editInventoryData"
    return client.call_api("PATCH", path, json_body=items)
