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

from .create_shared_services import create_shared_services
from .get_shared_services import get_shared_services
from .delete_shared_services import delete_shared_services
from .list_vf_interfaces import list_vf_interfaces
from .assign_vf_interface import assign_vf_interface
from .unassign_vf_interface import unassign_vf_interface

# --- Network policy (GFX-4721 / GFX-5327 / GFX-5394 / GFX-5536) ---
from .list_network_policies import list_network_policies
from .create_network_policy import create_network_policy
from .patch_network_policy import patch_network_policy
from .delete_network_policy import delete_network_policy
# --- Webhook registration: global, fabric and tenant scopes ---
from .register_webhook import register_webhook
from .list_webhooks import list_webhooks
from .get_webhook import get_webhook
from .update_webhook import update_webhook
from .delete_webhook import delete_webhook
from .register_fabric_webhook import register_fabric_webhook
from .get_fabric_webhook import get_fabric_webhook
from .update_fabric_webhook import update_fabric_webhook
from .delete_fabric_webhook import delete_fabric_webhook
from .register_tenant_webhook import register_tenant_webhook
from .get_tenant_webhook import get_tenant_webhook
from .update_tenant_webhook import update_tenant_webhook
from .delete_tenant_webhook import delete_tenant_webhook
# --- IPAM address queries (GFX-5504) ---
from .get_server_ips import get_server_ips
from .get_tenant_ips import get_tenant_ips
from .get_tenant_compute_ip import get_tenant_compute_ip
# --- Tenant DHCP (GFX-5504 / GFX-5611) ---
from .reconcile_fabric_dhcp import reconcile_fabric_dhcp
from .get_tenant_dhcp import get_tenant_dhcp
from .retry_tenant_dhcp import retry_tenant_dhcp
# --- IPAM configuration (GFX-5229 / GFX-5449) ---
from .get_ipam_config import get_ipam_config
from .update_ipam_config import update_ipam_config
# --- Rail CIDRPools (GFX-5264) ---
from .get_rail_cidr_pools import get_rail_cidr_pools
from .get_rail_cidr_pool import get_rail_cidr_pool
# --- Inventory delete by hostname (GFX-5012) ---

__all__ = [
    "create_shared_services",
    "get_shared_services",
    "delete_shared_services",
    "list_vf_interfaces",
    "assign_vf_interface",
    "unassign_vf_interface",

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
    # Network policy (GFX-4721 / GFX-5327 / GFX-5394 / GFX-5536)
    "list_network_policies",
    "create_network_policy",
    "patch_network_policy",
    "delete_network_policy",
    # Webhook registration: global, fabric and tenant scopes
    "register_webhook",
    "list_webhooks",
    "get_webhook",
    "update_webhook",
    "delete_webhook",
    "register_fabric_webhook",
    "get_fabric_webhook",
    "update_fabric_webhook",
    "delete_fabric_webhook",
    "register_tenant_webhook",
    "get_tenant_webhook",
    "update_tenant_webhook",
    "delete_tenant_webhook",
    # IPAM address queries (GFX-5504)
    "get_server_ips",
    "get_tenant_ips",
    "get_tenant_compute_ip",
    # Tenant DHCP (GFX-5504 / GFX-5611)
    "reconcile_fabric_dhcp",
    "get_tenant_dhcp",
    "retry_tenant_dhcp",
    # IPAM configuration (GFX-5229 / GFX-5449)
    "get_ipam_config",
    "update_ipam_config",
    # Rail CIDRPools (GFX-5264)
    "get_rail_cidr_pools",
    "get_rail_cidr_pool",
    # Inventory delete by hostname (GFX-5012)
]
