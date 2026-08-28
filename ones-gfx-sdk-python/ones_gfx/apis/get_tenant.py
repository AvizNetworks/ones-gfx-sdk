"""
GET /fabrics/{fabricName}/tenants/{tenantName}  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}

Tenant detail incl. GPU assignment.
"""
from __future__ import annotations


from ..client import get_client


def get_tenant(fabricName: str, tenantName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/tenants/{tenantName}"
    return client.call_api("GET", path)
