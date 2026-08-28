"""
GET /getAllGpusList/{fabricName}  ->  /api/fm/getAllGpusList/{fabricName}

All GPUs in a fabric.
"""
from __future__ import annotations


from .._types import GpuItem
from ..client import get_client


def get_all_gpus_list(fabricName: str) -> list[GpuItem]:
    client = get_client()
    path = f"getAllGpusList/{fabricName}"
    return client.call_api("GET", path)
