"""
GET /fabrics/{fabricName}/networkpolicy  ->  /api/fm/fabrics/{fabricName}/networkpolicy

Every network policy configured on the fabric.

The body is not the usual status/message envelope: it is
`{"status", "fabricName", "count", "policies"}`, and `policies` holds the
same `{"policy", "spec"}` shape `create_network_policy` accepts. FM answers
404 when the fabric is unknown.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def list_network_policies(client: Client, fabricName: str) -> dict[str, Any]:
    """List the network policies on `fabricName`."""
    path = f"fabrics/{fabricName}/networkpolicy"
    return client.call_api("GET", path)
