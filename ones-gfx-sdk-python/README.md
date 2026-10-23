# ONES Spectrum-X Python SDK

Python client library for the AVIZ **ONES Spectrum-X** tenant management API
(ONES 4.2). Provides a small, typed surface for managing tenants and GPU
allocations across an NVIDIA Spectrum-X fabric.

## Status

**v1.0.0** — initial release. Covers tenant CRUD, GPU allocate/deallocate,
fabric and operation reads, and VPC peering. JWT-based auth with automatic
token refresh.

## Installation

For most users:

```bash
pip install ones-gfx-sdk
```

From a local checkout (editable, picks up your code edits live):

```bash
pip install -e .
```

Requires Python 3.9 or newer. The only runtime dependency is
[`requests`](https://requests.readthedocs.io/).

### Using the Makefile

A `Makefile` wraps the common workflows:

```bash
make help         # list all targets
make install      # editable install
make build        # build wheel + sdist into dist/
make run-example  # run examples/usage_example.py
make check        # quick import sanity check
make clean        # remove build artifacts and caches
```

Override the Python interpreter when needed:

```bash
make PYTHON=python3.11 install
```

### Building a wheel

```bash
make build
# or directly:
pip install build && python -m build
```

Outputs land in `dist/`:

- `ones_gfx_sdk-1.0.0-py3-none-any.whl` — the wheel partners install
- `ones_gfx_sdk-1.0.0.tar.gz` — the source distribution

Ship the `.whl` to partners; they install with:

```bash
pip install ones_gfx_sdk-1.0.0-py3-none-any.whl
```

### Running the example without installing

The bundled `examples/usage_example.py` adds the repo root to `sys.path`
itself, so it can run directly against the source tree without
`pip install`. Any of these work from a fresh checkout:

```bash
python examples/usage_example.py        # from repo root
python usage_example.py                 # from inside examples/
python -m examples.usage_example        # module style, from repo root
```

## Quickstart

```python
from ones_gfx import ONESClient, JWTAuth, OperationMode, UNLIMITED_GPUS

auth = JWTAuth(
    access_token="eyJhbGc...",
    refresh_token="eyJhbGc...",
    refresh_url="https://10.4.5.76:8089/refresh",
)

with ONESClient(base_url="https://10.4.5.76:8089", auth=auth) as client:
    # List fabrics
    for fabric in client.fabrics.list():
        print(fabric.fabric_name, fabric.default_storage_name)

    # Inventory sync — UFM enabled fabrics
    client.fabrics.inventory_sync("UFM-fabric")

    # Create a tenant (synchronous — blocks until done)
    tenant = client.tenants.create(
        fabric_name="UI-CIT-fabric482",
        name="tenant2",
        description="Testing workloads",
        max_gpus_allowed=8,
    )
    print(tenant.id, tenant.vlan_id)

    # Allocate GPUs (PATCH /tenants — all fabric types)
    client.tenants.allocate_gpus(
        fabric_name="UI-CIT-fabric482",
        name="tenant2",
        servers=["hgx-su00-h00", "hgx-su00-h01"],
    )

    # Assign specific ports — UFM/NMXC fabrics (POST /gpus, always sync)
    client.tenants.assign_ports(
        fabric_name="NMXC",
        tenant_name="tenant2",
        server_names=["su00-rack00-node00", "su00-rack01-node18"],
        gpu_ids=[1, 2, 3],
    )

    # Unassign ports
    client.tenants.unassign_ports(
        fabric_name="NMXC",
        tenant_name="tenant2",
        server_names=["su00-rack00-node00", "su00-rack01-node18"],
        gpu_ids=[1, 2, 3],
    )
```

## Operation modes

Mutating methods (`create`, `delete`, `allocate_gpus`, `deallocate_gpus`)
accept a `mode` keyword argument:

| Mode | Behavior |
|---|---|
| `OperationMode.SYNCHRONOUS` *(default)* | Blocks until the server finishes; returns the parsed result. |
| `OperationMode.ASYNC_POLL` | Returns immediately with an `Operation` containing a job ID. Caller polls `client.operations.get(op.id)`. |
| `OperationMode.ASYNC_WEBHOOK` | Returns immediately with an `Operation`. Server POSTs the final result to your `webhook_url`. |

### Async polling

```python
op = client.tenants.create(
    fabric_name="UI-CIT-fabric482",
    name="tenantH",
    description="Testing workloads",
    max_gpus_allowed=8,
    mode=OperationMode.ASYNC_POLL,
)
print(op.id)         # "TENANT_1776512073596_d4f8_..."
print(op.status)     # "PENDING"

# Caller controls the polling cadence
import time
while True:
    current = client.operations.get(op.id)
    print(f"{current.status} ({current.progress}%)")
    if current.is_done:
        break
    time.sleep(5)

if current.is_success:
    print("Tenant created:", current.result)
else:
    print("Failed:", current.error_message)
```

### Async with webhook

```python
op = client.tenants.create(
    fabric_name="UI-CIT-fabric482",
    name="tenant4",
    description="Testing workloads",
    max_gpus_allowed=8,
    mode=OperationMode.ASYNC_WEBHOOK,
    webhook_url="http://my-receiver:8000/hook",
    webhook_events=["tenant.create"],
)
print(op.id, op.webhook_registered)
```

The SDK does not implement the receiver — you build that on your side.
The server will POST the operation result to `webhook_url` on completion.

## IPAM and rail CIDRPools

Seven functions in `ones_gfx.apis` cover per-fabric IPAM and the rail CIDRPool
manifests (typed responses are `TypedDict`s in `ones_gfx._types`).

```python
from ones_gfx import Client
from ones_gfx.apis import (
    get_rail_cidr_pools, get_rail_cidr_pool, update_ipam_config, get_ipam_config,
    update_tenant, get_tenant_compute_ip, get_tenant_ips, get_server_ips,
)

client = Client.initialize_with_token("https://10.4.5.76:8089", token)

# Rail CIDRPool manifests — raw YAML by default; accept= picks the representation:
# "application/yaml" | "application/x-yaml" | "text/plain" (all a str) or
# "application/json" ({"fabricName", "rails", "yaml"}).
manifest = get_rail_cidr_pools(client, "yml")
r0 = get_rail_cidr_pool(client, "yml", "r0")

# Per-fabric IPAM switches (externally managed fabrics only)
update_ipam_config(client, "ns", enableIpam=True, enableIpAssign=True)
get_ipam_config(client, "ns")
update_ipam_config(client, "ns", enableIpam=False, enableIpAssign=False, force=True)

# Adding a server to a tenant allocates (and, with enableIpAssign, assigns) its IP
res = update_tenant(client, "ns", "gorax", operation="ADD",
                    servers=[{"serverName": "su00-node00"}])
res["servers"][0]["ipAddress"]      # "192.3.0.2/32"
res["partial"]                      # True when any server failed (HTTP 207)

get_tenant_compute_ip(client, "ns", "gorax", "su00-node00")   # one allocation
get_tenant_ips(client, "ns", "gorax", page=0, size=50, status="ASSIGNED")
get_server_ips(client, "ns", "su00-node00")                   # across tenants
```

Behaviour worth knowing:

- Arguments are validated before any request: blank names, non-bool switches,
  `enableIpAssign` without `enableIpam`, `page < 0`, `size < 1`, unknown
  `status` / `accept` raise `ValueError`. Path parameters are URL-quoted.
- A partial IPAM result is HTTP **207**, which is a success status — no
  exception; check `res["partial"]`.
- Failures raise `requests.HTTPError`. The server's error code is in
  `exc.response.json()["error"]["code"]` (e.g. `IPAM_HAS_ACTIVE_POOLS`,
  `FABRIC_NOT_EXTERNALLY_MANAGED`, `IPAM_NOT_CONFIGURED`, `INVALID_RAIL`). For the
  YAML/text rail representations the error body is plain text
  `"<CODE>: <detail>"` (`exc.response.text`).
- `create_tenant(..., tenantDhcpCidrPool="192.3.0.0/17")` sets the tenant's IPAM
  pool; the server requires it when IPAM is on and the fabric has no tenant
  compute pool to derive one from (`IPAM_CIDR_REQUIRED`). It is validated before
  the request — IPv4, prefix /8–/30, surrounding whitespace stripped — and raises
  `ValueError` otherwise. The Terraform provider (`fabricapi_tenant.tenant_dhcp_cidr_pool`)
  applies the same rule, the same 207 handling and surfaces the same error codes.

Run the tests (unit + a fake Fabric Manager over real HTTP) with `make test`.

## Partial GPU allocation

Map or unmap specific GPU indices to a tenant on a shared server. This
complements `allocate_gpus` (whole-server attachment) with fine-grained
per-GPU control via `POST /fabrics/{fabric}/tenants/{tenant}/gpuAllocations`.

```python
from ones_gfx import ONESClient, JWTAuth

auth = JWTAuth(access_token=..., refresh_token=..., refresh_url=...)
with ONESClient(base_url=..., auth=auth) as client:
    # Map G0-G3 on hgx-su00-h00 (server index "0") to tenant-A
    result = client.fabrics.modify_gpu_allocations(
        fabric_name="my-fabric",
        tenant_name="tenant-A",
        operation="ADD",           # or "DELETE"
        suid={
            "0": {
                "hgx-su00-h00": {"gpus": ["G0", "G1", "G2", "G3"]}
            }
        },
    )
    print(result["status"], result["message"])
```

**Notes:**
- The server must already be attached to the tenant via `allocate_gpus` before mapping individual GPUs.
- Remove per-GPU mappings with `operation="DELETE"` before calling `deallocate_gpus` to detach the server.
- This endpoint is valid only for externally managed fabrics; on ONES-controlled fabrics it returns `409 FABRIC_NOT_EXTERNALLY_MANAGED` (use the tenant-update flow instead).
- Multiple tenants can share the same physical server (e.g. G0–G3 → tenant-A, G4–G7 → tenant-B) when the server is attached with `shared=True`.

---

## Unlimited GPU quota

Use the `UNLIMITED_GPUS` constant (or pass `-1` directly) to mean
"no quota cap":

```python
from ones_gfx import UNLIMITED_GPUS

client.tenants.create(
    fabric_name="...", name="...", description="...",
    max_gpus_allowed=UNLIMITED_GPUS,
)
```

`max_gpus_allowed` must be `-1` or a positive integer; `0` and other
negative values are rejected client-side with a clear error.

## Token refresh behavior

`JWTAuth` handles token rotation automatically:

- **Proactive**: before each API call, if the access token expires within
  10 seconds, the SDK refreshes first.
- **Reactive**: if a request returns 401 anyway (clock skew, server-side
  revocation), the SDK refreshes and retries the request once.

Both the access token and the refresh token are rotated on every refresh
(per the ONES contract). To persist the rotated tokens, pass an
`on_token_refresh` callback:

```python
def save_tokens(access, refresh, expires_in):
    my_secret_store.write({"access": access, "refresh": refresh})

auth = JWTAuth(
    access_token=...,
    refresh_token=...,
    refresh_url="https://10.4.5.76:8089/refresh",
    on_token_refresh=save_tokens,
)
```

If the refresh token itself is rejected, `AuthenticationError` is raised
and you must obtain new tokens out-of-band and construct a new client.

## Exception hierarchy

```
ONESError
├── TransportError              # network / TLS / timeout
├── APIError                    # any server-returned error
│   ├── AuthenticationError     # 401, or refresh failure
│   ├── BadRequestError         # 400 (invalid input)
│   ├── NotFoundError           # 404
│   ├── ConflictError           # 409 (duplicate, GPUs still allocated)
│   └── ServerError             # 5xx
└── OperationFailedError        # async job ended with status FAILURE
```

Catch `ONESError` to handle anything from the SDK; catch a specific
subclass for finer-grained control.

## TLS

Use `verify_tls` to control certificate verification:

```python
ONESClient(base_url=..., auth=..., verify_tls=True)              # default
ONESClient(base_url=..., auth=..., verify_tls="/etc/ssl/ca.pem") # custom CA
ONESClient(base_url=..., auth=..., verify_tls=False)             # dev only!
```

## Limitations / planned

- v1.0.0 does not include a login/logout flow — partners supply tokens
  obtained out-of-band. Login support is planned for v1.1.0
- v1.0.0 does not retry on transport errors beyond the single 401 retry.
  Wrap calls with your own retry policy if you need exponential backoff.
