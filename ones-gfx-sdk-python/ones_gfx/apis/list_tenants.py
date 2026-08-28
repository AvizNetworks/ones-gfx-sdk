"""
GET /fabrics/{fabricName}/tenants  ->  /api/fm/fabrics/{fabricName}/tenants

List tenants in a fabric.
"""
from __future__ import annotations


from ..client import get_client


def list_tenants(fabricName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/tenants"
    return client.call_api("GET", path)
