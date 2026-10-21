package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// InventoryResource covers inventory CRUD and UFM sync/query endpoints.
type InventoryResource struct {
	transport *ones_gfx.Transport
}

// NewInventoryResource constructs an InventoryResource.
func NewInventoryResource(transport *ones_gfx.Transport) *InventoryResource {
	return &InventoryResource{transport: transport}
}

// Ports returns UFM inventory ports grouped by host. The payload mixes
// fixed keys (success, error, message, fabricName) with dynamic per-host
// data, so it is surfaced as a raw map.
// Maps to GET /fabrics/{fabricName}/inventoryPorts.
func (r *InventoryResource) Ports(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/inventoryPorts", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// UFMHosts lists hosts directly from UFM. Dynamic payload → raw map.
// Maps to GET /fabrics/{fabricName}/inventoryHosts.
func (r *InventoryResource) UFMHosts(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/inventoryHosts", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Sync forces an immediate UFM inventory sync. Dynamic payload → raw map.
// Maps to POST /fabrics/{fabricName}/inventorySync.
func (r *InventoryResource) Sync(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "POST", "fabrics/"+fabricName+"/inventorySync", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
