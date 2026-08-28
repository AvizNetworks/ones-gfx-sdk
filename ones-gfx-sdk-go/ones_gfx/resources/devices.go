package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// DevicesResource covers device management: facts, versions, upgrades,
// aliveness and UFM credential validation.
type DevicesResource struct {
	transport *ones_gfx.Transport
}

// NewDevicesResource constructs a DevicesResource.
func NewDevicesResource(transport *ones_gfx.Transport) *DevicesResource {
	return &DevicesResource{transport: transport}
}

// DeviceDetail mirrors Helper/DeviceDetail.java (also the base of
// ImageUpgradeDetailsItem below).
type DeviceDetail struct {
	IP       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// ImageUpgradeDetailsItem mirrors Helper/ImageUpgradeDetails.java.
type ImageUpgradeDetailsItem struct {
	DeviceDetail
	PathToImage string `json:"pathToImage"`
}

// DeviceInventoryItem mirrors Helper/DeviceInventory.java.
type DeviceInventoryItem struct {
	IPAddress  string `json:"ipAddress,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	Layer      string `json:"layer,omitempty"`
	FabricName string `json:"fabricName,omitempty"`
}

// UfmCredsResult is the response of ValidateUFMCreds (POST /ValidateUfmCreds).
type UfmCredsResult struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// AddFacts registers devices. Maps to POST /addDeviceFacts.
func (r *DevicesResource) AddFacts(ctx context.Context, items []DeviceDetail) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "addDeviceFacts", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// RemoveFacts removes device records by IP. Maps to POST /removeDeviceFacts.
func (r *DevicesResource) RemoveFacts(ctx context.Context, devices []string) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "removeDeviceFacts", devices, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Versions returns device facts and NOS versions. Dynamic per-device payload.
// Maps to POST /getVersion.
func (r *DevicesResource) Versions(ctx context.Context, devices []string) ([]interface{}, error) {
	res, err := ones_gfx.Call[[]interface{}](r.transport, "POST", "getVersion", devices, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ImgmgmtStatus returns image-management status per device. Dynamic payload.
// Maps to POST /getImgmgmtStatus.
func (r *DevicesResource) ImgmgmtStatus(ctx context.Context, devices []string) ([]interface{}, error) {
	res, err := ones_gfx.Call[[]interface{}](r.transport, "POST", "getImgmgmtStatus", devices, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// EnableZTP enables ZTP and runs the upgrade. Maps to POST /enableZTPUpgrade.
func (r *DevicesResource) EnableZTP(ctx context.Context, devices []string) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "enableZTPUpgrade", devices, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// UpgradeNOS starts a NOS image upgrade. Maps to POST /upgradeNOSImage.
func (r *DevicesResource) UpgradeNOS(ctx context.Context, items []ImageUpgradeDetailsItem) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "upgradeNOSImage", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Reboot reboots devices. Maps to POST /rebootRequest.
func (r *DevicesResource) Reboot(ctx context.Context, devices []string) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "rebootRequest", devices, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// IsAlive pings one device. Maps to GET /isDeviceAlive?hostIP=...
func (r *DevicesResource) IsAlive(ctx context.Context, hostIP string) (bool, error) {
	q := url.Values{}
	q.Set("hostIP", hostIP)
	res, err := ones_gfx.Call[bool](r.transport, "GET", "isDeviceAlive", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// FMInventory returns the FM inventory list. Maps to GET /fm-Inventory.
func (r *DevicesResource) FMInventory(ctx context.Context) ([]DeviceInventoryItem, error) {
	res, err := ones_gfx.Call[[]DeviceInventoryItem](r.transport, "GET", "fm-Inventory", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ValidateUFMCreds checks UFM URL/username/password connectivity.
// Maps to POST /ValidateUfmCreds.
func (r *DevicesResource) ValidateUFMCreds(ctx context.Context, ufmURL, username, password string) (*UfmCredsResult, error) {
	body := struct {
		UfmURL   string `json:"ufmUrl"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{ufmURL, username, password}
	return ones_gfx.Call[UfmCredsResult](r.transport, "POST", "ValidateUfmCreds", body, ones_gfx.OperationModeSynchronous, nil)
}
