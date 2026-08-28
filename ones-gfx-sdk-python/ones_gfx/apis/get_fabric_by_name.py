"""
GET /getFabricByName/{name}  ->  /api/fm/getFabricByName/{name}

Single fabric.
"""
from __future__ import annotations


from .._types import FabricItem
from ..client import Client


def get_fabric_by_name(client: Client, name: str) -> FabricItem:
    path = f"getFabricByName/{name}"
    return client.call_api("GET", path)
