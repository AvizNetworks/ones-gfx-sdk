"""
GET /getBootstrapBatch/{batchName}  ->  /api/fm/getBootstrapBatch/{batchName}

Batch detail with its devices.
"""
from __future__ import annotations


from ..client import Client


def get_bootstrap_batch(client: Client, batchName: str) -> dict[str, object]:
    path = f"getBootstrapBatch/{batchName}"
    return client.call_api("GET", path)
