"""
GET /fabrics/{fabricName}/servers/{serverName}/vf-interfaces  ->  /api/fm/fabrics/{fabricName}/servers/{serverName}/vf-interfaces

List a GPU server's DPU VF interfaces and what they are assigned to.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def list_vf_interfaces(client: Client, fabricName: str, serverName: str) -> Any:
    """Return the VF interfaces on `serverName`'s DPU.

    A VF is identified either by its DPU interface name (`pf0vf0_if`) or by
    the host interface name (`enp193s0f1np1`); both appear in each entry.
    """
    return client.call_api("GET", f"fabrics/{fabricName}/servers/{serverName}/vf-interfaces")
