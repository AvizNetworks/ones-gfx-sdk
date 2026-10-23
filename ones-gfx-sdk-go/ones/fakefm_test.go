package ones_test

// fakefm_test.go is an in-process stand-in for the SES_FM Fabric Manager REST
// surface used by the IPAM and rail-CIDRPool endpoints. Status codes, error
// codes, Accept negotiation and response shapes are ported from SES_FM
// FabricManager.java (getRailCidrPools, getRailCidrPool, updateIpamConfig,
// getIpamConfig, updateTenant, getTenantComputeIp, getTenantIps, getServerIps)
// and IpamService.updateConfig; happy-path payloads come from a live instance.
// It mirrors ones-gfx-sdk-python/tests/fake_fm.py.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const testToken = "test-token"

var railYAML = map[string]string{
	"r0": `apiVersion: nv-ipam.nvidia.com/v1alpha1
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
`,
	"r2": `apiVersion: nv-ipam.nvidia.com/v1alpha1
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
`,
}

var railOrder = []string{"r0", "r2"}

func allRailsYAML() string { return railYAML["r0"] + "---\n" + railYAML["r2"] }

type recordedReq struct {
	Method string
	Path   string // escaped path as received
	Query  url.Values
	Header http.Header
	Body   string
}

type fakeFM struct {
	t   *testing.T
	srv *httptest.Server
	mu  sync.Mutex

	fabrics  map[string]bool // fabric -> isOnesControlled
	servers  map[string]map[string]bool
	tenants  map[string]map[string]bool
	cfg      map[string]map[string]any
	pools    map[[2]string]bool
	allocs   []map[string]any
	nextID   int
	nextHost int
	Requests []recordedReq
}

func newFakeFM(t *testing.T) *fakeFM {
	f := &fakeFM{
		t:        t,
		fabrics:  map[string]bool{"yml": false, "ns": false, "ones-ctl": true, "norails": false},
		servers:  map[string]map[string]bool{"ns": {"su00-node00": true, "su00-node01": true, "su00-node02": true, "full-node": true}},
		tenants:  map[string]map[string]bool{"ns": {"gorax": true}},
		cfg:      map[string]map[string]any{},
		pools:    map[[2]string]bool{{"ns", "gorax"}: true},
		nextID:   2,
		nextHost: 2,
	}
	f.srv = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeFM) last() recordedReq {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Requests[len(f.Requests)-1]
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	b, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	w.Write(b)
}

func writeText(w http.ResponseWriter, status int, contentType, body string) {
	w.Header().Set("Content-Type", contentType+";charset=UTF-8")
	w.WriteHeader(status)
	io.WriteString(w, body)
}

func okBody(msg string, data any) map[string]any {
	m := map[string]any{"status": "success", "message": msg}
	if data != nil {
		m["data"] = data
	}
	return m
}

func errBody(code, detail string) map[string]any {
	return map[string]any{"status": "error", "message": detail, "error": map[string]any{"code": code, "detail": detail}}
}

func (f *fakeFM) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	body, _ := io.ReadAll(r.Body)
	f.Requests = append(f.Requests, recordedReq{r.Method, r.URL.EscapedPath(), r.URL.Query(), r.Header.Clone(), string(body)})
	if r.Header.Get("authorization") != testToken {
		writeJSON(w, 401, errBody("UNAUTHORIZED", "Invalid token."))
		return
	}
	const prefix = "/api/fm/"
	if !strings.HasPrefix(r.URL.EscapedPath(), prefix) {
		writeJSON(w, 404, errBody("NOT_FOUND", "No such route."))
		return
	}
	p := strings.Split(strings.TrimPrefix(r.URL.EscapedPath(), prefix), "/")
	for i := range p {
		// the SDK escapes each segment; Go's server gives us the escaped form
		if u, err := urlUnescape(p[i]); err == nil {
			p[i] = u
		}
	}
	q := r.URL.Query()
	if len(p) < 3 || p[0] != "fabrics" {
		writeJSON(w, 404, errBody("NOT_FOUND", "No such route."))
		return
	}
	fabric, rest := p[1], p[2:]
	switch {
	case r.Method == "GET" && len(rest) == 2 && rest[0] == "rails" && rest[1] == "cidrpools":
		f.rails(w, fabric, "", r.Header.Get("Accept"))
	case r.Method == "GET" && len(rest) == 3 && rest[0] == "rails" && rest[2] == "cidrpool":
		f.rails(w, fabric, rest[1], r.Header.Get("Accept"))
	case len(rest) == 2 && rest[0] == "ipam" && rest[1] == "config" && r.Method == "GET":
		f.getConfig(w, fabric)
	case len(rest) == 2 && rest[0] == "ipam" && rest[1] == "config" && r.Method == "PATCH":
		f.patchConfig(w, fabric, q.Get("force"), body)
	case r.Method == "PATCH" && len(rest) == 2 && rest[0] == "tenants":
		f.patchTenant(w, fabric, rest[1], body)
	case r.Method == "GET" && len(rest) == 5 && rest[0] == "tenants" && rest[2] == "computes" && rest[4] == "ip-address":
		f.computeIP(w, fabric, rest[1], rest[3])
	case r.Method == "GET" && len(rest) == 3 && rest[0] == "tenants" && rest[2] == "ip-addresses":
		f.tenantIPs(w, fabric, rest[1], q.Get("page"), q.Get("size"), q.Get("status"))
	case r.Method == "GET" && len(rest) == 3 && rest[0] == "computes" && rest[2] == "ip-addresses":
		f.serverIPs(w, fabric, rest[1])
	default:
		writeJSON(w, 404, errBody("NOT_FOUND", "No such route."))
	}
}

