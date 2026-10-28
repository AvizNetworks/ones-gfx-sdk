package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// NetworkPolicyResource covers the fabric's network policies: which source is
// allowed to reach which destination. Policies apply on VXLAN (Cumulus, SONiC
// and DPU) and VRF-LITE (Cumulus and SONiC) fabrics alike.
type NetworkPolicyResource struct {
	transport *ones_gfx.Transport
}

// NewNetworkPolicyResource constructs a NetworkPolicyResource.
func NewNetworkPolicyResource(transport *ones_gfx.Transport) *NetworkPolicyResource {
	return &NetworkPolicyResource{transport: transport}
}

// Network policy endpoint types. FM matches these case-insensitively.
const (
	NetworkPolicyTypeTenant  = "tenant"
	NetworkPolicyTypeService = "service"
)

// MaxNetworkPolicyNameLength is the longest policy name FM accepts.
const MaxNetworkPolicyNameLength = 24

// NetworkPolicySide is one end of a policy — the source or the destination.
// ServiceName applies when Type is "service". Endpoints is empty on a
// tenant-to-tenant policy, which carries none.
type NetworkPolicySide struct {
	Type        string   `json:"type"`
	TenantName  string   `json:"tenantName"`
	ServiceName string   `json:"serviceName,omitempty"`
	Endpoints   []string `json:"endpoints,omitempty"`
}

// NetworkPolicySpec pairs the two ends of a policy.
type NetworkPolicySpec struct {
	Source      NetworkPolicySide `json:"source"`
	Destination NetworkPolicySide `json:"destination"`
}

// NetworkPolicyName carries the policy's own name.
type NetworkPolicyName struct {
	Name string `json:"name"`
}

// NetworkPolicy is a policy as FM stores and returns it.
type NetworkPolicy struct {
	Policy NetworkPolicyName `json:"policy"`
	Spec   NetworkPolicySpec `json:"spec"`
}

// NetworkPolicyList is the body of a List call. It is not wrapped in the
// usual status/message envelope.
type NetworkPolicyList struct {
	Status     string          `json:"status"`
	FabricName string          `json:"fabricName"`
	Count      int             `json:"count"`
	Policies   []NetworkPolicy `json:"policies"`
}

// NetworkPolicyRef names the policy a mutation acted on.
type NetworkPolicyRef struct {
	PolicyName string `json:"policyName"`
}

// List returns every network policy on the fabric.
// Maps to GET /fabrics/{fabricName}/networkpolicy.
// An unknown fabric answers 404.
func (r *NetworkPolicyResource) List(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (*NetworkPolicyList, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[NetworkPolicyList](r.transport, "GET", "fabrics/"+fabricName+"/networkpolicy", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Create adds a network policy allowing traffic from source to destination.
// Maps to POST /fabrics/{fabricName}/networkpolicy.
//
// policyName is at most MaxNetworkPolicyNameLength characters. Config is
// pushed to the border/combined leaves and DPUs, so FM runs a device licence
// check first and answers 403 when it fails.
func (r *NetworkPolicyResource) Create(ctx context.Context, fabricName, policyName string, source, destination NetworkPolicySide, opts ...ones_gfx.CallOption) (*NetworkPolicyRef, error) {
	if err := firstErr(pathName("fabricName", fabricName), newPolicyName(policyName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := NetworkPolicy{
		Policy: NetworkPolicyName{Name: policyName},
		Spec:   NetworkPolicySpec{Source: source, Destination: destination},
	}
	return ones_gfx.Call[NetworkPolicyRef](r.transport, "POST", "fabrics/"+fabricName+"/networkpolicy", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Patch adds endpoints to an existing policy, leaving the endpoints already
// configured on the devices untouched.
// Maps to PATCH /fabrics/{fabricName}/networkpolicy/{policyName}.
//
// Only endpoints may change: Type, TenantName and ServiceName on each side
// must match the existing policy exactly, and the Endpoints given are the NEW
// ones to add rather than the full set. A tenant-to-tenant policy has no
// endpoints and is rejected.
func (r *NetworkPolicyResource) Patch(ctx context.Context, fabricName, policyName string, source, destination NetworkPolicySide, opts ...ones_gfx.CallOption) (*NetworkPolicyRef, error) {
	if err := firstErr(pathName("fabricName", fabricName), pathName("policyName", policyName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	// policy.name is echoed into the body as well as the path: FM validates a
	// patch with the same check it uses for a create, which rejects a body
	// without it ("policy.name is required.").
	body := NetworkPolicy{
		Policy: NetworkPolicyName{Name: policyName},
		Spec:   NetworkPolicySpec{Source: source, Destination: destination},
	}
	return ones_gfx.Call[NetworkPolicyRef](r.transport, "PATCH", "fabrics/"+fabricName+"/networkpolicy/"+policyName, body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Delete removes a policy and its config from the devices it was pushed to.
// Maps to DELETE /fabrics/{fabricName}/networkpolicy/{policyName}.
func (r *NetworkPolicyResource) Delete(ctx context.Context, fabricName, policyName string, opts ...ones_gfx.CallOption) (*NetworkPolicyRef, error) {
	if err := firstErr(pathName("fabricName", fabricName), pathName("policyName", policyName)); err != nil {
		return nil, err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[NetworkPolicyRef](r.transport, "DELETE", "fabrics/"+fabricName+"/networkpolicy/"+policyName, nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}
