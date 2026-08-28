"""
POST /uploadfile  ->  /api/fm/uploadfile

Store an artifact.
"""
from __future__ import annotations

from typing import Any

from .._types import UploadFileResult
from ..client import get_client


def upload_file(*, filetype: str, file: str, version: str | None = None, vendor: str | None = None, tag: str | None = None) -> UploadFileResult:
    client = get_client()
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
