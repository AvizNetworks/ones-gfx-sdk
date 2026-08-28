"""
GET /getFabricSimByName/{name}  ->  /api/fm/getFabricSimByName/{name}

Simulation by fabric name.
"""
from __future__ import annotations


from .._types import FabricSimItem
from ..client import Client


def get_fabric_sim_by_name(client: Client, name: str) -> FabricSimItem:
    path = f"getFabricSimByName/{name}"
    return client.call_api("GET", path)
