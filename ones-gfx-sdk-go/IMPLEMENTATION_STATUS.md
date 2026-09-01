# Go SDK Implementation Status

The Go SDK mirrors the Python SDK's typed API surface: all 97 Fabric Manager
endpoints are exposed through typed resource methods (no `map[string]interface{}`
body builders, no `interface{}` returns except the documented open-payload leaves).

## Build and Test Commands

```bash
# From ones-gfx-sdk/ones-gfx-sdk-go/

# Build and vet everything
go build ./...
go vet ./...

# Format check
gofmt -l ones_gfx/ sdk/ examples/
```

## Package layout

- `ones_gfx/` — transport, auth, enums, errors, shared API types (`apitypes.go`).
- `ones_gfx/resources/` — one file per API domain, 97 endpoints total.
- `sdk/` — the top-level `Client` with one resource handle per domain.
- `examples/` — a runnable CLI tour (`go run examples/usage_examples.go`).

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

## Open-payload leaves (16)

The following returns fall back to `interface{}` / `map[string]interface{}`
because the Java service layer itself returns `Object`/`Map<String,Object>`
(see the Python SDK's TYPING_GAPS.md for the full analysis):

- bare `interface{}`: `ConfigMgmt.Get`, `ConfigMgmt.GetDiff`, `ConfigMgmt.Replace`
- `map[string]interface{}`: `Tenants.List/Get`, `Inventory.Ports/UFMHosts/Sync`,
  `NMXC.Inventory`, `Bootstrap.GetBatch`
- `[]interface{}`: `Devices.Versions/ImgmgmtStatus`, `Intents.Validation/Day1ConfigStatus/DerivationLogs`
- `[]map[string]interface{}`: `Bootstrap.ListBatches`
- `GetFilesResult.Files` (grouped map for "ALL" vs list for a typed filetype)

## Resource method index

| Resource | Methods |
|---|---|
| `Tenants` | Create/CreateAsync, List, Get, Delete/DeleteAsync, Update/UpdateAsync, ModifyAllocations/ModifyAllocationsAsync, AutoAllocate |
| `GPU` | AssignPorts, List, ByHost, TenantMappings, AllocationHistory, AvailableServers |
| `Fabrics` | Create, Update, Delete, List, Get, ListDTOs, DeviceIPs, DevicesByLayer, UpdateStatus |
| `Fabricsims` | Create, List, Get, Delete, UpdateStatus |
| `Inventory` | Add, Update, Edit, ListAll, ByFabric, Ports, UFMHosts, Sync |
| `Devices` | AddFacts, RemoveFacts, Versions, ImgmgmtStatus, EnableZTP, UpgradeNOS, Reboot, IsAlive, FMInventory, ValidateUFMCreds |
| `ConfigMgmt` | Get, GetDiff, Replace, Restore, Backup, FetchBackupFiles, ListToRestore, UploadDay1 |
| `Bootstrap` | Fill, Trigger, List, ListBatches, GetBatch, DeviceStages |
| `RMA` | Fill, Trigger, List, Status |
| `Files` | Upload, List, Delete |
| `System` | ControllerVersion, ControllerVersionInternal, Status, UploadStatus, StartStreaming, StopStreaming, UpdateRoleInfo, SetLogLevel, GetLogLevel |
| `NetOps` | Device, Fabric |
| `HostTenants` | Add, Delete, List, ListHosts, UpdateHosts |
| `NMXC` | UpdateDomains, ProbeDomains, Inventory, Reset, FactoryReset |
| `Intents` | UploadUIObject, GetUIObject, LastOrchestratedName, Validation, Day1ConfigStatus, DerivationLogs |
| `VPCPeering` | Create/CreateAsync, Delete |
| `Operations` | Get, WebhookStatus |
