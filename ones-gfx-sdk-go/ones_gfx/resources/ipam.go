package resources

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// IPAMResource covers the per-fabric IPAM configuration and the IP allocation
// read endpoints under /api/fm. (Allocation itself happens as a side effect of
// Tenants.Update; see GpuUpdateResult.)
type IPAMResource struct {
	transport *ones_gfx.Transport
}

// NewIPAMResource constructs an IPAMResource.
func NewIPAMResource(transport *ones_gfx.Transport) *IPAMResource {
	return &IPAMResource{transport: transport}
}

// IpAssignStatus is the state of an allocated address (Models/IpamAllocation.AssignStatus).
type IpAssignStatus string

const (
	IpAssignStatusAllocatedOnly IpAssignStatus = "ALLOCATED_ONLY"
	IpAssignStatusAssigned      IpAssignStatus = "ASSIGNED"
	IpAssignStatusAssignFailed  IpAssignStatus = "IP_ASSIGN_FAILED"
)

func (s IpAssignStatus) valid() bool {
	switch s {
	case IpAssignStatusAllocatedOnly, IpAssignStatusAssigned, IpAssignStatusAssignFailed:
		return true
	}
	return false
}

// IpCleanupStatus is the cleanup state of an allocation (Models/IpamAllocation.CleanupStatus).
type IpCleanupStatus string

const (
	IpCleanupStatusNone    IpCleanupStatus = "NONE"
	IpCleanupStatusPending IpCleanupStatus = "PENDING"
)

// IpamConfig mirrors Models/IpamServiceConfig.java. UpdatedBy is nil until the
// config has been written once; a never-configured fabric reports both
// switches false.
type IpamConfig struct {
	FabricName     string  `json:"fabricName"`
	EnableIpam     bool    `json:"enableIpam"`
	EnableIpAssign bool    `json:"enableIpAssign"`
	UpdatedAt      string  `json:"updatedAt"`
	UpdatedBy      *string `json:"updatedBy"`
}

// IpAllocation mirrors Models/IpamAllocation.java. IPAddress and PrefixLength
// are separate fields ("192.3.0.2" and 32). InterfaceName is set only when the
// address lives on an HBN VF rather than the N-S NIC.
type IpAllocation struct {
	ID            int64           `json:"id"`
	PoolID        int64           `json:"poolId"`
	FabricName    string          `json:"fabricName"`
	TenantName    string          `json:"tenantName"`
	TenantID      *int            `json:"tenantId"`
	ServerName    string          `json:"serverName"`
	IPAddress     string          `json:"ipAddress"`
	PrefixLength  int             `json:"prefixLength"`
	AssignStatus  IpAssignStatus  `json:"assignStatus"`
	CleanupStatus IpCleanupStatus `json:"cleanupStatus"`
	AllocatedAt   string          `json:"allocatedAt"`
	InterfaceName *string         `json:"interfaceName,omitempty"`
}

// TenantIpPool is the tenant's address pool and its usage.
type TenantIpPool struct {
	CIDR                    string `json:"cidr"`
	AllocatableAddressCount int64  `json:"allocatableAddressCount"`
	AllocatedCount          int64  `json:"allocatedCount"`
}

// Pagination is zero-based page metadata.
type Pagination struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

// TenantIpAllocations is the result of GET /fabrics/{f}/tenants/{t}/ip-addresses.
type TenantIpAllocations struct {
	FabricName  string         `json:"fabricName"`
	TenantName  string         `json:"tenantName"`
	Pool        TenantIpPool   `json:"pool"`
	Allocations []IpAllocation `json:"allocations"`
	Pagination  Pagination     `json:"pagination"`
}

// ServerIpAllocations is the result of GET /fabrics/{f}/computes/{s}/ip-addresses.
type ServerIpAllocations struct {
	FabricName  string         `json:"fabricName"`
	ServerName  string         `json:"serverName"`
	Allocations []IpAllocation `json:"allocations"`
}

// GetConfig returns the IPAM switches of a fabric.
// Maps to GET /fabrics/{fabricName}/ipam/config.
//
// Errors: *ones_gfx.NotFoundError (FABRIC_NOT_FOUND).
func (r *IPAMResource) GetConfig(ctx context.Context, fabricName string) (*IpamConfig, error) {
	segs, err := pathSegments("fabricName", fabricName)
	if err != nil {
		return nil, err
	}
	return ones_gfx.Call[IpamConfig](r.transport, "GET", "fabrics/"+segs[0]+"/ipam/config", nil, ones_gfx.OperationModeSynchronous, nil)
}

