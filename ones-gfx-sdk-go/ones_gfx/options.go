package ones_gfx

import (
	"crypto/tls"
	"time"
)

// CallOption configures optional parameters for SDK method calls.
type CallOption func(*callConfig)

// callConfig holds per-call configuration that can override client defaults.
type callConfig struct {
	timeout        *time.Duration
	webhookURL     string
	webhookEvents  []string
	idempotencyKey string
	requestOrigin  string
	anycastGateway *AnycastGateway
}

// defaultCallConfig returns a config with no overrides (uses client defaults).
func defaultCallConfig() callConfig {
	return callConfig{}
}

// WithIdempotencyKey sends an Idempotency-Key header so a retried async call
// replays the original operation instead of creating a duplicate.
//
// Example:
//
//	client.Tenants.CreateAsync(ctx, "fab1", "t1", nil, nil, nil,
//	    ones_gfx.WithIdempotencyKey("ONES-t1-CRT-20260828"))
func WithIdempotencyKey(key string) CallOption {
	return func(cfg *callConfig) {
		cfg.idempotencyKey = key
	}
}

// WithRequestOrigin sends the x-request-origin header. The server treats
// "ones-ui" as an INTERNAL caller; anything else is EXTERNAL.
func WithRequestOrigin(origin string) CallOption {
	return func(cfg *callConfig) {
		cfg.requestOrigin = origin
	}
}

// AnycastGateway is the request-side anycastGateway block accepted by tenant
// creation and by shared-services creation (GFX-5611). FM rejects any key but
// these two and requires Enabled to be a real JSON boolean; leave MACAddress
// empty to take the fabric's own MAC.
//
// The response field of the same name is a gateway address string, not this
// object.
type AnycastGateway struct {
	Enabled    bool   `json:"enabled"`
	MACAddress string `json:"macAddress,omitempty"`
}

// WithAnycastGateway sends an anycastGateway block on calls that accept one
// (tenant create, shared-services create). Applies to HBN DPU fabrics in
// tenant DHCP mode; elsewhere FM answers ANYCAST_GATEWAY_NOT_SUPPORTED.
//
// Example:
//
//	client.Tenants.Create(ctx, "fab1", "t1", nil, nil, nil,
//	    ones_gfx.WithAnycastGateway(true, "00:11:22:33:44:55"))
func WithAnycastGateway(enabled bool, macAddress string) CallOption {
	return func(cfg *callConfig) {
		cfg.anycastGateway = &AnycastGateway{Enabled: enabled, MACAddress: macAddress}
	}
}

// ApplyCallOptions resolves CallOption values into a concrete config.
// Intended for use by resource methods that need timeout/webhook values.
func ApplyCallOptions(opts ...CallOption) (timeout *time.Duration, webhookURL string, webhookEvents []string) {
	cfg := defaultCallConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg.timeout, cfg.webhookURL, cfg.webhookEvents
}

// ResolvedCall is the full per-call configuration after resolving CallOptions.
// Resource methods use it to build transport ReqOpts and webhook body fields.
type ResolvedCall struct {
	Timeout        *time.Duration
	WebhookURL     string
	WebhookEvents  []string
	IdempotencyKey string
	RequestOrigin  string
	AnycastGateway *AnycastGateway
}

// ReqOpts converts the resolved call config into transport-level ReqOpts.
func (rc ResolvedCall) ReqOpts() *ReqOpts {
	return &ReqOpts{
		Timeout:        rc.Timeout,
		IdempotencyKey: rc.IdempotencyKey,
		RequestOrigin:  rc.RequestOrigin,
	}
}

// ResolveCallOptions resolves CallOptions into the full per-call config.
func ResolveCallOptions(opts ...CallOption) ResolvedCall {
	cfg := defaultCallConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return ResolvedCall{
		Timeout:        cfg.timeout,
		WebhookURL:     cfg.webhookURL,
		WebhookEvents:  cfg.webhookEvents,
		IdempotencyKey: cfg.idempotencyKey,
		RequestOrigin:  cfg.requestOrigin,
		AnycastGateway: cfg.anycastGateway,
	}
}

// WithTimeout overrides the client-level timeout for a specific call.
// Use this for long-running operations like GPU allocation (10-15 minutes).
//
// Example:
//
//	client.Tenants.AllocateGPUs(ctx, fabricName, tenantName, servers,
//	    ones_gfx.WithTimeout(15*time.Minute))
func WithTimeout(d time.Duration) CallOption {
	return func(cfg *callConfig) {
		cfg.timeout = &d
	}
}

// WithWebhook configures webhook delivery for async operations.
// Only valid when used with async methods (CreateAsync, DeleteAsync, etc.).
//
// The ONES server will POST the operation result to webhookURL when the
// job completes. Events should be a list of event types to subscribe to,
// e.g., []string{"tenant.create", "tenant.delete"}.
//
// Example:
//
//	op, err := client.Tenants.CreateAsync(ctx, fabricName, tenantName, nil, nil, nil,
//	    ones_gfx.WithWebhook("http://receiver:8000/hook", []string{"tenant.create"}))
func WithWebhook(webhookURL string, events []string) CallOption {
	return func(cfg *callConfig) {
		cfg.webhookURL = webhookURL
		cfg.webhookEvents = events
	}
}

// ClientOption configures optional parameters for the Client.
type ClientOption func(*clientConfig)

// clientConfig holds client-wide configuration.
type clientConfig struct {
	timeout   time.Duration
	tlsConfig *tls.Config
	verifyTLS bool

	// syncPatchFallback re-polls tenant status when a synchronous PATCH times
	// out. Off by default: see WithSyncPatchTimeoutFallback.
	syncPatchFallback SyncPatchFallback
}

