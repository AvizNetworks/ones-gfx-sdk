"""
POST /fabrics/{fabricName}/tenants  ->  /api/fm/fabrics/{fabricName}/tenants

Create tenant.
"""
from __future__ import annotations

from typing import Any

from ..client import ONESClient


def create_tenant(client: ONESClient, fabricName: str, *, body, prefer=None, idempotency_key=None) -> Any:
    path = f"fabrics/{fabricName}/tenants"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    return client.call_api("POST", path, json_body=body, extra_headers=headers or None)
