package ones_gfx

import (
	"fmt"
	"math"
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

	// OperationStatusFailed is the spelling Fabric Manager actually sends for
	// a failed operation (terminal). It means the same as
	// OperationStatusFailure; use IsFailure to test for either.
	OperationStatusFailed OperationStatus = "FAILED"
)

// IsTerminal returns true if the operation has reached a final state.
func (s OperationStatus) IsTerminal() bool {
	return s == OperationStatusSuccess || s.IsFailure()
}

// IsFailure reports whether the operation ended with an error, under either
// spelling (FAILURE or FAILED).
func (s OperationStatus) IsFailure() bool {
	return s == OperationStatusFailure || s == OperationStatusFailed
}

// String returns the string representation of the status.
func (s OperationStatus) String() string {
	return string(s)
}

// Valid reports whether s is one of the states FM reports.
func (s OperationStatus) Valid() bool {
	switch s {
	case OperationStatusPending, OperationStatusRunning, OperationStatusSuccess, OperationStatusFailure, OperationStatusFailed:
		return true
	}
	return false
}

// ParseOperationStatus converts external input into an OperationStatus. It is
// case-insensitive; anything other than PENDING, RUNNING, SUCCESS, FAILURE or
// FAILED is an error.
func ParseOperationStatus(s string) (OperationStatus, error) {
	st := OperationStatus(strings.ToUpper(strings.TrimSpace(s)))
	if !st.Valid() {
		return "", fmt.Errorf("invalid operation status %q: want one of %s, %s, %s, %s, %s",
			s, OperationStatusPending, OperationStatusRunning, OperationStatusSuccess, OperationStatusFailure, OperationStatusFailed)
	}
	return st, nil
}

// ConfigStatus is a device's configuration state, as Fabric Manager reports it
// in the config_status field.
//
// The numbering is positional, not a progression, and reads wrong at a glance:
//
//	0  FAIL         configuration failed
//	1  PASS         configuration succeeded
//	2  NOT_STARTED  configuration has not begun   <- NOT a success
//	3  IN_PROGRESS  configuration is underway
//	4  BROWN_FIELD  skipped: the fabric deploys as brownfield
//	5  SKIPPED      skipped: pre-validation or minimal-config excluded the device
//
// The trap is 2. Because it sits above PASS it looks like a later, stronger
// state, and a "status >= 1 means fine" check silently treats a device that was
// never configured as configured. Compare against the named constants and never
// against the ordinal, and use IsSuccess rather than ordering.
//
// 4 and 5 both mean "deliberately not configured" and are terminal; neither is
// a failure and neither is a success.
//
// Fabric Manager sends these as JSON *strings* ("1", not 1). UnmarshalJSON
// accepts either; ParseConfigStatus is the strict entry point for values that
// arrive from somewhere else.
type ConfigStatus int

const (
	// ConfigStatusFail indicates configuration failed.
	ConfigStatusFail ConfigStatus = 0

	// ConfigStatusPass indicates configuration succeeded. This is the only
	// value that means the device is configured.
	ConfigStatusPass ConfigStatus = 1

	// ConfigStatusNotStarted indicates configuration has not begun. Despite
	// ordering above ConfigStatusPass, this is not a success state.
	ConfigStatusNotStarted ConfigStatus = 2

	// ConfigStatusInProgress indicates configuration is underway.
	ConfigStatusInProgress ConfigStatus = 3

	// ConfigStatusBrownField indicates the device was skipped because the
	// fabric is deployed as brownfield (AsyncConfig sets this instead of
	// pushing config).
	ConfigStatusBrownField ConfigStatus = 4

	// ConfigStatusSkipped indicates the device was skipped by pre-validation
	// or by the minimal-config path.
	ConfigStatusSkipped ConfigStatus = 5
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
	if v := ConfigStatus(n); v.IsKnown() {
		*c = v
	} else {
		*c = ConfigStatusNotStarted
	}
	return nil
}

// IsKnown reports whether c is a value Fabric Manager actually emits.
func (c ConfigStatus) IsKnown() bool {
	switch c {
	case ConfigStatusFail, ConfigStatusPass, ConfigStatusNotStarted,
		ConfigStatusInProgress, ConfigStatusBrownField, ConfigStatusSkipped:
		return true
	}
	return false
}

