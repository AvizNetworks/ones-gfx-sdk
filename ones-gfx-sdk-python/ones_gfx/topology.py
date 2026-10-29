"""
North-south (N-S) host interfaces from a fabric topology YAML.

ONES keeps each fabric's topology as the YAML it was deployed from, served by
the ONES gateway at ``GET /fabrics/{fabricName}`` (see `get_fabric_yaml`). It
describes every device and cable. This module reduces it to one question:
which interfaces on each host are N-S, and where do they plug in.

A host link is N-S when the leaf at its far end has a name matching one of the
leaf patterns. The defaults cover the three N-S designs ONES generates:
``^leaf-cn`` (converged user + storage), ``^leaf-fn`` (front-end network) and
``^leaf-sn`` (dedicated storage network). E-W leaves always carry an ``su``
token (``leaf-su00-r0``, ``leaf-p0-su00-r1``) and never match. Pass your own
patterns for a reference architecture that names its leaves differently.

How the topology is read:

* Leaves come from ``Connectivity.Leaf`` (``switchId`` -> ``switchName``),
  hosts from ``Connectivity.Host`` (``hostId`` -> ``hostName``), DPUs from
  ``Connectivity.DPU`` (``switchId`` -> ``switchName``).
* Links are strings such as ``"G1_enp8s0f1np0 | L39_swp9"``: ``G<hostId>`` is a
  host, ``L<switchId>`` a leaf, ``D<switchId>`` a DPU, and the rest of each
  side is the interface.
* Links are read from the Host, Leaf and DPU sections alike, so a cable listed
  on only one side is still found; a cable listed on both is reported once.
* HBN fabrics put a DPU between host and leaf: host ``vf0`` -> DPU
  ``pf0vf0_if``, and the DPU's uplink ``p0_if`` -> leaf. A VF on ``pf<N>``
  leaves through uplink ``p<N>``, so it is reported with that uplink's leaf
  port and ``via``/``viaInterface`` naming the DPU and its VF interface (the
  ``vfId`` the VF-assign API takes). If the PF cannot be told from the name,
  every N-S uplink of the DPU is reported.
* Endpoints that are not cluster hosts (``H<n>`` external routers on border
  leaves) are ignored.
* Only names and interfaces are copied out; the ``Credentials`` blocks in the
  input are never read into the result.
"""
from __future__ import annotations

import argparse
import re
import sys
from collections.abc import Iterable
from typing import Any

import yaml

#: Leaf-name patterns used when none are given: the converged (``leaf-cn``),
#: front-end (``leaf-fn``) and dedicated storage (``leaf-sn``) N-S leaves.
DEFAULT_NS_LEAF_PATTERNS: tuple[str, ...] = ("^leaf-cn", "^leaf-fn", "^leaf-sn")

_ENDPOINT = re.compile(r"^([A-Za-z]+)(\d+)_(.+)$")
_PF = re.compile(r"^pf(\d+)")
_UPLINK = re.compile(r"^p(\d+)(?:_|$)")


class TopologyError(ValueError):
    """The input is not a fabric topology this module can read."""


def _endpoint(text: str) -> tuple[str, str] | None:
    """``"G1_enp8s0f1np0"`` -> ``("G1", "enp8s0f1np0")``; None if malformed."""
    m = _ENDPOINT.match(text.strip())
    if not m:
        return None
    return f"{m.group(1).upper()}{int(m.group(2))}", m.group(3)


def _link_sides(link: Any) -> tuple[tuple[str, str], tuple[str, str]] | None:
    if not isinstance(link, dict) or not isinstance(link.get("link"), str):
        return None
    parts = link["link"].split("|")
    if len(parts) != 2:
        return None
    a, b = _endpoint(parts[0]), _endpoint(parts[1])
    if a is None or b is None:
        return None
    return a, b


def _load(source: str | bytes | dict) -> dict:
    if isinstance(source, dict):
        doc = source
    else:
        try:
            doc = yaml.safe_load(source)
        except yaml.YAMLError as e:
            raise TopologyError(f"not valid YAML: {e}") from e
    if not isinstance(doc, dict) or not isinstance(doc.get("Connectivity"), dict):
        raise TopologyError(
            "not a fabric topology: no 'Connectivity' section "
            "(an HTML page here usually means the gateway did not serve GET /fabrics/{name})"
        )
    return doc


