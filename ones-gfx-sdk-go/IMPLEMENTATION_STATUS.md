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

## Typed responses

Every resource method returns a typed result. The responses FM builds as a
`Map<String,Object>` are modelled from the Java source:

- `Tenants.List`, `Tenants.Get` → `[]ones_gfx.Tenant`, `*ones_gfx.Tenant`
- `Inventory.Ports` → `*InventoryPortsResult`
- `Inventory.UFMHosts` → `*UFMHostsResult`
- `Inventory.Sync` → `*InventorySyncResult`
- `NMXC.Inventory` → `*NmxcInventory`

Each keeps top-level keys it does not model in an `Extra` field, and has a
`...Raw` twin (`ListRaw`, `GetRaw`, `PortsRaw`, `UFMHostsRaw`, `SyncRaw`,
`InventoryRaw`) that returns the response as `map[string]interface{}`.

## Client-side validation

`resources/validation.go` checks every method's arguments before a request is
built and returns `*ones_gfx.ValidationError` (an `ONESError`). See the
README's "Client-side validation" section for the rules.

## Resource method index

| Resource | Methods |
|---|---|
| `Fabrics` | ListDTOs |
| `Tenants` | Create/CreateAsync, CreateWith/CreateAsyncWith, List/ListRaw, Get/GetRaw, Delete/DeleteAsync, Update/UpdateAsync, UpdateWith/UpdateAsyncWith, ModifyAllocations/ModifyAllocationsAsync, ModifyAllocationsWith/ModifyAllocationsAsyncWith |
| `GPU` | AssignPorts, AssignPortsWith, AvailableServers |
| `VPCPeering` | Create/CreateAsync, Delete |
| `Devices` | ValidateUFMCreds |
| `Inventory` | Ports/PortsRaw, UFMHosts/UFMHostsRaw, Sync/SyncRaw |
| `NMXC` | UpdateDomains, ProbeDomains, Inventory/InventoryRaw, Reset, FactoryReset |
| `Operations` | Get |
| `VFs` | List, Assign, Unassign |
