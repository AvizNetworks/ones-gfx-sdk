"""
POST /removeDeviceFacts  ->  /api/fm/removeDeviceFacts

Remove devices.
"""
from __future__ import annotations


from ..client import get_client


def remove_device_facts(*, payload: list[str]) -> bool:
    client = get_client()
    path = "removeDeviceFacts"
    return client.call_api("POST", path, json_body=payload)
