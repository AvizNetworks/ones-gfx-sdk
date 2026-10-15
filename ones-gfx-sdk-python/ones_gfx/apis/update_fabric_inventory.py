"""
POST /api/config/updateinventory  (ONES UI config service)

Push device credentials into the Fabric Manager inventory.
"""
from __future__ import annotations

from typing import Any

from .._types import FabricDevice
from ..client import Client
from ._fabric_config import apply_config, yes_no


def update_fabric_inventory(client: Client, fabricName: str, devices: list[FabricDevice]) -> Any:
    body = [
        {
            "fabricName": fabricName,
            "hostname": d["hostname"],
            "ipAddress": d["ip"],
            "username": d["username"],
            "password": d["password"],
            "configure": yes_no(apply_config(d)),
        }
        for d in devices
    ]
    return client.call_config_api("POST", "api/config/updateinventory", json_body=body)
