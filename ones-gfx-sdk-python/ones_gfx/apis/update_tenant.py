"""
PATCH /fabrics/{fabricName}/tenants/{tenantName}  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}

Add/remove GPU servers on a tenant.
"""
from __future__ import annotations

from typing import Any

from .._types import ApiResponseMessage, GpuAction, GpuServerInfo, OperationAccepted
from ..client import get_client


def update_tenant(fabricName: str, tenantName: str, *, servers: list[GpuServerInfo], prefer: str | None = None, idempotency_key: str | None = None, operation: GpuAction | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> ApiResponseMessage | OperationAccepted:
    client = get_client()
    path = f"fabrics/{fabricName}/tenants/{tenantName}"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    body: dict[str, Any] = {}
    if operation is not None:
        body["operation"] = operation
    body["servers"] = servers
    if enableWebhook is not None:
        body["enableWebhook"] = enableWebhook
    if webhookUrl is not None:
        body["webhookUrl"] = webhookUrl
    if webhookEvents is not None:
        body["webhookEvents"] = webhookEvents
    return client.call_api("PATCH", path, json_body=body or None, extra_headers=headers or None)
