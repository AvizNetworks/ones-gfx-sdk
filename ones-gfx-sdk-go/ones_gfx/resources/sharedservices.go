package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// SharedServicesResource covers the fabric's shared-services tenant
// (GFX-5709). A fabric has at most one: a tenant every host can reach,
// configured on every DPU the FM manages and on the border/storage leaves,
// with its own VLAN and, where IPAM/DHCP is enabled, its own Kea subnet.
type SharedServicesResource struct {
	transport *ones_gfx.Transport
}

// NewSharedServicesResource constructs a SharedServicesResource.
func NewSharedServicesResource(transport *ones_gfx.Transport) *SharedServicesResource {
	return &SharedServicesResource{transport: transport}
}

// SharedServicesDPU is one DPU's outcome in a shared-services response.
// Status is empty on a GET, which does not contact the DPUs. A skipped DPU's
// Reason is DPU_CONFIG_DISABLED (executeConfig=no) or DPU_NOT_MANAGED.
type SharedServicesDPU struct {
	DPUName string `json:"dpuName"`
	Status  string `json:"status,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// SharedServicesVF is one VF assigned to the shared-services tenant.
type SharedServicesVF struct {
	ServerName string `json:"serverName"`
	DPUName    string `json:"dpuName"`
	IfName     string `json:"ifName"`
	ServerIf   string `json:"serverIf"`
	Status     string `json:"status"`
}

// SharedServices is the shared-services tenant as FM returns it. The body is
// not wrapped in the status/message envelope. AnycastGateway is set only on
// HBN fabrics in tenant DHCP mode; Leaves only on the create response.
type SharedServices struct {
	FabricName     string              `json:"fabricName"`
	TenantName     string              `json:"tenantName"`
	TenantType     string              `json:"tenantType"`
	VlanID         int                 `json:"vlanId"`
	VlanVniID      int                 `json:"vlanVniId"`
	VrfVniID       int                 `json:"vrfVniId"`
	AnycastGateway string              `json:"anycastGateway,omitempty"`
	DPUs           []SharedServicesDPU `json:"dpus"`
	Leaves         []string            `json:"leaves,omitempty"`
	VFs            []SharedServicesVF  `json:"vfs"`
}

// Create makes the fabric's shared-services tenant.
// Maps to POST /fabrics/{fabricName}/shared-services, which returns 201.
//
// tenantName starts with a letter and is at most 12 characters of letters,
// digits, '-' and '_'; "Storage" is reserved. vlanID must be in 2-4094,
// excluding 500-1020, 3725-3999, 69, 100 and 110, and unused by another
// tenant. description is at most 60 characters. anycastGateway may be nil.
func (r *SharedServicesResource) Create(ctx context.Context, fabricName, tenantName string, vlanID int, description *string, anycastGateway *ones_gfx.AnycastGateway, opts ...ones_gfx.CallOption) (*SharedServices, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		TenantName     string                   `json:"tenantName"`
		VlanID         int                      `json:"vlanId"`
		Description    *string                  `json:"description,omitempty"`
		AnycastGateway *ones_gfx.AnycastGateway `json:"anycastGateway,omitempty"`
	}{tenantName, vlanID, description, anycastGateway}
	return ones_gfx.Call[SharedServices](r.transport, "POST", "fabrics/"+fabricName+"/shared-services", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Get returns the fabric's shared-services tenant and the VFs assigned to it.
// Maps to GET /fabrics/{fabricName}/shared-services.
// A fabric with none answers 404 SHARED_SERVICES_NOT_FOUND.
func (r *SharedServicesResource) Get(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (*SharedServices, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[SharedServices](r.transport, "GET", "fabrics/"+fabricName+"/shared-services", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Delete removes the shared-services tenant from the leaves and DPUs and
// releases its DHCP subnet and IPAM pool.
// Maps to DELETE /fabrics/{fabricName}/shared-services.
//
// Unassign every VF first: FM answers 409 VFS_ASSIGNED and lists them while
// any remain. If a delete fails partway the tenant is kept, so calling this
// again is safe and is how a partial failure is retried.
func (r *SharedServicesResource) Delete(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[ApiResponseMessage](r.transport, "DELETE", "fabrics/"+fabricName+"/shared-services", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}
