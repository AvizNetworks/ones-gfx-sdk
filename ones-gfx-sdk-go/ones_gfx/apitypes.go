package ones_gfx

// apitypes.go holds request/response types shared by more than one resource
// file, mirroring ones_gfx/_types.py in the Python SDK. Types used by a
// single resource are defined in that resource's file instead.

// GpuAction is the ADD/DELETE discriminator shared by GPU allocation,
// tenant-update and NMX-C domain endpoints (Java Helper/Enums.java).
type GpuAction = GPUOperation

const (
	GpuActionAdd    = OperationAdd
	GpuActionDelete = OperationDelete
)

// OperationAccepted is the 202 Accepted body of async-capable endpoints
// (ApiResponse.success(msg, asyncData) with the envelope unwrapped).
// OperationType appears only on GPU allocate/deallocate operations;
// WebhookRegistered is true only in ASYNC_WEBHOOK mode.
type OperationAccepted struct {
	OperationID       string `json:"operationId,omitempty"`
	Status            string `json:"status,omitempty"`
	OperationType     string `json:"operationType,omitempty"`
	WebhookRegistered bool   `json:"webhookRegistered,omitempty"`
}
