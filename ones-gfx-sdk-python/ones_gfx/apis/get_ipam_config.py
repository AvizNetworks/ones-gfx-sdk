"""
GET /fabrics/{fabricName}/ipam/config  ->  /api/fm/fabrics/{fabricName}/ipam/config

The fabric's IPAM configuration.

Returns the same shape `update_ipam_config` accepts:
`{"fabricName", "enableIpam", "static": {...}, "dhcp": {...}}`, so a body read
here can be sent straight back. 404 when the fabric is unknown.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_ipam_config(client: Client, fabricName: str) -> dict[str, Any]:
    """The IPAM configuration for `fabricName`."""
    path = f"fabrics/{fabricName}/ipam/config"
    return client.call_api("GET", path)
