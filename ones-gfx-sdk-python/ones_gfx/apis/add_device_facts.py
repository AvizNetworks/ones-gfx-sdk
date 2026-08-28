"""
POST /addDeviceFacts  ->  /api/fm/addDeviceFacts

Add devices.
"""
from __future__ import annotations


from .._types import DeviceDetail
from ..client import get_client


def add_device_facts(*, items: list[DeviceDetail]) -> bool:
    client = get_client()
    path = "addDeviceFacts"
    return client.call_api("POST", path, json_body=items)
