"""
GET /getFabricByName/{name}  ->  /api/fm/getFabricByName/{name}

Single fabric.
"""
from __future__ import annotations


from .._types import FabricItem
from ..client import get_client


def get_fabric_by_name(name: str) -> FabricItem:
    client = get_client()
    path = f"getFabricByName/{name}"
    return client.call_api("GET", path)
