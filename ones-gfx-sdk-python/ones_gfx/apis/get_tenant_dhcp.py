"""
GET /fabrics/{fabricName}/tenants/{tenantName}/dhcp  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/dhcp

The tenant's DHCP subnet, whether Kea is actually serving it, and a page of
the leases handed out in it.

Returns `{"subnet", "liveInKea", "leases", "pagination"}`. `liveInKea` is
reported separately from the subnet's stored status on purpose: the two
disagreeing is exactly the drift reconciliation exists to repair. It comes
back None, alongside a `leaseCheckError`, when Kea could not be queried.

Leases are ordered by IP address, so pages are stable while the lease set is.
404 (`DHCP_SUBNET_NOT_FOUND`) when the fabric does not assign addresses by
DHCP, or the tenant has no VLAN or IPAM pool.
"""
from __future__ import annotations

from typing import Any

from ..client import Client

#: Largest lease page FM will return.
MAX_DHCP_LEASE_PAGE_SIZE = 1000


def get_tenant_dhcp(
    client: Client,
    fabricName: str,
    tenantName: str,
    *,
    page: int | None = None,
    size: int | None = None,
) -> dict[str, Any]:
    """The DHCP subnet and leases for `tenantName`.

    `page` is zero-based and defaults to 0; `size` defaults to 50 and is
    capped at `MAX_DHCP_LEASE_PAGE_SIZE`.
    """
    path = f"fabrics/{fabricName}/tenants/{tenantName}/dhcp"
    params: dict[str, Any] = {}
    if page is not None:
        params["page"] = page
    if size is not None:
        params["size"] = size
    return client.call_api("GET", path, params=params or None)
