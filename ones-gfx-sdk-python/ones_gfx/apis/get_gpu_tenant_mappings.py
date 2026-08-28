"""
GET /fabrics/{fabricName}/gpuTenantMappings  ->  /api/fm/fabrics/{fabricName}/gpuTenantMappings

Per-GPU tenant ownership.
"""
from __future__ import annotations

from typing import Any, TypedDict

from ..client import get_client

class GpuTenantMappingItem(TypedDict, total=False):
    """Models/GpuTenantMapping.java JSON"""

    id: int
    fabricName: str
    serverName: str
    gpuIndex: int
    logicalGpuName: str
    tenantName: str
    configStatus: str
    createdAt: str
    updatedAt: str


def get_gpu_tenant_mappings(fabricName: str, *, tenantName: str | None = None) -> list[GpuTenantMappingItem]:
    client = get_client()
    path = f"fabrics/{fabricName}/gpuTenantMappings"
    params: dict[str, Any] = {}
    if tenantName is not None:
        params["tenantName"] = tenantName
    return client.call_api("GET", path, params=params or None)
