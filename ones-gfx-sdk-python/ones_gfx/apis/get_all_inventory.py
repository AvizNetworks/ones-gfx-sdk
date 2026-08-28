"""
GET /getAllInventory  ->  /api/fm/getAllInventory

All inventory rows.
"""
from __future__ import annotations


from .._types import InventoryRecord
from ..client import get_client


def get_all_inventory() -> list[InventoryRecord]:
    client = get_client()
    path = "getAllInventory"
    return client.call_api("GET", path)
