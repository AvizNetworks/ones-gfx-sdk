"""Unit tests: request construction, validation and response handling of the
IPAM / rail-CIDRPool endpoints, with the HTTP session mocked out."""
import json
from unittest.mock import MagicMock

import pytest
import requests

import ones_gfx
from ones_gfx import Client, NotAuthenticatedError
from ones_gfx import apis
from ones_gfx.apis import (
    create_tenant,
    get_ipam_config,
    get_rail_cidr_pool,
    get_rail_cidr_pools,
    get_server_ips,
    get_tenant_compute_ip,
    get_tenant_ips,
    update_ipam_config,
    update_tenant,
)

NEW_APIS = [
    "get_rail_cidr_pools",
    "get_rail_cidr_pool",
    "get_ipam_config",
    "update_ipam_config",
    "get_tenant_compute_ip",
    "get_tenant_ips",
    "get_server_ips",
]


def _response(status=200, body=None, text=None):
    r = requests.Response()
    r.status_code = status
    r._content = (json.dumps(body) if body is not None else (text or "")).encode()
    return r


@pytest.fixture
def client():
    c = Client("https://fm.example:8443/", authToken="tok")
    c._session = MagicMock()
    c._session.request.return_value = _response(body={"status": "success", "data": {"ok": 1}})
    return c


def sent(client):
    """kwargs of the single request the client made."""
    client._session.request.assert_called_once()
    return client._session.request.call_args.kwargs


# ── exports ──────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("name", NEW_APIS)
def test_exported_from_both_packages(name):
    assert callable(getattr(apis, name))
    assert getattr(ones_gfx, name) is getattr(apis, name)
    assert name in apis.__all__ and name in ones_gfx.__all__


def test_all_has_no_duplicates():
    assert len(apis.__all__) == len(set(apis.__all__)) == 104


# ── rails ────────────────────────────────────────────────────────────────────

def test_rail_cidr_pools_defaults_to_yaml_accept(client):
    client._session.request.return_value = _response(text="apiVersion: v1\nkind: CIDRPool\n")
    out = get_rail_cidr_pools(client, "yml")
    kw = sent(client)
    assert kw["method"] == "GET"
    assert kw["url"] == "https://fm.example:8443/api/fm/fabrics/yml/rails/cidrpools"
    assert kw["headers"]["Accept"] == "application/yaml"
    assert kw["headers"]["authorization"] == "tok"
    assert out == "apiVersion: v1\nkind: CIDRPool\n"  # raw manifest, not parsed


@pytest.mark.parametrize("accept", ["application/yaml", "application/x-yaml", "text/plain", "application/json"])
def test_rail_cidr_pools_accept_passthrough(client, accept):
    get_rail_cidr_pools(client, "yml", accept=accept)
    assert sent(client)["headers"]["Accept"] == accept


def test_rail_cidr_pools_json_mode_unwraps_envelope(client):
    data = {"fabricName": "yml", "rails": ["r0", "r2"], "yaml": "x"}
    client._session.request.return_value = _response(body={"status": "success", "data": data})
    assert get_rail_cidr_pools(client, "yml", accept="application/json") == data


@pytest.mark.parametrize("bad", ["application/xml", "*/*", "", None, "APPLICATION/YAML"])
def test_rail_accept_rejected_before_any_request(client, bad):
    with pytest.raises(ValueError, match="accept"):
        get_rail_cidr_pools(client, "yml", accept=bad)
    with pytest.raises(ValueError, match="accept"):
        get_rail_cidr_pool(client, "yml", "r0", accept=bad)
    client._session.request.assert_not_called()


def test_rail_cidr_pool_path(client):
    client._session.request.return_value = _response(text="kind: CIDRPool\n")
    assert get_rail_cidr_pool(client, "yml", "r0", accept="text/plain") == "kind: CIDRPool\n"
    kw = sent(client)
    assert kw["url"].endswith("/api/fm/fabrics/yml/rails/r0/cidrpool")
    assert kw["headers"]["Accept"] == "text/plain"


def test_rail_http_error_keeps_plain_text_body(client):
    client._session.request.return_value = _response(404, text="INVALID_RAIL: Rail 'r9' does not exist in fabric 'yml'.")
    with pytest.raises(requests.HTTPError) as ei:
        get_rail_cidr_pool(client, "yml", "r9")
    assert ei.value.response.status_code == 404
    assert ei.value.response.text.startswith("INVALID_RAIL: ")


# ── ipam config ──────────────────────────────────────────────────────────────

