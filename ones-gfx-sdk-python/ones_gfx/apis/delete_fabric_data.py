"""
DELETE /delFabricData/{name}  ->  /api/fm/delFabricData/{name}

Delete fabric and deallocate its LAAS license.
"""
from __future__ import annotations


from ..client import Client


def delete_fabric_data(client: Client, name: str) -> str:
    path = f"delFabricData/{name}"
    return client.call_api("DELETE", path)
