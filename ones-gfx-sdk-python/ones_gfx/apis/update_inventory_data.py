"""
POST /updateInventoryData  ->  /api/fm/updateInventoryData

Update credentials/hostname per device.
"""
from __future__ import annotations


from .._types import FabricInventoryUpdateItem
from ..client import Client


def update_inventory_data(client: Client, *, items: list[FabricInventoryUpdateItem]) -> str:
    path = "updateInventoryData"
    return client.call_api("POST", path, json_body=items)
