"""
DELETE /fabrics/{fabricName}/tenants/{tenantName}/webhook  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/webhook

Delete the tenant's webhook registration. Returns None (FM answers 204); 404
when none is registered.
"""
from __future__ import annotations

from ..client import Client


def delete_tenant_webhook(client: Client, fabricName: str, tenantName: str) -> None:
    """Delete the webhook registered for `tenantName` on `fabricName`."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/webhook"
    return client.call_api("DELETE", path)
