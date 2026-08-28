"""
POST /fillbootstrapconfig  ->  /api/fm/fillbootstrapconfig

Stage a bootstrap batch.
"""
from __future__ import annotations

from typing import Any, TypedDict

from .._types import BootstrapsubnetItem
from ..client import get_client

class BootstrapinfoItem(TypedDict):
    deviceMacAddress: str
    serial: str
    deviceIp: str
    hostname: str
    region: str
    layer: str
    azid: str
    rackid: str
    brickid: str
    groupid: str
    paramspath: str
    bootfilepath: str
    nosImagePath: str
    agentImagePath: str
    fmcliImagePath: str
    baseConfigDbPath: str
    baseConfigFmPath: str
    status: int
    batchName: str
    collectorIP: str
    fabricid: str
    vendor: str
    bootstrapsubnet: BootstrapsubnetItem


def fill_bootstrap_config(*, batchName: str | None = None, subnet: str | None = None, netmask: str | None = None, gateway: str | None = None, bootstrapinfo: list[BootstrapinfoItem] | None = None) -> bool:
    client = get_client()
    path = "fillbootstrapconfig"
    body: dict[str, Any] = {}
    if batchName is not None:
        body["batchName"] = batchName
    if subnet is not None:
        body["subnet"] = subnet
    if netmask is not None:
        body["netmask"] = netmask
    if gateway is not None:
        body["gateway"] = gateway
    if bootstrapinfo is not None:
        body["bootstrapinfo"] = bootstrapinfo
    return client.call_api("POST", path, json_body=body or None)
