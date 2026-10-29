"""
GET /fabrics/{fabricName}/tenants/{tenantName}/vpcpeering  ->  FM direct (the gateway does not proxy this read)

One page of the VPC peerings a tenant takes part in, including its storage
peering.

Returns `{"fabricName", "tenantName", "vpc", "peerings", "pagination"}`.
Because the tenant is named, each peering reports the tenant's own `vpc` and
the `peerVpc` on the other side (instead of the unordered `vpcs` pair).
FM answers 404 `TENANT_NOT_FOUND` / `FABRIC_NOT_FOUND` when either is unknown.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def list_tenant_vpc_peerings(client: Client, fabricName: str, tenantName: str, *, page: int | None = None, size: int | None = None) -> dict[str, Any]:
    """Peerings that `tenantName` on `fabricName` takes part in. `page` is
    zero-based (default 0); `size` defaults to 50 and FM caps it at 500."""
    params: dict[str, Any] = {}
    if page is not None:
        params["page"] = page
    if size is not None:
        params["size"] = size
    return client.call_api("GET", f"fabrics/{fabricName}/tenants/{tenantName}/vpcpeering", params=params or None)
