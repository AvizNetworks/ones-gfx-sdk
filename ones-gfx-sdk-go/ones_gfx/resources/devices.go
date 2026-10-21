package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// DevicesResource covers device management: facts, versions, upgrades,
// aliveness and UFM credential validation.
type DevicesResource struct {
	transport *ones_gfx.Transport
}

// NewDevicesResource constructs a DevicesResource.
func NewDevicesResource(transport *ones_gfx.Transport) *DevicesResource {
	return &DevicesResource{transport: transport}
}

// UfmCredsResult is the response of ValidateUFMCreds (POST /ValidateUfmCreds).
type UfmCredsResult struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ValidateUFMCreds checks UFM URL/username/password connectivity.
// Maps to POST /ValidateUfmCreds.
func (r *DevicesResource) ValidateUFMCreds(ctx context.Context, ufmURL, username, password string) (*UfmCredsResult, error) {
	body := struct {
		UfmURL   string `json:"ufmUrl"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{ufmURL, username, password}
	return ones_gfx.Call[UfmCredsResult](r.transport, "POST", "ValidateUfmCreds", body, ones_gfx.OperationModeSynchronous, nil)
}
