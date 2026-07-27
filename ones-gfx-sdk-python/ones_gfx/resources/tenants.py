from __future__ import annotations

from typing import Any
from urllib.parse import urlparse

from ..enums import OperationMode
from ..models import Operation, Tenant
from ..transport import Transport

_OPERATION_KEYS = {"operationId", "status"}


class TenantsResource:
    def __init__(self, transport: Transport):
        self._transport = transport

    def list(self, fabric_name: str) -> list[Tenant]:
        _require(fabric_name, "fabric_name")
        body = self._transport.get(f"fabrics/{fabric_name}/tenants")
        items = body.get("tenants", []) if isinstance(body, dict) else []
        return [Tenant.from_api(item) for item in items]

    def get(self, fabric_name: str, name: str) -> Tenant:
        _require(fabric_name, "fabric_name")
        _require(name, "name")
        body = self._transport.get(f"fabrics/{fabric_name}/tenants/{name}")
        tenant_payload = body.get("tenant", body) if isinstance(body, dict) else body
        return Tenant.from_api(tenant_payload)

    def available_servers(self, fabric_name: str) -> list[str]:
        _require(fabric_name, "fabric_name")
        body = self._transport.get(f"fabrics/{fabric_name}/available_servers")
        return body.get("availableGPUs", []) if isinstance(body, dict) else []

    def create(
        self,
        fabric_name: str,
        name: str,
        description: str,
        max_gpus_allowed: int,
        *,
        shared: bool = False,
        mode: OperationMode = OperationMode.SYNCHRONOUS,
        webhook_url: str | None = None,
        webhook_events: list[str] | None = None,
    ) -> Tenant | Operation:
        _require(fabric_name, "fabric_name")
        _require(name, "name")
        _validate_max_gpus(max_gpus_allowed)
        body = {"tenantName": name, "description": description or "", "maxGpusAllowed": max_gpus_allowed, "shared": shared}
        _attach_webhook_fields(body, mode, webhook_url, webhook_events)
        result = self._transport.post(f"fabrics/{fabric_name}/tenants", json_body=body, mode=mode)
        if _looks_like_operation(result):
            return Operation.from_submission(result)
        if isinstance(result, str):
            return Tenant.from_api({})
        tenant_payload = result.get("tenant", result) if isinstance(result, dict) else result
        return Tenant.from_api(tenant_payload)

    def delete(
        self,
        fabric_name: str,
        name: str,
        *,
        mode: OperationMode = OperationMode.SYNCHRONOUS,
        webhook_url: str | None = None,
        webhook_events: list[str] | None = None,
    ) -> Operation | None:
        _require(fabric_name, "fabric_name")
        _require(name, "name")
        body = None
        if mode is OperationMode.ASYNC_WEBHOOK:
            body = {}
            _attach_webhook_fields(body, mode, webhook_url, webhook_events)
        if body is None:
            result = self._transport.delete(f"fabrics/{fabric_name}/tenants/{name}", mode=mode)
        else:
            result = self._transport._request("DELETE", f"fabrics/{fabric_name}/tenants/{name}", json_body=body, mode=mode)  # noqa: SLF001
        if _looks_like_operation(result):
            return Operation.from_submission(result)
        return None

    def allocate_gpus(self, fabric_name: str, name: str, servers: list[str | dict], *, mode: OperationMode = OperationMode.SYNCHRONOUS, webhook_url: str | None = None, webhook_events: list[str] | None = None, timeout: int | None = None) -> Operation | None:
        return self._gpu_update(fabric_name, name, "ADD", servers, mode=mode, webhook_url=webhook_url, webhook_events=webhook_events, timeout=timeout)

    def deallocate_gpus(self, fabric_name: str, name: str, servers: list[str | dict], *, mode: OperationMode = OperationMode.SYNCHRONOUS, webhook_url: str | None = None, webhook_events: list[str] | None = None, timeout: int | None = None) -> Operation | None:
        return self._gpu_update(fabric_name, name, "DELETE", servers, mode=mode, webhook_url=webhook_url, webhook_events=webhook_events, timeout=timeout)

    def _gpu_update(self, fabric_name: str, name: str, operation: str, servers: list[str | dict], *, mode: OperationMode, webhook_url: str | None, webhook_events: list[str] | None, timeout: int | None = None) -> Operation | None:
        _require(fabric_name, "fabric_name")
        _require(name, "name")
        if not servers:
            raise ValueError("servers list must not be empty.")
        body = {"operation": operation, "servers": [_normalize_server(s) for s in servers]}
        _attach_webhook_fields(body, mode, webhook_url, webhook_events)
        result = self._transport.patch(f"fabrics/{fabric_name}/tenants/{name}", json_body=body, mode=mode, timeout=timeout)
        if _looks_like_operation(result):
            return Operation.from_submission(result)
        return None

    def assign_ports(self, fabric_name: str, tenant_name: str, server_names: list[str], gpu_ids: list[int] | None = None) -> None:
        self._post_gpus(fabric_name, tenant_name, "ADD", server_names, gpu_ids)

    def unassign_ports(self, fabric_name: str, tenant_name: str, server_names: list[str], gpu_ids: list[int] | None = None) -> None:
        self._post_gpus(fabric_name, tenant_name, "DELETE", server_names, gpu_ids)

    def _post_gpus(self, fabric_name: str, tenant_name: str, operation: str, server_names: list[str], gpu_ids: list[int] | None) -> None:
        _require(fabric_name, "fabric_name")
        _require(tenant_name, "tenant_name")
        if not server_names:
            raise ValueError("server_names cannot be empty.")
        body: dict = {
            "operation": operation,
            "serverNames": server_names,
        }
        if gpu_ids:
            body["gpuIds"] = gpu_ids
        self._transport.post(
            f"fabrics/{fabric_name}/tenants/{tenant_name}/gpus",
            json_body=body,
            mode=OperationMode.SYNCHRONOUS,
        )


