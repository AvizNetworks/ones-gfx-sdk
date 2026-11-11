from __future__ import annotations


class ONESError(Exception):
    pass


class TransportError(ONESError):
    pass


class APIError(ONESError):
    def __init__(self, message: str, status_code: int | None = None, response_body: object = None):
        super().__init__(message)
        self.status_code = status_code
        self.response_body = response_body


class AuthenticationError(APIError):
    pass


class BadRequestError(APIError):
    pass


class NotFoundError(APIError):
    pass


class ConflictError(APIError):
    pass


class ServerError(APIError):
    pass


class OperationFailedError(ONESError):
    def __init__(self, operation_id: str, error_message: str | None, operation: object = None):
        msg = f"Operation {operation_id} failed: {error_message or '<no message>'}"
        super().__init__(msg)
        self.operation_id = operation_id
        self.error_message = error_message
        self.operation = operation


class ConfigurationFailedError(ONESError):
    pass


class OperationNotDoneError(ONESError):
    """Raised by result-retrieval helpers (e.g. get_topology_audit_result) when
    the operation is still PENDING or RUNNING. Not a failure - poll again
    later."""

    def __init__(self, operation_id: str, status: str | None):
        super().__init__(f"Operation {operation_id} has not finished yet (status={status}).")
        self.operation_id = operation_id
        self.status = status
