"""
POST /api/nmxc/domains/{domainId}/factory-reset  ->  /api/fm/api/nmxc/domains/{domainId}/factory-reset

Factory-reset an NMX-C domain (destructive; creates backup first).
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class NmxcFactoryResetResult(TypedDict, total=False):
    """POST /api/nmxc/domains/{domainId}/factory-reset"""

    domainId: str
    action: str
    message: str
    backupApplied: str


def factory_reset_nmxc_domain(client: Client, domainId: str) -> NmxcFactoryResetResult:
    path = f"api/nmxc/domains/{domainId}/factory-reset"
    return client.call_api("POST", path)
