"""
ONES Spectrum-X Python SDK.

Public API:

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
from .client import ONESClient
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

__all__ = [
    "__version__",
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
]
