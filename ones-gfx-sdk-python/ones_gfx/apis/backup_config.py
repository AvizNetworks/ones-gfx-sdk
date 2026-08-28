"""
POST /backupConfig  ->  /api/fm/backupConfig

Backup with optional label per device.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class DeviceConfigBackup(TypedDict):
    """Helper/DeviceConfigBackup.java"""

    ip: str
    label: str


def backup_config(*, items: list[DeviceConfigBackup]) -> bool:
    client = get_client()
    path = "backupConfig"
    return client.call_api("POST", path, json_body=items)