// resolveRawMediaType ports FabricManager.resolveRailYamlMediaType: "" => JSON envelope.
func resolveRawMediaType(accept string) string {
	if strings.TrimSpace(accept) == "" {
		return ""
	}
	for _, item := range strings.Split(accept, ",") {
		mt := strings.ToLower(strings.TrimSpace(strings.Split(item, ";")[0]))
		switch mt {
		case "application/json", "application/*", "*/*":
			return ""
		case "text/*":
			return "text/plain"
		case "text/plain", "application/yaml", "application/x-yaml":
			return mt
		}
	}
	return ""
}

func (f *fakeFM) railErr(w http.ResponseWriter, raw string, status int, code, detail string) {
	if raw != "" {
		writeText(w, status, "text/plain", code+": "+detail)
		return
	}
	writeJSON(w, status, errBody(code, detail))
}

func (f *fakeFM) rails(w http.ResponseWriter, fabric, rail, accept string) {
	raw := resolveRawMediaType(accept)
	if _, ok := f.fabrics[fabric]; !ok {
		f.railErr(w, raw, 404, "INVALID_FABRIC", "Specified fabric does not exist.")
		return
	}
	if fabric == "norails" {
		f.railErr(w, raw, 404, "RAIL_CIDRPOOL_NOT_FOUND", "No rail CIDRPool data is available for fabric '"+fabric+"'.")
		return
	}
	var text string
	data := map[string]any{"fabricName": fabric}
	if rail == "" {
		text = allRailsYAML()
		data["rails"] = railOrder
	} else {
		y, ok := railYAML[rail]
		if !ok {
			f.railErr(w, raw, 404, "INVALID_RAIL", "Rail '"+rail+"' does not exist in fabric '"+fabric+"'.")
			return
		}
		text = y
		data["railName"] = rail
	}
	if raw != "" {
		writeText(w, 200, raw, text)
		return
	}
	data["yaml"] = text
	writeJSON(w, 200, okBody("Rail CIDRPools retrieved successfully.", data))
}

func (f *fakeFM) config(fabric string) map[string]any {
	if c, ok := f.cfg[fabric]; ok {
		return c
	}
	return map[string]any{"fabricName": fabric, "enableIpam": false, "enableIpAssign": false, "updatedAt": "2026-09-24T09:56:21.645313Z", "updatedBy": nil}
}

func (f *fakeFM) getConfig(w http.ResponseWriter, fabric string) {
	if _, ok := f.fabrics[fabric]; !ok {
		writeJSON(w, 404, errBody("FABRIC_NOT_FOUND", "Specified fabric does not exist."))
		return
	}
	writeJSON(w, 200, okBody("IPAM configuration retrieved successfully.", f.config(fabric)))
}

