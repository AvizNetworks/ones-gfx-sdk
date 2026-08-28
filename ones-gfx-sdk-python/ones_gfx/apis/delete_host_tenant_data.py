"""
DELETE /delHostTenantData/{fabricName}/{tenantName}  ->  /api/fm/delHostTenantData/{fabricName}/{tenantName}

Remove a host-tenant data record.
"""
from __future__ import annotations


from ..client import Client


def delete_host_tenant_data(client: Client, fabricName: str, tenantName: str) -> str:
    path = f"delHostTenantData/{fabricName}/{tenantName}"
    return client.call_api("DELETE", path)
