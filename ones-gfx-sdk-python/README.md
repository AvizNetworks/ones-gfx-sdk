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

    # Create a tenant (synchronous — blocks until done)
    tenant = client.tenants.create(
        fabric_name="UI-CIT-fabric482",
        name="tenant2",
        description="Testing workloads",
        max_gpus_allowed=8,
    )
    print(tenant.id, tenant.vlan_id)

    # Allocate GPUs
    client.tenants.allocate_gpus(
        fabric_name="UI-CIT-fabric482",
        name="tenant2",
        servers=["hgx-su00-h00", "hgx-su00-h01"],
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
