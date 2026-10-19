"""
DELETE /fabrics/{fabricName}/webhook  ->  /api/fm/fabrics/{fabricName}/webhook

Delete the fabric's webhook registration. Returns None (FM answers 204); 404
when none is registered.
"""
from __future__ import annotations

from ..client import Client


def delete_fabric_webhook(client: Client, fabricName: str) -> None:
    """Delete the webhook registered for `fabricName`."""
    path = f"fabrics/{fabricName}/webhook"
    return client.call_api("DELETE", path)
