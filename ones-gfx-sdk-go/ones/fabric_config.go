package ones

// fabric_config.go exposes the ONES UI config service — fabric generation,
// deployment and deletion — as package-level functions, alongside the
// generated Fabric Manager functions in api.go.

import (
	"context"
	"encoding/json"
)

// CreateFabricConfig generates a fabric (topology, YAML and inventory) in
// Draft state. Calls FabricConfig.Create (POST /api/config/addFabricData).
//
// Unlike CreateFabric (Fabric Manager addFabricData, which only stores the
// fabric record), this is what the ONES UI uses, and it is the starting
// point for DeployFabric.
func CreateFabricConfig(ctx context.Context, client *Client, name string, args FabricConfigArgs) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.Create(ctx, name, args)
}

// GetFabricYaml returns a fabric's generated YAML. Calls FabricConfig.GetYAML.
func GetFabricYaml(ctx context.Context, client *Client, name string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.GetYAML(ctx, name)
}

// DeleteFabricConfig deletes a fabric; already-deleted is not an error.
// Calls FabricConfig.Delete (DELETE /api/config/deletefabricbyname/{name}).
func DeleteFabricConfig(ctx context.Context, client *Client, name string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.Delete(ctx, name)
}

// UploadFabricDeviceIPs calls FabricConfig.UploadDeviceIPs.
func UploadFabricDeviceIPs(ctx context.Context, client *Client, fabricName string, devices []FabricDevice) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.UploadDeviceIPs(ctx, fabricName, devices)
}

// ValidateFabricSwitches calls FabricConfig.ValidateSwitches.
func ValidateFabricSwitches(ctx context.Context, client *Client, fabricName string, devices []FabricDevice) ([]ValidateDeviceResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.FabricConfig.ValidateSwitches(ctx, fabricName, devices)
}

// ValidateFabricServers calls FabricConfig.ValidateServers.
func ValidateFabricServers(ctx context.Context, client *Client, fabricName string, devices []FabricDevice) ([]ValidateDeviceResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.FabricConfig.ValidateServers(ctx, fabricName, devices)
}

// UpdateFabricInventory calls FabricConfig.UpdateInventory.
func UpdateFabricInventory(ctx context.Context, client *Client, fabricName string, devices []FabricDevice) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.UpdateInventory(ctx, fabricName, devices)
}

// PushFabricConfig calls FabricConfig.PushConfig.
func PushFabricConfig(ctx context.Context, client *Client, fabricName, instance string, data json.RawMessage) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.PushConfig(ctx, fabricName, instance, data)
}

// UpdateFabricConfigStatus calls FabricConfig.UpdateStatus.
func UpdateFabricConfigStatus(ctx context.Context, client *Client, name, status, description, deploymentType string) (string, error) {
	if client == nil {
		return "", ErrNilClient
	}
	return client.FabricConfig.UpdateStatus(ctx, name, status, description, deploymentType)
}

// DeployFabric runs the full deploy sequence for a fabric made with
// CreateFabricConfig. Calls FabricConfig.Deploy; see it for the steps.
func DeployFabric(ctx context.Context, client *Client, fabricName string, devices []FabricDevice, opts *FabricDeployOptions) (*FabricDeployResult, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return client.FabricConfig.Deploy(ctx, fabricName, devices, opts)
}
