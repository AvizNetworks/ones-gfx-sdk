"""
Deploy a fabric made with create_fabric_config — the ONES UI's "Deploy Fabric"
action, as a sequence of config-service calls:

1. upload_fabric_device_ips for every device
2. validate_fabric_switches for switches with apply_config (skipped if none)
3. validate_fabric_servers for hosts/servers/DPUs (skipped if none)
4. update_fabric_inventory for every device
5. fetch the fabric YAML (or use custom_yaml), fill in each device's IP and
   credentials, and push_fabric_config it
6. update_fabric_config_status to "Deployed"

Any failure stops the sequence. A validation failure raises
FabricValidationError before anything is configured.
"""
from __future__ import annotations

from typing import Any, Callable

from .._types import FabricDeployResult, FabricDevice, ValidateDeviceResult
from ..client import Client
from ..exceptions import FabricValidationError
from ._fabric_config import apply_config, fabric_yaml_to_document, is_server, is_switch
from .get_fabric_yaml import get_fabric_yaml
from .push_fabric_config import push_fabric_config
from .update_fabric_config_status import update_fabric_config_status
from .update_fabric_inventory import update_fabric_inventory
from .upload_fabric_device_ips import upload_fabric_device_ips
from .validate_fabric_servers import validate_fabric_servers
from .validate_fabric_switches import validate_fabric_switches


def _failures(results: list[ValidateDeviceResult], passed: Callable[[ValidateDeviceResult], bool]) -> list[str]:
    failures = []
    for res in results or []:
        where = f"{res.get('hostname')} ({res.get('ip')})"
        if (res.get("error") or "").strip():
            failures.append(f"{where}: {res['error']}")
        elif not passed(res):
            failures.append(f"{where}: validation did not report success")
    return failures


def deploy_fabric(client: Client, fabricName: str, devices: list[FabricDevice], *, description: str = "", deploymentType: str = "DEFAULT", instance: str = "fm", custom_yaml: str | None = None) -> FabricDeployResult:
    """Run the full deploy sequence; returns every step's response.

    Args:
        devices: every switch, server and DPU in the fabric, as
            {hostname, ip, username, password, device_type?, device_role?,
            apply_config?} — the devices.json format.
        deploymentType: DEFAULT, PARTIAL_CONFIG or BROWNFIELD.
        custom_yaml: use this fabric YAML instead of the generated one.

    Raises:
        FabricValidationError: a switch or server failed validation.
        requests.HTTPError: a step was rejected by the server.
    """
    if not devices:
        raise ValueError("devices must list every switch/server/DPU in the fabric")
    out: FabricDeployResult = {}
    out["uploadIps"] = upload_fabric_device_ips(client, fabricName, devices)

    switches = [d for d in devices if is_switch(d) and apply_config(d)]
    if switches:
        out["switchResults"] = validate_fabric_switches(client, fabricName, switches)
        failed = _failures(out["switchResults"], lambda r: bool((r.get("build") or "").strip()))
        if failed:
            raise FabricValidationError("switch", failed, out["switchResults"])

    servers = [d for d in devices if is_server(d)]
    if servers:
        out["serverResults"] = validate_fabric_servers(client, fabricName, servers)
        failed = _failures(out["serverResults"], lambda r: bool((r.get("os") or "").strip()))
        if failed:
            raise FabricValidationError("server", failed, out["serverResults"])

    out["inventory"] = update_fabric_inventory(client, fabricName, devices)

    raw_yaml = custom_yaml if custom_yaml and custom_yaml.strip() else get_fabric_yaml(client, fabricName)
    document: Any = fabric_yaml_to_document(raw_yaml, devices)
    out["configPushed"] = push_fabric_config(client, fabricName, document, instance=instance)

    out["status"] = update_fabric_config_status(
        client, name=fabricName, status="Deployed", description=description, deploymentType=deploymentType,
    )
    return out