func (f *fakeFM) patchConfig(w http.ResponseWriter, fabric, force string, raw []byte) {
	ones, ok := f.fabrics[fabric]
	if !ok {
		writeJSON(w, 404, errBody("FABRIC_NOT_FOUND", "Specified fabric does not exist."))
		return
	}
	if ones {
		writeJSON(w, 409, errBody("FABRIC_NOT_EXTERNALLY_MANAGED", "IPAM configuration can only be modified for externally managed (non-ONES-controlled) fabrics."))
		return
	}
	var body map[string]any
	_ = json.Unmarshal(raw, &body)
	ei, ok1 := body["enableIpam"]
	ea, ok2 := body["enableIpAssign"]
	if !ok1 || !ok2 {
		writeJSON(w, 400, errBody("INVALID_REQUEST_BODY", "enableIpam and enableIpAssign are required."))
		return
	}
	enable, b1 := ei.(bool)
	assign, b2 := ea.(bool)
	if !b1 || !b2 {
		writeJSON(w, 400, errBody("INVALID_REQUEST_BODY", "enableIpam and enableIpAssign must be booleans."))
		return
	}
	if assign && !enable {
		writeJSON(w, 400, errBody("INVALID_IPAM_CONFIG", "enableIpAssign requires enableIpam to be true."))
		return
	}
	active := 0
	for k := range f.pools {
		if k[0] == fabric {
			active++
		}
	}
	if !enable && !strings.EqualFold(force, "true") && active > 0 {
		writeJSON(w, 409, errBody("IPAM_HAS_ACTIVE_POOLS", "Active IPAM pools exist for this fabric."))
		return
	}
	if !enable && strings.EqualFold(force, "true") {
		for k := range f.pools {
			if k[0] == fabric {
				delete(f.pools, k)
			}
		}
	}
	c := map[string]any{"fabricName": fabric, "enableIpam": enable, "enableIpAssign": assign, "updatedAt": "2026-09-24T08:57:28.254792109Z", "updatedBy": "Aviz Admin"}
	f.cfg[fabric] = c
	writeJSON(w, 200, okBody("IPAM configuration updated successfully.", c))
}

func (f *fakeFM) findAlloc(fabric, tenant, server string) int {
	for i, a := range f.allocs {
		if a["fabricName"] == fabric && a["tenantName"] == tenant && a["serverName"] == server {
			return i
		}
	}
	return -1
}

func (f *fakeFM) patchTenant(w http.ResponseWriter, fabric, tenant string, raw []byte) {
	var body struct {
		Operation string `json:"operation"`
		Servers   []struct {
			ServerName string `json:"serverName"`
		} `json:"servers"`
	}
	_ = json.Unmarshal(raw, &body)
	if len(body.Servers) == 0 {
		writeJSON(w, 400, errBody("INVALID_SERVER_LIST", "Server list is null or empty."))
		return
	}
	if body.Operation != "ADD" && body.Operation != "DELETE" {
		writeJSON(w, 400, errBody("INVALID_OPERATION", "Operation field is missing or invalid."))
		return
	}
	if !f.tenants[fabric][tenant] {
		writeJSON(w, 404, errBody("TENANT_NOT_FOUND", "Tenant does not exist."))
		return
	}
	for _, s := range body.Servers {
		if !f.servers[fabric][s.ServerName] {
			writeJSON(w, 404, errBody("SERVER_NOT_FOUND", "One or more servers not found."))
			return
		}
	}
	c := f.config(fabric)
	if c["enableIpam"] != true {
		writeJSON(w, 200, okBody("GPU update processed successfully.", nil))
		return
	}
	partial := false
	out := []map[string]any{}
	for _, s := range body.Servers {
		n := s.ServerName
		i := f.findAlloc(fabric, tenant, n)
		if body.Operation == "ADD" {
			if n == "full-node" {
				out = append(out, map[string]any{"serverName": n, "status": "IP_POOL_EXHAUSTED", "ipAddress": nil, "ipAssignStatus": "IP_POOL_EXHAUSTED", "message": "No free address in the tenant pool."})
				partial = true
				continue
			}
			if i < 0 {
				status := "ALLOCATED_ONLY"
				if c["enableIpAssign"] == true {
					status = "ASSIGNED"
				}
				f.allocs = append(f.allocs, map[string]any{
					"id": f.nextID, "poolId": 1, "fabricName": fabric, "tenantName": tenant, "tenantId": 3, "serverName": n,
					"ipAddress": fmt.Sprintf("192.3.0.%d", f.nextHost), "prefixLength": 32, "assignStatus": status,
					"cleanupStatus": "NONE", "allocatedAt": "2026-09-24T09:56:21.645313Z",
				})
				f.nextID++
				f.nextHost++
				i = len(f.allocs) - 1
			}
			a := f.allocs[i]
			out = append(out, map[string]any{"serverName": n, "status": "SUCCESS", "ipAddress": fmt.Sprintf("%v/%v", a["ipAddress"], a["prefixLength"]), "ipAssignStatus": a["assignStatus"]})
			continue
		}
		if i < 0 {
			out = append(out, map[string]any{"serverName": n, "status": "FAILED", "ipAddress": nil, "ipAssignStatus": "NOT_ALLOCATED", "message": "Server has no active IP allocation."})
			partial = true
			continue
		}
		a := f.allocs[i]
		f.allocs = append(f.allocs[:i], f.allocs[i+1:]...)
		out = append(out, map[string]any{"serverName": n, "status": "SUCCESS", "ipAddress": fmt.Sprintf("%v/%v", a["ipAddress"], a["prefixLength"]), "ipAssignStatus": a["assignStatus"], "released": true})
	}
	data := map[string]any{"fabricName": fabric, "tenantName": tenant, "operation": body.Operation, "servers": out, "partial": partial}
	status := 200
	if partial {
		status = 207
	}
	writeJSON(w, status, okBody("GPU update processed successfully.", data))
}

