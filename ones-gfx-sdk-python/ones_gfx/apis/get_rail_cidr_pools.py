"""
GET /fabrics/{fabricName}/rails/cidrpools  ->  /api/fm/fabrics/{fabricName}/rails/cidrpools

The rail CIDRPool manifest for every East-West rail on the fabric.

Content negotiated. By default FM answers JSON and the SDK returns the parsed
body; ask for `application/yaml` (or the legacy `application/x-yaml`, or
`text/plain`) via `accept` and the raw manifest text is returned instead,
which is what you want when feeding it to a tool that consumes the YAML.

404 with `INVALID_FABRIC` when the fabric is unknown, or
`RAIL_CIDRPOOL_NOT_FOUND` when it has no East-West rails.
"""
from __future__ import annotations

from typing import Any

from ..client import Client

#: Media types FM will render the manifest as.
RAIL_YAML_MEDIA_TYPE = "application/yaml"
RAIL_LEGACY_YAML_MEDIA_TYPE = "application/x-yaml"


def get_rail_cidr_pools(
    client: Client, fabricName: str, *, accept: str | None = None
) -> Any:
    """The rail CIDRPool manifest for `fabricName`.

    Returns the parsed JSON body, or the raw manifest string when `accept`
    asks for YAML or plain text.
    """
    path = f"fabrics/{fabricName}/rails/cidrpools"
    headers: dict[str, str] = {}
    if accept is not None:
        headers["Accept"] = accept
    return client.call_api("GET", path, extra_headers=headers or None)
