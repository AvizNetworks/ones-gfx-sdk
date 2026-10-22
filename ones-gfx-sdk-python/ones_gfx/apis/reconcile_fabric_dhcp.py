"""
POST /fabrics/{fabricName}/dhcp/reconcile  ->  /api/fm/fabrics/{fabricName}/dhcp/reconcile

Re-apply every tenant DHCP subnet so Kea matches Fabric Manager's records.

The same work the scheduled reconciliation job does, for an operator who does
not want to wait for the next run. Only the tenant-owned part of the Kea
config is rewritten — bootstrap subnets are never touched. `fabricName` scopes
the log line, not the work: the re-apply covers every tenant FM knows.

404 when the fabric is unknown, 400 on a validation failure, 502 when Kea
cannot be reached.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def reconcile_fabric_dhcp(client: Client, fabricName: str) -> dict[str, Any]:
    """Reconcile tenant DHCP subnets into Kea. Returns the re-apply summary."""
    path = f"fabrics/{fabricName}/dhcp/reconcile"
    return client.call_api("POST", path)
