"""
POST /fabrics/{fabricName}/tenants/{tenantName}/gpus  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/gpus

Per-port GPU assignment.
"""
from __future__ import annotations

from typing import Any, TypedDict

from .._types import GpuAction
from ..client import Client

class GpuPortAssignmentResult(TypedDict, total=False):
    """POST .../tenants/{tenantName}/gpus"""

    fabricName: str
    tenantName: str
    servers: list[str]
    success: bool
    error: str
    message: str
    serversProcessed: int
    operation: str
    gpuIdsProcessed: int | None


def assign_gpu_ports(client: Client, fabricName: str, tenantName: str, *, operation: GpuAction, serverNames: list[str] | None = None, gpuIds: list[int] | None = None, membership: str | None = None) -> GpuPortAssignmentResult:
    path = f"fabrics/{fabricName}/tenants/{tenantName}/gpus"
    body: dict[str, Any] = {}
    body["operation"] = operation
    if serverNames is not None:
        body["serverNames"] = serverNames
    if gpuIds is not None:
        body["gpuIds"] = gpuIds
    if membership is not None:
        body["membership"] = membership
    return client.call_api("POST", path, json_body=body or None)
