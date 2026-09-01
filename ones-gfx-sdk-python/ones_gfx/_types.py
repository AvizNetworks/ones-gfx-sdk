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
    "BootstrapsubnetItem",
    "ConnectionItem",
    "ControllerVersion",
    "DeviceConfigRestore",
    "DeviceDetail",
    "DeviceItem",
    "DevicefactsItem",
    "FabricInventoryUpdateItem",
    "FabricItem",
    "FabricSimItem",
    "GpuAction",
    "GpuItem",
    "GpuServerInfo",
    "InventoryItem",
    "InventoryRecord",
    "NetOpsAction",
    "NmxcDomain",
    "OperationAccepted",
    "ParametersItem",
    "QosItem",
    "QosmappingItem",
    "RMAInfoItem",
    "SchedulerItem",
    "StatusItem",
    "SuidMap",
    "SvimapperItem",
    "UploadFileResult",
    "VlanvnimapperItem",
]

# ── Enums (Helper/Enums.java) ────────────────────────────────────────────────

GpuAction = Literal["ADD", "DELETE"]
NetOpsAction = Literal["vlan", "portchange", "vlanmember"]

# ── GPU allocation (Helper/GpuServerInfo.java, GpuStatusUpdate.suid) ────────


class GpuServerInfo(TypedDict):
    serverName: str
    shared: bool | None


# Java: TreeMap<Integer, Map<String, DeviceGpus{gpus: List<String>}>>
# JSON object keys are strings; Jackson coerces them to Integer server-side.
SuidMap = dict[str, dict[str, list[str]]]

# ── Device facts & inventory ────────────────────────────────────────────────


class DeviceDetail(TypedDict):
    """Helper/DeviceDetail.java"""

    ip: str
    user: str
    password: str


class DeviceConfigRestore(TypedDict):
    """Helper/DeviceConfigRestore.java"""

    ip: str
    timestamp: str  # format: ddMMyyyyHHmmss


class InventoryItem(TypedDict):
    """Models/Inventory.java"""

    ipAddress: str
    username: str
    password: str
    hostname: str
    deviceType: str
    deviceRole: str
    sku: str
    interfaceData: str
    status: str
    fabricName: str
    executeConfig: str
    ufmSystemName: str


class FabricInventoryUpdateItem(TypedDict):
    """Helper/FabricInventoryUpdate.java.

    fabricName/hostname always required; ipAddress/username/password required
    unless configure == "EXECUTE_CONFIG_NO" (case-insensitive).
    """

    fabricName: str
    hostname: str
    ipAddress: str
    username: str
    password: str
    configure: str


# ── RMA (Models/RMAInfo.java) ───────────────────────────────────────────────


class RMAInfoItem(TypedDict):
    replacingDeviceMac: str
    replacingDeviceIp: str
    replacingDeviceSerial: str
    replacetoDeviceMac: str
    replacetoDeviceIp: str
    replacetoDeviceSerial: str
    fabricId: str
    state: str
    markedTime: str
    ticketId: str
    backupSelected: str
    scheduledTime: str
    isScheduledForLater: bool
    triggertime: str


# ── Bootstrap (Models/Bootstrapinfo.java, Models/Bootstrapsubnet.java) ─────


class BootstrapsubnetItem(TypedDict):
    subnet: str
    netmask: str
    gateway: str


# ── NMX-C domains (FabricManager nmxc endpoints) ───────────────────────────


class NmxcDomain(TypedDict):
    host: str
    port: int | None
    useTls: bool | None


# ── Intent sub-objects (Models/*.java) ──────────────────────────────────────


class DevicefactsItem(TypedDict):
    ipAddress: str
    username: str
    password: str
    iv: str
    agentVersion: str
    nosVersion: str
    config: str
    isactive: bool


class ConnectionItem(TypedDict):
    localIp: str
    remoteIp: str
    localIntName: str
    remoteIntName: str
    localSwitchId: str
    remoteSwitchId: str
    remoteAsn: str
    localAsn: str
    networkAddress: str
    isMclagLink: bool
    isHostLink: bool
    localPOIp: str
    remotePOIp: str
    localPO: str
    remotePO: str
    vlans: str
    mode: str
    mclagPOGroup: str
    isMhEnabled: bool


class StatusItem(TypedDict):
    intentName: str
    subIntentName: str
    config_status: str
    verification_status: str
    config_logs: str
    verification_logs: str
    isactive: bool


class SvimapperItem(TypedDict):
    vlan: str
    network: str
    isactive: bool


class VlanvnimapperItem(TypedDict):
    vlan: str
    vni: str
    vrf: str
    network: str
    isIRB: bool
    isactive: bool


