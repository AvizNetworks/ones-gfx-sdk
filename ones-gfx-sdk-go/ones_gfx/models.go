package ones_gfx

import (
	"encoding/json"
	"strings"
	"time"
)

// Operation represents an asynchronous job submitted to the API.
type Operation struct {
	ID                string          `json:"id"`
	Type              string          `json:"type,omitempty"`
	Status            OperationStatus `json:"status,omitempty"`
	Progress          int             `json:"progress,omitempty"`
	Result            interface{}     `json:"result,omitempty"` // Parsed or raw string
	ErrorMessage      string          `json:"errorMessage,omitempty"`
	HTTPStatusCode    *int            `json:"httpStatusCode,omitempty"`
	TenantName        string          `json:"tenantName,omitempty"`
	FabricName        string          `json:"fabricName,omitempty"`
	OperationType     string          `json:"operationType,omitempty"` // GPU_ALLOCATE, GPU_DEALLOCATE
	WebhookRegistered bool            `json:"webhookRegistered,omitempty"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	CompletedAt       time.Time       `json:"-"`
	CreatedAtRaw      interface{}     `json:"createdAt,omitempty"`
	UpdatedAtRaw      interface{}     `json:"updatedAt,omitempty"`
	CompletedAtRaw    interface{}     `json:"completedAt,omitempty"`
}

// UnmarshalJSON parses Operation and handles result field (may be JSON string).
func (o *Operation) UnmarshalJSON(data []byte) error {
	type Alias Operation
	aux := &struct {
		ResultRaw   interface{} `json:"result"`
		OperationID string      `json:"operationId"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if o.ID == "" && aux.OperationID != "" {
		o.ID = aux.OperationID
	}
	// Parse result (may be JSON-encoded string or plain string)
	if aux.ResultRaw != nil {
		if str, ok := aux.ResultRaw.(string); ok && str != "" {
			var parsed interface{}
			if err := json.Unmarshal([]byte(str), &parsed); err == nil {
				o.Result = parsed
			} else {
				o.Result = str
			}
		} else {
			o.Result = aux.ResultRaw
		}
	}
	// Parse timestamps
	o.CreatedAt = parseTimestamp(o.CreatedAtRaw)
	o.UpdatedAt = parseTimestamp(o.UpdatedAtRaw)
	o.CompletedAt = parseTimestamp(o.CompletedAtRaw)
	return nil
}

// IsDone returns true if the operation has reached a terminal state.
func (o *Operation) IsDone() bool {
	return o.Status.IsTerminal()
}

// IsSuccess returns true if the operation completed successfully.
func (o *Operation) IsSuccess() bool {
	return o.Status == OperationStatusSuccess
}

// IsFailure returns true if the operation failed.
func (o *Operation) IsFailure() bool {
	return o.Status.IsFailure()
}

// parseTimestamp handles both ms-since-epoch and ISO-8601 timestamps.
func parseTimestamp(raw interface{}) time.Time {
	if raw == nil {
		return time.Time{}
	}
	// Try as float64 (ms-since-epoch)
	if ms, ok := raw.(float64); ok {
		return time.Unix(0, int64(ms)*int64(time.Millisecond))
	}
	// Try as string (ISO-8601)
	if str, ok := raw.(string); ok && str != "" {
		// Normalize trailing Z to +00:00 for time.Parse
		normalized := strings.Replace(str, "Z", "+00:00", 1)
		if t, err := time.Parse(time.RFC3339Nano, normalized); err == nil {
			return t
		}
		// Fallback: try as epoch string
		if t, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", str); err == nil {
			return t
		}
	}
	return time.Time{}
}

