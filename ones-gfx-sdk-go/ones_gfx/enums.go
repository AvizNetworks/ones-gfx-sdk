package ones_gfx

import (
	"strconv"
	"strings"
)

// OperationMode determines how the SDK submits mutating requests to the API.
type OperationMode string

const (
	// OperationModeSynchronous blocks until the server finishes. Returns the
	// parsed result directly (e.g., Tenant, nil error).
	OperationModeSynchronous OperationMode = "sync"

	// OperationModeAsyncPoll returns immediately with an Operation handle.
	// Caller polls client.Operations.Get(op.ID) until done.
	OperationModeAsyncPoll OperationMode = "async_poll"

	// OperationModeAsyncWebhook returns immediately with an Operation handle.
	// Server POSTs result to the webhook URL on completion.
	OperationModeAsyncWebhook OperationMode = "async_webhook"
)

// String returns the string representation of the mode.
func (m OperationMode) String() string {
	return string(m)
}

// OperationStatus represents the lifecycle state of an async operation.
type OperationStatus string

const (
	// OperationStatusPending means accepted but not yet started.
	OperationStatusPending OperationStatus = "PENDING"

	// OperationStatusRunning means in progress on the server.
	OperationStatusRunning OperationStatus = "RUNNING"

	// OperationStatusSuccess means completed successfully (terminal).
	OperationStatusSuccess OperationStatus = "SUCCESS"

	// OperationStatusFailure means completed with error (terminal).
	OperationStatusFailure OperationStatus = "FAILURE"
)

// IsTerminal returns true if the operation has reached a final state.
func (s OperationStatus) IsTerminal() bool {
	return s == OperationStatusSuccess || s == OperationStatusFailure
}

// String returns the string representation of the status.
func (s OperationStatus) String() string {
	return string(s)
}

// ConfigStatus represents tenant configuration state as returned by the API.
type ConfigStatus int

const (
	// ConfigStatusFail indicates configuration failed.
	ConfigStatusFail ConfigStatus = 0

	// ConfigStatusPass indicates configuration succeeded.
	ConfigStatusPass ConfigStatus = 1

	// ConfigStatusNotStarted indicates configuration has not begun.
	ConfigStatusNotStarted ConfigStatus = 2

	// ConfigStatusInProgress indicates configuration is underway.
	ConfigStatusInProgress ConfigStatus = 3
)

// UnmarshalJSON decodes a config status from either spelling Fabric Manager
// uses: the string "1" that the tenant payload carries, or a bare number.
//
// An unrecognised value decodes to ConfigStatusNotStarted rather than failing
// the whole response, matching ConfigStatus.from_api in the Python SDK: one
// malformed status should not make an otherwise good tenant unreadable.
func (c *ConfigStatus) UnmarshalJSON(data []byte) error {
	s := strings.Trim(strings.TrimSpace(string(data)), `"`)
	if s == "" || s == "null" {
		*c = ConfigStatusNotStarted
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		*c = ConfigStatusNotStarted
		return nil
	}
	switch ConfigStatus(n) {
	case ConfigStatusFail, ConfigStatusPass, ConfigStatusNotStarted, ConfigStatusInProgress:
		*c = ConfigStatus(n)
	default:
		*c = ConfigStatusNotStarted
	}
	return nil
}

// String returns the human-readable name of the status.
func (c ConfigStatus) String() string {
	switch c {
	case ConfigStatusFail:
		return "FAIL"
	case ConfigStatusPass:
		return "PASS"
	case ConfigStatusNotStarted:
		return "NOT_STARTED"
	case ConfigStatusInProgress:
		return "IN_PROGRESS"
	default:
		return "UNKNOWN"
	}
}

// ParseConfigStatus parses the API-returned config_status string (which
// arrives as "0", "1", "2", "3") into a ConfigStatus. Returns
// ConfigStatusNotStarted for unknown values.
func ParseConfigStatus(s string) ConfigStatus {
	switch s {
	case "0":
		return ConfigStatusFail
	case "1":
		return ConfigStatusPass
	case "2":
		return ConfigStatusNotStarted
	case "3":
		return ConfigStatusInProgress
	default:
		return ConfigStatusNotStarted
	}
}

// UnlimitedGPUs is the sentinel value for max_gpus_allowed meaning no cap.
// Pass this to CreateTenantRequest.MaxGPUsAllowed to indicate unlimited quota.
const UnlimitedGPUs = -1
