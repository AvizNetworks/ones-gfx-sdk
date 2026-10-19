"""
PUT /fabrics/{fabricName}/tenants/{tenantName}/webhook  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/webhook

Update the tenant's webhook registration.

`enabled` is sent on every call because FM reads it as a plain boolean here
rather than as "leave unchanged when absent"; it defaults to True, matching
the server's own default. Passing an empty `eventTypes` list clears the
filter, subscribing the registration to every event.
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._webhook_types import WebhookRegistration


def update_tenant_webhook(
    client: Client,
    fabricName: str,
    tenantName: str,
    *,
    callbackUrl: str | None = None,
    enabled: bool = True,
    eventTypes: list[str] | None = None,
) -> WebhookRegistration:
    """Update the webhook registered for `tenantName`. 404 when none is."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/webhook"
    body: dict[str, Any] = {"enabled": enabled}
    if callbackUrl is not None:
        body["callbackUrl"] = callbackUrl
    if eventTypes is not None:
        body["eventTypes"] = eventTypes
    return client.call_api("PUT", path, json_body=body)
