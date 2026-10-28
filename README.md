# ONES GFX SDK

Multi-language SDK for **AVIZ ONES Spectrum-X** tenant management API (v4.2.1).
Python and Go implementations expose the **same 130 Fabric Manager endpoints**
through the same `Client` shape, so code translates almost line-for-line between
them.

Covers the endpoints documented in *ONES Spectrum-X API Operations v3.0*:
tenant lifecycle, GPU allocation and per-port assignment, fabric listing,
NMX-C and UFM integration, VPC peering, and async operation tracking.

---

## Install

Neither SDK is published to a registry — install both from this checkout.

**Python** (3.9+) — editable install from the local directory:

```bash
pip install -e ./ones-gfx-sdk-python
```

**Go** (1.19+, zero dependencies) — the module lives in this repo, so point at
it with a `replace` directive in your own `go.mod`:

```
require github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go v0.0.0

replace github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go => /path/to/ones-gfx-sdk/ones-gfx-sdk-go
```

The import path stays `github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones`;
only resolution is local.

---

## The `Client` API

A `Client` is a connection to **one** ONES instance. It holds the base URL, the
auth token, and optionally the credentials used to obtain it — all per instance.
Nothing is global and nothing is written to disk, so you can hold several
clients for different hosts or users at once.

You then **pass the client into every API function**.

### Creating a client

| | Python | Go |
|---|---|---|
| From credentials | `Client.initialize_with_creds(baseUrl, username, password)` | `ones.InitializeWithCreds(baseUrl, username, password)` |
| From an existing token | `Client.initialize_with_token(baseUrl, token)` | `ones.InitializeWithToken(baseUrl, token)` |
| Empty, log in later | `Client(baseUrl)` | `ones.NewClient(baseUrl)` |

`initialize_with_creds` logs in immediately, so a returned client is always ready.
`initialize_with_token` **errors if the token is empty** (`ValueError` in Python,
an `error` in Go).

Pass the **root** base URL, e.g. `https://10.4.5.76:8089`. Both SDKs derive the
auth endpoints (`/api/user/...`) and the Fabric Manager base (`/api/fm/...`)
from it, and strip any trailing slash.

### Options

| Option | Python | Go |
|---|---|---|
| Request timeout | `requests` default | `ones.DefaultTimeout` = 20 min |
| Custom auth / transport options | build `Client(...)` directly | `ones.NewClientWithAuth(baseUrl, auth, opts...)` |

### Methods

Identical on both sides — Python is snake_case, Go is exported PascalCase.

| Purpose | Python | Go |
|---|---|---|
| Base URL | `get_base_url()` / `set_base_url(v)` | `GetBaseUrl()` / `SetBaseUrl(v) error` |
| Token | `get_auth_token()` / `set_auth_token(v)` | `GetAuthToken()` / `SetAuthToken(v)` |
| Username | `get_username()` / `set_username(v)` | `GetUsername()` / `SetUsername(v)` |
| Password | `get_password()` / `set_password(v)` | `GetPassword()` / `SetPassword(v)` |
| Rotate token (rejects empty) | `update_token(token)` | `UpdateToken(token) error` |
| Is a token held? | `is_authenticated()` | `IsAuthenticated()` |
| Log in (falls back to stored creds) | `login(username=None, password=None)` | `Login(username, password)` |
| Exchange token for a new one | `refresh_auth()` | `RefreshAuth()` |
| Log out, clear token | `logout()` | `Logout()` |
| Release HTTP resources | `close()` (or `with`) | `Close()` |
| Raw call for uncovered routes | `call_api(method, path, ...)` | `GetTransport()` |

### Calling the API

Every endpoint is a flat function taking the client. Go additionally takes a
`context.Context` first, and offers `...Async` variants for the 5 async-capable
endpoints (135 functions vs Python's 130; Python selects async with the `prefer`
argument on the same function).

```python
from ones_gfx.apis import get_fabrics
fabrics = get_fabrics(client)
```

```go
fabrics, err := ones.GetFabrics(ctx, client)
```

Names differ only in case convention: `get_fabrics` ↔ `GetFabrics`,
`create_tenant` ↔ `CreateTenant`, `modify_gpu_allocations` ↔
`ModifyGpuAllocations`. Optional fields are keyword args defaulting to `None` in
Python, and pointers (`ones.Ptr(v)`, `nil` to omit) in Go.

### Errors

