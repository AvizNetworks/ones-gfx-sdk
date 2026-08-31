package ones

// client.go mirrors the Python SDK's ones_gfx/client.py Client class.
//
// A Client is a configured connection to one ONES instance. It holds the base
// URL, the auth token, and (optionally) the credentials used to obtain it —
// all per instance, nothing global and nothing on disk. Pass it explicitly into
// the package-level API functions:
//
//	client, err := ones.InitializeWithCreds("https://host:3002", "admin", "secret")
//	fabrics, err := ones.GetAllFabrics(ctx, client)

import (
	"errors"
	"strings"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
)

// DefaultTimeout is the per-request timeout applied by the initialisers.
// Generous because GPU allocate/deallocate operations run for minutes.
const DefaultTimeout = 20 * time.Minute

// ErrNotAuthenticated is returned by API calls made without a token.
var ErrNotAuthenticated = ones_gfx.ErrNotAuthenticated

// ErrNilClient is returned when an API function is called with a nil client.
var ErrNilClient = errors.New("client is nil: build one with ones.InitializeWithCreds or ones.InitializeWithToken")

// Client is a configured connection to one ONES instance.
type Client struct {
	baseUrl  string
	username string
	password string

	auth      *ones_gfx.TokenAuth
	transport *ones_gfx.Transport

	Tenants     *resources.TenantsResource
	GPU         *resources.GPUResource
	Fabrics     *resources.FabricsResource
	Fabricsims  *resources.FabricsimsResource
	Inventory   *resources.InventoryResource
	Devices     *resources.DevicesResource
	ConfigMgmt  *resources.ConfigMgmtResource
	Bootstrap   *resources.BootstrapResource
	RMA         *resources.RMAResource
	Files       *resources.FilesResource
	System      *resources.SystemResource
	NetOps      *resources.NetOpsResource
	HostTenants *resources.HostTenantsResource
	NMXC        *resources.NMXCResource
	Intents     *resources.IntentsResource
	VPCPeering  *resources.VPCPeeringResource
	Operations  *resources.OperationsResource
}

// NewClient builds a Client for the given root base URL.
//
// TLS certificate verification is disabled — ONES deployments normally use
// self-signed certificates. Do not point this at an untrusted host: a
// man-in-the-middle can read and alter the traffic, including the credentials
// sent to /api/user/login.
func NewClient(baseUrl string) (*Client, error) {
	if baseUrl == "" {
		return nil, errors.New("baseUrl is required")
	}
	baseUrl = strings.TrimRight(baseUrl, "/")
	auth := ones_gfx.NewTokenAuth(baseUrl, false, DefaultTimeout)
	transport := ones_gfx.NewTransportWithOptions(
		ones_gfx.FMBaseURL(baseUrl),
		auth,
		ones_gfx.WithTLSVerify(false),
		ones_gfx.WithClientTimeout(DefaultTimeout),
	)
	return &Client{
		baseUrl:     baseUrl,
		auth:        auth,
		transport:   transport,
		Tenants:     resources.NewTenantsResource(transport),
		GPU:         resources.NewGPUResource(transport),
		Fabrics:     resources.NewFabricsResource(transport),
		Fabricsims:  resources.NewFabricsimsResource(transport),
		Inventory:   resources.NewInventoryResource(transport),
		Devices:     resources.NewDevicesResource(transport),
		ConfigMgmt:  resources.NewConfigMgmtResource(transport),
		Bootstrap:   resources.NewBootstrapResource(transport),
		RMA:         resources.NewRMAResource(transport),
		Files:       resources.NewFilesResource(transport),
		System:      resources.NewSystemResource(transport),
		NetOps:      resources.NewNetOpsResource(transport),
		HostTenants: resources.NewHostTenantsResource(transport),
		NMXC:        resources.NewNMXCResource(transport),
		Intents:     resources.NewIntentsResource(transport),
		VPCPeering:  resources.NewVPCPeeringResource(transport),
		Operations:  resources.NewOperationsResource(transport),
	}, nil
}

// NewClientWithAuth builds a Client around an arbitrary AuthProvider — use it
// for the JWTAuth (access/refresh pair, `Authorization: Bearer`) flow, or to
// supply custom TLS/timeout options.
//
// Login/RefreshAuth/Logout and the token getters/setters only work when auth is
// a *ones_gfx.TokenAuth; with any other provider they return ErrNoTokenAuth.
func NewClientWithAuth(baseUrl string, auth ones_gfx.AuthProvider, opts ...ones_gfx.ClientOption) *Client {
	root := strings.TrimRight(baseUrl, "/")
	transport := ones_gfx.NewTransportWithOptions(ones_gfx.FMBaseURL(root), auth, opts...)
	tokenAuth, _ := auth.(*ones_gfx.TokenAuth)
	c := &Client{
		baseUrl:     root,
		auth:        tokenAuth,
		transport:   transport,
		Tenants:     resources.NewTenantsResource(transport),
		GPU:         resources.NewGPUResource(transport),
		Fabrics:     resources.NewFabricsResource(transport),
		Fabricsims:  resources.NewFabricsimsResource(transport),
		Inventory:   resources.NewInventoryResource(transport),
		Devices:     resources.NewDevicesResource(transport),
		ConfigMgmt:  resources.NewConfigMgmtResource(transport),
		Bootstrap:   resources.NewBootstrapResource(transport),
		RMA:         resources.NewRMAResource(transport),
		Files:       resources.NewFilesResource(transport),
		System:      resources.NewSystemResource(transport),
		NetOps:      resources.NewNetOpsResource(transport),
		HostTenants: resources.NewHostTenantsResource(transport),
		NMXC:        resources.NewNMXCResource(transport),
		Intents:     resources.NewIntentsResource(transport),
		VPCPeering:  resources.NewVPCPeeringResource(transport),
		Operations:  resources.NewOperationsResource(transport),
	}
	return c
}

