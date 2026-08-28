package ones_gfx

// apitypes.go holds request/response types shared by more than one resource
// file, mirroring ones_gfx/_types.py in the Python SDK. Types used by a
// single resource are defined in that resource's file instead.

// GpuAction is the ADD/DELETE discriminator shared by GPU allocation,
// tenant-update and NMX-C domain endpoints (Java Helper/Enums.java).
type GpuAction = GPUOperation

const (
	GpuActionAdd    = OperationAdd
	GpuActionDelete = OperationDelete
)

// OperationAccepted is the 202 Accepted body of async-capable endpoints
// (ApiResponse.success(msg, asyncData) with the envelope unwrapped).
// OperationType appears only on GPU allocate/deallocate operations;
// WebhookRegistered is true only in ASYNC_WEBHOOK mode.
type OperationAccepted struct {
	OperationID       string `json:"operationId,omitempty"`
	Status            string `json:"status,omitempty"`
	OperationType     string `json:"operationType,omitempty"`
	WebhookRegistered bool   `json:"webhookRegistered,omitempty"`
}

// FabricItem mirrors the Fabrics entity JSON (Models/Fabrics.java) returned
// by GET /getAllFabrics and GET /getFabricByName/{name}.
type FabricItem struct {
	ID                    int    `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Type                  string `json:"type,omitempty"`
	Status                string `json:"status,omitempty"`
	Description           string `json:"description,omitempty"`
	OrchestrationStatus   string `json:"orchestrationStatus,omitempty"`
	NumOfSus              int    `json:"numOfSus,omitempty"`
	MaxNumOfSus           int    `json:"maxNumOfSus,omitempty"`
	Dedicated             *bool  `json:"dedicated,omitempty"`
	Hybrid                *bool  `json:"hybrid,omitempty"`
	IsDPUFabric           *bool  `json:"isDPUFabric,omitempty"`
	StartingSubnetGpu     *int   `json:"startingSubnetGpu,omitempty"`
	StartingSubnetCpu     string `json:"startingSubnetCpu,omitempty"`
	StartingSubnetTenants string `json:"startingSubnetTenants,omitempty"`
	StartingSubnetStorage string `json:"startingSubnetStorage,omitempty"`
	SimulationID          *int   `json:"simulationId,omitempty"`
	Intent                string `json:"intent,omitempty"`
	EWTenantAware         *bool  `json:"ewTenantAware,omitempty"`
	StorageTenantAware    *bool  `json:"storageTenantAware,omitempty"`
	IsOnesControlled      *bool  `json:"isOnesControlled,omitempty"`
	SuHostCnt             string `json:"suHostCnt,omitempty"`
	IsVxlanFabric         *bool  `json:"isVxlanFabric,omitempty"`
	NodeType              string `json:"nodeType,omitempty"`
	DeploymentType        string `json:"deploymentType,omitempty"`
	SpineEvpnConfigured   *bool  `json:"spineEvpnConfigured,omitempty"`
	Isimported            *bool  `json:"isimported,omitempty"`
	GpuScaleMode          string `json:"gpuScaleMode,omitempty"`
	Cnpq                  string `json:"cnpq,omitempty"`
	UfmURL                string `json:"ufmUrl,omitempty"`
	UfmUsername           string `json:"ufmUsername,omitempty"`
	UfmPasswordEncrypted  string `json:"ufmPasswordEncrypted,omitempty"`
	CreatedAt             string `json:"createdAt,omitempty"`
	UpdatedAt             string `json:"updatedAt,omitempty"`
}
