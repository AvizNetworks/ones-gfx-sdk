"""
ONES Spectrum-X SDK examples — a runnable tour of the Python SDK.

Everything is driven from the command line; no credentials are baked in.

Usage
-----
Read-only tour (fabrics, tenants, available servers):

    python examples/usage_examples.py \
        -u https://localhost:3002 -U superadmin -P 'Admin@123456' \
        -f my-fabric --action read-only --insecure

Full tenant lifecycle in each of the three execution modes:

    # blocks until the server finishes each step
    python examples/usage_examples.py -u ... -U ... -P ... -f my-fabric \
        --action lifecycle --mode sync --insecure

    # returns an operationId per step, then polls GET /operations/{id}
    python examples/usage_examples.py -u ... -U ... -P ... -f my-fabric \
        --action lifecycle --mode async-poll --insecure

    # returns an operationId, server POSTs the result to --webhook-url
    python examples/usage_examples.py -u ... -U ... -P ... -f my-fabric \
        --action lifecycle --mode async-webhook \
        --webhook-url http://my-host:8000/hook --insecure

Single actions (each honours --mode where the endpoint supports it):

    --action create           create one tenant
    --action delete           delete one tenant
    --action allocate         attach servers   (--servers hgx-su00-h00,hgx-su00-h01)
    --action deallocate       detach servers
    --action assign-ports     assign GPU ports (--servers ... --gpu-ids 1,2,3)
    --action unassign-ports   unassign GPU ports
    --action gpu-allocations  map/unmap individual GPUs (--gpus G0,G1 --su-id 0)
    --action inventory-sync   force a UFM inventory sync
    --action vpcpeering       create a VPC peering
    --action error-handling   deliberately trigger and catch SDK errors

    python examples/usage_examples.py -u ... -U ... -P ... -f my-fabric \
        --action allocate --tenant-name demo --servers hgx-su00-h00 --insecure

All flags:

    --url/-u           ONES base URL                      (required)
    --username/-U      login username                     (required)
    --password/-P      login password                     (required)
    --fabric/-f        fabric name                        (required)
    --action           which scenario to run              (default: lifecycle)
    --mode             sync | async-poll | async-webhook  (default: sync)
    --tenant-name      override the tenant name
    --max-gpus         maxGpusAllowed on create           (default: 8)
    --shared           create the tenant as shared
    --servers          comma-separated server hostnames
    --gpu-ids          comma-separated GPU port ids
    --gpus             comma-separated GPU names (G0,G1,...)
    --su-id            SU index for --action gpu-allocations  (default: 0)
    --peering-name     peering name for --action vpcpeering
    --vpc-name         VPC name for --action vpcpeering
    --peer-vpc-name    peer VPC name for --action vpcpeering
    --webhook-url      receiver URL, required for --mode async-webhook
    --poll-interval    seconds between polls              (default: 5)
    --poll-attempts    max poll attempts                  (default: 60)
    --insecure         skip TLS certificate verification
"""

from __future__ import annotations

import argparse
import logging
import os
import sys
import time
from typing import Any

# Run straight from a checkout without installing.
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from ones_gfx import Client, NotAuthenticatedError
from ones_gfx.apis import (
    assign_gpu_ports,
    create_tenant,
    create_vpc_peering,
    delete_tenant,
    get_all_fabrics,
    get_available_servers,
    get_operation,
    get_operation_webhook_status,
    get_tenant,
    inventory_sync,
    list_tenants,
    modify_gpu_allocations,
    update_tenant,
)

MODES = ("sync", "async-poll", "async-webhook")


# ---------------------------------------------------------------------------
# Mode plumbing
#
# `mode` only ever becomes a Prefer header plus (for webhooks) three body
# fields. Sync returns the finished result; both async modes return an
# OperationAccepted carrying an operationId.
# ---------------------------------------------------------------------------


def mode_kwargs(mode: str, webhook_url: str | None, events: list[str]) -> dict[str, Any]:
    """Translate a --mode value into the kwargs the api functions expect."""
    if mode == "sync":
        return {}
    kwargs: dict[str, Any] = {"prefer": "respond-async"}
    if mode == "async-webhook":
        kwargs.update(
            enableWebhook=True, webhookUrl=webhook_url, webhookEvents=events
        )
    return kwargs


