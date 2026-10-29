package resources

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// VPCPeeringResource covers VPC route-leaking (peering) configuration.
type VPCPeeringResource struct {
	transport *ones_gfx.Transport
}

// NewVPCPeeringResource constructs a VPCPeeringResource.
func NewVPCPeeringResource(transport *ones_gfx.Transport) *VPCPeeringResource {
	return &VPCPeeringResource{transport: transport}
}

// peeringPayload mirrors VpcPeeringRequest.java. name/vpcname/peervpcname are
// required server-side (Java 400s without them), plus optional webhook fields.
type peeringPayload struct {
	Name          string   `json:"name"`
	Vpcname       string   `json:"vpcname"`
	Peervpcname   string   `json:"peervpcname"`
	EnableWebhook bool     `json:"enableWebhook,omitempty"`
	WebhookURL    string   `json:"webhookUrl,omitempty"`
	WebhookEvents []string `json:"webhookEvents,omitempty"`
}

func (r *VPCPeeringResource) buildPayload(name, vpcname, peervpcname string, rc ones_gfx.ResolvedCall) peeringPayload {
	body := peeringPayload{Name: name, Vpcname: vpcname, Peervpcname: peervpcname}
	if rc.WebhookURL != "" {
		body.EnableWebhook = true
		body.WebhookURL = rc.WebhookURL
		body.WebhookEvents = rc.WebhookEvents
	}
	return body
}

