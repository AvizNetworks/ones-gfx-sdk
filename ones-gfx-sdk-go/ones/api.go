package ones

// api.go exposes every Fabric Manager endpoint as a package-level function,
// mirroring the flat apis surface of the Python SDK. Each function resolves
// the shared client configured by Configure and enforces the auth guard, so
// calls fail with ErrNotAuthenticated until Login succeeds.
//
// Generated from the resource method set — edit the generator, not this file.

import "context"

// AddDeviceFacts calls Devices.AddFacts.
func AddDeviceFacts(ctx context.Context, items []DeviceDetail) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.AddFacts(ctx, items)
}

// AddHostTenantData calls HostTenants.Add.
func AddHostTenantData(ctx context.Context, name, fabricName string, description *string, hostsAllocated *int, vniID *int, configStatus *string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.HostTenants.Add(ctx, name, fabricName, description, hostsAllocated, vniID, configStatus)
}

// AddInventoryData calls Inventory.Add.
func AddInventoryData(ctx context.Context, items []InventoryItem) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Inventory.Add(ctx, items)
}

// AssignGpuPorts calls GPU.AssignPorts.
func AssignGpuPorts(ctx context.Context, fabricName, tenantName string, operation GpuAction, serverNames []string, gpuIds []int, membership *string) (*GpuPortAssignmentResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.AssignPorts(ctx, fabricName, tenantName, operation, serverNames, gpuIds, membership)
}

// AutoAllocateGpusToTenants calls Tenants.AutoAllocate.
func AutoAllocateGpusToTenants(ctx context.Context, fabricName, tenantName string, autoAllocationDevicesNeed int, suid SuidMap) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Tenants.AutoAllocate(ctx, fabricName, tenantName, autoAllocationDevicesNeed, suid)
}

// BackupConfig calls ConfigMgmt.Backup.
func BackupConfig(ctx context.Context, items []DeviceConfigBackup) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.ConfigMgmt.Backup(ctx, items)
}

// ConfigsListToRestore calls ConfigMgmt.ListToRestore.
func ConfigsListToRestore(ctx context.Context, devices []string, onlylimited *bool) (*string, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.ConfigMgmt.ListToRestore(ctx, devices, onlylimited)
}

// CreateFabric calls Fabrics.Create.
func CreateFabric(ctx context.Context, name string, args *FabricCreateArgs) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabrics.Create(ctx, name, args)
}

// CreateFabricSim calls Fabricsims.Create.
func CreateFabricSim(ctx context.Context, id *int, fabricName, simulationID, username, token, orgUUID, status, uiLink *string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabricsims.Create(ctx, id, fabricName, simulationID, username, token, orgUUID, status, uiLink)
}

// CreateTenant calls Tenants.Create.
func CreateTenant(ctx context.Context, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...CallOption) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Tenants.Create(ctx, fabricName, tenantName, description, maxGpusAllowed, shared, opts...)
}

// CreateTenantAsync calls Tenants.CreateAsync.
func CreateTenantAsync(ctx context.Context, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...CallOption) (*OperationAccepted, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.CreateAsync(ctx, fabricName, tenantName, description, maxGpusAllowed, shared, opts...)
}