// SyncPatchFallback configures what a synchronous PATCH does when it times
// out.
//
// A timeout on a synchronous write is ambiguous: Fabric Manager may have
// finished, may still be applying the change to devices, or may never have
// received it. The transport cannot tell, so by default the timeout surfaces
// as a TransportError and the caller decides.
//
// With Enabled set, the SDK instead treats the operation as possibly still
// running and polls the tenant until its config_status reaches a terminal
// state, which is what callers were otherwise each writing by hand.
type SyncPatchFallback struct {
	// Enabled turns the fallback on. Default false, preserving the strict
	// behaviour of surfacing the timeout.
	Enabled bool

	// PollInterval is the gap between status reads. Default 5s.
	PollInterval time.Duration

	// MaxWait bounds the polling, measured from the moment the PATCH timed
	// out. Default 2m. Reaching it returns a SyncPatchTimeoutError rather
	// than a success or a failure, because the outcome is still unknown.
	MaxWait time.Duration
}

// withDefaults fills in the unset knobs.
func (f SyncPatchFallback) withDefaults() SyncPatchFallback {
	if f.PollInterval <= 0 {
		f.PollInterval = 5 * time.Second
	}
	if f.MaxWait <= 0 {
		f.MaxWait = 2 * time.Minute
	}
	return f
}

// defaultClientConfig returns sensible defaults.
func defaultClientConfig() clientConfig {
	return clientConfig{
		timeout:   30 * time.Second,
		verifyTLS: true,
	}
}

// WithSyncPatchTimeoutFallback enables re-polling tenant status when a
// synchronous PATCH times out.
//
// Pass zero for pollInterval or maxWait to take the defaults (5s and 2m).
//
//	client, _ := ones.InitializeWithCreds(url, user, pass,
//	    ones.WithSyncPatchTimeoutFallback(true, 5*time.Second, 2*time.Minute))
//
// Leaving this off keeps the existing behaviour exactly: a timeout surfaces as
// a TransportError and nothing is polled.
func WithSyncPatchTimeoutFallback(enabled bool, pollInterval, maxWait time.Duration) ClientOption {
	return func(c *clientConfig) {
		c.syncPatchFallback = SyncPatchFallback{
			Enabled:      enabled,
			PollInterval: pollInterval,
			MaxWait:      maxWait,
		}.withDefaults()
	}
}

// WithClientTimeout sets the default HTTP timeout for all requests.
// Individual calls can override with WithTimeout().
func WithClientTimeout(d time.Duration) ClientOption {
	return func(cfg *clientConfig) {
		cfg.timeout = d
	}
}

// WithTLSConfig provides a custom TLS configuration.
// Use this to specify a CA bundle for self-signed certificates.
func WithTLSConfig(tlsConf *tls.Config) ClientOption {
	return func(cfg *clientConfig) {
		cfg.tlsConfig = tlsConf
	}
}

// WithTLSVerify controls whether TLS certificates are verified.
// Set to false only in dev/lab environments with self-signed certs.
// Production deployments should always verify TLS.
func WithTLSVerify(verify bool) ClientOption {
	return func(cfg *clientConfig) {
		cfg.verifyTLS = verify
	}
}

// ResolvedClient is the client-wide configuration after applying ClientOptions.
// It lets callers outside this package read a resolved value — the auth client
// needs VerifyTLS because it builds its own http.Client separately from the
// transport.
type ResolvedClient struct {
	Timeout   time.Duration
	TLSConfig *tls.Config
	VerifyTLS bool
}

// ResolveClientOptions applies opts over the defaults and returns the result,
// mirroring ResolveCallOptions.
func ResolveClientOptions(opts ...ClientOption) ResolvedClient {
	cfg := defaultClientConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return ResolvedClient{
		Timeout:   cfg.timeout,
		TLSConfig: cfg.tlsConfig,
		VerifyTLS: cfg.verifyTLS,
	}
}

// AuthOption configures optional parameters for authentication.
type AuthOption func(*authConfig)

// authConfig holds auth-specific configuration.
type authConfig struct {
	onTokenRefresh         func(accessToken, refreshToken string, expiresIn int)
	verifyTLS              bool
	proactiveRefreshBuffer time.Duration
	refreshTimeout         time.Duration
}

// defaultAuthConfig returns sensible defaults for JWT auth.
func defaultAuthConfig() authConfig {
	return authConfig{
		verifyTLS:              true,
		proactiveRefreshBuffer: 10 * time.Second,
		refreshTimeout:         10 * time.Second,
	}
}

// WithTokenRefreshCallback registers a callback invoked after every successful
// token refresh. Use this to persist rotated tokens to a secret store.
//
// The callback receives (accessToken, refreshToken, expiresIn) where expiresIn
// is the number of seconds until the new access token expires.
func WithTokenRefreshCallback(cb func(string, string, int)) AuthOption {
	return func(cfg *authConfig) {
		cfg.onTokenRefresh = cb
	}
}

// WithAuthTLSVerify controls TLS verification for the refresh endpoint.
// Set to false only for dev/lab with self-signed certs.
func WithAuthTLSVerify(verify bool) AuthOption {
	return func(cfg *authConfig) {
		cfg.verifyTLS = verify
	}
}

// WithProactiveRefreshBuffer sets how many seconds before expiry the SDK
// should proactively refresh the access token. Default is 10 seconds.
func WithProactiveRefreshBuffer(d time.Duration) AuthOption {
	return func(cfg *authConfig) {
		cfg.proactiveRefreshBuffer = d
	}
}

// WithRefreshTimeout sets the HTTP timeout for refresh requests. Default 10s.
func WithRefreshTimeout(d time.Duration) AuthOption {
	return func(cfg *authConfig) {
		cfg.refreshTimeout = d
	}
}
