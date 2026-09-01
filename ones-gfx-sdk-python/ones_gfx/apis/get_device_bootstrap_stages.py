"""
GET /getDeviceBootstrapStages/{bootstrapId}  ->  /api/fm/getDeviceBootstrapStages/{bootstrapId}

Stage progress.
"""
from __future__ import annotations


from ..client import Client


def get_device_bootstrap_stages(client: Client, bootstrapId: str) -> list[dict[str, object]]:
    path = f"getDeviceBootstrapStages/{bootstrapId}"
    return client.call_api("GET", path)
