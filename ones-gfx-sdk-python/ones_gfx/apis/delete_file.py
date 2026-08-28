"""
DELETE /deletefile/{id}  ->  /api/fm/deletefile/{id}

Permanently delete a stored file.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class DeleteFileResult(TypedDict, total=False):
    """DELETE /deletefile/{id}"""

    success: bool
    message: str
    deletedId: int
    error: str


def delete_file(id: str) -> DeleteFileResult:
    client = get_client()
    path = f"deletefile/{id}"
    return client.call_api("DELETE", path)
