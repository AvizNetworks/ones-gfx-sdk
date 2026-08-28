package resources

import (
	"context"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// BootstrapResource covers bootstrap batch staging, triggering and queries.
type BootstrapResource struct {
	transport *ones_gfx.Transport
}

// NewBootstrapResource constructs a BootstrapResource.
func NewBootstrapResource(transport *ones_gfx.Transport) *BootstrapResource {
	return &BootstrapResource{transport: transport}
}

// BootstrapsubnetItem mirrors Models/Bootstrapsubnet.java.
type BootstrapsubnetItem struct {
	Subnet  string `json:"subnet,omitempty"`
	Netmask string `json:"netmask,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// BootstrapinfoItem mirrors Models/Bootstrapinfo.java as a request payload.
type BootstrapinfoItem struct {
	DeviceMacAddress string               `json:"deviceMacAddress,omitempty"`
	Serial           string               `json:"serial,omitempty"`
	DeviceIP         string               `json:"deviceIp,omitempty"`
	Hostname         string               `json:"hostname,omitempty"`
	Region           string               `json:"region,omitempty"`
	Layer            string               `json:"layer,omitempty"`
	Azid             string               `json:"azid,omitempty"`
	Rackid           string               `json:"rackid,omitempty"`
	Brickid          string               `json:"brickid,omitempty"`
	Groupid          string               `json:"groupid,omitempty"`
	Paramspath       string               `json:"paramspath,omitempty"`
	Bootfilepath     string               `json:"bootfilepath,omitempty"`
	NosImagePath     string               `json:"nosImagePath,omitempty"`
	AgentImagePath   string               `json:"agentImagePath,omitempty"`
	FmcliImagePath   string               `json:"fmcliImagePath,omitempty"`
	BaseConfigDbPath string               `json:"baseConfigDbPath,omitempty"`
	BaseConfigFmPath string               `json:"baseConfigFmPath,omitempty"`
	Status           *int                 `json:"status,omitempty"`
	BatchName        string               `json:"batchName,omitempty"`
	CollectorIP      string               `json:"collectorIP,omitempty"`
	Fabricid         string               `json:"fabricid,omitempty"`
	Vendor           string               `json:"vendor,omitempty"`
	Bootstrapsubnet  *BootstrapsubnetItem `json:"bootstrapsubnet,omitempty"`
}

// BootstrapinfoRecord mirrors Models/Bootstrapinfo.java JSON (responses).
type BootstrapinfoRecord struct {
	ID               int                  `json:"id,omitempty"`
	DeviceMacAddress string               `json:"deviceMacAddress,omitempty"`
	Serial           string               `json:"serial,omitempty"`
	DeviceIP         string               `json:"deviceIp,omitempty"`
	Hostname         string               `json:"hostname,omitempty"`
	Region           string               `json:"region,omitempty"`
	Layer            string               `json:"layer,omitempty"`
	Azid             string               `json:"azid,omitempty"`
	Rackid           string               `json:"rackid,omitempty"`
	Brickid          string               `json:"brickid,omitempty"`
	Groupid          string               `json:"groupid,omitempty"`
	Paramspath       string               `json:"paramspath,omitempty"`
	Bootfilepath     string               `json:"bootfilepath,omitempty"`
	NosImagePath     string               `json:"nosImagePath,omitempty"`
	AgentImagePath   string               `json:"agentImagePath,omitempty"`
	FmcliImagePath   string               `json:"fmcliImagePath,omitempty"`
	BaseConfigDbPath string               `json:"baseConfigDbPath,omitempty"`
	BaseConfigFmPath string               `json:"baseConfigFmPath,omitempty"`
	Status           *int                 `json:"status,omitempty"`
	BatchName        string               `json:"batchName,omitempty"`
	CollectorIP      string               `json:"collectorIP,omitempty"`
	Fabricid         string               `json:"fabricid,omitempty"`
	Vendor           string               `json:"vendor,omitempty"`
	Bootstrapsubnet  *BootstrapsubnetItem `json:"bootstrapsubnet,omitempty"`
	Lastupdated      string               `json:"lastupdated,omitempty"`
}

// TriggerBootstrapResult is the response of Trigger
// (POST /triggerbootstrapconfig).
type TriggerBootstrapResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// BootstrapStage is one entry of DeviceStages. Shape derived from
// ServiceProvider.getDeviceBootstrapStages (SES_FM ServiceProvider.java),
// which puts task/tasktitle/status/starttime/endtime/logs.
type BootstrapStage struct {
	Task      string `json:"task,omitempty"`
	TaskTitle string `json:"tasktitle,omitempty"`
	Status    *int   `json:"status,omitempty"`
	StartTime string `json:"starttime,omitempty"`
	EndTime   string `json:"endtime,omitempty"`
	Logs      string `json:"logs,omitempty"`
}

// bootstrapDetails mirrors Helper/BootstrapDetails.java.
type bootstrapDetails struct {
	BatchName     *string             `json:"batchName,omitempty"`
	Subnet        *string             `json:"subnet,omitempty"`
	Netmask       *string             `json:"netmask,omitempty"`
	Gateway       *string             `json:"gateway,omitempty"`
	Bootstrapinfo []BootstrapinfoItem `json:"bootstrapinfo,omitempty"`
}

// Fill stages a bootstrap batch without starting it.
// Maps to POST /fillbootstrapconfig.
func (r *BootstrapResource) Fill(ctx context.Context, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (bool, error) {
	body := bootstrapDetails{batchName, subnet, netmask, gateway, bootstrapinfo}
	res, err := ones_gfx.Call[bool](r.transport, "POST", "fillbootstrapconfig", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Trigger runs the bootstrap process synchronously.
// Maps to POST /triggerbootstrapconfig.
func (r *BootstrapResource) Trigger(ctx context.Context, batchName, subnet, netmask, gateway *string, bootstrapinfo []BootstrapinfoItem) (*TriggerBootstrapResult, error) {
	body := bootstrapDetails{batchName, subnet, netmask, gateway, bootstrapinfo}
	return ones_gfx.Call[TriggerBootstrapResult](r.transport, "POST", "triggerbootstrapconfig", body, ones_gfx.OperationModeSynchronous, nil)
}

// List returns all bootstrap devices. Maps to GET /getbootstrapinfo.
func (r *BootstrapResource) List(ctx context.Context) ([]BootstrapinfoRecord, error) {
	res, err := ones_gfx.Call[[]BootstrapinfoRecord](r.transport, "GET", "getbootstrapinfo", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ListBatches returns bootstrap batch summaries. Open map payload per batch.
// Maps to GET /getAllBootstrapBatches.
func (r *BootstrapResource) ListBatches(ctx context.Context) ([]map[string]interface{}, error) {
	res, err := ones_gfx.Call[[]map[string]interface{}](r.transport, "GET", "getAllBootstrapBatches", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// GetBatch returns one batch's details incl. devices. Open map payload.
// Maps to GET /getBootstrapBatch/{batchName}.
func (r *BootstrapResource) GetBatch(ctx context.Context, batchName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "getBootstrapBatch/"+batchName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// DeviceStages returns bootstrap stage progress for one device
// (STAGE10, STAGE20, ...). Maps to GET /getDeviceBootstrapStages/{bootstrapId}.
func (r *BootstrapResource) DeviceStages(ctx context.Context, bootstrapID int) ([]BootstrapStage, error) {
	res, err := ones_gfx.Call[[]BootstrapStage](r.transport, "GET", "getDeviceBootstrapStages/"+strconv.Itoa(bootstrapID), nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
