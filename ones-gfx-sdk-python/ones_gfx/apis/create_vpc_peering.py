"""
POST /fabrics/{fabricName}/vpcpeering  ->  /api/fm/fabrics/{fabricName}/vpcpeering

Create route leaking between two VPCs.
"""
from __future__ import annotations

from typing import Any

from .._types import OperationAccepted
from ..client import get_client


def create_vpc_peering(fabricName: str, *, prefer: str | None = None, idempotency_key: str | None = None, name: str | None = None, vpcname: str | None = None, peervpcname: str | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> str | OperationAccepted:
    client = get_client()
    path = f"fabrics/{fabricName}/vpcpeering"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
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
    return client.call_api("POST", path, json_body=body or None, extra_headers=headers or None)
