"""
GET /fabrics/{fabricName}/computes/{serverName}/ip-addresses  ->  /api/fm/fabrics/{fabricName}/computes/{serverName}/ip-addresses

Every IP IPAM has allocated to one compute server, across all its tenants.

Unpaginated, unlike the per-tenant listing: a server holds one address per
tenant it belongs to. FM answers 404 (`SERVER_NOT_FOUND`) when the server is
not in the fabric's inventory.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_server_ips(client: Client, fabricName: str, serverName: str) -> dict[str, Any]:
    """IP allocations for `serverName`, as `{"fabricName", "serverName", "allocations"}`."""
    path = f"fabrics/{fabricName}/computes/{serverName}/ip-addresses"
    return client.call_api("GET", path)
