"""Unit tests for the VPC peering read APIs (GFX-5658).

GET /fabrics/{f}/vpcpeering, GET /fabrics/{f}/vpcpeering/{name} and
GET /fabrics/{f}/tenants/{t}/vpcpeering. The gateway does not proxy these, so
they go straight to Fabric Manager; create and delete on the same path keep
going through the gateway.
"""
from __future__ import annotations

import json

import pytest
import requests

from ones_gfx import (
    Client,
    delete_vpc_peering,
    get_vpc_peering,
    is_direct_route,
    list_tenant_vpc_peerings,
    list_vpc_peerings,
)

BASE_URL = "https://ones.example:8089"
FM_URL = "https://ones.example:8443"  # derived FM-direct base


def _response(status: int, body: object) -> requests.Response:
    resp = requests.Response()
    resp.status_code = status
    resp._content = json.dumps(body).encode()
    resp.headers["Content-Type"] = "application/json"
    return resp


def _ok(data: dict) -> dict:
    return {"status": "success", "message": "VPC peerings retrieved successfully.", "data": data}


def _err(code: str, message: str) -> dict:
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


PAGE = {"page": 0, "size": 50, "totalElements": 1, "totalPages": 1}


def test_list_vpc_peerings_goes_direct_and_returns_page(client):
    data = {
        "fabricName": "Fab101",
        "peerings": [{
            "name": "t1-storage",
            "vpcs": ["t1-Fab101-north-south", "Fab101-Storage"],
            "origin": "VPC_PEERING",
            "createdAt": "2026-10-28T10:31:29Z",
        }],
        "pagination": PAGE,
    }
    session = _stub(client, 200, _ok(data))

    result = list_vpc_peerings(client, "Fab101")

    assert result == data
    call = session.calls[0]
    assert call["method"] == "GET"
    assert call["url"] == f"{FM_URL}/fabrics/Fab101/vpcpeering"
    assert call["params"] is None
    assert call["json"] is None
    assert call["headers"]["authorization"] == "tok"


@pytest.mark.parametrize(
    "kwargs, params",
    [({"page": 2, "size": 25}, {"page": 2, "size": 25}), ({"page": 0}, {"page": 0}), ({}, None)],
)
def test_list_vpc_peerings_paging(client, kwargs, params):
    session = _stub(client, 200, _ok({"fabricName": "F", "peerings": [], "pagination": PAGE}))
    list_vpc_peerings(client, "F", **kwargs)
    assert session.calls[0]["params"] == params


def test_get_vpc_peering_unwraps_peering(client):
    peering = {"name": "t1-storage", "vpcs": ["a", "b"], "origin": "NETWORK_POLICY", "createdAt": "x"}
    session = _stub(client, 200, _ok({"fabricName": "Fab101", "peering": peering}))

    assert get_vpc_peering(client, "Fab101", "t1-storage") == peering
    assert session.calls[0]["url"] == f"{FM_URL}/fabrics/Fab101/vpcpeering/t1-storage"


def test_list_tenant_vpc_peerings_tenant_shape(client):
    data = {
        "fabricName": "Fab101",
        "tenantName": "t1",
        "vpc": "t1-Fab101-north-south",
        "peerings": [{"name": "t1-storage", "vpc": "t1-Fab101-north-south", "peerVpc": "Fab101-Storage",
                      "origin": "VPC_PEERING", "createdAt": "x"}],
        "pagination": PAGE,
    }
    session = _stub(client, 200, _ok(data))

    result = list_tenant_vpc_peerings(client, "Fab101", "t1", page=0, size=10)

    assert result["peerings"][0]["peerVpc"] == "Fab101-Storage"
    call = session.calls[0]
    assert call["url"] == f"{FM_URL}/fabrics/Fab101/tenants/t1/vpcpeering"
    assert call["params"] == {"page": 0, "size": 10}


@pytest.mark.parametrize(
    "fn, args, code",
    [
        (list_vpc_peerings, ("Gone",), "FABRIC_NOT_FOUND"),
        (get_vpc_peering, ("Fab101", "nope"), "VPC_PEERING_NOT_FOUND"),
        (list_tenant_vpc_peerings, ("Fab101", "ghost"), "TENANT_NOT_FOUND"),
    ],
)
def test_reads_propagate_404(client, fn, args, code):
    _stub(client, 404, _err(code, f"{code} message"))
    with pytest.raises(requests.HTTPError) as excinfo:
        fn(client, *args)
    assert excinfo.value.response.status_code == 404
    assert excinfo.value.response.json()["error"]["code"] == code


@pytest.mark.parametrize("status", [401, 403, 500])
def test_reads_propagate_auth_and_server_errors(client, status):
    _stub(client, status, _err("ERROR", "x"))
    with pytest.raises(requests.HTTPError) as excinfo:
        list_vpc_peerings(client, "Fab101")
    assert excinfo.value.response.status_code == status


def test_delete_still_goes_through_gateway(client):
    session = _stub(client, 200, {"status": "success", "message": "deleted"})
    delete_vpc_peering(client, "F", name="p1", vpcname="t1-F-north-south", peervpcname="F-Storage")
    assert session.calls[0]["url"] == f"{BASE_URL}/api/fm/fabrics/F/vpcpeering"


@pytest.mark.parametrize(
    "method, path, direct",
    [
        ("GET", "fabrics/F/vpcpeering", True),
        ("get", "/fabrics/F/vpcpeering?page=1", True),
        ("GET", "fabrics/F/vpcpeering/p1", True),
        ("GET", "fabrics/F/tenants/t1/vpcpeering", True),
        ("POST", "fabrics/F/vpcpeering", False),
        ("DELETE", "fabrics/F/vpcpeering", False),
        ("GET", "fabrics/F/tenants/t1", False),
        ("POST", "fabrics/F/networkpolicy", True),  # path-only table applies to every method
    ],
)
def test_is_direct_route(method, path, direct):
    assert is_direct_route(method, path) is direct
