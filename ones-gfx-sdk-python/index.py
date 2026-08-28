"""
Quick manual test for the ONES SDK auth flow.

Run from the SDK root:

    python3 index.py

Edit BASE_URL / USERNAME / PASSWORD below to point at your ONES instance.
"""

import os
import sys

# Make sure we import the local ones_gfx package (no install required).
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from ones_gfx import Client
from ones_gfx.apis import add_fabric_data, get_all_fabrics

BASE_URL = 'https://localhost:3002'
USERNAME = 'superadmin'
PASSWORD = 'Admin@123456'


def main() -> None:
    print(f"Logging in to {BASE_URL} as {USERNAME!r} ...")
    client = Client.initialize_with_creds(
        BASE_URL, USERNAME, PASSWORD, verify_tls=False
    )
    print(f"Authenticated: {client.is_authenticated()} "
          f"(token {len(client.get_auth_token() or '')} chars)")

    print("\nFetching all fabrics ...")
    try:
        fabrics = get_all_fabrics(client)
        print(f"Fabrics ({len(fabrics)}):")
        for f in fabrics:
            print(f"  - id={f.get('id')} name={f.get('name')!r} "
                  f"type={f.get('type')!r} status={f.get('status')!r}")
    except Exception as exc:  # noqa: BLE001 - manual test script
        print("get_all_fabrics failed:", exc)

    print("\nAdding a fabric ...")
    try:
        result = add_fabric_data(
            client,
            name="CLI ASN Fabric",
            type="DNO ASN",
            description="fabric created fromcli",
            status="draft",
        )
        print("Add fabric response:", result)
    except Exception as exc:  # noqa: BLE001 - manual test script
        print("add_fabric_data failed:", exc)

    # print("\nRefreshing auth ...")
    # try:
    #     refreshed = client.refresh_auth()
    #     print("Refresh response:", refreshed)
    # except Exception as exc:  # noqa: BLE001 - manual test script
    #     print("refresh_auth failed:", exc)

    # print("\nLogging out ...")
    # try:
    #     print("Logout response:", client.logout())
    # except Exception as exc:  # noqa: BLE001 - manual test script
    #     print("logout failed:", exc)


if __name__ == "__main__":
    main()
