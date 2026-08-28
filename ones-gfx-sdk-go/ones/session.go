package ones

// session.go holds the package-level session: Init stores the base URL, Login
// stores the token, and every API function in api.go uses them implicitly.
// Mirrors the module-global model of the Python SDK's ones_gfx/client.py.

import (
	"errors"
	"sync"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

var (
	globalMu     sync.RWMutex
	globalClient *Client
)

// DefaultTimeout is the per-request timeout Init applies. Generous because GPU
// allocate/deallocate operations routinely run for minutes.
const DefaultTimeout = 20 * time.Minute

// ErrNotConfigured is returned by the package-level API functions when Init has
// not been called.
var ErrNotConfigured = errors.New("client not initialised: call ones.Init(baseURL) first")

// ErrNoTokenAuth is returned by Login/Refresh/Logout on a client that was built
// with a non-TokenAuth AuthProvider (e.g. JWTAuth), which has no login flow.
var ErrNoTokenAuth = errors.New("client has no TokenAuth: build it with ones.Init to use Login/Refresh/Logout")

// Init stores the base URL for every subsequent call. This is the only setup
// step; it returns nothing and performs no I/O.
//
// Pass the *root* base URL (e.g. https://host:3002): the auth endpoints live
// under /api/user/ and the Fabric Manager resources under /api/fm/, and Init
// wires both. Any token already in ones_gfx.SecretsFile is loaded, so a previous
// session's login is reused without calling Login again.
//
//	ones.Init("https://10.4.5.76:8089")
//	if _, err := ones.Login("user", "pass"); err != nil { ... }
//	fabrics, err := ones.GetAllFabrics(ctx)
//
// TLS certificate verification is disabled — ONES deployments normally use
// self-signed certificates. Do not point this at an untrusted host: a
// man-in-the-middle can read and alter the traffic, including the credentials
// sent to /api/user/login. Build a client with NewClient and
// ones_gfx.WithTLSConfig if you need certificate pinning or a CA bundle.
func Init(baseURL string) {
	auth := ones_gfx.NewTokenAuth(baseURL, false, DefaultTimeout)
	c := NewClient(
		ones_gfx.FMBaseURL(baseURL),
		auth,
		ones_gfx.WithTLSVerify(false),
		ones_gfx.WithClientTimeout(DefaultTimeout),
	)

	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()
}

// --- package-level session operations ---------------------------------------

// Login authenticates and stores the token, which is then sent on every
// subsequent API call and persisted to ones_gfx.SecretsFile.
func Login(username, password string) (*ones_gfx.AuthResponse, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}
	return c.Login(username, password)
}

// Refresh exchanges the stored token for a new one and stores that.
func Refresh() (*ones_gfx.AuthResponse, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}
	return c.Refresh()
}

// Logout invalidates the session server-side and clears the stored token.
func Logout() (*ones_gfx.AuthResponse, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}
	return c.Logout()
}

// Token returns the stored auth token ("" when not logged in or not initialised).
func Token() string {
	c, err := GetClient()
	if err != nil {
		return ""
	}
	return c.Token()
}

// Close releases the shared client's HTTP resources.
func Close() error {
	c, err := GetClient()
	if err != nil {
		return nil
	}
	return c.Close()
}

// GetClient returns the shared client, or ErrNotConfigured if Init was never
// called. Most callers do not need this — use the package-level functions.
func GetClient() (*Client, error) {
	globalMu.RLock()
	defer globalMu.RUnlock()
	if globalClient == nil {
		return nil, ErrNotConfigured
	}
	return globalClient, nil
}

// SetClient installs an explicit shared client, for tests or for callers that
// build the client themselves via NewClient.
func SetClient(c *Client) {
	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()
}

// requireClient resolves the shared client for the package-level API functions.
//
// It deliberately does NOT check for a token: calls are attempted even when not
// logged in, and the server decides. An unauthenticated request comes back as an
// *ones_gfx.AuthenticationError (HTTP 401/403) from the transport rather than a
// client-side error.
func requireClient() (*Client, error) {
	return GetClient()
}

// Login authenticates and persists the token. Subsequent API calls are
// authorised automatically.
func (c *Client) Login(username, password string) (*ones_gfx.AuthResponse, error) {
	if c.auth == nil {
		return nil, ErrNoTokenAuth
	}
	return c.auth.Login(username, password)
}

// Refresh exchanges the current token for a new one and persists it.
func (c *Client) Refresh() (*ones_gfx.AuthResponse, error) {
	if c.auth == nil {
		return nil, ErrNoTokenAuth
	}
	return c.auth.RefreshWithResponse()
}

// Logout invalidates the session server-side and clears the stored token.
func (c *Client) Logout() (*ones_gfx.AuthResponse, error) {
	if c.auth == nil {
		return nil, ErrNoTokenAuth
	}
	return c.auth.Logout()
}

// Token returns the currently held auth token ("" when not logged in).
func (c *Client) Token() string {
	if c.auth == nil {
		return ""
	}
	return c.auth.Token()
}

// Authenticated reports whether a token is currently held.
func (c *Client) Authenticated() bool { return ones_gfx.Authenticated() }
