"""
POST /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations

Assign/unassign GPUs to a tenant.
"""
from __future__ import annotations

from typing import Any, Literal

from .._types import ApiResponseMessage, GpuAction, OperationAccepted, SuidMap
from ..client import Client

ConfigScope = Literal["WHOLE_SERVER", "PARTICULAR_GPU"]


def modify_gpu_allocations(client: Client, fabricName: str, tenantName: str, *, suid: SuidMap, prefer: str | None = None, idempotency_key: str | None = None, operation: GpuAction | None = None, configScope: ConfigScope | None = None, unreachableDevices: list[str] | None = None, enableWebhook: bool | None = None, webhookUrl: str | None = None, webhookEvents: list[str] | None = None) -> ApiResponseMessage | OperationAccepted:
    path = f"fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations"
    headers: dict[str, str] = {}
    if prefer is not None:
        headers["Prefer"] = prefer
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    body: dict[str, Any] = {}
    body["suid"] = suid
    if operation is not None:
        body["operation"] = operation
    body["tenantName"] = tenantName
    body["fabricName"] = fabricName
    if configScope is not None:
        body["configScope"] = configScope
    if unreachableDevices is not None:
        body["unreachableDevices"] = unreachableDevices
    if enableWebhook is not None:
        body["enableWebhook"] = enableWebhook
    if webhookUrl is not None:
        body["webhookUrl"] = webhookUrl
    if webhookEvents is not None:
        body["webhookEvents"] = webhookEvents
    return client.call_api("POST", path, json_body=body or None, extra_headers=headers or None)