def _require(value: Any, field_name: str) -> None:
    if value is None or (isinstance(value, str) and not value.strip()):
        raise ValueError(f"{field_name} is required and cannot be empty.")


def _validate_max_gpus(value: int) -> None:
    if not isinstance(value, int) or isinstance(value, bool):
        raise ValueError(f"max_gpus_allowed must be an integer; got {type(value).__name__}.")
    if value == -1 or value >= 1:
        return
    raise ValueError(f"max_gpus_allowed must be a positive integer or -1 (unlimited); got {value}.")


def _normalize_server(spec: str | dict) -> dict:
    if isinstance(spec, str):
        return {"serverName": spec, "shared": False}
    if isinstance(spec, dict):
        if "serverName" not in spec:
            raise ValueError("Server spec dict must contain 'serverName'.")
        return {"serverName": spec["serverName"], "shared": bool(spec.get("shared", False))}
    raise TypeError(f"Server entries must be str or dict; got {type(spec).__name__}.")


def _attach_webhook_fields(body: dict, mode: OperationMode, webhook_url: str | None, webhook_events: list[str] | None) -> None:
    if mode is not OperationMode.ASYNC_WEBHOOK:
        if webhook_url is not None or webhook_events is not None:
            raise ValueError("webhook_url and webhook_events are only valid when mode=OperationMode.ASYNC_WEBHOOK.")
        return
    if not webhook_url:
        raise ValueError("webhook_url is required when mode is ASYNC_WEBHOOK.")
    parsed = urlparse(webhook_url)
    if not parsed.scheme or not parsed.netloc:
        raise ValueError(f"webhook_url must be an absolute URL with scheme and host; got {webhook_url!r}.")
    if not webhook_events or not isinstance(webhook_events, list):
        raise ValueError("webhook_events must be a non-empty list of strings.")
    for ev in webhook_events:
        if not isinstance(ev, str) or not ev.strip():
            raise ValueError(f"webhook_events entries must be non-empty strings; got {ev!r}.")
    body["enableWebhook"] = True
    body["webhookUrl"] = webhook_url
    body["webhookEvents"] = list(webhook_events)


def _looks_like_operation(result: Any) -> bool:
    return isinstance(result, dict) and _OPERATION_KEYS.issubset(result.keys())
