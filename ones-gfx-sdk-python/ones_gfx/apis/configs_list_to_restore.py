"""
POST /configslisttorestore  ->  /api/fm/configslisttorestore

List restorable backups.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def configs_list_to_restore(*, devices: list[str] | None = None, onlylimited: bool | None = None) -> str | None:
    client = get_client()
    path = "configslisttorestore"
    body: dict[str, Any] = {}
    if devices is not None:
        body["devices"] = devices
    if onlylimited is not None:
        body["onlylimited"] = onlylimited
    return client.call_api("POST", path, json_body=body or None)
