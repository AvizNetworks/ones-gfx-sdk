package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FabricsResource covers fabric CRUD, listing and device discovery endpoints.
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

// fabricPayload is the Fabrics request body shared by Create and Update
// (Models/Fabrics.java). All fields optional except name on create.
type fabricPayload struct {
	ID                    *int    `json:"id,omitempty"`
	Name                  *string `json:"name,omitempty"`
	Type                  *string `json:"type,omitempty"`
	Status                *string `json:"status,omitempty"`
	Description           *string `json:"description,omitempty"`
	OrchestrationStatus   *string `json:"orchestrationStatus,omitempty"`
	NumOfSus              *int    `json:"numOfSus,omitempty"`
	MaxNumOfSus           *int    `json:"maxNumOfSus,omitempty"`
	Dedicated             *bool   `json:"dedicated,omitempty"`
	Hybrid                *bool   `json:"hybrid,omitempty"`
	IsDPUFabric           *bool   `json:"isDPUFabric,omitempty"`
	StartingSubnetGpu     *int    `json:"startingSubnetGpu,omitempty"`
	StartingSubnetCpu     *string `json:"startingSubnetCpu,omitempty"`
	StartingSubnetTenants *string `json:"startingSubnetTenants,omitempty"`
	StartingSubnetStorage *string `json:"startingSubnetStorage,omitempty"`
	SimulationID          *int    `json:"simulationId,omitempty"`
	Intent                *string `json:"intent,omitempty"`
	EWTenantAware         *bool   `json:"ewTenantAware,omitempty"`
	StorageTenantAware    *bool   `json:"storageTenantAware,omitempty"`
	IsOnesControlled      *bool   `json:"isOnesControlled,omitempty"`
	SuHostCnt             *string `json:"suHostCnt,omitempty"`
	IsVxlanFabric         *bool   `json:"isVxlanFabric,omitempty"`
	NodeType              *string `json:"nodeType,omitempty"`
	DeploymentType        *string `json:"deploymentType,omitempty"`
	SpineEvpnConfigured   *bool   `json:"spineEvpnConfigured,omitempty"`
	Isimported            *bool   `json:"isimported,omitempty"`
	GpuScaleMode          *string `json:"gpuScaleMode,omitempty"`
	Cnpq                  *string `json:"cnpq,omitempty"`
	UfmURL                *string `json:"ufmUrl,omitempty"`
	UfmUsername           *string `json:"ufmUsername,omitempty"`
	UfmPasswordEncrypted  *string `json:"ufmPasswordEncrypted,omitempty"`
}

// FabricCreateArgs carries the optional fabric fields for Create/Update.
// Zero-value nils are omitted from the request body.
type FabricCreateArgs struct {
	ID                    *int
	Type                  *string
	Status                *string
	Description           *string
	OrchestrationStatus   *string
	NumOfSus              *int
	MaxNumOfSus           *int
	Dedicated             *bool
	Hybrid                *bool
	IsDPUFabric           *bool
	StartingSubnetGpu     *int
	StartingSubnetCpu     *string
	StartingSubnetTenants *string
	StartingSubnetStorage *string
	SimulationID          *int
	Intent                *string
	EWTenantAware         *bool
	StorageTenantAware    *bool
	IsOnesControlled      *bool
	SuHostCnt             *string
	IsVxlanFabric         *bool
	NodeType              *string
	DeploymentType        *string
	SpineEvpnConfigured   *bool
	Isimported            *bool
	GpuScaleMode          *string
	Cnpq                  *string
	UfmURL                *string
	UfmUsername           *string
	UfmPasswordEncrypted  *string
}