// VNetTenantIPPool is the tenant's address pool inside a VNet. Both fields are
// null until IPAM has assigned the tenant a subnet.
type VNetTenantIPPool struct {
	Gateway string `json:"gateway,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
}

// VNetInfo describes a tenant's virtual network.
//
// Note the wire key for the pool is "tenant IP Pool", spaces included.
type VNetInfo struct {
	Name         string           `json:"name,omitempty"`
	Description  string           `json:"description,omitempty"`
	TenantIPPool VNetTenantIPPool `json:"tenant IP Pool"`
}

// NMXCPartition is one NMX-C partition attached to a tenant. Present only on
// fabrics with NMX-C integration.
type NMXCPartition struct {
	Name       string   `json:"name,omitempty"`
	DomainID   string   `json:"domainId,omitempty"`
	DomainName string   `json:"domainName,omitempty"`
	GUIDs      []string `json:"guids,omitempty"`
}

// Tenant is a tenant as Fabric Manager returns it, mirroring the Python SDK's
// ones_gfx.models.Tenant so the two stay readable side by side.
//
// FM builds this payload per fabric: the base fields depend on which starting
// subnets the fabric defines, UFM fabrics add the pkey/partition fields, and
// NMX-C fabrics merge in their own. Anything not named here is kept in Extra
// rather than dropped, so a field this SDK does not model yet is still
// reachable.
type Tenant struct {
	ID             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	FabricName     string `json:"fabricName,omitempty"`
	MaxGpusAllowed int    `json:"maxGpusAllowed,omitempty"`
	GpusAllocated  int    `json:"gpusAllocated,omitempty"`

	// AllotedGpus is a comma-separated list of *server* names, not GPU ids,
	// and FM spells it with one "t". Use AllottedServers for a slice.
	AllotedGpus string `json:"allotedGpus,omitempty"`

	VniID             *int   `json:"vniId,omitempty"`
	VlanID            *int   `json:"vlanId,omitempty"`
	TransitVlanID     *int   `json:"transitVlanId,omitempty"`
	VlanVniID         *int   `json:"vlanVniId,omitempty"`
	NsVniID           *int   `json:"nsVniId,omitempty"`
	VlanSubnetCPU     string `json:"vlanSubnetCPU,omitempty"`
	VlanSubnetStorage string `json:"vlanSubnetStorage,omitempty"`

	// TenantDhcpCidrPool is set only on fabrics assigning addresses by DHCP.
	TenantDhcpCidrPool string `json:"tenantDhcpCidrPool,omitempty"`

	// AnycastGateway is set only on HBN DPU fabrics in tenant DHCP mode.
	AnycastGateway string `json:"anycastGateway,omitempty"`

	// ConfigStatus arrives as the string "0".."3"; see ConfigStatus.
	ConfigStatus ConfigStatus `json:"config_status,omitempty"`

	Shared   bool                `json:"shared,omitempty"`
	Networks []string            `json:"networks,omitempty"`
	Tags     []map[string]string `json:"tags,omitempty"`
	VNets    VNetInfo            `json:"vnets,omitempty"`

	// UFM fabrics only.
	PkeyValue           string `json:"pkeyValue,omitempty"`
	UfmPartitionName    string `json:"ufmPartitionName,omitempty"`
	UfmAllocatedPorts   int    `json:"ufmAllocatedPorts,omitempty"`
	UfmAllocatedServers int    `json:"ufmAllocatedServers,omitempty"`

	// NMX-C fabrics only.
	NmxcGpusAllocated    int             `json:"nmxcGpusAllocated,omitempty"`
	NmxcServersAllocated int             `json:"nmxcServersAllocated,omitempty"`
	NmxcPartitions       []NMXCPartition `json:"nmxcPartitions,omitempty"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	CreatedAtRaw interface{} `json:"createdAt,omitempty"`
	UpdatedAtRaw interface{} `json:"updatedAt,omitempty"`

	// Extra holds every key not matched above, so a payload FM extends does
	// not silently lose data. Empty when the response is fully modelled.
	Extra map[string]interface{} `json:"-"`
}

// knownTenantKeys are the wire keys Tenant models directly; everything else
// lands in Extra.
var knownTenantKeys = map[string]struct{}{
	"id": {}, "name": {}, "description": {}, "fabricName": {},
	"maxGpusAllowed": {}, "gpusAllocated": {}, "allotedGpus": {},
	"vniId": {}, "vlanId": {}, "transitVlanId": {}, "vlanVniId": {}, "nsVniId": {},
	"vlanSubnetCPU": {}, "vlanSubnetStorage": {}, "tenantDhcpCidrPool": {},
	"anycastGateway": {}, "config_status": {}, "shared": {}, "networks": {},
	"tags": {}, "vnets": {}, "pkeyValue": {}, "ufmPartitionName": {},
	"ufmAllocatedPorts": {}, "ufmAllocatedServers": {}, "nmxcGpusAllocated": {},
	"nmxcServersAllocated": {}, "nmxcPartitions": {}, "createdAt": {}, "updatedAt": {},
}

// UnmarshalJSON decodes a tenant, normalising the timestamps and collecting
// any unmodelled keys into Extra.
func (t *Tenant) UnmarshalJSON(data []byte) error {
	type Alias Tenant
	aux := (*Alias)(t)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	t.CreatedAt = parseTimestamp(t.CreatedAtRaw)
	t.UpdatedAt = parseTimestamp(t.UpdatedAtRaw)

	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		return nil // object decoded above; nothing more to salvage
	}
	for k, v := range all {
		if _, known := knownTenantKeys[k]; known {
			continue
		}
		if t.Extra == nil {
			t.Extra = map[string]interface{}{}
		}
		var val interface{}
		if err := json.Unmarshal(v, &val); err == nil {
			t.Extra[k] = val
		}
	}
	return nil
}

// AllottedServers splits AllotedGpus into the server names it holds.
func (t *Tenant) AllottedServers() []string {
	if t == nil || strings.TrimSpace(t.AllotedGpus) == "" {
		return nil
	}
	var out []string
	for _, s := range strings.Split(t.AllotedGpus, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// IsUnlimitedGPUs reports whether the tenant was created with no GPU ceiling.
func (t *Tenant) IsUnlimitedGPUs() bool {
	return t != nil && t.MaxGpusAllowed == UnlimitedGPUs
}
