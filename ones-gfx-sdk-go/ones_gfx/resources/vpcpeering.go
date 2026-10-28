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
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := deletePeeringPayload{Name: name, Vpcname: vpcname, Peervpcname: peervpcname}
	return ones_gfx.Call[ApiResponseMessage](r.transport, "DELETE", "fabrics/"+fabricName+"/vpcpeering", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}
