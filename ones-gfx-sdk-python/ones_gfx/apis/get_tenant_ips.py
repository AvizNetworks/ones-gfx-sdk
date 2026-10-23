"""
GET /fabrics/{fabricName}/tenants/{tenantName}/ip-addresses  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/ip-addresses

Paged IP allocations of a tenant, plus its pool usage.
"""
from __future__ import annotations

from typing import Any, get_args

from .._types import IpAssignStatus, TenantIpAllocations
from ..client import Client
from ._validation import path_segment


def get_tenant_ips(client: Client, fabricName: str, tenantName: str, *, page: int | None = None, size: int | None = None, status: IpAssignStatus | None = None) -> TenantIpAllocations:
    """List a tenant's allocations.

    ``page`` is zero-based (server default 0), ``size`` defaults to 50, and
    ``status`` filters on ``ALLOCATED_ONLY`` / ``ASSIGNED`` / ``IP_ASSIGN_FAILED``.
    Returns ``{"fabricName", "tenantName", "pool", "allocations", "pagination"}``.

    Raises:
        ValueError: ``page`` < 0, ``size`` < 1, or an unknown ``status``. (The
            server would silently clamp the numbers; the SDK reports them.)
        requests.HTTPError: 404 ``IPAM_NOT_CONFIGURED``.
    """
    if page is not None and (isinstance(page, bool) or not isinstance(page, int) or page < 0):
        raise ValueError(f"page must be an integer >= 0, got {page!r}.")
    if size is not None and (isinstance(size, bool) or not isinstance(size, int) or size < 1):
        raise ValueError(f"size must be an integer >= 1, got {size!r}.")
    if status is not None and status not in get_args(IpAssignStatus):
        raise ValueError(f"status must be one of {get_args(IpAssignStatus)}, got {status!r}.")
    path = f"fabrics/{path_segment(fabricName, 'fabricName')}/tenants/{path_segment(tenantName, 'tenantName')}/ip-addresses"
    params: dict[str, Any] = {}
    if page is not None:
        params["page"] = page
    if size is not None:
        params["size"] = size
    if status is not None:
        params["status"] = status
    return client.call_api("GET", path, params=params or None)
