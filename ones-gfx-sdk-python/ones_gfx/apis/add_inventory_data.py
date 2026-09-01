"""
POST /addInventoryData  ->  /api/fm/addInventoryData

Add inventory rows for a fabric.
"""
from __future__ import annotations


from .._types import InventoryItem
from ..client import Client


def add_inventory_data(client: Client, *, items: list[InventoryItem]) -> str:
    path = "addInventoryData"
    return client.call_api("POST", path, json_body=items)
