"""
DELETE /fabrics/{fabricName}/vpcpeering  ->  /api/fm/fabrics/{fabricName}/vpcpeering

Remove route leaking between two VPCs (the inverse of create_vpc_peering).

Synchronous only: FM has no async path for this route. The peering is
addressed by `name`, and FM rejects the call unless `vpcname`/`peervpcname`
match the pair recorded for it. Failures raise `requests.HTTPError`; the
response body carries a machine-readable `error.code`:

- 404 `VPC_PEERING_NOT_FOUND` / `FABRIC_NOT_FOUND`
- 400 `VPC_PEERING_PAIR_MISMATCH`, `VPC_PEERING_NAME_REQUIRED`, ...
- 500 `VPC_PEERING_DELETE_FAILED` when a switch could not be updated
"""
from __future__ import annotations

from .._types import ApiResponseMessage
from ..client import Client


def delete_vpc_peering(client: Client, fabricName: str, *, name: str, vpcname: str, peervpcname: str) -> ApiResponseMessage:
    """Delete peering `name` between `vpcname` and `peervpcname` on `fabricName`."""
    body = {"name": name, "vpcname": vpcname, "peervpcname": peervpcname}
    return client.call_api("DELETE", f"fabrics/{fabricName}/vpcpeering", json_body=body)
