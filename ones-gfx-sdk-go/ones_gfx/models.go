package ones_gfx

import (
	"encoding/json"
	"strings"
	"time"
)

// Fabric represents a top-level networking domain containing tenants.
type Fabric struct {
	ID                 int    `json:"id"`
	FabricName         string `json:"fabricName"`
	Description        string `json:"description"`
	NumOfSUs           int    `json:"numOfSUs"`
	MaxNumOfSUs        int    `json:"maxNumOfSUs"`
	EWTenantAware      bool   `json:"e-wTenantAware"`
	NSTenantAware      bool   `json:"n-sTenantAware"`
	DefaultStorageName string `json:"defaultStorageName,omitempty"`
}

// VNetInfo describes a virtual network attached to a tenant.
type VNetInfo struct {
	Name          string                 `json:"name,omitempty"`
	Description   string                 `json:"description,omitempty"`
	IPPoolGateway string                 `json:"-"` // Extracted from nested "tenant IP Pool"
	IPPoolSubnet  string                 `json:"-"`
	rawIPPool     map[string]interface{} // Temp storage during unmarshal
}

// UnmarshalJSON custom unmarshaler to handle the nested "tenant IP Pool" field.
func (v *VNetInfo) UnmarshalJSON(data []byte) error {
	type Alias VNetInfo
	aux := &struct {
		IPPool map[string]interface{} `json:"tenant IP Pool"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.IPPool != nil {
		if gw, ok := aux.IPPool["gateway"].(string); ok {
			v.IPPoolGateway = gw
		}
		if subnet, ok := aux.IPPool["subnet"].(string); ok {
			v.IPPoolSubnet = subnet
		}
	}
	return nil
}

// Tenant represents a tenant within a fabric.
type Tenant struct {
	ID                    int             `json:"id"`
	Name                  string          `json:"name"`
	Description           *string         `json:"description"`
	FabricName            string          `json:"fabricName"`
	MaxGPUsAllowed        int             `json:"maxGpusAllowed"`
	GPUsAllocated         int             `json:"gpusAllocated"`
	AllotedGPUs           string          `json:"allotedGpus"` // Comma-separated hostnames
	VNIID                 *int            `json:"vniId,omitempty"`
	VLANID                *int            `json:"vlanId,omitempty"`
	TransitVLANID         *int            `json:"transitVlanId,omitempty"`
	VLANVNIID             *int            `json:"vlanVniId,omitempty"`
	VLANSubnetCPU         *string         `json:"vlanSubnetCPU,omitempty"`
	VLANSubnetStorage     *string         `json:"vlanSubnetStorage,omitempty"`
	ConfigStatus          ConfigStatus    `json:"-"` // Parsed from string
	ConfigStatusRaw       string          `json:"config_status"`
	Networks              []string        `json:"networks,omitempty"`
	Tags                  []Tag           `json:"tags,omitempty"`
	VNets                 *VNetInfo       `json:"vnets,omitempty"`
	Shared                bool            `json:"shared,omitempty"`
	NSVNIID               *int            `json:"nsVniId,omitempty"`
	CreatedAt             time.Time       `json:"-"` // Parsed from ms-epoch or ISO
	UpdatedAt             time.Time       `json:"-"`
	CreatedAtRaw          interface{}     `json:"createdAt"`
	UpdatedAtRaw          interface{}     `json:"updatedAt"`

	// UFM-specific fields (IB_EW fabrics).
	PKeyValue             *string         `json:"pkeyValue"`
	UFMAllocatedPorts     int             `json:"ufmAllocatedPorts,omitempty"`
	UFMAllocatedServers   int             `json:"ufmAllocatedServers,omitempty"`

	// NMXC-specific fields.
	NMXCGPUsAllocated     int             `json:"nmxcGpusAllocated,omitempty"`
	NMXCServersAllocated  int             `json:"nmxcServersAllocated,omitempty"`
	NMXCPartitions        []NMXCPartition `json:"nmxcPartitions,omitempty"`
}

// NMXCPartition represents a partition assigned to a tenant on an NMXC fabric.
type NMXCPartition struct {
	DomainName     string `json:"domainName"`
	PartitionID    int    `json:"partitionId"`
	PartitionName  string `json:"partitionName"`
	ResiliencyMode string `json:"resiliencyMode"`
	Status         string `json:"status"`
}

// Tag represents a key-value tag on a tenant.
type Tag map[string]string

// UnmarshalJSON parses Tenant and handles config_status + timestamps.
func (t *Tenant) UnmarshalJSON(data []byte) error {
	type Alias Tenant
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Parse config_status
	t.ConfigStatus = ParseConfigStatus(t.ConfigStatusRaw)
	// Parse timestamps
	t.CreatedAt = parseTimestamp(t.CreatedAtRaw)
	t.UpdatedAt = parseTimestamp(t.UpdatedAtRaw)
	return nil
}

// AllotedServers returns the allocated GPU server hostnames as a slice.
func (t *Tenant) AllotedServers() []string {
	if t.AllotedGPUs == "" {
		return nil
	}
	parts := strings.Split(t.AllotedGPUs, ",")
	servers := make([]string, 0, len(parts))
	for _, s := range parts {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			servers = append(servers, trimmed)
		}
	}
	return servers
}

// IsUnlimitedGPUs returns true if this tenant has no GPU quota cap.
func (t *Tenant) IsUnlimitedGPUs() bool {
	return t.MaxGPUsAllowed == UnlimitedGPUs
}

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
	return o.Status == OperationStatusFailure
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
