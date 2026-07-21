from __future__ import annotations

from typing import Any

from ..models import Fabric
from ..transport import Transport


class FabricsResource:
    def __init__(self, transport: Transport):
        self._transport = transport

    def list(self) -> list[Fabric]:
        body = self._transport.get("fabrics")
        items = body.get("fabrics", []) if isinstance(body, dict) else []
        return [Fabric.from_api(item) for item in items]

    def modify_gpu_allocations(
        self,
        fabric_name: str,
        tenant_name: str,
        operation: str,
        suid: dict[str, dict[str, dict[str, list[str]]]],
    ) -> dict[str, Any]:
        """Map or unmap specific GPUs to a tenant using fine-grained suid addressing.

        POST /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations

        Args:
            fabric_name:  Target fabric.
            tenant_name:  Target tenant.
            operation:    ``"ADD"`` to map GPUs or ``"DELETE"`` to unmap.
            suid:         Nested map: server_index -> hostname -> {"gpus": ["G0", ...]}.

        Example::

            client.fabrics.modify_gpu_allocations(
                fabric_name="fab1",
                tenant_name="tenant1",
                operation="ADD",
                suid={
                    "0": {
                        "hgx-su00-h00": {"gpus": ["G0", "G1", "G2", "G3"]}
                    }
                },
            )

        Returns:
            dict with keys ``status``, ``operationId`` (if async), ``message``.
        """
        if not fabric_name:
            raise ValueError("fabric_name is required")
        if not tenant_name:
            raise ValueError("tenant_name is required")
        if operation not in ("ADD", "DELETE"):
            raise ValueError(f'operation must be "ADD" or "DELETE", got {operation!r}')
        if not suid:
            raise ValueError("suid must contain at least one server entry")

        body = {"operation": operation, "suid": suid}
        path = f"fabrics/{fabric_name}/tenants/{tenant_name}/gpuAllocations"
        result = self._transport.post(path, json_body=body)
        if isinstance(result, dict):
            return result
        if isinstance(result, str):
            return {"status": "success", "message": result}
        if result is None:
            return {"status": "success"}
        raise TypeError(
            f"unexpected response type for modify_gpu_allocations: {type(result).__name__}"
        )
