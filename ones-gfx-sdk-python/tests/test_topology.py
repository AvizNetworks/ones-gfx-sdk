"""Unit tests for N-S interface extraction and GET /fabrics/{fabricName}."""
from __future__ import annotations

import json
from pathlib import Path

import pytest
import requests
import yaml

from ones_gfx import (
    Client,
    NotAuthenticatedError,
    TopologyError,
    extract_ns_interfaces,
    get_fabric_ns_interfaces,
    get_fabric_yaml,
    ns_interfaces_to_yaml,
)
from ones_gfx.topology import main

FIXTURES = Path(__file__).parent / "fixtures"
TOPOLOGY = (FIXTURES / "ns_topology.yaml").read_text()
EXPECTED = (FIXTURES / "ns_topology_expected.yaml").read_text()
BASE_URL = "https://ones.example"


def _ifaces(result: dict, host: str) -> list[tuple[str, ...]]:
    for h in result["hosts"]:
        if h["host"] == host:
            return [
                (i["interface"], i["leaf"], i["leafPort"]) + ((i["via"], i["viaInterface"]) if "via" in i else ())
                for i in h["nsInterfaces"]
            ]
    return []


# --- extraction --------------------------------------------------------------


def test_matches_golden_output():
    # The Go SDK checks against the same golden file, so the two stay identical.
    result = extract_ns_interfaces(TOPOLOGY, fabric_name="EW")
    assert ns_interfaces_to_yaml(result) == EXPECTED


def test_reads_both_sides_and_deduplicates():
    result = extract_ns_interfaces(TOPOLOGY)
    # hgx-su00-h00: each cable is listed under the host and the leaf; reported once.
    assert _ifaces(result, "hgx-su00-h00") == [
        ("enp8s0f1np0", "leaf-cn-00", "swp9"),
        ("enp9s0f1np0", "leaf-cn-01", "swp9"),
    ]
    # hgx-su00-h01: one link has no destination, the other is only on the leaf side.
    assert _ifaces(result, "hgx-su00-h01") == [
        ("enp8s0f1np0", "leaf-cn-00", "swp10"),
        ("enp9s0f1np0", "leaf-cn-01", "swp10"),
    ]
    assert result["hostCount"] == 4 and result["interfaceCount"] == 10


def test_default_covers_cn_fn_and_sn():
    result = extract_ns_interfaces(TOPOLOGY)
    assert result["leafPatterns"] == ["^leaf-cn", "^leaf-fn", "^leaf-sn"]
    # Front-end and dedicated-storage leaves count without any option.
    assert _ifaces(result, "hgx-su00-h02") == [
        ("enp12s0f1np0", "leaf-sn-00", "swp5"),
        ("enp8s0f1np0", "leaf-fn-00", "swp5"),
    ]


def test_ignores_ew_leaves_routers_and_dpus_without_uplinks():
    flat = [i for h in extract_ns_interfaces(TOPOLOGY)["hosts"] for i in h["nsInterfaces"]]
    assert not any(i["leaf"] == "leaf-su00-r0" for i in flat)  # E-W
    assert not any(i["leaf"].startswith("leaf-border") for i in flat)  # external routers
    # hgx-su00-h00's DPU has no uplinks, so its vf0 has no N-S path.
    assert not any(i.get("via") == "hgx-su00-h00-dpu-00" for i in flat)


def test_follows_host_dpu_leaf_paths():
    dpu = "hgx-su00-h03-dpu-00"
    assert _ifaces(extract_ns_interfaces(TOPOLOGY), "hgx-su00-h03") == [
        # No PF index in the DPU interface: every N-S uplink.
        ("hostif", "leaf-cn-00", "swp20", dpu, "mgmt_if"),
        ("hostif", "leaf-cn-01", "swp20", dpu, "mgmt_if"),
        # pf0 -> uplink p0 -> leaf-cn-00; pf1 -> p1 -> leaf-cn-01.
        ("vf0", "leaf-cn-00", "swp20", dpu, "pf0vf0_if"),
        ("vf4", "leaf-cn-01", "swp20", dpu, "pf1vf0_if"),
    ]


def test_dpu_entry_key_order():
    h03 = next(h for h in extract_ns_interfaces(TOPOLOGY)["hosts"] if h["host"] == "hgx-su00-h03")
    assert list(h03["nsInterfaces"][0]) == ["interface", "via", "viaInterface", "leaf", "leafPort"]
    h00 = next(h for h in extract_ns_interfaces(TOPOLOGY)["hosts"] if h["host"] == "hgx-su00-h00")
    assert list(h00["nsInterfaces"][0]) == ["interface", "leaf", "leafPort"]


def test_custom_and_multiple_patterns():
    fn = extract_ns_interfaces(TOPOLOGY, leaf_patterns="^leaf-fn")
    assert [h["host"] for h in fn["hosts"]] == ["hgx-su00-h02"]
    assert fn["leafPatterns"] == ["^leaf-fn"]

    cn = extract_ns_interfaces(TOPOLOGY, leaf_patterns=["^leaf-cn"])
    assert [h["host"] for h in cn["hosts"]] == ["hgx-su00-h00", "hgx-su00-h01", "hgx-su00-h03"]


