# ONES GFX SDK

Multi-language SDK for **AVIZ ONES Spectrum-X** tenant management API (v4.2.1).
Python and Go implementations expose the **same 97 Fabric Manager endpoints**
through the same `Client` shape, so code translates almost line-for-line between
them.

Covers tenant lifecycle, GPU allocation and per-port assignment, fabric
discovery and configuration, inventory, bootstrap/ZTP, RMA, config
backup/restore, NMX-C and UFM integration, VPC peering, and async operation
tracking.

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
endpoints (102 functions vs Python's 97; Python selects async with the `prefer`
argument on the same function).

```python
from ones_gfx.apis import get_all_fabrics
fabrics = get_all_fabrics(client)
```

```go
fabrics, err := ones.GetAllFabrics(ctx, client)
```

Names differ only in case convention: `get_all_fabrics` ↔ `GetAllFabrics`,
`add_fabric_data` ↔ `CreateFabric`, `modify_gpu_allocations` ↔
`ModifyGpuAllocations`. Optional fields are keyword args defaulting to `None` in
Python, and pointers (`ones.Ptr(v)`, `nil` to omit) in Go.

### Errors

| Condition | Python | Go |
|---|---|---|
| Call with no token | raises `NotAuthenticatedError` | `ErrNotAuthenticated` |
| Nil/missing client | n/a | `ErrNilClient` |
| Session method on a non-token client | n/a | `ErrNoTokenAuth` |
| HTTP 4xx/5xx | `requests.HTTPError` | typed: `*AuthenticationError`, `*NotFoundError`, `*ConflictError`, `*ServerError`, … (use `errors.As`) |

---

## Example — Python

```python
from ones_gfx import Client
from ones_gfx.apis import add_fabric_data, get_all_fabrics

client = Client.initialize_with_creds("https://10.4.5.76:8089", "admin", "secret")
print("authenticated:", client.is_authenticated())

for f in get_all_fabrics(client):
    print(f"  {f.get('id')} {f.get('name')} type={f.get('type')}")

print(add_fabric_data(
    client,
    name="gpu-fabric-1",
    type="DNO ASN",
    description="Primary GPU fabric",
    status="draft",
))

client.close()
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

    fabrics, err := ones.GetAllFabrics(ctx, client)
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range fabrics {
        fmt.Printf("  %d %s type=%s\n", f.ID, f.Name, f.Type)
    }

    msg, err := ones.CreateFabric(ctx, client, "gpu-fabric-1", &ones.FabricCreateArgs{
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

`ones` is the only import you need — every type, constant and option is
re-exported there (`ones.FabricCreateArgs`, `ones.GpuActionAdd`, `ones.Ptr`, …).

Runnable versions of both live at
[`ones-gfx-sdk-python/index.py`](./ones-gfx-sdk-python/index.py) and
[`ones-gfx-sdk-go/examples/index`](./ones-gfx-sdk-go/examples/index):

```bash
cd ones-gfx-sdk-python && python3 index.py
cd ones-gfx-sdk-go     && go run ./examples/index
```

---

## Language notes

| | Python | Go |
|---|---|---|
| Endpoint coverage | 97 | 97 (+5 async variants = 102 functions) |
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
│   │   ├── apis/                # 97 endpoint functions
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
| v1.0.0 | v1.0.0     | 4.2.1                  |

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
