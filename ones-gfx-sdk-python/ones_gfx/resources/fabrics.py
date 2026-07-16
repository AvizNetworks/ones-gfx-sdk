from __future__ import annotations

from ..models import Fabric
from ..transport import Transport


class FabricsResource:
    def __init__(self, transport: Transport):
        self._transport = transport

    def list(self) -> list[Fabric]:
        body = self._transport.get("fabrics")
        items = body.get("fabrics", []) if isinstance(body, dict) else []
        return [Fabric.from_api(item) for item in items]

    def inventory_sync(self, fabric_name: str) -> None:
        # POST /fabrics/{fabricName}/inventorySync — UFM enabled fabrics only.
        if not fabric_name:
            raise ValueError("fabric_name is required")
        self._transport.post(f"fabrics/{fabric_name}/inventorySync")
