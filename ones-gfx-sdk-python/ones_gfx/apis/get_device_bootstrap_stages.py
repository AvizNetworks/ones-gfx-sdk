"""
GET /getDeviceBootstrapStages/{bootstrapId}  ->  /api/fm/getDeviceBootstrapStages/{bootstrapId}

Stage progress.
"""
from __future__ import annotations


from ..client import get_client


def get_device_bootstrap_stages(bootstrapId: str) -> list[dict[str, object]]:
    client = get_client()
    path = f"getDeviceBootstrapStages/{bootstrapId}"
    return client.call_api("GET", path)