def operation_id(result: Any) -> str | None:
    """Return the operationId when the server answered asynchronously."""
    if isinstance(result, dict):
        return result.get("operationId")
    return None


def poll_operation(
    client: Client, op_id: str, *, label: str, interval_s: int, attempts: int
) -> bool:
    """Poll GET /operations/{id} until it reaches a terminal state."""
    print(f"  -> polling {label} operation: {op_id}")
    for attempt in range(1, attempts + 1):
        current = get_operation(client, op_id)
        status = current.get("status")
        if status in {"PENDING", "RUNNING"}:
            print(".", end="", flush=True)
        else:
            print(f"\n  [poll {attempt}] status={status}")
        if status in {"SUCCESS", "FAILURE"}:
            if status == "SUCCESS":
                return True
            print(f"  -> FAILED: {current.get('errorMessage')}")
            return False
        time.sleep(interval_s)
    print(f"\n  -> {label} still not done after {attempts} polls; giving up")
    return False


def poll_webhook(
    client: Client, op_id: str, *, label: str, interval_s: int, attempts: int
) -> bool:
    """Poll GET /operations/{id}/webhook-status until delivery settles."""
    print(f"  -> polling {label} webhook delivery: {op_id}")
    for attempt in range(1, attempts + 1):
        current = get_operation_webhook_status(client, op_id)
        status = current.get("deliveryStatus")
        print(f"\n  [poll {attempt}] webhook deliveryStatus={status}")
        if status and status.upper() not in {"PENDING", "RETRYING"}:
            return status.upper() in {"DELIVERED", "SUCCESS"}
        time.sleep(interval_s)
    print(f"\n  -> {label} not delivered after {attempts} polls; giving up")
    return False


def settle(client: Client, result: Any, *, label: str, args: argparse.Namespace) -> bool:
    """Resolve a call result: print it if sync, poll it if async."""
    op_id = operation_id(result)
    if op_id is None:
        print(f"  -> {label}: {result}")
        return True
    print(f"  -> operation id: {op_id}")
    if args.mode == "async-webhook":
        if isinstance(result, dict):
            print(f"  -> webhook registered: {result.get('webhookRegistered')}")
        poller = poll_webhook
    else:
        poller = poll_operation
    return poller(
        client,
        op_id,
        label=label,
        interval_s=args.poll_interval,
        attempts=args.poll_attempts,
    )


# ---------------------------------------------------------------------------
# Scenarios
# ---------------------------------------------------------------------------


def scenario_read_only(client: Client, args: argparse.Namespace) -> None:
    print("\n--- Scenario: read-only ---")

    fabrics = get_all_fabrics(client)
    print(f"Fabrics ({len(fabrics)}):")
    for f in fabrics:
        print(f"  - id={f.get('id')} name={f.get('name')!r} type={f.get('type')!r}")

    print(f"\nTenants in {args.fabric!r}:")
    print(f"  {list_tenants(client, args.fabric)}")

    print(f"\nAvailable servers in {args.fabric!r}:")
    print(f"  {get_available_servers(client, args.fabric).get('availableGPUs')}")


def scenario_create(client: Client, args: argparse.Namespace) -> bool:
    name = args.tenant_name or default_tenant_name(args.mode)
    print(f"Creating tenant {name!r} ({args.mode})...")
    result = create_tenant(
        client,
        args.fabric,
        tenantName=name,
        description=f"Created by SDK example ({args.mode} mode)",
        maxGpusAllowed=args.max_gpus,
        shared=args.shared,
        **mode_kwargs(args.mode, args.webhook_url, ["tenant.create"]),
    )
    return settle(client, result, label="create", args=args)


def scenario_servers(client: Client, args: argparse.Namespace, operation: str) -> bool:
    """Attach (ADD) or detach (DELETE) whole servers on a tenant."""
    name = args.tenant_name or default_tenant_name(args.mode)
    verb = "Allocating" if operation == "ADD" else "Deallocating"
    if not args.servers:
        print(f"  -> no --servers given; skipping {verb.lower()}")
        return True
    print(f"{verb} {args.servers} on {name!r} ({args.mode})...")
    result = update_tenant(
        client,
        args.fabric,
        name,
        servers=[{"serverName": s, "shared": args.shared or None} for s in args.servers],
        operation=operation,
        **mode_kwargs(args.mode, args.webhook_url, ["tenant.update"]),
    )
    return settle(client, result, label=operation.lower(), args=args)