// CreateVpcPeering calls VPCPeering.Create.
func CreateVpcPeering(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.VPCPeering.Create(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// CreateVpcPeeringAsync calls VPCPeering.CreateAsync.
func CreateVpcPeeringAsync(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (*OperationAccepted, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.VPCPeering.CreateAsync(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// DeleteFabric calls Fabrics.Delete.
func DeleteFabric(ctx context.Context, name string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabrics.Delete(ctx, name)
}

// DeleteFabricSim calls Fabricsims.Delete.
func DeleteFabricSim(ctx context.Context, name string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabricsims.Delete(ctx, name)
}

// DeleteFile calls Files.Delete.
func DeleteFile(ctx context.Context, id int) (*DeleteFileResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Files.Delete(ctx, id)
}

// DeleteHostTenantData calls HostTenants.Delete.
func DeleteHostTenantData(ctx context.Context, fabricName, tenantName string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.HostTenants.Delete(ctx, fabricName, tenantName)
}

// DeleteTenant calls Tenants.Delete.
func DeleteTenant(ctx context.Context, fabricName, tenantName string, opts ...CallOption) (*ApiResponseMessage, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.Delete(ctx, fabricName, tenantName, opts...)
}

// DeleteTenantAsync calls Tenants.DeleteAsync.
func DeleteTenantAsync(ctx context.Context, fabricName, tenantName string, opts ...CallOption) (*OperationAccepted, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.DeleteAsync(ctx, fabricName, tenantName, opts...)
}

// DeleteVpcPeering calls VPCPeering.Delete.
func DeleteVpcPeering(ctx context.Context, fabricName, name, vpcname, peervpcname string, opts ...CallOption) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.VPCPeering.Delete(ctx, fabricName, name, vpcname, peervpcname, opts...)
}

// EditFabric calls Fabrics.Update.
func EditFabric(ctx context.Context, fabricName string, args *FabricCreateArgs) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabrics.Update(ctx, fabricName, args)
}

// EditInventoryData calls Inventory.Edit.
func EditInventoryData(ctx context.Context, fabricName string, items []InventoryItem) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Inventory.Edit(ctx, fabricName, items)
}

// EnableZTPUpgrade calls Devices.EnableZTP.
func EnableZTPUpgrade(ctx context.Context, devices []string) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.EnableZTP(ctx, devices)
}

// FactoryResetNmxcDomain calls NMXC.FactoryReset.
func FactoryResetNmxcDomain(ctx context.Context, domainID string) (*NmxcFactoryResetResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.NMXC.FactoryReset(ctx, domainID)
}

// FetchDeviceBackupFiles calls ConfigMgmt.FetchBackupFiles.
func FetchDeviceBackupFiles(ctx context.Context, ip, timestamp *string) (*FetchBackupFilesResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.ConfigMgmt.FetchBackupFiles(ctx, ip, timestamp)
}

// FillBootstrapConfig calls Bootstrap.Fill.
func FillBootstrapConfig(ctx context.Context, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Bootstrap.Fill(ctx, batchName, subnet, netmask, gateway, bootstrapinfo)
}

// FillRmaConfig calls RMA.Fill.
func FillRmaConfig(ctx context.Context, items []RMAInfoItem) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.RMA.Fill(ctx, items)
}

// GetAllBootstrapBatches calls Bootstrap.ListBatches.
func GetAllBootstrapBatches(ctx context.Context) ([]map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Bootstrap.ListBatches(ctx)
}

// GetAllFabricSims calls Fabricsims.List.
func GetAllFabricSims(ctx context.Context) ([]FabricSimItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabricsims.List(ctx)
}

// GetAllFabrics calls Fabrics.List.
func GetAllFabrics(ctx context.Context) ([]FabricItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabrics.List(ctx)
}

// GetAllGpusList calls GPU.List.
func GetAllGpusList(ctx context.Context, fabricName string) ([]GpuItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.List(ctx, fabricName)
}

// GetAllInventory calls Inventory.ListAll.
func GetAllInventory(ctx context.Context) ([]InventoryRecord, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Inventory.ListAll(ctx)
}

// GetAvailableServers calls GPU.AvailableServers.
func GetAvailableServers(ctx context.Context, fabricName string) (*AvailableServerResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.AvailableServers(ctx, fabricName)
}

// GetBootstrapBatch calls Bootstrap.GetBatch.
func GetBootstrapBatch(ctx context.Context, batchName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Bootstrap.GetBatch(ctx, batchName)
}

// GetBootstrapInfo calls Bootstrap.List.
func GetBootstrapInfo(ctx context.Context) ([]BootstrapinfoRecord, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Bootstrap.List(ctx)
}

// GetConfig calls ConfigMgmt.Get.
func GetConfig(ctx context.Context, deviceIP string) (interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.ConfigMgmt.Get(ctx, deviceIP)
}

// GetConfigDiff calls ConfigMgmt.GetDiff.
func GetConfigDiff(ctx context.Context, ip string) (interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.ConfigMgmt.GetDiff(ctx, ip)
}

// GetControllerVersion calls System.ControllerVersion.
func GetControllerVersion(ctx context.Context) (*ControllerVersion, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.System.ControllerVersion(ctx)
}

// GetControllerVersionInternal calls System.ControllerVersionInternal.
func GetControllerVersionInternal(ctx context.Context) (*ControllerVersion, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.System.ControllerVersionInternal(ctx)
}

// GetDay1ConfigStatus calls Intents.Day1ConfigStatus.
func GetDay1ConfigStatus(ctx context.Context, intentName string) ([]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Intents.Day1ConfigStatus(ctx, intentName)
}

