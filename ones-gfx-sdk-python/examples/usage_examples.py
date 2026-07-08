"""
End-to-end usage examples for the ONES Spectrum-X SDK.

This file is intended to be read top-to-bottom by partners as a tour of
the SDK's surface. It covers:

    1. Constructing the client with JWT auth.
    2. Reading fabrics and tenants.
    3. Running a full tenant lifecycle in a single mode (sync/async-poll/async-webhook).
    4. Error handling.

How to run
----------
1. Edit the CONFIG section below with your ONES base URL, JWT tokens,
   and a fabric name that exists in your environment.

2. Choose which mode to run by passing ``--mode`` (sync, async-poll,
    async-webhook). The default is sync.

3. Run the script. Any of these work:

   a) Without installing — directly against the source tree:

          python examples/usage_examples.py       # from repo root
          python usage_examples.py                # from inside examples/
          python -m examples.usage_examples       # module style, from repo root

      The bootstrap block below adds the repo root to sys.path, so the
      ``ones_gfx`` package is importable without ``pip install``.

   b) After installing the SDK (editable or wheel):

          pip install -e .          # from repo root, once
          python examples/usage_examples.py
"""

from __future__ import annotations

# ---------------------------------------------------------------------------
# Bootstrap — make the ``ones_gfx`` package importable when this
# file is run directly from a source checkout, without ``pip install``.
#
# This locates the repo root (the directory two levels above this file:
# usage_example.py -> examples/ -> repo root) and prepends it to
# sys.path, so ``import ones_gfx`` resolves to the local source.
#
# The block is a no-op when the SDK is already installed, because the
# installed package takes precedence over a path-based import only when
# the path entry comes later in sys.path. We use ``insert(0, ...)`` so
# the local source wins during development; remove this block in
# production scripts that should use the installed package.
# ---------------------------------------------------------------------------
import os
import sys

_REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))
if _REPO_ROOT not in sys.path:
    sys.path.insert(0, _REPO_ROOT)
# ---------------------------------------------------------------------------

import argparse
import logging
import time
import warnings

import requests

from urllib3.exceptions import InsecureRequestWarning

from ones_gfx import (
    AuthenticationError,
    ConflictError,
    JWTAuth,
    NotFoundError,
    ONESClient,
    ONESError,
    Operation,
    OperationMode,
    UNLIMITED_GPUS,
)

# ---------------------------------------------------------------------------
# CONFIG — fill these in for your environment.
# ---------------------------------------------------------------------------

BASE_URL = "https://10.4.5.71:8089"
REFRESH_URL = "https://10.4.5.71:8089/refresh"

ACCESS_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InN1cGVyYWRtaW4iLCJyb2xlIjoiU1VQRVJfQURNSU4iLCJwZXJtaXNzaW9ucyI6WyJSRUFEIiwiV1JJVEUiXSwidHlwIjoiYWNjZXNzIiwiaXNzIjoib25lcy1mbSIsImF1ZCI6Im9uZXMtZm0tY2xpZW50IiwianRpIjoiZjFjMGZmN2YtMWY1Yy00NWM4LWJmM2ItNGEwMDFiNzEyY2ExIiwiaWF0IjoxNzgzNTAyNzkwLCJleHAiOjE3ODM1MDk5OTB9.6EgDX-Ez3sCA8zr1SwYdzBuPJSnjBOHjvDLdMKEGNeY"
REFRESH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InN1cGVyYWRtaW4iLCJ0eXAiOiJyZWZyZXNoIiwiaXNzIjoib25lcy1mbSIsImF1ZCI6Im9uZXMtZm0tY2xpZW50IiwianRpIjoiZTk1NjcwYWYtNjFhMi00YzQ2LWI1MzEtMWU2YmQ4YWMzZmNiIiwiaWF0IjoxNzgzNTAyNzkwLCJleHAiOjE3ODM1ODkxOTB9.wQCP_HYSKEdDpeASV8buPDooyEMJy46AF2oMUyo90Qg"

LOGIN_USERNAME = "superadmin"
LOGIN_PASSWORD = "Admin@123456"

