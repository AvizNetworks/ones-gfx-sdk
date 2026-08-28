"""
DELETE /delFabricSim/{name}  ->  /api/fm/delFabricSim/{name}

Delete simulation.
"""
from __future__ import annotations


from ..client import Client


def delete_fabric_sim(client: Client, name: str) -> str:
    path = f"delFabricSim/{name}"
    return client.call_api("DELETE", path)
