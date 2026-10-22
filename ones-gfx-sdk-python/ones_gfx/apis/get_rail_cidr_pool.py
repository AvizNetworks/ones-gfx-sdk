"""
GET /fabrics/{fabricName}/rails/{railName}/cidrpool  ->  /api/fm/fabrics/{fabricName}/rails/{railName}/cidrpool

The rail CIDRPool document for a single rail.

Content negotiated exactly like `get_rail_cidr_pools`: JSON by default, or the
raw manifest text when `accept` asks for YAML or plain text.

404 with `INVALID_FABRIC` when the fabric is unknown, or `INVALID_RAIL` when
the rail is not one of its rails. `railName` is the rail's CIDRPool name.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_rail_cidr_pool(
    client: Client, fabricName: str, railName: str, *, accept: str | None = None
) -> Any:
    """The CIDRPool document for `railName` on `fabricName`."""
    path = f"fabrics/{fabricName}/rails/{railName}/cidrpool"
    headers: dict[str, str] = {}
    if accept is not None:
        headers["Accept"] = accept
    return client.call_api("GET", path, extra_headers=headers or None)
