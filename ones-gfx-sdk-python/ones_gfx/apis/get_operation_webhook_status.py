"""
GET /operations/{operationId}/webhook-status  ->  /api/fm/operations/{operationId}/webhook-status

Webhook delivery attempts for an operation.
"""
from __future__ import annotations

from typing import TypedDict


from ..client import get_client

class WebhookDeliveryStatus(TypedDict, total=False):
    """GET /operations/{operationId}/webhook-status"""

    operationId: str
    deliveryStatus: str
    totalAttempts: int
    lastStatusCode: int
    lastErrorMessage: str
    lastAttemptedAt: str
    nextRetryAt: str
    attempts: list[WebhookDeliveryAttempt]


class WebhookDeliveryAttempt(TypedDict, total=False):
    attemptNumber: int
    status: str
    statusCode: int
    errorMessage: str
    attemptedAt: str


def get_operation_webhook_status(operationId: str) -> WebhookDeliveryStatus:
    client = get_client()
    path = f"operations/{operationId}/webhook-status"
    return client.call_api("GET", path)
