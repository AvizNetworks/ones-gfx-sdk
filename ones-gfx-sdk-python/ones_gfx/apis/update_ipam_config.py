"""
PATCH /fabrics/{fabricName}/ipam/config  ->  /api/fm/fabrics/{fabricName}/ipam/config

Enable/disable IPAM and automatic IP assignment on a fabric.
"""
from __future__ import annotations

from typing import Any

from .._types import IpamConfig
from ..client import Client
from ._validation import path_segment, require_bool


def update_ipam_config(client: Client, fabricName: str, *, enableIpam: bool, enableIpAssign: bool, force: bool | None = None) -> IpamConfig:
    """Set both IPAM switches and return the stored config.

    ``force=True`` allows turning IPAM off while active pools exist; the pools
    are disabled and their allocations moved to history. Omit it (server default
    ``False``) to have that case rejected.

    Raises:
        ValueError: a switch is not a bool, or ``enableIpAssign`` is set without
            ``enableIpam`` (the server rejects that too).
        requests.HTTPError: 404 ``FABRIC_NOT_FOUND``; 409
            ``FABRIC_NOT_EXTERNALLY_MANAGED`` (ONES-controlled fabrics cannot be
            changed) or ``IPAM_HAS_ACTIVE_POOLS`` (disabling without ``force``);
            400 ``INVALID_REQUEST_BODY`` / ``INVALID_IPAM_CONFIG``. The error code
            is in ``exc.response.json()["error"]["code"]``.
    """
    require_bool(enableIpam, "enableIpam")
    require_bool(enableIpAssign, "enableIpAssign")
    if force is not None:
        require_bool(force, "force")
    if enableIpAssign and not enableIpam:
        raise ValueError("enableIpAssign requires enableIpam to be True.")
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/ipam/config"
    params: dict[str, Any] = {}
    if force is not None:
        params["force"] = "true" if force else "false"
    body: dict[str, Any] = {"enableIpam": enableIpam, "enableIpAssign": enableIpAssign}
    return client.call_api("PATCH", path, params=params or None, json_body=body)
