"""
GET /getAllFabricSims  ->  /api/fm/getAllFabricSims

All simulations.
"""
from __future__ import annotations


from .._types import FabricSimItem
from ..client import Client


def get_all_fabric_sims(client: Client) -> list[FabricSimItem]:
    path = "getAllFabricSims"
    return client.call_api("GET", path)
