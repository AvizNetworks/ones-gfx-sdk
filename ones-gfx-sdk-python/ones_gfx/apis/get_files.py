"""
GET /getfiles/{filetype}  ->  /api/fm/getfiles/{filetype}

Files of a type plus a count.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class GetFilesResult(TypedDict, total=False):
    """GET /getfiles/{filetype} — files is a grouped map (ALL) or a list (typed)"""

    success: bool
    files: object
    totalCount: int
    count: int
    error: str


def get_files(client: Client, filetype: str) -> GetFilesResult:
    path = f"getfiles/{filetype}"
    return client.call_api("GET", path)
