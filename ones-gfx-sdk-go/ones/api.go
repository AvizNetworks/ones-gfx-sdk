package ones

// api.go exposes every Fabric Manager endpoint as a package-level function,
// mirroring the flat apis surface of the Python SDK. Each function takes the
// *Client explicitly — there is no shared global. Calls fail with
// ErrNotAuthenticated until Login succeeds.
//
// Generated from the resource method set — edit the generator, not this file.

import "context"

// AddDeviceFacts calls Devices.AddFacts.
func AddDeviceFacts(ctx context.Context, client *Client, items []DeviceDetail) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.AddFacts(ctx, items)
}

// AddHostTenantData calls HostTenants.Add.
func AddHostTenantData(ctx context.Context, client *Client, name, fabricName string, description *string, hostsAllocated *int, vniID *int, configStatus *string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.HostTenants.Add(ctx, name, fabricName, description, hostsAllocated, vniID, configStatus)
}

// AddInventoryData calls Inventory.Add.
func AddInventoryData(ctx context.Context, client *Client, items []InventoryItem) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Inventory.Add(ctx, items)
}

// AssignGpuPorts calls GPU.AssignPorts.
func AssignGpuPorts(ctx context.Context, client *Client, fabricName, tenantName string, operation GpuAction, serverNames []string, gpuIds []int, membership *string) (*GpuPortAssignmentResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AssignPorts(ctx, fabricName, tenantName, operation, serverNames, gpuIds, membership)
}

// AutoAllocateGpusToTenants calls Tenants.AutoAllocate.
func AutoAllocateGpusToTenants(ctx context.Context, client *Client, fabricName, tenantName string, autoAllocationDevicesNeed int, suid SuidMap) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Tenants.AutoAllocate(ctx, fabricName, tenantName, autoAllocationDevicesNeed, suid)
}

// BackupConfig calls ConfigMgmt.Backup.
func BackupConfig(ctx context.Context, client *Client, items []DeviceConfigBackup) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.ConfigMgmt.Backup(ctx, items)
}

// ConfigsListToRestore calls ConfigMgmt.ListToRestore.
func ConfigsListToRestore(ctx context.Context, client *Client, devices []string, onlylimited *bool) (*string, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.ConfigMgmt.ListToRestore(ctx, devices, onlylimited)
}

// CreateFabric calls Fabrics.Create.
func CreateFabric(ctx context.Context, client *Client, name string, args *FabricCreateArgs) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabrics.Create(ctx, name, args)
}

// CreateFabricSim calls Fabricsims.Create.
func CreateFabricSim(ctx context.Context, client *Client, id *int, fabricName, simulationID, username, token, orgUUID, status, uiLink *string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabricsims.Create(ctx, id, fabricName, simulationID, username, token, orgUUID, status, uiLink)
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

// DeleteFabric calls Fabrics.Delete.
func DeleteFabric(ctx context.Context, client *Client, name string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabrics.Delete(ctx, name)
}

// DeleteFabricSim calls Fabricsims.Delete.
func DeleteFabricSim(ctx context.Context, client *Client, name string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabricsims.Delete(ctx, name)
}

// DeleteFile calls Files.Delete.
func DeleteFile(ctx context.Context, client *Client, id int) (*DeleteFileResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Files.Delete(ctx, id)
}

// DeleteHostTenantData calls HostTenants.Delete.
func DeleteHostTenantData(ctx context.Context, client *Client, fabricName, tenantName string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.HostTenants.Delete(ctx, fabricName, tenantName)
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
func DeleteVpcPeering(ctx context.Context, client *Client, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.VPCPeering.Delete(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// EditFabric calls Fabrics.Update.
func EditFabric(ctx context.Context, client *Client, fabricName string, args *FabricCreateArgs) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabrics.Update(ctx, fabricName, args)
}

// EditInventoryData calls Inventory.Edit.
func EditInventoryData(ctx context.Context, client *Client, fabricName string, items []InventoryItem) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Inventory.Edit(ctx, fabricName, items)
}

// EnableZTPUpgrade calls Devices.EnableZTP.
func EnableZTPUpgrade(ctx context.Context, client *Client, devices []string) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.EnableZTP(ctx, devices)
}

// FactoryResetNmxcDomain calls NMXC.FactoryReset.
func FactoryResetNmxcDomain(ctx context.Context, client *Client, domainID string) (*NmxcFactoryResetResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.FactoryReset(ctx, domainID)
}

// FetchDeviceBackupFiles calls ConfigMgmt.FetchBackupFiles.
func FetchDeviceBackupFiles(ctx context.Context, client *Client, ip, timestamp *string) (*FetchBackupFilesResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.ConfigMgmt.FetchBackupFiles(ctx, ip, timestamp)
}

// FillBootstrapConfig calls Bootstrap.Fill.
func FillBootstrapConfig(ctx context.Context, client *Client, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Bootstrap.Fill(ctx, batchName, subnet, netmask, gateway, bootstrapinfo)
}

// FillRmaConfig calls RMA.Fill.
func FillRmaConfig(ctx context.Context, client *Client, items []RMAInfoItem) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.RMA.Fill(ctx, items)
}

