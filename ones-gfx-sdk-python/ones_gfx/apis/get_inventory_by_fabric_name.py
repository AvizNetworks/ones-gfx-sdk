"""
GET /getInventoryByFabricName/{name}  ->  /api/fm/getInventoryByFabricName/{name}

Inventory for one fabric.
"""
from __future__ import annotations


from .._types import InventoryRecord
from ..client import get_client


def get_inventory_by_fabric_name(name: str) -> list[InventoryRecord]:
    client = get_client()
    path = f"getInventoryByFabricName/{name}"
    return client.call_api("GET", path)
