"""
POST /replaceConfig  ->  /api/fm/replaceConfig

Push a replacement config file.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def replace_config(*, deviceip: str, file: str, onlydiff: bool | None = None) -> object:
    client = get_client()
    path = "replaceConfig"
    params: dict[str, Any] = {}
    params["deviceip"] = deviceip
    if onlydiff is not None:
        params["onlydiff"] = onlydiff
    files: dict[str, Any] = {}
    files["file"] = open(file, "rb")
    return client.call_api("POST", path, params=params or None, files=files or None)
