"""
GET /fabrics/{fabricName}/tenants/{tenantName}/webhook  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/webhook

The tenant's webhook registration. FM answers 404 when none is registered.
"""
from __future__ import annotations

from ..client import Client
from ._webhook_types import WebhookRegistration


def get_tenant_webhook(client: Client, fabricName: str, tenantName: str) -> WebhookRegistration:
    """Fetch the webhook registered for `tenantName` on `fabricName`."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/webhook"
    return client.call_api("GET", path)
