package resources

import (
	"context"
	"encoding/json"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// NMXCResource covers NMX-C domain management, probing and recovery.
type NMXCResource struct {
	transport *ones_gfx.Transport
}

// NewNMXCResource constructs an NMXCResource.
func NewNMXCResource(transport *ones_gfx.Transport) *NMXCResource {
	return &NMXCResource{transport: transport}
}

// NmxcDomain is one NMX-C domain entry (host required; port/useTls optional).
type NmxcDomain struct {
	Host   string `json:"host"`
	Port   *int   `json:"port,omitempty"`
	UseTLS *bool  `json:"useTls,omitempty"`
}

// NmxcDomainsUpdateResult is the response of UpdateDomains
// (PATCH /fabrics/{fabricName}/nmxc/domains). Failure keys appear only on
// partial failure (HTTP 207).
type NmxcDomainsUpdateResult struct {
	Success         bool                     `json:"success,omitempty"`
	Operation       string                   `json:"operation,omitempty"`
	NotRegistered   []map[string]interface{} `json:"notRegistered,omitempty"`
	NotDeregistered []map[string]interface{} `json:"notDeregistered,omitempty"`
}

// NmxcProbeResult is one entry of ProbeDomains
// (POST /fabrics/{fabricName}/nmxc/domains/probe).
type NmxcProbeResult struct {
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	Reachable  bool    `json:"reachable"`
	Configured bool    `json:"configured"`
	Error      *string `json:"error"`
}

// NmxcResetResult is the 202 body of Reset
// (POST /api/nmxc/domains/{domainId}/reset).
type NmxcResetResult struct {
	DomainID string `json:"domainId"`
	Action   string `json:"action"`
	Message  string `json:"message"`
}

// NmxcFactoryResetResult is the body of FactoryReset
// (POST /api/nmxc/domains/{domainId}/factory-reset).
type NmxcFactoryResetResult struct {
	DomainID      string `json:"domainId,omitempty"`
	Action        string `json:"action,omitempty"`
	Message       string `json:"message,omitempty"`
	BackupApplied string `json:"backupApplied,omitempty"`
}

// UpdateDomains registers (ADD) or de-registers (DELETE) NMX-C domains.
// Maps to PATCH /fabrics/{fabricName}/nmxc/domains.
func (r *NMXCResource) UpdateDomains(ctx context.Context, fabricName string, domains []NmxcDomain, operation ones_gfx.GpuAction) (*NmxcDomainsUpdateResult, error) {
	if err := firstErr(pathName("fabricName", fabricName), gpuAction("operation", operation), nmxcDomains(domains)); err != nil {
		return nil, err
	}
	body := struct {
		Operation ones_gfx.GpuAction `json:"operation"`
		Domains   []NmxcDomain       `json:"domains"`
	}{operation, domains}
	return ones_gfx.Call[NmxcDomainsUpdateResult](r.transport, "PATCH", "fabrics/"+fabricName+"/nmxc/domains", body, ones_gfx.OperationModeSynchronous, nil)
}

// ProbeDomains performs a read-only Hello handshake against each domain.
// Maps to POST /fabrics/{fabricName}/nmxc/domains/probe.
func (r *NMXCResource) ProbeDomains(ctx context.Context, fabricName string, domains []NmxcDomain) ([]NmxcProbeResult, error) {
	if err := firstErr(pathName("fabricName", fabricName), nmxcDomains(domains)); err != nil {
		return nil, err
	}
	body := struct {
		Domains []NmxcDomain `json:"domains"`
	}{domains}
	res, err := ones_gfx.Call[[]NmxcProbeResult](r.transport, "POST", "fabrics/"+fabricName+"/nmxc/domains/probe", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// NmxcInventoryGPU is one GPU in an NMX-C tray. PartitionID and the tenant
// fields are set only while the GPU is allocated.
type NmxcInventoryGPU struct {
	GpuUID       string `json:"gpuUid,omitempty"`
	GpuID        *int   `json:"gpuId,omitempty"`
	HealthStatus string `json:"healthStatus,omitempty"`
	PartitionID  *int   `json:"partitionId,omitempty"`
	TenantID     *int   `json:"tenantId,omitempty"`
	TenantName   string `json:"tenantName,omitempty"`
	SystemUID    string `json:"systemUid,omitempty"`
}

// NmxcInventoryServer is one compute tray. ServerName is reconstructed by FM
// as <domainName>-node<NN>.
type NmxcInventoryServer struct {
	ServerName           string             `json:"serverName,omitempty"`
	TrayIndex            int                `json:"trayIndex"`
	TotalGpus            int                `json:"totalGpus"`
	FreeGpus             int                `json:"freeGpus"`
	AllocatedGpus        int                `json:"allocatedGpus"`
	TenantID             *int               `json:"tenantId,omitempty"`
	TenantName           string             `json:"tenantName,omitempty"`
	LastConfiguredTenant string             `json:"lastConfiguredTenant,omitempty"`
	Gpus                 []NmxcInventoryGPU `json:"gpus"`
}

// NmxcInventoryPartition is one NMX-C partition in a domain.
type NmxcInventoryPartition struct {
	PartitionID    *int   `json:"partitionId,omitempty"`
	PartitionName  string `json:"partitionName,omitempty"`
	TenantID       *int   `json:"tenantId,omitempty"`
	TenantName     string `json:"tenantName,omitempty"`
	ResiliencyMode string `json:"resiliencyMode,omitempty"`
	Status         string `json:"status,omitempty"`
	GpuCount       int    `json:"gpuCount"`
}

// NmxcInventoryDomain is one NMX-C domain with its partitions and trays.
// DomainID is what Reset and FactoryReset take: pass DomainID.String().
type NmxcInventoryDomain struct {
	DomainID          json.Number              `json:"domainId,omitempty"`
	DomainName        string                   `json:"domainName,omitempty"`
	ConnectionState   string                   `json:"connectionState,omitempty"`
	NmxcHost          string                   `json:"nmxcHost,omitempty"`
	NmxcPort          int                      `json:"nmxcPort,omitempty"`
	MaxGpus           *int                     `json:"maxGpus,omitempty"`
	MaxPartitions     *int                     `json:"maxPartitions,omitempty"`
	TopologyType      string                   `json:"topologyType,omitempty"`
	LastSyncedAt      ones_gfx.Timestamp       `json:"lastSyncedAt,omitempty"`
	GpuCount          int                      `json:"gpuCount"`
	FreeGpuCount      int                      `json:"freeGpuCount"`
	AllocatedGpuCount int                      `json:"allocatedGpuCount"`
	ServerCount       int                      `json:"serverCount"`
	Partitions        []NmxcInventoryPartition `json:"partitions"`
	Servers           []NmxcInventoryServer    `json:"servers"`
}

// NmxcInventorySummary totals the inventory across every domain.
type NmxcInventorySummary struct {
	DomainCount       int `json:"domainCount"`
	TotalGpuCount     int `json:"totalGpuCount"`
	FreeGpuCount      int `json:"freeGpuCount"`
	AllocatedGpuCount int `json:"allocatedGpuCount"`
	ServerCount       int `json:"serverCount"`
}

// NmxcInventory is the response of NMXC.Inventory.
//
// OverallHealth is HEALTHY when every domain is ACTIVE, DEGRADED when any is
// reconnecting, and UNHEALTHY when any is degraded.
type NmxcInventory struct {
	FabricName    string                `json:"fabricName,omitempty"`
	FabricType    string                `json:"fabricType,omitempty"`
	Success       bool                  `json:"success"`
	OverallHealth string                `json:"overallHealth,omitempty"`
	Summary       NmxcInventorySummary  `json:"summary"`
	Domains       []NmxcInventoryDomain `json:"domains"`

	// Extra holds any top-level key this SDK does not model yet.
	Extra map[string]json.RawMessage `json:"-"`
}

var nmxcInventoryKeys = keySet("fabricName", "fabricType", "success", "overallHealth", "summary", "domains")

// UnmarshalJSON decodes the inventory and collects unmodelled keys into Extra.
func (n *NmxcInventory) UnmarshalJSON(data []byte) error {
	type alias NmxcInventory
	if err := json.Unmarshal(data, (*alias)(n)); err != nil {
		return err
	}
	n.Extra = extraKeys(data, nmxcInventoryKeys)
	return nil
}

// Inventory returns the NMX-C inventory snapshot for a fabric. FM serves it
// from its own database; no NMX-C call is made.
// Maps to GET /fabrics/{fabricName}/nmxc/inventory.
func (r *NMXCResource) Inventory(ctx context.Context, fabricName string) (*NmxcInventory, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[NmxcInventory](r.transport, "GET", "fabrics/"+fabricName+"/nmxc/inventory", nil, ones_gfx.OperationModeSynchronous, nil)
}

// InventoryRaw returns the inventory response exactly as FM sent it.
// Prefer Inventory; this is the pre-typed behaviour, kept for migration.
func (r *NMXCResource) InventoryRaw(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/nmxc/inventory", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Reset soft-resets a DEGRADED NMX-C domain.
// Maps to POST /api/nmxc/domains/{domainId}/reset.
func (r *NMXCResource) Reset(ctx context.Context, domainID string) (*NmxcResetResult, error) {
	if err := pathName("domainId", domainID); err != nil {
		return nil, err
	}
	return ones_gfx.Call[NmxcResetResult](r.transport, "POST", "api/nmxc/domains/"+domainID+"/reset", nil, ones_gfx.OperationModeSynchronous, nil)
}

// FactoryReset destructively factory-resets a domain (returns its backup).
// Maps to POST /api/nmxc/domains/{domainId}/factory-reset.
func (r *NMXCResource) FactoryReset(ctx context.Context, domainID string) (*NmxcFactoryResetResult, error) {
	if err := pathName("domainId", domainID); err != nil {
		return nil, err
	}
	return ones_gfx.Call[NmxcFactoryResetResult](r.transport, "POST", "api/nmxc/domains/"+domainID+"/factory-reset", nil, ones_gfx.OperationModeSynchronous, nil)
}