def test_get_ipam_config(client):
    cfg = {"fabricName": "ns", "enableIpam": True, "enableIpAssign": True, "updatedAt": "t", "updatedBy": "Aviz Admin"}
    client._session.request.return_value = _response(body={"status": "success", "data": cfg})
    assert get_ipam_config(client, "ns") == cfg
    kw = sent(client)
    assert (kw["method"], kw["url"]) == ("GET", "https://fm.example:8443/api/fm/fabrics/ns/ipam/config")


def test_update_ipam_config_body_and_no_force_by_default(client):
    update_ipam_config(client, "ns", enableIpam=True, enableIpAssign=True)
    kw = sent(client)
    assert kw["method"] == "PATCH"
    assert kw["url"].endswith("/api/fm/fabrics/ns/ipam/config")
    assert kw["json"] == {"enableIpam": True, "enableIpAssign": True}
    assert kw["params"] is None


@pytest.mark.parametrize("force,wire", [(True, "true"), (False, "false")])
def test_update_ipam_config_force_is_lowercase_query(client, force, wire):
    update_ipam_config(client, "ns", enableIpam=False, enableIpAssign=False, force=force)
    assert sent(client)["params"] == {"force": wire}


@pytest.mark.parametrize(
    "kwargs",
    [
        {"enableIpam": "true", "enableIpAssign": False},
        {"enableIpam": True, "enableIpAssign": 1},
        {"enableIpam": None, "enableIpAssign": False},
        {"enableIpam": True, "enableIpAssign": True, "force": "yes"},
    ],
)
def test_update_ipam_config_rejects_non_bool(client, kwargs):
    with pytest.raises(ValueError, match="boolean"):
        update_ipam_config(client, "ns", **kwargs)
    client._session.request.assert_not_called()


def test_update_ipam_config_assign_requires_ipam(client):
    with pytest.raises(ValueError, match="enableIpAssign requires enableIpam"):
        update_ipam_config(client, "ns", enableIpam=False, enableIpAssign=True)
    client._session.request.assert_not_called()


def test_update_ipam_config_requires_keyword_switches(client):
    with pytest.raises(TypeError):
        update_ipam_config(client, "ns", True, True)  # type: ignore[misc]


def test_update_ipam_config_409_exposes_error_code(client):
    err = {"status": "error", "message": "m", "error": {"code": "IPAM_HAS_ACTIVE_POOLS", "detail": "m"}}
    client._session.request.return_value = _response(409, body=err)
    with pytest.raises(requests.HTTPError) as ei:
        update_ipam_config(client, "ns", enableIpam=False, enableIpAssign=False)
    assert ei.value.response.status_code == 409
    assert ei.value.response.json()["error"]["code"] == "IPAM_HAS_ACTIVE_POOLS"


# ── ip allocations ───────────────────────────────────────────────────────────

def test_get_tenant_compute_ip_path(client):
    get_tenant_compute_ip(client, "ns", "gorax", "su00-node00")
    kw = sent(client)
    assert kw["url"].endswith("/api/fm/fabrics/ns/tenants/gorax/computes/su00-node00/ip-address")


def test_get_tenant_ips_no_params_by_default(client):
    get_tenant_ips(client, "ns", "gorax")
    kw = sent(client)
    assert kw["url"].endswith("/api/fm/fabrics/ns/tenants/gorax/ip-addresses")
    assert kw["params"] is None


def test_get_tenant_ips_params(client):
    get_tenant_ips(client, "ns", "gorax", page=0, size=10, status="ASSIGNED")
    assert sent(client)["params"] == {"page": 0, "size": 10, "status": "ASSIGNED"}


@pytest.mark.parametrize(
    "kwargs",
    [
        {"page": -1},
        {"page": "0"},
        {"page": True},
        {"size": 0},
        {"size": -5},
        {"size": 1.5},
        {"status": "assigned"},
        {"status": "BOGUS"},
    ],
)
def test_get_tenant_ips_validation(client, kwargs):
    with pytest.raises(ValueError):
        get_tenant_ips(client, "ns", "gorax", **kwargs)
    client._session.request.assert_not_called()


def test_get_server_ips_path(client):
    get_server_ips(client, "ns", "su00-node00")
    assert sent(client)["url"].endswith("/api/fm/fabrics/ns/computes/su00-node00/ip-addresses")


# ── path-segment safety (applies to every new endpoint) ──────────────────────

