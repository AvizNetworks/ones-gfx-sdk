package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// TenantsResource handles tenant lifecycle operations.
type TenantsResource struct {
	transport *ones_gfx.Transport
}

// NewTenantsResource constructs a TenantsResource.
func NewTenantsResource(transport *ones_gfx.Transport) *TenantsResource {
	return &TenantsResource{transport: transport}
}

// CreateTenantRequest holds parameters for tenant creation.
type CreateTenantRequest struct {
	Name           string
	Description    string
	MaxGPUsAllowed int
	Shared         bool
}

// ServerSpec describes a GPU server to allocate/deallocate.
type ServerSpec struct {
	ServerName string
	Shared     bool
}

// GpuPortAssignmentRequest is the request for POST .../gpus.
// GPUIDs are integer port IDs (1, 2, 3, ...).
type GpuPortAssignmentRequest struct {
	Operation   string   // "ADD" or "DELETE"
	ServerNames []string
	GPUIDs      []int
}

// ServerSpecsFromNames converts server hostnames into ServerSpec entries.
func ServerSpecsFromNames(names []string) []ServerSpec {
	servers := make([]ServerSpec, 0, len(names))
	for _, name := range names {
		servers = append(servers, ServerSpec{ServerName: name})
	}
	return servers
}

// List retrieves all tenants in a fabric.
// Maps to GET /fabrics/{fabricName}/tenants.
func (r *TenantsResource) List(ctx context.Context, fabricName string) ([]ones_gfx.Tenant, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}

	path := fmt.Sprintf("fabrics/%s/tenants", fabricName)
	body, err := r.transport.Get(path, nil)
	if err != nil {
		return nil, err
	}

	obj, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", body)
	}

	tenantsRaw, ok := obj["tenants"]
	if !ok {
		return nil, fmt.Errorf("response missing 'tenants' field")
	}

	tenantsList, ok := tenantsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'tenants' field is not an array")
	}

	tenants := make([]ones_gfx.Tenant, 0, len(tenantsList))
	for _, item := range tenantsList {
		jsonBytes, err := remarshal(item)
		if err != nil {
			continue
		}
		var tenant ones_gfx.Tenant
		if err := unmarshalJSON(jsonBytes, &tenant); err != nil {
			continue
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// Get retrieves a single tenant by name.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) Get(ctx context.Context, fabricName, tenantName string) (*ones_gfx.Tenant, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return nil, fmt.Errorf("tenantName is required")
	}

	path := fmt.Sprintf("fabrics/%s/tenants/%s", fabricName, tenantName)
	body, err := r.transport.Get(path, nil)
	if err != nil {
		return nil, err
	}

	// Response wraps tenant in {"tenant": {...}}
	obj, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", body)
	}

	tenantData, ok := obj["tenant"]
	if !ok {
		tenantData = body
	}

	jsonBytes, err := remarshal(tenantData)
	if err != nil {
		return nil, err
	}

	var tenant ones_gfx.Tenant
	if err := unmarshalJSON(jsonBytes, &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

// AvailableServers lists GPU servers available for allocation.
// Maps to GET /fabrics/{fabricName}/available_servers.
func (r *TenantsResource) AvailableServers(ctx context.Context, fabricName string) ([]string, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}

	path := fmt.Sprintf("fabrics/%s/available_servers", fabricName)
	body, err := r.transport.Get(path, nil)
	if err != nil {
		return nil, err
	}

	obj, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", body)
	}

	serversRaw, ok := obj["availableGPUs"]
	if !ok {
		return []string{}, nil
	}

	serversList, ok := serversRaw.([]interface{})
	if !ok {
		return []string{}, nil
	}

	servers := make([]string, 0, len(serversList))
	for _, s := range serversList {
		if str, ok := s.(string); ok {
			servers = append(servers, str)
		}
	}

	return servers, nil
}

