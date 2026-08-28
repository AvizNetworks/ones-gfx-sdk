"""
POST /autoAllocateGpusToTenants  ->  /api/fm/autoAllocateGpusToTenants

Auto-pick GPUs for a tenant.
"""
from __future__ import annotations

from typing import Any

from .._types import SuidMap
from ..client import get_client


def auto_allocate_gpus_to_tenants(*, fabricName: str, tenantName: str, autoAllocationDevicesNeed: int, suid: SuidMap | None = None) -> bool:
    client = get_client()
    path = "autoAllocateGpusToTenants"
    body: dict[str, Any] = {}
    body["fabricName"] = fabricName
    body["tenantName"] = tenantName
    body["autoAllocationDevicesNeed"] = autoAllocationDevicesNeed
    if suid is not None:
        body["suid"] = suid
    return client.call_api("POST", path, json_body=body or None)
