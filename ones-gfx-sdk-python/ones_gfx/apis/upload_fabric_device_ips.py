"""
POST /api/config/uploadip  (ONES UI config service)

Record each device's real IP and credentials against a fabric.
"""
from __future__ import annotations

from typing import Any

from .._types import FabricDevice
from ..client import Client
from ._fabric_config import apply_config, yes_no


def upload_fabric_device_ips(client: Client, fabricName: str, devices: list[FabricDevice]) -> Any:
    entries = []
    for d in devices:
        yn = yes_no(apply_config(d))
        entries.append({
            "hostname": d["hostname"],
            "ip": d["ip"],
            "username": d["username"],
            "password": d["password"],
            # The UI sends both casings; so do we.
            "applyConfig": yn,
            "applyconfig": yn,
            "configure": yn,
        })
    body = {"updatedDevices": entries, "fabricName": fabricName}
    return client.call_config_api("POST", "api/config/uploadip", json_body=body)
