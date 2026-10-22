"""
GET /webhooks/{webhookId}  ->  /api/fm/webhooks/{webhookId}

One global webhook registration. FM answers 404 when `webhookId` is unknown.
"""
from __future__ import annotations

from ..client import Client
from ._webhook_types import WebhookRegistration


def get_webhook(client: Client, webhookId: str) -> WebhookRegistration:
    """Fetch the global registration `webhookId`."""
    path = f"webhooks/{webhookId}"
    return client.call_api("GET", path)