FABRIC_NAME = "NMXC"

# Sample server hostnames for UFM fabric.
SAMPLE_SERVERS = ["hgx-su00-h00"]

# Sample server hostnames for NMXC fabric (4 ports per server, IDs 1-4).
NMXC_SAMPLE_SERVERS = ["su00-rack00-node00", "su00-rack01-node00"]

# Webhook receiver URL for the async-webhook example. The SDK does not
# implement the receiver — point this at an HTTP endpoint you control.
WEBHOOK_URL = "http://10.4.5.124:5000/test/webhook-receiver"

# Disable TLS verification only against dev/lab deployments with self-
# signed certs. In production, leave this as True or supply a CA path.
VERIFY_TLS: bool | str = False

# Sync operations can take several minutes (allocate/deallocate up to ~15 min).
DEFAULT_TIMEOUT_S = 1200

# Silence noisy TLS warnings when intentionally disabling verification.
if VERIFY_TLS is False:
    warnings.filterwarnings("ignore", category=InsecureRequestWarning)


# ---------------------------------------------------------------------------
# Auth callback (optional) — persist rotated tokens so they survive a
# process restart. Replace the print() with whatever your secret store
# requires.
# ---------------------------------------------------------------------------

def on_token_refresh(access: str, refresh: str, expires_in: int) -> None:
    print(f"[auth] tokens rotated; new access expires in {expires_in}s")
    # Example:
    #   keyring.set_password("ones", "access_token", access)
    #   keyring.set_password("ones", "refresh_token", refresh)


def build_client() -> ONESClient:
    """Construct an ONESClient with JWT auth."""
    auth = JWTAuth(
        access_token=ACCESS_TOKEN,
        refresh_token=REFRESH_TOKEN,
        refresh_url=REFRESH_URL,
        on_token_refresh=on_token_refresh,
        verify_tls=VERIFY_TLS,
    )
    return ONESClient(
        base_url=BASE_URL,
        auth=auth,
        verify_tls=VERIFY_TLS,
        timeout=DEFAULT_TIMEOUT_S,
    )


# ---------------------------------------------------------------------------
# Scenario 1: list fabrics and tenants.
# ---------------------------------------------------------------------------

def scenario_read_only(client: ONESClient) -> None:
    print("\n--- Scenario: read-only ---")

    fabrics = client.fabrics.list()
    print(f"Found {len(fabrics)} fabric(s):")
    for fab in fabrics:
        print(
            f"  - {fab.fabric_name} (SUs={fab.num_of_sus}/{fab.max_num_of_sus}, "
            f"default storage={fab.default_storage_name})"
        )

    tenants = client.tenants.list(FABRIC_NAME)
    print(f"\nTenants in {FABRIC_NAME}: {len(tenants)}")
    for t in tenants:
        gpu_summary = "unlimited" if t.is_unlimited_gpus else f"{t.gpus_allocated}/{t.max_gpus_allowed}"
        print(f"  - {t.name} | gpus={gpu_summary} | status={t.config_status.name}")

    available = client.tenants.available_servers(FABRIC_NAME)
    print(f"\nAvailable servers in {FABRIC_NAME}: {available}")


def scenario_login(username: str, password: str) -> None:
    print("\n--- Scenario: login ---")
    payload = {"username": username, "password": password}
    url = f"{BASE_URL.rstrip('/')}/login"
    response = requests.post(
        url,
        json=payload,
        timeout=DEFAULT_TIMEOUT_S,
        verify=VERIFY_TLS,
    )
    response.raise_for_status()
    print(response.json())


# ---------------------------------------------------------------------------
# Scenario 2: tenant lifecycle in a single mode.
# ---------------------------------------------------------------------------

def _poll_operation(
    client: ONESClient,
    operation_id: str,
    *,
    label: str,
    max_attempts: int = 360,
    poll_interval_s: int = 5,
) -> bool:
    print(f"  -> polling {label} operation: {operation_id}")
    for attempt in range(1, max_attempts + 1):
        current = client.operations.get(operation_id)
        if current.status in {"RUNNING", "PENDING"}:
            print(".", end="", flush=True)
        else:
            print(f"\n  [poll {attempt}] status={current.status}")
        if current.is_done:
            print("")
            if current.is_success:
                return True
            print(f"  -> FAILED: {current.error_message}")
            return False
        time.sleep(poll_interval_s)
    print("  -> still not done after polling window; giving up")
    return False


