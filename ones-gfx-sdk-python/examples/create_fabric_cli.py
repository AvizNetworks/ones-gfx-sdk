"""
Create a fabric from the command line.

Usage
-----
    python examples/create_fabric_cli.py \
        --url https://localhost:3002 \
        --username superadmin \
        --password 'Admin@123456' \
        --name "gpu-fabric-1" \
        --type "DNO ASN" \
        --description "Primary GPU fabric" \
        --status draft \
        --verify-tls

Minimal (only the connection details and a fabric name are required):

    python examples/create_fabric_cli.py \
        -u https://localhost:3002 -U superadmin -P 'Admin@123456' \
        -n "gpu-fabric-1" --verify-tls

List the fabrics afterwards instead of creating one:

    python examples/create_fabric_cli.py \
        -u https://localhost:3002 -U superadmin -P 'Admin@123456' \
        -n unused --list --verify-tls

All flags:

    --url/-u          ONES base URL, e.g. https://host:3002     (required)
    --username/-U     login username                            (required)
    --password/-P     login password                            (required)
    --name/-n         fabric name                               (required)
    --type            fabric type, e.g. "DNO ASN"
    --description     free-text description
    --status          fabric status, e.g. draft
    --num-sus         number of SUs                             (int)
    --max-sus         maximum number of SUs                     (int)
    --dedicated       mark the fabric dedicated                 (flag)
    --list            list fabrics instead of creating one
           skip TLS certificate verification
"""

import argparse
import os
import sys

# Run straight from a checkout without installing.
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from ones_gfx import Client
from ones_gfx.apis import add_fabric_data, get_all_fabrics


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        description="Create a fabric on a ONES instance.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    conn = p.add_argument_group("connection")
    conn.add_argument("-u", "--url", required=True, help="ONES base URL")
    conn.add_argument("-U", "--username", required=True, help="login username")
    conn.add_argument("-P", "--password", required=True, help="login password")
    conn.add_argument(
        "--verify-tls", action="store_true",
        help="enforce strict TLS certificate checking (default: off, self-signed OK)"
    )

    fab = p.add_argument_group("fabric")
    fab.add_argument("-n", "--name", required=True, help="fabric name")
    fab.add_argument("--type", dest="fabric_type", help='fabric type, e.g. "DNO ASN"')
    fab.add_argument("--description", help="free-text description")
    fab.add_argument("--status", help="fabric status, e.g. draft")
    fab.add_argument("--num-sus", type=int, help="number of SUs")
    fab.add_argument("--max-sus", type=int, help="maximum number of SUs")
    fab.add_argument(
        "--dedicated", action="store_true", default=None, help="mark fabric dedicated"
    )

    p.add_argument(
        "--list", action="store_true", help="list fabrics instead of creating one"
    )
    return p.parse_args()


def create_fabric(client: Client, args: argparse.Namespace) -> str:
    """Create one fabric from the parsed CLI arguments.

    Only the flags actually supplied are sent — everything left unset stays
    None and is omitted from the request body.
    """
    return add_fabric_data(
        client,
        name=args.name,
        type=args.fabric_type,
        description=args.description,
        status=args.status,
        numOfSus=args.num_sus,
        maxNumOfSus=args.max_sus,
        dedicated=args.dedicated,
    )


def main() -> int:
    args = parse_args()

    client = Client.initialize_with_creds(
        args.url, args.username, args.password, verify_tls=args.verify_tls
    )
    try:
        if args.list:
            fabrics = get_all_fabrics(client)
            print(f"Fabrics ({len(fabrics)}):")
            for f in fabrics:
                print(
                    f"  - id={f.get('id')} name={f.get('name')!r} "
                    f"type={f.get('type')!r} status={f.get('status')!r}"
                )
            return 0

        print(f"Creating fabric {args.name!r} on {args.url} ...")
        print("Response:", create_fabric(client, args))
        return 0
    except Exception as exc:  # noqa: BLE001 - CLI: report and exit non-zero
        print(f"error: {exc}", file=sys.stderr)
        return 1
    finally:
        client.close()


if __name__ == "__main__":
    raise SystemExit(main())