// GetDeviceBootstrapStages calls Bootstrap.DeviceStages.
func GetDeviceBootstrapStages(ctx context.Context, bootstrapID int) ([]BootstrapStage, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Bootstrap.DeviceStages(ctx, bootstrapID)
}

// GetDevicesByLayer calls Fabrics.DevicesByLayer.
func GetDevicesByLayer(ctx context.Context, layer string) ([]string, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabrics.DevicesByLayer(ctx, layer)
}

// GetFabricByName calls Fabrics.Get.
func GetFabricByName(ctx context.Context, name string) (*FabricItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabrics.Get(ctx, name)
}

// GetFabricDevices calls Fabrics.DeviceIPs.
func GetFabricDevices(ctx context.Context, fabricName string) ([]string, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabrics.DeviceIPs(ctx, fabricName)
}

// GetFabricSimByName calls Fabricsims.Get.
func GetFabricSimByName(ctx context.Context, name string) (*FabricSimItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabricsims.Get(ctx, name)
}

// GetFabrics calls Fabrics.ListDTOs.
func GetFabrics(ctx context.Context) (*FabricsListResponse, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Fabrics.ListDTOs(ctx)
}

// GetFiles calls Files.List.
func GetFiles(ctx context.Context, filetype string) (*GetFilesResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Files.List(ctx, filetype)
}

// GetFmInventory calls Devices.FMInventory.
func GetFmInventory(ctx context.Context) ([]DeviceInventoryItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Devices.FMInventory(ctx)
}

// GetGpuAllocationHistory calls GPU.AllocationHistory.
func GetGpuAllocationHistory(ctx context.Context, fabricName, tenantName string) ([]GpuAllocationHistoryItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.AllocationHistory(ctx, fabricName, tenantName)
}

// GetGpuTenantMappings calls GPU.TenantMappings.
func GetGpuTenantMappings(ctx context.Context, fabricName string, tenantName *string) ([]GpuTenantMappingItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.TenantMappings(ctx, fabricName, tenantName)
}

// GetGpusByHost calls GPU.ByHost.
func GetGpusByHost(ctx context.Context, hostName, fabricName string) ([]GpuItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.GPU.ByHost(ctx, hostName, fabricName)
}

// GetHostTenantsList calls HostTenants.List.
func GetHostTenantsList(ctx context.Context, fabricName string) ([]HosttenantsRecord, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.HostTenants.List(ctx, fabricName)
}

// GetHostsList calls HostTenants.ListHosts.
func GetHostsList(ctx context.Context, fabricName string) ([]HostItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.HostTenants.ListHosts(ctx, fabricName)
}

// GetImgmgmtStatus calls Devices.ImgmgmtStatus.
func GetImgmgmtStatus(ctx context.Context, devices []string) ([]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Devices.ImgmgmtStatus(ctx, devices)
}

// GetIntentDerivationLogs calls Intents.DerivationLogs.
func GetIntentDerivationLogs(ctx context.Context, device string) ([]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Intents.DerivationLogs(ctx, device)
}

// GetIntentValidation calls Intents.Validation.
func GetIntentValidation(ctx context.Context, intentName string) ([]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Intents.Validation(ctx, intentName)
}

// GetInventoryByFabricName calls Inventory.ByFabric.
func GetInventoryByFabricName(ctx context.Context, name string) ([]InventoryRecord, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Inventory.ByFabric(ctx, name)
}

// GetInventoryHosts calls Inventory.UFMHosts.
func GetInventoryHosts(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Inventory.UFMHosts(ctx, fabricName)
}

// GetInventoryPorts calls Inventory.Ports.
func GetInventoryPorts(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Inventory.Ports(ctx, fabricName)
}

// GetLastOrchestratedIntentName calls Intents.LastOrchestratedName.
func GetLastOrchestratedIntentName(ctx context.Context) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Intents.LastOrchestratedName(ctx)
}

// GetLogLevel calls System.GetLogLevel.
func GetLogLevel(ctx context.Context) (map[string]string, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.System.GetLogLevel(ctx)
}

// GetNmxcInventory calls NMXC.Inventory.
func GetNmxcInventory(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.NMXC.Inventory(ctx, fabricName)
}

// GetOperation calls Operations.Get.
func GetOperation(ctx context.Context, operationID string) (*OperationStatusItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Operations.Get(ctx, operationID)
}

