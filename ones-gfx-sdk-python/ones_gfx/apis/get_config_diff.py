"""
POST /getConfigDiff  ->  /api/fm/getConfigDiff

Running-vs-saved diff for ip.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_config_diff(*, ip: str) -> object:
    client = get_client()
    path = "getConfigDiff"
    body: dict[str, Any] = {}
    body["ip"] = ip
    return client.call_api("POST", path, json_body=body or None)