CALLS = [
    lambda c, v: get_rail_cidr_pools(c, v),
    lambda c, v: get_rail_cidr_pool(c, v, "r0"),
    lambda c, v: get_rail_cidr_pool(c, "f", v),
    lambda c, v: get_ipam_config(c, v),
    lambda c, v: update_ipam_config(c, v, enableIpam=True, enableIpAssign=False),
    lambda c, v: get_tenant_compute_ip(c, v, "t", "s"),
    lambda c, v: get_tenant_compute_ip(c, "f", v, "s"),
    lambda c, v: get_tenant_compute_ip(c, "f", "t", v),
    lambda c, v: get_tenant_ips(c, v, "t"),
    lambda c, v: get_tenant_ips(c, "f", v),
    lambda c, v: get_server_ips(c, v, "s"),
    lambda c, v: get_server_ips(c, "f", v),
]


@pytest.mark.parametrize("call", CALLS)
@pytest.mark.parametrize("blank", ["", "   ", None, 7])
def test_blank_path_params_rejected(client, call, blank):
    with pytest.raises(ValueError, match="required"):
        call(client, blank)
    client._session.request.assert_not_called()


@pytest.mark.parametrize("call", CALLS[:1] + CALLS[3:4] + CALLS[8:9])
def test_path_params_are_quoted(client, call):
    call(client, "a/b?c#d")
    url = sent(client)["url"]
    assert "a%2Fb%3Fc%23d" in url
    assert "?" not in url and "#" not in url


@pytest.mark.parametrize("call", CALLS)
def test_unauthenticated_client_raises(call):
    c = Client("https://fm.example")
    with pytest.raises(NotAuthenticatedError):
        call(c, "x")


# ── modified existing endpoints ──────────────────────────────────────────────

def test_create_tenant_sends_dhcp_pool_only_when_given(client):
    create_tenant(client, "ns", tenantName="gorax")
    assert "tenantDhcpCidrPool" not in sent(client)["json"]
    client._session.request.reset_mock()
    create_tenant(client, "ns", tenantName="gorax", tenantDhcpCidrPool="192.3.0.0/17")
    assert sent(client)["json"] == {"tenantName": "gorax", "tenantDhcpCidrPool": "192.3.0.0/17"}


def test_update_tenant_returns_ipam_result_from_207(client):
    data = {
        "fabricName": "ns", "tenantName": "gorax", "operation": "ADD", "partial": True,
        "servers": [{"serverName": "full-node", "status": "IP_POOL_EXHAUSTED", "ipAddress": None, "ipAssignStatus": "IP_POOL_EXHAUSTED"}],
    }
    client._session.request.return_value = _response(207, body={"status": "success", "message": "m", "data": data})
    out = update_tenant(client, "ns", "gorax", operation="ADD", servers=[{"serverName": "full-node"}])
    assert out == data  # 207 is success, not an exception
    assert out["partial"] is True


@pytest.mark.parametrize("cidr", ["192.3.0.0/17", "10.0.0.0/8", "172.16.0.0/30", "192.3.0.5/17", " 192.3.0.0/24 "])
def test_create_tenant_accepts_server_valid_pools(client, cidr):
    create_tenant(client, "ns", tenantName="gorax", tenantDhcpCidrPool=cidr)
    assert sent(client)["json"]["tenantDhcpCidrPool"] == cidr.strip()


@pytest.mark.parametrize(
    "cidr,code",
    [
        ("", "INVALID_CIDR_FORMAT"),
        ("192.3.0.0", "INVALID_CIDR_FORMAT"),
        ("192.3.0.0/", "INVALID_CIDR_FORMAT"),
        ("192.3.0.0/123", "INVALID_CIDR_FORMAT"),
        ("192.3.0/24", "INVALID_CIDR_FORMAT"),
        ("192.3.0.256/24", "INVALID_CIDR_FORMAT"),
        ("a.b.c.d/24", "INVALID_CIDR_FORMAT"),
        ("fd00::/64", "INVALID_CIDR_FORMAT"),
        (7, "INVALID_CIDR_FORMAT"),
        ("192.3.0.0/7", "CIDR_POOL_TOO_SMALL"),
        ("192.3.0.0/31", "CIDR_POOL_TOO_SMALL"),
        ("192.3.0.0/32", "CIDR_POOL_TOO_SMALL"),
    ],
)
def test_create_tenant_rejects_invalid_pool_before_any_request(client, cidr, code):
    with pytest.raises(ValueError, match=code):
        create_tenant(client, "ns", tenantName="gorax", tenantDhcpCidrPool=cidr)
    client._session.request.assert_not_called()
