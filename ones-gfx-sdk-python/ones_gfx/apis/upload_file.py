"""
POST /uploadfile  ->  /api/fm/uploadfile

Store an artifact.
"""
from __future__ import annotations

from typing import Any

from .._types import UploadFileResult
from ..client import Client


def upload_file(client: Client, *, filetype: str, file: str, version: str | None = None, vendor: str | None = None, tag: str | None = None) -> UploadFileResult:
    path = "uploadfile"
    params: dict[str, Any] = {}
    params["filetype"] = filetype
    if version is not None:
        params["version"] = version
    if vendor is not None:
        params["vendor"] = vendor
    if tag is not None:
        params["tag"] = tag
    with open(file, "rb") as fh:
        files: dict[str, Any] = {"file": fh}
        return client.call_api("POST", path, params=params or None, files=files or None)
