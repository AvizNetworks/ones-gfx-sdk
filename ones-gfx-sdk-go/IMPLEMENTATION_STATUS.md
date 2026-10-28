# Go SDK Implementation Status

The Go SDK mirrors the Python SDK's typed API surface: the 20 Fabric Manager
endpoints documented in *ONES Spectrum-X API Operations v3.0* are exposed as
flat functions (25 including the `...Async` variants) backed by typed resource
methods. See [`../FM_API_LIST.md`](../FM_API_LIST.md) for the full list.

Not yet included: the guide's three HBN VF-interface endpoints
(`GET …/vf-interfaces`, `POST`/`DELETE …/vf-interfaces/{vfId}/assign`).

## Build and Test Commands

```bash
# From ones-gfx-sdk/ones-gfx-sdk-go/

# Build and vet everything
go build ./...
go vet ./...

# Format check
gofmt -l ones/ ones_gfx/ examples/
```

## Package layout

- `ones/` — the public package: `client.go` (the `Client`), `api.go` (one flat
  function per endpoint) and `types.go` (type aliases, constants and options, so
  `ones` is the only import a caller needs).
- `ones_gfx/` — transport, auth, enums, errors, shared API types (`apitypes.go`).
- `ones_gfx/resources/` — one file per API domain.
- `examples/` — `index/` (login + list fabrics) and `usage_examples.go`, a CLI
  tour of the tenant lifecycle (`go run ./examples`).

## Typing conventions

- **Path params** are positional; **body/query fields** are individual typed
  args. Required fields use value types, optional fields use pointers
  (`nil` = omit), matching the Python SDK's `| None = None` convention.
- **Async-capable endpoints** have separate Sync/Async methods (e.g.
  `Tenants.Create` / `Tenants.CreateAsync`); async methods return
  `*ones_gfx.OperationAccepted`.
- **Types used by more than one resource file** live in `ones_gfx/apitypes.go`;
  single-use types are defined next to their method.
- **JSON tags** match the Java field names exactly.

## Typed request/response transport

`ones_gfx.Call[T]` and `ones_gfx.CallMultipart[T]` decode responses directly
into a typed result, unwrapping the `{"data": ...}` envelope exactly as the
Python client does. The legacy `Get/Post/Patch/Delete` methods remain for
back-compat but are no longer used by the resources.

## Open-payload leaves (6)

The following returns fall back to `map[string]interface{}` because the Java
service layer itself returns `Map<String,Object>`:

- `Tenants.List`, `Tenants.Get`
- `Inventory.Ports`, `Inventory.UFMHosts`, `Inventory.Sync`
- `NMXC.Inventory`

## Resource method index

| Resource | Methods |
|---|---|
| `Fabrics` | ListDTOs |
| `Tenants` | Create/CreateAsync, List, Get, Delete/DeleteAsync, Update/UpdateAsync, ModifyAllocations/ModifyAllocationsAsync |
| `GPU` | AssignPorts, AvailableServers |
| `VPCPeering` | Create/CreateAsync, Delete |
| `Devices` | ValidateUFMCreds |
| `Inventory` | Ports, UFMHosts, Sync |
| `NMXC` | UpdateDomains, ProbeDomains, Inventory, Reset, FactoryReset |
| `Operations` | Get |
