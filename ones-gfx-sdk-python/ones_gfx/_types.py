"""Typed request payloads mirroring the FabricManager (SES_FM) Java DTOs.

Every type here maps 1:1 to a server-side request-body type referenced by
SES_FM/src/main/java/com/ONES/FabricManager. Server-managed fields (ids,
back-references, creation/update timestamps) are intentionally omitted since
the API ignores them on input.
"""
from __future__ import annotations

from typing import Literal, TypedDict

__all__ = [
    "ApiResponseMessage",
    "AuthData",
    "AuthResponse",
    "GpuAction",
    "GpuServerInfo",
    "NmxcDomain",
    "OperationAccepted",
    "SuidMap",
]

# ── Enums (Helper/Enums.java) ────────────────────────────────────────────────

GpuAction = Literal["ADD", "DELETE"]

# ── GPU allocation (Helper/GpuServerInfo.java, GpuStatusUpdate.suid) ────────


class GpuServerInfo(TypedDict):
    serverName: str
    shared: bool | None


# Java: TreeMap<Integer, Map<String, DeviceGpus{gpus: List<String>}>>
# JSON object keys are strings; Jackson coerces them to Integer server-side.
SuidMap = dict[str, dict[str, list[str]]]


# ── NMX-C domains (FabricManager nmxc endpoints) ───────────────────────────


class NmxcDomain(TypedDict):
    host: str
    port: int | None
    useTls: bool | None


# ─────────────────────────────────────────────────────────────────────────────
# Response types (FabricManager @ResponseBody shapes after the SDK client
# unwraps the {"data": ...} envelope). total=False everywhere a key can be
# absent or null on the wire.
# ─────────────────────────────────────────────────────────────────────────────


class OperationAccepted(TypedDict, total=False):
    """202 async bodies (ApiResponse.success(msg, asyncData) → data unwrapped).

    operationType appears only on GPU allocate/deallocate ops; webhookRegistered
    only in ASYNC_WEBHOOK mode.
    """

    operationId: str
    status: str
    operationType: str
    webhookRegistered: bool


class ApiResponseMessage(TypedDict):
    """ApiResponse.success(message) with no data — status/message always set."""

    status: str
    message: str


class _AuthDataOptional(TypedDict, total=False):
    """Auth ``data`` keys that only some endpoints return.

    Split out because TypedDict cannot mix required and optional keys in one
    class body, and ``NotRequired`` needs Python 3.11 (this package targets 3.9).
    """

    isPwdResetNeeded: bool


class AuthData(_AuthDataOptional):
    """The ``data`` object of POST /api/user/login and POST /api/user/refresh.

    ``message`` differs by endpoint — "Login Successful" vs "Token refreshed".
    ``isPwdResetNeeded`` is returned by login but not by refresh.
    """

    message: str
    token: str


class AuthResponse(TypedDict):
    """Full login/refresh response body.

    Unlike the ``ones_gfx.apis`` functions, ``ONESClient.login`` and
    ``ONESClient.refresh`` return the whole envelope without unwrapping ``data``.
    """

    data: AuthData


