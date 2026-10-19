"""
GET /webhooks  ->  /api/fm/webhooks

Every global webhook registration, newest first. Fabric- and tenant-scoped
registrations are not listed here; read those at their own paths.
"""
from __future__ import annotations

from ..client import Client
from ._webhook_types import WebhookRegistration


def list_webhooks(client: Client) -> list[WebhookRegistration]:
    """List the global webhook registrations (possibly empty)."""
    path = "webhooks"
    return client.call_api("GET", path)
