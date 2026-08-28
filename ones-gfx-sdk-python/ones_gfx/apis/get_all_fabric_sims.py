"""
GET /getAllFabricSims  ->  /api/fm/getAllFabricSims

All simulations.
"""
from __future__ import annotations


from .._types import FabricSimItem
from ..client import get_client


def get_all_fabric_sims() -> list[FabricSimItem]:
    client = get_client()
    path = "getAllFabricSims"
    return client.call_api("GET", path)
