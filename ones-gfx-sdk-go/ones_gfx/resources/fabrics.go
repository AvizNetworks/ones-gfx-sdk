package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FabricsResource covers the fabric listing endpoint (GET /fabrics).
type FabricsResource struct {
	transport *ones_gfx.Transport
}

// NewFabricsResource constructs a FabricsResource.
func NewFabricsResource(transport *ones_gfx.Transport) *FabricsResource {
	return &FabricsResource{transport: transport}
}

// FabricDtoItem mirrors Cumulus/dto/FabricDto.java (GET /fabrics entries).
type FabricDtoItem struct {
	ID                 int    `json:"id,omitempty"`
	FabricName         string `json:"fabricName,omitempty"`
	Description        string `json:"description,omitempty"`
	NumOfSUs           int    `json:"numOfSUs,omitempty"`
	MaxNumOfSUs        int    `json:"maxNumOfSUs,omitempty"`
	EWTenantAware      bool   `json:"ewTenantAware,omitempty"`
	NSTenantAware      bool   `json:"nsTenantAware,omitempty"`
	DefaultStorageName string `json:"defaultStorageName,omitempty"`
	Cnpq               string `json:"cnpq,omitempty"`
	CreatedAt          string `json:"createdAt,omitempty"`
	UpdatedAt          string `json:"updatedAt,omitempty"`
}

// FabricsListResponse is the GET /fabrics body: {"fabrics": [...]}.
type FabricsListResponse struct {
	Fabrics []FabricDtoItem `json:"fabrics"`
}

// ListDTOs returns the FabricDto projection. Maps to GET /fabrics.
func (r *FabricsResource) ListDTOs(ctx context.Context) (*FabricsListResponse, error) {
	return ones_gfx.Call[FabricsListResponse](r.transport, "GET", "fabrics", nil, ones_gfx.OperationModeSynchronous, nil)
}

// TopologyYAML returns the fabric's topology YAML (the file it was deployed
// from).
// Maps to GET /fabrics/{fabricName} on the ONES gateway, served at the root
// rather than under /api/fm. The gateway answers 404 NOT_FOUND when it has no
// topology for the fabric.
func (r *FabricsResource) TopologyYAML(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (string, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return "", err
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	ro := rc.ReqOpts()
	ro.Accept = "application/x-yaml, application/json"
	// The leading "/" resolves against the host root, not the /api/fm/ base.
	res, err := ones_gfx.Call[string](r.transport, "GET", "/fabrics/"+url.PathEscape(fabricName), nil, ones_gfx.OperationModeSynchronous, ro)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// NSInterfaces returns each host's north-south interfaces on a fabric: it
// fetches TopologyYAML and reduces it with ones_gfx.ExtractNSInterfaces.
// opts.FabricName defaults to fabricName; opts.LeafPatterns to
// ones_gfx.DefaultNSLeafPatterns (^leaf-cn, ^leaf-fn, ^leaf-sn).
// Render the result with its YAML method.
func (r *FabricsResource) NSInterfaces(ctx context.Context, fabricName string, opts ones_gfx.NSInterfaceOptions, callOpts ...ones_gfx.CallOption) (*ones_gfx.NSInterfaces, error) {
	text, err := r.TopologyYAML(ctx, fabricName, callOpts...)
	if err != nil {
		return nil, err
	}
	if opts.FabricName == "" {
		opts.FabricName = fabricName
	}
	return ones_gfx.ExtractNSInterfaces([]byte(text), opts)
}
