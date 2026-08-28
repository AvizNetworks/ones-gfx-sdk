"""
GET /getIntentValidation  ->  /api/fm/getIntentValidation

Validation results for an intent.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_intent_validation(client: Client, *, intentName: str) -> list[object]:
    path = "getIntentValidation"
    params: dict[str, Any] = {}
    params["intentName"] = intentName
    return client.call_api("GET", path, params=params or None)
