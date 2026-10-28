from __future__ import annotations

from typing import Any

from ..transport import Transport


class PeeringResource:
    def __init__(self, transport: Transport):
        self._transport = transport

    def create(self, fabric_name: str, name: str, vpc_name: str, peer_vpc_name: str) -> Any:
        if not fabric_name:
            raise ValueError("fabric_name is required.")
        if not name:
            raise ValueError("name is required.")
        if not vpc_name:
            raise ValueError("vpc_name is required.")
        if not peer_vpc_name:
            raise ValueError("peer_vpc_name is required.")
        body = {"name": name, "vpcname": vpc_name, "peervpcname": peer_vpc_name}
        return self._transport.post(f"fabrics/{fabric_name}/vpcpeering", json_body=body)

    def delete(self, fabric_name: str, name: str, vpc_name: str, peer_vpc_name: str) -> Any:
        if not fabric_name:
            raise ValueError("fabric_name is required.")
        if not name:
            raise ValueError("name is required.")
        if not vpc_name:
            raise ValueError("vpc_name is required.")
        if not peer_vpc_name:
            raise ValueError("peer_vpc_name is required.")
        body = {"name": name, "vpcname": vpc_name, "peervpcname": peer_vpc_name}
        return self._transport.delete(f"fabrics/{fabric_name}/vpcpeering", json_body=body)
