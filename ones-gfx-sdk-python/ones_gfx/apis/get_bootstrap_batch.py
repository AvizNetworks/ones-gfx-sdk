"""
GET /getBootstrapBatch/{batchName}  ->  /api/fm/getBootstrapBatch/{batchName}

Batch detail with its devices.
"""
from __future__ import annotations


from ..client import get_client


def get_bootstrap_batch(batchName: str) -> dict[str, object]:
    client = get_client()
    path = f"getBootstrapBatch/{batchName}"
    return client.call_api("GET", path)
