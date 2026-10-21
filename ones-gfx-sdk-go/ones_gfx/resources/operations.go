package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// OperationsResource handles async operation polling and webhook status.
type OperationsResource struct {
	transport *ones_gfx.Transport
}

// NewOperationsResource constructs an OperationsResource.
func NewOperationsResource(transport *ones_gfx.Transport) *OperationsResource {
	return &OperationsResource{transport: transport}
}

// OperationStatusItem mirrors Models/OperationStatus.java JSON
// (GET /operations/{operationId}).
type OperationStatusItem struct {
	ID             string `json:"id,omitempty"`
	Type           string `json:"type,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	CompletedAt    string `json:"completedAt,omitempty"`
	Progress       *int   `json:"progress,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	Result         string `json:"result,omitempty"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
	CachedResponse string `json:"cachedResponse,omitempty"`
	HTTPStatusCode *int   `json:"httpStatusCode,omitempty"`
	TenantName     string `json:"tenantName,omitempty"`
	FabricName     string `json:"fabricName,omitempty"`
}

// Get fetches the current state of an async operation.
// Maps to GET /operations/{operationId}.
func (r *OperationsResource) Get(ctx context.Context, operationID string) (*OperationStatusItem, error) {
	return ones_gfx.Call[OperationStatusItem](r.transport, "GET", "operations/"+operationID, nil, ones_gfx.OperationModeSynchronous, nil)
}
