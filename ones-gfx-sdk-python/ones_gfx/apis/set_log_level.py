"""
POST /log/level  ->  /api/fm/log/level

Change log levels at runtime.
"""
from __future__ import annotations


from ..client import Client


def set_log_level(client: Client, *, loggers: dict[str, str]) -> str:
    path = "log/level"
    return client.call_api("POST", path, json_body=loggers)
