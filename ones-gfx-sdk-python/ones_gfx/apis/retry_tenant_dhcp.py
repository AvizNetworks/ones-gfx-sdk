"""
POST /fabrics/{fabricName}/tenants/{tenantName}/dhcp/retry  ->  /api/fm/fabrics/{fabricName}/tenants/{tenantName}/dhcp/retry

Retry a tenant's failed DHCP provisioning by re-applying its desired state
to Kea.

Idempotent: the fragment is rendered from the database and applied whole, so
calling this when nothing is wrong re-applies the same configuration.

404 (`DHCP_SUBNET_NOT_FOUND`) when DHCP does not apply to the tenant, 400 on
a validation failure, 502 when Kea cannot be reached.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def retry_tenant_dhcp(client: Client, fabricName: str, tenantName: str) -> dict[str, Any]:
    """Re-apply `tenantName`'s DHCP subnet. Returns the re-apply summary."""
    path = f"fabrics/{fabricName}/tenants/{tenantName}/dhcp/retry"
    return client.call_api("POST", path)
