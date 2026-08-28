"""
POST /restoreconfig  ->  /api/fm/restoreconfig

Restore by ip + timestamp.
"""
from __future__ import annotations


from .._types import DeviceConfigRestore
from ..client import get_client


def restore_config(*, items: list[DeviceConfigRestore]) -> bool:
    client = get_client()
    path = "restoreconfig"
    return client.call_api("POST", path, json_body=items)
