package resources

import (
	"context"
	"fmt"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// RailsResource covers the rail CIDRPool manifest endpoints under /api/fm.
type RailsResource struct {
	transport *ones_gfx.Transport
}

// NewRailsResource constructs a RailsResource.
func NewRailsResource(transport *ones_gfx.Transport) *RailsResource {
	return &RailsResource{transport: transport}
}

// RailAccept selects the representation of a rail CIDRPool response.
type RailAccept string

const (
	// RailAcceptYAML returns the manifest verbatim (the default).
	RailAcceptYAML RailAccept = "application/yaml"
	// RailAcceptXYAML is the legacy YAML media type; same body as RailAcceptYAML.
	RailAcceptXYAML RailAccept = "application/x-yaml"
	// RailAcceptText returns the manifest verbatim as text/plain.
	RailAcceptText RailAccept = "text/plain"
	// RailAcceptJSON returns the standard JSON envelope, which also lists the rails.
	RailAcceptJSON RailAccept = "application/json"
)

func (a RailAccept) valid() bool {
	switch a {
	case RailAcceptYAML, RailAcceptXYAML, RailAcceptText, RailAcceptJSON:
		return true
	}
	return false
}

func (a RailAccept) raw() bool { return a != RailAcceptJSON }

// RailCidrPools is the rail CIDRPool manifest of a fabric: one multi-document
// YAML string with a CIDRPool document per rail.
//
// YAML is always set. Rails (the rail names) is only populated for
// RailAcceptJSON, because the server only lists them in the JSON envelope.
type RailCidrPools struct {
	FabricName string   `json:"fabricName"`
	Rails      []string `json:"rails,omitempty"`
	YAML       string   `json:"yaml"`
}

// RailCidrPool is the CIDRPool document of a single rail. RailName is set in
// every mode (it is the rail that was asked for).
type RailCidrPool struct {
	FabricName string `json:"fabricName"`
	RailName   string `json:"railName"`
	YAML       string `json:"yaml"`
}

func resolveRailAccept(accept *RailAccept) (RailAccept, error) {
	if accept == nil {
		return RailAcceptYAML, nil
	}
	if !accept.valid() {
		return "", &ones_gfx.ValidationError{
			Field:  "accept",
			Reason: fmt.Sprintf("must be one of %q, %q, %q, %q; got %q", RailAcceptYAML, RailAcceptXYAML, RailAcceptText, RailAcceptJSON, string(*accept)),
		}
	}
	return *accept, nil
}

// CidrPools returns the rail CIDRPool manifest of a fabric. A nil accept means
// RailAcceptYAML.
// Maps to GET /fabrics/{fabricName}/rails/cidrpools.
//
// Errors: *ones_gfx.NotFoundError for an unknown fabric (INVALID_FABRIC) or a
// fabric without East-West rails (RAIL_CIDRPOOL_NOT_FOUND). For the raw
// representations the server's error body is plain text "<CODE>: <detail>",
// which becomes the error message; ResponseBody holds the same text.
func (r *RailsResource) CidrPools(ctx context.Context, fabricName string, accept *RailAccept) (*RailCidrPools, error) {
	segs, err := pathSegments("fabricName", fabricName)
	if err != nil {
		return nil, err
	}
	a, err := resolveRailAccept(accept)
	if err != nil {
		return nil, err
	}
	path := "fabrics/" + segs[0] + "/rails/cidrpools"
	ro := &ones_gfx.ReqOpts{Accept: string(a)}
	if a.raw() {
		text, err := ones_gfx.CallText(r.transport, "GET", path, ro)
		if err != nil {
			return nil, err
		}
		return &RailCidrPools{FabricName: fabricName, YAML: text}, nil
	}
	return ones_gfx.Call[RailCidrPools](r.transport, "GET", path, nil, ones_gfx.OperationModeSynchronous, ro)
}

// CidrPool returns the CIDRPool document of one rail (e.g. "r0"). A nil accept
// means RailAcceptYAML.
// Maps to GET /fabrics/{fabricName}/rails/{railName}/cidrpool.
//
// Errors: *ones_gfx.NotFoundError for an unknown fabric (INVALID_FABRIC) or
// rail (INVALID_RAIL); see CidrPools for the error body format.
func (r *RailsResource) CidrPool(ctx context.Context, fabricName, railName string, accept *RailAccept) (*RailCidrPool, error) {
	segs, err := pathSegments("fabricName", fabricName, "railName", railName)
	if err != nil {
		return nil, err
	}
	a, err := resolveRailAccept(accept)
	if err != nil {
		return nil, err
	}
	path := "fabrics/" + segs[0] + "/rails/" + segs[1] + "/cidrpool"
	ro := &ones_gfx.ReqOpts{Accept: string(a)}
	if a.raw() {
		text, err := ones_gfx.CallText(r.transport, "GET", path, ro)
		if err != nil {
			return nil, err
		}
		return &RailCidrPool{FabricName: fabricName, RailName: railName, YAML: text}, nil
	}
	return ones_gfx.Call[RailCidrPool](r.transport, "GET", path, nil, ones_gfx.OperationModeSynchronous, ro)
}