def test_credentials_never_copied():
    out = ns_interfaces_to_yaml(extract_ns_interfaces(TOPOLOGY, fabric_name="EW"))
    for secret in ("secret", "password", "Credentials", "admin", "root"):
        assert secret not in out


def test_no_ns_links_gives_empty_result():
    result = extract_ns_interfaces(TOPOLOGY, leaf_patterns="^does-not-exist")
    assert result["hosts"] == [] and result["hostCount"] == 0 and result["interfaceCount"] == 0


def test_accepts_bytes_and_parsed_dict():
    as_dict = yaml.safe_load(TOPOLOGY)
    assert extract_ns_interfaces(TOPOLOGY.encode()) == extract_ns_interfaces(as_dict)


@pytest.mark.parametrize(
    "source",
    ["not: [valid", "<!doctype html><html></html>", "Inventory: {}\n", "- a list\n"],
)
def test_rejects_non_topology(source):
    with pytest.raises(TopologyError):
        extract_ns_interfaces(source)


@pytest.mark.parametrize("patterns", [[], ["("]])
def test_rejects_bad_patterns(patterns):
    with pytest.raises(ValueError):
        extract_ns_interfaces(TOPOLOGY, leaf_patterns=patterns)


# --- fetch: GET /fabrics/{fabricName} ---------------------------------------


class StubSession:
    def __init__(self, status: int, body: str, content_type: str = "application/x-yaml"):
        self.status, self.body, self.content_type = status, body, content_type
        self.calls: list[dict] = []

    def request(self, **kwargs):
        self.calls.append(kwargs)
        resp = requests.Response()
        resp.status_code = self.status
        resp._content = self.body.encode()
        resp.headers["Content-Type"] = self.content_type
        resp.encoding = "utf-8"
        return resp

    def close(self) -> None:
        pass


@pytest.fixture
def client():
    return Client(baseUrl=BASE_URL, authToken="tok")


def test_get_fabric_yaml_uses_gateway_root(client):
    client._session = session = StubSession(200, TOPOLOGY)
    assert get_fabric_yaml(client, "EW") == TOPOLOGY
    call = session.calls[0]
    assert call["method"] == "GET"
    assert call["url"] == f"{BASE_URL}/fabrics/EW"  # root, not /api/fm/
    assert call["headers"]["authorization"] == "tok"


def test_get_fabric_yaml_escapes_name(client):
    client._session = session = StubSession(200, TOPOLOGY)
    get_fabric_yaml(client, "full topo check")
    assert session.calls[0]["url"] == f"{BASE_URL}/fabrics/full%20topo%20check"


def test_get_fabric_ns_interfaces_end_to_end(client):
    client._session = StubSession(200, TOPOLOGY)
    result = get_fabric_ns_interfaces(client, "EW")
    assert ns_interfaces_to_yaml(result) == EXPECTED


def test_get_fabric_yaml_404(client):
    body = json.dumps({"status": "FAILED", "code": 404, "errorType": "NOT_FOUND",
                       "message": 'YAML file for fabric "nope" not found.'})
    client._session = StubSession(404, body, "application/json")
    with pytest.raises(requests.HTTPError) as excinfo:
        get_fabric_yaml(client, "nope")
    assert excinfo.value.response.status_code == 404
    assert excinfo.value.response.json()["errorType"] == "NOT_FOUND"


@pytest.mark.parametrize("status", [401, 500])
def test_get_fabric_yaml_errors(client, status):
    client._session = StubSession(status, '{"message":"x"}', "application/json")
    with pytest.raises(requests.HTTPError):
        get_fabric_yaml(client, "EW")


def test_get_fabric_ns_interfaces_html_fallback(client):
    # A gateway without the route serves the UI's index.html with HTTP 200.
    client._session = StubSession(200, "<!doctype html><html></html>", "text/html")
    with pytest.raises(TopologyError):
        get_fabric_ns_interfaces(client, "EW")


def test_requires_auth_and_name():
    with pytest.raises(NotAuthenticatedError):
        get_fabric_yaml(Client(baseUrl=BASE_URL), "EW")
    with pytest.raises(ValueError):
        get_fabric_yaml(Client(baseUrl=BASE_URL, authToken="tok"), " ")


# --- CLI ---------------------------------------------------------------------


def test_cli_file_to_output(tmp_path, capsys):
    out = tmp_path / "ns.yaml"
    assert main([str(FIXTURES / "ns_topology.yaml"), "--fabric", "EW", "-o", str(out)]) == 0
    assert out.read_text() == EXPECTED

    assert main([str(FIXTURES / "ns_topology.yaml"), "-p", "^leaf-fn"]) == 0
    out = capsys.readouterr().out
    assert "hgx-su00-h02" in out and "hgx-su00-h00" not in out


def test_cli_bad_input(tmp_path, capsys):
    bad = tmp_path / "bad.yaml"
    bad.write_text("<html></html>")
    assert main([str(bad)]) == 1
    assert "not a fabric topology" in capsys.readouterr().err
