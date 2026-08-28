"""
POST /fabrics/{fabricName}/nmxc/domains/probe  ->  /api/fm/fabrics/{fabricName}/nmxc/domains/probe

Read-only NMX-C domain probe.
"""
from __future__ import annotations

from typing import Any, TypedDict

from .._types import NmxcDomain
from ..client import get_client

class NmxcProbeResult(TypedDict):
    """POST .../nmxc/domains/probe — one entry per requested domain"""

    host: str
    port: int
    reachable: bool
    configured: bool
    error: str | None


def probe_nmxc_domains(fabricName: str, *, domains: list[NmxcDomain]) -> list[NmxcProbeResult]:
    client = get_client()
    path = f"fabrics/{fabricName}/nmxc/domains/probe"
    body: dict[str, Any] = {}
    body["domains"] = domains
    return client.call_api("POST", path, json_body=body or None)