def _compile(patterns: Iterable[str] | str | None) -> list[re.Pattern[str]]:
    if patterns is None:
        patterns = DEFAULT_NS_LEAF_PATTERNS
    elif isinstance(patterns, str):
        patterns = [patterns]
    patterns = list(patterns)
    if not patterns:
        raise ValueError("leaf_patterns must contain at least one pattern")
    try:
        return [re.compile(p) for p in patterns]
    except re.error as e:
        raise ValueError(f"invalid leaf pattern: {e}") from e


def _names(devices: list, id_key: str, name_key: str, prefix: str) -> dict[str, str]:
    out: dict[str, str] = {}
    for d in devices:
        if isinstance(d, dict) and d.get(id_key) is not None and d.get(name_key):
            try:
                out[f"{prefix}{int(d[id_key])}"] = str(d[name_key])
            except (TypeError, ValueError):
                continue
    return out


def _devices(section: Any) -> list:
    return section if isinstance(section, list) else []


def extract_ns_interfaces(
    source: str | bytes | dict,
    fabric_name: str | None = None,
    leaf_patterns: Iterable[str] | str | None = None,
) -> dict[str, Any]:
    """Reduce a fabric topology to each host's N-S interfaces.

    Args:
        source: the topology YAML (text or bytes), or an already-parsed dict.
        fabric_name: echoed into the result as ``fabric``.
        leaf_patterns: regular expressions matched (``re.search``) against leaf
            names; a link is N-S when its leaf matches any of them. Defaults to
            ``DEFAULT_NS_LEAF_PATTERNS``.

    Returns:
        ``{"fabric", "leafPatterns", "hostCount", "interfaceCount", "hosts"}``,
        where ``hosts`` lists only hosts with at least one N-S interface, in
        topology order, each as ``{"host", "nsInterfaces": [...]}``. An
        interface is ``{"interface", "leaf", "leafPort"}``, plus ``via`` and
        ``viaInterface`` when it reaches the leaf through a DPU.

    Raises:
        TopologyError: the input is not a fabric topology.
        ValueError: a leaf pattern is not a valid regular expression.
    """
    compiled = _compile(leaf_patterns)
    conn = _load(source)["Connectivity"]
    hosts_raw = _devices(conn.get("Host"))
    leaves_raw = _devices(conn.get("Leaf"))
    dpus_raw = _devices(conn.get("DPU"))

    host_names = _names(hosts_raw, "hostId", "hostName", "G")
    leaf_names = _names(leaves_raw, "switchId", "switchName", "L")
    dpu_names = _names(dpus_raw, "switchId", "switchName", "D")

    def is_ns(leaf_name: str) -> bool:
        return any(p.search(leaf_name) for p in compiled)

    direct: set[tuple[str, str, str, str]] = set()  # (hostKey, iface, leaf, port)
    vfs: set[tuple[str, str, str, str]] = set()  # (hostKey, iface, dpuKey, dpuIface)
    uplinks: dict[str, set[tuple[str, str, str]]] = {}  # dpuKey -> (dpuIface, leaf, port)

    for section, from_host in ((hosts_raw, True), (leaves_raw, False), (dpus_raw, False)):
        for device in section:
            if not isinstance(device, dict):
                continue
            for link in device.get("Links") or []:
                sides = _link_sides(link)
                if sides is None:
                    continue
                for (ka, pa), (kb, pb) in (sides, sides[::-1]):
                    if ka in host_names and kb.startswith("L"):
                        leaf = leaf_names.get(kb)
                        if leaf is None and from_host:
                            # In the Host section, destination names the leaf.
                            dest = (link.get("properties") or {}).get("destination")
                            leaf = str(dest) if dest else None
                        if leaf is not None:
                            direct.add((ka, pa, leaf, pb))
                    elif ka in host_names and kb in dpu_names:
                        vfs.add((ka, pa, kb, pb))
                    elif ka in dpu_names and kb in leaf_names:
                        uplinks.setdefault(ka, set()).add((pa, leaf_names[kb], pb))

    found: dict[str, set[tuple[str, str, str, str, str]]] = {}
    order: list[str] = []

    def add(host_key: str, iface: str, leaf: str, port: str, via: str = "", via_iface: str = "") -> None:
        if not is_ns(leaf):
            return
        name = host_names[host_key]
        if name not in found:
            found[name] = set()
            order.append(name)
        found[name].add((iface, leaf, port, via, via_iface))

    for host_key, iface, leaf, port in sorted(direct):
        add(host_key, iface, leaf, port)

    for host_key, iface, dpu_key, dpu_iface in sorted(vfs):
        ups = sorted(uplinks.get(dpu_key, ()))
        pf = _PF.match(dpu_iface)
        if pf:
            same = [u for u in ups if (m := _UPLINK.match(u[0])) and m.group(1) == pf.group(1)]
            ups = same or ups
        for _, leaf, port in ups:
            add(host_key, iface, leaf, port, dpu_names[dpu_key], dpu_iface)

    # Host-section order first; any other host follows in the order found.
    rank = {name: i for i, name in enumerate(host_names.values())}

    def entry(i: str, leaf: str, port: str, via: str, via_iface: str) -> dict[str, str]:
        e = {"interface": i}
        if via:
            e["via"] = via
            e["viaInterface"] = via_iface
        e["leaf"] = leaf
        e["leafPort"] = port
        return e

    hosts = [
        {"host": name, "nsInterfaces": [entry(*t) for t in sorted(found[name])]}
        for name in sorted(order, key=lambda n: rank.get(n, len(rank)))
    ]
    return {
        "fabric": fabric_name,
        "leafPatterns": [p.pattern for p in compiled],
        "hostCount": len(hosts),
        "interfaceCount": sum(len(h["nsInterfaces"]) for h in hosts),
        "hosts": hosts,
    }


