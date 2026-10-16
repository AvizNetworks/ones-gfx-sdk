"""
DELETE /fabrics/{fabricName}/shared-services  ->  /api/fm/fabrics/{fabricName}/shared-services

Delete the fabric's shared-services tenant (GFX-5709).

Removes it from the leaves and DPUs, releases its DHCP subnet and IPAM pool.
Every VF must be unassigned first: FM answers 409 VFS_ASSIGNED and lists them
in `data` while any remain. If a delete fails partway the tenant is kept, so
calling this again is safe and is how a partial failure is retried.
"""
from __future__ import annotations

from typing import Any

from .._types import ApiResponseMessage
from ..client import Client


def delete_shared_services(client: Client, fabricName: str, *, idempotency_key: str | None = None) -> ApiResponseMessage:
    """Delete the shared-services tenant on `fabricName`."""
    headers: dict[str, str] = {}
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    return client.call_api("DELETE", f"fabrics/{fabricName}/shared-services", extra_headers=headers or None)
