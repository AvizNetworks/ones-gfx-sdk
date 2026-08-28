"""
POST /getImgmgmtStatus  ->  /api/fm/getImgmgmtStatus

Image-management progress per device.
"""
from __future__ import annotations


from ..client import Client


def get_imgmgmt_status(client: Client, *, payload: list[str]) -> list[object]:
    path = "getImgmgmtStatus"
    return client.call_api("POST", path, json_body=payload)
