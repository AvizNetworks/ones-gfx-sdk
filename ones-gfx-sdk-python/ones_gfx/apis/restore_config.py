"""
POST /restoreconfig  ->  /api/fm/restoreconfig

Restore by ip + timestamp.
"""
from __future__ import annotations


from .._types import DeviceConfigRestore
from ..client import Client


def restore_config(client: Client, *, items: list[DeviceConfigRestore]) -> bool:
    path = "restoreconfig"
    return client.call_api("POST", path, json_body=items)