// Create creates a tenant (synchronous mode).
// Maps to POST /fabrics/{fabricName}/tenants.
func (r *TenantsResource) Create(ctx context.Context, fabricName string, req CreateTenantRequest, opts ...ones_gfx.CallOption) (*ones_gfx.Tenant, error) {
	if err := validateCreateRequest(fabricName, req); err != nil {
		return nil, err
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	if webhookURL != "" || len(webhookEvents) > 0 {
		return nil, fmt.Errorf("webhook options are only valid for async methods")
	}
	body := buildCreateBody(req)
	path := fmt.Sprintf("fabrics/%s/tenants", fabricName)

	result, err := r.transport.Post(path, body, ones_gfx.OperationModeSynchronous, timeout)
	if err != nil {
		return nil, err
	}

	return extractTenant(result)
}

// CreateAsync creates a tenant (async mode).
// Returns an Operation handle for polling or webhook delivery.
func (r *TenantsResource) CreateAsync(ctx context.Context, fabricName string, req CreateTenantRequest, opts ...ones_gfx.CallOption) (*ones_gfx.Operation, error) {
	if err := validateCreateRequest(fabricName, req); err != nil {
		return nil, err
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	body := buildCreateBody(req)
	if err := attachWebhookFields(body, webhookURL, webhookEvents); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("fabrics/%s/tenants", fabricName)
	mode := ones_gfx.OperationModeAsyncPoll
	if webhookURL != "" {
		mode = ones_gfx.OperationModeAsyncWebhook
	}

	result, err := r.transport.Post(path, body, mode, timeout)
	if err != nil {
		return nil, err
	}

	return extractOperation(result)
}

// Delete deletes a tenant (synchronous mode).
// Maps to DELETE /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) Delete(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) error {
	if fabricName == "" {
		return fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return fmt.Errorf("tenantName is required")
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	if webhookURL != "" || len(webhookEvents) > 0 {
		return fmt.Errorf("webhook options are only valid for async methods")
	}
	path := fmt.Sprintf("fabrics/%s/tenants/%s", fabricName, tenantName)

	_, err := r.transport.Delete(path, nil, ones_gfx.OperationModeSynchronous, timeout)
	return err
}

// DeleteAsync deletes a tenant (async mode).
func (r *TenantsResource) DeleteAsync(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) (*ones_gfx.Operation, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return nil, fmt.Errorf("tenantName is required")
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	path := fmt.Sprintf("fabrics/%s/tenants/%s", fabricName, tenantName)

	var body map[string]interface{}
	if webhookURL != "" || len(webhookEvents) > 0 {
		body = make(map[string]interface{})
		if err := attachWebhookFields(body, webhookURL, webhookEvents); err != nil {
			return nil, err
		}
	}

	mode := ones_gfx.OperationModeAsyncPoll
	if webhookURL != "" {
		mode = ones_gfx.OperationModeAsyncWebhook
	}

	result, err := r.transport.Delete(path, body, mode, timeout)
	if err != nil {
		return nil, err
	}

	return extractOperation(result)
}

// AllocateGPUs adds GPU servers to a tenant (synchronous mode).
// Maps to PATCH /fabrics/{fabricName}/tenants/{tenantName}.
func (r *TenantsResource) AllocateGPUs(ctx context.Context, fabricName, tenantName string, servers []ServerSpec, opts ...ones_gfx.CallOption) error {
	return r.gpuUpdate(ctx, fabricName, tenantName, "ADD", servers, opts...)
}

// AllocateGPUsAsync adds GPU servers (async mode).
func (r *TenantsResource) AllocateGPUsAsync(ctx context.Context, fabricName, tenantName string, servers []ServerSpec, opts ...ones_gfx.CallOption) (*ones_gfx.Operation, error) {
	return r.gpuUpdateAsync(ctx, fabricName, tenantName, "ADD", servers, opts...)
}

// DeallocateGPUs removes GPU servers from a tenant (synchronous mode).
func (r *TenantsResource) DeallocateGPUs(ctx context.Context, fabricName, tenantName string, servers []ServerSpec, opts ...ones_gfx.CallOption) error {
	return r.gpuUpdate(ctx, fabricName, tenantName, "DELETE", servers, opts...)
}

// DeallocateGPUsAsync removes GPU servers (async mode).
func (r *TenantsResource) DeallocateGPUsAsync(ctx context.Context, fabricName, tenantName string, servers []ServerSpec, opts ...ones_gfx.CallOption) (*ones_gfx.Operation, error) {
	return r.gpuUpdateAsync(ctx, fabricName, tenantName, "DELETE", servers, opts...)
}

// gpuUpdate handles sync GPU allocate/deallocate.
func (r *TenantsResource) gpuUpdate(ctx context.Context, fabricName, tenantName, operation string, servers []ServerSpec, opts ...ones_gfx.CallOption) error {
	if fabricName == "" {
		return fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return fmt.Errorf("tenantName is required")
	}
	if len(servers) == 0 {
		return fmt.Errorf("servers list cannot be empty")
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	if webhookURL != "" || len(webhookEvents) > 0 {
		return fmt.Errorf("webhook options are only valid for async methods")
	}
	body := buildGPUUpdateBody(operation, servers)
	path := fmt.Sprintf("fabrics/%s/tenants/%s", fabricName, tenantName)

	_, err := r.transport.Patch(path, body, ones_gfx.OperationModeSynchronous, timeout)
	return err
}

// gpuUpdateAsync handles async GPU allocate/deallocate.
func (r *TenantsResource) gpuUpdateAsync(ctx context.Context, fabricName, tenantName, operation string, servers []ServerSpec, opts ...ones_gfx.CallOption) (*ones_gfx.Operation, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return nil, fmt.Errorf("tenantName is required")
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("servers list cannot be empty")
	}

	timeout, webhookURL, webhookEvents := ones_gfx.ApplyCallOptions(opts...)
	body := buildGPUUpdateBody(operation, servers)
	if err := attachWebhookFields(body, webhookURL, webhookEvents); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("fabrics/%s/tenants/%s", fabricName, tenantName)
	mode := ones_gfx.OperationModeAsyncPoll
	if webhookURL != "" {
		mode = ones_gfx.OperationModeAsyncWebhook
	}

	result, err := r.transport.Patch(path, body, mode, timeout)
	if err != nil {
		return nil, err
	}

	return extractOperation(result)
}

// AssignPorts assigns specific ports to a tenant.
// Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/gpus with operation="ADD".
func (r *TenantsResource) AssignPorts(ctx context.Context, fabricName, tenantName string, req GpuPortAssignmentRequest) error {
	req.Operation = "ADD"
	return r.postGPUs(ctx, fabricName, tenantName, req)
}

// UnassignPorts removes specific ports from a tenant.
// Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/gpus with operation="DELETE".
func (r *TenantsResource) UnassignPorts(ctx context.Context, fabricName, tenantName string, req GpuPortAssignmentRequest) error {
	req.Operation = "DELETE"
	return r.postGPUs(ctx, fabricName, tenantName, req)
}

// postGPUs is the shared implementation for AssignPorts and UnassignPorts.
func (r *TenantsResource) postGPUs(ctx context.Context, fabricName, tenantName string, req GpuPortAssignmentRequest) error {
	if fabricName == "" {
		return fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return fmt.Errorf("tenantName is required")
	}
	if len(req.ServerNames) == 0 {
		return fmt.Errorf("serverNames cannot be empty")
	}

	path := fmt.Sprintf("fabrics/%s/tenants/%s/gpus", fabricName, tenantName)
	_, err := r.transport.Post(path, buildAssignPortsBody(req), ones_gfx.OperationModeSynchronous, nil)
	return err
}

// Helper functions

func validateCreateRequest(fabricName string, req CreateTenantRequest) error {
	if fabricName == "" {
		return fmt.Errorf("fabricName is required")
	}
	if req.Name == "" {
		return fmt.Errorf("tenant name is required")
	}
	if req.MaxGPUsAllowed != ones_gfx.UnlimitedGPUs && req.MaxGPUsAllowed < 1 {
		return fmt.Errorf("maxGPUsAllowed must be -1 (unlimited) or >= 1, got %d", req.MaxGPUsAllowed)
	}
	return nil
}

func buildCreateBody(req CreateTenantRequest) map[string]interface{} {
	return map[string]interface{}{
		"tenantName":     req.Name,
		"description":    req.Description,
		"maxGpusAllowed": req.MaxGPUsAllowed,
		"shared":         req.Shared,
	}
}

func buildAssignPortsBody(req GpuPortAssignmentRequest) map[string]interface{} {
	body := map[string]interface{}{
		"operation":   req.Operation,
		"serverNames": req.ServerNames,
	}
	if len(req.GPUIDs) > 0 {
		body["gpuIds"] = req.GPUIDs
	}
	return body
}

func buildGPUUpdateBody(operation string, servers []ServerSpec) map[string]interface{} {
	serverMaps := make([]map[string]interface{}, len(servers))
	for i, s := range servers {
		serverMaps[i] = map[string]interface{}{
			"serverName": s.ServerName,
			"shared":     s.Shared,
		}
	}
	return map[string]interface{}{
		"operation": operation,
		"servers":   serverMaps,
	}
}

func attachWebhookFields(body map[string]interface{}, webhookURL string, webhookEvents []string) error {
	if webhookURL == "" && len(webhookEvents) == 0 {
		return nil
	}
	if webhookURL == "" || len(webhookEvents) == 0 {
		return fmt.Errorf("webhookURL and webhookEvents must be provided together")
	}
	parsed, err := url.Parse(webhookURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("webhookURL must be an absolute URL with scheme and host, got %q", webhookURL)
	}
	for _, event := range webhookEvents {
		if strings.TrimSpace(event) == "" {
			return fmt.Errorf("webhookEvents must contain non-empty strings")
		}
	}
	body["enableWebhook"] = true
	body["webhookUrl"] = webhookURL
	body["webhookEvents"] = webhookEvents
	return nil
}

func extractTenant(result interface{}) (*ones_gfx.Tenant, error) {
	if _, ok := result.(string); ok {
		return &ones_gfx.Tenant{}, nil
	}
	obj, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	tenantData, ok := obj["tenant"]
	if !ok {
		tenantData = result
	}

	jsonBytes, err := remarshal(tenantData)
	if err != nil {
		return nil, err
	}

	var tenant ones_gfx.Tenant
	if err := unmarshalJSON(jsonBytes, &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

func extractOperation(result interface{}) (*ones_gfx.Operation, error) {
	jsonBytes, err := remarshal(result)
	if err != nil {
		return nil, err
	}

	var op ones_gfx.Operation
	if err := unmarshalJSON(jsonBytes, &op); err != nil {
		return nil, err
	}

	return &op, nil
}

// remarshal converts interface{} back to JSON bytes for unmarshaling into structs.
func remarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
