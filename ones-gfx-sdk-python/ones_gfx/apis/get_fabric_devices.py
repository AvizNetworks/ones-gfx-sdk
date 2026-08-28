"""
GET /getFabricDevices/{fabricName}  ->  /api/fm/getFabricDevices/{fabricName}

Device IPs in the fabric.
"""
from __future__ import annotations


from ..client import get_client


def get_fabric_devices(fabricName: str) -> list[str]:
    client = get_client()
    path = f"getFabricDevices/{fabricName}"
    return client.call_api("GET", path)
