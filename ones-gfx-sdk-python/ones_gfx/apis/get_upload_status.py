"""
GET /uploadStatus  ->  /api/fm/uploadStatus

Returns the literal string 'uploadStatus'.
"""
from __future__ import annotations


from ..client import get_client


def get_upload_status() -> str:
    client = get_client()
    path = "uploadStatus"
    return client.call_api("GET", path)
