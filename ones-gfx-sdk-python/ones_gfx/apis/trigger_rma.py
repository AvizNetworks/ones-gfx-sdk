"""
POST /triggerrma  ->  /api/fm/triggerrma

Execute the RMA swap.
"""
from __future__ import annotations


from .._types import RMAInfoItem
from ..client import get_client


def trigger_rma(*, items: list[RMAInfoItem]) -> bool:
    client = get_client()
    path = "triggerrma"
    return client.call_api("POST", path, json_body=items)
