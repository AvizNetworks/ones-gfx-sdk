"""
GET /fabrics/{fabricName}/ipam/config  ->  /api/fm/fabrics/{fabricName}/ipam/config

Per-fabric IPAM switches.
"""
from __future__ import annotations

from .._types import IpamConfig
from ..client import Client
from ._validation import path_segment


def get_ipam_config(client: Client, fabricName: str) -> IpamConfig:
    """Return ``{"fabricName", "enableIpam", "enableIpAssign", "updatedAt", "updatedBy"}``.

    A fabric that was never configured reports both switches ``False``.

    Raises:
        ValueError: ``fabricName`` is blank.
        requests.HTTPError: 404 ``FABRIC_NOT_FOUND``.
    """
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/ipam/config"
    return client.call_api("GET", path)
