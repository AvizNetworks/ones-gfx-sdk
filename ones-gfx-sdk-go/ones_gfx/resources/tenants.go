package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// TenantsResource covers the tenant CRUD, GPU allocation and auto-allocation
// endpoints under /api/fm.
type TenantsResource struct {
	transport *ones_gfx.Transport
}

// NewTenantsResource constructs a TenantsResource.
func NewTenantsResource(transport *ones_gfx.Transport) *TenantsResource {
	return &TenantsResource{transport: transport}
}

// SuidMap mirrors the Java GpuStatusUpdate.suid / AutoAllocateGpuRequest.suid:
// server index (as string; Jackson coerces to Integer) → hostname → GPU names.
type SuidMap = map[string]map[string][]string

// GpuServerInfo is one entry of the UpdateTenant.servers list
// (Helper/GpuServerInfo.java). Shared nil keeps the existing sharing mode.
type GpuServerInfo struct {
	ServerName string `json:"serverName"`
	Shared     *bool  `json:"shared,omitempty"`
}

// ConfigScope discriminates whole-server vs per-GPU configuration
// (GpuStatusUpdate.ConfigScope).
type ConfigScope string

const (
	ConfigScopeWholeServer   ConfigScope = "WHOLE_SERVER"
	ConfigScopeParticularGPU ConfigScope = "PARTICULAR_GPU"
)

// ApiResponseMessage is the {"status", "message"} body returned by
// ApiResponse.success(message) calls that carry no data payload.
type ApiResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// tenantWebhookFields is embedded by every body that supports webhook params.
type tenantWebhookFields struct {
	EnableWebhook bool     `json:"enableWebhook,omitempty"`
	WebhookURL    string   `json:"webhookUrl,omitempty"`
	WebhookEvents []string `json:"webhookEvents,omitempty"`
}