// ErrNoTokenAuth is returned by the session methods on a client built with a
// non-TokenAuth provider (e.g. JWTAuth), which has no single-token login flow.
var ErrNoTokenAuth = errors.New("client has no TokenAuth: build it with InitializeWithCreds/InitializeWithToken")

// ---------------------------------------------------------------------------
// Initialisers
// ---------------------------------------------------------------------------

// InitializeWithCreds builds a client, logs in, and stores the returned token.
// A returned client is always usable.
func InitializeWithCreds(baseUrl, username, password string) (*Client, error) {
	c, err := NewClient(baseUrl)
	if err != nil {
		return nil, err
	}
	c.username = username
	c.password = password
	if _, err := c.Login("", ""); err != nil {
		return nil, err
	}
	return c, nil
}

// InitializeWithToken builds a client from an existing token. The token is sent
// as the `authorization` header on every call.
//
// Returns an error when token is empty.
func InitializeWithToken(baseUrl, token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("token is required when initializing with a token")
	}
	c, err := NewClient(baseUrl)
	if err != nil {
		return nil, err
	}
	c.auth.SetToken(token)
	return c, nil
}

// ---------------------------------------------------------------------------
// Getters / setters
// ---------------------------------------------------------------------------

// GetBaseUrl returns the root base URL.
func (c *Client) GetBaseUrl() string { return c.baseUrl }

// SetBaseUrl points this client at a different host.
func (c *Client) SetBaseUrl(baseUrl string) error {
	if baseUrl == "" {
		return errors.New("baseUrl is required")
	}
	baseUrl = strings.TrimRight(baseUrl, "/")
	c.baseUrl = baseUrl
	if c.auth != nil {
		c.auth.SetBaseURL(baseUrl)
	}
	c.transport.SetBaseURL(ones_gfx.FMBaseURL(baseUrl))
	return nil
}

// GetAuthToken returns the stored token ("" when not logged in).
func (c *Client) GetAuthToken() string {
	if c.auth == nil {
		return ""
	}
	return c.auth.Token()
}

// SetAuthToken replaces the stored token without validation.
func (c *Client) SetAuthToken(token string) {
	if c.auth != nil {
		c.auth.SetToken(token)
	}
}

// GetUsername returns the stored username.
func (c *Client) GetUsername() string { return c.username }

// SetUsername stores the username used by Login.
func (c *Client) SetUsername(username string) { c.username = username }

// GetPassword returns the stored password.
func (c *Client) GetPassword() string { return c.password }

// SetPassword stores the password used by Login.
func (c *Client) SetPassword(password string) { c.password = password }

// ---------------------------------------------------------------------------
// Token management
// ---------------------------------------------------------------------------

// UpdateToken replaces the stored token, rejecting an empty value. Prefer this
// over SetAuthToken when rotating a token.
func (c *Client) UpdateToken(token string) error {
	if c.auth == nil {
		return ErrNoTokenAuth
	}
	return c.auth.UpdateToken(token)
}

// IsAuthenticated reports whether a token is held.
func (c *Client) IsAuthenticated() bool {
	return c.auth != nil && c.auth.Authenticated()
}

// ---------------------------------------------------------------------------
// Session endpoints (/api/user/...)
// ---------------------------------------------------------------------------

// Login authenticates and stores the returned token, falling back to the
// credentials already on the client when username/password are empty.
func (c *Client) Login(username, password string) (*ones_gfx.AuthResponse, error) {
	user, pwd := username, password
	if user == "" {
		user = c.username
	}
	if pwd == "" {
		pwd = c.password
	}
	if c.auth == nil {
		return nil, ErrNoTokenAuth
	}
	if user == "" || pwd == "" {
		return nil, errors.New("username and password are required to log in")
	}
	// Remember them so a later Login/RefreshAuth works without re-supplying.
	c.username, c.password = user, pwd
	return c.auth.Login(user, pwd)
}

// RefreshAuth exchanges the current token for a new one and stores it.
// Returns ErrNotAuthenticated when there is no current token to exchange.
func (c *Client) RefreshAuth() (*ones_gfx.AuthResponse, error) {
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

// ---------------------------------------------------------------------------
// Misc
// ---------------------------------------------------------------------------

// GetTransport returns the underlying transport for raw calls not covered by
// the generated API functions.
func (c *Client) GetTransport() *ones_gfx.Transport { return c.transport }

// Close releases HTTP resources. Safe to call multiple times.
func (c *Client) Close() error { return c.transport.Close() }
