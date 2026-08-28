"""
DELETE /delFabricSim/{name}  ->  /api/fm/delFabricSim/{name}

Delete simulation.
"""
from __future__ import annotations


from ..client import get_client


def delete_fabric_sim(name: str) -> str:
    client = get_client()
    path = f"delFabricSim/{name}"
    return client.call_api("DELETE", path)
