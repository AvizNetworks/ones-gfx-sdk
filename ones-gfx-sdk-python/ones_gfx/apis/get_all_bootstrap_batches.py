"""
GET /getAllBootstrapBatches  ->  /api/fm/getAllBootstrapBatches

Batch summaries for the landing page.
"""
from __future__ import annotations


from ..client import Client


def get_all_bootstrap_batches(client: Client) -> list[dict[str, object]]:
    path = "getAllBootstrapBatches"
    return client.call_api("GET", path)
