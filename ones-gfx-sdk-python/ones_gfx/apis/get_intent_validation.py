"""
GET /getIntentValidation  ->  /api/fm/getIntentValidation

Validation results for an intent.
"""
from __future__ import annotations

from typing import Any

from ..client import get_client


def get_intent_validation(*, intentName: str) -> list[object]:
    client = get_client()
    path = "getIntentValidation"
    params: dict[str, Any] = {}
    params["intentName"] = intentName
    return client.call_api("GET", path, params=params or None)
