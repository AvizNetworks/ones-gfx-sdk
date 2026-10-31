package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// Request structs for the calls whose positional form takes several optional
// pointers. They are an alternative, not a replacement: the positional
// methods stay, and each ...With method validates the struct and then calls
// its positional twin, so both send identical requests.
//
// Optional fields use the zero value for "not set" wherever the zero value is
// not itself meaningful. Where it is (Shared: false differs from "keep the
// server default"), the field is a pointer; build it with ones.Ptr.

// CreateTenantRequest is the input of TenantsResource.CreateWith.
type CreateTenantRequest struct {
	FabricName string // required
	TenantName string // required

	// Description is omitted when empty.
	Description string

	// MaxGpusAllowed is the GPU quota: a positive number, or
	// ones_gfx.UnlimitedGPUs (-1) for no cap. 0 omits it, leaving the server
	// default.
	MaxGpusAllowed int

	// Shared is nil to leave the server default, or ones.Ptr(true/false).
	Shared *bool
}

func (q CreateTenantRequest) args() (description *string, maxGpusAllowed *int, shared *bool) {
	if q.Description != "" {
		d := q.Description
		description = &d
	}
	if q.MaxGpusAllowed != 0 {
		m := q.MaxGpusAllowed
		maxGpusAllowed = &m
	}
	return description, maxGpusAllowed, q.Shared
}

// Validate reports the first problem with q, without contacting the server.
func (q CreateTenantRequest) Validate() error {
	_, m, _ := q.args()
	return firstErr(fabricTenant(q.FabricName, q.TenantName), maxGpus(m))
}

// CreateWith creates a tenant (synchronous).
func (r *TenantsResource) CreateWith(ctx context.Context, req CreateTenantRequest, opts ...ones_gfx.CallOption) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	d, m, s := req.args()
	return r.Create(ctx, req.FabricName, req.TenantName, d, m, s, opts...)
}

// CreateAsyncWith creates a tenant and returns the 202 operation handle.
func (r *TenantsResource) CreateAsyncWith(ctx context.Context, req CreateTenantRequest, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	d, m, s := req.args()
	return r.CreateAsync(ctx, req.FabricName, req.TenantName, d, m, s, opts...)
}

// UpdateTenantServersRequest is the input of TenantsResource.UpdateWith:
// attach (GpuActionAdd) or detach (GpuActionDelete) whole servers.
type UpdateTenantServersRequest struct {
	FabricName string          // required
	TenantName string          // required
	Servers    []GpuServerInfo // required, at least one

	// Action is required. There is deliberately no default: a detach that
	// forgot to say so must not turn into an attach.
	Action ones_gfx.GpuAction
}

// Validate reports the first problem with q, without contacting the server.
func (q UpdateTenantServersRequest) Validate() error {
	return firstErr(fabricTenant(q.FabricName, q.TenantName), gpuServers(q.Servers), gpuAction("action", q.Action))
}

// UpdateWith attaches or detaches whole servers (synchronous).
func (r *TenantsResource) UpdateWith(ctx context.Context, req UpdateTenantServersRequest, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	a := req.Action
	return r.Update(ctx, req.FabricName, req.TenantName, req.Servers, &a, opts...)
}

// UpdateAsyncWith attaches or detaches whole servers and returns the 202
// operation handle.
func (r *TenantsResource) UpdateAsyncWith(ctx context.Context, req UpdateTenantServersRequest, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	a := req.Action
	return r.UpdateAsync(ctx, req.FabricName, req.TenantName, req.Servers, &a, opts...)
}

// ModifyAllocationsRequest is the input of TenantsResource.ModifyAllocationsWith:
// map (GpuActionAdd) or unmap (GpuActionDelete) specific GPUs.
type ModifyAllocationsRequest struct {
	FabricName string  // required
	TenantName string  // required
	Suid       SuidMap // required: server index -> hostname -> GPUs

	// Action is required; see UpdateTenantServersRequest.Action.
	Action ones_gfx.GpuAction

	// Scope is empty to leave the server default (WHOLE_SERVER).
	Scope ConfigScope

	// UnreachableDevices is omitted when empty.
	UnreachableDevices []string
}

func (q ModifyAllocationsRequest) scope() *ConfigScope {
	if q.Scope == "" {
		return nil
	}
	s := q.Scope
	return &s
}

// Validate reports the first problem with q, without contacting the server.
func (q ModifyAllocationsRequest) Validate() error {
	return firstErr(fabricTenant(q.FabricName, q.TenantName), suidMap(q.Suid),
		gpuAction("action", q.Action), optConfigScope(q.scope()))
}

// ModifyAllocationsWith maps or unmaps specific GPUs (synchronous).
func (r *TenantsResource) ModifyAllocationsWith(ctx context.Context, req ModifyAllocationsRequest, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	a := req.Action
	return r.ModifyAllocations(ctx, req.FabricName, req.TenantName, req.Suid, &a, req.scope(), req.UnreachableDevices, opts...)
}

// ModifyAllocationsAsyncWith maps or unmaps specific GPUs and returns the 202
// operation handle.
func (r *TenantsResource) ModifyAllocationsAsyncWith(ctx context.Context, req ModifyAllocationsRequest, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	a := req.Action
	return r.ModifyAllocationsAsync(ctx, req.FabricName, req.TenantName, req.Suid, &a, req.scope(), req.UnreachableDevices, opts...)
}

// AssignPortsRequest is the input of GPUResource.AssignPortsWith.
type AssignPortsRequest struct {
	FabricName  string             // required
	TenantName  string             // required
	Action      ones_gfx.GpuAction // required: GpuActionAdd or GpuActionDelete
	ServerNames []string           // required, at least one

	// GpuIDs is empty to act on every port of each server.
	GpuIDs []int

	// Membership is omitted when empty.
	Membership string
}

// Validate reports the first problem with q, without contacting the server.
func (q AssignPortsRequest) Validate() error {
	return firstErr(fabricTenant(q.FabricName, q.TenantName), gpuAction("action", q.Action),
		nonEmptyNames("serverNames", q.ServerNames), gpuIDs(q.GpuIDs))
}

// AssignPortsWith assigns or removes InfiniBand GPU ports (always synchronous).
func (r *GPUResource) AssignPortsWith(ctx context.Context, req AssignPortsRequest) (*GpuPortAssignmentResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var membership *string
	if req.Membership != "" {
		m := req.Membership
		membership = &m
	}
	return r.AssignPorts(ctx, req.FabricName, req.TenantName, req.Action, req.ServerNames, req.GpuIDs, membership)
}
