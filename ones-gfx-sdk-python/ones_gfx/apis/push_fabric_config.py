"""
POST /api/config  (ONES UI config service)

Send the fabric document (its YAML, parsed) to be configured on the devices.
Build ``data`` with fabric_yaml_to_document so the keys keep their YAML order.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def push_fabric_config(client: Client, fabricName: str, data: Any, *, instance: str = "fm") -> Any:
    body = {"data": data, "fabricName": fabricName, "instance": instance}
    return client.call_config_api("POST", "api/config", json_body=body)
