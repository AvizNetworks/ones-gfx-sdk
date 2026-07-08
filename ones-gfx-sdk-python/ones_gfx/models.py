from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any

from .enums import ConfigStatus


def _ms_to_datetime(value: Any) -> datetime | None:
    if value is None:
        return None
    try:
        return datetime.fromtimestamp(int(value) / 1000.0, tz=timezone.utc)
    except (ValueError, TypeError, OSError):
        return None


def _parse_timestamp(value: Any) -> datetime | None:
    if value is None:
        return None
    if isinstance(value, (int, float)):
        return _ms_to_datetime(value)
    if isinstance(value, str):
        iso = value.replace("Z", "+00:00")
        try:
            return datetime.fromisoformat(iso)
        except ValueError:
            return _ms_to_datetime(value)
    return None


@dataclass
class Fabric:
    id: int
    fabric_name: str
    description: str
    num_of_sus: int
    max_num_of_sus: int
    ew_tenant_aware: bool
    ns_tenant_aware: bool
    default_storage_name: str | None = None

    @classmethod
    def from_api(cls, payload: dict) -> "Fabric":
        return cls(
            id=payload.get("id"),
            fabric_name=payload.get("fabricName"),
            description=payload.get("description", ""),
            num_of_sus=payload.get("numOfSUs", 0),
            max_num_of_sus=payload.get("maxNumOfSUs", 0),
            ew_tenant_aware=payload.get("e-wTenantAware", False),
            ns_tenant_aware=payload.get("n-sTenantAware", False),
            default_storage_name=payload.get("defaultStorageName"),
        )


@dataclass
class NMXCPartition:
    domain_name: str
    partition_id: int
    partition_name: str
    resiliency_mode: str
    status: str

    @classmethod
    def from_api(cls, payload: dict) -> "NMXCPartition":
        return cls(
            domain_name=payload.get("domainName", ""),
            partition_id=payload.get("partitionId", 0),
            partition_name=payload.get("partitionName", ""),
            resiliency_mode=payload.get("resiliencyMode", ""),
            status=payload.get("status", ""),
        )



@dataclass
class VNetInfo:
    name: str | None = None
    description: str | None = None
    ip_pool_gateway: str | None = None
    ip_pool_subnet: str | None = None

    @classmethod
    def from_api(cls, payload: dict | None) -> "VNetInfo":
        if not payload:
            return cls()
        ip_pool = payload.get("tenant IP Pool", {}) or {}
        return cls(
            name=payload.get("name"),
            description=payload.get("description"),
            ip_pool_gateway=ip_pool.get("gateway"),
            ip_pool_subnet=ip_pool.get("subnet"),
        )


@dataclass
class Tenant:
    id: int
    name: str
    description: str
    fabric_name: str
    max_gpus_allowed: int
    gpus_allocated: int
    alloted_gpus: str
    vni_id: int | None
    vlan_id: int | None
    vlan_vni_id: int | None
    vlan_subnet_cpu: str | None
    vlan_subnet_storage: str | None
    config_status: ConfigStatus
    networks: list[str] = field(default_factory=list)
    tags: list[dict] = field(default_factory=list)
    vnets: VNetInfo = field(default_factory=VNetInfo)
    shared: bool = False
    transit_vlan_id: int | None = None
    ns_vni_id: int | None = None
    created_at: datetime | None = None
    updated_at: datetime | None = None
    # UFM-specific fields
    pkey_value: str | None = None
    ufm_allocated_ports: int = 0
    ufm_allocated_servers: int = 0
    # NMXC-specific fields
    nmxc_gpus_allocated: int = 0
    nmxc_servers_allocated: int = 0
    nmxc_partitions: list[NMXCPartition] = field(default_factory=list)

    @classmethod
    def from_api(cls, payload: dict) -> "Tenant":
        return cls(
            id=payload.get("id"),
            name=payload.get("name"),
            description=payload.get("description", ""),
            fabric_name=payload.get("fabricName"),
            max_gpus_allowed=payload.get("maxGpusAllowed", 0),
            gpus_allocated=payload.get("gpusAllocated", 0) or 0,
            alloted_gpus=payload.get("allotedGpus", "") or "",
            vni_id=payload.get("vniId"),
            vlan_id=payload.get("vlanId"),
            vlan_vni_id=payload.get("vlanVniId"),
            vlan_subnet_cpu=payload.get("vlanSubnetCPU"),
            vlan_subnet_storage=payload.get("vlanSubnetStorage"),
            config_status=ConfigStatus.from_api(payload.get("config_status")),
            networks=payload.get("networks", []) or [],
            tags=payload.get("tags", []) or [],
            vnets=VNetInfo.from_api(payload.get("vnets")),
            shared=payload.get("shared", False),
            transit_vlan_id=payload.get("transitVlanId"),
            ns_vni_id=payload.get("nsVniId"),
            created_at=_parse_timestamp(payload.get("createdAt")),
            updated_at=_parse_timestamp(payload.get("updatedAt")),
            pkey_value=payload.get("pkeyValue"),
            ufm_allocated_ports=payload.get("ufmAllocatedPorts", 0),
            ufm_allocated_servers=payload.get("ufmAllocatedServers", 0),
            nmxc_gpus_allocated=payload.get("nmxcGpusAllocated", 0),
            nmxc_servers_allocated=payload.get("nmxcServersAllocated", 0),
            nmxc_partitions=[NMXCPartition.from_api(p) for p in payload.get("nmxcPartitions", []) or []],
        )

    @property
    def alloted_servers(self) -> list[str]:
        return [s for s in self.alloted_gpus.split(",") if s]

    @property
    def is_unlimited_gpus(self) -> bool:
        return self.max_gpus_allowed == -1


@dataclass
class Operation:
    id: str
    type: str | None = None
    status: str | None = None
    progress: int | None = None
    result: Any = None
    error_message: str | None = None
    http_status_code: int | None = None
    tenant_name: str | None = None
    fabric_name: str | None = None
    operation_type: str | None = None
    webhook_registered: bool = False
    created_at: datetime | None = None
    updated_at: datetime | None = None
    completed_at: datetime | None = None

    @classmethod
    def from_submission(cls, data: dict) -> "Operation":
        return cls(
            id=data.get("operationId"),
            status=data.get("status"),
            operation_type=data.get("operationType"),
            webhook_registered=bool(data.get("webhookRegistered", False)),
        )

    @classmethod
    def from_polled(cls, payload: dict) -> "Operation":
        import json
        raw_result = payload.get("result")
        parsed_result: Any = raw_result
        if isinstance(raw_result, str):
            try:
                parsed_result = json.loads(raw_result)
            except (json.JSONDecodeError, ValueError):
                parsed_result = raw_result
        return cls(
            id=payload.get("id"),
            type=payload.get("type"),
            status=payload.get("status"),
            progress=payload.get("progress"),
            result=parsed_result,
            error_message=payload.get("errorMessage"),
            http_status_code=payload.get("httpStatusCode"),
            tenant_name=payload.get("tenantName"),
            fabric_name=payload.get("fabricName"),
            created_at=_parse_timestamp(payload.get("createdAt")),
            updated_at=_parse_timestamp(payload.get("updatedAt")),
            completed_at=_parse_timestamp(payload.get("completedAt")),
        )

    @property
    def is_done(self) -> bool:
        return self.status in ("SUCCESS", "FAILURE")

    @property
    def is_success(self) -> bool:
        return self.status == "SUCCESS"

    @property
    def is_failure(self) -> bool:
        return self.status == "FAILURE"
