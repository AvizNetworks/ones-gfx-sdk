"""
POST /netopsfabric/{fabricName}  ->  /api/fm/netopsfabric/{fabricName}

Execute a NetOps operation across an entire fabric.
"""
from __future__ import annotations

from typing import Any

from .._types import NetOpsAction
from ..client import get_client


def netops_fabric(fabricName: str, *, action: NetOpsAction, params: dict[str, str | int | list[str]]) -> bool:
    client = get_client()
    path = f"netopsfabric/{fabricName}"
    body: dict[str, Any] = {}
    body["action"] = action
    body["params"] = params
    return client.call_api("POST", path, json_body=body or None)
