"""Integration tests: the real SDK + `requests` stack against a fake Fabric
Manager that replays the SES_FM controller behaviour over HTTP (see
fake_fm.py). Payloads of the happy paths are the ones captured from a live
instance."""
import json

import pytest
import requests

from ones_gfx import Client
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

from .fake_fm import ALL_RAILS_YAML, RAIL_YAML, TOKEN


def http_error(callable_, *a, **kw) -> requests.Response:
    with pytest.raises(requests.HTTPError) as ei:
        callable_(*a, **kw)
    return ei.value.response


def code_of(resp: requests.Response) -> str:
    return resp.json()["error"]["code"]


# ── rails ────────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("accept", ["application/yaml", "application/x-yaml", "text/plain"])
def test_all_rails_manifest_verbatim(live_client, accept):
    text = get_rail_cidr_pools(live_client, "yml", accept=accept)
    assert text == ALL_RAILS_YAML
    assert text.count("kind: CIDRPool") == 2


def test_default_accept_returns_yaml_text(live_client, fm):
    assert get_rail_cidr_pools(live_client, "yml") == ALL_RAILS_YAML
    assert fm.requests[-1]["headers"]["accept"] == "application/yaml"


def test_all_rails_json_envelope(live_client):
    data = get_rail_cidr_pools(live_client, "yml", accept="application/json")
    assert data == {"fabricName": "yml", "rails": ["r0", "r2"], "yaml": ALL_RAILS_YAML}


def test_single_rail_yaml_and_json(live_client):
    assert get_rail_cidr_pool(live_client, "yml", "r0", accept="text/plain") == RAIL_YAML["r0"]
    data = get_rail_cidr_pool(live_client, "yml", "r2", accept="application/json")
    assert data == {"fabricName": "yml", "railName": "r2", "yaml": RAIL_YAML["r2"]}


def test_rail_errors_are_plain_text_for_yaml_and_json_for_json(live_client):
    resp = http_error(get_rail_cidr_pool, live_client, "yml", "r9", accept="text/plain")
    assert resp.status_code == 404
    assert resp.text == "INVALID_RAIL: Rail 'r9' does not exist in fabric 'yml'."
    assert resp.headers["Content-Type"].startswith("text/plain")

    resp = http_error(get_rail_cidr_pool, live_client, "yml", "r9", accept="application/json")
    assert (resp.status_code, code_of(resp)) == (404, "INVALID_RAIL")


@pytest.mark.parametrize("accept", ["application/yaml", "application/json"])
def test_rail_unknown_fabric_and_no_rails(live_client, accept):
    assert http_error(get_rail_cidr_pools, live_client, "nope", accept=accept).status_code == 404
    assert http_error(get_rail_cidr_pools, live_client, "norails", accept=accept).status_code == 404
    resp = http_error(get_rail_cidr_pools, live_client, "norails", accept="application/json")
    assert code_of(resp) == "RAIL_CIDRPOOL_NOT_FOUND"


# ── ipam config ──────────────────────────────────────────────────────────────

def test_unconfigured_fabric_reports_disabled(live_client):
    cfg = get_ipam_config(live_client, "ns")
    assert cfg["enableIpam"] is False and cfg["enableIpAssign"] is False
    assert cfg["updatedBy"] is None


def test_update_then_get_roundtrip(live_client):
    out = update_ipam_config(live_client, "ns", enableIpam=True, enableIpAssign=True)
    assert out == {
        "fabricName": "ns", "enableIpam": True, "enableIpAssign": True,
        "updatedAt": "2026-09-24T08:57:28.254792109Z", "updatedBy": "Aviz Admin",
    }
    got = get_ipam_config(live_client, "ns")
    assert (got["enableIpam"], got["enableIpAssign"], got["updatedBy"]) == (True, True, "Aviz Admin")


def test_ipam_config_unknown_fabric(live_client):
    assert code_of(http_error(get_ipam_config, live_client, "nope")) == "FABRIC_NOT_FOUND"
    resp = http_error(update_ipam_config, live_client, "nope", enableIpam=True, enableIpAssign=False)
    assert (resp.status_code, code_of(resp)) == (404, "FABRIC_NOT_FOUND")


