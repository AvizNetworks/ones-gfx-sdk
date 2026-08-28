"""
POST /rebootRequest  ->  /api/fm/rebootRequest

Reboot devices.
"""
from __future__ import annotations


from ..client import get_client


def reboot_request(*, payload: list[str]) -> bool:
    client = get_client()
    path = "rebootRequest"
    return client.call_api("POST", path, json_body=payload)