// GetAllBootstrapBatches calls Bootstrap.ListBatches.
func GetAllBootstrapBatches(ctx context.Context, client *Client) ([]map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Bootstrap.ListBatches(ctx)
}

// GetAllFabricSims calls Fabricsims.List.
func GetAllFabricSims(ctx context.Context, client *Client) ([]FabricSimItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabricsims.List(ctx)
}

// GetAllFabrics calls Fabrics.List.
func GetAllFabrics(ctx context.Context, client *Client) ([]FabricItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.List(ctx)
}

// GetAllGpusList calls GPU.List.
func GetAllGpusList(ctx context.Context, client *Client, fabricName string) ([]GpuItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.List(ctx, fabricName)
}

// GetAllInventory calls Inventory.ListAll.
func GetAllInventory(ctx context.Context, client *Client) ([]InventoryRecord, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.ListAll(ctx)
}

// GetAvailableServers calls GPU.AvailableServers.
func GetAvailableServers(ctx context.Context, client *Client, fabricName string) (*AvailableServerResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AvailableServers(ctx, fabricName)
}

// GetBootstrapBatch calls Bootstrap.GetBatch.
func GetBootstrapBatch(ctx context.Context, client *Client, batchName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Bootstrap.GetBatch(ctx, batchName)
}

// GetBootstrapInfo calls Bootstrap.List.
func GetBootstrapInfo(ctx context.Context, client *Client) ([]BootstrapinfoRecord, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Bootstrap.List(ctx)
}

// GetConfig calls ConfigMgmt.Get.
func GetConfig(ctx context.Context, client *Client, deviceIP string) (interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.ConfigMgmt.Get(ctx, deviceIP)
}

// GetConfigDiff calls ConfigMgmt.GetDiff.
func GetConfigDiff(ctx context.Context, client *Client, ip string) (interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.ConfigMgmt.GetDiff(ctx, ip)
}

// GetControllerVersion calls System.ControllerVersion.
func GetControllerVersion(ctx context.Context, client *Client) (*ControllerVersion, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.System.ControllerVersion(ctx)
}

// GetControllerVersionInternal calls System.ControllerVersionInternal.
func GetControllerVersionInternal(ctx context.Context, client *Client) (*ControllerVersion, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.System.ControllerVersionInternal(ctx)
}

// GetDay1ConfigStatus calls Intents.Day1ConfigStatus.
func GetDay1ConfigStatus(ctx context.Context, client *Client, intentName string) ([]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Intents.Day1ConfigStatus(ctx, intentName)
}

// GetDeviceBootstrapStages calls Bootstrap.DeviceStages.
func GetDeviceBootstrapStages(ctx context.Context, client *Client, bootstrapID int) ([]BootstrapStage, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Bootstrap.DeviceStages(ctx, bootstrapID)
}

// GetDevicesByLayer calls Fabrics.DevicesByLayer.
func GetDevicesByLayer(ctx context.Context, client *Client, layer string) ([]string, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.DevicesByLayer(ctx, layer)
}

// GetFabricByName calls Fabrics.Get.
func GetFabricByName(ctx context.Context, client *Client, name string) (*FabricItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.Get(ctx, name)
}

// GetFabricDevices calls Fabrics.DeviceIPs.
func GetFabricDevices(ctx context.Context, client *Client, fabricName string) ([]string, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.DeviceIPs(ctx, fabricName)
}

// GetFabricSimByName calls Fabricsims.Get.
func GetFabricSimByName(ctx context.Context, client *Client, name string) (*FabricSimItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabricsims.Get(ctx, name)
}

// GetFabrics calls Fabrics.ListDTOs.
func GetFabrics(ctx context.Context, client *Client) (*FabricsListResponse, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Fabrics.ListDTOs(ctx)
}

// GetFiles calls Files.List.
func GetFiles(ctx context.Context, client *Client, filetype string) (*GetFilesResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Files.List(ctx, filetype)
}

// GetFmInventory calls Devices.FMInventory.
func GetFmInventory(ctx context.Context, client *Client) ([]DeviceInventoryItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Devices.FMInventory(ctx)
}

// GetGpuAllocationHistory calls GPU.AllocationHistory.
func GetGpuAllocationHistory(ctx context.Context, client *Client, fabricName, tenantName string) ([]GpuAllocationHistoryItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.AllocationHistory(ctx, fabricName, tenantName)
}

