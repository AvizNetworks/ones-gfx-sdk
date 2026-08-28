"""
POST /fetchdevicebackupfiles  ->  /api/fm/fetchdevicebackupfiles

Backup file contents for one ip + timestamp.
"""
from __future__ import annotations

from typing import Any, TypedDict

from ..client import get_client

class FetchBackupFilesResult(TypedDict, total=False):
    """POST /fetchdevicebackupfiles"""

    fetchedConfig: str


def fetch_device_backup_files(*, ip: str | None = None, timestamp: str | None = None) -> FetchBackupFilesResult:
    client = get_client()
    path = "fetchdevicebackupfiles"
    body: dict[str, Any] = {}
    if ip is not None:
        body["ip"] = ip
    if timestamp is not None:
        body["timestamp"] = timestamp
    return client.call_api("POST", path, json_body=body or None)
