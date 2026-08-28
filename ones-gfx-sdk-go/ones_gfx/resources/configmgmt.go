package resources

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// ConfigMgmtResource covers device configuration backup/restore/diff and
// Day-1 config upload.
type ConfigMgmtResource struct {
	transport *ones_gfx.Transport
}

// NewConfigMgmtResource constructs a ConfigMgmtResource.
func NewConfigMgmtResource(transport *ones_gfx.Transport) *ConfigMgmtResource {
	return &ConfigMgmtResource{transport: transport}
}

// DeviceConfigBackup mirrors Helper/DeviceConfigBackup.java.
type DeviceConfigBackup struct {
	IP    string `json:"ip"`
	Label string `json:"label,omitempty"`
}

// DeviceConfigRestore mirrors Helper/DeviceConfigRestore.java
// (timestamp format ddMMyyyyHHmmss).
type DeviceConfigRestore struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
}

// FetchBackupFilesResult is the response of FetchBackupFiles
// (POST /fetchdevicebackupfiles).
type FetchBackupFilesResult struct {
	FetchedConfig string `json:"fetchedConfig,omitempty"`
}

// Get returns the current running config of a device. The shape depends on
// the device NOS, so it is surfaced raw.
// Maps to GET /getconfig?deviceip=...
func (r *ConfigMgmtResource) Get(ctx context.Context, deviceIP string) (interface{}, error) {
	q := url.Values{}
	q.Set("deviceip", deviceIP)
	return ones_gfx.Call[interface{}](r.transport, "GET", "getconfig", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
}

// GetDiff returns the running-vs-saved config diff. Dynamic shape → raw.
// Maps to POST /getConfigDiff.
func (r *ConfigMgmtResource) GetDiff(ctx context.Context, ip string) (interface{}, error) {
	body := struct {
		IP string `json:"ip"`
	}{ip}
	return ones_gfx.Call[interface{}](r.transport, "POST", "getConfigDiff", body, ones_gfx.OperationModeSynchronous, nil)
}

// Replace pushes a replacement config file to a device.
// Maps to POST /replaceConfig (multipart). Dynamic result → raw.
func (r *ConfigMgmtResource) Replace(ctx context.Context, deviceIP, filePath string, onlydiff *bool) (interface{}, error) {
	fields := map[string]string{"deviceip": deviceIP}
	if onlydiff != nil {
		fields["onlydiff"] = strconv.FormatBool(*onlydiff)
	}
	return ones_gfx.CallMultipart[interface{}](r.transport, "POST", "replaceConfig", fields, map[string]string{"file": filePath}, ones_gfx.OperationModeSynchronous, nil)
}

// Restore restores saved configs by ip + timestamp.
// Maps to POST /restoreconfig.
func (r *ConfigMgmtResource) Restore(ctx context.Context, items []DeviceConfigRestore) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "restoreconfig", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Backup triggers a config backup with optional label per device.
// Maps to POST /backupConfig.
func (r *ConfigMgmtResource) Backup(ctx context.Context, items []DeviceConfigBackup) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "backupConfig", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// FetchBackupFiles returns backed-up config content for one ip + timestamp.
// Maps to POST /fetchdevicebackupfiles.
func (r *ConfigMgmtResource) FetchBackupFiles(ctx context.Context, ip, timestamp *string) (*FetchBackupFilesResult, error) {
	body := struct {
		IP        *string `json:"ip,omitempty"`
		Timestamp *string `json:"timestamp,omitempty"`
	}{ip, timestamp}
	return ones_gfx.Call[FetchBackupFilesResult](r.transport, "POST", "fetchdevicebackupfiles", body, ones_gfx.OperationModeSynchronous, nil)
}

// ListToRestore lists restorable backups. Returns nil when the server
// responds with null.
// Maps to POST /configslisttorestore.
func (r *ConfigMgmtResource) ListToRestore(ctx context.Context, devices []string, onlylimited *bool) (*string, error) {
	body := struct {
		Devices     []string `json:"devices,omitempty"`
		Onlylimited *bool    `json:"onlylimited,omitempty"`
	}{devices, onlylimited}
	return ones_gfx.Call[string](r.transport, "POST", "configslisttorestore", body, ones_gfx.OperationModeSynchronous, nil)
}

// UploadDay1 uploads a Day-1 config intent file. The server answers with the
// stored file name. Use WithRequestOrigin("ones-ui") to mark internal calls.
// Maps to POST /uploadDay1Config (multipart).
func (r *ConfigMgmtResource) UploadDay1(ctx context.Context, filePath string, opts ...ones_gfx.CallOption) (string, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.CallMultipart[string](r.transport, "POST", "uploadDay1Config", nil, map[string]string{"file": filePath}, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}