// GetOperationWebhookStatus calls Operations.WebhookStatus.
func GetOperationWebhookStatus(ctx context.Context, operationID string) (*WebhookDeliveryStatus, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Operations.WebhookStatus(ctx, operationID)
}

// GetRmaInfo calls RMA.List.
func GetRmaInfo(ctx context.Context) ([]RMAInfoRecord, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.RMA.List(ctx)
}

// GetRmaStatus calls RMA.Status.
func GetRmaStatus(ctx context.Context, rmaInfoID int) ([]RMAStatusItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.RMA.Status(ctx, rmaInfoID)
}

// GetStatus calls System.Status.
func GetStatus(ctx context.Context, fileName *string) ([]string, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.System.Status(ctx, fileName)
}

// GetTenant calls Tenants.Get.
func GetTenant(ctx context.Context, fabricName, tenantName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.Get(ctx, fabricName, tenantName)
}

// GetUIObject calls Intents.GetUIObject.
func GetUIObject(ctx context.Context, name *string) (*IntentItem, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Intents.GetUIObject(ctx, name)
}

// GetUploadStatus calls System.UploadStatus.
func GetUploadStatus(ctx context.Context) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.System.UploadStatus(ctx)
}

// GetVersion calls Devices.Versions.
func GetVersion(ctx context.Context, devices []string) ([]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Devices.Versions(ctx, devices)
}

// InventorySync calls Inventory.Sync.
func InventorySync(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Inventory.Sync(ctx, fabricName)
}

// IsDeviceAlive calls Devices.IsAlive.
func IsDeviceAlive(ctx context.Context, hostIP string) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.IsAlive(ctx, hostIP)
}

// ListTenants calls Tenants.List.
func ListTenants(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.List(ctx, fabricName)
}

// ModifyGpuAllocations calls Tenants.ModifyAllocations.
func ModifyGpuAllocations(ctx context.Context, fabricName, tenantName string, suid SuidMap, operation *GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...CallOption) (*ApiResponseMessage, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.ModifyAllocations(ctx, fabricName, tenantName, suid, operation, configScope, unreachableDevices, opts...)
}

// ModifyGpuAllocationsAsync calls Tenants.ModifyAllocationsAsync.
func ModifyGpuAllocationsAsync(ctx context.Context, fabricName, tenantName string, suid SuidMap, operation *GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...CallOption) (*OperationAccepted, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.ModifyAllocationsAsync(ctx, fabricName, tenantName, suid, operation, configScope, unreachableDevices, opts...)
}

// NetOpsDevice calls NetOps.Device.
func NetOpsDevice(ctx context.Context, ipAddress string, action NetOpsAction, params NetOpsParams) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.NetOps.Device(ctx, ipAddress, action, params)
}

// NetOpsFabric calls NetOps.Fabric.
func NetOpsFabric(ctx context.Context, fabricName string, action NetOpsAction, params NetOpsParams) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.NetOps.Fabric(ctx, fabricName, action, params)
}

// ProbeNmxcDomains calls NMXC.ProbeDomains.
func ProbeNmxcDomains(ctx context.Context, fabricName string, domains []NmxcDomain) ([]NmxcProbeResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.NMXC.ProbeDomains(ctx, fabricName, domains)
}

// RebootRequest calls Devices.Reboot.
func RebootRequest(ctx context.Context, devices []string) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.Reboot(ctx, devices)
}

// RemoveDeviceFacts calls Devices.RemoveFacts.
func RemoveDeviceFacts(ctx context.Context, devices []string) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.RemoveFacts(ctx, devices)
}

// ReplaceConfig calls ConfigMgmt.Replace.
func ReplaceConfig(ctx context.Context, deviceIP, filePath string, onlydiff *bool) (interface{}, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.ConfigMgmt.Replace(ctx, deviceIP, filePath, onlydiff)
}

// ResetNmxcDomain calls NMXC.Reset.
func ResetNmxcDomain(ctx context.Context, domainID string) (*NmxcResetResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.NMXC.Reset(ctx, domainID)
}

// RestoreConfig calls ConfigMgmt.Restore.
func RestoreConfig(ctx context.Context, items []DeviceConfigRestore) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.ConfigMgmt.Restore(ctx, items)
}

// SetLogLevel calls System.SetLogLevel.
func SetLogLevel(ctx context.Context, loggers map[string]string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.System.SetLogLevel(ctx, loggers)
}

