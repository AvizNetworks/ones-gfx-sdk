"""
POST /api/config/updatefabricstatus  (ONES UI config service)

Set a fabric's status, e.g. "Deployed". deploymentType is DEFAULT,
PARTIAL_CONFIG or BROWNFIELD.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def update_fabric_config_status(client: Client, *, name: str, status: str, description: str = "", deploymentType: str = "DEFAULT") -> Any:
    body = {"name": name, "status": status, "description": description, "deploymentType": deploymentType}
    return client.call_config_api("POST", "api/config/updatefabricstatus", json_body=body)