def _poll_webhook_status(
    client: ONESClient,
    operation_id: str,
    *,
    label: str,
    max_attempts: int = 360,
    poll_interval_s: int = 5,
) -> bool:
    """Poll the webhook delivery status endpoint until SUCCESS or failure."""
    path = f"operations/{operation_id}/webhook-status"
    print(f"  -> polling {label} webhook delivery: {operation_id}")
    for attempt in range(1, max_attempts + 1):
        result = client._transport.get(path)
        status = result.get("deliveryStatus", "PENDING") if isinstance(result, dict) else "PENDING"
        if status == "SUCCESS":
            print(f"\n  [poll {attempt}] webhook deliveryStatus={status}")
            return True
        if status in {"FAILED", "EXHAUSTED"}:
            print(f"\n  [poll {attempt}] webhook deliveryStatus={status}")
            print(f"  -> FAILED: webhook delivery did not succeed")
            return False
        print(".", end="", flush=True)
        time.sleep(poll_interval_s)
    print("\n  -> still not delivered after polling window; giving up")
    return False


def _tenant_name_for_mode(mode: OperationMode) -> str:
    if mode is OperationMode.SYNCHRONOUS:
        return "sdk_sync"
    if mode is OperationMode.ASYNC_POLL:
        return "sdk_async"
    return "sdk_hook"


def _mode_label(mode: OperationMode) -> str:
    if mode is OperationMode.SYNCHRONOUS:
        return "sync"
    if mode is OperationMode.ASYNC_POLL:
        return "async-poll"
    return "async-webhook"


def _default_vpc_name(tenant_name: str) -> str:
    return f"{tenant_name}-{FABRIC_NAME}-north-south"


def _default_peer_vpc_name() -> str:
    return f"{FABRIC_NAME}-Storage-VPC"


def scenario_tenant_lifecycle(client: ONESClient, mode: OperationMode) -> None:
    label = _mode_label(mode)

    print(f"\n--- Scenario: tenant lifecycle ({label}) ---")

    tenant_name = _tenant_name_for_mode(mode)

    print(f"Creating tenant {tenant_name!r} ({label})...")
    create_result = client.tenants.create(
        fabric_name=FABRIC_NAME,
        name=tenant_name,
        description=f"Created by SDK example ({label} mode)",
        max_gpus_allowed=8,
        mode=mode,
        webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
        webhook_events=["tenant.create"] if mode is OperationMode.ASYNC_WEBHOOK else None,
    )

    if isinstance(create_result, Operation):
        print(f"  -> operation id: {create_result.id}")
        if mode is OperationMode.ASYNC_WEBHOOK:
            print(f"  -> webhook registered: {create_result.webhook_registered}")
        poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
        if not poll_fn(client, create_result.id, label="create"):
            return
    else:
        print(f"  -> id={create_result.id}, vlan={create_result.vlan_id}, vni={create_result.vni_id}")

    if SAMPLE_SERVERS:
        print(f"Allocating GPUs {SAMPLE_SERVERS} to {tenant_name} ({label})...")
        allocate_result = client.tenants.allocate_gpus(
            fabric_name=FABRIC_NAME,
            name=tenant_name,
            servers=SAMPLE_SERVERS,
            timeout=500,
            mode=mode,
            webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
            webhook_events=["tenant.allocate"] if mode is OperationMode.ASYNC_WEBHOOK else None,
        )
        if isinstance(allocate_result, Operation):
            print(f"  -> operation id: {allocate_result.id}")
            poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
            if not poll_fn(client, allocate_result.id, label="allocate"):
                return
        else:
            print("  -> allocate done")

        refreshed = client.tenants.get(FABRIC_NAME, tenant_name)
        print(f"  Tenant now has servers: {refreshed.alloted_servers}")

        # Peering after allocate so tenant VPC / NS path is live with GPUs attached.
        peering_name = f"{tenant_name}-storage-route-leak"
        vpc_name = _default_vpc_name(tenant_name)
        peer_vpc_name = _default_peer_vpc_name()
        print(
            f"Creating VPC peering {peering_name!r} between {vpc_name!r} and {peer_vpc_name!r} (sync-only)..."
        )
        peering_result = client.peering.create(
            fabric_name=FABRIC_NAME,
            name=peering_name,
            vpc_name=vpc_name,
            peer_vpc_name=peer_vpc_name,
        )
        print(f"  -> response: {peering_result}")

        print(f"Deallocating GPUs {SAMPLE_SERVERS} ({label})...")
        deallocate_result = client.tenants.deallocate_gpus(
            fabric_name=FABRIC_NAME,
            name=tenant_name,
            servers=SAMPLE_SERVERS,
            timeout=500,
            mode=mode,
            webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
            webhook_events=["tenant.deallocate"] if mode is OperationMode.ASYNC_WEBHOOK else None,
        )
        if isinstance(deallocate_result, Operation):
            print(f"  -> operation id: {deallocate_result.id}")
            poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
            if not poll_fn(client, deallocate_result.id, label="deallocate"):
                return
        else:
            print("  -> deallocate done")

    print(f"Deleting tenant {tenant_name!r} ({label})...")
    delete_result = client.tenants.delete(
        fabric_name=FABRIC_NAME,
        name=tenant_name,
        mode=mode,
        webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
        webhook_events=["tenant.delete"] if mode is OperationMode.ASYNC_WEBHOOK else None,
    )
    if isinstance(delete_result, Operation):
        print(f"  -> operation id: {delete_result.id}")
        poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
        poll_fn(client, delete_result.id, label="delete")
    else:
        print("  -> delete done")


