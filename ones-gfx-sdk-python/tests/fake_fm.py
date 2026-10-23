"""In-process stand-in for the SES_FM Fabric Manager REST surface used by the
IPAM and rail-CIDRPool endpoints.

Behaviour is ported from SES_FM ``FabricManager.java`` (controller methods
``getRailCidrPools``, ``getRailCidrPool``, ``updateIpamConfig``,
``getIpamConfig``, ``updateTenant``, ``getTenantComputeIp``, ``getTenantIps``,
``getServerIps``) and ``IpamService.updateConfig`` — status codes, error codes,
Accept negotiation and response shapes — and the happy-path payloads are the
ones captured from a live instance. It serves real HTTP so the SDK's actual
``requests`` stack is exercised end to end.
"""
from __future__ import annotations

import json
import re
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import parse_qs, unquote, urlparse

TOKEN = "test-token"

RAIL_YAML = {
    "r0": """apiVersion: nv-ipam.nvidia.com/v1alpha1
kind: CIDRPool
metadata:
  name: r0
  namespace: aviz-network-operator
spec:
  cidr: 172.16.0.0/15
  gatewayIndex: 0
  perNodeNetworkPrefix: 31
  routes:
  - dst: 172.16.0.0/15
  - dst: 172.16.0.0/12
  staticAllocations:
  - gateway: 172.16.0.1
    nodeName: hgx-su00-h00
    prefix: 172.16.0.0/31
  - gateway: 172.16.0.3
    nodeName: hgx-su00-h01
    prefix: 172.16.0.2/31
""",
    "r2": """apiVersion: nv-ipam.nvidia.com/v1alpha1
kind: CIDRPool
metadata:
  name: r2
  namespace: aviz-network-operator
spec:
  cidr: 172.20.0.0/15
  gatewayIndex: 0
  perNodeNetworkPrefix: 31
  routes:
  - dst: 172.20.0.0/15
  - dst: 172.16.0.0/12
  staticAllocations:
  - gateway: 172.20.0.1
    nodeName: hgx-su00-h00
    prefix: 172.20.0.0/31
  - gateway: 172.20.0.3
    nodeName: hgx-su00-h01
    prefix: 172.20.0.2/31
""",
}
ALL_RAILS_YAML = "---\n".join(RAIL_YAML.values())

POOL_CIDR = "192.3.0.0/17"
POOL_CAPACITY = 32765


def _iso(n: int = 0) -> str:
    return f"2026-09-24T09:56:{21 + n:02d}.645313Z"


