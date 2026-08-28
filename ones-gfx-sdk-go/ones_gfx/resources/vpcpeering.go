package resources

import (
	"context"

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
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := r.buildPayload(name, vpcname, peervpcname, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "POST", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// Delete removes route leaking between two VPCs.
// Maps to DELETE /fabrics/{fabricName}/vpcpeering.
func (r *VPCPeeringResource) Delete(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...ones_gfx.CallOption) (string, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := r.buildPayload(name, vpcname, peervpcname, rc)
	res, err := ones_gfx.Call[string](r.transport, "DELETE", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}
