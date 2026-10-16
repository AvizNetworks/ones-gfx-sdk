"""
Multi-call workflows.

Everything else in this SDK is one function per endpoint. These are the
sequences whose *order* is forced by the Fabric Manager, so getting them
wrong strands state rather than returning an error.
"""
from __future__ import annotations

from typing import Any

from .apis.delete_shared_services import delete_shared_services
from .apis.delete_tenant import delete_tenant
from .apis.get_gpu_tenant_mappings import get_gpu_tenant_mappings
from .apis.get_shared_services import get_shared_services
from .apis.list_vf_interfaces import list_vf_interfaces
from .apis.modify_gpu_allocations import modify_gpu_allocations
from .apis.unassign_vf_interface import unassign_vf_interface
from .client import Client
from .exceptions import APIError, NotFoundError

__all__ = ["CleanUpResult", "clean_up_tenant"]


class CleanUpResult(dict):
    """What `clean_up_tenant` did, and what it could not finish.

    `clean` is True when nothing failed. `failures` names the steps that did,
    so the caller can re-run; `shared_services` records whether the tenant was
    the fabric's shared-services tenant.
    """

    @property
    def clean(self) -> bool:
        return not self["failures"]


def _tenant_vf_ids(client: Client, fabric_name: str, server: str, tenant_name: str) -> list[str]:
    """The VF ids on `server` currently assigned to `tenant_name`."""
    try:
        entries = list_vf_interfaces(client, fabric_name, server)
    except (NotFoundError, APIError):
        return []
    rows = entries if isinstance(entries, list) else (entries or {}).get("data") or []
    ids: list[str] = []
    for row in rows:
        if not isinstance(row, dict):
            continue
        if row.get("tenantName") != tenant_name:
            continue
        vf_id = row.get("ifName") or row.get("if_name") or row.get("serverIf")
        if vf_id:
            ids.append(vf_id)
    return ids


def _discover_servers(client: Client, fabric_name: str, tenant_name: str) -> list[str]:
    """The GPU servers this tenant holds, from gpuTenantMappings."""
    try:
        rows = get_gpu_tenant_mappings(client, fabric_name, tenantName=tenant_name)
    except (NotFoundError, APIError):
        return []
    servers: list[str] = []
    for row in rows or []:
        if not isinstance(row, dict):
            continue
        if row.get("tenantName") not in (None, tenant_name):
            continue
        name = row.get("gpuHostname") or row.get("serverName") or row.get("hostname")
        if name and name not in servers:
            servers.append(name)
    return servers


def clean_up_tenant(client: Client, fabric_name: str, tenant_name: str, *, servers: list[str] | None = None) -> CleanUpResult:
    """Tear a DPU tenant down in the one order the FM accepts.

    `DELETE /tenants/{t}` answers 409 GPUS_ALLOCATED while GPUs are still
    attached, and a VF cannot be freed once its server has left the tenant.
    So the order is forced, and is the reverse of provisioning:

        1. unassign every VF the tenant holds
        2. detach every GPU server from the tenant
        3. delete the tenant

    A shared-services tenant is deleted through its own route instead, which
    also releases its Kea subnet and IPAM pool.

    Every step tolerates "already gone", so re-running a half-finished
    teardown finishes it rather than failing. Steps that genuinely fail are
    listed in the result's `failures` instead of raising, so one bad step
    never strands the rest.
    """
    failures: list[str] = []
    found = servers if servers is not None else _discover_servers(client, fabric_name, tenant_name)

    # 1. VFs ----------------------------------------------------------------
    freed: list[str] = []
    for server in found:
        for vf_id in _tenant_vf_ids(client, fabric_name, server, tenant_name):
            try:
                unassign_vf_interface(client, fabric_name, server, vf_id)
                freed.append(f"{server}/{vf_id}")
            except NotFoundError:
                pass  # already free
            except APIError:
                failures.append(f"vf {server}/{vf_id}")

    # 2. GPU servers --------------------------------------------------------
    # Addressed as fabric -> SU -> server -> port names, so the ports come
    # from the fabric's own GPU inventory, grouped the same way.
    if found:
        # {su: {hostname: {"gpus": [...]}}} — the Java DeviceGpus shape.
        suid: dict[str, dict[str, dict[str, list[str]]]] = {}
        try:
            for row in get_gpu_tenant_mappings(client, fabric_name, tenantName=tenant_name) or []:
                if not isinstance(row, dict):
                    continue
                host = row.get("gpuHostname") or row.get("serverName")
                if host not in found:
                    continue
                su = str(row.get("suId", 0))
                port = row.get("gpuPort") or row.get("gpuName")
                if port:
                    suid.setdefault(su, {}).setdefault(host, {"gpus": []})["gpus"].append(port)
        except APIError:
            suid = {}
        if suid:
            try:
                modify_gpu_allocations(client, fabric_name, tenant_name, suid=suid,
                                       operation="DELETE", configScope="WHOLE_SERVER")
            except NotFoundError:
                pass
            except APIError:
                failures.append("gpu detach")

    # 3. Tenant -------------------------------------------------------------
    is_shared = False
    try:
        shared = get_shared_services(client, fabric_name)
        is_shared = isinstance(shared, dict) and shared.get("tenantName") == tenant_name
    except (NotFoundError, APIError):
        is_shared = False
    try:
        if is_shared:
            delete_shared_services(client, fabric_name)
        else:
            delete_tenant(client, fabric_name, tenant_name)
    except NotFoundError:
        pass  # already gone
    except APIError:
        failures.append("shared-services delete" if is_shared else "tenant delete")

    return CleanUpResult(
        fabricName=fabric_name,
        tenantName=tenant_name,
        servers=found,
        vfsFreed=freed,
        sharedServices=is_shared,
        failures=failures,
    )
