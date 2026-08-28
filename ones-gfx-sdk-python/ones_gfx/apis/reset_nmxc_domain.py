"""
POST /api/nmxc/domains/{domainId}/reset  ->  /api/fm/api/nmxc/domains/{domainId}/reset

Soft-reset an NMX-C domain.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class NmxcResetResult(TypedDict):
    """POST /api/nmxc/domains/{domainId}/reset (202)"""

    domainId: str
    action: str
    message: str


def reset_nmxc_domain(client: Client, domainId: str) -> NmxcResetResult:
    path = f"api/nmxc/domains/{domainId}/reset"
    return client.call_api("POST", path)
