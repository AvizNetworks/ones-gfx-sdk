"""Fabric Manager API functions — one module per endpoint.

Every function routes through ONESClient.call_api, which enforces the
auth guard (login required) and the /api/fm/ path prefix.
"""
from __future__ import annotations

from .create_tenant import create_tenant
from .list_tenants import list_tenants
from .get_tenant import get_tenant
from .update_tenant import update_tenant
from .delete_tenant import delete_tenant
from .get_operation import get_operation
from .get_fabrics import get_fabrics
from .get_inventory_ports import get_inventory_ports
from .get_inventory_hosts import get_inventory_hosts
from .inventory_sync import inventory_sync
from .validate_ufm_creds import validate_ufm_creds
from .update_nmxc_domains import update_nmxc_domains
from .probe_nmxc_domains import probe_nmxc_domains
from .get_nmxc_inventory import get_nmxc_inventory
from .create_vpc_peering import create_vpc_peering
from .modify_gpu_allocations import modify_gpu_allocations
from .assign_gpu_ports import assign_gpu_ports
from .get_available_servers import get_available_servers
from .reset_nmxc_domain import reset_nmxc_domain
from .factory_reset_nmxc_domain import factory_reset_nmxc_domain

__all__ = [
    "create_tenant",
    "list_tenants",
    "get_tenant",
    "update_tenant",
    "delete_tenant",
    "get_operation",
    "get_fabrics",
    "get_inventory_ports",
    "get_inventory_hosts",
    "inventory_sync",
    "validate_ufm_creds",
    "update_nmxc_domains",
    "probe_nmxc_domains",
    "get_nmxc_inventory",
    "create_vpc_peering",
    "modify_gpu_allocations",
    "assign_gpu_ports",
    "get_available_servers",
    "reset_nmxc_domain",
    "factory_reset_nmxc_domain",
]
