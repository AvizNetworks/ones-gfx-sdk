"""
POST /enableZTPUpgrade  ->  /api/fm/enableZTPUpgrade

Enable ZTP and run.
"""
from __future__ import annotations


from ..client import get_client


def enable_ztp_upgrade(*, payload: list[str]) -> bool:
    client = get_client()
    path = "enableZTPUpgrade"
    return client.call_api("POST", path, json_body=payload)
