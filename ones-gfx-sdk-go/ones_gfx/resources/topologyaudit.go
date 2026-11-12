package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// TopologyAuditResource triggers and reads back the LLDP-vs-intent Topology
// Audit (GFX-6219/GFX-6225). It owns exactly one endpoint — the trigger.
// Status/Result otherwise ride the existing Operations API: Status is a
// named alias for Operations.Get, and Result parses the completed
// operation's JSON payload into a typed report. There is no synchronous
// variant; Fabric Manager always runs the audit in the background.
type TopologyAuditResource struct {
	transport *ones_gfx.Transport
	ops       *OperationsResource
}

// NewTopologyAuditResource constructs a TopologyAuditResource.
func NewTopologyAuditResource(transport *ones_gfx.Transport, ops *OperationsResource) *TopologyAuditResource {
	return &TopologyAuditResource{transport: transport, ops: ops}
}

// ErrTopologyAuditNotDone is returned by Result while the operation is still
// PENDING or RUNNING. It is not a failure — poll again, or use
// Operations.WaitForOperation to block until terminal.
var ErrTopologyAuditNotDone = errors.New("topology audit operation has not finished yet")

// TopologyAuditLink is one link's outcome (SES_FM
// TopologyCheckResult.LinkCheckResult).
//
// Status is one of PASS (expected and observed agree), MISSING_LINK (LLDP
// collected fine but the expected neighbor wasn't seen), MISMATCH (a
// different device and/or interface was observed), UNVERIFIABLE (LLDP
// couldn't be collected/trusted at all), UNEXPECTED_LINK (LLDP shows a
// neighbor with no corresponding intent entry — RemoteDevice/RemoteInterface
// are empty in this case), or SKIPPED (host/DPU/out-of-intent peer, not
// auditable by name). ObservedRemoteDevice/ObservedRemoteInterface are only
// populated once a neighbor was seen at all.
type TopologyAuditLink struct {
	Device                  string `json:"device"`
	LocalInterface          string `json:"localInterface"`
	RemoteDevice            string `json:"remoteDevice,omitempty"`
	RemoteInterface         string `json:"remoteInterface,omitempty"`
	ObservedRemoteDevice    string `json:"observedRemoteDevice,omitempty"`
	ObservedRemoteInterface string `json:"observedRemoteInterface,omitempty"`
	Status                  string `json:"status"`
	Message                 string `json:"message,omitempty"`
}

// TopologyAuditReport is the parsed `result` payload of a completed (SUCCESS)
// topology-audit operation (SES_FM TopologyAuditService.buildReport).
//
// OverallStatus is INCOMPLETE if any link is UNVERIFIABLE (a collection
// failure must never be reported as a clean PASS), else FAIL if any link is
// MISSING_LINK/MISMATCH/UNEXPECTED_LINK, else PASS. DeviceIPMap maps
// hostname -> management IP for devices the audit observed.
type TopologyAuditReport struct {
	FabricName    string              `json:"fabricName"`
	IntentVersion string              `json:"intentVersion"`
	OverallStatus string              `json:"overallStatus"`
	Links         []TopologyAuditLink `json:"links"`
	DeviceIPMap   map[string]string   `json:"deviceIpMap,omitempty"`
}

// Trigger kicks off an async topology audit for fabricName and returns the
// 202 operation handle. Poll it with Status or Result, or block with
// Operations.WaitForOperation.
// Maps to POST /fabrics/{fabricName}/topology-audit.
//
// FM answers 404 when the fabric or its intent cannot be resolved.
func (r *TopologyAuditResource) Trigger(ctx context.Context, fabricName string) (*ones_gfx.OperationAccepted, error) {
	if err := pathName("fabricName", fabricName); err != nil {
		return nil, err
	}
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "POST", "fabrics/"+fabricName+"/topology-audit", nil, ones_gfx.OperationModeAsyncPoll, nil)
}

// Status polls the audit's operation state. It is a named alias for
// Operations.Get — Topology Audit has no status endpoint of its own.
// Maps to GET /operations/{operationId}.
func (r *TopologyAuditResource) Status(ctx context.Context, operationID string) (*OperationStatusItem, error) {
	return r.ops.Get(ctx, operationID)
}

// Result fetches the operation and parses its completed report.
//
// Returns ErrTopologyAuditNotDone while the operation is PENDING or
// RUNNING — call again later, or use Operations.WaitForOperation to block
// until terminal. Returns a plain error (wrapping the operation's
// ErrorMessage) if the audit itself FAILED/FAILURE (e.g. an unresolvable
// fabric/intent, or an exception during LLDP collection) — there is no
// report to parse in that case.
func (r *TopologyAuditResource) Result(ctx context.Context, operationID string) (*TopologyAuditReport, error) {
	op, err := r.ops.Get(ctx, operationID)
	if err != nil {
		return nil, err
	}
	status := ones_gfx.OperationStatus(op.Status)
	switch {
	case status.IsFailure():
		return nil, fmt.Errorf("topology audit operation %s failed: %s", operationID, op.ErrorMessage)
	case !status.IsTerminal():
		return nil, ErrTopologyAuditNotDone
	}
	var report TopologyAuditReport
	if err := json.Unmarshal([]byte(op.Result), &report); err != nil {
		return nil, fmt.Errorf("topology audit operation %s: malformed result: %w", operationID, err)
	}
	return &report, nil
}
