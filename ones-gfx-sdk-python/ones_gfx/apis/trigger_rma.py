"""
POST /triggerrma  ->  /api/fm/triggerrma

Execute the RMA swap.
"""
from __future__ import annotations


from .._types import RMAInfoItem
from ..client import Client


def trigger_rma(client: Client, *, items: list[RMAInfoItem]) -> bool:
    path = "triggerrma"
    return client.call_api("POST", path, json_body=items)
