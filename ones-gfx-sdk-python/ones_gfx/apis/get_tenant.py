"""
GET /fabrics/{fabricName}/tenants/{tenantName}  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}

Tenant detail incl. GPU assignment.
"""
from __future__ import annotations


from ..client import Client


def get_tenant(client: Client, fabricName: str, tenantName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/tenants/{tenantName}"
    return client.call_api("GET", path)