// Create registers a new tenant on a fabric (synchronous).
// Maps to POST /fabrics/{fabricName}/tenants.
// Returns the unwrapped success payload (tenant details string).
func (r *TenantsResource) Create(ctx context.Context, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...ones_gfx.CallOption) (string, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		TenantName     string  `json:"tenantName"`
		Description    *string `json:"description,omitempty"`
		MaxGpusAllowed *int    `json:"maxGpusAllowed,omitempty"`
		Shared         *bool   `json:"shared,omitempty"`
	}{tenantName, description, maxGpusAllowed, shared}
	res, err := ones_gfx.Call[string](r.transport, "POST", "fabrics/"+fabricName+"/tenants", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// CreateAsync registers a new tenant and returns the 202 operation handle.
// Poll with Operations.Get; use WithWebhook for push delivery and
// WithIdempotencyKey to make retries replay-safe.
func (r *TenantsResource) CreateAsync(ctx context.Context, fabricName, tenantName string, description *string, maxGpusAllowed *int, shared *bool, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		TenantName     string  `json:"tenantName"`
		Description    *string `json:"description,omitempty"`
		MaxGpusAllowed *int    `json:"maxGpusAllowed,omitempty"`
		Shared         *bool   `json:"shared,omitempty"`
		tenantWebhookFields
	}{TenantName: tenantName, Description: description, MaxGpusAllowed: maxGpusAllowed, Shared: shared}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "POST", "fabrics/"+fabricName+"/tenants", body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// List returns the tenants of a fabric.
// Maps to GET /fabrics/{fabricName}/tenants.
func (r *TenantsResource) List(ctx context.Context, fabricName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/tenants", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Get returns one tenant (detail incl. GPU assignment).
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) Get(ctx context.Context, fabricName, tenantName string) (map[string]interface{}, error) {
	res, err := ones_gfx.Call[map[string]interface{}](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName, nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Delete removes a tenant (synchronous).
// Maps to DELETE /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) Delete(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		tenantWebhookFields
	}{}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ApiResponseMessage](r.transport, "DELETE", "fabrics/"+fabricName+"/tenants/"+tenantName, body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// DeleteAsync removes a tenant and returns the 202 operation handle.
func (r *TenantsResource) DeleteAsync(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		tenantWebhookFields
	}{}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "DELETE", "fabrics/"+fabricName+"/tenants/"+tenantName, body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// Update adds or removes whole GPU servers on a tenant (synchronous).
// Maps to PATCH /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) Update(ctx context.Context, fabricName, tenantName string, servers []GpuServerInfo, operation *ones_gfx.GpuAction, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		Operation *ones_gfx.GpuAction `json:"operation,omitempty"`
		Servers   []GpuServerInfo     `json:"servers"`
		tenantWebhookFields
	}{Operation: operation, Servers: servers}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ApiResponseMessage](r.transport, "PATCH", "fabrics/"+fabricName+"/tenants/"+tenantName, body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// UpdateAsync adds or removes whole GPU servers and returns the 202 handle.
func (r *TenantsResource) UpdateAsync(ctx context.Context, fabricName, tenantName string, servers []GpuServerInfo, operation *ones_gfx.GpuAction, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		Operation *ones_gfx.GpuAction `json:"operation,omitempty"`
		Servers   []GpuServerInfo     `json:"servers"`
		tenantWebhookFields
	}{Operation: operation, Servers: servers}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "PATCH", "fabrics/"+fabricName+"/tenants/"+tenantName, body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// ModifyAllocations records per-GPU allocations for servers already attached
// to a tenant on an externally managed fabric (synchronous).
// Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations.
func (r *TenantsResource) ModifyAllocations(ctx context.Context, fabricName, tenantName string, suid SuidMap, operation *ones_gfx.GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...ones_gfx.CallOption) (*ApiResponseMessage, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		Suid               SuidMap             `json:"suid"`
		Operation          *ones_gfx.GpuAction `json:"operation,omitempty"`
		TenantName         string              `json:"tenantName"`
		FabricName         string              `json:"fabricName"`
		ConfigScope        *ConfigScope        `json:"configScope,omitempty"`
		UnreachableDevices []string            `json:"unreachableDevices,omitempty"`
		tenantWebhookFields
	}{Suid: suid, Operation: operation, TenantName: tenantName, FabricName: fabricName, ConfigScope: configScope, UnreachableDevices: unreachableDevices}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ApiResponseMessage](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/gpuAllocations", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// ModifyAllocationsAsync is the async variant of ModifyAllocations.
func (r *TenantsResource) ModifyAllocationsAsync(ctx context.Context, fabricName, tenantName string, suid SuidMap, operation *ones_gfx.GpuAction, configScope *ConfigScope, unreachableDevices []string, opts ...ones_gfx.CallOption) (*ones_gfx.OperationAccepted, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		Suid               SuidMap             `json:"suid"`
		Operation          *ones_gfx.GpuAction `json:"operation,omitempty"`
		TenantName         string              `json:"tenantName"`
		FabricName         string              `json:"fabricName"`
		ConfigScope        *ConfigScope        `json:"configScope,omitempty"`
		UnreachableDevices []string            `json:"unreachableDevices,omitempty"`
		tenantWebhookFields
	}{Suid: suid, Operation: operation, TenantName: tenantName, FabricName: fabricName, ConfigScope: configScope, UnreachableDevices: unreachableDevices}
	attachTenantWebhook(&body.tenantWebhookFields, rc)
	return ones_gfx.Call[ones_gfx.OperationAccepted](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/gpuAllocations", body, ones_gfx.OperationModeAsyncPoll, rc.ReqOpts())
}

// AutoAllocate automatically picks GPUs for a tenant.
// Maps to POST /autoAllocateGpusToTenants.
func (r *TenantsResource) AutoAllocate(ctx context.Context, fabricName, tenantName string, autoAllocationDevicesNeed int, suid SuidMap) (bool, error) {
	body := struct {
		FabricName                string  `json:"fabricName"`
		TenantName                string  `json:"tenantName"`
		AutoAllocationDevicesNeed int     `json:"autoAllocationDevicesNeed"`
		Suid                      SuidMap `json:"suid,omitempty"`
	}{fabricName, tenantName, autoAllocationDevicesNeed, suid}
	res, err := ones_gfx.Call[bool](r.transport, "POST", "autoAllocateGpusToTenants", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// attachTenantWebhook fills webhook body fields when WithWebhook was supplied.
func attachTenantWebhook(fields *tenantWebhookFields, rc ones_gfx.ResolvedCall) {
	if rc.WebhookURL == "" {
		return
	}
	fields.EnableWebhook = true
	fields.WebhookURL = rc.WebhookURL
	fields.WebhookEvents = rc.WebhookEvents
}
