"""
GET /fabrics/{fabricName}/vpcpeering  ->  FM direct (the gateway does not proxy this read)

One page of a fabric's VPC peerings, oldest first.

Returns `{"fabricName", "peerings", "pagination"}`. Each peering carries
`name`, `vpcs` (the pair, unordered: the leak is configured on both VRFs),
`origin` (`VPC_PEERING`, or `NETWORK_POLICY` for one a tenant-to-tenant
network policy created on a VXLAN fabric) and `createdAt`. A fabric with no
peerings returns an empty list; FM answers 404 `FABRIC_NOT_FOUND` for an
unknown fabric. Peerings made by a VRF-LITE network policy are not recorded
here; `list_network_policies` is their system of record.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def list_vpc_peerings(client: Client, fabricName: str, *, page: int | None = None, size: int | None = None) -> dict[str, Any]:
    """Peerings on `fabricName`. `page` is zero-based (default 0); `size`
    defaults to 50 and FM caps it at 500."""
    params: dict[str, Any] = {}
    if page is not None:
        params["page"] = page
    if size is not None:
        params["size"] = size
    return client.call_api("GET", f"fabrics/{fabricName}/vpcpeering", params=params or None)
