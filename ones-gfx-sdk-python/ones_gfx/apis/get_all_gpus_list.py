"""
GET /getAllGpusList/{fabricName}  ->  /api/fm/getAllGpusList/{fabricName}

All GPUs in a fabric.
"""
from __future__ import annotations


from .._types import GpuItem
from ..client import Client


def get_all_gpus_list(client: Client, fabricName: str) -> list[GpuItem]:
    path = f"getAllGpusList/{fabricName}"
    return client.call_api("GET", path)
