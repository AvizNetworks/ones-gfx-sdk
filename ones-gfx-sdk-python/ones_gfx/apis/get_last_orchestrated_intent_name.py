"""
GET /getLastOrchestratedIntentName  ->  /api/fm/getLastOrchestratedIntentName

Name of the most recently orchestrated intent.
"""
from __future__ import annotations


from ..client import get_client


def get_last_orchestrated_intent_name() -> str:
    client = get_client()
    path = "getLastOrchestratedIntentName"
    return client.call_api("GET", path)
