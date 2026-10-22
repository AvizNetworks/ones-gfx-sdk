"""
POST /fabrics/{fabricName}/tenants/{tenantName}/webhook  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/webhook

Register (or replace) the webhook for one tenant. FM answers 201.

A tenant holds at most one registration, so posting again replaces it. Omit
`eventTypes` to receive every event for the tenant.
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._webhook_types import WebhookRegistration


def register_tenant_webhook(
    client: Client,
    fabricName: str,
    tenantName: str,
    *,
    callbackUrl: str,
    enabled: bool | None = None,
    eventTypes: list[str] | None = None,
) -> WebhookRegistration:
    """Register the webhook for `tenantName` on `fabricName`."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/webhook"
    body: dict[str, Any] = {"callbackUrl": callbackUrl}
    if enabled is not None:
        body["enabled"] = enabled
    if eventTypes is not None:
        body["eventTypes"] = eventTypes
    return client.call_api("POST", path, json_body=body)
