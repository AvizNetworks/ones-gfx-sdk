"""
GET /getrmastatus  ->  /api/fm/getrmastatus

Per-stage RMA progress.
"""
from __future__ import annotations

from typing import Any, TypedDict

from ..client import Client

class RMAStatusItem(TypedDict, total=False):
    """Models/RMAStatus.java JSON"""

    id: int
    rmainfoId: int
    task: str
    status: int
    starttime: str
    endtime: str
    logs: str
    tasktitle: str


def get_rma_status(client: Client, *, rmaInfoId: int) -> list[RMAStatusItem]:
    path = "getrmastatus"
    params: dict[str, Any] = {}
    params["rmaInfoId"] = rmaInfoId
    return client.call_api("GET", path, params=params or None)
