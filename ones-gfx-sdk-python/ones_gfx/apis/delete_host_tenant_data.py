"""
DELETE /delHostTenantData/{fabricName}/{tenantName}  ->  /api/fm/delHostTenantData/{fabricName}/{tenantName}

Remove a host-tenant data record.
"""
from __future__ import annotations


from ..client import get_client


def delete_host_tenant_data(fabricName: str, tenantName: str) -> str:
    client = get_client()
    path = f"delHostTenantData/{fabricName}/{tenantName}"
    return client.call_api("DELETE", path)
