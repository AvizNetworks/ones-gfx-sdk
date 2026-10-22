"""
DELETE /fabrics/{fabricName}/networkpolicy/{policyName}  ->  /api/fm/fabrics/{fabricName}/networkpolicy/{policyName}

Delete a network policy and remove its config from the devices it was
pushed to.
"""
from __future__ import annotations

from typing import Any

from ..client import Client


def delete_network_policy(client: Client, fabricName: str, policyName: str) -> dict[str, Any]:
    """Delete `policyName` from `fabricName`. Returns `{"policyName": ...}`."""
    path = f"fabrics/{fabricName}/networkpolicy/{policyName}"
    return client.call_api("DELETE", path)
