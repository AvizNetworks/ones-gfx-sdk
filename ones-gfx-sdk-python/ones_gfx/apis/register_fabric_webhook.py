"""
POST /fabrics/{fabricName}/webhook  ->  /api/fm/fabrics/{fabricName}/webhook

Register (or replace) the webhook for one fabric. FM answers 201.

A fabric holds at most one registration, so posting again replaces it. Omit
`eventTypes` to receive every event for the fabric.
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._webhook_types import WebhookRegistration


def register_fabric_webhook(
    client: Client,
    fabricName: str,
    *,
    callbackUrl: str,
    enabled: bool | None = None,
    eventTypes: list[str] | None = None,
) -> WebhookRegistration:
    """Register the webhook for `fabricName`.

    `callbackUrl` must be an HTTP/HTTPS URL FM accepts; see
    `WEBHOOK_EVENT_TYPES` for the event names.
    """
    path = f"fabrics/{fabricName}/webhook"
    body: dict[str, Any] = {"callbackUrl": callbackUrl}
    if enabled is not None:
        body["enabled"] = enabled
    if eventTypes is not None:
        body["eventTypes"] = eventTypes
    return client.call_api("POST", path, json_body=body)
