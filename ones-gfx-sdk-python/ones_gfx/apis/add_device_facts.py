"""
POST /addDeviceFacts  ->  /api/fm/addDeviceFacts

Add devices.
"""
from __future__ import annotations


from .._types import DeviceDetail
from ..client import Client


def add_device_facts(client: Client, *, items: list[DeviceDetail]) -> bool:
    path = "addDeviceFacts"
    return client.call_api("POST", path, json_body=items)
