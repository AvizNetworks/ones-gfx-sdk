"""
GET /fabrics/{fabricName}/webhook  ->  /api/fm/fabrics/{fabricName}/webhook

The fabric's webhook registration. FM answers 404 when none is registered.
"""
from __future__ import annotations

from ..client import Client
from ._webhook_types import WebhookRegistration


def get_fabric_webhook(client: Client, fabricName: str) -> WebhookRegistration:
    """Fetch the webhook registered for `fabricName`."""
    path = f"fabrics/{fabricName}/webhook"
    return client.call_api("GET", path)
