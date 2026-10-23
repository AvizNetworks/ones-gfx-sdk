"""
GET /fabrics/{fabricName}/computes/{serverName}/ip-addresses  ->  /api/fm/fabrics/{fabricName}/computes/{serverName}/ip-addresses

Every IP allocation of a compute server, across tenants.
"""
from __future__ import annotations

from .._types import ServerIpAllocations
from ..client import Client
from ._validation import path_segment


def get_server_ips(client: Client, fabricName: str, serverName: str) -> ServerIpAllocations:
    """Return ``{"fabricName", "serverName", "allocations"}`` (empty list when none).

    Raises:
        ValueError: a name is blank.
        requests.HTTPError: 404 ``SERVER_NOT_FOUND`` (not in the fabric inventory).
    """
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/computes/{path_segment(serverName, 'serverName')}/ip-addresses"
    return client.call_api("GET", path)
