"""
GET /getLastOrchestratedIntentName  ->  /api/fm/getLastOrchestratedIntentName

Name of the most recently orchestrated intent.
"""
from __future__ import annotations


from ..client import Client


def get_last_orchestrated_intent_name(client: Client) -> str:
    path = "getLastOrchestratedIntentName"
    return client.call_api("GET", path)
