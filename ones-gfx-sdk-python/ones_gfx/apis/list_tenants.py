"""
GET /fabrics/{fabricName}/tenants  ->  /api/fm/fabrics/{fabricName}/tenants

List tenants in a fabric.
"""
from __future__ import annotations


from ..client import Client


def list_tenants(client: Client, fabricName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/tenants"
    return client.call_api("GET", path)