def _parse_servers(raw: str | None) -> list[str]:
    if not raw:
        return []
    return [item.strip() for item in raw.split(",") if item.strip()]


def _parse_ports(raw: str | None) -> list[int]:
    if not raw:
        return []
    result = []
    for item in raw.split(","):
        item = item.strip()
        if item.isdigit():
            result.append(int(item))
    return result


def _report_sdk_error(err: ONESError) -> None:
    status = f" (status={err.status_code})" if err.status_code is not None else ""
    print(f"[SDK error] {err}{status}")


def scenario_tenant_action(
    client: ONESClient,
    mode: OperationMode,
    action: str,
    tenant_name: str | None,
    servers: list[str] | None,
    port_ids: list[int] | None,
    peering_name: str | None,
    vpc_name: str | None,
    peer_vpc_name: str | None,
) -> None:
    label = _mode_label(mode)
    tenant_name = tenant_name or _tenant_name_for_mode(mode)
    servers = servers or SAMPLE_SERVERS

    if action == "read-only":
        scenario_read_only(client)
        return

    print(f"\n--- Scenario: tenant {action} ({label}) ---")

    try:
        if action == "create":
            print(f"Creating tenant {tenant_name!r} ({label})...")
            create_result = client.tenants.create(
                fabric_name=FABRIC_NAME,
                name=tenant_name,
                description=f"Created by SDK example ({label} mode)",
                max_gpus_allowed=8,
                mode=mode,
                webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
                webhook_events=["tenant.create"] if mode is OperationMode.ASYNC_WEBHOOK else None,
            )

            if isinstance(create_result, Operation):
                print(f"  -> operation id: {create_result.id}")
                if mode is OperationMode.ASYNC_WEBHOOK:
                    print(f"  -> webhook registered: {create_result.webhook_registered}")
                poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
                poll_fn(client, create_result.id, label="create")
            else:
                print(f"  -> id={create_result.id}, vlan={create_result.vlan_id}, vni={create_result.vni_id}")
            return

        if action == "allocate":
            if not servers:
                raise ValueError("No servers provided for allocate. Use --servers or update SAMPLE_SERVERS.")
            print(f"Allocating GPUs {servers} to {tenant_name} ({label})...")
            allocate_result = client.tenants.allocate_gpus(
                fabric_name=FABRIC_NAME,
                name=tenant_name,
                servers=servers,
                timeout=500,
                mode=mode,
                webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
                webhook_events=["tenant.allocate"] if mode is OperationMode.ASYNC_WEBHOOK else None,
            )
            if isinstance(allocate_result, Operation):
                print(f"  -> operation id: {allocate_result.id}")
                poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
                poll_fn(client, allocate_result.id, label="allocate")
            else:
                print("  -> allocate done")
            return

        if action == "deallocate":
            if not servers:
                raise ValueError("No servers provided for deallocate. Use --servers or update SAMPLE_SERVERS.")
            print(f"Deallocating GPUs {servers} ({label})...")
            deallocate_result = client.tenants.deallocate_gpus(
                fabric_name=FABRIC_NAME,
                name=tenant_name,
                servers=servers,
                timeout=500,
                mode=mode,
                webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
                webhook_events=["tenant.deallocate"] if mode is OperationMode.ASYNC_WEBHOOK else None,
            )
            if isinstance(deallocate_result, Operation):
                print(f"  -> operation id: {deallocate_result.id}")
                poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
                poll_fn(client, deallocate_result.id, label="deallocate")
            else:
                print("  -> deallocate done")
            return

        if action == "delete":
            print(f"Deleting tenant {tenant_name!r} ({label})...")
            delete_result = client.tenants.delete(
                fabric_name=FABRIC_NAME,
                name=tenant_name,
                mode=mode,
                webhook_url=WEBHOOK_URL if mode is OperationMode.ASYNC_WEBHOOK else None,
                webhook_events=["tenant.delete"] if mode is OperationMode.ASYNC_WEBHOOK else None,
            )
            if isinstance(delete_result, Operation):
                print(f"  -> operation id: {delete_result.id}")
                poll_fn = _poll_webhook_status if mode is OperationMode.ASYNC_WEBHOOK else _poll_operation
                poll_fn(client, delete_result.id, label="delete")
            else:
                print("  -> delete done")
            return

        if action == "vpcpeering":
            if mode is not OperationMode.SYNCHRONOUS:
                print("VPC peering is supported only in sync mode.")
                return
            vpc_name = vpc_name or _default_vpc_name(tenant_name)
            peer_vpc_name = peer_vpc_name or _default_peer_vpc_name()
            peering_name = peering_name or f"{tenant_name}-storage-route-leak"
            print(
                f"Creating VPC peering {peering_name!r} between {vpc_name!r} and {peer_vpc_name!r}..."
            )
            result = client.peering.create(
                fabric_name=FABRIC_NAME,
                name=peering_name,
                vpc_name=vpc_name,
                peer_vpc_name=peer_vpc_name,
            )
            print(f"  -> response: {result}")
            return

        if action == "assign-ports":
            if not servers:
                raise ValueError("No servers provided. Use --servers.")
            print(f"AssignPorts {servers} ports={port_ids} to {tenant_name!r} on {FABRIC_NAME!r}...")
            client.tenants.assign_ports(
                fabric_name=FABRIC_NAME,
                tenant_name=tenant_name,
                server_names=servers,
                gpu_ids=port_ids,
            )
            print("  -> assign-ports done")
            return

        if action == "unassign-ports":
            if not servers:
                raise ValueError("No servers provided. Use --servers.")
            print(f"UnassignPorts {servers} ports={port_ids} from {tenant_name!r} on {FABRIC_NAME!r}...")
            client.tenants.unassign_ports(
                fabric_name=FABRIC_NAME,
                tenant_name=tenant_name,
                server_names=servers,
                gpu_ids=port_ids,
            )
            print("  -> unassign-ports done")
            return

    except ONESError as e:
        _report_sdk_error(e)
        return