// IsSuccess reports whether the device is configured.
//
// Only ConfigStatusPass qualifies. BrownField and Skipped are deliberate
// non-configuration, not success, so a caller gating on "is it configured"
// must not treat them as such.
func (c ConfigStatus) IsSuccess() bool { return c == ConfigStatusPass }

// IsTerminal reports whether the status can still change on its own.
//
// NotStarted and InProgress are the only non-terminal values: the first is
// waiting to run, the second is running. Fail, Pass, BrownField and Skipped
// are all settled, which is why polling should stop on any of them rather
// than only on Pass.
func (c ConfigStatus) IsTerminal() bool {
	switch c {
	case ConfigStatusFail, ConfigStatusPass, ConfigStatusBrownField, ConfigStatusSkipped:
		return true
	}
	return false
}

// IsSkipped reports whether the device was deliberately not configured.
func (c ConfigStatus) IsSkipped() bool {
	return c == ConfigStatusBrownField || c == ConfigStatusSkipped
}

// ParseConfigStatus converts an external value into a ConfigStatus, rejecting
// anything Fabric Manager does not emit.
//
// This is the strict counterpart to UnmarshalJSON. Decoding a response stays
// lenient, because one unreadable status should not discard an otherwise good
// payload; a value arriving from configuration, a CLI flag or another system
// gets no such benefit of the doubt and fails loudly instead.
//
// Accepts an int, a numeric string ("1"), or the canonical name ("PASS").
func ParseConfigStatus(v interface{}) (ConfigStatus, error) {
	switch x := v.(type) {
	case ConfigStatus:
		if x.IsKnown() {
			return x, nil
		}
		return ConfigStatusNotStarted, fmt.Errorf("unknown config status %d", int(x))
	case int:
		return ParseConfigStatus(ConfigStatus(x))
	case int64:
		return ParseConfigStatus(ConfigStatus(x))
	case float64:
		// JSON numbers arrive as float64. Truncating would quietly turn 1.5
		// into PASS, so a non-integral value is rejected rather than rounded.
		if x != math.Trunc(x) {
			return ConfigStatusNotStarted, fmt.Errorf("config status %v is not a whole number", x)
		}
		return ParseConfigStatus(ConfigStatus(int(x)))
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return ConfigStatusNotStarted, fmt.Errorf("empty config status")
		}
		if n, err := strconv.Atoi(s); err == nil {
			return ParseConfigStatus(ConfigStatus(n))
		}
		switch strings.ToUpper(strings.ReplaceAll(s, "-", "_")) {
		case "FAIL":
			return ConfigStatusFail, nil
		case "PASS":
			return ConfigStatusPass, nil
		case "NOT_STARTED":
			return ConfigStatusNotStarted, nil
		case "IN_PROGRESS":
			return ConfigStatusInProgress, nil
		case "BROWN_FIELD", "BROWNFIELD":
			return ConfigStatusBrownField, nil
		case "SKIPPED":
			return ConfigStatusSkipped, nil
		}
		return ConfigStatusNotStarted, fmt.Errorf("unknown config status %q", x)
	}
	return ConfigStatusNotStarted, fmt.Errorf("cannot read config status from %T", v)
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
	case ConfigStatusBrownField:
		return "BROWN_FIELD"
	case ConfigStatusSkipped:
		return "SKIPPED"
	default:
		return "UNKNOWN"
	}
}

// ParseConfigStatusLenient converts an API config_status string into a
// ConfigStatus, falling back to ConfigStatusNotStarted for anything it does
// not recognise.
//
// Prefer ParseConfigStatus, which reports unknown values instead of hiding
// them. This exists for callers that genuinely want a total function.
func ParseConfigStatusLenient(s string) ConfigStatus {
	if v, err := ParseConfigStatus(s); err == nil {
		return v
	}
	return ConfigStatusNotStarted
}

// UnlimitedGPUs is the sentinel value for max_gpus_allowed meaning no cap.
// Pass this to CreateTenantRequest.MaxGPUsAllowed to indicate unlimited quota.
const UnlimitedGPUs = -1