class _IndentedDumper(yaml.SafeDumper):
    """SafeDumper that indents list items under their key, as Go's yaml.v3
    does, so both SDKs print the same bytes for the same result."""

    def increase_indent(self, flow: bool = False, indentless: bool = False) -> None:
        return super().increase_indent(flow, False)


def ns_interfaces_to_yaml(result: dict[str, Any]) -> str:
    """Render an `extract_ns_interfaces` result as clean, stable YAML.

    Two-space indent with list items indented under their key; identical to
    the Go SDK's NSInterfaces.YAML() output.
    """
    return yaml.dump(
        result, Dumper=_IndentedDumper, sort_keys=False, default_flow_style=False,
        allow_unicode=True, indent=2,
    )


def main(argv: list[str] | None = None) -> int:
    """``ones-gfx-ns-interfaces``: print each host's N-S interfaces as YAML.

    Reads a topology file, or fetches it from ONES with --url/--username/--password.
    """
    p = argparse.ArgumentParser(
        prog="ones-gfx-ns-interfaces",
        description="List each host's north-south interfaces from a fabric topology YAML.",
    )
    p.add_argument("file", nargs="?", help="topology YAML (as downloaded from GET /fabrics/{name})")
    p.add_argument("--fabric", "-f", help="fabric name (required with --url)")
    p.add_argument("--leaf-pattern", "-p", action="append", dest="patterns",
                   help=f"leaf-name regex; repeatable (default: {' '.join(DEFAULT_NS_LEAF_PATTERNS)})")
    p.add_argument("--output", "-o", help="write YAML here instead of stdout")
    p.add_argument("--url", "-u", help="ONES base URL, to fetch the topology instead of reading a file")
    p.add_argument("--username", "-U")
    p.add_argument("--password", "-P")
    p.add_argument("--verify-tls", action="store_true", help="enforce TLS verification (default: off)")
    args = p.parse_args(argv)

    if args.url:
        if not (args.fabric and args.username and args.password):
            p.error("--url needs --fabric, --username and --password")
        from .apis.get_fabric_ns_interfaces import get_fabric_yaml
        from .client import Client

        client = Client.initialize_with_creds(args.url, args.username, args.password, verify_tls=args.verify_tls)
        try:
            text = get_fabric_yaml(client, args.fabric)
        finally:
            client.close()
    elif args.file:
        with open(args.file, encoding="utf-8") as fh:
            text = fh.read()
    else:
        p.error("give a topology file, or --url with --fabric, --username and --password")

    try:
        result = extract_ns_interfaces(text, fabric_name=args.fabric, leaf_patterns=args.patterns)
    except (TopologyError, ValueError) as e:
        print(f"error: {e}", file=sys.stderr)
        return 1
    out = ns_interfaces_to_yaml(result)
    if args.output:
        with open(args.output, "w", encoding="utf-8") as fh:
            fh.write(out)
    else:
        sys.stdout.write(out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
