"""
GET /log/level  ->  /api/fm/log/level

Map of logger name -> current level.
"""
from __future__ import annotations


from ..client import get_client


def get_log_level() -> dict[str, str]:
    client = get_client()
    path = "log/level"
    return client.call_api("GET", path)
