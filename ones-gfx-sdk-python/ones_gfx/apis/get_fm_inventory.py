"""
GET /fm-Inventory  ->  /api/fm/fm-Inventory

FM's own device inventory view.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class DeviceInventoryItem(TypedDict, total=False):
    """Helper/DeviceInventory.java"""

    ipAddress: str
    hostname: str
    layer: str
    fabricName: str


def get_fm_inventory(client: Client) -> list[DeviceInventoryItem]:
    path = "fm-Inventory"
    return client.call_api("GET", path)
