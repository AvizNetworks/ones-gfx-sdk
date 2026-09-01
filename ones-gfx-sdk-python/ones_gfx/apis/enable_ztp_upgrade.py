"""
POST /enableZTPUpgrade  ->  /api/fm/enableZTPUpgrade

Enable ZTP and run.
"""
from __future__ import annotations


from ..client import Client


def enable_ztp_upgrade(client: Client, *, payload: list[str]) -> bool:
    path = "enableZTPUpgrade"
    return client.call_api("POST", path, json_body=payload)
