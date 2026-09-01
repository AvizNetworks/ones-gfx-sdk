"""
GET /getGpusByHost/{hostName}/{fabricName}  ->  /api/fm/getGpusByHost/{hostName}/{fabricName}

GPUs on one host.
"""
from __future__ import annotations


from .._types import GpuItem
from ..client import Client


def get_gpus_by_host(client: Client, hostName: str, fabricName: str) -> list[GpuItem]:
    path = f"getGpusByHost/{hostName}/{fabricName}"
    return client.call_api("GET", path)
