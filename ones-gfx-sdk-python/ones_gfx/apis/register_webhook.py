"""
POST /webhooks  ->  /api/fm/webhooks

Register a global webhook, covering every fabric and tenant unless narrowed
by filters.

Idempotent on `callbackUrl`: re-registering the same URL updates that
subscription instead of adding a duplicate. FM answers 201 when it created a
new one and 200 when it returned or updated an existing one, so a configurator
that re-registers on every restart neither fails nor accumulates rows.

Unlike the fabric and tenant scopes, `eventTypes` is required and must be
non-empty here — a global subscription to every event is almost never
intended, so FM rejects it rather than delivering it.
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._webhook_types import WebhookRegistration


def register_webhook(
    client: Client,
    *,
    callbackUrl: str,
    eventTypes: list[str],
    enabled: bool | None = None,
    fabricNames: list[str] | None = None,
    tenantNames: list[str] | None = None,
) -> WebhookRegistration:
    """Register (or idempotently update) the global webhook for `callbackUrl`.

    `eventTypes` is required and non-empty; see `WEBHOOK_EVENT_TYPES` for the
    canonical names and `WEBHOOK_EVENT_ALIASES` for the accepted `fabric.*`
    spellings. `fabricNames` and `tenantNames` narrow delivery; omit both to
    receive matching events from everywhere.
    """
    path = "webhooks"
    body: dict[str, Any] = {"callbackUrl": callbackUrl, "eventTypes": eventTypes}
    if enabled is not None:
        body["enabled"] = enabled
    filters: dict[str, Any] = {}
    if fabricNames is not None:
        filters["fabricNames"] = fabricNames
    if tenantNames is not None:
        filters["tenantNames"] = tenantNames
    if filters:
        body["filters"] = filters
    return client.call_api("POST", path, json_body=body)
