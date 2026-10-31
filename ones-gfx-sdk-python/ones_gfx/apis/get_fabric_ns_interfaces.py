"""
GET /fabrics/{fabricName}  (ONES gateway, served at the root, not under /api/fm)

The fabric's topology YAML (the file it was deployed from), and the per-host
north-south interfaces derived from it.

The gateway answers 404 `NOT_FOUND` (JSON) when it has no topology for the
fabric. Failures raise `requests.HTTPError`, like every other API function.
"""
from __future__ import annotations

from collections.abc import Iterable
from typing import Any
from urllib.parse import quote

from ..client import Client, NotAuthenticatedError
from ..topology import extract_ns_interfaces


def get_fabric_yaml(client: Client, fabricName: str) -> str:
    """The topology YAML of `fabricName`, as text."""
    if not client.authToken:
        raise NotAuthenticatedError(
            "Not authenticated: call login() before calling Fabric Manager APIs."
        )
    if not fabricName or not fabricName.strip():
        raise ValueError("fabricName is required.")
    headers = client._auth_headers()
    headers["Accept"] = "application/x-yaml, application/json"
    response = client._session.request(
        method="GET",
        url=f"{client.baseUrl}/fabrics/{quote(fabricName, safe='')}",
        headers=headers,
        verify=client.verify_tls,
    )
    response.raise_for_status()
    return response.text


def get_fabric_ns_interfaces(
    client: Client,
    fabricName: str,
    *,
    leaf_patterns: Iterable[str] | str | None = None,
) -> dict[str, Any]:
    """Each host's N-S interfaces on `fabricName`.

    Fetches the topology and reduces it with `extract_ns_interfaces`; see that
    function for the result shape. `leaf_patterns` defaults to
    ``DEFAULT_NS_LEAF_PATTERNS`` (``^leaf-cn``, ``^leaf-fn``, ``^leaf-sn``).
    Render the result with `ones_gfx.ns_interfaces_to_yaml`.
    """
    return extract_ns_interfaces(
        get_fabric_yaml(client, fabricName), fabric_name=fabricName, leaf_patterns=leaf_patterns
    )
