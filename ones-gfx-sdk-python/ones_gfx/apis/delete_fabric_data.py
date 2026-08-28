"""
DELETE /delFabricData/{name}  ->  /api/fm/delFabricData/{name}

Delete fabric and deallocate its LAAS license.
"""
from __future__ import annotations


from ..client import get_client


def delete_fabric_data(name: str) -> str:
    client = get_client()
    path = f"delFabricData/{name}"
    return client.call_api("DELETE", path)
