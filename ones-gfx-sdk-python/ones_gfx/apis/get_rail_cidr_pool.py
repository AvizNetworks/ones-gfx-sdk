"""
GET /fabrics/{fabricName}/rails/{railName}/cidrpool  ->  /api/fm/fabrics/{fabricName}/rails/{railName}/cidrpool

CIDRPool manifest (single YAML document) for one rail of a fabric.
"""
from __future__ import annotations

from typing import get_args

from .._types import RailCidrPoolAccept, RailCidrPoolResult
from ..client import Client
from ._validation import path_segment


def get_rail_cidr_pool(client: Client, fabricName: str, railName: str, *, accept: RailCidrPoolAccept = "application/yaml") -> str | RailCidrPoolResult:
    """Fetch the CIDRPool document of rail ``railName`` (e.g. ``"r0"``).

    ``accept`` behaves as in :func:`get_rail_cidr_pools`; the JSON form is
    ``{"fabricName", "railName", "yaml"}``.

    Raises:
        ValueError: a name is blank or ``accept`` is unsupported.
        requests.HTTPError: 404 ``INVALID_FABRIC`` or ``INVALID_RAIL``. For the
            YAML/text types the error body is plain text ``"<CODE>: <detail>"``.
    """
    if accept not in get_args(RailCidrPoolAccept):
        raise ValueError(f"accept must be one of {get_args(RailCidrPoolAccept)}, got {accept!r}.")
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/rails/{path_segment(railName, 'railName')}/cidrpool"
    return client.call_api("GET", path, extra_headers={"Accept": accept})
