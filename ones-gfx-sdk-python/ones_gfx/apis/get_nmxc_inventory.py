"""
GET /fabrics/{fabricName}/nmxc/inventory  ->  /api/fm/fabrics/{fabricName}/nmxc/inventory

Full NMX-C inventory snapshot.
"""
from __future__ import annotations


from ..client import get_client


def get_nmxc_inventory(fabricName: str) -> dict[str, object]:
    client = get_client()
    path = f"fabrics/{fabricName}/nmxc/inventory"
    return client.call_api("GET", path)
