"""
POST /netops/{ipAddress}  ->  /api/fm/netops/{ipAddress}

Execute a NetOps operation on a single device.
"""
from __future__ import annotations

from typing import Any

from .._types import NetOpsAction
from ..client import get_client


def netops_device(ipAddress: str, *, action: NetOpsAction, params: dict[str, str | int | list[str]]) -> bool:
    client = get_client()
    path = f"netops/{ipAddress}"
    body: dict[str, Any] = {}
    body["action"] = action
    body["params"] = params
    return client.call_api("POST", path, json_body=body or None)
