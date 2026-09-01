"""
GET /getInventoryByFabricName/{name}  ->  /api/fm/getInventoryByFabricName/{name}

Inventory for one fabric.
"""
from __future__ import annotations


from .._types import InventoryRecord
from ..client import Client


def get_inventory_by_fabric_name(client: Client, name: str) -> list[InventoryRecord]:
    path = f"getInventoryByFabricName/{name}"
    return client.call_api("GET", path)