// GetGpuTenantMappings calls GPU.TenantMappings.
func GetGpuTenantMappings(ctx context.Context, client *Client, fabricName string, tenantName *string) ([]GpuTenantMappingItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.TenantMappings(ctx, fabricName, tenantName)
}

// GetGpusByHost calls GPU.ByHost.
func GetGpusByHost(ctx context.Context, client *Client, hostName, fabricName string) ([]GpuItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.GPU.ByHost(ctx, hostName, fabricName)
}

// GetHostTenantsList calls HostTenants.List.
func GetHostTenantsList(ctx context.Context, client *Client, fabricName string) ([]HosttenantsRecord, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.HostTenants.List(ctx, fabricName)
}

// GetHostsList calls HostTenants.ListHosts.
func GetHostsList(ctx context.Context, client *Client, fabricName string) ([]HostItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.HostTenants.ListHosts(ctx, fabricName)
}

// GetImgmgmtStatus calls Devices.ImgmgmtStatus.
func GetImgmgmtStatus(ctx context.Context, client *Client, devices []string) ([]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Devices.ImgmgmtStatus(ctx, devices)
}

// GetIntentDerivationLogs calls Intents.DerivationLogs.
func GetIntentDerivationLogs(ctx context.Context, client *Client, device string) ([]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Intents.DerivationLogs(ctx, device)
}

// GetIntentValidation calls Intents.Validation.
func GetIntentValidation(ctx context.Context, client *Client, intentName string) ([]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Intents.Validation(ctx, intentName)
}

// GetInventoryByFabricName calls Inventory.ByFabric.
func GetInventoryByFabricName(ctx context.Context, client *Client, name string) ([]InventoryRecord, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.ByFabric(ctx, name)
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

// GetLastOrchestratedIntentName calls Intents.LastOrchestratedName.
func GetLastOrchestratedIntentName(ctx context.Context, client *Client) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Intents.LastOrchestratedName(ctx)
}

// GetLogLevel calls System.GetLogLevel.
func GetLogLevel(ctx context.Context, client *Client) (map[string]string, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.System.GetLogLevel(ctx)
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

// GetOperationWebhookStatus calls Operations.WebhookStatus.
func GetOperationWebhookStatus(ctx context.Context, client *Client, operationID string) (*WebhookDeliveryStatus, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Operations.WebhookStatus(ctx, operationID)
}

// GetRmaInfo calls RMA.List.
func GetRmaInfo(ctx context.Context, client *Client) ([]RMAInfoRecord, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.RMA.List(ctx)
}

// GetRmaStatus calls RMA.Status.
func GetRmaStatus(ctx context.Context, client *Client, rmaInfoID int) ([]RMAStatusItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.RMA.Status(ctx, rmaInfoID)
}

// GetStatus calls System.Status.
func GetStatus(ctx context.Context, client *Client, fileName *string) ([]string, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.System.Status(ctx, fileName)
}

// GetTenant calls Tenants.Get.
func GetTenant(ctx context.Context, client *Client, fabricName, tenantName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Tenants.Get(ctx, fabricName, tenantName)
}

// GetUIObject calls Intents.GetUIObject.
func GetUIObject(ctx context.Context, client *Client, name *string) (*IntentItem, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Intents.GetUIObject(ctx, name)
}

// GetUploadStatus calls System.UploadStatus.
func GetUploadStatus(ctx context.Context, client *Client) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.System.UploadStatus(ctx)
}

// GetVersion calls Devices.Versions.
func GetVersion(ctx context.Context, client *Client, devices []string) ([]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Devices.Versions(ctx, devices)
}

// InventorySync calls Inventory.Sync.
func InventorySync(ctx context.Context, client *Client, fabricName string) (map[string]interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Inventory.Sync(ctx, fabricName)
}

// IsDeviceAlive calls Devices.IsAlive.
func IsDeviceAlive(ctx context.Context, client *Client, hostIP string) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.IsAlive(ctx, hostIP)
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

// NetOpsDevice calls NetOps.Device.
func NetOpsDevice(ctx context.Context, client *Client, ipAddress string, action NetOpsAction, params NetOpsParams) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.NetOps.Device(ctx, ipAddress, action, params)
}

// NetOpsFabric calls NetOps.Fabric.
func NetOpsFabric(ctx context.Context, client *Client, fabricName string, action NetOpsAction, params NetOpsParams) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.NetOps.Fabric(ctx, fabricName, action, params)
}

// ProbeNmxcDomains calls NMXC.ProbeDomains.
func ProbeNmxcDomains(ctx context.Context, client *Client, fabricName string, domains []NmxcDomain) ([]NmxcProbeResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.ProbeDomains(ctx, fabricName, domains)
}

