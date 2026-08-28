"""
DELETE /fabrics/{fabricName}/vpcpeering  ->  /api/fm/fabrics/{fabricName}/vpcpeering

Remove route-leak entries.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def delete_vpc_peering(fabricName: str, *, name: str | None = None, vpcname: str | None = None, peervpcname: str | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> str:
    client = get_client()
    path = f"fabrics/{fabricName}/vpcpeering"
    body: dict[str, Any] = {}
    if name is not None:
        body["name"] = name
    if vpcname is not None:
        body["vpcname"] = vpcname
    if peervpcname is not None:
        body["peervpcname"] = peervpcname
    if enableWebhook is not None:
        body["enableWebhook"] = enableWebhook
    if webhookUrl is not None:
        body["webhookUrl"] = webhookUrl
    if webhookEvents is not None:
        body["webhookEvents"] = webhookEvents
    return client.call_api("DELETE", path, json_body=body or None)
