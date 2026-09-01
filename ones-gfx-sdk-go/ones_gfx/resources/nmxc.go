package resources

import (
	"context"

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
	body := struct {
		Operation ones_gfx.GpuAction `json:"operation"`
		Domains   []NmxcDomain       `json:"domains"`
	}{operation, domains}
	return ones_gfx.Call[NmxcDomainsUpdateResult](r.transport, "PATCH", "fabrics/"+fabricName+"/nmxc/domains", body, ones_gfx.OperationModeSynchronous, nil)
}

// ProbeDomains performs a read-only Hello handshake against each domain.
// Maps to POST /fabrics/{fabricName}/nmxc/domains/probe.
func (r *NMXCResource) ProbeDomains(ctx context.Context, fabricName string, domains []NmxcDomain) ([]NmxcProbeResult, error) {
	body := struct {
		Domains []NmxcDomain `json:"domains"`
	}{domains}
	res, err := ones_gfx.Call[[]NmxcProbeResult](r.transport, "POST", "fabrics/"+fabricName+"/nmxc/domains/probe", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Inventory returns the NMX-C inventory snapshot for a fabric. The payload
// mixes fixed keys (success, error, message, fabricName) with dynamic domain/
// GPU data, so it is surfaced as a raw map.
// Maps to GET /fabrics/{fabricName}/nmxc/inventory.
func (r *NMXCResource) Inventory(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/nmxc/inventory", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Reset soft-resets a DEGRADED NMX-C domain.
// Maps to POST /api/nmxc/domains/{domainId}/reset.
func (r *NMXCResource) Reset(ctx context.Context, domainID string) (*NmxcResetResult, error) {
	return ones_gfx.Call[NmxcResetResult](r.transport, "POST", "api/nmxc/domains/"+domainID+"/reset", nil, ones_gfx.OperationModeSynchronous, nil)
}

// FactoryReset destructively factory-resets a domain (returns its backup).
// Maps to POST /api/nmxc/domains/{domainId}/factory-reset.
func (r *NMXCResource) FactoryReset(ctx context.Context, domainID string) (*NmxcFactoryResetResult, error) {
	return ones_gfx.Call[NmxcFactoryResetResult](r.transport, "POST", "api/nmxc/domains/"+domainID+"/factory-reset", nil, ones_gfx.OperationModeSynchronous, nil)
}
