"""
POST /uploadUIObject  ->  /api/fm/uploadUIObject

Store a UI-built intent and trigger orchestration.
"""
from __future__ import annotations

from typing import Any

from .._types import DeviceItem, ParametersItem, QosItem, QosmappingItem, SchedulerItem
from ..client import Client


def upload_ui_object(client: Client, *, x_request_origin: str | None = None, id: int | None = None, name: str | None = None, fabricId: int | None = None, orchestrationMode: str | None = None, sspineCount: int | None = None, spineCount: int | None = None, leafCount: int | None = None, fec: str | None = None, mtu: str | None = None, adminStatus: str | None = None, asnSSpine: str | None = None, isAsnSSpineUnique: bool | None = None, asnLeaf: str | None = None, isAsnLeafUnique: bool | None = None, asnSpine: str | None = None, isAsnSpineUnique: bool | None = None, ntpServer: str | None = None, timezone: str | None = None, sysLogServer: str | None = None, snmpServer: str | None = None, isBGP_U: bool | None = None, isTorPresent: bool | None = None, ND_RA: int | None = None, isactive: bool | None = None, logs: str | None = None, issag: bool | None = None, dhcpServerIps: str | None = None, dhcpRelaySrcInterface: str | None = None, dhcpRelaySrcInterfaceIP: str | None = None, devices: list[DeviceItem] | None = None, qosmappings: list[QosmappingItem] | None = None, schedulers: list[SchedulerItem] | None = None, QoS: QosItem | None = None, parameters: ParametersItem | None = None) -> str:
    path = "uploadUIObject"
    headers: dict[str, str] = {}
    if x_request_origin is not None:
        headers["x-request-origin"] = x_request_origin
    body: dict[str, Any] = {}
    if id is not None:
        body["id"] = id
    if name is not None:
        body["name"] = name
    if fabricId is not None:
        body["fabricId"] = fabricId
    if orchestrationMode is not None:
        body["orchestrationMode"] = orchestrationMode
    if sspineCount is not None:
        body["sspineCount"] = sspineCount
    if spineCount is not None:
        body["spineCount"] = spineCount
    if leafCount is not None:
        body["leafCount"] = leafCount
    if fec is not None:
        body["fec"] = fec
    if mtu is not None:
        body["mtu"] = mtu
    if adminStatus is not None:
        body["adminStatus"] = adminStatus
    if asnSSpine is not None:
        body["asnSSpine"] = asnSSpine
    if isAsnSSpineUnique is not None:
        body["isAsnSSpineUnique"] = isAsnSSpineUnique
    if asnLeaf is not None:
        body["asnLeaf"] = asnLeaf
    if isAsnLeafUnique is not None:
        body["isAsnLeafUnique"] = isAsnLeafUnique
    if asnSpine is not None:
        body["asnSpine"] = asnSpine
    if isAsnSpineUnique is not None:
        body["isAsnSpineUnique"] = isAsnSpineUnique
    if ntpServer is not None:
        body["ntpServer"] = ntpServer
    if timezone is not None:
        body["timezone"] = timezone
    if sysLogServer is not None:
        body["sysLogServer"] = sysLogServer
    if snmpServer is not None:
        body["snmpServer"] = snmpServer
    if isBGP_U is not None:
        body["isBGP_U"] = isBGP_U
    if isTorPresent is not None:
        body["isTorPresent"] = isTorPresent
    if ND_RA is not None:
        body["ND_RA"] = ND_RA
    if isactive is not None:
        body["isactive"] = isactive
    if logs is not None:
        body["logs"] = logs
    if issag is not None:
        body["issag"] = issag
    if dhcpServerIps is not None:
        body["dhcpServerIps"] = dhcpServerIps
    if dhcpRelaySrcInterface is not None:
        body["dhcpRelaySrcInterface"] = dhcpRelaySrcInterface
    if dhcpRelaySrcInterfaceIP is not None:
        body["dhcpRelaySrcInterfaceIP"] = dhcpRelaySrcInterfaceIP
    if devices is not None:
        body["devices"] = devices
    if qosmappings is not None:
        body["qosmappings"] = qosmappings
    if schedulers is not None:
        body["schedulers"] = schedulers
    if QoS is not None:
        body["QoS"] = QoS
    if parameters is not None:
        body["parameters"] = parameters
    return client.call_api("POST", path, json_body=body or None, extra_headers=headers or None)
