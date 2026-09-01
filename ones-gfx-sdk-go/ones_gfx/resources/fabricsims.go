package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FabricsimsResource covers fabric simulation endpoints.
type FabricsimsResource struct {
	transport *ones_gfx.Transport
}

// NewFabricsimsResource constructs a FabricsimsResource.
func NewFabricsimsResource(transport *ones_gfx.Transport) *FabricsimsResource {
	return &FabricsimsResource{transport: transport}
}

// FabricSimItem mirrors Models/Fabricsims.java JSON.
type FabricSimItem struct {
	ID           int    `json:"id,omitempty"`
	FabricName   string `json:"fabricName,omitempty"`
	SimulationID string `json:"simulationId,omitempty"`
	Username     string `json:"username,omitempty"`
	Token        string `json:"token,omitempty"`
	OrgUUID      string `json:"orgUuid,omitempty"`
	Status       string `json:"status,omitempty"`
	UILink       string `json:"uiLink,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}

// Create registers a simulation. Maps to POST /addFabricSim.
func (r *FabricsimsResource) Create(ctx context.Context, id *int, fabricName, simulationID, username, token, orgUUID, status, uiLink *string) (string, error) {
	body := struct {
		ID           *int    `json:"id,omitempty"`
		FabricName   *string `json:"fabricName,omitempty"`
		SimulationID *string `json:"simulationId,omitempty"`
		Username     *string `json:"username,omitempty"`
		Token        *string `json:"token,omitempty"`
		OrgUUID      *string `json:"orgUuid,omitempty"`
		Status       *string `json:"status,omitempty"`
		UILink       *string `json:"uiLink,omitempty"`
	}{id, fabricName, simulationID, username, token, orgUUID, status, uiLink}
	res, err := ones_gfx.Call[string](r.transport, "POST", "addFabricSim", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// List returns all fabric simulations. Maps to GET /getAllFabricSims.
func (r *FabricsimsResource) List(ctx context.Context) ([]FabricSimItem, error) {
	res, err := ones_gfx.Call[[]FabricSimItem](r.transport, "GET", "getAllFabricSims", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Get returns one simulation by fabric name. Maps to GET /getFabricSimByName/{name}.
func (r *FabricsimsResource) Get(ctx context.Context, name string) (*FabricSimItem, error) {
	return ones_gfx.Call[FabricSimItem](r.transport, "GET", "getFabricSimByName/"+name, nil, ones_gfx.OperationModeSynchronous, nil)
}

// Delete removes simulations by fabric name. Maps to DELETE /delFabricSim/{name}.
func (r *FabricsimsResource) Delete(ctx context.Context, name string) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "DELETE", "delFabricSim/"+name, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// UpdateStatus updates a simulation's operational status.
// Maps to POST /updateFabricSimStatus.
func (r *FabricsimsResource) UpdateStatus(ctx context.Context, name string, status *string) (string, error) {
	body := struct {
		Name   string  `json:"name"`
		Status *string `json:"status,omitempty"`
	}{name, status}
	res, err := ones_gfx.Call[string](r.transport, "POST", "updateFabricSimStatus", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}
