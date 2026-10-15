"""
GET /fabrics/{name}  (ONES UI config service)

The fabric's generated YAML, as text.
"""
from __future__ import annotations

from urllib.parse import quote

from ..client import Client


def get_fabric_yaml(client: Client, name: str) -> str:
    return client.call_config_api("GET", f"fabrics/{quote(name, safe='')}", accept="*/*")