// Create enables route leaking between two VPCs (synchronous).
// Maps to POST /fabrics/{fabricName}/vpcpeering.
func (r *VPCPeeringResource) Create(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...ones_gfx.CallOption) (string, error) {
	if err := firstErr(pathName("fabricName", fabricName), peeringFields(name, vpcname, peervpcname)); err != nil {
		return "", err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := r.buildPayload(name, vpcname, peervpcname, rc)
	res, err := ones_gfx.Call[string](r.transport, "POST", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// CreateAsync enables route leaking and returns the 202 operation handle.
func (r *VPCPeeringResource) CreateAsync(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	if err := firstErr(pathName("fabricName", fabricName), peeringFields(name, vpcname, peervpcname)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := r.buildPayload(name, vpcname, peervpcname, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "POST", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// deletePeeringPayload is the DELETE body. FM addresses a peering by name and
// rejects the call unless vpcname/peervpcname match the pair recorded for it.
type deletePeeringPayload struct {
	Name        string `json:"name"`
	Vpcname     string `json:"vpcname"`
	Peervpcname string `json:"peervpcname"`
}

// Delete removes route leaking between two VPCs (synchronous only; FM has no
// async path for this route). It is the inverse of Create.
// Maps to DELETE /fabrics/{fabricName}/vpcpeering.
//
// FM answers 404 VPC_PEERING_NOT_FOUND for an unknown peering name (or
// FABRIC_NOT_FOUND for an unknown fabric), 400 VPC_PEERING_PAIR_MISMATCH when
// the VPC pair is not the one recorded for name, and 500
// VPC_PEERING_DELETE_FAILED when a switch could not be updated. These surface
// as *NotFoundError, *BadRequestError and *ServerError respectively, with the
// full envelope (including error.code) in ResponseBody.
func (r *VPCPeeringResource) Delete(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	if err := firstErr(pathName("fabricName", fabricName), peeringFields(name, vpcname, peervpcname)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := deletePeeringPayload{Name: name, Vpcname: vpcname, Peervpcname: peervpcname}
	return ones_gfx.Call[ApiResponseMessage](r.transport, "DELETE", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// VPCPeeringOrigin says which API created a peering.
type VPCPeeringOrigin string

const (
	// VPCPeeringOriginVPCPeering is a peering created by VPCPeering.Create.
	VPCPeeringOriginVPCPeering VPCPeeringOrigin = "VPC_PEERING"
	// VPCPeeringOriginNetworkPolicy is a peering created by a tenant-to-tenant
	// network policy on a VXLAN fabric. FM refuses to delete it through
	// VPCPeering.Delete (409 VPC_PEERING_OWNED_BY_POLICY); remove the policy
	// instead.
	VPCPeeringOriginNetworkPolicy VPCPeeringOrigin = "NETWORK_POLICY"
)

// VPCPeering is one peering as the read API reports it.
//
// The fabric-wide calls (List, Get) report the pair as Vpcs, unordered,
// because the route leak is configured on both VRFs. ListForTenant knows
// which side is the tenant's, so it sets Vpc (the tenant's VPC) and PeerVpc
// instead, and leaves Vpcs empty.
type VPCPeering struct {
	Name      string             `json:"name"`
	Vpcs      []string           `json:"vpcs,omitempty"`
	Vpc       string             `json:"vpc,omitempty"`
	PeerVpc   string             `json:"peerVpc,omitempty"`
	Origin    VPCPeeringOrigin   `json:"origin,omitempty"`
	CreatedAt ones_gfx.Timestamp `json:"createdAt,omitempty"`
}

// ListPagination is the page block FM returns with a paged list.
type ListPagination struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
}

// VPCPeeringList is one page of peerings. TenantName and Vpc are set only by
// ListForTenant.
type VPCPeeringList struct {
	FabricName string         `json:"fabricName"`
	TenantName string         `json:"tenantName,omitempty"`
	Vpc        string         `json:"vpc,omitempty"`
	Peerings   []VPCPeering   `json:"peerings"`
	Pagination ListPagination `json:"pagination"`
}

// peeringPageQuery builds the page/size query. page is zero-based and omitted
// when negative; size is omitted when not positive (FM defaults to 50 and
// caps it at 500).
func peeringPageQuery(page, size int) url.Values {
	q := url.Values{}
	if page >= 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}
	return q
}

// List returns one page of a fabric's peerings, oldest first.
// Maps to GET /fabrics/{fabricName}/vpcpeering.
//
// A fabric with no peerings returns an empty list. Peerings made by a
// VRF-LITE network policy are not recorded here; NetworkPolicy.List is their
// system of record.
func (r *VPCPeeringResource) List(ctx context.Context, fabricName string, page, size int, opts ...ones_gfx.CallOption) (*VPCPeeringList, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	ro.Query = peeringPageQuery(page, size)
	return ones_gfx.Call[VPCPeeringList](r.transport, "GET", "fabrics/"+fabricName+"/vpcpeering", nil, ones_gfx.OperationModeSynchronous, ro)
}

// Get returns one peering by name.
// Maps to GET /fabrics/{fabricName}/vpcpeering/{peeringName}.
//
// An unknown peering returns *NotFoundError (VPC_PEERING_NOT_FOUND).
func (r *VPCPeeringResource) Get(ctx context.Context, fabricName, peeringName string, opts ...ones_gfx.CallOption) (*VPCPeering, error) {
	if err := firstErr(pathName("fabricName", fabricName), pathName("peeringName", peeringName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[struct {
		Peering VPCPeering `json:"peering"`
	}](r.transport, "GET", "fabrics/"+fabricName+"/vpcpeering/"+peeringName, nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return &res.Peering, nil
}

// ListForTenant returns one page of the peerings a tenant takes part in,
// including its storage peering. Each entry reports the tenant's own Vpc and
// the PeerVpc on the other side.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/vpcpeering.
//
// An unknown tenant returns *NotFoundError (TENANT_NOT_FOUND).
func (r *VPCPeeringResource) ListForTenant(ctx context.Context, fabricName, tenantName string, page, size int, opts ...ones_gfx.CallOption) (*VPCPeeringList, error) {
	if err := fabricTenant(fabricName, tenantName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	ro.Query = peeringPageQuery(page, size)
	return ones_gfx.Call[VPCPeeringList](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/vpcpeering", nil, ones_gfx.OperationModeSynchronous, ro)
}
