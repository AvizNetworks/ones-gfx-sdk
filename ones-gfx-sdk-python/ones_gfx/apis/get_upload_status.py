"""
GET /uploadStatus  ->  /api/fm/uploadStatus

Returns the literal string 'uploadStatus'.
"""
from __future__ import annotations


from ..client import Client


def get_upload_status(client: Client) -> str:
    path = "uploadStatus"
    return client.call_api("GET", path)
