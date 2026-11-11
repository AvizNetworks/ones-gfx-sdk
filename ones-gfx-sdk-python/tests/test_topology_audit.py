"""Unit tests for Topology Audit (GFX-6219/GFX-6225): trigger + status/result
retrieval over the existing operation-poll endpoint.

No network: the client's requests session is replaced with a stub that
records each request and replies with a canned response.
"""
from __future__ import annotations

import json

import pytest
import requests

from ones_gfx import (
    Client,
    NotAuthenticatedError,
    OperationFailedError,
    OperationNotDoneError,
    ONESError,
    get_topology_audit_result,
    get_topology_audit_status,
    trigger_topology_audit,
)

BASE_URL = "https://ones.example:8089"
FABRIC = "fbrc-01"
OPERATION_ID = "TOPOAUDIT_fbrc-01_abc123"


def _response(status: int, body: object) -> requests.Response:
    resp = requests.Response()
    resp.status_code = status
    resp._content = json.dumps(body).encode() if body is not None else b""
    resp.headers["Content-Type"] = "application/json"
    return resp


class StubSession:
    def __init__(self, response: requests.Response):
        self.response = response
        self.calls: list[dict] = []

    def request(self, **kwargs):
        self.calls.append(kwargs)
        return self.response

    def close(self) -> None:
        pass


@pytest.fixture
def client():
    return Client(baseUrl=BASE_URL, authToken="tok")


def _stub(client: Client, status: int, body: object) -> StubSession:
    session = StubSession(_response(status, body))
    client._session = session
    return session


def _report(overall_status: str, links: list[dict], device_ip_map: dict | None = None) -> dict:
    return {
        "fabricName": FABRIC,
        "intentVersion": "42:2026-11-10T00:00:00Z",
        "overallStatus": overall_status,
        "links": links,
        "deviceIpMap": device_ip_map or {},
    }


# ---------------------------------------------------------------------------
# trigger_topology_audit
# ---------------------------------------------------------------------------


def test_trigger_topology_audit_sends_request_and_returns_envelope(client):
    accepted = {"operationId": OPERATION_ID, "status": "PENDING"}
    session = _stub(client, 202, {"status": "success", "message": "Topology audit accepted.", "data": accepted})

    result = trigger_topology_audit(client, FABRIC)

    assert result == accepted
    assert len(session.calls) == 1
    call = session.calls[0]
    assert call["method"] == "POST"
    assert call["url"] == f"{BASE_URL}/api/fm/fabrics/{FABRIC}/topology-audit"
    assert call["headers"]["authorization"] == "tok"


def test_trigger_topology_audit_invalid_fabric_propagates_404(client):
    body = {"status": "error", "message": "No fabric or intent found for 'fbrc-01'."}
    _stub(client, 404, body)

    with pytest.raises(requests.HTTPError) as excinfo:
        trigger_topology_audit(client, FABRIC)
    assert excinfo.value.response.status_code == 404


def test_trigger_topology_audit_requires_authentication():
    unauthenticated = Client(baseUrl=BASE_URL)
    session = _stub(unauthenticated, 202, {"operationId": OPERATION_ID, "status": "PENDING"})

    with pytest.raises(NotAuthenticatedError):
        trigger_topology_audit(unauthenticated, FABRIC)
    assert session.calls == []


# ---------------------------------------------------------------------------
# get_topology_audit_status (alias of get_operation)
# ---------------------------------------------------------------------------


def test_get_topology_audit_status_in_progress(client):
    _stub(client, 200, {"id": OPERATION_ID, "status": "RUNNING", "progress": 10})

    status = get_topology_audit_status(client, OPERATION_ID)

    assert status["status"] == "RUNNING"


# ---------------------------------------------------------------------------
# get_topology_audit_result
# ---------------------------------------------------------------------------


@pytest.mark.parametrize("status", ["PENDING", "RUNNING"])
def test_get_topology_audit_result_raises_while_in_progress(client, status):
    _stub(client, 200, {"id": OPERATION_ID, "status": status})

    with pytest.raises(OperationNotDoneError) as excinfo:
        get_topology_audit_result(client, OPERATION_ID)
    assert excinfo.value.operation_id == OPERATION_ID


@pytest.mark.parametrize("status", ["FAILED", "FAILURE"])
def test_get_topology_audit_result_raises_on_failure(client, status):
    _stub(client, 200, {
        "id": OPERATION_ID,
        "status": status,
        "errorMessage": "No intent found for fabric 'fbrc-01'.",
    })

    with pytest.raises(OperationFailedError) as excinfo:
        get_topology_audit_result(client, OPERATION_ID)
    assert "No intent found" in str(excinfo.value)


def test_get_topology_audit_result_raises_on_unknown_operation_id(client):
    _stub(client, 404, {"status": "error", "message": "Operation not found."})

    with pytest.raises(requests.HTTPError) as excinfo:
        get_topology_audit_result(client, "does-not-exist")
    assert excinfo.value.response.status_code == 404


def test_get_topology_audit_result_pass(client):
    report = _report("PASS", [
        {"device": "spine1", "localInterface": "Ethernet1", "remoteDevice": "leaf1",
         "remoteInterface": "Ethernet1", "observedRemoteDevice": "leaf1",
         "observedRemoteInterface": "Ethernet1", "status": "PASS", "message": "OK"},
    ])
    _stub(client, 200, {"id": OPERATION_ID, "status": "SUCCESS", "result": json.dumps(report)})

    result = get_topology_audit_result(client, OPERATION_ID)

    assert result == report
    assert result["overallStatus"] == "PASS"


def test_get_topology_audit_result_fail_rollup(client):
    report = _report("FAIL", [
        {"device": "spine1", "localInterface": "Ethernet3", "remoteDevice": "leaf2",
         "remoteInterface": "Ethernet1", "observedRemoteDevice": None,
         "observedRemoteInterface": None, "status": "MISSING_LINK",
         "message": "No LLDP neighbor observed on spine1:Ethernet3; expected leaf2:Ethernet1"},
    ])
    _stub(client, 200, {"id": OPERATION_ID, "status": "SUCCESS", "result": json.dumps(report)})

    result = get_topology_audit_result(client, OPERATION_ID)

    assert result["overallStatus"] == "FAIL"
    assert result["links"][0]["status"] == "MISSING_LINK"


def test_get_topology_audit_result_incomplete_rollup(client):
    report = _report("INCOMPLETE", [
        {"device": "spine1", "localInterface": "Ethernet3", "status": "UNVERIFIABLE",
         "message": "No LLDP neighbor observed and device unreachable"},
    ])
    _stub(client, 200, {"id": OPERATION_ID, "status": "SUCCESS", "result": json.dumps(report)})

    result = get_topology_audit_result(client, OPERATION_ID)

    assert result["overallStatus"] == "INCOMPLETE"


def test_get_topology_audit_result_malformed_result_raises(client):
    _stub(client, 200, {"id": OPERATION_ID, "status": "SUCCESS", "result": "{not valid json"})

    with pytest.raises(ONESError):
        get_topology_audit_result(client, OPERATION_ID)
