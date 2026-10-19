"""
DELETE /fabrics/{fabricName}/inventory  ->  /api/fm/fabrics/{fabricName}/inventory

Delete inventory devices by hostname, scoped to one fabric.

The body is a bare JSON list of hostnames, not an object. Duplicates and
blanks are dropped before the delete, and at most
`MAX_INVENTORY_DELETE_HOSTNAMES` unique names may be sent.

The response always carries a `notFound` list of hostnames that matched
nothing, because the status alone does not say which: FM answers 200 when
every hostname matched, 207 when only some did, and 404 when none did.
"""
from __future__ import annotations

from typing import Any

from ..client import Client

#: Largest number of unique hostnames FM accepts in one call.
MAX_INVENTORY_DELETE_HOSTNAMES = 200


def delete_inventory_by_hostnames(
    client: Client, fabricName: str, hostnames: list[str]
) -> dict[str, Any]:
    """Delete the named devices from `fabricName`'s inventory.

    Returns `{"notFound": [...], "message": ...}`; check `notFound` even on a
    2xx, since a 207 means only some hostnames matched.
    """
    path = f"fabrics/{fabricName}/inventory"
    return client.call_api("DELETE", path, json_body=hostnames)
