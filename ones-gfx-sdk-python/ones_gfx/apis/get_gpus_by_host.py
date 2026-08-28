"""
GET /getGpusByHost/{hostName}/{fabricName}  ->  /api/fm/getGpusByHost/{hostName}/{fabricName}

GPUs on one host.
"""
from __future__ import annotations


from .._types import GpuItem
from ..client import get_client


def get_gpus_by_host(hostName: str, fabricName: str) -> list[GpuItem]:
    client = get_client()
    path = f"getGpusByHost/{hostName}/{fabricName}"
    return client.call_api("GET", path)
