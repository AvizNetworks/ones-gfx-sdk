"""
GET /fabrics/{fabricName}/available_servers  ->  /api/fm/fabrics/{fabricName}/available_servers

Servers with free GPUs.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class AvailableServerResult(TypedDict, total=False):
    """Cumulus/dto/AvailableServer.java"""

    availableGPUs: list[str]


def get_available_servers(client: Client, fabricName: str) -> AvailableServerResult:
    path = f"fabrics/{fabricName}/available_servers"
    return client.call_api("GET", path)
