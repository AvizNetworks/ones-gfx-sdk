"""Fabric Manager API functions — one module per endpoint.

Every function routes through ONESClient.call_api, which enforces the
auth guard (login required) and the /api/fm/ path prefix.
"""
from __future__ import annotations

from .create_tenant import create_tenant
from .get_all_fabrics import get_all_fabrics
from .get_gpu_tenant_mappings import get_gpu_tenant_mappings
from .get_tenant import get_tenant
from .upload_file import upload_file

__all__ = [
    "create_tenant",
    "get_all_fabrics",
    "get_gpu_tenant_mappings",
    "get_tenant",
    "upload_file",
]
