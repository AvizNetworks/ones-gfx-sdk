"""
PATCH /fabrics/{fabricName}/nmxc/domains  ->  /api/fm/fabrics/{fabricName}/nmxc/domains

Add/delete NMX-C domains.
"""
from __future__ import annotations

from typing import Any, TypedDict

from .._types import GpuAction, NmxcDomain
from ..client import get_client

class NmxcDomainsUpdateResult(TypedDict, total=False):
    """PATCH .../nmxc/domains — failure keys present only on partial failure (207)"""

    success: bool
    operation: str
    notRegistered: list[dict[str, object]]
    notDeregistered: list[dict[str, object]]


def update_nmxc_domains(fabricName: str, *, domains: list[NmxcDomain], operation: GpuAction) -> NmxcDomainsUpdateResult:
    client = get_client()
    path = f"fabrics/{fabricName}/nmxc/domains"
    body: dict[str, Any] = {}
    body["operation"] = operation
    body["domains"] = domains
    return client.call_api("PATCH", path, json_body=body or None)
