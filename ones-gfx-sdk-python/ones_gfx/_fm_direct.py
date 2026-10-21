"""
Which Fabric Manager routes the SDK must reach on Fabric Manager directly.

The ONES gateway proxies ``/api/fm/*`` through an explicit per-route allowlist
(``fm-router.js``). Routes missing from that allowlist do not 404 — they fall
through to the UI's single-page-app catch-all and come back as **HTML with
HTTP 200**, which is far more confusing than a refusal.

Thirty-four routes are in that state as of the 4.4.0 release: the endpoints
added for network policy, webhook registration, IPAM configuration and address
queries, tenant DHCP and rail CIDRPools, plus the
shared-services and VF-interface routes that have been unreachable since
GFX-5709 merged. The 4.4.0 train has shipped, so the gateway cannot be changed
for this release and the SDK routes these calls to Fabric Manager itself.

Everything else keeps going through the gateway, which is still the supported
path and the one the published API spec documents.

**This module is meant to be deleted.** Once a gateway release registers these
routes, drop the entries here and every call goes back through the gateway with
no other change. An empty table restores the original behaviour exactly.
"""
from __future__ import annotations

import re
from urllib.parse import urlsplit, urlunsplit

#: Port Fabric Manager serves on. The deployment publishes it on the same host
#: as the gateway, so the direct base URL is derived rather than configured.
FM_DIRECT_PORT = 8443

#: Path templates, as passed to ``Client.call_api``, that bypass the gateway.
#: ``{}`` stands for one path segment.
DIRECT_FM_PATHS: tuple[str, ...] = (
    # Network policy
    "fabrics/{}/networkpolicy",
    "fabrics/{}/networkpolicy/{}",
    # Webhook registration: global, fabric and tenant scopes
    "webhooks",
    "webhooks/{}",
    "fabrics/{}/webhook",
    "fabrics/{}/tenants/{}/webhook",
    # IPAM address queries
    "fabrics/{}/tenants/{}/ip-addresses",
    "fabrics/{}/computes/{}/ip-addresses",
    "fabrics/{}/tenants/{}/computes/{}/ip-address",
    # IPAM configuration
    "fabrics/{}/ipam/config",
    # Tenant DHCP
    "fabrics/{}/tenants/{}/dhcp",
    "fabrics/{}/tenants/{}/dhcp/retry",
    "fabrics/{}/dhcp/reconcile",
    # Rail CIDRPools
    "fabrics/{}/rails/cidrpools",
    "fabrics/{}/rails/{}/cidrpool",
    # Fabric listing
    "fabrics",
    # Shared services and VF interfaces (GFX-5709)
    "fabrics/{}/shared-services",
    "fabrics/{}/servers/{}/vf-interfaces",
    "fabrics/{}/servers/{}/vf-interfaces/{}/assign",
)


def _compile(template: str) -> re.Pattern[str]:
    # A path segment is anything but "/", so a template matches one concrete
    # path and nothing deeper: "fabrics/{}" must not swallow "fabrics/x/y".
    parts = [re.escape(p) for p in template.split("{}")]
    return re.compile("^" + "[^/]+".join(parts) + "$")


_PATTERNS: tuple[re.Pattern[str], ...] = tuple(_compile(t) for t in DIRECT_FM_PATHS)


def is_direct_path(path: str) -> bool:
    """True when `path` must bypass the gateway and go straight to FM."""
    p = path.strip("/").split("?", 1)[0]
    return any(pat.match(p) for pat in _PATTERNS)


def derive_fm_base_url(base_url: str, port: int = FM_DIRECT_PORT) -> str:
    """Fabric Manager's own base URL, from the gateway's.

    The deployment publishes FM on the same host as the gateway, so only the
    port changes. Pass ``fm_direct_url`` to :class:`~ones_gfx.client.Client`
    when a deployment puts it somewhere else.
    """
    parts = urlsplit(base_url)
    host = parts.hostname or parts.netloc
    netloc = f"[{host}]:{port}" if host and ":" in host else f"{host}:{port}"
    return urlunsplit((parts.scheme or "https", netloc, "", "", "")).rstrip("/")
