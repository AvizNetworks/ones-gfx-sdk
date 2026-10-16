"""
POST /fabrics/{fabricName}/shared-services  ->  /api/fm/fabrics/{fabricName}/shared-services

Create the fabric's shared-services tenant (GFX-5709).

A fabric has at most one. It is configured on every DPU the FM manages and on
the border/storage leaves, with its own VLAN and, where IPAM/DHCP is enabled,
its own Kea subnet. Hosts join afterwards through `assign_vf_interface`.

Returns 201 with the tenant and the outcome on each DPU; the body is not
wrapped in the status/message envelope.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def create_shared_services(client: Client, fabricName: str, *, tenantName: str, vlanId: int, description: str | None = None, anycastGateway: dict[str, Any] | None = None, idempotency_key: str | None = None) -> dict[str, Any]:
    """Create the shared-services tenant on `fabricName`.

    `tenantName` starts with a letter, at most 12 characters of letters,
    digits, `-` and `_`; `Storage` is reserved. `vlanId` must be a JSON
    integer in 2-4094, excluding 500-1020, 3725-3999, 69, 100 and 110, and
    unused by another tenant. `description` is at most 60 characters.

    `anycastGateway` is `{"enabled": bool}` plus an optional `"macAddress"`;
    FM rejects any other key, and `enabled` must be a real JSON boolean.
    """
    path = f"fabrics/{fabricName}/shared-services"
    headers: dict[str, str] = {}
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    body: dict[str, Any] = {}
    body["tenantName"] = tenantName
    body["vlanId"] = vlanId
    if description is not None:
        body["description"] = description
    if anycastGateway is not None:
        body["anycastGateway"] = anycastGateway
    return client.call_api("POST", path, json_body=body, extra_headers=headers or None)
