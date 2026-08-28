"""
POST /ValidateUfmCreds  ->  /api/fm/ValidateUfmCreds

Validate UFM URL/username/password.
"""
from __future__ import annotations

from typing import Any, TypedDict

from ..client import get_client

class UfmCredsResult(TypedDict, total=False):
    """POST /ValidateUfmCreds"""

    success: bool
    message: str
    error: str


def validate_ufm_creds(*, ufmUrl: str, username: str, password: str) -> UfmCredsResult:
    client = get_client()
    path = "ValidateUfmCreds"
    body: dict[str, Any] = {}
    body["ufmUrl"] = ufmUrl
    body["username"] = username
    body["password"] = password
    return client.call_api("POST", path, json_body=body or None)
