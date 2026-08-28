"""
GET /fm-Inventory  ->  /api/fm/fm-Inventory

FM's own device inventory view.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class DeviceInventoryItem(TypedDict, total=False):
    """Helper/DeviceInventory.java"""

    ipAddress: str
    hostname: str
    layer: str
    fabricName: str


def get_fm_inventory() -> list[DeviceInventoryItem]:
    client = get_client()
    path = "fm-Inventory"
    return client.call_api("GET", path)
