"""
GET /log/level  ->  /api/fm/log/level

Map of logger name -> current level.
"""
from __future__ import annotations


from ..client import Client


def get_log_level(client: Client) -> dict[str, str]:
    path = "log/level"
    return client.call_api("GET", path)
