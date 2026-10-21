package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// VFInterfacesResource covers the DPU VF interface endpoints under /api/fm.
// Assigning a VF is how a host joins a tenant, including the shared-services
// tenant.
type VFInterfacesResource struct {
	transport *ones_gfx.Transport
}

// NewVFInterfacesResource constructs a VFInterfacesResource.
func NewVFInterfacesResource(transport *ones_gfx.Transport) *VFInterfacesResource {
	return &VFInterfacesResource{transport: transport}
}

// VFInterface is one VF on a server's DPU. A VF is addressed either by
// IfName (the DPU interface, "pf0vf0_if") or by ServerIf (the host
// interface, "enp193s0f1np1"). TenantName is empty while the VF is free.
type VFInterface struct {
	ServerName string `json:"serverName,omitempty"`
	DPUName    string `json:"dpuName,omitempty"`
	IfName     string `json:"ifName"`
	ServerIf   string `json:"serverIf,omitempty"`
	TenantName string `json:"tenantName,omitempty"`
	Status     string `json:"status,omitempty"`
	VlanID     int    `json:"vlanId,omitempty"`
	VniID      int    `json:"vniId,omitempty"`
	Gateway    string `json:"gateway,omitempty"`
}

// VFAssignment is the result of assigning a VF. The Ip* fields are set only
// when IPAM or DHCP is enabled for the fabric; a 207 Multi-Status carries
// IPAssignStatus = IP_ASSIGN_FAILED, meaning the VF was assigned but its IP
// setup was not.
type VFAssignment struct {
	VfID           string `json:"vfId"`
	ServerName     string `json:"serverName"`
	FabricName     string `json:"fabricName"`
	TenantName     string `json:"tenantName"`
	Status         string `json:"status"`
	VlanID         int    `json:"vlanId,omitempty"`
	VniID          int    `json:"vniId,omitempty"`
	Gateway        string `json:"gateway,omitempty"`
	ProvisionedAt  string `json:"provisionedAt,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	IPAssignStatus string `json:"ipAssignStatus,omitempty"`
	IPInterface    string `json:"ipInterface,omitempty"`
}

// VFDpuEntry groups a DPU with the VFs it presents to the server.
type VFDpuEntry struct {
	DpuName      string        `json:"dpuName"`
	VfInterfaces []VFInterface `json:"vfInterfaces"`
}

// VFInterfaceList is what FM returns for a server's VF interfaces: the VFs
// grouped by the DPU that backs them, not a flat array.
type VFInterfaceList struct {
	ServerName string       `json:"serverName"`
	FabricName string       `json:"fabricName"`
	DpuCount   int          `json:"dpuCount"`
	Dpus       []VFDpuEntry `json:"dpus"`
}

// List returns the VF interfaces on a server, grouped by DPU.
// Maps to GET /fabrics/{fabricName}/servers/{serverName}/vf-interfaces.
func (r *VFInterfacesResource) List(ctx context.Context, fabricName, serverName string, opts ...ones_gfx.CallOption) (*VFInterfaceList, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[VFInterfaceList](r.transport, "GET", "fabrics/"+fabricName+"/servers/"+serverName+"/vf-interfaces", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Flatten returns every VF across all DPUs, for callers that do not care which
// DPU presents which.
func (l *VFInterfaceList) Flatten() []VFInterface {
	if l == nil {
		return nil
	}
	var out []VFInterface
	for _, d := range l.Dpus {
		out = append(out, d.VfInterfaces...)
	}
	return out
}

// Assign attaches a VF to a tenant.
// Maps to POST /fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign.
//
// FM requires the tenant to exist and the server to be in it first. vfID is
// the DPU interface name or the host interface name.
func (r *VFInterfacesResource) Assign(ctx context.Context, fabricName, serverName, vfID, tenantName string, opts ...ones_gfx.CallOption) (*VFAssignment, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		TenantName string `json:"tenantName"`
	}{tenantName}
	path := "fabrics/" + fabricName + "/servers/" + serverName + "/vf-interfaces/" + vfID + "/assign"
	return ones_gfx.Call[VFAssignment](r.transport, "POST", path, body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Unassign returns a VF to free status.
// Maps to DELETE /fabrics/{fabricName}/servers/{serverName}/vf-interfaces/{vfId}/assign.
//
// Order matters on teardown: a VF cannot be freed once its server has left
// the tenant, so unassign every VF before detaching GPU servers or deleting
// the tenant. CleanUpTenant does this in order.
func (r *VFInterfacesResource) Unassign(ctx context.Context, fabricName, serverName, vfID string, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	path := "fabrics/" + fabricName + "/servers/" + serverName + "/vf-interfaces/" + vfID + "/assign"
	return ones_gfx.Call[ApiResponseMessage](r.transport, "DELETE", path, nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}
