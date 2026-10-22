"""
DELETE /fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign  ->  /api/fm/...

Return a DPU VF interface to free status.

Order matters on teardown: a VF cannot be freed once its server has left the
tenant, so unassign every VF before detaching GPU servers or deleting the
tenant. `clean_up_tenant` in ones_gfx.workflows does this in order.
"""
from __future__ import annotations

from .._types import ApiResponseMessage
from ..client import Client


def unassign_vf_interface(client: Client, fabricName: str, serverName: str, vfId: str, *, idempotency_key: str | None = None) -> ApiResponseMessage:
    """Unassign VF `vfId` on `serverName`, returning it to free status."""
    path = f"fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign"
    headers: dict[str, str] = {}
    if idempotency_key is not None:
        headers["Idempotency-Key"] = idempotency_key
    return client.call_api("DELETE", path, extra_headers=headers or None)
