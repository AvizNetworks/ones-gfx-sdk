"""
GET /operations/{operationId}  ->  /api/fm/operations/{operationId}

Poll a long-running operation.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import Client

class OperationStatusItem(TypedDict, total=False):
    """Models/OperationStatus.java — GET /operations/{operationId}"""

    id: str
    type: str
    status: str
    createdAt: str
    updatedAt: str
    completedAt: str
    progress: int
    errorMessage: str
    result: str
    idempotencyKey: str
    cachedResponse: str
    httpStatusCode: int
    tenantName: str
    fabricName: str


def get_operation(client: Client, operationId: str) -> OperationStatusItem:
    path = f"operations/{operationId}"
    return client.call_api("GET", path)
