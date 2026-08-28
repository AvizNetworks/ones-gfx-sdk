"""
POST /triggerbootstrapconfig  ->  /api/fm/triggerbootstrapconfig

Start the bootstrap run for a batch.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class TriggerBootstrapResult(TypedDict):
    """POST /triggerbootstrapconfig"""

    success: bool
    message: str


def trigger_bootstrap_config(client: Client) -> TriggerBootstrapResult:
    path = "triggerbootstrapconfig"
    return client.call_api("POST", path)
