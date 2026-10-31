"""
GET /fabrics/{fabricName}/vpcpeering/{peeringName}  ->  FM direct (the gateway does not proxy this read)

One VPC peering by name.

FM wraps it as `{"fabricName", "peering"}`; this returns the inner `peering`
(`name`, `vpcs`, `origin`, `createdAt`). FM answers 404 `VPC_PEERING_NOT_FOUND`
for an unknown peering and 404 `FABRIC_NOT_FOUND` for an unknown fabric.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_vpc_peering(client: Client, fabricName: str, peeringName: str) -> dict[str, Any]:
    """Peering `peeringName` on `fabricName`."""
    data = client.call_api("GET", f"fabrics/{fabricName}/vpcpeering/{peeringName}")
    if isinstance(data, dict) and "peering" in data:
        return data["peering"]
    return data
