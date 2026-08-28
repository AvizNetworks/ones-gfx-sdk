"""
POST /fillrmaconfig  ->  /api/fm/fillrmaconfig

Stage RMA replacement records.
"""
from __future__ import annotations


from .._types import RMAInfoItem
from ..client import get_client


def fill_rma_config(*, items: list[RMAInfoItem]) -> bool:
    client = get_client()
    path = "fillrmaconfig"
    return client.call_api("POST", path, json_body=items)
