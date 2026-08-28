"""
GET /getControllerVersion  ->  /api/fm/getControllerVersion

Controller version info.
"""
from __future__ import annotations


from .._types import ControllerVersion
from ..client import Client


def get_controller_version(client: Client) -> ControllerVersion:
    path = "getControllerVersion"
    return client.call_api("GET", path)
