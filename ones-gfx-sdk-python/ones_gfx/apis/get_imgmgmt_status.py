"""
POST /getImgmgmtStatus  ->  /api/fm/getImgmgmtStatus

Image-management progress per device.
"""
from __future__ import annotations


from ..client import get_client


def get_imgmgmt_status(*, payload: list[str]) -> list[object]:
    client = get_client()
    path = "getImgmgmtStatus"
    return client.call_api("POST", path, json_body=payload)
