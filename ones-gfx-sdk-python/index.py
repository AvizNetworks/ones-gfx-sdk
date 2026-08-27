"""
Quick manual test for the ONES SDK auth flow.

Run from the SDK root:

    python3 index.py

Edit BASE_URL / USERNAME / PASSWORD below (or set them as env vars) to point
at your ONES instance.
"""

import os
import sys

# Make sure we import the local ones_gfx package (no install required).
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from ones_gfx.client import ONESClient
from ones_gfx.apis import get_all_fabrics
import ones_gfx.client as client_module

# # --- configuration -----------------------------------------------------------
# BASE_URL = os.environ.get("ONES_BASE_URL", "https://10.4.5.76:8089")
# USERNAME = os.environ.get("ONES_USERNAME", "admin")
# PASSWORD = os.environ.get("ONES_PASSWORD", "changeme")
# # -----------------------------------------------------------------------------

BASE_URL = 'https://localhost:3002'
USERNAME = 'superadmin'
PASSWORD = 'Admin@123456'

def main() -> None:
    # print(f"Token on startup (from secrets.json): {client_module.auth_token!r}")

    client = ONESClient(base_url=BASE_URL, verify_tls=False)

    print(f"\nLogging in to {BASE_URL} as {USERNAME!r} ...")
    # result = client.login(USERNAME, PASSWORD)
    # print("Login response:", result)
    # print("Global auth_token now:", client_module.auth_token)

    print("\nFetching inventory devices ...")
    # try:
    #     devices = client.getTeleDevices()
    #     print("Devices:", devices)
    # except Exception as exc:  # noqa: BLE001 - manual test script
    #     print("getTeleDevices failed:", exc)

    # print("\nFetching all fabrics ...")
    try:
        fabrics = get_all_fabrics(client)
        print("Fabrics:", fabrics)
    except Exception as exc:  # noqa: BLE001 - manual test script
        print("get_all_fabrics failed:", exc)

    # print("\nLogging out ...")
    # try:
    #     logout_result = client.logout()
    #     print("Logout response:", logout_result)
    # except Exception as exc:  # noqa: BLE001 - manual test script
    #     print("Logout failed:", exc)
    # print("Global auth_token now:", client_module.auth_token)


if __name__ == "__main__":
    main()