# ---------------------------------------------------------------------------
# Scenario 5: error handling.
# ---------------------------------------------------------------------------

def scenario_error_handling(client: ONESClient) -> None:
    print("\n--- Scenario: error handling ---")

    # Invalid input — caught client-side before any HTTP call.
    try:
        client.tenants.create(
            fabric_name=FABRIC_NAME,
            name="bad_quota_demo",
            description="will not be sent",
            max_gpus_allowed=0,  # disallowed: must be -1 or >= 1
        )
    except ValueError as e:
        print(f"[client-side validation] {e}")

    # Server-side 404 — looking up a tenant that doesn't exist.
    try:
        client.tenants.get(FABRIC_NAME, "this_tenant_does_not_exist_xyz")
    except NotFoundError as e:
        print(f"[NotFoundError] {e} (status={e.status_code})")
    except ONESError as e:
        # Catch-all if the server returns a different shape than expected.
        print(f"[ONESError] {e}")


# ---------------------------------------------------------------------------
# Entry point.
# ---------------------------------------------------------------------------

def main() -> None:
    # Turn on debug logging if you want to see the SDK's internal HTTP
    # activity (token refreshes, retry-on-401, etc.).
    # logging.basicConfig(level=logging.DEBUG)
    logging.basicConfig(level=logging.INFO)

    parser = argparse.ArgumentParser(description="ONES Spectrum-X SDK examples")
    parser.add_argument(
        "--mode",
        default="sync",
        choices=["sync", "async-poll", "async-webhook"],
        help="Tenant lifecycle mode to run (default: sync)",
    )
    parser.add_argument(
        "--action",
        default="lifecycle",
        choices=[
            "lifecycle",
            "read-only",
            "login",
            "create",
            "allocate",
            "deallocate",
            "delete",
            "assign-ports",
            "unassign-ports",
            "vpcpeering",
        ],
        help="Action to run (default: lifecycle)",
    )
    parser.add_argument(
        "--tenant-name",
        default=None,
        help="Override tenant name for create/delete/allocate/deallocate",
    )
    parser.add_argument(
        "--username",
        default=None,
        help="Username for login action (default: LOGIN_USERNAME)",
    )
    parser.add_argument(
        "--password",
        default=None,
        help="Password for login action (default: LOGIN_PASSWORD)",
    )
    parser.add_argument(
        "--servers",
        default=None,
        help="Comma-separated server list for allocate/deallocate/assign-ports/unassign-ports",
    )
    parser.add_argument(
        "--ports",
        default=None,
        help="Comma-separated port IDs for assign-ports/unassign-ports (e.g. 1,2,3,4)",
    )
    parser.add_argument(
        "--peering-name",
        default=None,
        help="Peering name for vpcpeering (default: <tenant>-storage-route-leak)",
    )
    parser.add_argument(
        "--vpc-name",
        default=None,
        help="Tenant VPC name for vpcpeering (default: <tenant>-<fabric>-north-south)",
    )
    parser.add_argument(
        "--peer-vpc-name",
        default=None,
        help="Peer VPC name for vpcpeering (default: <fabric>-Storage-VPC)",
    )
    args = parser.parse_args()

    mode_map = {
        "sync": OperationMode.SYNCHRONOUS,
        "async-poll": OperationMode.ASYNC_POLL,
        "async-webhook": OperationMode.ASYNC_WEBHOOK,
    }

    # The ``with`` block ensures the underlying HTTP session is closed
    # cleanly even if a scenario raises.
    try:
        if args.action == "login":
            username = args.username or LOGIN_USERNAME
            password = args.password or LOGIN_PASSWORD
            scenario_login(username, password)
            return

        with build_client() as client:
            if args.action == "lifecycle":
                scenario_read_only(client)
                scenario_tenant_lifecycle(client, mode_map[args.mode])
                scenario_error_handling(client)
            else:
                scenario_tenant_action(
                    client,
                    mode_map[args.mode],
                    args.action,
                    args.tenant_name,
                    _parse_servers(args.servers),
                    _parse_ports(args.ports),
                    args.peering_name,
                    args.vpc_name,
                    args.peer_vpc_name,
                )
    except AuthenticationError as e:
        print(f"\nAuth failed: {e}")
        print("Update ACCESS_TOKEN and REFRESH_TOKEN at the top of this file.")
    except ONESError as e:
        _report_sdk_error(e)


if __name__ == "__main__":
    main()
