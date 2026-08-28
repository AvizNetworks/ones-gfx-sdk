"""
GET /fabrics/{fabricName}/nmxc/inventory  ->  /api/fm/fabrics/{fabricName}/nmxc/inventory

Full NMX-C inventory snapshot.
"""
from __future__ import annotations


from ..client import Client


def get_nmxc_inventory(client: Client, fabricName: str) -> dict[str, object]:
    path = f"fabrics/{fabricName}/nmxc/inventory"
    return client.call_api("GET", path)
