"""
GET /getControllerVersionInternal  ->  /api/fm/getControllerVersionInternal

Controller version info (internal callers).
"""
from __future__ import annotations


from .._types import ControllerVersion
from ..client import get_client


def get_controller_version_internal() -> ControllerVersion:
    client = get_client()
    path = "getControllerVersionInternal"
    return client.call_api("GET", path)
