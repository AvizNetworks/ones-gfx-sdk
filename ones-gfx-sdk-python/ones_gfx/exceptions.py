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


class FabricValidationError(ONESError):
    """deploy_fabric: switch or server validation reported a failed device.

    Raised before anything is configured. ``kind`` is "switch" or "server";
    ``failures`` holds one "hostname (ip): reason" line per failed device and
    ``results`` the full validation response.
    """

    def __init__(self, kind: str, failures: list[str], results: list[dict]):
        super().__init__(f"{kind} validation failed:\n- " + "\n- ".join(failures))
        self.kind = kind
        self.failures = failures
        self.results = results
