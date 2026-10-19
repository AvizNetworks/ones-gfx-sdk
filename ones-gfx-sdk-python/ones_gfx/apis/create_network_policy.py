"""
POST /fabrics/{fabricName}/networkpolicy  ->  /api/fm/fabrics/{fabricName}/networkpolicy

Create a network policy: allow traffic from one source to one destination.

Works on VXLAN (Cumulus, SONiC and DPU) and VRF-LITE (Cumulus and SONiC).
Config is pushed to the border/combined leaves and DPUs, so FM runs a device
licence check first and answers 403 when it fails.
"""
from __future__ import annotations

from typing import Any

from ..client import Client

#: Accepted values for a source/destination ``type`` (FM matches case-insensitively).
NETWORK_POLICY_TYPES = ("tenant", "service")


def create_network_policy(
    client: Client,
    fabricName: str,
    *,
    policyName: str,
    source: dict[str, Any],
    destination: dict[str, Any],
) -> dict[str, Any]:
    """Create network policy `policyName` on `fabricName`.

    `policyName` is at most 24 characters. `source` and `destination` are each
    `{"type", "tenantName"}` plus an optional `"serviceName"` and a list of
    `"endpoints"`; `type` is `"tenant"` or `"service"`.

    Returns `{"policyName": ...}` on success.
    """
    path = f"fabrics/{fabricName}/networkpolicy"
    body: dict[str, Any] = {
        "policy": {"name": policyName},
        "spec": {"source": source, "destination": destination},
    }
    return client.call_api("POST", path, json_body=body)
