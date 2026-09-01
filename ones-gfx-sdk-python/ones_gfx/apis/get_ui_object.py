"""
GET /getUIObject  ->  /api/fm/getUIObject

Fetch stored intent.
"""
from __future__ import annotations

from typing import Any, TypedDict

from .._types import DeviceItem, ParametersItem, QosItem, QosmappingItem, SchedulerItem
from ..client import Client

class IntentItem(TypedDict, total=False):
    """Models/Intent.java JSON (GET /getUIObject)"""

    id: int
    name: str
    fabricId: int
    orchestrationMode: str
    sspineCount: int
    spineCount: int
    leafCount: int
    fec: str
    mtu: str
    adminStatus: str
    asnSSpine: str
    isAsnSSpineUnique: bool
    asnLeaf: str
    isAsnLeafUnique: bool
    asnSpine: str
    isAsnSpineUnique: bool
    ntpServer: str
    timezone: str
    sysLogServer: str
    snmpServer: str
    isBGP_U: bool
    isTorPresent: bool
    ND_RA: int
    isactive: bool
    logs: str
    issag: bool
    dhcpServerIps: str
    dhcpRelaySrcInterface: str
    dhcpRelaySrcInterfaceIP: str
    devices: list[DeviceItem]
    qosmappings: list[QosmappingItem]
    schedulers: list[SchedulerItem]
    QoS: QosItem
    parameters: ParametersItem
    createdAt: str
    updatedAt: str


def get_ui_object(client: Client, *, name: str | None = None) -> IntentItem:
    path = "getUIObject"
    params: dict[str, Any] = {}
    if name is not None:
        params["name"] = name
    return client.call_api("GET", path, params=params or None)
