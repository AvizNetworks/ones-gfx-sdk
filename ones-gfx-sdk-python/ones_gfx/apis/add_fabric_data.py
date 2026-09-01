"""
POST /addFabricData  ->  /api/fm/addFabricData

Create fabric.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def add_fabric_data(client: Client, *, name: str, id: int | None = None, type: str | None = None, status: str | None = None, description: str | None = None, orchestrationStatus: str | None = None, numOfSus: int | None = None, maxNumOfSus: int | None = None, dedicated: bool | None = None, hybrid: bool | None = None, isDPUFabric: bool | None = None, startingSubnetGpu: int | None = None, startingSubnetCpu: str | None = None, startingSubnetTenants: str | None = None, startingSubnetStorage: str | None = None, simulationId: int | None = None, intent: str | None = None, ewTenantAware: bool | None = None, storageTenantAware: bool | None = None, isOnesControlled: bool | None = None, suHostCnt: str | None = None, isVxlanFabric: bool | None = None, nodeType: str | None = None, deploymentType: str | None = None, spineEvpnConfigured: bool | None = None, isimported: bool | None = None, gpuScaleMode: str | None = None, cnpq: str | None = None, ufmUrl: str | None = None, ufmUsername: str | None = None, ufmPasswordEncrypted: str | None = None) -> str:
    path = "addFabricData"
    body: dict[str, Any] = {}
    if id is not None:
        body["id"] = id
    body["name"] = name
    if type is not None:
        body["type"] = type
    if status is not None:
        body["status"] = status
    if description is not None:
        body["description"] = description
    if orchestrationStatus is not None:
        body["orchestrationStatus"] = orchestrationStatus
    if numOfSus is not None:
        body["numOfSus"] = numOfSus
    if maxNumOfSus is not None:
        body["maxNumOfSus"] = maxNumOfSus
    if dedicated is not None:
        body["dedicated"] = dedicated
    if hybrid is not None:
        body["hybrid"] = hybrid
    if isDPUFabric is not None:
        body["isDPUFabric"] = isDPUFabric
    if startingSubnetGpu is not None:
        body["startingSubnetGpu"] = startingSubnetGpu
    if startingSubnetCpu is not None:
        body["startingSubnetCpu"] = startingSubnetCpu
    if startingSubnetTenants is not None:
        body["startingSubnetTenants"] = startingSubnetTenants
    if startingSubnetStorage is not None:
        body["startingSubnetStorage"] = startingSubnetStorage
    if simulationId is not None:
        body["simulationId"] = simulationId
    if intent is not None:
        body["intent"] = intent
    if ewTenantAware is not None:
        body["ewTenantAware"] = ewTenantAware
    if storageTenantAware is not None:
        body["storageTenantAware"] = storageTenantAware
    if isOnesControlled is not None:
        body["isOnesControlled"] = isOnesControlled
    if suHostCnt is not None:
        body["suHostCnt"] = suHostCnt
    if isVxlanFabric is not None:
        body["isVxlanFabric"] = isVxlanFabric
    if nodeType is not None:
        body["nodeType"] = nodeType
    if deploymentType is not None:
        body["deploymentType"] = deploymentType
    if spineEvpnConfigured is not None:
        body["spineEvpnConfigured"] = spineEvpnConfigured
    if isimported is not None:
        body["isimported"] = isimported
    if gpuScaleMode is not None:
        body["gpuScaleMode"] = gpuScaleMode
    if cnpq is not None:
        body["cnpq"] = cnpq
    if ufmUrl is not None:
        body["ufmUrl"] = ufmUrl
    if ufmUsername is not None:
        body["ufmUsername"] = ufmUsername
    if ufmPasswordEncrypted is not None:
        body["ufmPasswordEncrypted"] = ufmPasswordEncrypted
    return client.call_api("POST", path, json_body=body or None)
