"""Quick manual test of the Topology Audit SDK surface (GFX-6219/GFX-6225):
trigger an on-demand audit, then poll until it finishes.

Run from the SDK root:

    python3 examples/topology_audit_example.py --fabric Fab101
"""
from __future__ import annotations

import argparse
import time

from ones_gfx import Client, OperationFailedError, OperationNotDoneError
from ones_gfx.apis import get_topology_audit_result, get_topology_audit_status, trigger_topology_audit

# --- configuration -----------------------------------------------------------

BASE_URL = "https://localhost:3002"
USERNAME = "superadmin"
PASSWORD = "Admin@123456"

# -----------------------------------------------------------------------------


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--fabric", default="Fab101", help="fabric to audit")
    args = parser.parse_args()

    print(f"Logging in to {BASE_URL} as {USERNAME!r} ...")
    client = Client.initialize_with_creds(BASE_URL, USERNAME, PASSWORD)

    print(f"\nTriggering topology audit for fabric {args.fabric!r} ...")
    accepted = trigger_topology_audit(client, args.fabric)
    operation_id = accepted["operationId"]
    print(f"Accepted: operationId={operation_id} status={accepted.get('status')}")

    print("\nPolling for the result ...")
    while True:
        try:
            report = get_topology_audit_result(client, operation_id)
        except OperationNotDoneError:
            status = get_topology_audit_status(client, operation_id)
            print(f"  still running (status={status.get('status')}) ...")
            time.sleep(5)
            continue
        except OperationFailedError as e:
            print("audit failed:", e.error_message)
            return

        print(f"\nOverall status: {report['overallStatus']} (intent version {report['intentVersion']})")
        for link in report["links"]:
            print(
                f"  {link['status']:<15} {link['device']}:{link['localInterface']} -> "
                f"{link.get('remoteDevice')}:{link.get('remoteInterface')} "
                f"(observed {link.get('observedRemoteDevice')}:{link.get('observedRemoteInterface')})  "
                f"{link['message']}"
            )
        return


if __name__ == "__main__":
    main()
