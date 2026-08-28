"""
POST /addInventoryData  ->  /api/fm/addInventoryData

Add inventory rows for a fabric.
"""
from __future__ import annotations


from .._types import InventoryItem
from ..client import get_client


def add_inventory_data(*, items: list[InventoryItem]) -> str:
    client = get_client()
    path = "addInventoryData"
    return client.call_api("POST", path, json_body=items)
