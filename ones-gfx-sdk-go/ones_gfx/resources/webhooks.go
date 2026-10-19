package resources

import (
	"context"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// WebhooksResource covers webhook registration at all three scopes: global
// (/webhooks), per-fabric and per-tenant.
//
// The scopes are independent and there is no precedence between them: every
// registration matching an event receives it, so a callback registered at two
// scopes is delivered twice.
type WebhooksResource struct {
	transport *ones_gfx.Transport
}

// NewWebhooksResource constructs a WebhooksResource.
func NewWebhooksResource(transport *ones_gfx.Transport) *WebhooksResource {
	return &WebhooksResource{transport: transport}
}

// Canonical event types FM delivers.
const (
	WebhookEventTenantCreate     = "tenant.create"
	WebhookEventTenantDelete     = "tenant.delete"
	WebhookEventTenantAllocate   = "tenant.allocate"
	WebhookEventTenantDeallocate = "tenant.deallocate"
	WebhookEventTenantVPCPeering = "tenant.vpcpeering"
	WebhookEventFabricDeploy     = "fabric.deploy"
	WebhookEventFabricDelete     = "fabric.delete"
)

// WebhookEventTypes lists every canonical event name. FM also accepts the
// fabric.tenant_* aliases and resolves them to these on write.
var WebhookEventTypes = []string{
	WebhookEventTenantCreate,
	WebhookEventTenantDelete,
	WebhookEventTenantAllocate,
	WebhookEventTenantDeallocate,
	WebhookEventTenantVPCPeering,
	WebhookEventFabricDeploy,
	WebhookEventFabricDelete,
}

// WebhookFilters narrows which events a global registration receives. Empty
// lists mean no filter.
type WebhookFilters struct {
	FabricNames []string `json:"fabricNames,omitempty"`
	TenantNames []string `json:"tenantNames,omitempty"`
}

// WebhookRegistration is a registration as FM returns it.
//
// ID is the surrogate key, present on tenant- and fabric-scoped rows.
// WebhookID is the external id a global row is addressed by and is the one to
// hold on to; the two never both appear. TenantName and FabricName are
// omitted at the scopes they do not apply to, and Filters appears only on a
// global registration.
type WebhookRegistration struct {
	ID          int64           `json:"id,omitempty"`
	WebhookID   string          `json:"webhookId,omitempty"`
	TenantName  string          `json:"tenantName,omitempty"`
	FabricName  string          `json:"fabricName,omitempty"`
	CallbackURL string          `json:"callbackUrl"`
	Enabled     bool            `json:"enabled"`
	EventTypes  []string        `json:"eventTypes,omitempty"`
	Filters     *WebhookFilters `json:"filters,omitempty"`
	CreatedAt   string          `json:"createdAt,omitempty"`
	UpdatedAt   string          `json:"updatedAt,omitempty"`
}

// ---------------------------------------------------------------- global ---

// Register registers a global webhook covering every fabric and tenant unless
// narrowed by filters.
// Maps to POST /webhooks.
//
// Idempotent on callbackURL: re-registering the same URL updates that
// subscription rather than adding a duplicate. FM answers 201 when it created
// one and 200 when it returned or updated an existing one, so a configurator
// that re-registers on every restart neither fails nor accumulates rows.
//
// Unlike the fabric and tenant scopes, eventTypes is required and must be
// non-empty here: a global subscription to every event is almost never
// intended, so FM rejects it rather than delivering it. Pass a nil filters to
// receive matching events from everywhere.
func (r *WebhooksResource) Register(ctx context.Context, callbackURL string, eventTypes []string, enabled *bool, filters *WebhookFilters, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL string          `json:"callbackUrl"`
		EventTypes  []string        `json:"eventTypes"`
		Enabled     *bool           `json:"enabled,omitempty"`
		Filters     *WebhookFilters `json:"filters,omitempty"`
	}{callbackURL, eventTypes, enabled, filters}
	return ones_gfx.Call[WebhookRegistration](r.transport, "POST", "webhooks", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// List returns every global webhook registration, newest first. Fabric- and
// tenant-scoped registrations are not listed here.
// Maps to GET /webhooks.
func (r *WebhooksResource) List(ctx context.Context, opts ...ones_gfx.CallOption) ([]WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[[]WebhookRegistration](r.transport, "GET", "webhooks", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Get returns one global webhook registration, or 404 when webhookID is
// unknown.
// Maps to GET /webhooks/{webhookId}.
func (r *WebhooksResource) Get(ctx context.Context, webhookID string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[WebhookRegistration](r.transport, "GET", "webhooks/"+webhookID, nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Update partially updates a global registration: only the fields sent change.
// Maps to PATCH /webhooks/{webhookId}.
//
// The semantics differ from the fabric and tenant scopes, which replace the
// whole body with PUT. Here a nil eventTypes leaves the subscription
// untouched, while an explicitly empty slice clears it to "every event type";
// the same distinction applies to each filter list.
func (r *WebhooksResource) Update(ctx context.Context, webhookID string, callbackURL *string, enabled *bool, eventTypes []string, filters *WebhookFilters, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL *string         `json:"callbackUrl,omitempty"`
		Enabled     *bool           `json:"enabled,omitempty"`
		EventTypes  []string        `json:"eventTypes,omitempty"`
		Filters     *WebhookFilters `json:"filters,omitempty"`
	}{callbackURL, enabled, eventTypes, filters}
	return ones_gfx.Call[WebhookRegistration](r.transport, "PATCH", "webhooks/"+webhookID, body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// Delete removes a global registration. FM answers 204, so the result is nil
// on success; 404 when webhookID is unknown.
// Maps to DELETE /webhooks/{webhookId}.
func (r *WebhooksResource) Delete(ctx context.Context, webhookID string, opts ...ones_gfx.CallOption) error {
	rc := ones_gfx.ResolveCallOptions(opts...)
	_, err := ones_gfx.Call[WebhookRegistration](r.transport, "DELETE", "webhooks/"+webhookID, nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	return err
}

// ---------------------------------------------------------------- fabric ---

// RegisterFabric registers, or replaces, the webhook for one fabric; FM
// answers 201. A fabric holds at most one registration. Pass a nil eventTypes
// to receive every event for the fabric.
// Maps to POST /fabrics/{fabricName}/webhook.
func (r *WebhooksResource) RegisterFabric(ctx context.Context, fabricName, callbackURL string, enabled *bool, eventTypes []string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL string   `json:"callbackUrl"`
		Enabled     *bool    `json:"enabled,omitempty"`
		EventTypes  []string `json:"eventTypes,omitempty"`
	}{callbackURL, enabled, eventTypes}
	return ones_gfx.Call[WebhookRegistration](r.transport, "POST", "fabrics/"+fabricName+"/webhook", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// GetFabric returns the fabric's registration, or 404 when none is registered.
// Maps to GET /fabrics/{fabricName}/webhook.
func (r *WebhooksResource) GetFabric(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[WebhookRegistration](r.transport, "GET", "fabrics/"+fabricName+"/webhook", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// UpdateFabric updates the fabric's registration; 404 when none is registered.
// Maps to PUT /fabrics/{fabricName}/webhook.
//
// enabled is sent on every call because FM reads it as a plain boolean here
// rather than as "leave unchanged when absent". An empty eventTypes slice
// clears the filter, subscribing the registration to every event.
func (r *WebhooksResource) UpdateFabric(ctx context.Context, fabricName string, callbackURL *string, enabled bool, eventTypes []string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL *string  `json:"callbackUrl,omitempty"`
		Enabled     bool     `json:"enabled"`
		EventTypes  []string `json:"eventTypes,omitempty"`
	}{callbackURL, enabled, eventTypes}
	return ones_gfx.Call[WebhookRegistration](r.transport, "PUT", "fabrics/"+fabricName+"/webhook", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// DeleteFabric removes the fabric's registration. FM answers 204; 404 when
// none is registered.
// Maps to DELETE /fabrics/{fabricName}/webhook.
func (r *WebhooksResource) DeleteFabric(ctx context.Context, fabricName string, opts ...ones_gfx.CallOption) error {
	rc := ones_gfx.ResolveCallOptions(opts...)
	_, err := ones_gfx.Call[WebhookRegistration](r.transport, "DELETE", "fabrics/"+fabricName+"/webhook", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	return err
}

// ---------------------------------------------------------------- tenant ---

// RegisterTenant registers, or replaces, the webhook for one tenant; FM
// answers 201. A tenant holds at most one registration. Pass a nil eventTypes
// to receive every event for the tenant.
// Maps to POST /fabrics/{fabricName}/tenants/{tenantName}/webhook.
func (r *WebhooksResource) RegisterTenant(ctx context.Context, fabricName, tenantName, callbackURL string, enabled *bool, eventTypes []string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL string   `json:"callbackUrl"`
		Enabled     *bool    `json:"enabled,omitempty"`
		EventTypes  []string `json:"eventTypes,omitempty"`
	}{callbackURL, enabled, eventTypes}
	return ones_gfx.Call[WebhookRegistration](r.transport, "POST", "fabrics/"+fabricName+"/tenants/"+tenantName+"/webhook", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// GetTenant returns the tenant's registration, or 404 when none is registered.
// Maps to GET /fabrics/{fabricName}/tenants/{tenantName}/webhook.
func (r *WebhooksResource) GetTenant(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	return ones_gfx.Call[WebhookRegistration](r.transport, "GET", "fabrics/"+fabricName+"/tenants/"+tenantName+"/webhook", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// UpdateTenant updates the tenant's registration; 404 when none is registered.
// Maps to PUT /fabrics/{fabricName}/tenants/{tenantName}/webhook.
//
// enabled is sent on every call because FM reads it as a plain boolean here
// rather than as "leave unchanged when absent". An empty eventTypes slice
// clears the filter, subscribing the registration to every event.
func (r *WebhooksResource) UpdateTenant(ctx context.Context, fabricName, tenantName string, callbackURL *string, enabled bool, eventTypes []string, opts ...ones_gfx.CallOption) (*WebhookRegistration, error) {
	rc := ones_gfx.ResolveCallOptions(opts...)
	body := struct {
		CallbackURL *string  `json:"callbackUrl,omitempty"`
		Enabled     bool     `json:"enabled"`
		EventTypes  []string `json:"eventTypes,omitempty"`
	}{callbackURL, enabled, eventTypes}
	return ones_gfx.Call[WebhookRegistration](r.transport, "PUT", "fabrics/"+fabricName+"/tenants/"+tenantName+"/webhook", body, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
}

// DeleteTenant removes the tenant's registration. FM answers 204; 404 when
// none is registered.
// Maps to DELETE /fabrics/{fabricName}/tenants/{tenantName}/webhook.
func (r *WebhooksResource) DeleteTenant(ctx context.Context, fabricName, tenantName string, opts ...ones_gfx.CallOption) error {
	rc := ones_gfx.ResolveCallOptions(opts...)
	_, err := ones_gfx.Call[WebhookRegistration](r.transport, "DELETE", "fabrics/"+fabricName+"/tenants/"+tenantName+"/webhook", nil, ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	return err
}