class FakeFabricManager:
    """Mutable server state plus the HTTP server that fronts it."""

    def __init__(self) -> None:
        # fabric -> isOnesControlled
        self.fabrics = {"yml": False, "ns": False, "ones-ctl": True, "norails": False}
        self.servers = {"ns": {"su00-node00", "su00-node01", "su00-node02", "full-node"}}
        self.tenants = {"ns": {"gorax"}}
        self.ipam_config: dict[str, dict] = {}
        self.pools: set[tuple[str, str]] = {("ns", "gorax")}
        self.allocs: list[dict] = []
        self.next_id = 2
        self.next_host = 2
        self.requests: list[dict] = []  # every request received, for assertions
        self._server: ThreadingHTTPServer | None = None

    # -- lifecycle ----------------------------------------------------------
    def start(self) -> str:
        fm = self

        class Handler(BaseHTTPRequestHandler):
            def log_message(self, *a):  # keep test output quiet
                pass

            def do_GET(self):
                fm._handle(self, "GET")

            def do_PATCH(self):
                fm._handle(self, "PATCH")

            def do_POST(self):  # no POST route is modelled: records the request, answers 404
                fm._handle(self, "POST")

        self._server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        threading.Thread(target=self._server.serve_forever, kwargs={"poll_interval": 0.01}, daemon=True).start()
        return f"http://127.0.0.1:{self._server.server_address[1]}"

    def stop(self) -> None:
        if self._server:
            self._server.shutdown()
            self._server.server_close()

    # -- plumbing -----------------------------------------------------------
    @staticmethod
    def _send(h, status: int, body, content_type="application/json"):
        raw = body if isinstance(body, str) else json.dumps(body)
        data = raw.encode()
        h.send_response(status)
        h.send_header("Content-Type", content_type + "; charset=UTF-8")
        h.send_header("Content-Length", str(len(data)))
        h.end_headers()
        h.wfile.write(data)

    @staticmethod
    def _ok(message, data=None):
        out = {"status": "success", "message": message}
        if data is not None:
            out["data"] = data
        return out

    @staticmethod
    def _err(code, detail):
        return {"status": "error", "message": detail, "error": {"code": code, "detail": detail}}

    def _handle(self, h, method: str) -> None:
        url = urlparse(h.path)
        length = int(h.headers.get("Content-Length") or 0)
        raw_body = h.rfile.read(length) if length else b""
        self.requests.append(
            {
                "method": method,
                "path": url.path,
                "query": parse_qs(url.query),
                "headers": {k.lower(): v for k, v in h.headers.items()},
                "body": raw_body.decode() if raw_body else None,
            }
        )
        if h.headers.get("authorization") != TOKEN:
            return self._send(h, 401, self._err("UNAUTHORIZED", "Invalid token."))
        m = re.fullmatch(r"/api/fm/(.*)", url.path)
        if not m:
            return self._send(h, 404, self._err("NOT_FOUND", "No such route."))
        parts = [unquote(p) for p in m.group(1).split("/")]
        query = parse_qs(url.query)
        try:
            status, body, ctype = self._route(method, parts, query, raw_body, h.headers.get("Accept"))
        except Exception as exc:  # surface stub bugs as 500s
            status, body, ctype = 500, self._err("INTERNAL_ERROR", str(exc)), "application/json"
        self._send(h, status, body, ctype)

    # -- routing ------------------------------------------------------------
    def _route(self, method, p, q, raw, accept):
        J = "application/json"
        if p[0] != "fabrics" or len(p) < 3:
            return 404, self._err("NOT_FOUND", "No such route."), J
        fabric, rest = p[1], p[2:]
        if method == "GET" and rest == ["rails", "cidrpools"]:
            return self._rails(fabric, None, accept)
        if method == "GET" and len(rest) == 3 and rest[0] == "rails" and rest[2] == "cidrpool":
            return self._rails(fabric, rest[1], accept)
        if rest == ["ipam", "config"]:
            if method == "GET":
                return self._get_ipam_config(fabric)
            return self._patch_ipam_config(fabric, q, raw)
        if method == "PATCH" and len(rest) == 2 and rest[0] == "tenants":
            return self._patch_tenant(fabric, rest[1], raw)
        if method == "GET" and len(rest) == 5 and rest[0] == "tenants" and rest[2] == "computes" and rest[4] == "ip-address":
            return self._tenant_compute_ip(fabric, rest[1], rest[3])
        if method == "GET" and len(rest) == 3 and rest[0] == "tenants" and rest[2] == "ip-addresses":
            return self._tenant_ips(fabric, rest[1], q)
        if method == "GET" and len(rest) == 3 and rest[0] == "computes" and rest[2] == "ip-addresses":
            return self._server_ips(fabric, rest[1])
        return 404, self._err("NOT_FOUND", "No such route."), J

    # -- rails (FabricManager.getRailCidrPools / getRailCidrPool) ------------
    @staticmethod
    def _resolve_raw_media_type(accept):
        """Port of resolveRailYamlMediaType: None => JSON envelope."""
        if not accept or not accept.strip():
            return None
        for item in accept.split(","):
            mt = item.split(";")[0].strip().lower()
            if mt in ("application/json", "application/*", "*/*"):
                return None
            if mt in ("text/plain", "text/*", "application/yaml", "application/x-yaml"):
                return "text/plain" if mt == "text/*" else mt
        return None

    def _rail_err(self, raw_mt, status, code, detail):
        if raw_mt is not None:
            return status, f"{code}: {detail}", "text/plain"
        return status, self._err(code, detail), "application/json"

    def _rails(self, fabric, rail, accept):
        raw_mt = self._resolve_raw_media_type(accept)
        if fabric not in self.fabrics:
            return self._rail_err(raw_mt, 404, "INVALID_FABRIC", "Specified fabric does not exist.")
        if fabric == "norails":
            return self._rail_err(raw_mt, 404, "RAIL_CIDRPOOL_NOT_FOUND", f"No rail CIDRPool data is available for fabric '{fabric}'.")
        if rail is None:
            text, extra = ALL_RAILS_YAML, {"rails": list(RAIL_YAML)}
        else:
            if rail not in RAIL_YAML:
                return self._rail_err(raw_mt, 404, "INVALID_RAIL", f"Rail '{rail}' does not exist in fabric '{fabric}'.")
            text, extra = RAIL_YAML[rail], {"railName": rail}
        if raw_mt is not None:
            return 200, text, raw_mt
        data = {"fabricName": fabric, **extra, "yaml": text}
        msg = "Rail CIDRPools retrieved successfully." if rail is None else "Rail CIDRPool retrieved successfully."
        return 200, self._ok(msg, data), "application/json"

    # -- IPAM config ----------------------------------------------------------
    def _cfg(self, fabric):
        return self.ipam_config.get(
            fabric,
            {"fabricName": fabric, "enableIpam": False, "enableIpAssign": False, "updatedAt": _iso(), "updatedBy": None},
        )

    def _get_ipam_config(self, fabric):
        if fabric not in self.fabrics:
            return 404, self._err("FABRIC_NOT_FOUND", "Specified fabric does not exist."), "application/json"
        return 200, self._ok("IPAM configuration retrieved successfully.", self._cfg(fabric)), "application/json"

    def _patch_ipam_config(self, fabric, q, raw):
        J = "application/json"
        if fabric not in self.fabrics:
            return 404, self._err("FABRIC_NOT_FOUND", "Specified fabric does not exist."), J
        if self.fabrics[fabric]:
            return 409, self._err("FABRIC_NOT_EXTERNALLY_MANAGED", "IPAM configuration can only be modified for externally managed (non-ONES-controlled) fabrics."), J
        body = json.loads(raw or "null")
        if not isinstance(body, dict) or "enableIpam" not in body or "enableIpAssign" not in body:
            return 400, self._err("INVALID_REQUEST_BODY", "enableIpam and enableIpAssign are required."), J
        if not isinstance(body["enableIpam"], bool) or not isinstance(body["enableIpAssign"], bool):
            return 400, self._err("INVALID_REQUEST_BODY", "enableIpam and enableIpAssign must be booleans."), J
        enable, assign = body["enableIpam"], body["enableIpAssign"]
        force = (q.get("force", ["false"])[0]).lower() == "true"
        if assign and not enable:
            return 400, self._err("INVALID_IPAM_CONFIG", "enableIpAssign requires enableIpam to be true."), J
        active = [pool for pool in self.pools if pool[0] == fabric]
        if not enable and not force and active:
            return 409, self._err("IPAM_HAS_ACTIVE_POOLS", "Active IPAM pools exist for this fabric."), J
        if not enable and force:
            self.pools = {pool for pool in self.pools if pool[0] != fabric}
            self.allocs = [a for a in self.allocs if a["fabricName"] != fabric]
        cfg = {"fabricName": fabric, "enableIpam": enable, "enableIpAssign": assign, "updatedAt": "2026-09-24T08:57:28.254792109Z", "updatedBy": "Aviz Admin"}
        self.ipam_config[fabric] = cfg
        return 200, self._ok("IPAM configuration updated successfully.", cfg), J

    # -- tenant update (IPAM branch of updateTenant) --------------------------
    def _alloc_for(self, fabric, tenant, server):
        return next((a for a in self.allocs if (a["fabricName"], a["tenantName"], a["serverName"]) == (fabric, tenant, server)), None)

    def _patch_tenant(self, fabric, tenant, raw):
        J = "application/json"
        body = json.loads(raw or "null") or {}
        servers = body.get("servers")
        if not servers:
            return 400, self._err("INVALID_SERVER_LIST", "Server list is null or empty."), J
        if body.get("operation") not in ("ADD", "DELETE"):
            return 400, self._err("INVALID_OPERATION", "Operation field is missing or invalid."), J
        if tenant not in self.tenants.get(fabric, set()):
            return 404, self._err("TENANT_NOT_FOUND", "Tenant does not exist."), J
        names = [s["serverName"] for s in servers]
        if any(n not in self.servers.get(fabric, set()) for n in names):
            return 404, self._err("SERVER_NOT_FOUND", "One or more servers not found."), J
        cfg = self._cfg(fabric)
        if not cfg["enableIpam"]:
            return 200, self._ok("GPU update processed successfully."), J
        partial, out = False, []
        for n in names:
            if body["operation"] == "ADD":
                if n == "full-node":  # pool exhausted
                    out.append({"serverName": n, "status": "IP_POOL_EXHAUSTED", "ipAddress": None, "ipAssignStatus": "IP_POOL_EXHAUSTED", "message": "No free address in the tenant pool."})
                    partial = True
                    continue
                a = self._alloc_for(fabric, tenant, n)
                if a is None:
                    a = {
                        "id": self.next_id, "poolId": 1, "fabricName": fabric, "tenantName": tenant, "tenantId": 3,
                        "serverName": n, "ipAddress": f"192.3.0.{self.next_host}", "prefixLength": 32,
                        "assignStatus": "ASSIGNED" if cfg["enableIpAssign"] else "ALLOCATED_ONLY",
                        "cleanupStatus": "NONE", "allocatedAt": _iso(self.next_id),
                    }
                    self.next_id += 1
                    self.next_host += 1
                    self.allocs.append(a)
                out.append({"serverName": n, "status": "SUCCESS", "ipAddress": f"{a['ipAddress']}/{a['prefixLength']}", "ipAssignStatus": a["assignStatus"]})
            else:
                a = self._alloc_for(fabric, tenant, n)
                if a is None:
                    out.append({"serverName": n, "status": "FAILED", "ipAddress": None, "ipAssignStatus": "NOT_ALLOCATED", "message": "Server has no active IP allocation."})
                    partial = True
                    continue
                self.allocs.remove(a)
                out.append({"serverName": n, "status": "SUCCESS", "ipAddress": f"{a['ipAddress']}/{a['prefixLength']}", "ipAssignStatus": a["assignStatus"], "released": True})
        data = {"fabricName": fabric, "tenantName": tenant, "operation": body["operation"], "servers": out, "partial": partial}
        return (207 if partial else 200), self._ok("GPU update processed successfully.", data), J

    # -- reads ------------------------------------------------------------------
    def _tenant_compute_ip(self, fabric, tenant, server):
        J = "application/json"
        a = self._alloc_for(fabric, tenant, server)
        if a is None:
            if (fabric, tenant) not in self.pools:
                return 404, self._err("IPAM_NOT_CONFIGURED", "No IPAM pool exists for the tenant."), J
            return 404, self._err("SERVER_NOT_ALLOCATED_TO_TENANT", "Server has no IP allocation."), J
        return 200, self._ok("IP allocation retrieved successfully.", a), J

    def _tenant_ips(self, fabric, tenant, q):
        J = "application/json"
        if (fabric, tenant) not in self.pools:
            return 404, self._err("IPAM_NOT_CONFIGURED", "No IPAM pool exists for the tenant."), J
        status = q.get("status", [None])[0]
        if status is not None and status not in ("ALLOCATED_ONLY", "ASSIGNED", "IP_ASSIGN_FAILED"):
            return 400, self._err("BAD_REQUEST", f"Invalid status '{status}'."), J
        page = max(0, int(q.get("page", ["0"])[0]))
        size = max(1, int(q.get("size", ["50"])[0]))
        rows = [a for a in self.allocs if a["fabricName"] == fabric and a["tenantName"] == tenant and (status is None or a["assignStatus"] == status)]
        rows.sort(key=lambda a: tuple(int(x) for x in a["ipAddress"].split(".")))
        total = len(rows)
        pages = -(-total // size)
        data = {
            "fabricName": fabric, "tenantName": tenant,
            "pool": {"cidr": POOL_CIDR, "allocatableAddressCount": POOL_CAPACITY, "allocatedCount": len([a for a in self.allocs if a["tenantName"] == tenant])},
            "allocations": rows[page * size : (page + 1) * size],
            "pagination": {"page": page, "size": size, "totalElements": total, "totalPages": pages},
        }
        return 200, self._ok("IP allocations retrieved successfully.", data), J

    def _server_ips(self, fabric, server):
        J = "application/json"
        if server not in self.servers.get(fabric, set()):
            return 404, self._err("SERVER_NOT_FOUND", "Server was not found in the fabric inventory."), J
        data = {"fabricName": fabric, "serverName": server, "allocations": [a for a in self.allocs if a["fabricName"] == fabric and a["serverName"] == server]}
        return 200, self._ok("IP allocations retrieved successfully.", data), J
