"""
Shared vocabulary for the webhook registration endpoints.

Three independent scopes register callbacks: global (`/webhooks`), per-fabric
(`/fabrics/{f}/webhook`) and per-tenant (`/fabrics/{f}/tenants/{t}/webhook`).
There is no precedence between them — every registration matching an event
receives it, so a callback registered at two scopes is delivered twice.
"""
from __future__ import annotations

from typing import TypedDict

#: Canonical event names FM delivers.
WEBHOOK_EVENT_TYPES = (
    "tenant.create",
    "tenant.delete",
    "tenant.allocate",
    "tenant.deallocate",
    "tenant.vpcpeering",
    "fabric.deploy",
    "fabric.delete",
)

#: `fabric.*` spellings FM also accepts, resolved to canonical on write.
WEBHOOK_EVENT_ALIASES = {
    "fabric.tenant_create": "tenant.create",
    "fabric.tenant_delete": "tenant.delete",
    "fabric.tenant_allocate": "tenant.allocate",
    "fabric.tenant_deallocate": "tenant.deallocate",
    "fabric.tenant_vpcpeering": "tenant.vpcpeering",
}


class WebhookFilters(TypedDict, total=False):
    """Scope filters on a global registration. Empty lists mean "no filter"."""

    fabricNames: list[str]
    tenantNames: list[str]


class WebhookRegistration(TypedDict, total=False):
    """A registration as FM returns it.

    `id` is the surrogate key, present for tenant- and fabric-scoped rows.
    `webhookId` is the external id a global row is addressed by and is the one
    to hold on to; the two never both appear. `tenantName`/`fabricName` are
    omitted at the scopes they do not apply to, and `filters` appears only on
    a global registration.
    """

    id: int
    webhookId: str
    tenantName: str
    fabricName: str
    callbackUrl: str
    enabled: bool
    eventTypes: list[str]
    filters: WebhookFilters
    createdAt: str
    updatedAt: str
