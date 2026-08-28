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
from .get_operation_webhook_status import get_operation_webhook_status
from .add_fabric_data import add_fabric_data
from .edit_fabric_data import edit_fabric_data
from .delete_fabric_data import delete_fabric_data
from .get_all_fabrics import get_all_fabrics
from .get_fabric_by_name import get_fabric_by_name
from .get_fabrics import get_fabrics
from .update_fabric_status import update_fabric_status
from .get_fabric_devices import get_fabric_devices
from .get_devices_by_layer import get_devices_by_layer
from .add_fabric_sim import add_fabric_sim
from .get_all_fabric_sims import get_all_fabric_sims
from .get_fabric_sim_by_name import get_fabric_sim_by_name
from .delete_fabric_sim import delete_fabric_sim
from .update_fabric_sim_status import update_fabric_sim_status
from .add_inventory_data import add_inventory_data
from .update_inventory_data import update_inventory_data
from .edit_inventory_data import edit_inventory_data
from .get_all_inventory import get_all_inventory
from .get_inventory_by_fabric_name import get_inventory_by_fabric_name
from .get_fm_inventory import get_fm_inventory
from .get_inventory_ports import get_inventory_ports
from .get_inventory_hosts import get_inventory_hosts
from .inventory_sync import inventory_sync
from .validate_ufm_creds import validate_ufm_creds
from .update_nmxc_domains import update_nmxc_domains
from .probe_nmxc_domains import probe_nmxc_domains
from .get_nmxc_inventory import get_nmxc_inventory
from .create_vpc_peering import create_vpc_peering
from .delete_vpc_peering import delete_vpc_peering
from .modify_gpu_allocations import modify_gpu_allocations
from .assign_gpu_ports import assign_gpu_ports
from .auto_allocate_gpus_to_tenants import auto_allocate_gpus_to_tenants
from .get_all_gpus_list import get_all_gpus_list
from .get_gpus_by_host import get_gpus_by_host
from .get_gpu_tenant_mappings import get_gpu_tenant_mappings
from .get_gpu_allocation_history import get_gpu_allocation_history
from .get_available_servers import get_available_servers
from .upload_day1_config import upload_day1_config
from .upload_ui_object import upload_ui_object
from .get_ui_object import get_ui_object
from .get_last_orchestrated_intent_name import get_last_orchestrated_intent_name
from .get_day1_config_status import get_day1_config_status
from .get_intent_validation import get_intent_validation
from .get_intent_derivation_logs import get_intent_derivation_logs
from .get_upload_status import get_upload_status
from .is_device_alive import is_device_alive
from .add_device_facts import add_device_facts
from .remove_device_facts import remove_device_facts
from .get_version import get_version
from .get_imgmgmt_status import get_imgmgmt_status
from .upgrade_nos_image import upgrade_nos_image
from .enable_ztp_upgrade import enable_ztp_upgrade
from .reboot_request import reboot_request
from .backup_config import backup_config
from .restore_config import restore_config
from .configs_list_to_restore import configs_list_to_restore
from .fetch_device_backup_files import fetch_device_backup_files
from .get_config_diff import get_config_diff
from .replace_config import replace_config
from .get_config import get_config
from .fill_bootstrap_config import fill_bootstrap_config
from .trigger_bootstrap_config import trigger_bootstrap_config
from .get_bootstrap_info import get_bootstrap_info
from .get_all_bootstrap_batches import get_all_bootstrap_batches
from .get_bootstrap_batch import get_bootstrap_batch
from .get_device_bootstrap_stages import get_device_bootstrap_stages
from .fill_rma_config import fill_rma_config
from .trigger_rma import trigger_rma
from .get_rma_info import get_rma_info
from .get_rma_status import get_rma_status
from .upload_file import upload_file
from .get_files import get_files
from .delete_file import delete_file
from .update_role_info import update_role_info
from .get_log_level import get_log_level
from .set_log_level import set_log_level
from .get_status import get_status
from .start_streaming import start_streaming
from .stop_streaming import stop_streaming
from .get_controller_version import get_controller_version
from .get_controller_version_internal import get_controller_version_internal
from .netops_fabric import netops_fabric
from .netops_device import netops_device
from .add_host_tenant_data import add_host_tenant_data
from .delete_host_tenant_data import delete_host_tenant_data
from .get_host_tenants_list import get_host_tenants_list
from .update_hosts import update_hosts
from .get_hosts_list import get_hosts_list
from .reset_nmxc_domain import reset_nmxc_domain
from .factory_reset_nmxc_domain import factory_reset_nmxc_domain

__all__ = [
    "create_tenant",
    "list_tenants",
    "get_tenant",
    "update_tenant",
    "delete_tenant",
    "get_operation",
    "get_operation_webhook_status",
    "add_fabric_data",
    "edit_fabric_data",
    "delete_fabric_data",
    "get_all_fabrics",
    "get_fabric_by_name",
    "get_fabrics",
    "update_fabric_status",
    "get_fabric_devices",
    "get_devices_by_layer",
    "add_fabric_sim",
    "get_all_fabric_sims",
    "get_fabric_sim_by_name",
    "delete_fabric_sim",
    "update_fabric_sim_status",
    "add_inventory_data",
    "update_inventory_data",
    "edit_inventory_data",
    "get_all_inventory",
    "get_inventory_by_fabric_name",
    "get_fm_inventory",
    "get_inventory_ports",
    "get_inventory_hosts",
    "inventory_sync",
    "validate_ufm_creds",
    "update_nmxc_domains",
    "probe_nmxc_domains",
    "get_nmxc_inventory",
    "create_vpc_peering",
    "delete_vpc_peering",
    "modify_gpu_allocations",
    "assign_gpu_ports",
    "auto_allocate_gpus_to_tenants",
    "get_all_gpus_list",
    "get_gpus_by_host",
    "get_gpu_tenant_mappings",
    "get_gpu_allocation_history",
    "get_available_servers",
    "upload_day1_config",
    "upload_ui_object",
    "get_ui_object",
    "get_last_orchestrated_intent_name",
    "get_day1_config_status",
    "get_intent_validation",
    "get_intent_derivation_logs",
    "get_upload_status",
    "is_device_alive",
    "add_device_facts",
    "remove_device_facts",
    "get_version",
    "get_imgmgmt_status",
    "upgrade_nos_image",
    "enable_ztp_upgrade",
    "reboot_request",
    "backup_config",
    "restore_config",
    "configs_list_to_restore",
    "fetch_device_backup_files",
    "get_config_diff",
    "replace_config",
    "get_config",
    "fill_bootstrap_config",
    "trigger_bootstrap_config",
    "get_bootstrap_info",
    "get_all_bootstrap_batches",
    "get_bootstrap_batch",
    "get_device_bootstrap_stages",
    "fill_rma_config",
    "trigger_rma",
    "get_rma_info",
    "get_rma_status",
    "upload_file",
    "get_files",
    "delete_file",
    "update_role_info",
    "get_log_level",
    "set_log_level",
    "get_status",
    "start_streaming",
    "stop_streaming",
    "get_controller_version",
    "get_controller_version_internal",
    "netops_fabric",
    "netops_device",
    "add_host_tenant_data",
    "delete_host_tenant_data",
    "get_host_tenants_list",
    "update_hosts",
    "get_hosts_list",
    "reset_nmxc_domain",
    "factory_reset_nmxc_domain",
]
