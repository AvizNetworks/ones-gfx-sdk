package resources

import (
	"context"

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
