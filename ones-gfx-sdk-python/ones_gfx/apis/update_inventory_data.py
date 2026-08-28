"""
POST /updateInventoryData  ->  /api/fm/updateInventoryData

Update credentials/hostname per device.
"""
from __future__ import annotations


from .._types import FabricInventoryUpdateItem
from ..client import get_client


def update_inventory_data(*, items: list[FabricInventoryUpdateItem]) -> str:
    client = get_client()
    path = "updateInventoryData"
    return client.call_api("POST", path, json_body=items)
