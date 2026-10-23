"""
POST /fabrics/{fabricName}/tenants  ->  /api/fm/fabrics/{fabricName}/tenants

Create tenant.
"""
from __future__ import annotations

from typing import Any

from .._types import OperationAccepted
from ..client import Client
from ._validation import validate_pool_cidr


def create_tenant(client: Client, fabricName: str, *, tenantName: str, prefer: str | None = None, idempotency_key: str | None = None, description: str | None = None, maxGpusAllowed: int | None = None, shared: bool | None = None, tenantDhcpCidrPool: str | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> str | OperationAccepted:
    """Create a tenant.

    ``tenantDhcpCidrPool`` (e.g. ``"192.3.0.0/17"``) is the tenant's IPAM pool.
    It is required when IPAM is enabled on the fabric and the fabric has no
    tenant compute pool to derive one from (server error ``IPAM_CIDR_REQUIRED``).
    It is validated before the request: IPv4, prefix /8–/30, otherwise ``ValueError``.
    """
    if tenantDhcpCidrPool is not None:
        tenantDhcpCidrPool = validate_pool_cidr(tenantDhcpCidrPool, "tenantDhcpCidrPool")
    path = f"fabrics/{fabricName}/tenants"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    body: dict[str, Any] = {}
    body["tenantName"] = tenantName
    if description is not None:
        body["description"] = description
    if maxGpusAllowed is not None:
        body["maxGpusAllowed"] = maxGpusAllowed
    if shared is not None:
        body["shared"] = shared
    if tenantDhcpCidrPool is not None:
        body["tenantDhcpCidrPool"] = tenantDhcpCidrPool
    if enableWebhook is not None:
        body["enableWebhook"] = enableWebhook
    if webhookUrl is not None:
        body["webhookUrl"] = webhookUrl
    if webhookEvents is not None:
        body["webhookEvents"] = webhookEvents
    return client.call_api("POST", path, json_body=body or None, extra_headers=headers or None)
