"""
GET /getIntentDerivationLogs  ->  /api/fm/getIntentDerivationLogs

Intent derivation logs for one device.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_intent_derivation_logs(client: Client, *, device: str) -> list[object]:
    path = "getIntentDerivationLogs"
    params: dict[str, Any] = {}
    params["device"] = device
    return client.call_api("GET", path, params=params or None)
