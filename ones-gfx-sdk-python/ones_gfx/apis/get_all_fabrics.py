"""
GET /getAllFabrics  ->  /api/fm/getAllFabrics

All fabrics as Fabrics[].
"""
from __future__ import annotations


from .._types import FabricItem
from ..client import Client


def get_all_fabrics(client: Client) -> list[FabricItem]:
    path = "getAllFabrics"
    return client.call_api("GET", path)
