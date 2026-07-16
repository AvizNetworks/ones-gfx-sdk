# ONES GFX SDK — Go

Go client library for **AVIZ ONES Spectrum-X** tenant management API (v4.2).

## Status

**Core Library:** ✅ Production-ready (v1.0.0)  
**CLI Binary:** 🚧 In development (coming in v1.1.0)

The core library is fully functional and can be imported into Go applications now. The CLI wrapper (`ones-gfx-sdk-mod`) is under development.

---

## Requirements

- **Go 1.19** or newer
- **Zero external dependencies** (stdlib only)

---

## Installation

### As a Library (Recommended)

```bash
go get github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go
```

### Build from Source

```bash
# Clone the monorepo
git clone https://github.com/aviznetworks/ones-gfx-sdk.git
cd ones-gfx-sdk/ones-gfx-sdk-go

# Verify the build
go build ./ones_gfx
go build ./ones_gfx/resources
go build ./sdk
```

---

## Quick Start — Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    ones_gfx "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/sdk"
)

func main() {
    // 1. Set up JWT authentication
    auth, err := ones_gfx.NewJWTAuth(
        "your-access-token",
        "your-refresh-token",
        "https://10.4.5.76:8089/refresh",
        ones_gfx.WithTokenRefreshCallback(func(access, refresh string, expiresIn int) {
            fmt.Printf("Tokens refreshed, expires in %d seconds\n", expiresIn)
        }),
        ones_gfx.WithAuthTLSVerify(false), // Only for dev/lab with self-signed certs
    )
    if err != nil {
        log.Fatal(err)
    }

    // 2. Create the client
    client := sdk.NewClient(
        "https://10.4.5.76:8089",
        auth,
        ones_gfx.WithClientTimeout(20*time.Minute), // Default 30s, increase for long ops
        ones_gfx.WithTLSVerify(false),               // Only for dev/lab
    )
    defer client.Close()

    ctx := context.Background()

    // 3. List fabrics
    fabrics, err := client.Fabrics.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, fabric := range fabrics {
        fmt.Printf("Fabric: %s (Storage VPC: %s)\n",
            fabric.FabricName, fabric.DefaultStorageName)
    }

    // 4. Create a tenant (synchronous mode)
    tenant, err := client.Tenants.Create(ctx, "sdk", resources.CreateTenantRequest{
        Name:           "demo-tenant",
        Description:    "Created via Go SDK",
        MaxGPUsAllowed: 8,
        Shared:         false,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created tenant ID: %d, VLAN: %d\n", tenant.ID, *tenant.VLANID)

    // 5. Allocate GPUs with timeout override (15 minutes)
    servers := resources.ServerSpecsFromNames([]string{"hgx-su00-h00"})
    err = client.Tenants.AllocateGPUs(ctx, "sdk", "demo-tenant", servers,
        ones_gfx.WithTimeout(15*time.Minute), // Override client timeout for this call
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("GPUs allocated successfully")

    //Assign specific ports — UFM / NMXC fabrics only (always synchronous mode)
    req := resources.GpuPortAssignmentRequest{
        ServerNames: []string{"su00-rack00-node00"},
        GPUIDs:      []int{1, 2, 3},
    }
    err = client.Tenants.AssignPorts(ctx, "sdk", "demo-tenant", req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Ports assigned successfully")

    // 6. Get tenant details
    updatedTenant, err := client.Tenants.Get(ctx, "sdk", "demo-tenant")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Allocated servers: %v\n", updatedTenant.AllotedServers())

    // 7. Deallocate GPUs
    err = client.Tenants.DeallocateGPUs(ctx, "sdk", "demo-tenant", servers)
    if err != nil {
        log.Fatal(err)
    }

    //Unassign ports
    err = client.Tenants.UnassignPorts(ctx, "sdk", "demo-tenant", req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Ports unassigned successfully")

    // 8. Delete tenant
    err = client.Tenants.Delete(ctx, "sdk", "demo-tenant")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tenant deleted successfully")
}
```

## API Reference

### Client Construction

```go
// Minimal
client := sdk.NewClient(baseURL, auth)

// With options
client := sdk.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute),
    ones_gfx.WithTLSVerify(false),
    ones_gfx.WithTLSConfig(customTLSConfig),
)
```

### Authentication

```go
auth, err := ones_gfx.NewJWTAuth(accessToken, refreshToken, refreshURL,
    ones_gfx.WithTokenRefreshCallback(func(access, refresh string, expiresIn int) {
        // Persist tokens
    }),
    ones_gfx.WithAuthTLSVerify(false),
    ones_gfx.WithProactiveRefreshBuffer(10*time.Second),
    ones_gfx.WithRefreshTimeout(10*time.Second),
)
if err != nil {
    log.Fatal(err)
}
```

### Fabrics

```go
// List
fabrics, err := client.Fabrics.List(ctx)

// Inventory sync — UFM enabled fabrics only
err := client.Fabrics.InventorySync(ctx, fabricName)
```

### Tenants

```go
// List
tenants, err := client.Tenants.List(ctx, fabricName)

// Get
tenant, err := client.Tenants.Get(ctx, fabricName, tenantName)

// Available servers
servers, err := client.Tenants.AvailableServers(ctx, fabricName)

// Convert hostnames to ServerSpec entries
serverSpecs := resources.ServerSpecsFromNames([]string{"hgx-su00-h00"})

// Request body for Specific Gpus UFM / NMXC fabrics (But for UFM serverNames is {"hgx-su00-h00"} )
req := resources.GpuPortAssignmentRequest{
    ServerNames: []string{"su00-rack00-node00"},
    GPUIDs:      []int{1, 2, 3},
}

// Create (sync)
tenant, err := client.Tenants.Create(ctx, fabricName, resources.CreateTenantRequest{
    Name: "tenant1", Description: "...", MaxGPUsAllowed: 8, Shared: false,
})

// Create (async)
op, err := client.Tenants.CreateAsync(ctx, fabricName, resources.CreateTenantRequest{
    Name: "tenant1", Description: "...", MaxGPUsAllowed: 8, Shared: false,
})

// Delete (sync)
err := client.Tenants.Delete(ctx, fabricName, tenantName)

// Delete (async)
op, err := client.Tenants.DeleteAsync(ctx, fabricName, tenantName)

// Allocate GPUs (sync)
err := client.Tenants.AllocateGPUs(ctx, fabricName, tenantName,
    serverSpecs,
    ones_gfx.WithTimeout(15*time.Minute),
)

// Allocate GPUs (async)
op, err := client.Tenants.AllocateGPUsAsync(ctx, fabricName, tenantName, serverSpecs)

// Deallocate GPUs (sync)
err := client.Tenants.DeallocateGPUs(ctx, fabricName, tenantName, serverSpecs)

// Deallocate GPUs (async)
op, err := client.Tenants.DeallocateGPUsAsync(ctx, fabricName, tenantName, serverSpecs)

// Assign specific ports — UFM / NMXC fabrics only (always sync)
err := client.Tenants.AssignPorts(ctx, fabricName, tenantName, req)

// Assign all ports on each server (omit GPUIDs)
err := client.Tenants.AssignPorts(ctx, fabricName, tenantName, resources.GpuPortAssignmentRequest{ServerNames: req.ServerNames})

// Unassign specific ports
err := client.Tenants.UnassignPorts(ctx, fabricName, tenantName, req)

// Unassign all ports on each server (omit GPUIDs)
err := client.Tenants.UnassignPorts(ctx, fabricName, tenantName, resources.GpuPortAssignmentRequest{ServerNames: req.ServerNames})
```

---

## Operation Modes

All mutating methods support two execution modes:

### Synchronous (Default)

Blocks until the server completes the operation. Returns the result directly.

```go
// Returns *ones_gfx.Tenant when done
tenant, err := client.Tenants.Create(ctx, fabricName, req)
```

### Asynchronous (Poll)

Returns immediately with an Operation handle. You poll for completion.

```go
// Returns *ones_gfx.Operation immediately
op, err := client.Tenants.CreateAsync(ctx, fabricName, req)

// Poll until done
for {
    current, _ := client.Operations.Get(ctx, op.ID)
    if current.IsDone() {
        if current.IsSuccess() {
            fmt.Println("Operation succeeded:", current.Result)
        } else {
            fmt.Println("Operation failed:", current.ErrorMessage)
        }
        break
    }
    time.Sleep(5 * time.Second)
}
```

### Asynchronous (Webhook)

Returns immediately. Server POSTs result to your webhook URL.

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, req,
    ones_gfx.WithWebhook("http://receiver:8000/hook", []string{"tenant.create"}))

fmt.Printf("Operation submitted: %s (webhook: %v)\n", 
    op.ID, op.WebhookRegistered)
// No polling needed — result sent to webhook
```

---

## Timeout Override

Long-running operations (GPU allocate/deallocate) can take 10-15 minutes. Override the default timeout:

```go
// Client-level default: 20 minutes
client := sdk.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute))

// Per-call override: 15 minutes for this specific allocation
err := client.Tenants.AllocateGPUs(ctx, fabricName, tenantName, servers,
    ones_gfx.WithTimeout(15*time.Minute))
```

---

## Error Handling

All errors implement the `ones_gfx.ONESError` interface. Use `errors.As` to check types:

```go
import "errors"

tenant, err := client.Tenants.Get(ctx, fabricName, "nonexistent")
if err != nil {
    // Check for specific error types
    var notFound *ones_gfx.NotFoundError
    if errors.As(err, &notFound) {
        fmt.Printf("Tenant not found (HTTP %d)\n", notFound.StatusCode)
        return
    }
    
    var conflict *ones_gfx.ConflictError
    if errors.As(err, &conflict) {
        fmt.Println("Conflict:", conflict.Message)
        return
    }
    
    // Generic ONESError check
    var onesErr ones_gfx.ONESError
    if errors.As(err, &onesErr) {
        fmt.Println("SDK error:", err)
        return
    }
    
    // Non-SDK error
    log.Fatal(err)
}
```

**Error Types:**
- `TransportError` — Network/TLS/timeout
- `AuthenticationError` — HTTP 401/403 or refresh failure
- `BadRequestError` — HTTP 400
- `NotFoundError` — HTTP 404
- `ConflictError` — HTTP 409
- `ServerError` — HTTP 5xx
- `OperationFailedError` — Async job ended with status FAILURE

---

## Building the SDK

### Verify the Build

```bash
cd ones-gfx-sdk/ones-gfx-sdk-go

# Build core library
go build ./ones_gfx

# Build resources
go build ./ones_gfx/resources

# Build client package
go build ./sdk

# Run go vet (static analysis)
go vet ./...

# Format code
go fmt ./...
```

### Check for Issues

```bash
# Static analysis
go vet ./...

# Unused code detection (requires staticcheck)
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

---

## Testing the SDK

### Manual Testing (Requires ONES Instance)

Create a test file `test_sdk.go`:

```go
package main

import (
    "context"
    "log"
    "time"

    ones_gfx "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/sdk"
)

func main() {
    auth, err := ones_gfx.NewJWTAuth(
        "YOUR_ACCESS_TOKEN",
        "YOUR_REFRESH_TOKEN",
        "https://YOUR_ONES_IP:8089/refresh",
        ones_gfx.WithAuthTLSVerify(false),
    )
    if err != nil {
        log.Fatal(err)
    }

    client := sdk.NewClient(
        "https://YOUR_ONES_IP:8089",
        auth,
        ones_gfx.WithTLSVerify(false),
        ones_gfx.WithClientTimeout(20*time.Minute),
    )
    defer client.Close()

    ctx := context.Background()

    // Test 1: List fabrics
    log.Println("Test 1: Listing fabrics...")
    fabrics, err := client.Fabrics.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found %d fabric(s)\n", len(fabrics))

    if len(fabrics) == 0 {
        log.Fatal("No fabrics found - cannot proceed with tests")
    }

    fabricName := fabrics[0].FabricName
    log.Printf("Using fabric: %s\n", fabricName)

    // Test 2: List tenants
    log.Println("Test 2: Listing tenants...")
    tenants, err := client.Tenants.List(ctx, fabricName)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found %d tenant(s)\n", len(tenants))

    // Test 3: Available servers
    log.Println("Test 3: Checking available servers...")
    servers, err := client.Tenants.AvailableServers(ctx, fabricName)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Available servers: %v\n", servers)

    log.Println("All tests passed!")
}
```

Run it:

```bash
go run test_sdk.go
```

### Expected Output

```
Test 1: Listing fabrics...
Found 1 fabric(s)
Using fabric: sdk
Test 2: Listing tenants...
Found 3 tenant(s)
Test 3: Checking available servers...
Available servers: [hgx-su00-h00 hgx-su00-h01 hgx-su00-h02]
All tests passed!
```

---

## File Structure

```
ones-gfx-sdk/ones-gfx-sdk-go/
├── go.mod                      # Module definition
├── README.md                   # This file
├── IMPLEMENTATION_STATUS.md    # Build status
├── ones_gfx/                   # Core library (importable)
│   ├── auth.go                 # JWT authentication
│   ├── transport.go            # HTTP client
│   ├── enums.go                # Operation modes, statuses
│   ├── errors.go               # Typed errors
│   ├── models.go               # Fabric, Tenant, Operation structs
│   ├── options.go              # Functional options
│   └── resources/
│       ├── fabrics.go
│       ├── tenants.go
│       ├── operations.go
│       └── peering.go
├── sdk/                        # Client entry point
│   └── client.go               # Client struct
└── cmd/                        # CLI (coming in v1.1.0)
    └── ones-gfx-sdk-mod/
        └── main.go (WIP)
```

---

## Troubleshooting

### Issue: `certificate signed by unknown authority`

**Solution:** Disable TLS verification for dev/lab (self-signed certs):

```go
client := sdk.NewClient(baseURL, auth,
    ones_gfx.WithTLSVerify(false))
```

For production, provide a CA bundle:

```go
import "crypto/tls"
import "crypto/x509"

caCert, _ := os.ReadFile("/etc/ssl/certs/ones-ca.pem")
caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

tlsConfig := &tls.Config{
    RootCAs: caCertPool,
}

client := sdk.NewClient(baseURL, auth,
    ones_gfx.WithTLSConfig(tlsConfig))
```

### Issue: `context deadline exceeded`

**Solution:** Increase timeout for long operations:

```go
// Client-level
client := sdk.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute))

// Per-call
err := client.Tenants.AllocateGPUs(ctx, ...,
    ones_gfx.WithTimeout(15*time.Minute))
```

### Issue: `authentication error: token refresh failed`

**Solution:** Verify refresh token is valid and refresh URL is correct. Both tokens rotate on every refresh.

---

## Known Limitations (v1.0.0)

- **No login endpoint** — Partners supply tokens obtained out-of-band. Login support planned for v1.1.
- **CLI not yet complete** — Core library is fully functional. CLI wrapper coming in v1.1.


---


---

## Support

- **Documentation:** This README
- **Issues:** [GitHub Issues](https://github.com/aviznetworks/ones-gfx-sdk/issues)
- **API Reference:** [api_reference.txt](./api_reference.txt) (language-agnostic curl examples)
