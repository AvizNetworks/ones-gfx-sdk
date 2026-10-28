package resources

import (
	"context"

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
	if err := firstErr(fabricTenant(fabricName, tenantName), gpuAction("operation", operation), nonEmptyNames("serverNames", serverNames), gpuIDs(gpuIds)); err != nil {
		return nil, err
	}
	body := struct {
		Operation   ones_gfx.GpuAction `json:"operation"`
		ServerNames []string           `json:"serverNames,omitempty"`
		GpuIds      []int              `json:"gpuIds,omitempty"`
		Membership  *string            `json:"membership,omitempty"`
	}{operation, serverNames, gpuIds, membership}
	return ones_gfx.Call[GpuPortAssignmentResult](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/gpus", body, ones_gfx.OperationModeSynchronous, nil)
}

// AvailableServers lists servers with free GPUs in a fabric.
// Maps to GET /fabrics/{fabricName}/available_servers.
func (r *GPUResource) AvailableServers(ctx context.Context, fabricName string) (*AvailableServerResult, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[AvailableServerResult](r.transport, "GET", "fabrics/"+fabricName+"/available_servers", nil, ones_gfx.OperationModeSynchronous, nil)
}