func (f *fakeFM) computeIP(w http.ResponseWriter, fabric, tenant, server string) {
	i := f.findAlloc(fabric, tenant, server)
	if i < 0 {
		if !f.pools[[2]string{fabric, tenant}] {
			writeJSON(w, 404, errBody("IPAM_NOT_CONFIGURED", "No IPAM pool exists for the tenant."))
			return
		}
		writeJSON(w, 404, errBody("SERVER_NOT_ALLOCATED_TO_TENANT", "Server has no IP allocation."))
		return
	}
	writeJSON(w, 200, okBody("IP allocation retrieved successfully.", f.allocs[i]))
}

func (f *fakeFM) tenantIPs(w http.ResponseWriter, fabric, tenant, pageS, sizeS, status string) {
	if !f.pools[[2]string{fabric, tenant}] {
		writeJSON(w, 404, errBody("IPAM_NOT_CONFIGURED", "No IPAM pool exists for the tenant."))
		return
	}
	if status != "" && status != "ALLOCATED_ONLY" && status != "ASSIGNED" && status != "IP_ASSIGN_FAILED" {
		writeJSON(w, 400, errBody("BAD_REQUEST", "Invalid status '"+status+"'."))
		return
	}
	page, size := 0, 50
	if pageS != "" {
		page, _ = strconv.Atoi(pageS)
	}
	if sizeS != "" {
		size, _ = strconv.Atoi(sizeS)
	}
	if page < 0 {
		page = 0
	}
	if size < 1 {
		size = 1
	}
	rows := []map[string]any{}
	for _, a := range f.allocs {
		if a["fabricName"] == fabric && a["tenantName"] == tenant && (status == "" || a["assignStatus"] == status) {
			rows = append(rows, a)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i]["id"].(int) < rows[j]["id"].(int) })
	total := len(rows)
	lo, hi := page*size, (page+1)*size
	if lo > total {
		lo = total
	}
	if hi > total {
		hi = total
	}
	data := map[string]any{
		"fabricName": fabric, "tenantName": tenant,
		"pool":        map[string]any{"cidr": "192.3.0.0/17", "allocatableAddressCount": 32765, "allocatedCount": len(f.allocs)},
		"allocations": rows[lo:hi],
		"pagination":  map[string]any{"page": page, "size": size, "totalElements": total, "totalPages": (total + size - 1) / size},
	}
	writeJSON(w, 200, okBody("IP allocations retrieved successfully.", data))
}

func (f *fakeFM) serverIPs(w http.ResponseWriter, fabric, server string) {
	if !f.servers[fabric][server] {
		writeJSON(w, 404, errBody("SERVER_NOT_FOUND", "Server was not found in the fabric inventory."))
		return
	}
	rows := []map[string]any{}
	for _, a := range f.allocs {
		if a["fabricName"] == fabric && a["serverName"] == server {
			rows = append(rows, a)
		}
	}
	writeJSON(w, 200, okBody("IP allocations retrieved successfully.", map[string]any{"fabricName": fabric, "serverName": server, "allocations": rows}))
}

func urlUnescape(s string) (string, error) { return url.PathUnescape(s) }
