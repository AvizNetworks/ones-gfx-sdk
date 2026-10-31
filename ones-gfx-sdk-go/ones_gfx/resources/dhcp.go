package resources

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// DhcpResource covers the tenant DHCP subnets Fabric Manager maintains in
// Kea: reading one, retrying a failed one, and reconciling them all.
type DhcpResource struct {
	transport *ones_gfx.Transport
}

// NewDhcpResource constructs a DhcpResource.
func NewDhcpResource(transport *ones_gfx.Transport) *DhcpResource {
	return &DhcpResource{transport: transport}
}

// MaxDhcpLeasePageSize is the largest lease page FM will return.
const MaxDhcpLeasePageSize = 1000

// TenantDhcp is the tenant's DHCP subnet, whether Kea is serving it, and one
// page of the leases handed out in it.
//
// LiveInKea is reported separately from the subnet's stored status on
// purpose: the two disagreeing is exactly the drift reconciliation exists to
// repair. It is nil, with LeaseCheckError set, when Kea could not be queried.
// Leases are ordered by IP address, so pages are stable while the lease set
// is unchanged.
type TenantDhcp struct {
	Subnet          map[string]interface{}   `json:"subnet"`
	LiveInKea       *bool                    `json:"liveInKea"`
	Leases          []map[string]interface{} `json:"leases"`
	LeaseCheckError string                   `json:"leaseCheckError,omitempty"`
	Pagination      *IpamPagination          `json:"pagination,omitempty"`
}

// DhcpReapplySummary is the outcome of a reconcile or retry. The payload
// varies with what FM had to change, so it is surfaced as a raw map.
type DhcpReapplySummary = map[string]interface{}

// GetTenant returns the tenant's DHCP subnet and a page of its leases.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/dhcp.
//
// page is zero-based and defaults to 0 when negative; size defaults to 50
// when not positive and is capped at MaxDhcpLeasePageSize.
//
// 404 DHCP_SUBNET_NOT_FOUND when the fabric does not assign addresses by
// DHCP, or the tenant has no VLAN or IPAM pool.
func (r *DhcpResource) GetTenant(ctx context.Context, fabricName, tenantName string, page, size int, opts ...ones_gfx.CallOption) (*TenantDhcp, error) {
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
	if len(q) > 0 {
		ro.Query = q
	}
	return ones_gfx.Call[TenantDhcp](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/dhcp", nil, ones_gfx.OperationModeSynchronous, ro)
}

// RetryTenant re-applies a tenant's desired DHCP state to Kea.
// Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/dhcp/retry.
//
// Idempotent: the fragment is rendered from the database and applied whole,
// so calling this when nothing is wrong re-applies the same configuration.
//
// 404 DHCP_SUBNET_NOT_FOUND when DHCP does not apply to the tenant, 400 on a
// validation failure, 502 when Kea cannot be reached.
func (r *DhcpResource) RetryTenant(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) (DhcpReapplySummary, error) {
	if err := fabricTenant(fabricName, tenantName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[DhcpReapplySummary](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/dhcp/retry", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Reconcile re-applies every tenant DHCP subnet so Kea matches Fabric
// Manager's records.
// Maps to POST /fabrics/{fabricName}/dhcp/reconcile.
//
// The same work the scheduled reconciliation job does, for an operator who
// does not want to wait for the next run. Only the tenant-owned part of the
// Kea config is rewritten — bootstrap subnets are never touched. fabricName
// scopes the log line, not the work: the re-apply covers every tenant FM
// knows.
//
// 404 when the fabric is unknown, 400 on a validation failure, 502 when Kea
// cannot be reached.
func (r *DhcpResource) Reconcile(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (DhcpReapplySummary, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[DhcpReapplySummary](r.transport, "POST", "fabrics/"+fabricName+"/dhcp/reconcile", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
