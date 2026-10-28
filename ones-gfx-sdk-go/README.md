# ONES GFX SDK — Go

Go client library for **AVIZ ONES Spectrum-X** tenant management API (v4.2).

## Status

**Core Library:** ✅ The 20 Fabric Manager endpoints documented in *ONES
Spectrum-X API Operations v3.0* are exposed as 25 typed resource methods (20
plus the `...Async` variants) across 8 resource handles.
**CLI Binary:** 🚧 Not started.

Typed request params and typed return values throughout; the only untyped
returns are the 6 documented open-payload leaves (see
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

    fabrics, err := ones.GetFabrics(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range fabrics.Fabrics {
        fmt.Printf("%d %s\n", f.ID, f.FabricName)
    }
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

Names flatten `Resource.Method` — e.g. `Tenants.Create` → `ones.CreateTenant`,
`Fabrics.ListDTOs` → `ones.GetFabrics`, `Tenants.ModifyAllocations` →
`ones.ModifyGpuAllocations`. The full index is in
[IMPLEMENTATION_STATUS.md](./IMPLEMENTATION_STATUS.md).

---

## Conventions

Before the examples, the rules that apply to every resource method:

1. **Path params are positional**; body/query fields are individual typed args.
2. **Optional fields are pointers.** `nil` omits them from the request, which
   mirrors the Python SDK's `| None = None`. Build them with `ones.Ptr`:
   `ones.Ptr(8)`, `ones.Ptr(false)`, `ones.Ptr(ones.GpuActionAdd)`. You do not
   need a helper of your own.
3. **Calls with several optionals also take a request struct.** `CreateWith`,
   `UpdateWith`, `ModifyAllocationsWith` and `AssignPortsWith` (and their
   `ones.*With` flat forms) take a struct with named fields instead, so
   arguments cannot be passed in the wrong order. See
   [Request structs](#request-structs).
4. **Async-capable endpoints have a separate `...Async` method** returning
   `*ones_gfx.OperationAccepted`.
5. **Arguments are validated before anything is sent.** A bad argument returns
   `*ones.ValidationError` and no request reaches the server. See
   [Client-side validation](#client-side-validation).
6. **Enums are typed.** `GpuAction`, `ConfigScope`, `OperationStatus` and
   `ConfigStatus` have named constants. For values from user input, use the
   `ones.Parse…` functions (case-insensitive, error on unknown input) or
   `.Valid()`.

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
    fabrics, err := client.Fabrics.ListDTOs(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, fabric := range fabrics.Fabrics {
        fmt.Printf("Fabric: %s (%s)\n", fabric.FabricName, fabric.Description)
    }

    // 4. Create a tenant (synchronous). Returns the server's message string.
    msg, err := client.Tenants.Create(ctx, "sdk", "demo-tenant",
        ones.Ptr("Created via Go SDK"), // description
        ones.Ptr(8),                    // maxGpusAllowed
        ones.Ptr(false),                // shared
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Create tenant:", msg)

    // 5. Attach GPU servers to the tenant, with a per-call timeout override
    resp, err := client.Tenants.Update(ctx, "sdk", "demo-tenant",
        []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
        ones.Ptr(ones_gfx.GpuActionAdd),
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
        ones.Ptr(ones_gfx.GpuActionDelete),
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

Resource handles: `Tenants`, `GPU`, `Fabrics`, `Inventory`, `Devices`, `NMXC`,
`VPCPeering`, `Operations`.

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
// List fabrics (GET /fabrics)
dtos, err := client.Fabrics.ListDTOs(ctx)         // *resources.FabricsListResponse

// Inventory sync — UFM enabled fabrics only (lives on Inventory)
res, err := client.Inventory.Sync(ctx, fabricName) // *resources.InventorySyncResult
fmt.Println(res.Message)
```

### Tenants

```go
// List / Get: typed ([]ones.Tenant, *ones.Tenant). Fields FM adds that the
// SDK does not model yet are kept in Tenant.Extra; ListRaw/GetRaw return the
// response as a raw map.
tenants, err := client.Tenants.List(ctx, fabricName)
tenant, err := client.Tenants.Get(ctx, fabricName, tenantName)
fmt.Println(tenant.Name, tenant.ConfigStatus, tenant.AllottedServers())

// Servers with free GPUs (lives on GPU)
servers, err := client.GPU.AvailableServers(ctx, fabricName)
fmt.Println(servers.AvailableGPUs)

// Create (sync / async)
msg, err := client.Tenants.Create(ctx, fabricName, "tenant1",
    ones.Ptr("description"), ones.Ptr(8), ones.Ptr(false))
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1",
    ones.Ptr("description"), ones.Ptr(8), ones.Ptr(false))

// Delete (sync / async)
resp, err := client.Tenants.Delete(ctx, fabricName, tenantName)
op, err = client.Tenants.DeleteAsync(ctx, fabricName, tenantName)

// Attach / detach whole servers (replaces the old AllocateGPUs/DeallocateGPUs)
resp, err = client.Tenants.Update(ctx, fabricName, tenantName,
    []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
    ones.Ptr(ones_gfx.GpuActionAdd),
    ones_gfx.WithTimeout(15*time.Minute),
)
op, err = client.Tenants.UpdateAsync(ctx, fabricName, tenantName,
    []resources.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
    ones.Ptr(ones_gfx.GpuActionDelete),
)

// Share a server between tenants
servers := []resources.GpuServerInfo{{ServerName: "hgx-su00-h00", Shared: ones.Ptr(true)}}

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

### VPC Peering

```go
// Create route leaking between two VPCs (sync / async)
msg, err := client.VPCPeering.Create(ctx, fabricName,
    "tenant1-storage-route-leak", "tenant1-Fab101-north-south", "Fab101-Storage")
op, err := client.VPCPeering.CreateAsync(ctx, fabricName, name, vpcname, peervpcname)

// Delete it again — synchronous only. name addresses the peering; FM rejects
// the call (400 VPC_PEERING_PAIR_MISMATCH) unless the VPC pair matches it.
resp, err := client.VPCPeering.Delete(ctx, fabricName,
    "tenant1-storage-route-leak", "tenant1-Fab101-north-south", "Fab101-Storage")
// or: ones.DeleteVpcPeering(ctx, client, fabricName, name, vpcname, peervpcname)
fmt.Println(resp.Message)                          // *resources.ApiResponseMessage
```

An unknown peering or fabric returns `*NotFoundError`. A wrong pair or invalid
name returns `*BadRequestError`, and a switch failure returns `*ServerError`.
The FM `error.code` (for example `VPC_PEERING_NOT_FOUND`) is in the error's
`ResponseBody`.

### Request structs

The calls that take several optional pointers also accept a request struct.
Fields are named, so nothing can be passed in the wrong position. Optional
fields use the zero value for "not set" unless the zero value means something
on its own (`Shared: false` differs from "keep the server default"), in which
case they are pointers.

```go
// Instead of Tenants.Create(ctx, "f1", "t1", ones.Ptr("demo"), ones.Ptr(8), nil)
msg, err := ones.CreateTenantWith(ctx, client, ones.CreateTenantRequest{
    FabricName:     "f1",
    TenantName:     "t1",
    Description:    "demo",
    MaxGpusAllowed: 8,                  // or ones.UnlimitedGPUs; 0 = server default
})

// Action is required. There is no default, so a detach can never turn into an attach.
resp, err := ones.UpdateTenantWith(ctx, client, ones.UpdateTenantServersRequest{
    FabricName: "f1", TenantName: "t1",
    Servers:    []ones.GpuServerInfo{{ServerName: "hgx-su00-h00"}},
    Action:     ones.GpuActionDelete,
})

// Instead of AssignPorts(ctx, f, t, ones.GpuActionAdd, servers, nil, nil)
ports, err := ones.AssignGpuPortsWith(ctx, client, ones.AssignPortsRequest{
    FabricName: "f1", TenantName: "t1",
    Action:      ones.GpuActionAdd,
    ServerNames: []string{"su00-rack00-node00"}, // GpuIDs empty = every port
})

res, err := ones.ModifyGpuAllocationsWith(ctx, client, ones.ModifyAllocationsRequest{
    FabricName: "f1", TenantName: "t1",
    Suid:   ones.SuidMap{"0": {"hgx-su00-h00": {Gpus: []string{"G0", "G1"}}}},
    Action: ones.GpuActionAdd,
    Scope:  ones.ConfigScopeParticularGPU, // empty = server default
})
```

Each request type has a `Validate()` method for checking it ahead of time. The
`Async` forms (`CreateTenantAsyncWith`, `UpdateTenantAsyncWith`,
`ModifyGpuAllocationsAsyncWith`) take the same structs. The positional methods
are unchanged.

### Inventory (typed)

```go
// UFM ports grouped by host (GET /fabrics/{f}/inventoryPorts)
ports, err := ones.GetInventoryPorts(ctx, client, "UFM1")   // *ones.InventoryPortsResult
fmt.Println(ports.AvailablePortCount, "of", ports.TotalPortCount, "ports free")
if host, ok := ports.Host("hgx-su00-h00"); ok {              // hostname or UFM system name
    for _, p := range host.Ports {
        fmt.Println(p.PortName, p.Allocated, p.TenantName)
    }
}

hosts, err := ones.GetInventoryHosts(ctx, client, "UFM1")   // *ones.UFMHostsResult (live from UFM)
sync, err := ones.InventorySync(ctx, client, "UFM1")        // *ones.InventorySyncResult

// NMX-C domains, partitions, trays and GPUs (GET /fabrics/{f}/nmxc/inventory)
inv, err := ones.GetNmxcInventory(ctx, client, "NMXC1")     // *ones.NmxcInventory
for _, d := range inv.Domains {
    fmt.Println(d.DomainName, d.ConnectionState, d.FreeGpuCount)
    // d.DomainID.String() is what ResetNmxcDomain takes
}
```

Timestamps are `ones.Timestamp`. FM sends some with a zone and some without,
so the value is kept as sent; `.Time()` parses it and `.String()` returns it
unchanged. Top-level keys the SDK does not model yet are kept in `Extra`.
`...Raw` variants (`GetInventoryPortsRaw`, `GetInventoryHostsRaw`,
`InventorySyncRaw`, `GetNmxcInventoryRaw`) return the response as a raw map.

### HBN VF interfaces

```go
// GET /fabrics/{f}/servers/{s}/vf-interfaces (VFs grouped by DPU)
list, err := ones.ListVfInterfaces(ctx, client, "f1", "hgx-su00-h00")
for _, vf := range list.Flatten() {
    fmt.Println(vf.IfName, vf.ServerIf, vf.Status, vf.TenantName)
}

// POST .../vf-interfaces/{vfId}/assign. vfId is the DPU or host interface name.
a, err := ones.AssignVfInterface(ctx, client, "f1", "hgx-su00-h00", "pf0vf0_if", "t1")
fmt.Println(a.VlanID, a.IPAddress, a.IPAssignStatus)

// DELETE .../vf-interfaces/{vfId}/assign. Unassign VFs before detaching the
// server from the tenant; FM refuses once the server has left it.
_, err = ones.UnassignVfInterface(ctx, client, "f1", "hgx-su00-h00", "pf0vf0_if")
```

The resource forms are `client.VFs.List`, `Assign` and `Unassign`.

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
    ones.Ptr(ones_gfx.GpuActionAdd), // or ones_gfx.GpuActionDelete
    nil,                        // configScope (default WHOLE_SERVER)
    nil,                        // unreachableDevices
)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("status=%s  msg=%s\n", resp.Status, resp.Message)

// Async variant
op, err := client.Tenants.ModifyAllocationsAsync(ctx, fabricName, tenantName,
    suid, ones.Ptr(ones_gfx.GpuActionAdd), nil, nil)
```

For per-GPU scope, pass `ones.Ptr(resources.ConfigScopeParticularGPU)` as
`configScope`.

**Notes:**
- The server must already be attached to the tenant via `Tenants.Update` before mapping individual GPUs.
- Remove per-GPU mappings with `GpuActionDelete` before detaching the server.
- This endpoint is valid only for externally managed fabrics; on ONES-controlled fabrics it returns `409 FABRIC_NOT_EXTERNALLY_MANAGED` (use the tenant-update flow instead).
- Multiple tenants can share the same physical server (e.g. G0–G3 → tenant-A, G4–G7 → tenant-B) when the server is attached with `Shared: ones.Ptr(true)`.

---

## Operation Modes

Async-capable endpoints (`Tenants.Create/Delete/Update/ModifyAllocations`,
`VPCPeering.Create`) come in Sync and Async pairs. `VPCPeering.Delete` is
synchronous only.

### Synchronous (Default)

Blocks until the server completes the operation and returns the result.

```go
msg, err := client.Tenants.Create(ctx, fabricName, "tenant1", nil, ones.Ptr(8), nil)
```

### Asynchronous (Poll)

Returns immediately with an `*ones_gfx.OperationAccepted`. You poll for completion.

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ones.Ptr(8), nil)
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
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ones.Ptr(8), nil,
    ones_gfx.WithWebhook("http://receiver:8000/hook", []string{"tenant.create"}))

fmt.Printf("Operation submitted: %s (webhook: %v)\n",
    op.OperationID, op.WebhookRegistered)
// No polling needed — result sent to webhook
```

### Idempotency

Retrying an async call with the same key replays the original operation instead
of creating a duplicate:

```go
op, err := client.Tenants.CreateAsync(ctx, fabricName, "tenant1", nil, ones.Ptr(8), nil,
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
    ones.Ptr(ones_gfx.GpuActionAdd),
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
- `ValidationError` — rejected by the SDK before sending; see below
- `TransportError` — Network/TLS/timeout
- `AuthenticationError` — HTTP 401/403 or refresh failure
- `BadRequestError` — HTTP 400
- `NotFoundError` — HTTP 404
- `ConflictError` — HTTP 409
- `ServerError` — HTTP 5xx
- `OperationFailedError` — Async job ended with status FAILURE

### Client-side validation

Before building a request, every method checks its arguments. A bad argument
returns `*ones.ValidationError` (it also satisfies `ONESError`), and no request
reaches the server:

```go
_, err := ones.CreateTenant(ctx, client, "f1", "t1", nil, ones.Ptr(0), nil)
var ve *ones.ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Field)   // maxGpusAllowed
    fmt.Println(err)        // invalid maxGpusAllowed: 0 is not allowed: use -1 for unlimited or a positive number
}
```

| Check | Applies to |
|---|---|
| Not empty; no `/ ? # \ %`; not `.` or `..` | every name that goes into the URL path: fabric, tenant, server, VF, rail, policy, domain, webhook and operation ids |
| `-1` (`ones.UnlimitedGPUs`) or a positive number | `maxGpusAllowed` |
| `ADD` or `DELETE` | GPU actions (required on `AssignPorts`, `UpdateNmxcDomains` and the `...With` structs; optional elsewhere) |
| `WHOLE_SERVER` or `PARTICULAR_GPU` | `configScope` |
| At least one entry, none blank | `servers`, `serverNames`, `suid` (numeric indexes), NMX-C `domains` |
| No negative ids | `gpuIds` |
| Not empty, at most 24 characters | network policy name on create |
| Not empty | VPC peering `name` / `vpcname` / `peervpcname`, VF-assign and shared-services `tenantName` |

The checks only cover what is wrong on every fabric type. Rules that vary by
fabric, such as tenant-name length and character set, are left to Fabric
Manager, so the SDK never rejects a name FM would accept.

---

## Breaking changes and migration

Most of this release is additive: request structs, `...With` methods,
`Valid`/`Parse…` enum helpers, VF wrappers and validation. Two changes can
affect existing code:

**1. Inventory methods return typed results.**

| Call | Before | Now | Old behaviour |
|---|---|---|---|
| `Inventory.Ports` / `ones.GetInventoryPorts` | `map[string]interface{}` | `*InventoryPortsResult` | `PortsRaw` / `GetInventoryPortsRaw` |
| `Inventory.UFMHosts` / `ones.GetInventoryHosts` | `map[string]interface{}` | `*UFMHostsResult` | `UFMHostsRaw` / `GetInventoryHostsRaw` |
| `Inventory.Sync` / `ones.InventorySync` | `map[string]interface{}` | `*InventorySyncResult` | `SyncRaw` / `InventorySyncRaw` |
| `NMXC.Inventory` / `ones.GetNmxcInventory` | `map[string]interface{}` | `*NmxcInventory` | `InventoryRaw` / `GetNmxcInventoryRaw` |

Code that only checked the error (`_, err := ...`) is unaffected. Code that
indexed the map should use the struct fields
(`res["hosts"].([]interface{})` → `res.Hosts`), or switch to the `...Raw`
method to keep the old shape unchanged.

**2. Invalid arguments now fail before sending.** Calls that previously reached
Fabric Manager and came back as `*BadRequestError` (or hit the wrong endpoint
because a name contained `/`) now return `*ValidationError` without making a
request. Code that matched `*BadRequestError` for these cases should also
match `*ValidationError`, or match `ONESError`, which covers both. Valid
requests send exactly the same bytes as before.

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

Unit tests run against an in-process `httptest` server, so they need no ONES
instance:

```bash
go test ./...      # or: make test
```

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
    fabrics, err := client.Fabrics.ListDTOs(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found %d fabric(s)\n", len(fabrics.Fabrics))
    if len(fabrics.Fabrics) == 0 {
        log.Fatal("No fabrics found - cannot proceed with tests")
    }

    fabricName := fabrics.Fabrics[0].FabricName
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
│   ├── api.go                  # 25 flat functions (generated), auth-guarded
│   ├── types.go                # Re-exported types/consts/options (generated)
│   ├── session.go              # Init / GetClient / Login / Refresh / Logout
│   └── client.go               # Client struct + resource handles
├── ones_gfx/                   # Internals: transport, auth, shared types
│   ├── auth.go                 # JWTAuth (bearer access/refresh pair)
│   ├── client.go               # TokenAuth — single-token login flow
│   ├── transport.go            # HTTP client + Call[T]/CallMultipart[T]
│   ├── apitypes.go             # Types shared by >1 resource
│   ├── enums.go                # Operation modes, statuses
│   ├── errors.go               # Typed errors, incl. ValidationError
│   ├── gpu_allocation.go       # GPUOperation ADD/DELETE enum
│   ├── models.go               # Operation struct + timestamp parsing
│   ├── options.go              # Functional call options
│   ├── timestamp.go            # Timestamp: FM time values, kept as sent
│   └── resources/              # One file per API domain (20 endpoints)
│       ├── devices.go
│       ├── fabrics.go
│       ├── gpu.go
│       ├── inventory.go
│       ├── nmxc.go
│       ├── operations.go
│       ├── requests.go         # Request structs for the ...With methods
│       ├── tenants.go
│       ├── validation.go       # Client-side argument checks
│       ├── vfinterfaces.go
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
    ones.Ptr(ones_gfx.GpuActionAdd),
    ones_gfx.WithTimeout(15*time.Minute))
```

### Issue: `authentication error: token refresh failed`

**Solution:** Verify refresh token is valid and refresh URL is correct. Both tokens rotate on every refresh.

---

## Known Limitations

- **No login endpoint** — partners supply tokens obtained out-of-band.
- **No CLI binary** — the core library is the deliverable.
- **6 open-payload leaves** — a handful of responses are returned as
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