| Condition | Python | Go |
|---|---|---|
| Call with no token | raises `NotAuthenticatedError` | `ErrNotAuthenticated` |
| Nil/missing client | n/a | `ErrNilClient` |
| Session method on a non-token client | n/a | `ErrNoTokenAuth` |
| HTTP 4xx/5xx | `requests.HTTPError` | typed: `*AuthenticationError`, `*NotFoundError`, `*ConflictError`, `*ServerError`, … (use `errors.As`) |
| Invalid argument, caught before sending | n/a (sent to the server) | `*ValidationError`, naming the field ([details](./ones-gfx-sdk-go/README.md#client-side-validation)) |

### Example: deleting a VPC peering

`DELETE /fabrics/{fabricName}/vpcpeering` is the inverse of creating one. It is
synchronous only. The peering is addressed by `name`, and Fabric Manager rejects
the call unless `vpcname`/`peervpcname` match the pair recorded for it. On
success both SDKs return the `{"status": "success", "message": ...}` envelope.

```python
import requests
from ones_gfx.apis import delete_vpc_peering

try:
    resp = delete_vpc_peering(
        client, "Fab101",
        name="tenant1-storage-route-leak",
        vpcname="tenant1-Fab101-north-south",
        peervpcname="Fab101-Storage",
    )
    print(resp["message"])
except requests.HTTPError as e:
    code = e.response.json().get("error", {}).get("code")  # e.g. VPC_PEERING_NOT_FOUND
    print(e.response.status_code, code)
```

```go
resp, err := ones.DeleteVpcPeering(ctx, client, "Fab101",
    "tenant1-storage-route-leak", "tenant1-Fab101-north-south", "Fab101-Storage")
var notFound *ones.NotFoundError
switch {
case errors.As(err, &notFound):
    // 404 VPC_PEERING_NOT_FOUND or FABRIC_NOT_FOUND
case err != nil:
    log.Fatal(err)
default:
    fmt.Println(resp.Message)
}
```

| Status | `error.code` | Go error type |
|---|---|---|
| 404 | `VPC_PEERING_NOT_FOUND`, `FABRIC_NOT_FOUND` | `*NotFoundError` |
| 400 | `VPC_PEERING_PAIR_MISMATCH`, `VPC_PEERING_NAME_REQUIRED`, `VPC_NAME_REQUIRED`, `PEER_VPC_NAME_REQUIRED`, `INVALID_VPC_NAME`, `TENANT_NOT_FOUND`, … | `*BadRequestError` |
| 401 / 403 | — | `*AuthenticationError` |
| 500 | `VPC_PEERING_DELETE_FAILED` (a switch could not be updated) | `*ServerError` |

In Python the same body is on `e.response`. In Go it is on the error's
`ResponseBody`.

---

## Example — Python

```python
from ones_gfx import Client
from ones_gfx.apis import get_fabrics, list_tenants

client = Client.initialize_with_creds("https://10.4.5.76:8089", "admin", "secret")
print("authenticated:", client.is_authenticated())

for f in get_fabrics(client)["fabrics"]:
    print(f"  {f.get('id')} {f.get('fabricName')}")
    print(list_tenants(client, f["fabricName"]))
```

## Example — Go

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

    client, err := ones.InitializeWithCreds("https://10.4.5.76:8089", "admin", "secret")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    fmt.Println("authenticated:", client.IsAuthenticated())

    fabrics, err := ones.GetFabrics(ctx, client)
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range fabrics.Fabrics {
        fmt.Printf("  %d %s\n", f.ID, f.FabricName)
    }
}
```

`ones` is the only import you need — every type, constant and option is
re-exported there (`ones.GpuServerInfo`, `ones.GpuActionAdd`, `ones.Ptr`, …).

Runnable versions of both live at
[`ones-gfx-sdk-python/index.py`](./ones-gfx-sdk-python/index.py) and
[`ones-gfx-sdk-go/examples/index`](./ones-gfx-sdk-go/examples/index):

```bash
cd ones-gfx-sdk-python && python3 index.py
cd ones-gfx-sdk-go     && go run ./examples/index
```

---

## Asynchronous work: which contract answers what

Fabric Manager exposes two async contracts. They are **layered, not
alternatives**, and the common mistake is treating either one as the whole
answer.

| | Question it answers | How |
|---|---|---|
| **Operations API** | Did FM accept and process my request? | 202 → `operationId` → poll `Operations.Get`, or receive a webhook |
| **`config_status`** | Did the configuration reach the devices? | Poll the tenant, watch `config_status` |

**An operation reaches SUCCESS before the devices are configured.** Creating a
tenant asynchronously, `AsyncTenantService` marks the operation SUCCESS as soon
as the tenant row is saved and its IPAM pool provisioned — at that moment the
tenant's `config_status` is still `2` (NOT_STARTED). The device push runs
afterwards and is what drives `2 → 3 → 1`, or `0` on failure.

So an operation-only wait tells you FM took the request, not that the fabric
carries the change. A `config_status`-only wait misses the case where FM
rejected the request outright. Most workflows want both, in that order:

```go
res, err := client.Tenants.AwaitTenantChange(
    ctx, client.Operations,
    accepted.OperationID,       // "" if the call was synchronous
    "fabric-1", "tenant-a",
    resources.WaitOptions{Interval: 5 * time.Second, Timeout: 10 * time.Minute},
)
if res.Applied() {
    // operation processed AND config_status == PASS
}
```

Either layer can be waited on alone when that is genuinely what you want:

```go
op, _ := client.Operations.WaitForOperation(ctx, opID, opts)          // layer 1
tn, _ := client.Tenants.WaitForTenantConfigApplied(ctx, f, t, opts)   // layer 2
```

A **synchronous** call collapses layer 1 only: it returns once FM has processed
the request, and `config_status` may still be settling. That is the same gap
the sync-PATCH timeout fallback covers.

Running out of time returns `WaitTimeoutError`, which is deliberately not a
failure — the work may still be in progress.

## Configuration status

`config_status` is the one field in this API whose numbering actively misleads.
It is positional, not a progression:

| Value | Name | Meaning |
|---|---|---|
| 0 | `FAIL` | configuration failed |
| 1 | `PASS` | **configured** — the only success value |
| 2 | `NOT_STARTED` | configuration has not begun — **not** a success |
| 3 | `IN_PROGRESS` | configuration is underway |
| 4 | `BROWN_FIELD` | skipped: the fabric deploys as brownfield |
| 5 | `SKIPPED` | skipped: pre-validation or minimal-config excluded the device |

The trap is **2**. It sits above `PASS`, so it looks like a later, stronger
state — and a `status >= 1` check silently reports a device that was never
configured as configured. Compare against the named values, never the ordinal:

```python
from ones_gfx import ConfigStatus

if tenant.config_status.is_success:        # PASS only
    ...
if tenant.config_status.is_terminal:       # stop polling
    ...
ConfigStatus.parse("PASS")                 # strict: raises on anything unknown
ConfigStatus.from_api(raw)                 # lenient: unknown -> NOT_STARTED
```

```go
if tn.ConfigStatus.IsSuccess() { ... }     // PASS only
if tn.ConfigStatus.IsTerminal() { ... }    // stop polling
ones.ParseConfigStatus("PASS")             // strict: returns an error
```

Fabric Manager sends these as JSON **strings** (`"1"`, not `1`). Decoding a
response is lenient — an unreadable status becomes `NOT_STARTED` rather than
discarding an otherwise good payload. Values arriving from anywhere else should
go through the strict parser, which reports unknown input instead of hiding it.

`BROWN_FIELD` and `SKIPPED` both mean *deliberately not configured*. They are
terminal, and they are neither a failure nor a success — treat them as their
own outcome.

## Timeouts on synchronous writes

A timeout on a synchronous `PATCH` is ambiguous. Fabric Manager may have
finished, may still be pushing config to devices, or may never have received
the request — the HTTP layer cannot tell which.

By default the SDK says so plainly: the timeout surfaces as a `TransportError`
and the caller decides. Nothing is retried, because retrying a write that may
have landed is its own hazard.

Where that decision is always "wait and see", the SDK can do it for you:

```go
client, _ := ones.InitializeWithCreds(url, user, pass,
    ones.WithSyncPatchTimeoutFallback(true, 5*time.Second, 2*time.Minute))
```

On timeout the tenant's `config_status` is then polled until it settles:

| Outcome | Result |
|---|---|
| `PASS` | success — the change landed after all |
| `FAIL` | `APIError` naming the terminal status |
| `BROWN_FIELD` / `SKIPPED` | success — deliberately not configured |
| still running at `MaxWait` | `SyncPatchTimeoutError` — outcome still unknown |

That last case is deliberately **not** a failure. The change may yet complete,
so the error says the outcome is unknown rather than inviting a blind retry.

`ones.IsTimeout(err)` reports whether any error was a timeout, at any depth,
without matching on strings.

## How requests are routed

Point the client at your ONES host and it handles the rest — but it is worth
knowing that not every call leaves by the same door.

Most Fabric Manager calls go through the ONES gateway under `/api/fm/`. The
gateway proxies that prefix through an explicit per-route allowlist, and as of
the 4.4.0 release 35 routes are missing from it — network policy, webhook
registration, IPAM configuration and address queries, tenant DHCP, rail
CIDRPools, and the shared-services and VF-interface routes.
A call to an unlisted route does not fail cleanly: it falls through to the UI's
single-page app and returns **HTML with HTTP 200**.

So the SDK sends those 34 straight to Fabric Manager, which serves them at the
root with no prefix. Everything else continues through the gateway.

```
auth + most FM calls   ->  https://<ones-host>/api/fm/...
the 34 listed routes   ->  https://<ones-host>:8443/...
```

The Fabric Manager base URL is derived from the one you pass — same host, port
8443, which is how the deployment publishes it — so the common case needs no
extra configuration. Override it where that does not hold:

```python
client = Client.initialize_with_creds(
    "https://ones-host", "user", "pass",
    fm_direct_url="https://fm-host:8443",      # optional
)
```

```go
client, _ := ones.InitializeWithCreds("https://ones-host", "user", "pass")
client.GetTransport().SetFMDirectURL("https://fm-host:8443")   // optional
```

Authentication is unchanged: the token comes from the gateway and Fabric
Manager accepts it, because both validate against the same secret.

**This split is temporary.** It exists because the 4.4.0 train had shipped when
the gap was found, so neither the gateway nor Fabric Manager could change. The
route list lives in one place in each SDK — `ones_gfx/_fm_direct.py` and
`ones_gfx/fmdirect.go`. Once a gateway release registers these routes, empty
that list and every call goes back through the gateway with no other change.

## Language notes

| | Python | Go |
|---|---|---|
| Endpoint coverage | 130 | 130 (+5 async variants = 135 functions) |
| Errors | exceptions | returned `error` values, typed taxonomy |
| Type checking | TypedDicts, checked by mypy/pyright | compile-time |
| Context/cancellation | — | takes `context.Context` (currently accepted but **not** yet honoured by the transport) |
| Import surface | `ones_gfx` + `ones_gfx.apis` | single `ones` package |

Per-SDK detail: [Python README](./ones-gfx-sdk-python/README.md) ·
[Go README](./ones-gfx-sdk-go/README.md) ·
[Go implementation status](./ones-gfx-sdk-go/IMPLEMENTATION_STATUS.md)

---

## Repository structure

```
ones-gfx-sdk/
├── README.md                    # This file
├── CONTRIBUTING.md
├── ones-gfx-sdk-python/
│   ├── ones_gfx/
│   │   ├── client.py            # Client class
│   │   ├── apis/                # 130 endpoint functions
│   │   └── _types.py            # TypedDicts
│   ├── index.py                 # Runnable walkthrough
│   └── pyproject.toml
└── ones-gfx-sdk-go/
    ├── ones/                    # Public package: client.go, api.go, types.go
    ├── ones_gfx/                # Internals: transport, auth, resources/
    ├── examples/index/          # Runnable walkthrough
    └── go.mod
```

---

## Versioning Policy

This repository follows [Semantic Versioning](https://semver.org/)
(`MAJOR.MINOR.PATCH`) and is maintained **independently** of AVIZ ONES
Spectrum-X platform releases — a version bump here does not imply a
corresponding change in the ONES platform version, and vice versa.

- **MAJOR:** breaking, backward-incompatible changes
- **MINOR:** new backward-compatible features
- **PATCH:** backward-compatible fixes

### Compatibility Matrix

| Go SDK | Python SDK | Supported ONES Version |
|--------|------------|------------------------|
| v1.0.0 | v1.0.0     | 4.2.1, 4.3             |

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup, coding
conventions, and the PR workflow.

## Support

- **Issues:** [Open an issue](https://github.com/aviznetworks/ones-gfx-sdk/issues/new/choose)
- **Security:** report vulnerabilities to security@aviznetworks.com

## License

Apache 2.0 — see [LICENSE](./LICENSE).

---

**Developed by AVIZ Networks** — https://aviznetworks.com
