"""
POST /fabrics/{fabricName}/tenants  ->  /api/fm/fabrics/{fabricName}/tenants

Create tenant.
"""
from __future__ import annotations

from typing import Any

from .._types import OperationAccepted
from ..client import Client


def create_tenant(client: Client, fabricName: str, *, tenantName: str, prefer: str | None = None, idempotency_key: str | None = None, description: str | None = None, maxGpusAllowed: int | None = None, shared: bool | None = None, anycastGateway: dict[str, Any] | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> str | OperationAccepted:
    """Create a tenant on `fabricName`.

    `anycastGateway` (GFX-5611, HBN DPU fabrics in tenant DHCP mode) is
    `{"enabled": bool}` plus an optional `"macAddress"`; FM rejects any other
    key and requires `enabled` to be a real JSON boolean. Omit `macAddress`
    to take the fabric's own MAC.
    """
    path = f"fabrics/{fabricName}/tenants"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    body: dict[str, Any] = {}
    body["tenantName"] = tenantName
    if description is not None:
        body["description"] = description
    if maxGpusAllowed is not None:
        body["maxGpusAllowed"] = maxGpusAllowed
    if shared is not None:
        body["shared"] = shared
    if anycastGateway is not None:
        body["anycastGateway"] = anycastGateway
    if enableWebhook is not None:
        body["enableWebhook"] = enableWebhook
    if webhookUrl is not None:
        body["webhookUrl"] = webhookUrl
    if webhookEvents is not None:
        body["webhookEvents"] = webhookEvents
    return client.call_api("POST", path, json_body=body or None, extra_headers=headers or None)
