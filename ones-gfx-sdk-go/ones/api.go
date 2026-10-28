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

// DeleteVpcPeering calls VPCPeering.Delete.
func DeleteVpcPeering(ctx context.Context, client *Client, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.VPCPeering.Delete(ctx, fabricName, name, vpcname, peervpcname, opts...)
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
func GetInventoryHosts(ctx context.Context, client *Client, fabricName string) (*UFMHostsResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.UFMHosts(ctx, fabricName)
}

// GetInventoryPorts calls Inventory.Ports.
func GetInventoryPorts(ctx context.Context, client *Client, fabricName string) (*InventoryPortsResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.Ports(ctx, fabricName)
}

// GetNmxcInventory calls NMXC.Inventory.
func GetNmxcInventory(ctx context.Context, client *Client, fabricName string) (*NmxcInventory, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.Inventory(ctx, fabricName)
}

// GetInventoryHostsRaw calls Inventory.UFMHostsRaw.
func GetInventoryHostsRaw(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.UFMHostsRaw(ctx, fabricName)
}

// GetInventoryPortsRaw calls Inventory.PortsRaw.
func GetInventoryPortsRaw(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.PortsRaw(ctx, fabricName)
}

// GetNmxcInventoryRaw calls NMXC.InventoryRaw.
func GetNmxcInventoryRaw(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.InventoryRaw(ctx, fabricName)
}

// GetOperation calls Operations.Get.
func GetOperation(ctx context.Context, client *Client, operationID string) (*OperationStatusItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Operations.Get(ctx, operationID)
}

// GetTenant calls Tenants.Get.
func GetTenant(ctx context.Context, client *Client, fabricName, tenantName string) (*Tenant, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.Get(ctx, fabricName, tenantName)
}

// GetTenantRaw calls Tenants.GetRaw.
func GetTenantRaw(ctx context.Context, client *Client, fabricName, tenantName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.GetRaw(ctx, fabricName, tenantName)
}

// InventorySync calls Inventory.Sync.
func InventorySync(ctx context.Context, client *Client, fabricName string) (*InventorySyncResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.Sync(ctx, fabricName)
}

// InventorySyncRaw calls Inventory.SyncRaw.
func InventorySyncRaw(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.SyncRaw(ctx, fabricName)
}

// ListTenants calls Tenants.List.
func ListTenants(ctx context.Context, client *Client, fabricName string) ([]Tenant, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.List(ctx, fabricName)
}

// ListTenantsRaw calls Tenants.ListRaw.
func ListTenantsRaw(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.ListRaw(ctx, fabricName)
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

// --- network policy ---------------------------------------------------------

// ListNetworkPolicies calls NetworkPolicy.List.
func ListNetworkPolicies(ctx context.Context, client *Client, fabricName string, opts ...CallOption) (*NetworkPolicyList, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NetworkPolicy.List(ctx, fabricName, opts...)
}

// CreateNetworkPolicy calls NetworkPolicy.Create.
func CreateNetworkPolicy(ctx context.Context, client *Client, fabricName, policyName string, source, destination NetworkPolicySide, opts ...CallOption) (*NetworkPolicyRef, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NetworkPolicy.Create(ctx, fabricName, policyName, source, destination, opts...)
}

// PatchNetworkPolicy calls NetworkPolicy.Patch.
func PatchNetworkPolicy(ctx context.Context, client *Client, fabricName, policyName string, source, destination NetworkPolicySide, opts ...CallOption) (*NetworkPolicyRef, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NetworkPolicy.Patch(ctx, fabricName, policyName, source, destination, opts...)
}

// DeleteNetworkPolicy calls NetworkPolicy.Delete.
func DeleteNetworkPolicy(ctx context.Context, client *Client, fabricName, policyName string, opts ...CallOption) (*NetworkPolicyRef, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NetworkPolicy.Delete(ctx, fabricName, policyName, opts...)
}

// --- webhook registration ---------------------------------------------------

// RegisterWebhook calls Webhooks.Register.
func RegisterWebhook(ctx context.Context, client *Client, callbackURL string, eventTypes []string, enabled *bool, filters *WebhookFilters, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.Register(ctx, callbackURL, eventTypes, enabled, filters, opts...)
}

// ListWebhooks calls Webhooks.List.
func ListWebhooks(ctx context.Context, client *Client, opts ...CallOption) ([]WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.List(ctx, opts...)
}

// GetWebhook calls Webhooks.Get.
func GetWebhook(ctx context.Context, client *Client, webhookID string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.Get(ctx, webhookID, opts...)
}

// UpdateWebhook calls Webhooks.Update.
func UpdateWebhook(ctx context.Context, client *Client, webhookID string, callbackURL *string, enabled *bool, eventTypes []string, filters *WebhookFilters, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.Update(ctx, webhookID, callbackURL, enabled, eventTypes, filters, opts...)
}

// DeleteWebhook calls Webhooks.Delete.
func DeleteWebhook(ctx context.Context, client *Client, webhookID string, opts ...CallOption) error {
	if client == nil {
		return ErrNilClient
	}
	return client.Webhooks.Delete(ctx, webhookID, opts...)
}

// RegisterFabricWebhook calls Webhooks.RegisterFabric.
func RegisterFabricWebhook(ctx context.Context, client *Client, fabricName, callbackURL string, enabled *bool, eventTypes []string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.RegisterFabric(ctx, fabricName, callbackURL, enabled, eventTypes, opts...)
}

// GetFabricWebhook calls Webhooks.GetFabric.
func GetFabricWebhook(ctx context.Context, client *Client, fabricName string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.GetFabric(ctx, fabricName, opts...)
}

// UpdateFabricWebhook calls Webhooks.UpdateFabric.
func UpdateFabricWebhook(ctx context.Context, client *Client, fabricName string, callbackURL *string, enabled bool, eventTypes []string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.UpdateFabric(ctx, fabricName, callbackURL, enabled, eventTypes, opts...)
}

// DeleteFabricWebhook calls Webhooks.DeleteFabric.
func DeleteFabricWebhook(ctx context.Context, client *Client, fabricName string, opts ...CallOption) error {
	if client == nil {
		return ErrNilClient
	}
	return client.Webhooks.DeleteFabric(ctx, fabricName, opts...)
}

// RegisterTenantWebhook calls Webhooks.RegisterTenant.
func RegisterTenantWebhook(ctx context.Context, client *Client, fabricName, tenantName, callbackURL string, enabled *bool, eventTypes []string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.RegisterTenant(ctx, fabricName, tenantName, callbackURL, enabled, eventTypes, opts...)
}

// GetTenantWebhook calls Webhooks.GetTenant.
func GetTenantWebhook(ctx context.Context, client *Client, fabricName, tenantName string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.GetTenant(ctx, fabricName, tenantName, opts...)
}

// UpdateTenantWebhook calls Webhooks.UpdateTenant.
func UpdateTenantWebhook(ctx context.Context, client *Client, fabricName, tenantName string, callbackURL *string, enabled bool, eventTypes []string, opts ...CallOption) (*WebhookRegistration, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Webhooks.UpdateTenant(ctx, fabricName, tenantName, callbackURL, enabled, eventTypes, opts...)
}

// DeleteTenantWebhook calls Webhooks.DeleteTenant.
func DeleteTenantWebhook(ctx context.Context, client *Client, fabricName, tenantName string, opts ...CallOption) error {
	if client == nil {
		return ErrNilClient
	}
	return client.Webhooks.DeleteTenant(ctx, fabricName, tenantName, opts...)
}

// --- IPAM configuration and address queries ---------------------------------

// GetIpamConfig calls IPAM.GetConfig.
func GetIpamConfig(ctx context.Context, client *Client, fabricName string, opts ...CallOption) (*IpamConfig, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.IPAM.GetConfig(ctx, fabricName, opts...)
}

// UpdateIpamConfig calls IPAM.UpdateConfig.
func UpdateIpamConfig(ctx context.Context, client *Client, fabricName string, cfg IpamConfig, force bool, opts ...CallOption) (*IpamConfig, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.IPAM.UpdateConfig(ctx, fabricName, cfg, force, opts...)
}

// GetTenantIps calls IPAM.TenantIPs.
func GetTenantIps(ctx context.Context, client *Client, fabricName, tenantName string, page, size int, status string, opts ...CallOption) (*TenantIPList, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.IPAM.TenantIPs(ctx, fabricName, tenantName, page, size, status, opts...)
}

// GetServerIps calls IPAM.ServerIPs.
func GetServerIps(ctx context.Context, client *Client, fabricName, serverName string, opts ...CallOption) (*ServerIPList, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.IPAM.ServerIPs(ctx, fabricName, serverName, opts...)
}

// GetTenantComputeIp calls IPAM.TenantComputeIP.
func GetTenantComputeIp(ctx context.Context, client *Client, fabricName, tenantName, serverName string, opts ...CallOption) (IpamAllocation, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.IPAM.TenantComputeIP(ctx, fabricName, tenantName, serverName, opts...)
}

// --- tenant DHCP ------------------------------------------------------------

// GetTenantDhcp calls DHCP.GetTenant.
func GetTenantDhcp(ctx context.Context, client *Client, fabricName, tenantName string, page, size int, opts ...CallOption) (*TenantDhcp, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.DHCP.GetTenant(ctx, fabricName, tenantName, page, size, opts...)
}

// RetryTenantDhcp calls DHCP.RetryTenant.
func RetryTenantDhcp(ctx context.Context, client *Client, fabricName, tenantName string, opts ...CallOption) (DhcpReapplySummary, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.DHCP.RetryTenant(ctx, fabricName, tenantName, opts...)
}

// ReconcileFabricDhcp calls DHCP.Reconcile.
func ReconcileFabricDhcp(ctx context.Context, client *Client, fabricName string, opts ...CallOption) (DhcpReapplySummary, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.DHCP.Reconcile(ctx, fabricName, opts...)
}

// --- rail CIDRPools ---------------------------------------------------------

// GetRailCidrPools calls Rails.CidrPools.
func GetRailCidrPools(ctx context.Context, client *Client, fabricName string, opts ...CallOption) (RailCidrPools, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Rails.CidrPools(ctx, fabricName, opts...)
}

// GetRailCidrPoolsYAML calls Rails.CidrPoolsYAML.
func GetRailCidrPoolsYAML(ctx context.Context, client *Client, fabricName, mediaType string, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Rails.CidrPoolsYAML(ctx, fabricName, mediaType, opts...)
}

// GetRailCidrPool calls Rails.CidrPool.
func GetRailCidrPool(ctx context.Context, client *Client, fabricName, railName string, opts ...CallOption) (RailCidrPools, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Rails.CidrPool(ctx, fabricName, railName, opts...)
}

// GetRailCidrPoolYAML calls Rails.CidrPoolYAML.
func GetRailCidrPoolYAML(ctx context.Context, client *Client, fabricName, railName, mediaType string, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Rails.CidrPoolYAML(ctx, fabricName, railName, mediaType, opts...)
}

// --- request-struct variants -------------------------------------------------

// CreateTenantWith calls Tenants.CreateWith.
func CreateTenantWith(ctx context.Context, client *Client, req CreateTenantRequest, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Tenants.CreateWith(ctx, req, opts...)
}

// CreateTenantAsyncWith calls Tenants.CreateAsyncWith.
func CreateTenantAsyncWith(ctx context.Context, client *Client, req CreateTenantRequest, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.CreateAsyncWith(ctx, req, opts...)
}

// UpdateTenantWith calls Tenants.UpdateWith.
func UpdateTenantWith(ctx context.Context, client *Client, req UpdateTenantServersRequest, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.UpdateWith(ctx, req, opts...)
}

// UpdateTenantAsyncWith calls Tenants.UpdateAsyncWith.
func UpdateTenantAsyncWith(ctx context.Context, client *Client, req UpdateTenantServersRequest, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.UpdateAsyncWith(ctx, req, opts...)
}

// ModifyGpuAllocationsWith calls Tenants.ModifyAllocationsWith.
func ModifyGpuAllocationsWith(ctx context.Context, client *Client, req ModifyAllocationsRequest, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.ModifyAllocationsWith(ctx, req, opts...)
}

// ModifyGpuAllocationsAsyncWith calls Tenants.ModifyAllocationsAsyncWith.
func ModifyGpuAllocationsAsyncWith(ctx context.Context, client *Client, req ModifyAllocationsRequest, opts ...CallOption) (*OperationAccepted, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.ModifyAllocationsAsyncWith(ctx, req, opts...)
}

// AssignGpuPortsWith calls GPU.AssignPortsWith.
func AssignGpuPortsWith(ctx context.Context, client *Client, req AssignPortsRequest) (*GpuPortAssignmentResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AssignPortsWith(ctx, req)
}

// --- HBN VF interfaces -------------------------------------------------------

// ListVfInterfaces calls VFs.List.
func ListVfInterfaces(ctx context.Context, client *Client, fabricName, serverName string, opts ...CallOption) (*VFInterfaceList, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.VFs.List(ctx, fabricName, serverName, opts...)
}

// AssignVfInterface calls VFs.Assign.
func AssignVfInterface(ctx context.Context, client *Client, fabricName, serverName, vfID, tenantName string, opts ...CallOption) (*VFAssignment, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.VFs.Assign(ctx, fabricName, serverName, vfID, tenantName, opts...)
}

// UnassignVfInterface calls VFs.Unassign.
func UnassignVfInterface(ctx context.Context, client *Client, fabricName, serverName, vfID string, opts ...CallOption) (*ApiResponseMessage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.VFs.Unassign(ctx, fabricName, serverName, vfID, opts...)
}
