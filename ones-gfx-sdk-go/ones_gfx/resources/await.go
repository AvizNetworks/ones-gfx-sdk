package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// Fabric Manager exposes two async contracts, and they are layered rather
// than alternatives. Which one answers your question depends on what you are
// asking.
//
//	Operations API  — did Fabric Manager accept and process my request?
//	config_status   — did the resulting configuration reach the devices?
//
// They are not interchangeable, because an operation reaches SUCCESS before
// the device configuration is applied. Creating a tenant asynchronously,
// AsyncTenantService marks the operation SUCCESS as soon as the tenant row is
// saved and its IPAM pool provisioned; at that moment the tenant's
// config_status is still "2" (NOT_STARTED). The device push runs afterwards
// and is what drives config_status 2 -> 3 -> 1 (or 0 on failure).
//
// So:
//
//   - A 202 gives an operationId. Poll it (WaitForOperation) or receive the
//     webhook. Terminal means the request was processed, nothing more.
//   - Whether the fabric actually carries the change is config_status.
//     Poll it (WaitForTenantConfigApplied).
//   - Needing both is the normal case. AwaitTenantChange does them in order.
//
// A synchronous call collapses the first layer only: it returns once Fabric
// Manager has processed the request, and config_status may still be settling.

// WaitOptions bounds a polling loop. Zero values take the defaults.
type WaitOptions struct {
	// Interval between polls. Default 5s.
	Interval time.Duration
	// Timeout for the whole wait. Default 10m.
	Timeout time.Duration
}

func (w WaitOptions) withDefaults() WaitOptions {
	if w.Interval <= 0 {
		w.Interval = 5 * time.Second
	}
	if w.Timeout <= 0 {
		w.Timeout = 10 * time.Minute
	}
	return w
}

// WaitTimeoutError reports that a wait ran out of time before the thing it
// watched reached a terminal state. The work may still be in progress.
type WaitTimeoutError struct {
	What       string
	LastStatus string
	Waited     time.Duration
}

func (e *WaitTimeoutError) Error() string {
	return fmt.Sprintf("timed out after %s waiting for %s; last status %q (it may still be in progress)",
		e.Waited, e.What, e.LastStatus)
}

func (e *WaitTimeoutError) ONESError() {}

// WaitForOperation polls an operation until it reaches a terminal status.
//
// This answers "did Fabric Manager process my request". It does **not** mean
// the device configuration has been applied — see WaitForTenantConfigApplied.
//
// A terminal FAILURE is returned as the operation plus a nil error: the wait
// succeeded, the operation did not. Inspect Status and ErrorMessage.
func (r *OperationsResource) WaitForOperation(ctx context.Context, operationID string, opts WaitOptions) (*OperationStatusItem, error) {
	if err := pathName("operationId", operationID); err != nil {
		return nil, err
	}
	o := opts.withDefaults()
	started := time.Now()
	deadline := started.Add(o.Timeout)
	last := ""

	for {
		op, err := r.Get(ctx, operationID)
		if isValidation(err) {
			return nil, err // bad input will not fix itself by polling
		}
		if err == nil && op != nil {
			last = op.Status
			if ones_gfx.OperationStatus(op.Status).IsTerminal() {
				return op, nil
			}
		}
		if !time.Now().Add(o.Interval).Before(deadline) {
			return nil, &WaitTimeoutError{
				What:       "operation " + operationID,
				LastStatus: last,
				Waited:     time.Since(started).Round(time.Second),
			}
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(o.Interval):
		}
	}
}

// WaitForTenantConfigApplied polls a tenant until its config_status settles.
//
// This answers "did the configuration reach the devices", which an operation
// status cannot tell you. It returns the tenant at the point it settled;
// check ConfigStatus to see how.
//
// A terminal status is any of PASS, FAIL, BROWN_FIELD or SKIPPED. The last
// two mean the device was deliberately not configured, which is an outcome,
// not a failure — so this returns them rather than waiting for a PASS that
// will never come.
func (r *TenantsResource) WaitForTenantConfigApplied(ctx context.Context, fabricName, tenantName string, opts WaitOptions) (*ones_gfx.Tenant, error) {
	if err := fabricTenant(fabricName, tenantName); err != nil {
		return nil, err
	}
	o := opts.withDefaults()
	started := time.Now()
	deadline := started.Add(o.Timeout)
	last := ones_gfx.ConfigStatusNotStarted

	for {
		tn, err := r.Get(ctx, fabricName, tenantName)
		if isValidation(err) {
			return nil, err // bad input will not fix itself by polling
		}
		if err == nil && tn != nil {
			last = tn.ConfigStatus
			if last.IsTerminal() {
				return tn, nil
			}
		}
		if !time.Now().Add(o.Interval).Before(deadline) {
			return nil, &WaitTimeoutError{
				What:       fmt.Sprintf("tenant %s/%s config_status", fabricName, tenantName),
				LastStatus: last.String(),
				Waited:     time.Since(started).Round(time.Second),
			}
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(o.Interval):
		}
	}
}

// TenantChangeResult is the outcome of both async layers.
type TenantChangeResult struct {
	// Operation is the processed operation, nil when the change was made
	// synchronously and there was no operationId to follow.
	Operation *OperationStatusItem
	// Tenant is the tenant once its config_status settled.
	Tenant *ones_gfx.Tenant
}

// Applied reports whether the configuration actually reached the devices.
func (r *TenantChangeResult) Applied() bool {
	return r != nil && r.Tenant != nil && r.Tenant.ConfigStatus.IsSuccess()
}

// AwaitTenantChange runs both async layers in the order Fabric Manager
// completes them: first the operation, then the tenant's config_status.
//
// Pass the operationId from a 202 response, or "" for a change made
// synchronously — in that case only config_status is awaited, which is the
// layer a synchronous call leaves unresolved.
//
// An operation that ends FAILURE stops the wait: there is no point watching
// config_status for a change Fabric Manager never attempted.
func (r *TenantsResource) AwaitTenantChange(
	ctx context.Context,
	ops *OperationsResource,
	operationID, fabricName, tenantName string,
	opts WaitOptions,
) (*TenantChangeResult, error) {
	out := &TenantChangeResult{}

	if operationID != "" && ops != nil {
		op, err := ops.WaitForOperation(ctx, operationID, opts)
		if err != nil {
			return out, err
		}
		out.Operation = op
		if ones_gfx.OperationStatus(op.Status).IsFailure() {
			return out, fmt.Errorf("operation %s failed: %s", operationID, op.ErrorMessage)
		}
	}

	tn, err := r.WaitForTenantConfigApplied(ctx, fabricName, tenantName, opts)
	if err != nil {
		return out, err
	}
	out.Tenant = tn
	return out, nil
}
