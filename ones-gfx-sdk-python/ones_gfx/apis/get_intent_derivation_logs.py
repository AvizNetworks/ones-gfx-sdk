"""
GET /getIntentDerivationLogs  ->  /api/fm/getIntentDerivationLogs

Intent derivation logs for one device.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_intent_derivation_logs(*, device: str) -> list[object]:
    client = get_client()
    path = "getIntentDerivationLogs"
    params: dict[str, Any] = {}
    params["device"] = device
    return client.call_api("GET", path, params=params or None)