def scenario_ports(client: Client, args: argparse.Namespace, operation: str) -> bool:
    """Assign/unassign specific GPU ports — UFM / NMXC fabrics, always sync."""
    name = args.tenant_name or default_tenant_name(args.mode)
    if not args.servers:
        print("  -> no --servers given; skipping port assignment")
        return True
    print(f"{operation} ports on {name!r} servers={args.servers} gpuIds={args.gpu_ids}...")
    result = assign_gpu_ports(
        client,
        args.fabric,
        name,
        operation=operation,
        serverNames=args.servers,
        gpuIds=args.gpu_ids or None,
    )
    print(f"  -> {result}")
    return bool(result.get("success", True))


def scenario_gpu_allocations(client: Client, args: argparse.Namespace) -> bool:
    """Map or unmap individual GPUs on a shared server."""
    name = args.tenant_name or default_tenant_name(args.mode)
    if not args.servers or not args.gpus:
        print("  -> --servers and --gpus are both required for gpu-allocations")
        return False
    suid = {str(args.su_id): {args.servers[0]: args.gpus}}
    print(f"Mapping {args.gpus} on {args.servers[0]} to {name!r} ({args.mode})...")
    result = modify_gpu_allocations(
        client,
        args.fabric,
        name,
        suid=suid,
        operation="ADD",
        **mode_kwargs(args.mode, args.webhook_url, ["tenant.gpuAllocations"]),
    )
    return settle(client, result, label="gpu-allocations", args=args)


def scenario_delete(client: Client, args: argparse.Namespace) -> bool:
    name = args.tenant_name or default_tenant_name(args.mode)
    print(f"Deleting tenant {name!r} ({args.mode})...")
    result = delete_tenant(
        client,
        args.fabric,
        name,
        **mode_kwargs(args.mode, args.webhook_url, ["tenant.delete"]),
    )
    return settle(client, result, label="delete", args=args)


def scenario_inventory_sync(client: Client, args: argparse.Namespace) -> None:
    print(f"Forcing inventory sync on {args.fabric!r} (UFM fabrics only)...")
    print(f"  -> {inventory_sync(client, args.fabric)}")


def scenario_vpcpeering(client: Client, args: argparse.Namespace) -> bool:
    name = args.tenant_name or default_tenant_name(args.mode)
    vpc = args.vpc_name or f"{name}-{args.fabric}-north-south"
    peer = args.peer_vpc_name or f"{args.fabric}-Storage-VPC"
    peering = args.peering_name or f"{name}-storage-route-leak"
    print(f"Creating VPC peering {peering!r}: {vpc!r} <-> {peer!r} ({args.mode})...")
    result = create_vpc_peering(
        client,
        args.fabric,
        name=peering,
        vpcname=vpc,
        peervpcname=peer,
        **mode_kwargs(args.mode, args.webhook_url, ["vpcpeering.create"]),
    )
    return settle(client, result, label="vpcpeering", args=args)


def scenario_error_handling(client: Client, args: argparse.Namespace) -> None:
    """Show what the SDK raises for the common failure shapes."""
    print("\n--- Scenario: error handling ---")

    print("1. Fetching a tenant that does not exist:")
    try:
        get_tenant(client, args.fabric, "definitely-not-a-real-tenant")
    except Exception as exc:  # noqa: BLE001 - demonstrating the raised type
        print(f"   {type(exc).__name__}: {exc}")

    print("2. Calling an API with a token-less client:")
    try:
        get_all_fabrics(Client(args.url, verify_tls=not args.insecure))
    except NotAuthenticatedError as exc:
        print(f"   NotAuthenticatedError: {exc}")

    print("3. Creating a tenant with an invalid name:")
    try:
        create_tenant(client, args.fabric, tenantName="", maxGpusAllowed=1)
    except Exception as exc:  # noqa: BLE001 - demonstrating the raised type
        print(f"   {type(exc).__name__}: {exc}")


