"""
POST /backupConfig  ->  /api/fm/backupConfig

Backup with optional label per device.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class DeviceConfigBackup(TypedDict):
    """Helper/DeviceConfigBackup.java"""

    ip: str
    label: str


def backup_config(client: Client, *, items: list[DeviceConfigBackup]) -> bool:
    path = "backupConfig"
    return client.call_api("POST", path, json_body=items)
