"""
POST /api/config/validateswitch  (ONES UI config service)

Log in to each switch over SSH and check it is reachable and running the
expected NOS. A switch passed when its result has "build"; "error" explains a
failure.
"""
from __future__ import annotations

from .._types import FabricDevice, ValidateDeviceResult
from ..client import Client
from ._fabric_config import apply_config, yes_no


def validation_entries(fabricName: str, devices: list[FabricDevice]) -> list[dict[str, str]]:
    entries = []
    for d in devices:
        entry = {"hostname": d["hostname"], "ip": d["ip"], "username": d["username"], "password": d["password"]}
        if d.get("device_type"):
            entry["deviceType"] = d["device_type"]
        if d.get("device_role"):
            entry["deviceRole"] = d["device_role"]
        entry["applyConfig"] = yes_no(apply_config(d))
        entry["fabricName"] = fabricName
        entries.append(entry)
    return entries


def validate_fabric_switches(client: Client, fabricName: str, devices: list[FabricDevice]) -> list[ValidateDeviceResult]:
    return client.call_config_api("POST", "api/config/validateswitch", json_body=validation_entries(fabricName, devices))
