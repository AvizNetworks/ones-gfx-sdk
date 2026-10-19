package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// RailsResource covers the rail CIDRPool manifests for a fabric's East-West
// rails.
type RailsResource struct {
	transport *ones_gfx.Transport
}

// NewRailsResource constructs a RailsResource.
func NewRailsResource(transport *ones_gfx.Transport) *RailsResource {
	return &RailsResource{transport: transport}
}

// Media types FM will render a rail CIDRPool manifest as. Ask for one of the
// YAML types through the *YAML methods when feeding a tool that consumes the
// manifest text.
const (
	RailYamlMediaType       = "application/yaml"
	RailLegacyYamlMediaType = "application/x-yaml"
	RailTextMediaType       = "text/plain"
)

// RailCidrPools is the manifest for every East-West rail on a fabric. The
// document shape follows the CIDRPool manifest rather than a fixed schema, so
// it is surfaced as a raw map.
type RailCidrPools = map[string]interface{}

// CidrPools returns the rail CIDRPool manifest for every East-West rail, as
// JSON.
// Maps to GET /fabrics/{fabricName}/rails/cidrpools.
//
// 404 INVALID_FABRIC when the fabric is unknown, or RAIL_CIDRPOOL_NOT_FOUND
// when it has no East-West rails.
func (r *RailsResource) CidrPools(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (RailCidrPools, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[RailCidrPools](r.transport, "GET", "fabrics/"+fabricName+"/rails/cidrpools", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// CidrPoolsYAML returns the same manifest as raw text, rendered as the media
// type given — one of RailYamlMediaType, RailLegacyYamlMediaType or
// RailTextMediaType. Pass "" for RailYamlMediaType.
// Maps to GET /fabrics/{fabricName}/rails/cidrpools.
func (r *RailsResource) CidrPoolsYAML(ctx context.Context, fabricName, mediaType string, opts ...ones_gfx.CallOption) (string, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	if mediaType == "" {
		mediaType = RailYamlMediaType
	}
	ro.Accept = mediaType
	res, err := ones_gfx.Call[string](r.transport, "GET", "fabrics/"+fabricName+"/rails/cidrpools", nil, ones_gfx.OperationModeSynchronous, ro)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// CidrPool returns the CIDRPool document for one rail, as JSON. railName is
// the rail's CIDRPool name.
// Maps to GET /fabrics/{fabricName}/rails/{railName}/cidrpool.
//
// 404 INVALID_FABRIC when the fabric is unknown, or INVALID_RAIL when the
// rail is not one of its rails.
func (r *RailsResource) CidrPool(ctx context.Context, fabricName, railName string, opts ...ones_gfx.CallOption) (RailCidrPools, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[RailCidrPools](r.transport, "GET", "fabrics/"+fabricName+"/rails/"+railName+"/cidrpool", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// CidrPoolYAML returns one rail's document as raw text, rendered as the media
// type given. Pass "" for RailYamlMediaType.
// Maps to GET /fabrics/{fabricName}/rails/{railName}/cidrpool.
func (r *RailsResource) CidrPoolYAML(ctx context.Context, fabricName, railName, mediaType string, opts ...ones_gfx.CallOption) (string, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	if mediaType == "" {
		mediaType = RailYamlMediaType
	}
	ro.Accept = mediaType
	res, err := ones_gfx.Call[string](r.transport, "GET", "fabrics/"+fabricName+"/rails/"+railName+"/cidrpool", nil, ones_gfx.OperationModeSynchronous, ro)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}
