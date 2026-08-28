# ONES GFX SDK — Go

Go client library for **AVIZ ONES Spectrum-X** tenant management API (v4.2).

## Status

**Core Library:** ✅ Typed surface complete — all 97 Fabric Manager endpoints
exposed as 102 typed resource methods across 17 resource handles.
**CLI Binary:** 🚧 Not started.

Typed request params and typed return values throughout; the only untyped
returns are the 16 documented open-payload leaves (see
[IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md)).

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
go build ./...
go vet ./...
```

---

## Quick Start — recommended (flat API, one import)

Every endpoint is also a package-level function in `ones`, and every type,
constant and option is re-exported there — so **`ones` is the only import you
need**:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

func main() {
    ctx := context.Background()

    // One setup call with the ROOT base URL. This wires the /api/user/ auth
    // endpoints and the /api/fm/ resource base, and loads any persisted token.
    // TLS verification is disabled (ONES uses self-signed certs).
    ones.Init("https://10.4.5.76:8089")
    defer ones.Close()

    // Log in once; the token is stored and reused on later runs.
    if _, err := ones.Login("admin", "secret"); err != nil {
        log.Fatal(err)
    }

    fabrics, err := ones.GetAllFabrics(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range fabrics {
        fmt.Printf("%d %s (type=%s status=%s)\n", f.ID, f.Name, f.Type, f.Status)
    }

    msg, err := ones.CreateFabric(ctx, "gpu-fabric-1", &ones.FabricCreateArgs{
        Type:        ones.Ptr("DNO ASN"),
        Description: ones.Ptr("Primary GPU fabric"),
        Status:      ones.Ptr("draft"),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(msg)
}
```

`ones.Ptr` is a generic helper for the optional (pointer) arguments.

### Session methods

```go
ones.Init(baseURL)              // stores the URL; returns nothing, no I/O
ones.Login(username, password)  // POST /api/user/login,   stores the token
ones.Refresh()                  // POST /api/user/refresh, rotates it
ones.Logout()                   // POST /api/user/logout,  clears it
ones.Token()                    // stored token ("" when logged out)
ones.Authenticated()            // bool
ones.Close()                    // release HTTP resources
```

All state lives in the package — there is no client object to pass around. If
you do need the underlying client (for the resource-handle style below), call
`ones.GetClient()`.

The token is persisted to `secrets.json` (see `ones_gfx.SecretsFile`) and
reloaded by `Init`, so a token from a previous run is reused automatically.
The transport also refreshes reactively on a 401.

### Flat function names

Names flatten `Resource.Method` — e.g. `Fabrics.Create` → `ones.CreateFabric`,
`Fabrics.List` → `ones.GetAllFabrics`, `Tenants.ModifyAllocations` →
`ones.ModifyGpuAllocations`. The full index is in
[IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md).

---

## Conventions

Before the examples, three rules that apply to every resource method:

1. **Path params are positional**; body/query fields are individual typed args.
2. **Optional fields are pointers** — `nil` omits them from the request. This
   mirrors the Python SDK's `| None = None`.
3. **Async-capable endpoints have a separate `...Async` method** returning
   `*ones_gfx.OperationAccepted`.

A tiny helper makes the pointer args readable (Go 1.19 generics):

```go
func ptr[T any](v T) *T { return &v }
```

---

## Quick Start — alternative (resource handles, JWT auth)

The resource handles remain available on the client if you prefer grouping by
domain, or need `JWTAuth` (access/refresh pair, `Authorization: Bearer`) instead
of the single-token login flow. This style needs the `ones_gfx` and `resources`
imports too.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    ones_gfx "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

