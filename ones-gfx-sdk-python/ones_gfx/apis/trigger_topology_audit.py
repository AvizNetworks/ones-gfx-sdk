"""
POST /fabrics/{fabricName}/topology-audit  ->  /api/fm/fabrics/{fabricName}/topology-audit

Kick off an async LLDP-vs-intent Topology Audit (GFX-6219/GFX-6225) for a
fabric. There is no synchronous variant - Fabric Manager always runs this in
the background. Poll the returned operationId with get_topology_audit_status()
or get_topology_audit_result().
"""
from __future__ import annotations

from .._types import OperationAccepted
from ..client import Client


def trigger_topology_audit(client: Client, fabricName: str) -> OperationAccepted:
    """Trigger a topology audit for `fabricName`.

    FM answers 404 (raised as NotFoundError) when the fabric or its intent
    cannot be resolved.
    """
    path = f"fabrics/{fabricName}/topology-audit"
    return client.call_api("POST", path)
