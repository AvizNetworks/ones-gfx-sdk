"""
GET /fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address

The single IP allocated to one server within one tenant.

FM answers 404 either way when there is nothing to return, distinguished by
the error code: `IPAM_NOT_CONFIGURED` when the tenant has no pool at all,
`SERVER_NOT_ALLOCATED_TO_TENANT` when it has one but this server holds no
address in it.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_tenant_compute_ip(
    client: Client, fabricName: str, tenantName: str, serverName: str
) -> dict[str, Any]:
    """The IP allocation for `serverName` within `tenantName`."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address"
    return client.call_api("GET", path)