func ptr[T any](v T) *T { return &v }

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
    client := ones.NewClient(
        "https://10.4.5.76:8089",
        auth,
        ones_gfx.WithClientTimeout(20*time.Minute), // Default 30s, increase for long ops
        ones_gfx.WithTLSVerify(false),              // Only for dev/lab
    )
    defer client.Close()

    ctx := context.Background()

    // 3. List fabrics
    fabrics, err := client.Fabrics.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, fabric := range fabrics {
        fmt.Printf("Fabric: %s (type=%s status=%s)\n",
            fabric.Name, fabric.Type, fabric.Status)
    }

    // 4. Create a tenant (synchronous). Returns the server's message string.
    msg, err := client.Tenants.Create(ctx, "sdk", "demo-tenant",
        ptr("Created via Go SDK"), // description
        ptr(8),                    // maxGpusAllowed
        ptr(false),                // shared
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Create tenant:", msg)

    // 5. Attach GPU servers to the tenant, with a per-call timeout override
    resp, err := client.Tenants.Update(ctx, "sdk", "demo-tenant",
        []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
        ptr(ones_gfx.GpuActionAdd),
        ones_gfx.WithTimeout(15*time.Minute),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Servers attached:", resp.Message)

    // 6. Assign specific GPU ports — UFM / NMXC fabrics only (always sync)
    ports, err := client.GPU.AssignPorts(ctx, "sdk", "demo-tenant",
        ones_gfx.GpuActionAdd,
        []string{"su00-rack00-node00"},
        []int{1, 2, 3},
        nil, // membership
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Ports assigned: success=%v %s\n", ports.Success, ports.Message)

    // 7. Get tenant details (open payload — see the leaves list)
    detail, err := client.Tenants.Get(ctx, "sdk", "demo-tenant")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tenant detail: %v\n", detail)

    // 8. Detach the servers again
    if _, err = client.Tenants.Update(ctx, "sdk", "demo-tenant",
        []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
        ptr(ones_gfx.GpuActionDelete),
    ); err != nil {
        log.Fatal(err)
    }

    // 9. Delete the tenant
    del, err := client.Tenants.Delete(ctx, "sdk", "demo-tenant")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tenant deleted:", del.Message)
}
```

## API Reference

The full resource → method index lives in
[IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md). Highlights below.

### Client Construction

```go
// Minimal
client := ones.NewClient(baseURL, auth)

// With options
client := ones.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute),
    ones_gfx.WithTLSVerify(false),
    ones_gfx.WithTLSConfig(customTLSConfig),
)
```

Resource handles: `Tenants`, `GPU`, `Fabrics`, `Fabricsims`, `Inventory`,
`Devices`, `ConfigMgmt`, `Bootstrap`, `RMA`, `Files`, `System`, `NetOps`,
`HostTenants`, `NMXC`, `Intents`, `VPCPeering`, `Operations`.

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
// List (full entity form)
fabrics, err := client.Fabrics.List(ctx)          // []ones_gfx.FabricItem

// Get one
fabric, err := client.Fabrics.Get(ctx, name)      // *ones_gfx.FabricItem

// Grouped DTO form
dtos, err := client.Fabrics.ListDTOs(ctx)

// Device IPs
ips, err := client.Fabrics.DeviceIPs(ctx, fabricName)
ips, err = client.Fabrics.DevicesByLayer(ctx, "spine")

// Inventory sync — UFM enabled fabrics only (lives on Inventory)
res, err := client.Inventory.Sync(ctx, fabricName)
```

### Tenants

```go
// List / Get (open payloads)
tenants, err := client.Tenants.List(ctx, fabricName)
tenant, err := client.Tenants.Get(ctx, fabricName, tenantName)

// Servers with free GPUs (lives on GPU)
servers, err := client.GPU.AvailableServers(ctx, fabricName)
fmt.Println(servers.AvailableGPUs)

// Create (sync / async)
msg, err := client.Tenants.Create(ctx, fabricName, "tenant1",
    ptr("description"), ptr(8), ptr(false))
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1",
    ptr("description"), ptr(8), ptr(false))

// Delete (sync / async)
resp, err := client.Tenants.Delete(ctx, fabricName, tenantName)
op, err = client.Tenants.DeleteAsync(ctx, fabricName, tenantName)

// Attach / detach whole servers (replaces the old AllocateGPUs/DeallocateGPUs)
resp, err = client.Tenants.Update(ctx, fabricName, tenantName,
    []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
    ptr(ones_gfx.GpuActionAdd),
    ones_gfx.WithTimeout(15*time.Minute),
)
op, err = client.Tenants.UpdateAsync(ctx, fabricName, tenantName,
    []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
    ptr(ones_gfx.GpuActionDelete),
)

// Share a server between tenants
servers := []resources.GpuServerInfo{{ServerName: "hgx-su00-h00", Shared: ptr(true)}}

// Assign / unassign specific ports — UFM / NMXC only, always sync
ports, err := client.GPU.AssignPorts(ctx, fabricName, tenantName,
    ones_gfx.GpuActionAdd, []string{"su00-rack00-node00"}, []int{1, 2, 3}, nil)

// Assign all ports on each server — pass nil for gpuIds
ports, err = client.GPU.AssignPorts(ctx, fabricName, tenantName,
    ones_gfx.GpuActionAdd, []string{"su00-rack00-node00"}, nil, nil)

// Unassign — same call with GpuActionDelete
ports, err = client.GPU.AssignPorts(ctx, fabricName, tenantName,
    ones_gfx.GpuActionDelete, []string{"su00-rack00-node00"}, []int{1, 2, 3}, nil)
```

