package resources

import (
	"context"
	"encoding/json"

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

// InventoryPort is one InfiniBand port as UfmServiceProvider.getInventoryPortsByHost
// reports it. TenantID and TenantName are set only while the port is allocated.
type InventoryPort struct {
	PortGUID        string             `json:"portGuid,omitempty"`
	DName           string             `json:"dname,omitempty"`
	PortName        string             `json:"portName,omitempty"`
	PortState       string             `json:"portState,omitempty"`
	LinkSpeedActive string             `json:"linkSpeedActive,omitempty"`
	TenantID        *int               `json:"tenantId,omitempty"`
	TenantName      string             `json:"tenantName,omitempty"`
	PkeyValue       string             `json:"pkeyValue,omitempty"`
	Allocated       bool               `json:"allocated"`
	AllocatedAt     ones_gfx.Timestamp `json:"allocatedAt,omitempty"`
	LastSyncedAt    ones_gfx.Timestamp `json:"lastSyncedAt,omitempty"`
}

// InventoryHost groups a host's ports. Hostname is the north-south hostname
// when FM can map it, otherwise the UFM system name; UFMSystemName is always
// the UFM one.
type InventoryHost struct {
	Hostname             string          `json:"hostname,omitempty"`
	UFMSystemName        string          `json:"ufmSystemName,omitempty"`
	IPAddress            string          `json:"ipAddress,omitempty"`
	SystemGUID           string          `json:"systemGuid,omitempty"`
	PortCount            int             `json:"portCount"`
	AvailableCount       int             `json:"availableCount"`
	AllocatedCount       int             `json:"allocatedCount"`
	LastConfiguredTenant string          `json:"lastConfiguredTenant,omitempty"`
	Ports                []InventoryPort `json:"ports"`
}

// InventoryPortsResult is the response of Inventory.Ports.
//
// A failed lookup (unknown fabric, non-UFM fabric) arrives as an HTTP error
// and is returned as one, so Success is true on every value you receive.
type InventoryPortsResult struct {
	FabricName         string          `json:"fabricName,omitempty"`
	Success            bool            `json:"success"`
	HostCount          int             `json:"hostCount"`
	TotalPortCount     int             `json:"totalPortCount"`
	AvailablePortCount int             `json:"availablePortCount"`
	AllocatedPortCount int             `json:"allocatedPortCount"`
	Hosts              []InventoryHost `json:"hosts"`

	// Extra holds any top-level key this SDK does not model yet.
	Extra map[string]json.RawMessage `json:"-"`
}

var inventoryPortsKeys = keySet("fabricName", "success", "hostCount", "totalPortCount",
	"availablePortCount", "allocatedPortCount", "hosts")

// UnmarshalJSON decodes the result and collects unmodelled keys into Extra.
func (r *InventoryPortsResult) UnmarshalJSON(data []byte) error {
	type alias InventoryPortsResult
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	r.Extra = extraKeys(data, inventoryPortsKeys)
	return nil
}

// Host returns the host whose Hostname or UFMSystemName is name.
func (r *InventoryPortsResult) Host(name string) (*InventoryHost, bool) {
	for i := range r.Hosts {
		if r.Hosts[i].Hostname == name || r.Hosts[i].UFMSystemName == name {
			return &r.Hosts[i], true
		}
	}
	return nil, false
}

// UFMHost is one host as UFM reports it live.
type UFMHost struct {
	SystemGUID string `json:"systemGuid,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	IPAddress  string `json:"ipAddress,omitempty"`
	State      string `json:"state,omitempty"`
	Vendor     string `json:"vendor,omitempty"`
	Role       string `json:"role,omitempty"`
	Type       string `json:"type,omitempty"`
	PortCount  int    `json:"portCount"`
}

// UFMHostsResult is the response of Inventory.UFMHosts.
type UFMHostsResult struct {
	FabricName string             `json:"fabricName,omitempty"`
	Timestamp  ones_gfx.Timestamp `json:"timestamp,omitempty"`
	Success    bool               `json:"success"`
	HostCount  int                `json:"hostCount"`
	Hosts      []UFMHost          `json:"hosts"`
	// Source is "UFM_LIVE": the list comes straight from UFM, not FM's cache.
	Source string `json:"source,omitempty"`

	// Extra holds any top-level key this SDK does not model yet.
	Extra map[string]json.RawMessage `json:"-"`
}

var ufmHostsKeys = keySet("fabricName", "timestamp", "success", "hostCount", "hosts", "source")

// UnmarshalJSON decodes the result and collects unmodelled keys into Extra.
func (r *UFMHostsResult) UnmarshalJSON(data []byte) error {
	type alias UFMHostsResult
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	r.Extra = extraKeys(data, ufmHostsKeys)
	return nil
}

// InventorySyncResult is the response of Inventory.Sync.
type InventorySyncResult struct {
	FabricName string             `json:"fabricName,omitempty"`
	Timestamp  ones_gfx.Timestamp `json:"timestamp,omitempty"`
	Success    bool               `json:"success"`
	Message    string             `json:"message,omitempty"`

	// Extra holds any top-level key this SDK does not model yet.
	Extra map[string]json.RawMessage `json:"-"`
}

var inventorySyncKeys = keySet("fabricName", "timestamp", "success", "message")

// UnmarshalJSON decodes the result and collects unmodelled keys into Extra.
func (r *InventorySyncResult) UnmarshalJSON(data []byte) error {
	type alias InventorySyncResult
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	r.Extra = extraKeys(data, inventorySyncKeys)
	return nil
}

// Ports returns UFM inventory ports grouped by host.
// Maps to GET /fabrics/{fabricName}/inventoryPorts.
func (r *InventoryResource) Ports(ctx context.Context, fabricName string) (*InventoryPortsResult, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[InventoryPortsResult](r.transport, "GET", "fabrics/"+fabricName+"/inventoryPorts", nil, ones_gfx.OperationModeSynchronous, nil)
}

// PortsRaw returns the inventoryPorts response exactly as FM sent it.
// Prefer Ports; this is the pre-typed behaviour, kept for migration.
func (r *InventoryResource) PortsRaw(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	return r.raw("GET", fabricName, "inventoryPorts")
}

// UFMHosts lists hosts directly from UFM.
// Maps to GET /fabrics/{fabricName}/inventoryHosts.
func (r *InventoryResource) UFMHosts(ctx context.Context, fabricName string) (*UFMHostsResult, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[UFMHostsResult](r.transport, "GET", "fabrics/"+fabricName+"/inventoryHosts", nil, ones_gfx.OperationModeSynchronous, nil)
}

// UFMHostsRaw returns the inventoryHosts response exactly as FM sent it.
// Prefer UFMHosts; this is the pre-typed behaviour, kept for migration.
func (r *InventoryResource) UFMHostsRaw(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	return r.raw("GET", fabricName, "inventoryHosts")
}

// Sync forces an immediate UFM inventory sync.
// Maps to POST /fabrics/{fabricName}/inventorySync.
func (r *InventoryResource) Sync(ctx context.Context, fabricName string) (*InventorySyncResult, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[InventorySyncResult](r.transport, "POST", "fabrics/"+fabricName+"/inventorySync", nil, ones_gfx.OperationModeSynchronous, nil)
}

// SyncRaw returns the inventorySync response exactly as FM sent it.
// Prefer Sync; this is the pre-typed behaviour, kept for migration.
func (r *InventoryResource) SyncRaw(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	return r.raw("POST", fabricName, "inventorySync")
}

func (r *InventoryResource) raw(method, fabricName, endpoint string) (map[string]interface{}, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, method, "fabrics/"+fabricName+"/"+endpoint, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// keySet builds a lookup of the wire keys a type models.
func keySet(keys ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

// extraKeys returns the top-level keys of data that are not in known, or nil
// when there are none.
func extraKeys(data []byte, known map[string]struct{}) map[string]json.RawMessage {
	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		return nil
	}
	var extra map[string]json.RawMessage
	for k, v := range all {
		if _, ok := known[k]; ok {
			continue
		}
		if extra == nil {
			extra = map[string]json.RawMessage{}
		}
		extra[k] = v
	}
	return extra
}
