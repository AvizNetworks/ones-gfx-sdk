"""
GET /getFabricSimByName/{name}  ->  /api/fm/getFabricSimByName/{name}

Simulation by fabric name.
"""
from __future__ import annotations


from .._types import FabricSimItem
from ..client import get_client


def get_fabric_sim_by_name(name: str) -> FabricSimItem:
    client = get_client()
    path = f"getFabricSimByName/{name}"
    return client.call_api("GET", path)
