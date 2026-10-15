"""
POST /api/config/addFabricData  (ONES UI config service)

Generate a fabric — topology, YAML and inventory — in Draft state. This is the
call the ONES UI makes, and the starting point for deploy_fabric. (By
contrast, add_fabric_data posts to Fabric Manager and only stores the record.)
"""
from __future__ import annotations

from typing import Any

from ..client import Client
from ._fabric_config import host_map_from_hosts_per_su


def create_fabric_config(client: Client, *, name: str, type: str, hostsPerSu: str, numOfSus: int = 1, maxNumOfSus: int | None = None, status: str = "Draft", description: str = "", hostMap: dict[str, str] | None = None, nodeType: str = "dgx", tenantCtrl: str = "ones", simulationId: int = 1, instance: str = "fm", enableEW: bool = False, startingSubnetGpu: str = "", enableNS: bool = False, nsOperatingSystem: str = "cumulus", dedicatedStorage: bool = False, startingSubnetCpu: str | None = None, startingSubnetStorage: str | None = None) -> Any:
    """Create a fabric through the config service.

    Args:
        type: "NVIDIA SpX RA 1.3", "NVIDIA SpX RA 2.1" or "Aviz RA 1.0".
        hostsPerSu: SU index -> host count, e.g. "{0:2}" or "{0:1,1:2}";
            hostMap is derived from it unless given.
        maxNumOfSus: defaults to numOfSus.
        nodeType: dgx, gb200, gb300, b300_32, b300_64, rtxpro_4, rtxpro_8.
        enableEW / startingSubnetGpu: east-west (GPU) network; the subnet is
            required when enabled.
        enableNS / startingSubnetCpu: north-south (front-end) network; the
            subnet is required when enabled.
        nsOperatingSystem: "cumulus" or "sonic" for the front-end/storage
            switches; must match the real devices or deploy validation fails.
        dedicatedStorage / startingSubnetStorage: separate storage network;
            the subnet is required when enabled.

    Returns the server's response (it lists the generated inventory).
    """
    if enableEW and not startingSubnetGpu.strip():
        raise ValueError("startingSubnetGpu is required when enableEW is True")
    if enableNS and not (startingSubnetCpu or "").strip():
        raise ValueError("startingSubnetCpu is required when enableNS is True")
    if dedicatedStorage and not (startingSubnetStorage or "").strip():
        raise ValueError("startingSubnetStorage is required when dedicatedStorage is True")
    body: dict[str, Any] = {
        "name": name,
        "type": type,
        "status": status,
        "description": description,
        "numOfSus": numOfSus,
        "maxNumOfSus": maxNumOfSus if maxNumOfSus is not None else numOfSus,
        "hostMap": hostMap or host_map_from_hosts_per_su(hostsPerSu),
        "startingSubnetGpu": startingSubnetGpu,
        "simulationId": simulationId,
        "enableEW": enableEW,
        "suHostCnt": hostsPerSu,
        "tenant": tenantCtrl,
        "instance": instance,
        "enableNS": enableNS,
        # Always ONES-controlled, matching the UI default for RA fabrics.
        "isOnesControlled": True,
        # The UI toggles dedicatedStorage and frontendStorage together;
        # userandstorage is the shared (NS on, no dedicated storage) case.
        "userandstorage": enableNS and not dedicatedStorage,
        "dedicatedStorage": dedicatedStorage,
        "frontendStorage": dedicatedStorage,
        "nodeType": nodeType,
        # East-west is always Cumulus in the UI.
        "operatingSystemE": "cumulus",
        "operatingSystemF": nsOperatingSystem,
        "operatingSystemS": nsOperatingSystem,
    }
    if startingSubnetCpu:
        body["startingSubnetCpu"] = startingSubnetCpu
    if startingSubnetStorage:
        body["startingSubnetStorage"] = startingSubnetStorage
    return client.call_config_api("POST", "api/config/addFabricData", json_body=body)
