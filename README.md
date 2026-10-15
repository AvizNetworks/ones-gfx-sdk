# ONES GFX SDK

Multi-language SDK for **AVIZ ONES Spectrum-X** tenant management API (v4.2.1).
Python and Go implementations expose the **same 97 Fabric Manager endpoints**
through the same `Client` shape, so code translates almost line-for-line between
them.

Covers tenant lifecycle, GPU allocation and per-port assignment, fabric
discovery and configuration, inventory, bootstrap/ZTP, RMA, config
backup/restore, NMX-C and UFM integration, VPC peering, and async operation
tracking — plus fabric creation and deployment through the ONES UI config
service (see [Fabric creation and deployment](#fabric-creation-and-deployment)).

---

## Install

Neither SDK is published to a registry — install both from this checkout.

**Python** (3.9+) — editable install from the local directory:

```bash
pip install -e ./ones-gfx-sdk-python
```

**Go** (1.19+; one dependency, `gopkg.in/yaml.v3`, for fabric deploy) — the module lives in this repo, so point at
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

## Fabric creation and deployment

Creating and deploying a fabric goes through the **ONES UI config service**
(`/api/config/...` and `/fabrics/{name}` on the same host), the same calls the
ONES UI and `terraform-provider-fabricapi` make. This is separate from the
Fabric Manager `add_fabric_data` / `CreateFabric`, which only stores a fabric
record without generating its topology.

| Step | Python | Go | Endpoint |
|---|---|---|---|
| Create (Draft; generates topology, YAML, inventory) | `create_fabric_config` | `CreateFabricConfig` | `POST /api/config/addFabricData` |
| Generated YAML | `get_fabric_yaml` | `GetFabricYaml` | `GET /fabrics/{name}` |
| **Deploy** (runs the six steps below) | `deploy_fabric` | `DeployFabric` | |
| 1. Record device IPs and credentials | `upload_fabric_device_ips` | `UploadFabricDeviceIPs` | `POST /api/config/uploadip` |
| 2. SSH-check switches | `validate_fabric_switches` | `ValidateFabricSwitches` | `POST /api/config/validateswitch` |
| 3. SSH-check hosts/DPUs | `validate_fabric_servers` | `ValidateFabricServers` | `POST /api/config/validateserver` |
| 4. Update FM inventory | `update_fabric_inventory` | `UpdateFabricInventory` | `POST /api/config/updateinventory` |
| 5. Push the config (YAML with credentials filled in) | `push_fabric_config` | `PushFabricConfig` | `POST /api/config` |
| 6. Mark Deployed | `update_fabric_config_status` | `UpdateFabricConfigStatus` | `POST /api/config/updatefabricstatus` |
| Delete (404 is not an error) | `delete_fabric_config` | `DeleteFabricConfig` | `DELETE /api/config/deletefabricbyname/{name}` |

Deploy needs every device in the fabric, in the `devices.json` format used by
the terraform provider — the hostnames are the ones in the generated YAML
(`spine-sNN`, `leaf-suNN-rN`, `hgx-suNN-hNN`):

```json
[{"hostname": "spine-s00",    "ip": "10.0.0.1", "username": "cumulus", "password": "...", "device_role": "spine"},
 {"hostname": "leaf-su00-r0", "ip": "10.0.0.2", "username": "cumulus", "password": "...", "device_role": "leaf"},
 {"hostname": "hgx-su00-h00", "ip": "10.0.0.3", "username": "root",    "password": "...", "device_role": "host"}]
```

`device_type` (`switch`, `server`, `dpu`, `host`) can stand in for
`device_role`; `apply_config: false` leaves a device unconfigured and skips its
switch validation. If any device fails validation, deploy stops before pushing
anything (`FabricValidationError` in both SDKs).

```python
import json
from ones_gfx import Client
from ones_gfx.apis import create_fabric_config, deploy_fabric

client = Client.initialize_with_creds("https://localhost", "superadmin", "secret")
create_fabric_config(client, name="gpu-fabric-1", type="Aviz RA 1.0",
                     hostsPerSu="{0:2}", enableEW=True, startingSubnetGpu="192")
result = deploy_fabric(client, "gpu-fabric-1", json.load(open("devices.json")))
```

```go
var devices []ones.FabricDevice
raw, _ := os.ReadFile("devices.json")
_ = json.Unmarshal(raw, &devices)

_, err := ones.CreateFabricConfig(ctx, client, "gpu-fabric-1", ones.FabricConfigArgs{
    Type: "Aviz RA 1.0", HostsPerSU: "{0:2}", EnableEW: true, StartingSubnetGPU: "192",
})
result, err := ones.DeployFabric(ctx, client, "gpu-fabric-1", devices, nil)
```

Things to know:

- **Deploy pushes real configuration** to the switches and servers listed.
- **There is no undeploy.** Deleting a Deployed fabric removes the record, not
  the configuration already on the devices.
- Delete leaves the fabric's inventory rows behind, so recreating a fabric with
  the same device IPs can produce duplicate inventory entries.
- `nsOperatingSystem` / `NSOperatingSystem` (`cumulus` or `sonic`) must match
  the real front-end and storage switches, or switch validation fails.
- Deploy fills in each device's IP and credentials in the generated YAML and
  sends it with its keys in the original order.

---

## Language notes

| | Python | Go |
|---|---|---|
| Endpoint coverage | 97 | 97 (+5 async variants = 102 functions) |
| Fabric create/deploy (config service) | 10 functions | 10 functions |
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
