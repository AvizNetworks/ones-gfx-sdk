"""Unit tests for VPC peering deletion (DELETE /fabrics/{fabricName}/vpcpeering).

No network: the client's requests session is replaced with a stub that
records each request and replies with a canned response.
"""
from __future__ import annotations

import json

import pytest
import requests

from ones_gfx import Client, NotAuthenticatedError, delete_vpc_peering
from ones_gfx.exceptions import (
    AuthenticationError,
    BadRequestError,
    NotFoundError,
    ServerError,
)
from ones_gfx.resources import PeeringResource
from ones_gfx.transport import Transport

BASE_URL = "https://ones.example:8089"
FABRIC = "Fab101"
NAME = "tenant1-storage-route-leak"
VPC = "tenant1-Fab101-north-south"
PEER = "Fab101-Storage"


def _response(status: int, body: object) -> requests.Response:
    resp = requests.Response()
    resp.status_code = status
    resp._content = json.dumps(body).encode() if body is not None else b""
    resp.headers["Content-Type"] = "application/json"
    return resp


def _peering_error(code: str, message: str) -> dict:
    """The envelope FM's peeringError() returns."""
    return {"status": "error", "message": message, "error": {"code": code, "detail": message}}


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


# ---------------------------------------------------------------------------
# Flat API: ones_gfx.apis.delete_vpc_peering
# ---------------------------------------------------------------------------


def test_delete_vpc_peering_success_sends_request_and_returns_envelope(client):
    ok = {
        "status": "success",
        "message": f"VPC peering '{NAME}' deleted successfully.",
    }
    session = _stub(client, 200, ok)

    result = delete_vpc_peering(client, FABRIC, name=NAME, vpcname=VPC, peervpcname=PEER)

    assert result == ok
    assert len(session.calls) == 1
    call = session.calls[0]
    assert call["method"] == "DELETE"
    assert call["url"] == f"{BASE_URL}/api/fm/fabrics/{FABRIC}/vpcpeering"
    assert call["json"] == {"name": NAME, "vpcname": VPC, "peervpcname": PEER}
    assert call["headers"]["authorization"] == "tok"


def test_delete_vpc_peering_empty_body_returns_none(client):
    _stub(client, 204, None)
    assert delete_vpc_peering(client, FABRIC, name=NAME, vpcname=VPC, peervpcname=PEER) is None


@pytest.mark.parametrize(
    "status, code",
    [
        (404, "VPC_PEERING_NOT_FOUND"),
        (404, "FABRIC_NOT_FOUND"),
        (400, "VPC_PEERING_PAIR_MISMATCH"),
        (400, "VPC_PEERING_NAME_REQUIRED"),
        (401, "UNAUTHORIZED"),
        (403, "FORBIDDEN"),
        (500, "VPC_PEERING_DELETE_FAILED"),
    ],
)
def test_delete_vpc_peering_propagates_api_errors(client, status, code):
    body = _peering_error(code, f"{code} message")
    _stub(client, status, body)

    with pytest.raises(requests.HTTPError) as excinfo:
        delete_vpc_peering(client, FABRIC, name=NAME, vpcname=VPC, peervpcname=PEER)

    resp = excinfo.value.response
    assert resp.status_code == status
    assert resp.json()["error"]["code"] == code


def test_delete_vpc_peering_requires_authentication():
    unauthenticated = Client(baseUrl=BASE_URL)
    session = _stub(unauthenticated, 200, {"status": "success", "message": "ok"})

    with pytest.raises(NotAuthenticatedError):
        delete_vpc_peering(unauthenticated, FABRIC, name=NAME, vpcname=VPC, peervpcname=PEER)
    assert session.calls == []


def test_delete_vpc_peering_arguments_are_keyword_only(client):
    _stub(client, 200, {"status": "success", "message": "ok"})
    with pytest.raises(TypeError):
        delete_vpc_peering(client, FABRIC, NAME, VPC, PEER)  # type: ignore[misc]


# ---------------------------------------------------------------------------
# Resource API: PeeringResource.delete (typed exceptions via Transport)
# ---------------------------------------------------------------------------


class _NoAuth:
    def apply(self, headers: dict) -> None:
        headers["Authorization"] = "Bearer tok"

    def needs_proactive_refresh(self) -> bool:
        return False

    def refresh(self) -> None:
        raise AuthenticationError("refresh failed", status_code=401)

    def close(self) -> None:
        pass


@pytest.fixture
def peering():
    transport = Transport(f"{BASE_URL}/api/fm", _NoAuth())
    return PeeringResource(transport), transport


def test_peering_resource_delete_success(peering):
    resource, transport = peering
    session = StubSession(_response(200, {"status": "success", "message": "deleted"}))
    transport._session = session

    assert resource.delete(FABRIC, NAME, VPC, PEER) == "deleted"
    call = session.calls[0]
    assert call["method"] == "DELETE"
    assert call["url"] == f"{BASE_URL}/api/fm/fabrics/{FABRIC}/vpcpeering"
    assert call["json"] == {"name": NAME, "vpcname": VPC, "peervpcname": PEER}
    assert call["headers"]["Prefer"] == "respond-sync"


@pytest.mark.parametrize(
    "status, code, exc_cls",
    [
        (400, "VPC_PEERING_PAIR_MISMATCH", BadRequestError),
        (403, "FORBIDDEN", AuthenticationError),
        (404, "VPC_PEERING_NOT_FOUND", NotFoundError),
        (500, "VPC_PEERING_DELETE_FAILED", ServerError),
    ],
)
def test_peering_resource_delete_maps_errors(peering, status, code, exc_cls):
    resource, transport = peering
    body = _peering_error(code, f"{code} message")
    transport._session = StubSession(_response(status, body))

    with pytest.raises(exc_cls) as excinfo:
        resource.delete(FABRIC, NAME, VPC, PEER)

    assert excinfo.value.status_code == status
    assert excinfo.value.response_body["error"]["code"] == code


def test_peering_resource_delete_401_refresh_failure_raises_auth_error(peering):
    resource, transport = peering
    transport._session = StubSession(_response(401, _peering_error("UNAUTHORIZED", "expired")))

    with pytest.raises(AuthenticationError):
        resource.delete(FABRIC, NAME, VPC, PEER)


@pytest.mark.parametrize(
    "args, missing",
    [
        (("", NAME, VPC, PEER), "fabric_name"),
        ((FABRIC, "", VPC, PEER), "name"),
        ((FABRIC, NAME, "", PEER), "vpc_name"),
        ((FABRIC, NAME, VPC, ""), "peer_vpc_name"),
    ],
)
def test_peering_resource_delete_rejects_missing_identifiers(peering, args, missing):
    resource, transport = peering
    session = StubSession(_response(200, {"status": "success", "message": "x"}))
    transport._session = session

    with pytest.raises(ValueError, match=f"^{missing} is required"):
        resource.delete(*args)
    assert session.calls == []
