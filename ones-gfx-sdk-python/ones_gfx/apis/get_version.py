"""
POST /getVersion  ->  /api/fm/getVersion

Device facts/versions.
"""
from __future__ import annotations


from ..client import Client


def get_version(client: Client, *, payload: list[str]) -> list[object]:
    path = "getVersion"
    return client.call_api("POST", path, json_body=payload)
