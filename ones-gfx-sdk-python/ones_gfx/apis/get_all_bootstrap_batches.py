"""
GET /getAllBootstrapBatches  ->  /api/fm/getAllBootstrapBatches

Batch summaries for the landing page.
"""
from __future__ import annotations


from ..client import get_client


def get_all_bootstrap_batches() -> list[dict[str, object]]:
    client = get_client()
    path = "getAllBootstrapBatches"
    return client.call_api("GET", path)
