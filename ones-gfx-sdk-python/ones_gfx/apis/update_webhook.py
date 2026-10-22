"""
PATCH /webhooks/{webhookId}  ->  /api/fm/webhooks/{webhookId}

Partially update a global webhook registration: only the fields you send
change.

The semantics differ from the fabric and tenant scopes, which replace the
whole body with PUT. Here an omitted `eventTypes` leaves the subscription
untouched, while an explicitly empty list clears it to "every event type" —
and the same distinction applies to each filter list.
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._webhook_types import WebhookRegistration


def update_webhook(
    client: Client,
    webhookId: str,
    *,
    callbackUrl: str | None = None,
    enabled: bool | None = None,
    eventTypes: list[str] | None = None,
    fabricNames: list[str] | None = None,
    tenantNames: list[str] | None = None,
) -> WebhookRegistration:
    """Update the global registration `webhookId` in place.

    Every argument is optional; those left out keep their stored value. Pass
    an empty list to `eventTypes`, `fabricNames` or `tenantNames` to clear
    that filter rather than to leave it alone.
    """
    path = f"webhooks/{webhookId}"
    body: dict[str, Any] = {}
    if callbackUrl is not None:
        body["callbackUrl"] = callbackUrl
    if enabled is not None:
        body["enabled"] = enabled
    if eventTypes is not None:
        body["eventTypes"] = eventTypes
    filters: dict[str, Any] = {}
    if fabricNames is not None:
        filters["fabricNames"] = fabricNames
    if tenantNames is not None:
        filters["tenantNames"] = tenantNames
    if filters:
        body["filters"] = filters
    return client.call_api("PATCH", path, json_body=body)
