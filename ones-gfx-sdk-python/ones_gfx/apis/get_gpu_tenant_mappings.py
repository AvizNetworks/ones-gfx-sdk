"""
GET /fabrics/{fabricName}/gpuTenantMappings  ->  /api/fm/fabrics/{fabricName}/gpuTenantMappings

Per-GPU tenant ownership.
"""
from __future__ import annotations

from typing import Any

from ..client import ONESClient


def get_gpu_tenant_mappings(client: ONESClient, fabricName: str, *, tenantName=None) -> Any:
    path = f"fabrics/{fabricName}/gpuTenantMappings"
    params: dict[str, Any] = {}
    if tenantName is not None:
        params["tenantName"] = tenantName
    return client.call_api("GET", path, params=params or None)
