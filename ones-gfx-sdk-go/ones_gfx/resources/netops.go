package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// NetOpsResource covers device/fabric NetOps operations (VLAN, port change,
// VLAN member configuration).
type NetOpsResource struct {
	transport *ones_gfx.Transport
}

// NewNetOpsResource constructs a NetOpsResource.
func NewNetOpsResource(transport *ones_gfx.Transport) *NetOpsResource {
	return &NetOpsResource{transport: transport}
}

// NetOpsAction is the operation selector of NetOpsRequest
// (Helper/NetOpsRequest.java).
type NetOpsAction string

const (
	NetOpsActionVlan       NetOpsAction = "vlan"
	NetOpsActionPortChange NetOpsAction = "portchange"
	NetOpsActionVlanMember NetOpsAction = "vlanmember"
)

// NetOpsParams carries the operation-specific parameters of a NetOpsRequest.
// The accepted keys depend on Action — the Java service validates them:
//
//	vlan:       vlanId, operation
//	portchange: portList, state
//	vlanmember: vlanId, memberPorts, operation, taggingMode
//
// Values are strings, ints or lists of strings.
type NetOpsParams = map[string]interface{}

// Device executes a NetOps operation on one device.
// Maps to POST /netops/{ipAddress}.
func (r *NetOpsResource) Device(ctx context.Context, ipAddress string, action NetOpsAction, params NetOpsParams) (bool, error) {
	body := struct {
		Action NetOpsAction `json:"action"`
		Params NetOpsParams `json:"params"`
	}{action, params}
	res, err := ones_gfx.Call[bool](r.transport, "POST", "netops/"+ipAddress, body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Fabric executes a NetOps operation across an entire fabric.
// Maps to POST /netopsfabric/{fabricName}.
func (r *NetOpsResource) Fabric(ctx context.Context, fabricName string, action NetOpsAction, params NetOpsParams) (bool, error) {
	body := struct {
		Action NetOpsAction `json:"action"`
		Params NetOpsParams `json:"params"`
	}{action, params}
	res, err := ones_gfx.Call[bool](r.transport, "POST", "netopsfabric/"+fabricName, body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}