### Partial GPU Allocation

Map or unmap specific GPUs to a tenant on a shared server. This is distinct
from attaching a whole server via `Tenants.Update`: `ModifyAllocations` gives
fine-grained control over which GPU indices a tenant can access.

```go
// SuidMap is map[serverIndex]map[hostname][]gpuID
suid := resources.SuidMap{
    "0": {
        "hgx-su00-h00": []string{"G0", "G1", "G2", "G3"},
    },
}

resp, err := client.Tenants.ModifyAllocations(ctx, fabricName, tenantName,
    suid,
    ptr(ones_gfx.GpuActionAdd), // or ones_gfx.GpuActionDelete
    nil,                        // configScope (default WHOLE_SERVER)
    nil,                        // unreachableDevices
)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("status=%s  msg=%s\n", resp.Status, resp.Message)

// Async variant
op, err := client.Tenants.ModifyAllocationsAsync(ctx, fabricName, tenantName,
    suid, ptr(ones_gfx.GpuActionAdd), nil, nil)
```

For per-GPU scope, pass `ptr(resources.ConfigScopeParticularGPU)` as
`configScope`.

**Notes:**
- The server must already be attached to the tenant via `Tenants.Update` before mapping individual GPUs.
- Remove per-GPU mappings with `GpuActionDelete` before detaching the server.
- This endpoint is valid only for externally managed fabrics; on ONES-controlled fabrics it returns `409 FABRIC_NOT_EXTERNALLY_MANAGED` (use the tenant-update flow instead).
- Multiple tenants can share the same physical server (e.g. G0–G3 → tenant-A, G4–G7 → tenant-B) when the server is attached with `Shared: ptr(true)`.

---

## Operation Modes

Async-capable endpoints (`Tenants.Create/Delete/Update/ModifyAllocations`,
`VPCPeering.Create`) come in Sync and Async pairs.

### Synchronous (Default)

Blocks until the server completes the operation and returns the result.

```go
msg, err := client.Tenants.Create(ctx, fabricName, "tenant1", nil, ptr(8), nil)
```

### Asynchronous (Poll)

Returns immediately with an `*ones_gfx.OperationAccepted`. You poll for completion.

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ptr(8), nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println("operation:", op.OperationID, op.Status)

