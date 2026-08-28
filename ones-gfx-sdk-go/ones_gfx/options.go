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
}

// defaultClientConfig returns sensible defaults.
func defaultClientConfig() clientConfig {
	return clientConfig{
		timeout:   30 * time.Second,
		verifyTLS: true,
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
