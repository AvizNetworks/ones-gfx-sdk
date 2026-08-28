"""
GET /getAllInventory  ->  /api/fm/getAllInventory

All inventory rows.
"""
from __future__ import annotations


from .._types import InventoryRecord
from ..client import Client


def get_all_inventory(client: Client) -> list[InventoryRecord]:
    path = "getAllInventory"
    return client.call_api("GET", path)
