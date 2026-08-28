"""
GET /getrmainfo  ->  /api/fm/getrmainfo

All RMA-marked devices.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class RMAInfoRecord(TypedDict, total=False):
    """Models/RMAInfo.java JSON (response)"""

    id: int
    replacingDeviceMac: str
    replacingDeviceIp: str
    replacingDeviceSerial: str
    replacetoDeviceMac: str
    replacetoDeviceIp: str
    replacetoDeviceSerial: str
    fabricId: str
    state: str
    markedTime: str
    ticketId: str
    backupSelected: str
    scheduledTime: str
    isScheduledForLater: bool
    triggertime: str


def get_rma_info(client: Client) -> list[RMAInfoRecord]:
    path = "getrmainfo"
    return client.call_api("GET", path)
