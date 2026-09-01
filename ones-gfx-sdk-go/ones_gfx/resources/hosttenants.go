package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// HostTenantsResource covers host-based tenants, host listing and bulk
// host multi-tenancy updates.
type HostTenantsResource struct {
	transport *ones_gfx.Transport
}

// NewHostTenantsResource constructs a HostTenantsResource.
func NewHostTenantsResource(transport *ones_gfx.Transport) *HostTenantsResource {
	return &HostTenantsResource{transport: transport}
}

// HostAction is the ADD/DELETE discriminator of HostStatusUpdate
// (Helper/Enums.java).
type HostAction string

const (
	HostActionAdd    HostAction = "ADD"
	HostActionDelete HostAction = "DELETE"
)

// HosttenantsRecord mirrors Models/Hosttenants.java JSON.
type HosttenantsRecord struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	HostsAllocated *int   `json:"hostsAllocated,omitempty"`
	VniID          *int   `json:"vniId,omitempty"`
	FabricName     string `json:"fabricName,omitempty"`
	ConfigStatus   string `json:"config_status,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
}

// HostItem mirrors Models/Hosts.java JSON.
type HostItem struct {
	ID           int    `json:"id,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	HostStatus   string `json:"hostStatus,omitempty"`
	TenantName   string `json:"tenantName,omitempty"`
	FabricName   string `json:"fabricName,omitempty"`
	ConfigStatus string `json:"config_status,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}

// Add registers a host-based tenant record.
// Maps to POST /addHostTenantData.
func (r *HostTenantsResource) Add(ctx context.Context, name, fabricName string, description *string, hostsAllocated *int, vniID *int, configStatus *string) (string, error) {
	body := struct {
		Name           string  `json:"name"`
		FabricName     string  `json:"fabricName"`
		Description    *string `json:"description,omitempty"`
		HostsAllocated *int    `json:"hostsAllocated,omitempty"`
		VniID          *int    `json:"vniId,omitempty"`
		ConfigStatus   *string `json:"config_status,omitempty"`
	}{name, fabricName, description, hostsAllocated, vniID, configStatus}
	res, err := ones_gfx.Call[string](r.transport, "POST", "addHostTenantData", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Delete removes a host-tenant record.
// Maps to DELETE /delHostTenantData/{fabricName}/{tenantName}.
func (r *HostTenantsResource) Delete(ctx context.Context, fabricName, tenantName string) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "DELETE", "delHostTenantData/"+fabricName+"/"+tenantName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// List returns host tenants of a fabric.
// Maps to GET /getHostTenantsList/{fabricName}.
func (r *HostTenantsResource) List(ctx context.Context, fabricName string) ([]HosttenantsRecord, error) {
	res, err := ones_gfx.Call[[]HosttenantsRecord](r.transport, "GET", "getHostTenantsList/"+fabricName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ListHosts returns hosts of a fabric.
// Maps to GET /getHostsList/{fabricName}.
func (r *HostTenantsResource) ListHosts(ctx context.Context, fabricName string) ([]HostItem, error) {
	res, err := ones_gfx.Call[[]HostItem](r.transport, "GET", "getHostsList/"+fabricName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// UpdateHosts bulk-adds or removes hosts for a tenant. All four fields are
// required server-side. Maps to POST /updateHosts.
func (r *HostTenantsResource) UpdateHosts(ctx context.Context, hostnames []string, hostAction HostAction, tenantName, fabricName string) (bool, error) {
	body := struct {
		Hostnames  []string   `json:"hostnames"`
		HostAction HostAction `json:"hostAction"`
		TenantName string     `json:"tenantName"`
		FabricName string     `json:"fabricName"`
	}{hostnames, hostAction, tenantName, fabricName}
	res, err := ones_gfx.Call[bool](r.transport, "POST", "updateHosts", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}
