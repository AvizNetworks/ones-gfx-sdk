"""
GET /fabrics  ->  /api/fm/fabrics

Fabrics as FabricDto, grouped in a map.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class FabricsListResponse(TypedDict):
    """GET /fabrics"""

    fabrics: list[FabricDtoItem]


class FabricDtoItem(TypedDict, total=False):
    """Cumulus/dto/FabricDto.java"""

    id: int
    fabricName: str
    description: str
    numOfSUs: int
    maxNumOfSUs: int
    ewTenantAware: bool
    nsTenantAware: bool
    defaultStorageName: str
    cnpq: str
    createdAt: str
    updatedAt: str


def get_fabrics() -> FabricsListResponse:
    client = get_client()
    path = "fabrics"
    return client.call_api("GET", path)
