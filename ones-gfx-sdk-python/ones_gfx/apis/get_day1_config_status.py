"""
GET /getDay1ConfigStatus  ->  /api/fm/getDay1ConfigStatus

Per-device Day-1 config status.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_day1_config_status(*, intentName: str) -> list[object]:
    client = get_client()
    path = "getDay1ConfigStatus"
    params: dict[str, Any] = {}
    params["intentName"] = intentName
    return client.call_api("GET", path, params=params or None)
