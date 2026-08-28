"""
GET /getFabricDevices/{fabricName}  ->  /api/fm/getFabricDevices/{fabricName}

Device IPs in the fabric.
"""
from __future__ import annotations


from ..client import Client


def get_fabric_devices(client: Client, fabricName: str) -> list[str]:
    path = f"getFabricDevices/{fabricName}"
    return client.call_api("GET", path)
