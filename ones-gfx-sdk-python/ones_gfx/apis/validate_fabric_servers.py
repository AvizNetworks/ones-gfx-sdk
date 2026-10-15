"""
POST /api/config/validateserver  (ONES UI config service)

Log in to each host/DPU over SSH and check its OS. A server passed when its
result has "os"; "error" explains a failure.
"""
from __future__ import annotations

from .._types import FabricDevice, ValidateDeviceResult
from ..client import Client
from .validate_fabric_switches import validation_entries


def validate_fabric_servers(client: Client, fabricName: str, devices: list[FabricDevice]) -> list[ValidateDeviceResult]:
    return client.call_config_api("POST", "api/config/validateserver", json_body=validation_entries(fabricName, devices))
