"""
POST /fillrmaconfig  ->  /api/fm/fillrmaconfig

Stage RMA replacement records.
"""
from __future__ import annotations


from .._types import RMAInfoItem
from ..client import Client


def fill_rma_config(client: Client, *, items: list[RMAInfoItem]) -> bool:
    path = "fillrmaconfig"
    return client.call_api("POST", path, json_body=items)
