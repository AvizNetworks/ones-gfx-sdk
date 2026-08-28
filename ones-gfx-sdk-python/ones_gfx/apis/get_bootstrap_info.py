"""
GET /getbootstrapinfo  ->  /api/fm/getbootstrapinfo

All bootstrap device records.
"""
from __future__ import annotations

from typing import TypedDict


from .._types import BootstrapsubnetItem
from ..client import get_client

class BootstrapinfoRecord(TypedDict, total=False):
    """Models/Bootstrapinfo.java JSON (response)"""

    id: int
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
    lastupdated: str


def get_bootstrap_info() -> list[BootstrapinfoRecord]:
    client = get_client()
    path = "getbootstrapinfo"
    return client.call_api("GET", path)