func (a *FabricCreateArgs) payload(name *string) *fabricPayload {
	return &fabricPayload{
		ID: a.ID, Name: name, Type: a.Type, Status: a.Status, Description: a.Description,
		OrchestrationStatus: a.OrchestrationStatus, NumOfSus: a.NumOfSus, MaxNumOfSus: a.MaxNumOfSus,
		Dedicated: a.Dedicated, Hybrid: a.Hybrid, IsDPUFabric: a.IsDPUFabric,
		StartingSubnetGpu: a.StartingSubnetGpu, StartingSubnetCpu: a.StartingSubnetCpu,
		StartingSubnetTenants: a.StartingSubnetTenants, StartingSubnetStorage: a.StartingSubnetStorage,
		SimulationID: a.SimulationID, Intent: a.Intent, EWTenantAware: a.EWTenantAware,
		StorageTenantAware: a.StorageTenantAware, IsOnesControlled: a.IsOnesControlled,
		SuHostCnt: a.SuHostCnt, IsVxlanFabric: a.IsVxlanFabric, NodeType: a.NodeType,
		DeploymentType: a.DeploymentType, SpineEvpnConfigured: a.SpineEvpnConfigured,
		Isimported: a.Isimported, GpuScaleMode: a.GpuScaleMode, Cnpq: a.Cnpq,
		UfmURL: a.UfmURL, UfmUsername: a.UfmUsername, UfmPasswordEncrypted: a.UfmPasswordEncrypted,
	}
}

// Create adds a new fabric. Maps to POST /addFabricData. Returns the
// "Fabric Added" success message.
func (r *FabricsResource) Create(ctx context.Context, name string, args *FabricCreateArgs) (string, error) {
	if args == nil {
		args = &FabricCreateArgs{}
	}
	res, err := ones_gfx.Call[string](r.transport, "POST", "addFabricData", args.payload(&name), ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Update partially updates a fabric. Maps to PATCH /fabrics/{fabricName}/editFabricData.
func (r *FabricsResource) Update(ctx context.Context, fabricName string, args *FabricCreateArgs) (string, error) {
	if args == nil {
		args = &FabricCreateArgs{}
	}
	res, err := ones_gfx.Call[string](r.transport, "PATCH", "fabrics/"+fabricName+"/editFabricData", args.payload(nil), ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Delete removes a fabric by name. Maps to DELETE /delFabricData/{name}.
func (r *FabricsResource) Delete(ctx context.Context, name string) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "DELETE", "delFabricData/"+name, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// List returns all fabrics. Maps to GET /getAllFabrics.
func (r *FabricsResource) List(ctx context.Context) ([]ones_gfx.FabricItem, error) {
	res, err := ones_gfx.Call[[]ones_gfx.FabricItem](r.transport, "GET", "getAllFabrics", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Get returns one fabric by name. Maps to GET /getFabricByName/{name}.
func (r *FabricsResource) Get(ctx context.Context, name string) (*ones_gfx.FabricItem, error) {
	return ones_gfx.Call[ones_gfx.FabricItem](r.transport, "GET", "getFabricByName/"+name, nil, ones_gfx.OperationModeSynchronous, nil)
}

// ListDTOs returns the FabricDto projection. Maps to GET /fabrics.
func (r *FabricsResource) ListDTOs(ctx context.Context) (*FabricsListResponse, error) {
	return ones_gfx.Call[FabricsListResponse](r.transport, "GET", "fabrics", nil, ones_gfx.OperationModeSynchronous, nil)
}

// DeviceIPs lists device IP addresses of a fabric.
// Maps to GET /getFabricDevices/{fabricName}.
func (r *FabricsResource) DeviceIPs(ctx context.Context, fabricName string) ([]string, error) {
	res, err := ones_gfx.Call[[]string](r.transport, "GET", "getFabricDevices/"+fabricName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// DevicesByLayer lists device IPs for a layer ("spine", "leaf", "tor", "host").
// Maps to GET /getDevicesByLayer?layer=...
func (r *FabricsResource) DevicesByLayer(ctx context.Context, layer string) ([]string, error) {
	q := url.Values{}
	q.Set("layer", layer)
	res, err := ones_gfx.Call[[]string](r.transport, "GET", "getDevicesByLayer", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// UpdateStatus updates a fabric's operational status.
// Maps to POST /updateFabricStatus.
func (r *FabricsResource) UpdateStatus(ctx context.Context, name string, status, intent, description, deploymentType *string) (string, error) {
	body := struct {
		Name           string  `json:"name"`
		Status         *string `json:"status,omitempty"`
		Intent         *string `json:"intent,omitempty"`
		Description    *string `json:"description,omitempty"`
		DeploymentType *string `json:"deploymentType,omitempty"`
	}{name, status, intent, description, deploymentType}
	res, err := ones_gfx.Call[string](r.transport, "POST", "updateFabricStatus", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}
