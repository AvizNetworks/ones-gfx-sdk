"""
GET /fabrics/{fabricName}/shared-services  ->  /api/fm/fabrics/{fabricName}/shared-services

Read the fabric's shared-services tenant (GFX-5709).

Returns 200 with the tenant, its SVI gateway on each DPU and the VFs assigned
to it, or 404 SHARED_SERVICES_NOT_FOUND when the fabric has none. This call
does not contact the DPUs, so configured DPUs carry no `status` and there is
no `leaves` list.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def get_shared_services(client: Client, fabricName: str) -> dict[str, Any]:
    """Return the shared-services tenant on `fabricName`.

    Raises NotFoundError when the fabric has no shared-services tenant.
    """
    return client.call_api("GET", f"fabrics/{fabricName}/shared-services")
