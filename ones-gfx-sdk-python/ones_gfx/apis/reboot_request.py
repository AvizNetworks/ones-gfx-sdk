"""
POST /rebootRequest  ->  /api/fm/rebootRequest

Reboot devices.
"""
from __future__ import annotations


from ..client import Client


def reboot_request(client: Client, *, payload: list[str]) -> bool:
    path = "rebootRequest"
    return client.call_api("POST", path, json_body=payload)
