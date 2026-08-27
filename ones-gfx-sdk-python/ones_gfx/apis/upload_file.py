"""
POST /uploadfile  ->  /api/fm/uploadfile

Store an artifact.
"""
from __future__ import annotations

from typing import Any

from ..client import ONESClient


def upload_file(client: ONESClient, *, file, filetype, version=None, vendor=None, tag=None) -> Any:
    path = "uploadfile"
    params: dict[str, Any] = {}
    params["filetype"] = filetype
    if version is not None:
        params["version"] = version
    if vendor is not None:
        params["vendor"] = vendor
    if tag is not None:
        params["tag"] = tag
    files: dict[str, Any] = {}
    files["file"] = open(file, "rb")
    return client.call_api("POST", path, params=params or None, files=files or None)
