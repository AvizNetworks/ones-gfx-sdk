package ones

import (
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
)

// Client is the main entry point for the SDK. It holds the transport and
// exposes one resource handle per API domain.
type Client struct {
	transport *ones_gfx.Transport

	// auth is set when the client was built by Configure (or NewClient with a
	// *ones_gfx.TokenAuth). It backs Login/Refresh/Logout; nil for clients
	// constructed with a different AuthProvider such as JWTAuth.
	auth *ones_gfx.TokenAuth

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

// NewClient constructs a Client with the given base URL and authentication.
//
// Example:
//
//	auth, err := ones_gfx.NewJWTAuth(accessToken, refreshToken, refreshURL)
//	if err != nil {
//	    panic(err)
//	}
//	client := sdk.NewClient("https://10.4.5.76:8089", auth,
//	    ones_gfx.WithClientTimeout(20*time.Minute),
//	    ones_gfx.WithTLSVerify(false))
//	defer client.Close()
func NewClient(baseURL string, auth ones_gfx.AuthProvider, opts ...ones_gfx.ClientOption) *Client {
	transport := ones_gfx.NewTransportWithOptions(baseURL, auth, opts...)
	tokenAuth, _ := auth.(*ones_gfx.TokenAuth)
	return &Client{
		transport:   transport,
		auth:        tokenAuth,
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
}

// GetTransport returns the underlying transport for advanced use cases
// such as making raw API calls not covered by resource methods.
func (c *Client) GetTransport() *ones_gfx.Transport {
	return c.transport
}

// Close releases HTTP resources. Safe to call multiple times.
func (c *Client) Close() error {
	return c.transport.Close()
}
