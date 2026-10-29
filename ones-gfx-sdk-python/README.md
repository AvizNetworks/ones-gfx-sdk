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

### Running the unit tests

The tests under `tests/` stub the HTTP session, so they need no ONES instance:

```bash
make install-dev   # pytest + ruff
make test          # or: python -m pytest tests
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

## VPC peering

```python
import requests
from ones_gfx import Client
from ones_gfx.apis import create_vpc_peering, delete_vpc_peering

client = Client.initialize_with_creds("https://10.4.5.76:8089", "admin", "secret")

create_vpc_peering(
    client, "Fab101",
    name="tenant1-storage-route-leak",
    vpcname="tenant1-Fab101-north-south",
    peervpcname="Fab101-Storage",
)

# Delete: synchronous only. `name` addresses the peering, and FM rejects the
# call unless vpcname/peervpcname match the pair recorded for it.
try:
    resp = delete_vpc_peering(
        client, "Fab101",
        name="tenant1-storage-route-leak",
        vpcname="tenant1-Fab101-north-south",
        peervpcname="Fab101-Storage",
    )
    print(resp["message"])
except requests.HTTPError as e:
    # 404 VPC_PEERING_NOT_FOUND / FABRIC_NOT_FOUND, 400 VPC_PEERING_PAIR_MISMATCH,
    # 409 VPC_PEERING_OWNED_BY_POLICY (delete the network policy instead),
    # 401/403 auth, 500 VPC_PEERING_DELETE_FAILED
    print(e.response.status_code, e.response.json().get("error", {}).get("code"))
```

Reading peerings (FM with GFX-5658). The gateway does not proxy these GETs
yet, so the SDK sends them straight to Fabric Manager (`DIRECT_FM_GET_PATHS`):

```python
from ones_gfx.apis import get_vpc_peering, list_tenant_vpc_peerings, list_vpc_peerings

page = list_vpc_peerings(client, "Fab101", page=0, size=100)    # size defaults to 50, max 500
for p in page["peerings"]:
    print(p["name"], p["vpcs"], p["origin"])                       # origin: VPC_PEERING | NETWORK_POLICY
print(page["pagination"]["totalElements"], "total")

peering = get_vpc_peering(client, "Fab101", "tenant1-storage-route-leak")

for p in list_tenant_vpc_peerings(client, "Fab101", "tenant1")["peerings"]:
    print(p["name"], p["vpc"], "<->", p["peerVpc"])               # tenant's own side first
```

From the CLI example:

```bash
python examples/usage_examples.py -u ... -U ... -P ... -f Fab101 \
    --action vpcpeering-delete --peering-name tenant1-storage-route-leak \
    --vpc-name tenant1-Fab101-north-south --peer-vpc-name Fab101-Storage
```

## Host N-S interfaces

`GET /fabrics/{fabricName}` on the ONES gateway returns the fabric's topology
YAML (the file it was deployed from): every device, cable, credential and
setting. These calls fetch it and keep only what most automation needs: **each
host's north-south (N-S) interfaces**, and which leaf port each one is cabled to.

A host link is N-S when the leaf at its far end matches a name pattern. The
defaults cover the three N-S designs ONES generates:

| Pattern | N-S design |
|---|---|
| `^leaf-cn` | converged user + storage network |
| `^leaf-fn` | front-end network |
| `^leaf-sn` | dedicated storage network |

E-W leaves always carry an `su` token (`leaf-su00-r0`, `leaf-p0-su00-r1`) and
never match. Pass your own patterns, one or several, for a reference
architecture that names its leaves differently.

The Host, Leaf and DPU sections are all read, so a cable listed on only one side
is still found, and duplicates are reported once. On HBN fabrics the path goes
host VF → DPU → leaf. A VF on DPU interface `pf<N>…` leaves through uplink
`p<N>`, and the entry also names the DPU (`via`) and its VF interface
(`viaInterface`, the `vfId` that VF-assign takes). Links to external routers
(`H<n>`) are ignored, and **credentials are never copied into the output**.

```python
from ones_gfx import Client, extract_ns_interfaces, ns_interfaces_to_yaml
from ones_gfx.apis import get_fabric_ns_interfaces, get_fabric_yaml

client = Client.initialize_with_creds("https://10.4.5.76", "admin", "secret")

result = get_fabric_ns_interfaces(client, "Test")                     # default ^leaf-cn, ^leaf-fn, ^leaf-sn
result = get_fabric_ns_interfaces(client, "Test", leaf_patterns=["^leaf-cn"])         # converged only
print(ns_interfaces_to_yaml(result))

# Or from a topology file you already downloaded:
result = extract_ns_interfaces(open("Test_yaml.yaml").read(), fabric_name="Test")
```

Output:

```yaml
fabric: Test
leafPatterns:
  - ^leaf-cn
  - ^leaf-fn
  - ^leaf-sn
hostCount: 128
interfaceCount: 256
hosts:
  - host: hgx-su00-h00
    nsInterfaces:
      - interface: enp8s0f1np0
        leaf: leaf-cn-00
        leafPort: swp9
      - interface: enp9s0f1np0
        leaf: leaf-cn-01
        leafPort: swp9
  ...
```

On an HBN fabric, a path through a DPU looks like this:

```yaml
  - host: hgx-su00-h00
    nsInterfaces:
      - interface: vf0
        via: hgx-su00-h00-dpu-00
        viaInterface: pf0vf0_if
        leaf: leaf-cn-00
        leafPort: swp9
```

`get_fabric_ns_interfaces` raises `requests.HTTPError` when the gateway has no
topology for the fabric (404), and `ones_gfx.TopologyError` when the response is
not a topology. That usually means an HTML page, from a gateway without the
route.

Command line (installed with the package):

```bash
ones-gfx-ns-interfaces Test_yaml.yaml --fabric Test -o ns.yaml           # from a file
ones-gfx-ns-interfaces --url https://<ones> -U admin -P '<pw>' --fabric Test \
    --leaf-pattern '^leaf-cn'                                         # fetch; converged only
```

This adds **PyYAML** (`pyyaml>=5.4`) as a dependency.

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
