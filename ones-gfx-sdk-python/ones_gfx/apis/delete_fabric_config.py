"""
DELETE /api/config/deletefabricbyname/{name}  (ONES UI config service)

Delete a fabric. Already deleted (404) is not an error. A Deployed fabric has
no undeploy call: this removes the record, not the configuration already on
the switches. Its inventory rows are left behind, so recreating a fabric with
the same device IPs can produce duplicate inventory entries.
"""
from __future__ import annotations

from typing import Any
from urllib.parse import quote

import requests

from ..client import Client


def delete_fabric_config(client: Client, name: str) -> Any:
    try:
        return client.call_config_api("DELETE", f"api/config/deletefabricbyname/{quote(name, safe='')}")
    except requests.HTTPError as exc:
        if exc.response is not None and exc.response.status_code == 404:
            return None
        raise