for {
    current, err := client.Operations.Get(ctx, op.OperationID)
    if err != nil {
        log.Fatal(err)
    }
    if ones_gfx.OperationStatus(current.Status).IsTerminal() {
        if current.Status == string(ones_gfx.OperationStatusSuccess) {
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

Returns immediately. Server POSTs the result to your webhook URL.

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ptr(8), nil,
    ones_gfx.WithWebhook("http://receiver:8000/hook", []string{"tenant.create"}))

fmt.Printf("Operation submitted: %s (webhook: %v)\n",
    op.OperationID, op.WebhookRegistered)
// No polling needed — result sent to webhook
```

### Idempotency

Retrying an async call with the same key replays the original operation instead
of creating a duplicate:

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ptr(8), nil,
    ones_gfx.WithIdempotencyKey("ONES-tenant1-CRT-20260828"))
```

---

## Timeout Override

Long-running operations (GPU attach/detach) can take 10–15 minutes. Override the
default timeout:

```go
// Client-level default: 20 minutes
client := ones.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute))

// Per-call override for this specific operation
resp, err := client.Tenants.Update(ctx, fabricName, tenantName, servers,
    ptr(ones_gfx.GpuActionAdd),
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

```bash
cd ones-gfx-sdk/ones-gfx-sdk-go

# Build and vet everything
go build ./...
go vet ./...

# Or via make
make build-lib
make vet

# Format check
gofmt -l ones_gfx/ sdk/ examples/
```

### Check for Issues

```bash
# Unused code detection (requires staticcheck)
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

---

## Testing the SDK

A runnable CLI tour lives in `examples/`:

```bash
go run examples/usage_examples.go   # full CLI tour
go run ./examples/index             # minimal login + fabrics walkthrough
```

### Manual smoke test (requires a ONES instance)

```go
package main

import (
    "context"
    "log"
    "time"

    ones_gfx "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
    "github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
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

    client := ones.NewClient(
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

    fabricName := fabrics[0].Name
    log.Printf("Using fabric: %s\n", fabricName)

    // Test 2: List tenants
    log.Println("Test 2: Listing tenants...")
    tenants, err := client.Tenants.List(ctx, fabricName)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Tenants payload: %v\n", tenants)

    // Test 3: Available servers
    log.Println("Test 3: Checking available servers...")
    servers, err := client.GPU.AvailableServers(ctx, fabricName)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Available servers: %v\n", servers.AvailableGPUs)

    // Test 4: Controller version
    ver, err := client.System.ControllerVersion(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Controller: %s %s\n", ver.AppName, ver.AppVersion)

    log.Println("All tests passed!")
}
```

---

## File Structure

```
ones-gfx-sdk/ones-gfx-sdk-go/
├── go.mod                      # Module definition
├── Makefile                    # build-lib / vet / fmt / check targets
├── README.md                   # This file
├── IMPLEMENTATION_STATUS.md    # Typed surface + method index
├── api_reference.txt           # Language-agnostic curl reference
├── ones/                       # ★ Public API — the only package you import
│   ├── api.go                  # 102 flat functions (generated), auth-guarded
│   ├── types.go                # Re-exported types/consts/options (generated)
│   ├── session.go              # Init / GetClient / Login / Refresh / Logout
│   └── client.go               # Client struct + resource handles
├── ones_gfx/                   # Internals: transport, auth, shared types
│   ├── auth.go                 # JWTAuth (bearer access/refresh pair)
│   ├── client.go               # TokenAuth — single-token login flow
│   ├── transport.go            # HTTP client + Call[T]/CallMultipart[T]
│   ├── apitypes.go             # Types shared by >1 resource
│   ├── enums.go                # Operation modes, statuses
│   ├── errors.go               # Typed errors
│   ├── gpu_allocation.go       # GPUOperation ADD/DELETE enum
│   ├── models.go               # Operation struct + timestamp parsing
│   ├── options.go              # Functional call options
│   └── resources/              # One file per API domain (97 endpoints)
│       ├── bootstrap.go
│       ├── configmgmt.go
│       ├── devices.go
│       ├── fabrics.go
│       ├── fabricsims.go
│       ├── files.go
│       ├── gpu.go
│       ├── hosttenants.go
│       ├── intents.go
│       ├── inventory.go
│       ├── netops.go
│       ├── nmxc.go
│       ├── operations.go
│       ├── rma.go
│       ├── system.go
│       ├── tenants.go
│       └── vpcpeering.go
└── examples/
    ├── usage_examples.go       # Runnable CLI tour
    └── index/main.go           # Minimal login + fabrics walkthrough
```

Note the module root holds no Go package, so `go get` fetches a library rather
than installing a binary.

---

## Troubleshooting

### Issue: `certificate signed by unknown authority`

**Solution:** Disable TLS verification for dev/lab (self-signed certs):

```go
client := ones.NewClient(baseURL, auth,
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

client := ones.NewClient(baseURL, auth,
    ones_gfx.WithTLSConfig(tlsConfig))
```

### Issue: `context deadline exceeded`

**Solution:** Increase timeout for long operations:

```go
// Client-level
client := ones.NewClient(baseURL, auth,
    ones_gfx.WithClientTimeout(20*time.Minute))

// Per-call
resp, err := client.Tenants.Update(ctx, fabricName, tenantName, servers,
    ptr(ones_gfx.GpuActionAdd),
    ones_gfx.WithTimeout(15*time.Minute))
```

### Issue: `authentication error: token refresh failed`

**Solution:** Verify refresh token is valid and refresh URL is correct. Both tokens rotate on every refresh.

---

## Known Limitations

- **No login endpoint** — partners supply tokens obtained out-of-band.
- **No CLI binary** — the core library is the deliverable.
- **16 open-payload leaves** — a handful of responses are returned as
  `interface{}` / `map[string]interface{}` because the upstream Java service
  returns `Object` / `Map<String,Object>` on those paths. The full list is in
  [IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md).
- **Legacy transport verbs** — `Transport.Get/Post/Patch/Delete` remain for
  back-compat but are unused by the resources; prefer the typed resource methods.

---

## Support

- **Documentation:** This README and [IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md)
- **Issues:** [GitHub Issues](https://github.com/aviznetworks/ones-gfx-sdk/issues)
- **API Reference:** [api_reference.txt](./api_reference.txt) (language-agnostic curl examples)
