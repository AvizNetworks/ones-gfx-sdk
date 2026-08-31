"""
POST /replaceConfig  ->  /api/fm/replaceConfig

Push a replacement config file.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def replace_config(client: Client, *, deviceip: str, file: str, onlydiff: bool | None = None) -> object:
    path = "replaceConfig"
    params: dict[str, Any] = {}
    params["deviceip"] = deviceip
    if onlydiff is not None:
        params["onlydiff"] = onlydiff
    with open(file, "rb") as fh:
        files: dict[str, Any] = {"file": fh}
        return client.call_api("POST", path, params=params or None, files=files or None)
