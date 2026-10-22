"""
DELETE /webhooks/{webhookId}  ->  /api/fm/webhooks/{webhookId}

Delete a global webhook registration. Returns None (FM answers 204); 404 when
`webhookId` is unknown.
"""
from __future__ import annotations

from ..client import Client


def delete_webhook(client: Client, webhookId: str) -> None:
    """Delete the global registration `webhookId`."""
    path = f"webhooks/{webhookId}"
    return client.call_api("DELETE", path)