def test_ones_controlled_fabric_cannot_be_modified(live_client):
    resp = http_error(update_ipam_config, live_client, "ones-ctl", enableIpam=True, enableIpAssign=False)
    assert (resp.status_code, code_of(resp)) == (409, "FABRIC_NOT_EXTERNALLY_MANAGED")


def test_disabling_with_active_pools_needs_force(live_client, fm):
    update_ipam_config(live_client, "ns", enableIpam=True, enableIpAssign=True)
    resp = http_error(update_ipam_config, live_client, "ns", enableIpam=False, enableIpAssign=False)
    assert (resp.status_code, code_of(resp)) == (409, "IPAM_HAS_ACTIVE_POOLS")
    assert get_ipam_config(live_client, "ns")["enableIpam"] is True  # unchanged

    out = update_ipam_config(live_client, "ns", enableIpam=False, enableIpAssign=False, force=True)
    assert out["enableIpam"] is False
    assert fm.requests[-1]["query"] == {"force": ["true"]}
    assert not fm.pools  # force disabled the pool


def test_server_side_validation_matches_sdk_validation(live_client, fm):
    """The SDK refuses these locally; confirm the server really does reject them,
    so the local checks are not stricter than the contract."""
    for body in ('{"enableIpam": "yes", "enableIpAssign": false}', '{"enableIpam": true}'):
        r = requests.patch(
            f"{fm.base_url}/api/fm/fabrics/ns/ipam/config",
            data=body, headers={"authorization": TOKEN, "Content-Type": "application/json"},
        )
        assert (r.status_code, r.json()["error"]["code"]) == (400, "INVALID_REQUEST_BODY")
    r = requests.patch(
        f"{fm.base_url}/api/fm/fabrics/ns/ipam/config",
        json={"enableIpam": False, "enableIpAssign": True}, headers={"authorization": TOKEN},
    )
    assert (r.status_code, r.json()["error"]["code"]) == (400, "INVALID_IPAM_CONFIG")


# ── tenant GPU update + IP allocation lifecycle ──────────────────────────────

@pytest.fixture
def ipam_on(live_client):
    update_ipam_config(live_client, "ns", enableIpam=True, enableIpAssign=True)
    return live_client


def test_add_server_allocates_and_assigns_ip(ipam_on):
    out = update_tenant(ipam_on, "ns", "gorax", operation="ADD", servers=[{"serverName": "su00-node00"}])
    assert out == {
        "fabricName": "ns", "tenantName": "gorax", "operation": "ADD", "partial": False,
        "servers": [{"serverName": "su00-node00", "status": "SUCCESS", "ipAddress": "192.3.0.2/32", "ipAssignStatus": "ASSIGNED"}],
    }


def test_ipam_disabled_update_returns_plain_message(live_client):
    out = update_tenant(live_client, "ns", "gorax", operation="ADD", servers=[{"serverName": "su00-node00"}])
    assert out == {"status": "success", "message": "GPU update processed successfully."}


def test_partial_update_is_207_but_returns_normally(ipam_on):
    out = update_tenant(
        ipam_on, "ns", "gorax", operation="ADD",
        servers=[{"serverName": "su00-node01"}, {"serverName": "full-node"}],
    )
    assert out["partial"] is True
    by_name = {s["serverName"]: s for s in out["servers"]}
    assert by_name["su00-node01"]["status"] == "SUCCESS"
    assert by_name["full-node"]["status"] == "IP_POOL_EXHAUSTED"
    assert by_name["full-node"]["ipAddress"] is None and by_name["full-node"]["message"]


def test_read_endpoints_reflect_allocation(ipam_on):
    update_tenant(ipam_on, "ns", "gorax", operation="ADD", servers=[{"serverName": "su00-node00"}])

    one = get_tenant_compute_ip(ipam_on, "ns", "gorax", "su00-node00")
    assert one["ipAddress"] == "192.3.0.2" and one["prefixLength"] == 32
    assert one["assignStatus"] == "ASSIGNED" and one["cleanupStatus"] == "NONE"
    assert one["tenantName"] == "gorax" and one["serverName"] == "su00-node00"

    tenant = get_tenant_ips(ipam_on, "ns", "gorax")
    assert tenant["pool"] == {"cidr": "192.3.0.0/17", "allocatableAddressCount": 32765, "allocatedCount": 1}
    assert tenant["pagination"] == {"page": 0, "size": 50, "totalElements": 1, "totalPages": 1}
    assert tenant["allocations"] == [one]

    server = get_server_ips(ipam_on, "ns", "su00-node00")
    assert server == {"fabricName": "ns", "serverName": "su00-node00", "allocations": [one]}


