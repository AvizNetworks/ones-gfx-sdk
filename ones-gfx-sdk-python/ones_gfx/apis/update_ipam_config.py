"""
PATCH /fabrics/{fabricName}/ipam/config  ->  /api/fm/fabrics/{fabricName}/ipam/config

Update the fabric's IPAM configuration.

IPAM hands out addresses in one of two modes: `static` (an address is picked
per server and, with `enableIpAssign`, pushed to the host over SSH) or `dhcp`
(the DHCP server leases it). With IPAM on, exactly one mode is enabled; with
IPAM off, neither. A mode object left out counts as disabled, so send only
the mode you use. Settings may only accompany an enabled mode, and an enabled
mode's omitted settings keep their saved values.

`leaseTimeSeconds` is the only DHCP setting — renew and rebind follow it at
the RFC 2131 defaults of 50% and 87.5%.

FM rejects any field outside this set rather than ignoring it, so a mode
setting sent at the top level fails clearly. Only externally managed
(non-ONES-controlled) fabrics may be changed; others get 409
(`FABRIC_NOT_EXTERNALLY_MANAGED`). 404 when the fabric is unknown.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def update_ipam_config(
    client: Client,
    fabricName: str,
    *,
    enableIpam: bool,
    static: dict[str, Any] | None = None,
    dhcp: dict[str, Any] | None = None,
    force: bool | None = None,
) -> dict[str, Any]:
    """Update IPAM on `fabricName`.

    `static` is `{"enabled": bool, "enableIpAssign": bool}` and `dhcp` is
    `{"enabled": bool, "leaseTimeSeconds": int}`; pass at most one enabled.
    `force` pushes the change through checks FM would otherwise fail on.

    Returns the stored configuration in the same shape.
    """
    path = f"fabrics/{fabricName}/ipam/config"
    params: dict[str, Any] = {}
    if force is not None:
        # Spelled out rather than passed as a bool: requests would render
        # Python's True as "True", and the lowercase form is what FM's query
        # binding documents.
        params["force"] = "true" if force else "false"
    body: dict[str, Any] = {"enableIpam": enableIpam}
    if static is not None:
        body["static"] = static
    if dhcp is not None:
        body["dhcp"] = dhcp
    return client.call_api("PATCH", path, params=params or None, json_body=body)
