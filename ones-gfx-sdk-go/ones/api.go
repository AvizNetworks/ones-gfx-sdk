package ones

// api.go exposes every Fabric Manager endpoint as a package-level function,
// mirroring the flat apis surface of the Python SDK. Each function takes the
// *Client explicitly — there is no shared global. Calls fail with
// ErrNotAuthenticated until Login succeeds.
//
// Generated from the resource method set — edit the generator, not this file.

import "context"

// AssignGpuPorts calls GPU.AssignPorts.
func AssignGpuPorts(ctx context.Context, client *Client, fabricName, tenantName string, operation GpuAction, serverNames []string, gpuIds []int, membership *string) (*GpuPortAssignmentResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AssignPorts(ctx, fabricName, tenantName, operation, serverNames, gpuIds, membership)
}

// CreateTenant calls Tenants.Create.
func CreateTenant(ctx context.Context, client *Client, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Tenants.Create(ctx, fabricName, tenantName, description, maxGpusAllowed, shared, opts...)
}

// CreateTenantAsync calls Tenants.CreateAsync.
func CreateTenantAsync(ctx context.Context, client *Client, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.CreateAsync(ctx, fabricName, tenantName, description, maxGpusAllowed, shared, opts...)
}

// CreateVpcPeering calls VPCPeering.Create.
func CreateVpcPeering(ctx context.Context, client *Client, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.VPCPeering.Create(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// CreateVpcPeeringAsync calls VPCPeering.CreateAsync.
func CreateVpcPeeringAsync(ctx context.Context, client *Client, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.VPCPeering.CreateAsync(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// DeleteTenant calls Tenants.Delete.
func DeleteTenant(ctx context.Context, client *Client, fabricName, tenantName string, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.Delete(ctx, fabricName, tenantName, opts...)
}

// DeleteTenantAsync calls Tenants.DeleteAsync.
func DeleteTenantAsync(ctx context.Context, client *Client, fabricName, tenantName string, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.DeleteAsync(ctx, fabricName, tenantName, opts...)
}

// FactoryResetNmxcDomain calls NMXC.FactoryReset.
func FactoryResetNmxcDomain(ctx context.Context, client *Client, domainID string) (*NmxcFactoryResetResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.FactoryReset(ctx, domainID)
}

// GetAvailableServers calls GPU.AvailableServers.
func GetAvailableServers(ctx context.Context, client *Client, fabricName string) (*AvailableServerResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AvailableServers(ctx, fabricName)
}

// GetFabrics calls Fabrics.ListDTOs.
func GetFabrics(ctx context.Context, client *Client) (*FabricsListResponse, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.ListDTOs(ctx)
}

// GetInventoryHosts calls Inventory.UFMHosts.
func GetInventoryHosts(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.UFMHosts(ctx, fabricName)
}

// GetInventoryPorts calls Inventory.Ports.
func GetInventoryPorts(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.Ports(ctx, fabricName)
}

// GetNmxcInventory calls NMXC.Inventory.
func GetNmxcInventory(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.Inventory(ctx, fabricName)
}

// GetOperation calls Operations.Get.
func GetOperation(ctx context.Context, client *Client, operationID string) (*OperationStatusItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Operations.Get(ctx, operationID)
}

// GetTenant calls Tenants.Get.
func GetTenant(ctx context.Context, client *Client, fabricName, tenantName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.Get(ctx, fabricName, tenantName)
}

// InventorySync calls Inventory.Sync.
func InventorySync(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.Sync(ctx, fabricName)
}

// ListTenants calls Tenants.List.
func ListTenants(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.List(ctx, fabricName)
}

// ModifyGpuAllocations calls Tenants.ModifyAllocations.
func ModifyGpuAllocations(ctx context.Context, client *Client, fabricName, tenantName string, suid SuidMap, operation *GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.ModifyAllocations(ctx, fabricName, tenantName, suid, operation, configScope, unreachableDevices, opts...)
}

// ModifyGpuAllocationsAsync calls Tenants.ModifyAllocationsAsync.
func ModifyGpuAllocationsAsync(ctx context.Context, client *Client, fabricName, tenantName string, suid SuidMap, operation *GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.ModifyAllocationsAsync(ctx, fabricName, tenantName, suid, operation, configScope, unreachableDevices, opts...)
}

// ProbeNmxcDomains calls NMXC.ProbeDomains.
func ProbeNmxcDomains(ctx context.Context, client *Client, fabricName string, domains []NmxcDomain) ([]NmxcProbeResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.ProbeDomains(ctx, fabricName, domains)
}

// ResetNmxcDomain calls NMXC.Reset.
func ResetNmxcDomain(ctx context.Context, client *Client, domainID string) (*NmxcResetResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.Reset(ctx, domainID)
}

// UpdateNmxcDomains calls NMXC.UpdateDomains.
func UpdateNmxcDomains(ctx context.Context, client *Client, fabricName string, domains []NmxcDomain, operation GpuAction) (*NmxcDomainsUpdateResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.UpdateDomains(ctx, fabricName, domains, operation)
}

// UpdateTenant calls Tenants.Update.
func UpdateTenant(ctx context.Context, client *Client, fabricName, tenantName string, servers []GpuServerInfo, operation *GpuAction, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.Update(ctx, fabricName, tenantName, servers, operation, opts...)
}

// UpdateTenantAsync calls Tenants.UpdateAsync.
func UpdateTenantAsync(ctx context.Context, client *Client, fabricName, tenantName string, servers []GpuServerInfo, operation *GpuAction, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.UpdateAsync(ctx, fabricName, tenantName, servers, operation, opts...)
}

// ValidateUfmCreds calls Devices.ValidateUFMCreds.
func ValidateUfmCreds(ctx context.Context, client *Client, ufmURL, username, password string) (*UfmCredsResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Devices.ValidateUFMCreds(ctx, ufmURL, username, password)
}