def test_tenant_ips_pagination_and_status_filter(ipam_on):
    update_tenant(
        ipam_on, "ns", "gorax", operation="ADD",
        servers=[{"serverName": "su00-node00"}, {"serverName": "su00-node01"}, {"serverName": "su00-node02"}],
    )
    p0 = get_tenant_ips(ipam_on, "ns", "gorax", page=0, size=2)
    p1 = get_tenant_ips(ipam_on, "ns", "gorax", page=1, size=2)
    assert [a["ipAddress"] for a in p0["allocations"]] == ["192.3.0.2", "192.3.0.3"]
    assert [a["ipAddress"] for a in p1["allocations"]] == ["192.3.0.4"]
    assert p0["pagination"] == {"page": 0, "size": 2, "totalElements": 3, "totalPages": 2}

    assigned = get_tenant_ips(ipam_on, "ns", "gorax", status="ASSIGNED")
    assert assigned["pagination"]["totalElements"] == 3
    failed = get_tenant_ips(ipam_on, "ns", "gorax", status="IP_ASSIGN_FAILED")
    assert failed["allocations"] == [] and failed["pagination"]["totalElements"] == 0


def test_delete_server_releases_ip(ipam_on):
    update_tenant(ipam_on, "ns", "gorax", operation="ADD", servers=[{"serverName": "su00-node00"}])
    out = update_tenant(ipam_on, "ns", "gorax", operation="DELETE", servers=[{"serverName": "su00-node00"}])
    assert out["partial"] is False
    assert out["servers"][0]["released"] is True and out["servers"][0]["ipAddress"] == "192.3.0.2/32"
    assert get_server_ips(ipam_on, "ns", "su00-node00")["allocations"] == []
    assert code_of(http_error(get_tenant_compute_ip, ipam_on, "ns", "gorax", "su00-node00")) == "SERVER_NOT_ALLOCATED_TO_TENANT"


def test_delete_without_allocation_is_partial(ipam_on):
    out = update_tenant(ipam_on, "ns", "gorax", operation="DELETE", servers=[{"serverName": "su00-node02"}])
    assert out["partial"] is True
    assert out["servers"][0]["ipAssignStatus"] == "NOT_ALLOCATED"


def test_ip_read_errors(ipam_on):
    r = http_error(get_tenant_compute_ip, ipam_on, "ns", "gorax", "su00-node00")
    assert (r.status_code, code_of(r)) == (404, "SERVER_NOT_ALLOCATED_TO_TENANT")
    r = http_error(get_tenant_compute_ip, ipam_on, "ns", "ghost", "su00-node00")
    assert (r.status_code, code_of(r)) == (404, "IPAM_NOT_CONFIGURED")
    r = http_error(get_tenant_ips, ipam_on, "ns", "ghost")
    assert (r.status_code, code_of(r)) == (404, "IPAM_NOT_CONFIGURED")
    r = http_error(get_server_ips, ipam_on, "ns", "no-such-server")
    assert (r.status_code, code_of(r)) == (404, "SERVER_NOT_FOUND")


def test_create_tenant_forwards_dhcp_pool_on_the_wire(live_client, fm):
    # The fake has no POST /tenants route (404): only the request matters here.
    http_error(create_tenant, live_client, "ns", tenantName="t1", tenantDhcpCidrPool="10.9.0.0/24")
    assert json.loads(fm.requests[-1]["body"]) == {"tenantName": "t1", "tenantDhcpCidrPool": "10.9.0.0/24"}


def test_bad_token_is_a_401(fm):
    bad = Client.initialize_with_token(fm.base_url, "wrong")
    assert http_error(get_ipam_config, bad, "ns").status_code == 401
