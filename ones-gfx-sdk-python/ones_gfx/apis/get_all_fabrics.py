"""
GET /getAllFabrics  ->  /api/fm/getAllFabrics

All fabrics as Fabrics[].
"""
from __future__ import annotations

from typing import Any

from ..client import ONESClient


def get_all_fabrics(client: ONESClient) -> Any:
    path = "getAllFabrics"
    return client.call_api("GET", path)