def scenario_lifecycle(client: Client, args: argparse.Namespace) -> None:
    """create -> allocate -> assign ports -> unassign -> deallocate -> delete."""
    name = args.tenant_name or default_tenant_name(args.mode)
    print(f"\n--- Scenario: tenant lifecycle ({args.mode}) on {name!r} ---")

    if not scenario_create(client, args):
        return
    print(f"\nTenant detail: {get_tenant(client, args.fabric, name)}")

    if args.servers:
        if not scenario_servers(client, args, "ADD"):
            return
        scenario_ports(client, args, "ADD")
        scenario_ports(client, args, "DELETE")
        if not scenario_servers(client, args, "DELETE"):
            return

    scenario_delete(client, args)


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------


def default_tenant_name(mode: str) -> str:
    return {"sync": "sdk_sync", "async-poll": "sdk_async"}.get(mode, "sdk_hook")


def csv_list(raw: str | None) -> list[str]:
    return [p.strip() for p in raw.split(",") if p.strip()] if raw else []


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        description="ONES Spectrum-X SDK examples",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    conn = p.add_argument_group("connection")
    conn.add_argument("-u", "--url", required=True, help="ONES base URL")
    conn.add_argument("-U", "--username", required=True, help="login username")
    conn.add_argument("-P", "--password", required=True, help="login password")
    conn.add_argument("-f", "--fabric", required=True, help="fabric name")
    conn.add_argument("--insecure", action="store_true", help="skip TLS verification")

    p.add_argument(
        "--action",
        default="lifecycle",
        choices=[
            "lifecycle", "read-only", "create", "delete", "allocate", "deallocate",
            "assign-ports", "unassign-ports", "gpu-allocations", "inventory-sync",
            "vpcpeering", "error-handling",
        ],
        help="scenario to run",
    )
    p.add_argument("--mode", default="sync", choices=list(MODES), help="execution mode")

    tenant = p.add_argument_group("tenant")
    tenant.add_argument("--tenant-name", help="override the tenant name")
    tenant.add_argument("--max-gpus", type=int, default=8, help="maxGpusAllowed")
    tenant.add_argument("--shared", action="store_true", help="create as shared")
    tenant.add_argument("--servers", type=csv_list, default=[], help="server hostnames")
    tenant.add_argument("--gpu-ids", type=lambda s: [int(x) for x in csv_list(s)],
                        default=[], help="GPU port ids")
    tenant.add_argument("--gpus", type=csv_list, default=[], help="GPU names, G0,G1,...")
    tenant.add_argument("--su-id", type=int, default=0, help="SU index")
    tenant.add_argument("--peering-name", help="peering name for vpcpeering")
    tenant.add_argument("--vpc-name", help="VPC name for vpcpeering")
    tenant.add_argument("--peer-vpc-name", help="peer VPC name for vpcpeering")

    poll = p.add_argument_group("async")
    poll.add_argument("--webhook-url", help="receiver URL (async-webhook)")
    poll.add_argument("--poll-interval", type=int, default=5, help="seconds between polls")
    poll.add_argument("--poll-attempts", type=int, default=60, help="max poll attempts")

    args = p.parse_args()
    if args.mode == "async-webhook" and not args.webhook_url:
        p.error("--webhook-url is required when --mode async-webhook")
    return args


def main() -> int:
    logging.basicConfig(level=logging.INFO)
    args = parse_args()

    client = Client.initialize_with_creds(
        args.url, args.username, args.password, verify_tls=not args.insecure
    )
    try:
        if args.action == "read-only":
            scenario_read_only(client, args)
        elif args.action == "lifecycle":
            scenario_lifecycle(client, args)
        elif args.action == "create":
            scenario_create(client, args)
        elif args.action == "delete":
            scenario_delete(client, args)
        elif args.action == "allocate":
            scenario_servers(client, args, "ADD")
        elif args.action == "deallocate":
            scenario_servers(client, args, "DELETE")
        elif args.action == "assign-ports":
            scenario_ports(client, args, "ADD")
        elif args.action == "unassign-ports":
            scenario_ports(client, args, "DELETE")
        elif args.action == "gpu-allocations":
            scenario_gpu_allocations(client, args)
        elif args.action == "inventory-sync":
            scenario_inventory_sync(client, args)
        elif args.action == "vpcpeering":
            scenario_vpcpeering(client, args)
        elif args.action == "error-handling":
            scenario_error_handling(client, args)
        return 0
    except Exception as exc:  # noqa: BLE001 - CLI: report and exit non-zero
        print(f"error: {exc}", file=sys.stderr)
        return 1
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
