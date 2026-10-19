"""
GET /fabrics/{fabricName}/tenants/{tenantName}/ip-addresses  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/ip-addresses

One page of the tenant's IPAM pool allocations, with the pool's own counters.

Returns `{"fabricName", "tenantName", "pool", "allocations", "pagination"}`,
where `pool` carries `cidr`, `allocatableAddressCount` and `allocatedCount`.
FM answers 404 (`IPAM_NOT_CONFIGURED`) when the tenant has no pool.
"""
from __future__ import annotations

from typing import Any

from ..client import Client

#: Values accepted by the ``status`` filter.
IP_ASSIGN_STATUSES = ("ALLOCATED_ONLY", "ASSIGNED", "IP_ASSIGN_FAILED")


def get_tenant_ips(
    client: Client,
    fabricName: str,
    tenantName: str,
    *,
    page: int | None = None,
    size: int | None = None,
    status: str | None = None,
) -> dict[str, Any]:
    """IP allocations for `tenantName`.

    `page` is zero-based and defaults to 0; `size` defaults to 50. `status`
    filters by assignment state — see `IP_ASSIGN_STATUSES`.
    """
    path = f"fabrics/{fabricName}/tenants/{tenantName}/ip-addresses"
    params: dict[str, Any] = {}
    if page is not None:
        params["page"] = page
    if size is not None:
        params["size"] = size
    if status is not None:
        params["status"] = status
    return client.call_api("GET", path, params=params or None)
