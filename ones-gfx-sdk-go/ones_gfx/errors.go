package ones_gfx

import "fmt"

// ONESError is the interface implemented by all SDK errors. Partners can check
// if an error originated from the SDK with errors.As(&ones.ONESError{}).
type ONESError interface {
	error
	// ONESError is a marker method to identify SDK errors.
	ONESError()
}

// TransportError represents network-level failures: connection refused, TLS
// error, timeout, DNS resolution failure.
type TransportError struct {
	Message string
	Cause   error
}

func (e *TransportError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("transport error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("transport error: %s", e.Message)
}

func (e *TransportError) ONESError() {}

func (e *TransportError) Unwrap() error {
	return e.Cause
}

// APIError is the base for errors returned by the ONES API (i.e., the request
// reached the server and got an HTTP response, but the response indicates failure).
type APIError struct {
	Message      string
	StatusCode   int
	ResponseBody interface{} // Decoded JSON or raw text
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) ONESError() {}

// AuthenticationError represents HTTP 401/403 or refresh token failure.
type AuthenticationError struct {
	APIError
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error (HTTP %d): %s", e.StatusCode, e.Message)
}

// BadRequestError represents HTTP 400 (invalid input).
type BadRequestError struct {
	APIError
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request (HTTP %d): %s", e.StatusCode, e.Message)
}

// NotFoundError represents HTTP 404 (fabric, tenant, or operation not found).
type NotFoundError struct {
	APIError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found (HTTP %d): %s", e.StatusCode, e.Message)
}

// ConflictError represents HTTP 409 (tenant already exists, GPUs still allocated, etc.).
type ConflictError struct {
	APIError
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict (HTTP %d): %s", e.StatusCode, e.Message)
}

// ServerError represents HTTP 5xx (ONES server internal failure).
type ServerError struct {
	APIError
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("server error (HTTP %d): %s", e.StatusCode, e.Message)
}

// OperationFailedError is raised when an async operation completes with
// status FAILURE. The operation's error_message field contains the reason.
type OperationFailedError struct {
	OperationID  string
	ErrorMessage string
	Operation    *Operation // Full operation object for inspection
}

func (e *OperationFailedError) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("operation %s failed: %s", e.OperationID, e.ErrorMessage)
	}
	return fmt.Sprintf("operation %s failed (no error message)", e.OperationID)
}

func (e *OperationFailedError) ONESError() {}

// ConfigurationFailedError is reserved for future helpers that watch
// config_status transitions. Not used in v1.0.0
type ConfigurationFailedError struct {
	TenantName   string
	FabricName   string
	ConfigStatus ConfigStatus
}

func (e *ConfigurationFailedError) Error() string {
	return fmt.Sprintf(
		"tenant %s in fabric %s configuration failed (status: %s)",
		e.TenantName, e.FabricName, e.ConfigStatus,
	)
}

func (e *ConfigurationFailedError) ONESError() {}
