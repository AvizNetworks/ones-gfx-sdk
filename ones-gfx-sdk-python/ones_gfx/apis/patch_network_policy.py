"""
PATCH /fabrics/{fabricName}/networkpolicy/{policyName}  ->  /api/fm/fabrics/{fabricName}/networkpolicy/{policyName}

Add endpoints to an existing network policy, leaving the ones already on the
devices untouched.

Only endpoints may change: `type`, `tenantName` and `serviceName` on each side
must match the existing policy exactly. The endpoints you send are the NEW
ones to add, not the full set. A tenant-to-tenant policy carries no endpoints
and is rejected.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def patch_network_policy(
    client: Client,
    fabricName: str,
    policyName: str,
    *,
    source: dict[str, Any],
    destination: dict[str, Any],
) -> dict[str, Any]:
    """Add endpoints to `policyName` on `fabricName`.

    `source` and `destination` mirror :func:`create_network_policy`, with
    `endpoints` holding only the additions.

    Returns `{"policyName": ...}` on success.
    """
    path = f"fabrics/{fabricName}/networkpolicy/{policyName}"
    # `policy.name` is echoed into the body as well as the path: FM validates
    # a patch with the same check it uses for a create, which rejects a body
    # without it ("policy.name is required.").
    body: dict[str, Any] = {
        "policy": {"name": policyName},
        "spec": {"source": source, "destination": destination},
    }
    return client.call_api("PATCH", path, json_body=body)
