"""
Shared helpers for the ONES UI config-service APIs (fabric generation and
deployment): device classification, the "Yes"/"No" convention, and turning the
generated fabric YAML into the JSON document that POST /api/config expects.
"""
from __future__ import annotations

import json
import re
from typing import Any

from .._types import FabricDevice


def yes_no(value: bool) -> str:
    """The backend expects "Yes"/"No" strings, as in the UI's device spreadsheet."""
    return "Yes" if value else "No"


def apply_config(device: FabricDevice) -> bool:
    return device.get("apply_config", True) is not False


def is_switch(device: FabricDevice) -> bool:
    kind = (device.get("device_type") or "").strip().lower()
    role = (device.get("device_role") or "").strip().lower()
    return kind == "switch" or role in ("spine", "leaf")


def is_server(device: FabricDevice) -> bool:
    kind = (device.get("device_type") or "").strip().lower()
    role = (device.get("device_role") or "").strip().lower()
    return kind in ("server", "dpu", "host") or role == "host"


def host_map_from_hosts_per_su(raw: str) -> dict[str, str]:
    """Parse the "{0:1,1:2}" SU host-count shorthand into {"0": "1", "1": "2"}."""
    trimmed = raw.strip()
    trimmed = trimmed[1:] if trimmed.startswith("{") else trimmed
    trimmed = (trimmed[:-1] if trimmed.endswith("}") else trimmed).strip()
    host_map: dict[str, str] = {}
    if not trimmed:
        return host_map
    for pair in trimmed.split(","):
        key, sep, value = pair.partition(":")
        if not sep or not key.strip() or not value.strip():
            raise ValueError(f'expected "su:count" pairs like "{{0:1,1:2}}", got {pair!r}')
        host_map[key.strip()] = value.strip()
    return host_map


def _yaml():
    try:
        import yaml
    except ImportError as exc:  # pragma: no cover - pyyaml is a declared dependency
        raise ImportError("fabric YAML handling needs PyYAML: pip install pyyaml") from exc
    return yaml


def _loader() -> type:
    """A SafeLoader that keeps YAML 1.1 surprises as plain strings.

    PyYAML follows YAML 1.1, where yes/no/on/off become booleans, values like
    12:30 become base-60 integers, 0755 is octal and dates become datetime
    objects. The fabric YAML has to round-trip to the server unchanged, so only
    true/false, plain decimal and hex integers, and floats with a dot or
    exponent are typed; everything else stays a string.
    """
    yaml = _yaml()
    if hasattr(_loader, "cached"):
        return _loader.cached  # type: ignore[attr-defined]

    class FabricYamlLoader(yaml.SafeLoader):
        pass

    typed = {"tag:yaml.org,2002:bool", "tag:yaml.org,2002:int", "tag:yaml.org,2002:float",
             "tag:yaml.org,2002:timestamp", "tag:yaml.org,2002:value"}
    FabricYamlLoader.yaml_implicit_resolvers = {
        first: [r for r in resolvers if r[0] not in typed]
        for first, resolvers in yaml.SafeLoader.yaml_implicit_resolvers.items()
    }
    FabricYamlLoader.add_implicit_resolver(
        "tag:yaml.org,2002:bool", re.compile(r"^(?:true|True|TRUE|false|False|FALSE)$"), list("tTfF"))
    FabricYamlLoader.add_implicit_resolver(
        "tag:yaml.org,2002:int", re.compile(r"^(?:[-+]?(?:0|[1-9][0-9]*)|0x[0-9a-fA-F]+)$"), list("-+0123456789"))
    FabricYamlLoader.add_implicit_resolver(
        "tag:yaml.org,2002:float",
        re.compile(r"^[-+]?(?:[0-9]+\.[0-9]*|\.[0-9]+|[0-9]+(?=[eE]))(?:[eE][-+]?[0-9]+)?$"),
        list("-+0123456789."))
    _loader.cached = FabricYamlLoader  # type: ignore[attr-defined]
    return FabricYamlLoader


def inject_device_credentials(doc: Any, devices: list[FabricDevice]) -> None:
    """Set ipAddress and Credentials on every device in doc["Connectivity"].

    Connectivity holds tiers (SSpine, Spine, Leaf, Tor, Host), each a list of
    devices named by switchName or hostName. Devices are matched by hostname,
    case-insensitively; only keys already present are overwritten.
    """
    if not isinstance(doc, dict) or not isinstance(doc.get("Connectivity"), dict):
        return
    by_name = {(d.get("hostname") or "").strip().lower(): d for d in devices}
    for tier in doc["Connectivity"].values():
        if not isinstance(tier, list):
            continue
        for entry in tier:
            if not isinstance(entry, dict):
                continue
            name = entry["switchName"] if "switchName" in entry else entry.get("hostName")
            match = by_name.get(str(name).strip().lower()) if name is not None else None
            if match is None:
                continue
            if "ipAddress" in entry:
                entry["ipAddress"] = match.get("ip", "")
            creds = entry.get("Credentials")
            if isinstance(creds, dict):
                if "user" in creds:
                    creds["user"] = match.get("username", "")
                if "password" in creds:
                    creds["password"] = match.get("password", "")


def fabric_yaml_to_document(raw_yaml: str, devices: list[FabricDevice]) -> Any:
    """Parse fabric YAML, fill in device credentials, and return the document.

    Mapping key order is preserved (dicts keep insertion order), so the JSON
    sent to POST /api/config matches the YAML's layout.
    """
    doc = _yaml().load(raw_yaml, Loader=_loader())
    inject_device_credentials(doc, devices)
    json.dumps(doc)  # fail early on anything that cannot be sent as JSON
    return doc
