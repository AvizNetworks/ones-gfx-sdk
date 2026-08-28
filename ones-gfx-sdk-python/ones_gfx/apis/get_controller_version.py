"""
GET /getControllerVersion  ->  /api/fm/getControllerVersion

Controller version info.
"""
from __future__ import annotations


from .._types import ControllerVersion
from ..client import get_client


def get_controller_version() -> ControllerVersion:
    client = get_client()
    path = "getControllerVersion"
    return client.call_api("GET", path)
