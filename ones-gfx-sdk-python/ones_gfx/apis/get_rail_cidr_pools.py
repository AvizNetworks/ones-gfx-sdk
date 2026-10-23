"""
GET /fabrics/{fabricName}/rails/cidrpools  ->  /api/fm/fabrics/{fabricName}/rails/cidrpools

CIDRPool manifest (multi-document YAML, one document per rail) for a fabric.
"""
from __future__ import annotations

from typing import get_args

from .._types import RailCidrPoolAccept, RailCidrPoolsResult
from ..client import Client
from ._validation import path_segment


def get_rail_cidr_pools(client: Client, fabricName: str, *, accept: RailCidrPoolAccept = "application/yaml") -> str | RailCidrPoolsResult:
    """Fetch the rail CIDRPool manifest of ``fabricName``.

    ``accept`` selects the representation. ``application/yaml`` (default),
    ``application/x-yaml`` and ``text/plain`` return the manifest verbatim as a
    ``str``; ``application/json`` returns ``{"fabricName", "rails", "yaml"}``.

    Raises:
        ValueError: ``fabricName`` is blank or ``accept`` is unsupported.
        requests.HTTPError: 404 ``INVALID_FABRIC`` (unknown fabric) or
            ``RAIL_CIDRPOOL_NOT_FOUND`` (no East-West rails). For the YAML/text
            types the error body is plain text ``"<CODE>: <detail>"``
            (``exc.response.text``), not JSON.
    """
    if accept not in get_args(RailCidrPoolAccept):
        raise ValueError(f"accept must be one of {get_args(RailCidrPoolAccept)}, got {accept!r}.")
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/rails/cidrpools"
    return client.call_api("GET", path, extra_headers={"Accept": accept})
