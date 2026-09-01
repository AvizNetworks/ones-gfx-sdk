"""
POST /removeDeviceFacts  ->  /api/fm/removeDeviceFacts

Remove devices.
"""
from __future__ import annotations


from ..client import Client


def remove_device_facts(client: Client, *, payload: list[str]) -> bool:
    path = "removeDeviceFacts"
    return client.call_api("POST", path, json_body=payload)
