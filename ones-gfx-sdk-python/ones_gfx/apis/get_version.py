"""
POST /getVersion  ->  /api/fm/getVersion

Device facts/versions.
"""
from __future__ import annotations


from ..client import get_client


def get_version(*, payload: list[str]) -> list[object]:
    client = get_client()
    path = "getVersion"
    return client.call_api("POST", path, json_body=payload)