class DeviceItem(TypedDict, total=False):
    """Models/Device.java"""

    hostname: str
    switchId: int
    asn: str
    loopbackIp: str
    loopbackIntName: str
    isSSpine: bool
    isSpine: bool
    isLeaf: bool
    isTor: bool
    isHost: bool
    isDpu: bool
    isMclagEnabled: bool
    mclagKeepaliveVlan: str
    mclagPeerIdentifier: str
    isMclagOverL3: bool
    vtepLoopbackIdentifier: str
    vtepLoopbackInterface: str
    vtepLoopbackInterfaceIp: str
    vrfs: str
    vlans: str
    isSflowEnabled: bool
    sFlowAgent: str
    sFlowPollInterval: str
    isactive: bool
    ufmSystemName: str
    devicefacts: DevicefactsItem
    connections: list[ConnectionItem]
    statuses: list[StatusItem]
    svis: list[SvimapperItem]


class QosmappingItem(TypedDict):
    """Models/Qosmapping.java"""

    traffic_class: int
    dscp: str
    dot1p: str
    queues: int
    priority_group: int


class SchedulerItem(TypedDict):
    """Models/Scheduler.java"""

    queues: int
    profile_name: str
    weights: int


class QosItem(TypedDict):
    """Models/Qos.java"""

    dscp_tc_map: str
    dot1p_tc_map: str
    tc_pg_map: str
    tc_queue_map: str
    pfc_enabled_queues: str
    pfcwd_action: str
    pfcwd_detection_time: int
    pfcwd_restoration_time: int
    is_brs_enabled: bool
    wred_policy: str
    ecn_mode: str
    ecn_cnp_queue: int
    ecn_gmin: int
    ecn_gmax: int
    ecn_gmark: int
    ecn_rmin: int
    ecn_rmax: int
    ecn_rmark: int
    ecn_ymin: int
    ecn_ymax: int
    ecn_ymark: int
    scheduler_type: str


class ParametersItem(TypedDict):
    """Models/Parameters.java"""

    vlan: str
    vni: str
    anycastGateway: str
    anycastMac: str
    hostsPerVlan: int
    hostsPerL3Subnet: int
    networkAddress: str
    prefixLen: int
    possibleSubNetworks: int
    possibleHost: int
    logs: str
    irbVlans: str
    irbVnis: str
    routing_symmetric: bool
    isactive: bool
    vlanVniMapper: list[VlanvnimapperItem]


# ─────────────────────────────────────────────────────────────────────────────
# Response types (FabricManager @ResponseBody shapes after the SDK client
# unwraps the {"data": ...} envelope). total=False everywhere a key can be
# absent or null on the wire.
# ─────────────────────────────────────────────────────────────────────────────


class ControllerVersion(TypedDict):
    """GET /getControllerVersion — Map<String, String> build info"""

    appName: str
    appArtifactId: str
    appVersion: str
    appBuildDateTime: str


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


class FabricItem(TypedDict, total=False):
    """Models/Fabrics.java JSON"""

    id: int
    name: str
    type: str
    status: str
    description: str
    orchestrationStatus: str
    numOfSus: int
    maxNumOfSus: int
    dedicated: bool
    hybrid: bool
    isDPUFabric: bool
    startingSubnetGpu: int
    startingSubnetCpu: str
    startingSubnetTenants: str
    startingSubnetStorage: str
    simulationId: int
    intent: str
    ewTenantAware: bool
    storageTenantAware: bool
    isOnesControlled: bool
    suHostCnt: str
    isVxlanFabric: bool
    nodeType: str
    deploymentType: str
    spineEvpnConfigured: bool
    isimported: bool
    gpuScaleMode: str
    cnpq: str
    ufmUrl: str
    ufmUsername: str
    ufmPasswordEncrypted: str
    createdAt: str
    updatedAt: str


class FabricSimItem(TypedDict, total=False):
    """Models/Fabricsims.java JSON"""

    id: int
    fabricName: str
    simulationId: str
    username: str
    token: str
    orgUuid: str
    status: str
    uiLink: str
    createdAt: str
    updatedAt: str


class InventoryRecord(TypedDict, total=False):
    """Models/Inventory.java JSON"""

    id: int
    ipAddress: str
    username: str
    password: str
    hostname: str
    deviceType: str
    deviceRole: str
    sku: str
    interfaceData: str
    status: str
    fabricName: str
    executeConfig: str
    ufmSystemName: str
    createdAt: str
    updatedAt: str


class GpuItem(TypedDict, total=False):
    """Models/Gpus.java JSON"""

    id: int
    gpuPort: str
    gpuHostname: str
    gpuStatus: str
    gpuLeafInt: str
    suId: int
    leafHostname: str
    leafIpAddress: str
    tenantName: str
    fabricName: str
    config_status: str
    lastConfiguredTenant: str
    createdAt: str
    updatedAt: str


class UploadFileResult(TypedDict, total=False):
    """POST /uploadfile"""

    success: bool
    message: str
    filepath: str
    checksum: str
    filesize: int
    filename: str
    error: str