// StartStreaming calls System.StartStreaming.
func StartStreaming(ctx context.Context, filename string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.System.StartStreaming(ctx, filename)
}

// StopStreaming calls System.StopStreaming.
func StopStreaming(ctx context.Context, filename string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.System.StopStreaming(ctx, filename)
}

// TriggerBootstrapConfig calls Bootstrap.Trigger.
func TriggerBootstrapConfig(ctx context.Context, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (*TriggerBootstrapResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Bootstrap.Trigger(ctx, batchName, subnet, netmask, gateway, bootstrapinfo)
}

// TriggerRma calls RMA.Trigger.
func TriggerRma(ctx context.Context, items []RMAInfoItem) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.RMA.Trigger(ctx, items)
}

// UpdateFabricSimStatus calls Fabricsims.UpdateStatus.
func UpdateFabricSimStatus(ctx context.Context, name string, status *string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabricsims.UpdateStatus(ctx, name, status)
}

// UpdateFabricStatus calls Fabrics.UpdateStatus.
func UpdateFabricStatus(ctx context.Context, name string, status, intent, description, deploymentType *string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Fabrics.UpdateStatus(ctx, name, status, intent, description, deploymentType)
}

// UpdateHosts calls HostTenants.UpdateHosts.
func UpdateHosts(ctx context.Context, hostnames []string, hostAction HostAction, tenantName, fabricName string) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.HostTenants.UpdateHosts(ctx, hostnames, hostAction, tenantName, fabricName)
}

// UpdateInventoryData calls Inventory.Update.
func UpdateInventoryData(ctx context.Context, items []FabricInventoryUpdateItem) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Inventory.Update(ctx, items)
}

// UpdateNmxcDomains calls NMXC.UpdateDomains.
func UpdateNmxcDomains(ctx context.Context, fabricName string, domains []NmxcDomain, operation GpuAction) (*NmxcDomainsUpdateResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.NMXC.UpdateDomains(ctx, fabricName, domains, operation)
}

// UpdateRoleInfo calls System.UpdateRoleInfo.
func UpdateRoleInfo(ctx context.Context, layer int, currentName string) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.System.UpdateRoleInfo(ctx, layer, currentName)
}

// UpdateTenant calls Tenants.Update.
func UpdateTenant(ctx context.Context, fabricName, tenantName string, servers []GpuServerInfo, operation *GpuAction, opts ...CallOption) (*ApiResponseMessage, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.Update(ctx, fabricName, tenantName, servers, operation, opts...)
}

// UpdateTenantAsync calls Tenants.UpdateAsync.
func UpdateTenantAsync(ctx context.Context, fabricName, tenantName string, servers []GpuServerInfo, operation *GpuAction, opts ...CallOption) (*OperationAccepted, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Tenants.UpdateAsync(ctx, fabricName, tenantName, servers, operation, opts...)
}

// UpgradeNOSImage calls Devices.UpgradeNOS.
func UpgradeNOSImage(ctx context.Context, items []ImageUpgradeDetailsItem) (bool, error) {
	c, err := requireClient()
	if err != nil {
		return false, err
	}
	return c.Devices.UpgradeNOS(ctx, items)
}

// UploadDay1Config calls ConfigMgmt.UploadDay1.
func UploadDay1Config(ctx context.Context, filePath string, opts ...CallOption) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.ConfigMgmt.UploadDay1(ctx, filePath, opts...)
}

// UploadFile calls Files.Upload.
func UploadFile(ctx context.Context, filePath, filetype string, version, vendor, tag *string) (*UploadFileResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Files.Upload(ctx, filePath, filetype, version, vendor, tag)
}

// UploadUIObject calls Intents.UploadUIObject.
func UploadUIObject(ctx context.Context, args *UploadUIObjectArgs, opts ...CallOption) (string, error) {
	c, err := requireClient()
	if err != nil {
		return "", err
	}
	return c.Intents.UploadUIObject(ctx, args, opts...)
}

// ValidateUfmCreds calls Devices.ValidateUFMCreds.
func ValidateUfmCreds(ctx context.Context, ufmURL, username, password string) (*UfmCredsResult, error) {
	c, err := requireClient()
	if err != nil {
		return nil, err
	}
	return c.Devices.ValidateUFMCreds(ctx, ufmURL, username, password)
}
