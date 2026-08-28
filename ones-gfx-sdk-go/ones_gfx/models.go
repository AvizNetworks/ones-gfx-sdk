package ones_gfx

import (
	"encoding/json"
	"strings"
	"time"
)

// Operation represents an asynchronous job submitted to the API.
type Operation struct {
	ID                string          `json:"id"`
	Type              string          `json:"type,omitempty"`
	Status            OperationStatus `json:"status,omitempty"`
	Progress          int             `json:"progress,omitempty"`
	Result            interface{}     `json:"result,omitempty"` // Parsed or raw string
	ErrorMessage      string          `json:"errorMessage,omitempty"`
	HTTPStatusCode    *int            `json:"httpStatusCode,omitempty"`
	TenantName        string          `json:"tenantName,omitempty"`
	FabricName        string          `json:"fabricName,omitempty"`
	OperationType     string          `json:"operationType,omitempty"` // GPU_ALLOCATE, GPU_DEALLOCATE
	WebhookRegistered bool            `json:"webhookRegistered,omitempty"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	CompletedAt       time.Time       `json:"-"`
	CreatedAtRaw      interface{}     `json:"createdAt,omitempty"`
	UpdatedAtRaw      interface{}     `json:"updatedAt,omitempty"`
	CompletedAtRaw    interface{}     `json:"completedAt,omitempty"`
}

// UnmarshalJSON parses Operation and handles result field (may be JSON string).
func (o *Operation) UnmarshalJSON(data []byte) error {
	type Alias Operation
	aux := &struct {
		ResultRaw   interface{} `json:"result"`
		OperationID string      `json:"operationId"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if o.ID == "" && aux.OperationID != "" {
		o.ID = aux.OperationID
	}
	// Parse result (may be JSON-encoded string or plain string)
	if aux.ResultRaw != nil {
		if str, ok := aux.ResultRaw.(string); ok && str != "" {
			var parsed interface{}
			if err := json.Unmarshal([]byte(str), &parsed); err == nil {
				o.Result = parsed
			} else {
				o.Result = str
			}
		} else {
			o.Result = aux.ResultRaw
		}
	}
	// Parse timestamps
	o.CreatedAt = parseTimestamp(o.CreatedAtRaw)
	o.UpdatedAt = parseTimestamp(o.UpdatedAtRaw)
	o.CompletedAt = parseTimestamp(o.CompletedAtRaw)
	return nil
}

// IsDone returns true if the operation has reached a terminal state.
func (o *Operation) IsDone() bool {
	return o.Status.IsTerminal()
}

// IsSuccess returns true if the operation completed successfully.
func (o *Operation) IsSuccess() bool {
	return o.Status == OperationStatusSuccess
}

// IsFailure returns true if the operation failed.
func (o *Operation) IsFailure() bool {
	return o.Status == OperationStatusFailure
}

// parseTimestamp handles both ms-since-epoch and ISO-8601 timestamps.
func parseTimestamp(raw interface{}) time.Time {
	if raw == nil {
		return time.Time{}
	}
	// Try as float64 (ms-since-epoch)
	if ms, ok := raw.(float64); ok {
		return time.Unix(0, int64(ms)*int64(time.Millisecond))
	}
	// Try as string (ISO-8601)
	if str, ok := raw.(string); ok && str != "" {
		// Normalize trailing Z to +00:00 for time.Parse
		normalized := strings.Replace(str, "Z", "+00:00", 1)
		if t, err := time.Parse(time.RFC3339Nano, normalized); err == nil {
			return t
		}
		// Fallback: try as epoch string
		if t, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", str); err == nil {
			return t
		}
	}
	return time.Time{}
}
