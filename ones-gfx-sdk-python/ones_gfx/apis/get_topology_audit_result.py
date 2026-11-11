"""
GET /operations/{operationId}  ->  /api/fm/operations/{operationId}

Status/result retrieval for a Topology Audit (GFX-6219/GFX-6225). Topology
Audit has no endpoints of its own beyond the trigger POST - both of these
ride the existing generic operation-poll endpoint.
"""
from __future__ import annotations

import json
from typing import List, Optional, TypedDict

from ..enums import OperationStatus
from ..exceptions import ONESError, OperationFailedError, OperationNotDoneError
from .get_operation import OperationStatusItem, get_operation
from ..client import Client


class TopologyAuditLink(TypedDict, total=False):
    """SES_FM TopologyCheckResult.LinkCheckResult.

    `status` is one of PASS (expected and observed agree), MISSING_LINK (LLDP
    collected fine but the expected neighbor wasn't seen), MISMATCH (a
    different device and/or interface was observed), UNVERIFIABLE (LLDP
    couldn't be collected/trusted at all), UNEXPECTED_LINK (LLDP shows a
    neighbor with no corresponding intent entry - remoteDevice/remoteInterface
    are None in this case), or SKIPPED (host/DPU/out-of-intent peer, not
    auditable by name). observedRemoteDevice/observedRemoteInterface are only
    populated once a neighbor was seen at all.
    """

    device: str
    localInterface: str
    remoteDevice: Optional[str]
    remoteInterface: Optional[str]
    observedRemoteDevice: Optional[str]
    observedRemoteInterface: Optional[str]
    status: str
    message: str


class TopologyAuditReport(TypedDict, total=False):
    """Parsed `result` payload of a completed (SUCCESS) topology-audit
    operation (SES_FM TopologyAuditService.buildReport).

    `overallStatus` is INCOMPLETE if any link is UNVERIFIABLE (a collection
    failure must never be reported as a clean PASS), else FAIL if any link is
    MISSING_LINK/MISMATCH/UNEXPECTED_LINK, else PASS. `deviceIpMap` maps
    hostname -> management IP for devices the audit observed.
    """

    fabricName: str
    intentVersion: str
    overallStatus: str
    links: List[TopologyAuditLink]
    deviceIpMap: dict[str, str]


def get_topology_audit_status(client: Client, operationId: str) -> OperationStatusItem:
    """Poll the audit's operation state.

    This is a named alias for get_operation() - Topology Audit has no status
    endpoint of its own.
    """
    return get_operation(client, operationId)


def get_topology_audit_result(client: Client, operationId: str) -> TopologyAuditReport:
    """Fetch the operation and parse its completed report.

    Raises OperationNotDoneError while the operation is still PENDING or
    RUNNING - call again later. Raises OperationFailedError if the audit
    itself failed (e.g. an unresolvable fabric/intent, or an exception during
    LLDP collection) - there is no report to parse in that case.
    """
    op = get_operation(client, operationId)
    status = op.get("status")
    if OperationStatus.is_failure(status):
        raise OperationFailedError(operationId, op.get("errorMessage"), op)
    if not OperationStatus.is_terminal(status):
        raise OperationNotDoneError(operationId, status)
    raw_result = op.get("result")
    if isinstance(raw_result, str):
        try:
            return json.loads(raw_result)
        except json.JSONDecodeError as e:
            raise ONESError(
                f"Operation {operationId}: malformed topology audit result: {e}"
            ) from e
    return raw_result or {}
