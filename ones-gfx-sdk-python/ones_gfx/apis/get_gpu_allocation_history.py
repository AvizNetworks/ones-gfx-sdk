"""
GET /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocationHistory  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/gpuAllocationHistory

Allocation history records.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class GpuAllocationHistoryItem(TypedDict, total=False):
    """Models/GpuAllocationHistory.java JSON"""

    id: int
    suNumber: int
    hostName: str
    gpusAdded: str
    gpusRemoved: str
    tenantName: str
    fabricName: str
    createdAt: str
    updatedAt: str


def get_gpu_allocation_history(client: Client, fabricName: str, tenantName: str) -> list[GpuAllocationHistoryItem]:
    path = f"fabrics/{fabricName}/tenants/{tenantName}/gpuAllocationHistory"
    return client.call_api("GET", path)
