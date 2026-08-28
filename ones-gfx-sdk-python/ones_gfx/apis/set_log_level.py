"""
POST /log/level  ->  /api/fm/log/level

Change log levels at runtime.
"""
from __future__ import annotations


from ..client import get_client


def set_log_level(*, loggers: dict[str, str]) -> str:
    client = get_client()
    path = "log/level"
    return client.call_api("POST", path, json_body=loggers)
