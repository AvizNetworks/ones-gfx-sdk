"""
ONES Spectrum-X Python SDK.

Public API:

    The 20 Fabric Manager API functions documented in the ONES Spectrum-X API
    Operations guide are re-exported here (create_tenant, get_fabrics,
    update_tenant, ...), mirroring ones_gfx.apis.

    ONESClient         - main entry point
    JWTAuth            - JWT-based authentication
    OperationMode      - sync / async-poll / async-webhook
    OperationStatus    - PENDING / RUNNING / SUCCESS / FAILURE
    ConfigStatus       - tenant configuration state enum
    UNLIMITED_GPUS     - sentinel value for max_gpus_allowed=-1

    Models:
        Fabric, Tenant, VNetInfo, Operation, NMXCPartition

    Exceptions:
        ONESError (base), AuthenticationError, BadRequestError,
        NotFoundError, ConflictError, ServerError, TransportError,
        OperationFailedError, APIError
"""

from ._version import __version__
from .auth import AuthBase, JWTAuth
from .client import Client, NotAuthenticatedError, ONESClient
from .enums import (
    UNLIMITED_GPUS,
    ConfigStatus,
    OperationMode,
    OperationStatus,
)
from .exceptions import (
    APIError,
    AuthenticationError,
    BadRequestError,
    ConfigurationFailedError,
    ConflictError,
    NotFoundError,
    ONESError,
    OperationFailedError,
    ServerError,
    TransportError,
)
from .models import Fabric, NMXCPartition, Operation, Tenant, VNetInfo
from .apis import (
    create_tenant,
    list_tenants,
    get_tenant,
    update_tenant,
    delete_tenant,
    get_operation,
    get_fabrics,
    get_inventory_ports,
    get_inventory_hosts,
    inventory_sync,
    validate_ufm_creds,
    update_nmxc_domains,
    probe_nmxc_domains,
    get_nmxc_inventory,
    create_vpc_peering,
    modify_gpu_allocations,
    assign_gpu_ports,
    get_available_servers,
    reset_nmxc_domain,
    factory_reset_nmxc_domain,
)

__all__ = [
    "__version__",
    "Client",
    "NotAuthenticatedError",
    "ONESClient",
    "AuthBase",
    "JWTAuth",
    "OperationMode",
    "OperationStatus",
    "ConfigStatus",
    "UNLIMITED_GPUS",
    "Fabric",
    "Tenant",
    "VNetInfo",
    "Operation",
    "NMXCPartition",
    "ONESError",
    "APIError",
    "AuthenticationError",
    "BadRequestError",
    "NotFoundError",
    "ConflictError",
    "ServerError",
    "TransportError",
    "OperationFailedError",
    "ConfigurationFailedError",
    "create_tenant",
    "list_tenants",
    "get_tenant",
    "update_tenant",
    "delete_tenant",
    "get_operation",
    "get_fabrics",
    "get_inventory_ports",
    "get_inventory_hosts",
    "inventory_sync",
    "validate_ufm_creds",
    "update_nmxc_domains",
    "probe_nmxc_domains",
    "get_nmxc_inventory",
    "create_vpc_peering",
    "modify_gpu_allocations",
    "assign_gpu_ports",
    "get_available_servers",
    "reset_nmxc_domain",
    "factory_reset_nmxc_domain",
]
