"""
GET /fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address

IP allocation of one compute server within a tenant.
"""
from __future__ import annotations

from .._types import IpAllocation
from ..client import Client
from ._validation import path_segment


def get_tenant_compute_ip(client: Client, fabricName: str, tenantName: str, serverName: str) -> IpAllocation:
    """Return the allocation row (``ipAddress`` and ``prefixLength`` are separate fields).

    Raises:
        ValueError: a name is blank.
        requests.HTTPError: 404 ``IPAM_NOT_CONFIGURED`` (tenant has no IPAM pool)
            or ``SERVER_NOT_ALLOCATED_TO_TENANT`` (server has no allocation).
    """
    path = (
        f"fabrics/{path_segment(fabricName, 'fabricName')}"
        f"/tenants/{path_segment(tenantName, 'tenantName')}"
        f"/computes/{path_segment(serverName, 'serverName')}/ip-address"
    )
    return client.call_api("GET", path)