// UpdateConfig sets both IPAM switches and returns the stored config. A nil
// force means the server default (false); force=true allows turning IPAM off
// while active pools exist — the pools are disabled and their allocations
// moved to history.
// Maps to PATCH /fabrics/{fabricName}/ipam/config.
//
// enableIpAssign without enableIpam is rejected locally with a
// *ones_gfx.ValidationError (the server rejects it too).
// Errors: *ones_gfx.NotFoundError (FABRIC_NOT_FOUND); *ones_gfx.ConflictError
// (FABRIC_NOT_EXTERNALLY_MANAGED for ONES-controlled fabrics, or
// IPAM_HAS_ACTIVE_POOLS when disabling without force — the code is in
// ResponseBody's error.code); *ones_gfx.BadRequestError.
func (r *IPAMResource) UpdateConfig(ctx context.Context, fabricName string, enableIpam, enableIpAssign bool, force *bool) (*IpamConfig, error) {
	segs, err := pathSegments("fabricName", fabricName)
	if err != nil {
		return nil, err
	}
	if enableIpAssign && !enableIpam {
		return nil, &ones_gfx.ValidationError{Field: "enableIpAssign", Reason: "requires enableIpam to be true"}
	}
	body := struct {
		EnableIpam     bool `json:"enableIpam"`
		EnableIpAssign bool `json:"enableIpAssign"`
	}{enableIpam, enableIpAssign}
	var ro *ones_gfx.ReqOpts
	if force != nil {
		ro = &ones_gfx.ReqOpts{Query: url.Values{"force": {strconv.FormatBool(*force)}}}
	}
	return ones_gfx.Call[IpamConfig](r.transport, "PATCH", "fabrics/"+segs[0]+"/ipam/config", body, ones_gfx.OperationModeSynchronous, ro)
}

// TenantComputeIP returns the IP allocation of one compute server in a tenant.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address.
//
// Errors: *ones_gfx.NotFoundError — IPAM_NOT_CONFIGURED (tenant has no pool) or
// SERVER_NOT_ALLOCATED_TO_TENANT (server has no allocation).
func (r *IPAMResource) TenantComputeIP(ctx context.Context, fabricName, tenantName, serverName string) (*IpAllocation, error) {
	segs, err := pathSegments("fabricName", fabricName, "tenantName", tenantName, "serverName", serverName)
	if err != nil {
		return nil, err
	}
	path := "fabrics/" + segs[0] + "/tenants/" + segs[1] + "/computes/" + segs[2] + "/ip-address"
	return ones_gfx.Call[IpAllocation](r.transport, "GET", path, nil, ones_gfx.OperationModeSynchronous, nil)
}

// TenantIPs lists a tenant's IP allocations and pool usage, one page at a
// time. page is zero-based (server default 0), size defaults to 50, status
// filters on the assign status; nil omits each. Out-of-range values are
// rejected locally with a *ones_gfx.ValidationError (the server would silently
// clamp them).
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/ip-addresses.
//
// Errors: *ones_gfx.NotFoundError (IPAM_NOT_CONFIGURED).
func (r *IPAMResource) TenantIPs(ctx context.Context, fabricName, tenantName string, page, size *int, status *IpAssignStatus) (*TenantIpAllocations, error) {
	segs, err := pathSegments("fabricName", fabricName, "tenantName", tenantName)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if page != nil {
		if *page < 0 {
			return nil, &ones_gfx.ValidationError{Field: "page", Reason: fmt.Sprintf("must be >= 0, got %d", *page)}
		}
		q.Set("page", strconv.Itoa(*page))
	}
	if size != nil {
		if *size < 1 {
			return nil, &ones_gfx.ValidationError{Field: "size", Reason: fmt.Sprintf("must be >= 1, got %d", *size)}
		}
		q.Set("size", strconv.Itoa(*size))
	}
	if status != nil {
		if !status.valid() {
			return nil, &ones_gfx.ValidationError{Field: "status", Reason: fmt.Sprintf("must be one of %q, %q, %q; got %q", IpAssignStatusAllocatedOnly, IpAssignStatusAssigned, IpAssignStatusAssignFailed, string(*status))}
		}
		q.Set("status", string(*status))
	}
	var ro *ones_gfx.ReqOpts
	if len(q) > 0 {
		ro = &ones_gfx.ReqOpts{Query: q}
	}
	path := "fabrics/" + segs[0] + "/tenants/" + segs[1] + "/ip-addresses"
	return ones_gfx.Call[TenantIpAllocations](r.transport, "GET", path, nil, ones_gfx.OperationModeSynchronous, ro)
}

// ServerIPs lists every IP allocation of a compute server across tenants
// (empty when it has none).
// Maps to GET /fabrics/{fabricName}/computes/{serverName}/ip-addresses.
//
// Errors: *ones_gfx.NotFoundError (SERVER_NOT_FOUND).
func (r *IPAMResource) ServerIPs(ctx context.Context, fabricName, serverName string) (*ServerIpAllocations, error) {
	segs, err := pathSegments("fabricName", fabricName, "serverName", serverName)
	if err != nil {
		return nil, err
	}
	path := "fabrics/" + segs[0] + "/computes/" + segs[1] + "/ip-addresses"
	return ones_gfx.Call[ServerIpAllocations](r.transport, "GET", path, nil, ones_gfx.OperationModeSynchronous, nil)
}
