"""
POST /fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign  ->  /api/fm/...

Assign a DPU VF interface to a tenant.

This is how a host joins any tenant, including the shared-services tenant.
FM requires the tenant to exist and the server to be in it first. Returns 200,
or 207 Multi-Status with `ipAssignStatus` = IP_ASSIGN_FAILED when the VF is
assigned but its IP setup failed. The `ip*` fields appear only when IPAM or
DHCP is enabled for the fabric.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def assign_vf_interface(client: Client, fabricName: str, serverName: str, vfId: str, *, tenantName: str, idempotency_key: str | None = None) -> dict[str, Any]:
    """Assign VF `vfId` on `serverName` to `tenantName`.

    `vfId` is the DPU interface name (`pf0vf0_if`) or the host interface
    name (`enp193s0f1np1`).
    """
    path = f"fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign"
    headers: dict[str, str] = {}
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    return client.call_api("POST", path, json_body={"tenantName": tenantName}, extra_headers=headers or None)
