package resources

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// IpamResource covers the fabric's IPAM configuration and the addresses IPAM
// has handed out.
type IpamResource struct {
	transport *ones_gfx.Transport
}

// NewIpamResource constructs an IpamResource.
func NewIpamResource(transport *ones_gfx.Transport) *IpamResource {
	return &IpamResource{transport: transport}
}

// IP assignment states, used by the status filter on TenantIPs.
const (
	IPAssignStatusAllocatedOnly = "ALLOCATED_ONLY"
	IPAssignStatusAssigned      = "ASSIGNED"
	IPAssignStatusFailed        = "IP_ASSIGN_FAILED"
)

// IpamStaticConfig is the static address-assignment mode. With EnableIPAssign
// the chosen address is pushed to the host over SSH.
type IpamStaticConfig struct {
	Enabled        bool  `json:"enabled"`
	EnableIPAssign *bool `json:"enableIpAssign,omitempty"`
}

// IpamDhcpConfig is the DHCP address-assignment mode. LeaseTimeSeconds is its
// only setting; renew and rebind follow it at the RFC 2131 defaults of 50%
// and 87.5%.
type IpamDhcpConfig struct {
	Enabled          bool `json:"enabled"`
	LeaseTimeSeconds *int `json:"leaseTimeSeconds,omitempty"`
}

// IpamConfig is the fabric's IPAM configuration, in the shape GET returns and
// Update accepts.
//
// IPAM hands out addresses in one of two modes. With IPAM on, exactly one of
// Static and Dhcp is enabled; with IPAM off, neither. A mode left nil counts
// as disabled, so send only the mode in use.
type IpamConfig struct {
	FabricName string            `json:"fabricName,omitempty"`
	EnableIpam bool              `json:"enableIpam"`
	Static     *IpamStaticConfig `json:"static,omitempty"`
	Dhcp       *IpamDhcpConfig   `json:"dhcp,omitempty"`
}

// IpamPoolSummary is the tenant pool's counters, returned alongside its
// allocations.
type IpamPoolSummary struct {
	Cidr                    string `json:"cidr"`
	AllocatableAddressCount int64  `json:"allocatableAddressCount"`
	AllocatedCount          int64  `json:"allocatedCount"`
}

// IpamPagination is the paging block shared by the paged IPAM listings.
type IpamPagination struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

// IpamAllocation is one address IPAM has allocated. The payload varies by
// fabric and mode, so it is surfaced as a raw map.
type IpamAllocation = map[string]interface{}

// TenantIPList is one page of a tenant's allocations, with its pool counters.
type TenantIPList struct {
	FabricName  string           `json:"fabricName"`
	TenantName  string           `json:"tenantName"`
	Pool        IpamPoolSummary  `json:"pool"`
	Allocations []IpamAllocation `json:"allocations"`
	Pagination  IpamPagination   `json:"pagination"`
}

// ServerIPList is every address allocated to one server, across its tenants.
// Unpaginated: a server holds one address per tenant it belongs to.
type ServerIPList struct {
	FabricName  string           `json:"fabricName"`
	ServerName  string           `json:"serverName"`
	Allocations []IpamAllocation `json:"allocations"`
}

// GetConfig returns the fabric's IPAM configuration, in the same shape Update
// accepts, so a body read here can be sent straight back.
// Maps to GET /fabrics/{fabricName}/ipam/config.
// An unknown fabric answers 404.
func (r *IpamResource) GetConfig(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (*IpamConfig, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[IpamConfig](r.transport, "GET", "fabrics/"+fabricName+"/ipam/config", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// UpdateConfig changes the fabric's IPAM configuration.
// Maps to PATCH /fabrics/{fabricName}/ipam/config.
//
// Settings may only accompany an enabled mode, and an enabled mode's omitted
// settings keep their saved values. FM rejects any field outside this set
// rather than ignoring it, so a mode setting sent at the top level fails
// clearly. force pushes the change through checks FM would otherwise fail on.
//
// Only externally managed (non-ONES-controlled) fabrics may be changed;
// others answer 409 FABRIC_NOT_EXTERNALLY_MANAGED. 404 when the fabric is
// unknown.
func (r *IpamResource) UpdateConfig(ctx context.Context, fabricName string, cfg IpamConfig, force bool, opts ...ones_gfx.CallOption) (*IpamConfig, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	if force {
		ro.Query = url.Values{"force": []string{"true"}}
	}
	// FM rejects a body whose fabricName disagrees with the URL, so the field
	// is dropped: it exists only so a body read from GetConfig can be sent
	// straight back, and carrying a stale one here would fail the call.
	cfg.FabricName = ""
	return ones_gfx.Call[IpamConfig](r.transport, "PATCH", "fabrics/"+fabricName+"/ipam/config", cfg, ones_gfx.OperationModeSynchronous, ro)
}

// TenantIPs returns one page of a tenant's IPAM allocations with the pool's
// counters.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/ip-addresses.
//
// page is zero-based and defaults to 0 when negative; size defaults to 50
// when not positive. status filters by assignment state — pass "" for no
// filter. A tenant with no pool answers 404 IPAM_NOT_CONFIGURED.
func (r *IpamResource) TenantIPs(ctx context.Context, fabricName, tenantName string, page, size int, status string, opts ...ones_gfx.CallOption) (*TenantIPList, error) {
	if err := fabricTenant(fabricName, tenantName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	q := url.Values{}
	if page >= 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}
	if status != "" {
		q.Set("status", status)
	}
	if len(q) > 0 {
		ro.Query = q
	}
	return ones_gfx.Call[TenantIPList](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/ip-addresses", nil, ones_gfx.OperationModeSynchronous, ro)
}

// ServerIPs returns every address allocated to one compute server, across all
// its tenants.
// Maps to GET /fabrics/{fabricName}/computes/{serverName}/ip-addresses.
// A server missing from the fabric inventory answers 404 SERVER_NOT_FOUND.
func (r *IpamResource) ServerIPs(ctx context.Context, fabricName, serverName string, opts ...ones_gfx.CallOption) (*ServerIPList, error) {
	if err := firstErr(pathName("fabricName", fabricName), pathName("serverName", serverName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[ServerIPList](r.transport, "GET", "fabrics/"+fabricName+"/computes/"+serverName+"/ip-addresses", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// TenantComputeIP returns the single address allocated to one server within
// one tenant.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/computes/{serverName}/ip-address.
//
// FM answers 404 either way when there is nothing to return, distinguished by
// the error code: IPAM_NOT_CONFIGURED when the tenant has no pool at all,
// SERVER_NOT_ALLOCATED_TO_TENANT when it has one but this server holds no
// address in it.
func (r *IpamResource) TenantComputeIP(ctx context.Context, fabricName, tenantName, serverName string, opts ...ones_gfx.CallOption) (IpamAllocation, error) {
	if err := firstErr(fabricTenant(fabricName, tenantName), pathName("serverName", serverName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[IpamAllocation](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/computes/"+serverName+"/ip-address", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
