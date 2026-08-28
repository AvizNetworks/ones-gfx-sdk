"""
DELETE /deletefile/{id}  ->  /api/fm/deletefile/{id}

Permanently delete a stored file.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class DeleteFileResult(TypedDict, total=False):
    """DELETE /deletefile/{id}"""

    success: bool
    message: str
    deletedId: int
    error: str


def delete_file(client: Client, id: str) -> DeleteFileResult:
    path = f"deletefile/{id}"
    return client.call_api("DELETE", path)
