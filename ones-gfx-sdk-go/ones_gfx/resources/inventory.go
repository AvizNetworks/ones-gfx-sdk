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

// InventoryItem mirrors Models/Inventory.java as a request payload
// (shared by Add and Edit).
type InventoryItem struct {
	IPAddress     string `json:"ipAddress,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	DeviceType    string `json:"deviceType,omitempty"`
	DeviceRole    string `json:"deviceRole,omitempty"`
	Sku           string `json:"sku,omitempty"`
	InterfaceData string `json:"interfaceData,omitempty"`
	Status        string `json:"status,omitempty"`
	FabricName    string `json:"fabricName,omitempty"`
	ExecuteConfig string `json:"executeConfig,omitempty"`
	UfmSystemName string `json:"ufmSystemName,omitempty"`
}

// FabricInventoryUpdateItem mirrors Helper/FabricInventoryUpdate.java.
// fabricName/hostname are always required server-side; ipAddress/username/
// password are required unless Configure == "EXECUTE_CONFIG_NO".
type FabricInventoryUpdateItem struct {
	FabricName string  `json:"fabricName"`
	Hostname   string  `json:"hostname"`
	IPAddress  *string `json:"ipAddress,omitempty"`
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	Configure  *string `json:"configure,omitempty"`
}

// InventoryRecord mirrors Models/Inventory.java JSON (responses).
type InventoryRecord struct {
	ID            int    `json:"id,omitempty"`
	IPAddress     string `json:"ipAddress,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	DeviceType    string `json:"deviceType,omitempty"`
	DeviceRole    string `json:"deviceRole,omitempty"`
	Sku           string `json:"sku,omitempty"`
	InterfaceData string `json:"interfaceData,omitempty"`
	Status        string `json:"status,omitempty"`
	FabricName    string `json:"fabricName,omitempty"`
	ExecuteConfig string `json:"executeConfig,omitempty"`
	UfmSystemName string `json:"ufmSystemName,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
}

// Add inserts inventory rows. Maps to POST /addInventoryData.
func (r *InventoryResource) Add(ctx context.Context, items []InventoryItem) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "POST", "addInventoryData", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Update updates credentials/hostname per device.
// Maps to POST /updateInventoryData.
func (r *InventoryResource) Update(ctx context.Context, items []FabricInventoryUpdateItem) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "POST", "updateInventoryData", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Edit partially updates inventory scoped to a fabric.
// Maps to PATCH /fabrics/{fabricName}/editInventoryData.
func (r *InventoryResource) Edit(ctx context.Context, fabricName string, items []InventoryItem) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "PATCH", "fabrics/"+fabricName+"/editInventoryData", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// ListAll returns every inventory entry. Maps to GET /getAllInventory.
func (r *InventoryResource) ListAll(ctx context.Context) ([]InventoryRecord, error) {
	res, err := ones_gfx.Call[[]InventoryRecord](r.transport, "GET", "getAllInventory", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ByFabric returns inventory rows of one fabric.
// Maps to GET /getInventoryByFabricName/{name}.
func (r *InventoryResource) ByFabric(ctx context.Context, name string) ([]InventoryRecord, error) {
	res, err := ones_gfx.Call[[]InventoryRecord](r.transport, "GET", "getInventoryByFabricName/"+name, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
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

// MaxInventoryDeleteHostnames is the largest number of unique hostnames FM
// accepts in one DeleteByHostnames call.
const MaxInventoryDeleteHostnames = 200

// InventoryDeleteResult reports which hostnames a delete did not match.
//
// NotFound is always present, because the status alone does not say which
// names missed: FM answers 200 when every hostname matched, 207 when only
// some did, and 404 when none did. Check NotFound even on success.
type InventoryDeleteResult struct {
	NotFound []string `json:"notFound"`
	Message  string   `json:"message,omitempty"`
}

// DeleteByHostnames removes inventory devices by hostname, scoped to one
// fabric. The request body is a bare JSON array of hostnames.
// Maps to DELETE /fabrics/{fabricName}/inventory.
//
// Duplicates and blanks are dropped before the delete, and at most
// MaxInventoryDeleteHostnames unique names may be sent.
func (r *InventoryResource) DeleteByHostnames(ctx context.Context, fabricName string, hostnames []string, opts ...ones_gfx.CallOption) (*InventoryDeleteResult, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[InventoryDeleteResult](r.transport, "DELETE", "fabrics/"+fabricName+"/inventory", hostnames, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}
