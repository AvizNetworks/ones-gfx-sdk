"""Client-side argument checks shared by the IPAM / rail endpoints.

These mirror rules the server enforces anyway, so a bad call fails fast with a
``ValueError`` (and never builds a malformed URL) instead of a round-trip 4xx.
"""
from __future__ import annotations

import re
from urllib.parse import quote


def path_segment(value: object, label: str) -> str:
    """Return ``value`` URL-quoted for use as one path segment.

    Raises:
        ValueError: if ``value`` is not a non-blank string.
    """
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{label} is required and must be a non-empty string.")
    return quote(value, safe="")


def require_bool(value: object, label: str) -> bool:
    """Reject non-bool values (the server also rejects them with a 400)."""
    if not isinstance(value, bool):
        raise ValueError(f"{label} must be a boolean, got {type(value).__name__}.")
    return value


_POOL_CIDR_RE = re.compile(r"[^/]+/[0-9]{1,2}")


def validate_pool_cidr(value: object, label: str) -> str:
    """Check a tenant IPAM pool the way the server's IpamService does.

    An IPv4 network with a prefix length between /8 and /30 (host bits are fine —
    the server canonicalises ``192.3.0.5/17`` to ``192.3.0.0/17``). Surrounding
    whitespace is stripped, and the stripped value is returned — the server matches
    the CIDR untrimmed, so it must not be sent padded.

    Raises:
        ValueError: with the server's error code (``INVALID_CIDR_FORMAT`` /
            ``CIDR_POOL_TOO_SMALL``) in the message.
    """
    if not isinstance(value, str) or not _POOL_CIDR_RE.fullmatch(value.strip()):
        raise ValueError(
            f"{label} must be an IPv4 network with a prefix length, e.g. 192.3.0.0/17 (INVALID_CIDR_FORMAT)."
        )
    value = value.strip()
    addr, prefix_s = value.split("/")
    octets = addr.split(".")
    if len(octets) != 4 or not all(o.isascii() and o.isdigit() and int(o) <= 255 for o in octets):
        raise ValueError(f"{label} {value!r} is not a valid IPv4 network (INVALID_CIDR_FORMAT).")
    prefix = int(prefix_s)
    if not 8 <= prefix <= 30:
        raise ValueError(f"{label} prefix length must be between /8 and /30, got /{prefix} (CIDR_POOL_TOO_SMALL).")
    return value