// RebootRequest calls Devices.Reboot.
func RebootRequest(ctx context.Context, client *Client, devices []string) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.Reboot(ctx, devices)
}

// RemoveDeviceFacts calls Devices.RemoveFacts.
func RemoveDeviceFacts(ctx context.Context, client *Client, devices []string) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.RemoveFacts(ctx, devices)
}

// ReplaceConfig calls ConfigMgmt.Replace.
func ReplaceConfig(ctx context.Context, client *Client, deviceIP, filePath string, onlydiff *bool) (interface{}, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.ConfigMgmt.Replace(ctx, deviceIP, filePath, onlydiff)
}

// ResetNmxcDomain calls NMXC.Reset.
func ResetNmxcDomain(ctx context.Context, client *Client, domainID string) (*NmxcResetResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.Reset(ctx, domainID)
}

// RestoreConfig calls ConfigMgmt.Restore.
func RestoreConfig(ctx context.Context, client *Client, items []DeviceConfigRestore) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.ConfigMgmt.Restore(ctx, items)
}

// SetLogLevel calls System.SetLogLevel.
func SetLogLevel(ctx context.Context, client *Client, loggers map[string]string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.System.SetLogLevel(ctx, loggers)
}

// StartStreaming calls System.StartStreaming.
func StartStreaming(ctx context.Context, client *Client, filename string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.System.StartStreaming(ctx, filename)
}

// StopStreaming calls System.StopStreaming.
func StopStreaming(ctx context.Context, client *Client, filename string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.System.StopStreaming(ctx, filename)
}

// TriggerBootstrapConfig calls Bootstrap.Trigger.
func TriggerBootstrapConfig(ctx context.Context, client *Client, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (*TriggerBootstrapResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Bootstrap.Trigger(ctx, batchName, subnet, netmask, gateway, bootstrapinfo)
}

// TriggerRma calls RMA.Trigger.
func TriggerRma(ctx context.Context, client *Client, items []RMAInfoItem) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.RMA.Trigger(ctx, items)
}

// UpdateFabricSimStatus calls Fabricsims.UpdateStatus.
func UpdateFabricSimStatus(ctx context.Context, client *Client, name string, status *string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabricsims.UpdateStatus(ctx, name, status)
}

// UpdateFabricStatus calls Fabrics.UpdateStatus.
func UpdateFabricStatus(ctx context.Context, client *Client, name string, status, intent, description, deploymentType *string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Fabrics.UpdateStatus(ctx, name, status, intent, description, deploymentType)
}

// UpdateHosts calls HostTenants.UpdateHosts.
func UpdateHosts(ctx context.Context, client *Client, hostnames []string, hostAction HostAction, tenantName, fabricName string) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.HostTenants.UpdateHosts(ctx, hostnames, hostAction, tenantName, fabricName)
}

// UpdateInventoryData calls Inventory.Update.
func UpdateInventoryData(ctx context.Context, client *Client, items []FabricInventoryUpdateItem) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Inventory.Update(ctx, items)
}

// UpdateNmxcDomains calls NMXC.UpdateDomains.
func UpdateNmxcDomains(ctx context.Context, client *Client, fabricName string, domains []NmxcDomain, operation GpuAction) (*NmxcDomainsUpdateResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.NMXC.UpdateDomains(ctx, fabricName, domains, operation)
}

// UpdateRoleInfo calls System.UpdateRoleInfo.
func UpdateRoleInfo(ctx context.Context, client *Client, layer int, currentName string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.System.UpdateRoleInfo(ctx, layer, currentName)
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

// UpgradeNOSImage calls Devices.UpgradeNOS.
func UpgradeNOSImage(ctx context.Context, client *Client, items []ImageUpgradeDetailsItem) (bool, error) {
	if client == nil {
		return false, ErrNilClient
	}
	return client.Devices.UpgradeNOS(ctx, items)
}

// UploadDay1Config calls ConfigMgmt.UploadDay1.
func UploadDay1Config(ctx context.Context, client *Client, filePath string, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.ConfigMgmt.UploadDay1(ctx, filePath, opts...)
}

// UploadFile calls Files.Upload.
func UploadFile(ctx context.Context, client *Client, filePath, filetype string, version, vendor, tag *string) (*UploadFileResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Files.Upload(ctx, filePath, filetype, version, vendor, tag)
}

// UploadUIObject calls Intents.UploadUIObject.
func UploadUIObject(ctx context.Context, client *Client, args *UploadUIObjectArgs, opts ...CallOption) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.Intents.UploadUIObject(ctx, args, opts...)
}

// ValidateUfmCreds calls Devices.ValidateUFMCreds.
func ValidateUfmCreds(ctx context.Context, client *Client, ufmURL, username, password string) (*UfmCredsResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.Devices.ValidateUFMCreds(ctx, ufmURL, username, password)
}
