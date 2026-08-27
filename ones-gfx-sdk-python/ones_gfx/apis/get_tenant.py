"""
GET /fabrics/{fabricName}/tenants/{tenantName}  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}

Tenant detail incl. GPU assignment.
"""
from __future__ import annotations

from typing import Any

from ..client import ONESClient


def get_tenant(client: ONESClient, fabricName: str, tenantName: str) -> Any:
    path = f"fabrics/{fabricName}/tenants/{tenantName}"
    return client.call_api("GET", path)
