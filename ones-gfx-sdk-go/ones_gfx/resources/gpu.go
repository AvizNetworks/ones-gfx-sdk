package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// GPUResource covers GPU inventory, mapping and port-assignment endpoints.
type GPUResource struct {
	transport *ones_gfx.Transport
}

// NewGPUResource constructs a GPUResource.
func NewGPUResource(transport *ones_gfx.Transport) *GPUResource {
	return &GPUResource{transport: transport}
}

// GpuItem mirrors Models/Gpus.java JSON.
type GpuItem struct {
	ID                   int    `json:"id,omitempty"`
	GpuPort              string `json:"gpuPort,omitempty"`
	GpuHostname          string `json:"gpuHostname,omitempty"`
	GpuStatus            string `json:"gpuStatus,omitempty"`
	GpuLeafInt           string `json:"gpuLeafInt,omitempty"`
	SuID                 *int   `json:"suId,omitempty"`
	LeafHostname         string `json:"leafHostname,omitempty"`
	LeafIPAddress        string `json:"leafIpAddress,omitempty"`
	TenantName           string `json:"tenantName,omitempty"`
	FabricName           string `json:"fabricName,omitempty"`
	ConfigStatus         string `json:"config_status,omitempty"`
	LastConfiguredTenant string `json:"lastConfiguredTenant,omitempty"`
	CreatedAt            string `json:"createdAt,omitempty"`
	UpdatedAt            string `json:"updatedAt,omitempty"`
}

// GpuTenantMappingItem mirrors Models/GpuTenantMapping.java JSON.
type GpuTenantMappingItem struct {
	ID             int    `json:"id,omitempty"`
	FabricName     string `json:"fabricName,omitempty"`
	ServerName     string `json:"serverName,omitempty"`
	GPUIndex       *int   `json:"gpuIndex,omitempty"`
	LogicalGpuName string `json:"logicalGpuName,omitempty"`
	TenantName     string `json:"tenantName,omitempty"`
	ConfigStatus   string `json:"configStatus,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
}

// GpuAllocationHistoryItem mirrors Models/GpuAllocationHistory.java JSON.
type GpuAllocationHistoryItem struct {
	ID          int    `json:"id,omitempty"`
	SuNumber    *int   `json:"suNumber,omitempty"`
	HostName    string `json:"hostName,omitempty"`
	GPUsAdded   string `json:"gpusAdded,omitempty"`
	GPUsRemoved string `json:"gpusRemoved,omitempty"`
	TenantName  string `json:"tenantName,omitempty"`
	FabricName  string `json:"fabricName,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// AvailableServerResult mirrors Cumulus/dto/AvailableServer.java.
type AvailableServerResult struct {
	AvailableGPUs []string `json:"availableGPUs,omitempty"`
}

// GpuPortAssignmentResult is the response of AssignPorts
// (POST /fabrics/{fabricName}/tenants/{tenantName}/gpus).
type GpuPortAssignmentResult struct {
	FabricName       string   `json:"fabricName,omitempty"`
	TenantName       string   `json:"tenantName,omitempty"`
	Servers          []string `json:"servers,omitempty"`
	Success          bool     `json:"success,omitempty"`
	Error            string   `json:"error,omitempty"`
	Message          string   `json:"message,omitempty"`
	ServersProcessed int      `json:"serversProcessed,omitempty"`
	Operation        string   `json:"operation,omitempty"`
	GpuIdsProcessed  *int     `json:"gpuIdsProcessed,omitempty"`
}

// AssignPorts assigns or removes InfiniBand GPU ports for a tenant's UFM
// partition. Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/gpus.
func (r *GPUResource) AssignPorts(ctx context.Context, fabricName, tenantName string, operation ones_gfx.GpuAction, serverNames []string, gpuIds []int, membership *string) (*GpuPortAssignmentResult, error) {
	body := struct {
		Operation   ones_gfx.GpuAction `json:"operation"`
		ServerNames []string           `json:"serverNames,omitempty"`
		GpuIds      []int              `json:"gpuIds,omitempty"`
		Membership  *string            `json:"membership,omitempty"`
	}{operation, serverNames, gpuIds, membership}
	return ones_gfx.Call[GpuPortAssignmentResult](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/gpus", body, ones_gfx.OperationModeSynchronous, nil)
}

// List returns every GPU record under a fabric.
// Maps to GET /getAllGpusList/{fabricName}.
func (r *GPUResource) List(ctx context.Context, fabricName string) ([]GpuItem, error) {
	res, err := ones_gfx.Call[[]GpuItem](r.transport, "GET", "getAllGpusList/"+fabricName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ByHost returns the GPU devices of one host.
// Maps to GET /getGpusByHost/{hostName}/{fabricName}.
func (r *GPUResource) ByHost(ctx context.Context, hostName, fabricName string) ([]GpuItem, error) {
	res, err := ones_gfx.Call[[]GpuItem](r.transport, "GET", "getGpusByHost/"+hostName+"/"+fabricName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// TenantMappings returns per-GPU logical tenant mappings, optionally filtered
// by tenant. Maps to GET /fabrics/{fabricName}/gpuTenantMappings.
func (r *GPUResource) TenantMappings(ctx context.Context, fabricName string, tenantName *string) ([]GpuTenantMappingItem, error) {
	q := url.Values{}
	if tenantName != nil {
		q.Set("tenantName", *tenantName)
	}
	res, err := ones_gfx.Call[[]GpuTenantMappingItem](r.transport, "GET", "fabrics/"+fabricName+"/gpuTenantMappings", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// AllocationHistory returns the GPU allocation history of a tenant.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocationHistory.
func (r *GPUResource) AllocationHistory(ctx context.Context, fabricName, tenantName string) ([]GpuAllocationHistoryItem, error) {
	res, err := ones_gfx.Call[[]GpuAllocationHistoryItem](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/gpuAllocationHistory", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// AvailableServers lists servers with free GPUs in a fabric.
// Maps to GET /fabrics/{fabricName}/available_servers.
func (r *GPUResource) AvailableServers(ctx context.Context, fabricName string) (*AvailableServerResult, error) {
	return ones_gfx.Call[AvailableServerResult](r.transport, "GET", "fabrics/"+fabricName+"/available_servers", nil, ones_gfx.OperationModeSynchronous, nil)
}
