package ones_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

var ctx = context.Background()

func newClient(t *testing.T) (*ones.Client, *fakeFM) {
	t.Helper()
	fm := newFakeFM(t)
	c, err := ones.InitializeWithToken(fm.srv.URL, testToken)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c, fm
}

func ipamOn(t *testing.T, c *ones.Client) {
	t.Helper()
	if _, err := ones.UpdateIpamConfig(ctx, c, "ns", true, true, nil); err != nil {
		t.Fatal(err)
	}
}

func addServers(t *testing.T, c *ones.Client, names ...string) *ones.GpuUpdateResult {
	t.Helper()
	var servers []ones.GpuServerInfo
	for _, n := range names {
		servers = append(servers, ones.GpuServerInfo{ServerName: n})
	}
	out, err := ones.UpdateTenant(ctx, c, "ns", "gorax", servers, ones.Ptr(ones.GpuActionAdd))
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func wantErr[E error](t *testing.T, err error) E {
	t.Helper()
	var target E
	if !errors.As(err, &target) {
		t.Fatalf("want error of type %T, got %T: %v", target, err, err)
	}
	return target
}

func errorCode(t *testing.T, body any) string {
	t.Helper()
	m, _ := body.(map[string]any)
	e, _ := m["error"].(map[string]any)
	code, _ := e["code"].(string)
	return code
}

// ── rails ────────────────────────────────────────────────────────────────────

func TestRailCidrPoolsRawRepresentationsReturnManifestVerbatim(t *testing.T) {
	c, fm := newClient(t)
	for _, accept := range []ones.RailAccept{ones.RailAcceptYAML, ones.RailAcceptXYAML, ones.RailAcceptText} {
		a := accept
		got, err := ones.GetRailCidrPools(ctx, c, "yml", &a)
		if err != nil {
			t.Fatalf("%s: %v", accept, err)
		}
		if got.YAML != allRailsYAML() || got.FabricName != "yml" || got.Rails != nil {
			t.Errorf("%s: unexpected result %+v", accept, got)
		}
		if h := fm.last().Header.Get("Accept"); h != string(accept) {
			t.Errorf("Accept header = %q, want %q", h, accept)
		}
	}
}

func TestRailCidrPoolsNilAcceptDefaultsToYAML(t *testing.T) {
	c, fm := newClient(t)
	got, err := ones.GetRailCidrPools(ctx, c, "yml", nil)
	if err != nil || got.YAML != allRailsYAML() {
		t.Fatalf("got %+v, %v", got, err)
	}
	if h := fm.last().Header.Get("Accept"); h != "application/yaml" {
		t.Errorf("Accept = %q", h)
	}
}

func TestRailCidrPoolsJSONEnvelopeListsRails(t *testing.T) {
	c, _ := newClient(t)
	got, err := ones.GetRailCidrPools(ctx, c, "yml", ones.Ptr(ones.RailAcceptJSON))
	if err != nil {
		t.Fatal(err)
	}
	want := &ones.RailCidrPools{FabricName: "yml", Rails: []string{"r0", "r2"}, YAML: allRailsYAML()}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestRailCidrPoolSingle(t *testing.T) {
	c, _ := newClient(t)
	raw, err := ones.GetRailCidrPool(ctx, c, "yml", "r0", ones.Ptr(ones.RailAcceptText))
	if err != nil {
		t.Fatal(err)
	}
	if want := (&ones.RailCidrPool{FabricName: "yml", RailName: "r0", YAML: railYAML["r0"]}); !reflect.DeepEqual(raw, want) {
		t.Errorf("raw: got %+v", raw)
	}
	js, err := ones.GetRailCidrPool(ctx, c, "yml", "r2", ones.Ptr(ones.RailAcceptJSON))
	if err != nil {
		t.Fatal(err)
	}
	if want := (&ones.RailCidrPool{FabricName: "yml", RailName: "r2", YAML: railYAML["r2"]}); !reflect.DeepEqual(js, want) {
		t.Errorf("json: got %+v", js)
	}
}

func TestRailErrorsArePlainTextForYAMLAndStructuredForJSON(t *testing.T) {
	c, _ := newClient(t)

	_, err := ones.GetRailCidrPool(ctx, c, "yml", "r9", ones.Ptr(ones.RailAcceptText))
	nf := wantErr[*ones.NotFoundError](t, err)
	if want := "INVALID_RAIL: Rail 'r9' does not exist in fabric 'yml'."; nf.Message != want {
		t.Errorf("message = %q, want %q", nf.Message, want)
	}
	if nf.StatusCode != 404 {
		t.Errorf("status = %d", nf.StatusCode)
	}

	_, err = ones.GetRailCidrPool(ctx, c, "yml", "r9", ones.Ptr(ones.RailAcceptJSON))
	nf = wantErr[*ones.NotFoundError](t, err)
	if code := errorCode(t, nf.ResponseBody); code != "INVALID_RAIL" {
		t.Errorf("code = %q", code)
	}

	for _, accept := range []ones.RailAccept{ones.RailAcceptYAML, ones.RailAcceptJSON} {
		a := accept
		_, err = ones.GetRailCidrPools(ctx, c, "nope", &a)
		wantErr[*ones.NotFoundError](t, err)
		_, err = ones.GetRailCidrPools(ctx, c, "norails", &a)
		wantErr[*ones.NotFoundError](t, err)
	}
}

func TestRailInvalidAcceptRejectedBeforeRequest(t *testing.T) {
	c, fm := newClient(t)
	for _, bad := range []ones.RailAccept{"application/xml", "*/*", "", "APPLICATION/YAML"} {
		b := bad
		_, err := ones.GetRailCidrPools(ctx, c, "yml", &b)
		if v := wantErr[*ones.ValidationError](t, err); v.Field != "accept" {
			t.Errorf("field = %q", v.Field)
		}
		_, err = ones.GetRailCidrPool(ctx, c, "yml", "r0", &b)
		wantErr[*ones.ValidationError](t, err)
	}
	if n := len(fm.Requests); n != 0 {
		t.Errorf("%d requests sent despite invalid accept", n)
	}
}

// ── ipam config ──────────────────────────────────────────────────────────────

func TestIpamConfigUnconfiguredAndRoundTrip(t *testing.T) {
	c, _ := newClient(t)
	cfg, err := ones.GetIpamConfig(ctx, c, "ns")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.EnableIpam || cfg.EnableIpAssign || cfg.UpdatedBy != nil {
		t.Errorf("unconfigured fabric reported %+v", cfg)
	}

	out, err := ones.UpdateIpamConfig(ctx, c, "ns", true, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !out.EnableIpam || !out.EnableIpAssign || out.FabricName != "ns" || out.UpdatedBy == nil || *out.UpdatedBy != "Aviz Admin" {
		t.Errorf("update result %+v", out)
	}
	if out.UpdatedAt != "2026-09-24T08:57:28.254792109Z" {
		t.Errorf("updatedAt = %q", out.UpdatedAt)
	}
	got, _ := ones.GetIpamConfig(ctx, c, "ns")
	if !got.EnableIpam || !got.EnableIpAssign {
		t.Errorf("get after update: %+v", got)
	}
}

func TestIpamConfigErrors(t *testing.T) {
	c, _ := newClient(t)
	_, err := ones.GetIpamConfig(ctx, c, "nope")
	wantErr[*ones.NotFoundError](t, err)
	_, err = ones.UpdateIpamConfig(ctx, c, "nope", true, false, nil)
	wantErr[*ones.NotFoundError](t, err)

	_, err = ones.UpdateIpamConfig(ctx, c, "ones-ctl", true, false, nil)
	ce := wantErr[*ones.ConflictError](t, err)
	if code := errorCode(t, ce.ResponseBody); code != "FABRIC_NOT_EXTERNALLY_MANAGED" {
		t.Errorf("code = %q", code)
	}
}

func TestDisablingIpamWithActivePoolsNeedsForce(t *testing.T) {
	c, fm := newClient(t)
	ipamOn(t, c)

	_, err := ones.UpdateIpamConfig(ctx, c, "ns", false, false, nil)
	ce := wantErr[*ones.ConflictError](t, err)
	if code := errorCode(t, ce.ResponseBody); code != "IPAM_HAS_ACTIVE_POOLS" {
		t.Errorf("code = %q", code)
	}
	if cfg, _ := ones.GetIpamConfig(ctx, c, "ns"); !cfg.EnableIpam {
		t.Error("config changed despite the 409")
	}
	if q := fm.Requests[len(fm.Requests)-2].Query; len(q) != 0 {
		t.Errorf("force must be omitted when nil, got %v", q)
	}

	out, err := ones.UpdateIpamConfig(ctx, c, "ns", false, false, ones.Ptr(true))
	if err != nil || out.EnableIpam {
		t.Fatalf("forced disable: %+v, %v", out, err)
	}
	forced := fm.Requests[len(fm.Requests)-1]
	if got := forced.Query["force"]; len(got) != 1 || got[0] != "true" {
		t.Errorf("force query = %v", forced.Query)
	}
	var body map[string]bool
	if err := json.Unmarshal([]byte(forced.Body), &body); err != nil || len(body) != 2 || body["enableIpam"] || body["enableIpAssign"] {
		t.Errorf("body = %q", forced.Body)
	}
}

func TestUpdateIpamConfigAssignRequiresIpamLocally(t *testing.T) {
	c, fm := newClient(t)
	_, err := ones.UpdateIpamConfig(ctx, c, "ns", false, true, nil)
	if v := wantErr[*ones.ValidationError](t, err); v.Field != "enableIpAssign" {
		t.Errorf("field = %q", v.Field)
	}
	if len(fm.Requests) != 0 {
		t.Error("request was sent")
	}
}

// ── tenant update + allocations ──────────────────────────────────────────────

func TestUpdateTenantWithoutIpamKeepsPlainMessage(t *testing.T) {
	c, _ := newClient(t)
	out := addServers(t, c, "su00-node00")
	if out.Status != "success" || out.Message != "GPU update processed successfully." || len(out.Servers) != 0 || out.Partial {
		t.Errorf("got %+v", out)
	}
}

func TestUpdateTenantAddAllocatesAndAssignsIP(t *testing.T) {
	c, _ := newClient(t)
	ipamOn(t, c)
	out := addServers(t, c, "su00-node00")
	if out.FabricName != "ns" || out.TenantName != "gorax" || out.Operation != "ADD" || out.Partial {
		t.Errorf("envelope fields: %+v", out)
	}
	if len(out.Servers) != 1 {
		t.Fatalf("servers: %+v", out.Servers)
	}
	s := out.Servers[0]
	if s.ServerName != "su00-node00" || s.Status != "SUCCESS" || s.IPAddress == nil || *s.IPAddress != "192.3.0.2/32" || s.IPAssignStatus != "ASSIGNED" {
		t.Errorf("server result: %+v", s)
	}
}

func TestPartialUpdateIs207ButNotAnError(t *testing.T) {
	c, _ := newClient(t)
	ipamOn(t, c)
	out := addServers(t, c, "su00-node01", "full-node")
	if !out.Partial {
		t.Fatal("expected Partial")
	}
	by := map[string]ones.IpamServerResult{}
	for _, s := range out.Servers {
		by[s.ServerName] = s
	}
	if by["su00-node01"].Status != "SUCCESS" {
		t.Errorf("node01: %+v", by["su00-node01"])
	}
	fn := by["full-node"]
	if fn.Status != "IP_POOL_EXHAUSTED" || fn.IPAddress != nil || fn.Message == "" {
		t.Errorf("full-node: %+v", fn)
	}
}

func TestReadEndpointsReflectAllocation(t *testing.T) {
	c, _ := newClient(t)
	ipamOn(t, c)
	addServers(t, c, "su00-node00")

	one, err := ones.GetTenantComputeIp(ctx, c, "ns", "gorax", "su00-node00")
	if err != nil {
		t.Fatal(err)
	}
	if one.IPAddress != "192.3.0.2" || one.PrefixLength != 32 || one.AssignStatus != ones.IpAssignStatusAssigned ||
		one.CleanupStatus != ones.IpCleanupStatusNone || one.TenantName != "gorax" || one.ServerName != "su00-node00" ||
		one.TenantID == nil || *one.TenantID != 3 || one.InterfaceName != nil {
		t.Errorf("allocation: %+v", one)
	}

	tenant, err := ones.GetTenantIps(ctx, c, "ns", "gorax", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if tenant.Pool != (ones.TenantIpPool{CIDR: "192.3.0.0/17", AllocatableAddressCount: 32765, AllocatedCount: 1}) {
		t.Errorf("pool: %+v", tenant.Pool)
	}
	if tenant.Pagination != (ones.Pagination{Page: 0, Size: 50, TotalElements: 1, TotalPages: 1}) {
		t.Errorf("pagination: %+v", tenant.Pagination)
	}
	if len(tenant.Allocations) != 1 || !reflect.DeepEqual(tenant.Allocations[0], *one) {
		t.Errorf("allocations: %+v", tenant.Allocations)
	}

	server, err := ones.GetServerIps(ctx, c, "ns", "su00-node00")
	if err != nil {
		t.Fatal(err)
	}
	if server.FabricName != "ns" || server.ServerName != "su00-node00" || len(server.Allocations) != 1 || !reflect.DeepEqual(server.Allocations[0], *one) {
		t.Errorf("server ips: %+v", server)
	}
}

func TestTenantIpsPaginationAndStatusFilter(t *testing.T) {
	c, fm := newClient(t)
	ipamOn(t, c)
	addServers(t, c, "su00-node00", "su00-node01", "su00-node02")

	p0, err := ones.GetTenantIps(ctx, c, "ns", "gorax", ones.Ptr(0), ones.Ptr(2), nil)
	if err != nil {
		t.Fatal(err)
	}
	p1, _ := ones.GetTenantIps(ctx, c, "ns", "gorax", ones.Ptr(1), ones.Ptr(2), nil)
	ips := func(r *ones.TenantIpAllocations) []string {
		var out []string
		for _, a := range r.Allocations {
			out = append(out, a.IPAddress)
		}
		return out
	}
	if !reflect.DeepEqual(ips(p0), []string{"192.3.0.2", "192.3.0.3"}) || !reflect.DeepEqual(ips(p1), []string{"192.3.0.4"}) {
		t.Errorf("pages: %v / %v", ips(p0), ips(p1))
	}
	if p0.Pagination != (ones.Pagination{Page: 0, Size: 2, TotalElements: 3, TotalPages: 2}) {
		t.Errorf("pagination: %+v", p0.Pagination)
	}

	failed, err := ones.GetTenantIps(ctx, c, "ns", "gorax", nil, nil, ones.Ptr(ones.IpAssignStatusAssignFailed))
	if err != nil || len(failed.Allocations) != 0 || failed.Pagination.TotalElements != 0 {
		t.Errorf("failed filter: %+v, %v", failed, err)
	}
	q := fm.last().Query
	if q.Get("status") != "IP_ASSIGN_FAILED" || q.Get("page") != "" || q.Get("size") != "" {
		t.Errorf("query = %v", q)
	}
}

func TestDeleteReleasesIP(t *testing.T) {
	c, _ := newClient(t)
	ipamOn(t, c)
	addServers(t, c, "su00-node00")
	out, err := ones.UpdateTenant(ctx, c, "ns", "gorax", []ones.GpuServerInfo{{ServerName: "su00-node00"}}, ones.Ptr(ones.GpuActionDelete))
	if err != nil {
		t.Fatal(err)
	}
	s := out.Servers[0]
	if out.Partial || s.Released == nil || !*s.Released || s.IPAddress == nil || *s.IPAddress != "192.3.0.2/32" {
		t.Errorf("release: %+v", out)
	}
	srv, _ := ones.GetServerIps(ctx, c, "ns", "su00-node00")
	if len(srv.Allocations) != 0 {
		t.Errorf("allocations remain: %+v", srv.Allocations)
	}
	_, err = ones.GetTenantComputeIp(ctx, c, "ns", "gorax", "su00-node00")
	nf := wantErr[*ones.NotFoundError](t, err)
	if code := errorCode(t, nf.ResponseBody); code != "SERVER_NOT_ALLOCATED_TO_TENANT" {
		t.Errorf("code = %q", code)
	}

	// releasing a server that has nothing is partial, not an error
	out, err = ones.UpdateTenant(ctx, c, "ns", "gorax", []ones.GpuServerInfo{{ServerName: "su00-node02"}}, ones.Ptr(ones.GpuActionDelete))
	if err != nil || !out.Partial || out.Servers[0].IPAssignStatus != "NOT_ALLOCATED" {
		t.Errorf("release of unallocated: %+v, %v", out, err)
	}
}

func TestIPReadErrorCodes(t *testing.T) {
	c, _ := newClient(t)
	ipamOn(t, c)
	check := func(name string, err error, code string) {
		t.Helper()
		nf := wantErr[*ones.NotFoundError](t, err)
		if got := errorCode(t, nf.ResponseBody); got != code {
			t.Errorf("%s: code = %q, want %q", name, got, code)
		}
	}
	_, err := ones.GetTenantComputeIp(ctx, c, "ns", "gorax", "su00-node00")
	check("no allocation", err, "SERVER_NOT_ALLOCATED_TO_TENANT")
	_, err = ones.GetTenantComputeIp(ctx, c, "ns", "ghost", "su00-node00")
	check("no pool", err, "IPAM_NOT_CONFIGURED")
	_, err = ones.GetTenantIps(ctx, c, "ns", "ghost", nil, nil, nil)
	check("tenant ips", err, "IPAM_NOT_CONFIGURED")
	_, err = ones.GetServerIps(ctx, c, "ns", "no-such-server")
	check("server ips", err, "SERVER_NOT_FOUND")
}

// ── validation, escaping, wiring ─────────────────────────────────────────────

func TestInvalidArgumentsFailBeforeAnyRequest(t *testing.T) {
	c, fm := newClient(t)
	calls := map[string]func() error{
		"blank fabric (rails)":   func() error { _, e := ones.GetRailCidrPools(ctx, c, "", nil); return e },
		"blank rail":             func() error { _, e := ones.GetRailCidrPool(ctx, c, "yml", "  ", nil); return e },
		"blank fabric (config)":  func() error { _, e := ones.GetIpamConfig(ctx, c, ""); return e },
		"blank fabric (update)":  func() error { _, e := ones.UpdateIpamConfig(ctx, c, "", true, false, nil); return e },
		"blank tenant":           func() error { _, e := ones.GetTenantComputeIp(ctx, c, "ns", "", "s"); return e },
		"blank server":           func() error { _, e := ones.GetTenantComputeIp(ctx, c, "ns", "t", ""); return e },
		"blank tenant (ips)":     func() error { _, e := ones.GetTenantIps(ctx, c, "ns", "", nil, nil, nil); return e },
		"blank server (ips)":     func() error { _, e := ones.GetServerIps(ctx, c, "ns", ""); return e },
		"negative page":          func() error { _, e := ones.GetTenantIps(ctx, c, "ns", "t", ones.Ptr(-1), nil, nil); return e },
		"zero size":              func() error { _, e := ones.GetTenantIps(ctx, c, "ns", "t", nil, ones.Ptr(0), nil); return e },
		"unknown status":         func() error { return statusErr(c, "assigned") },
		"unknown status (upper)": func() error { return statusErr(c, "BOGUS") },
	}
	for name, call := range calls {
		if v := wantErr[*ones.ValidationError](t, call()); v.Field == "" || v.Reason == "" {
			t.Errorf("%s: empty validation error %+v", name, v)
		}
	}
	if n := len(fm.Requests); n != 0 {
		t.Errorf("%d requests were sent", n)
	}
}

func statusErr(c *ones.Client, s string) error {
	st := ones.IpAssignStatus(s)
	_, err := ones.GetTenantIps(ctx, c, "ns", "t", nil, nil, &st)
	return err
}

func TestPathParamsAreEscapedIntoOneSegment(t *testing.T) {
	c, fm := newClient(t)
	_, _ = ones.GetIpamConfig(ctx, c, "a/b?c#d")
	got := fm.last().Path
	if got != "/api/fm/fabrics/a%2Fb%3Fc%23d/ipam/config" {
		t.Errorf("path = %q", got)
	}
}

func TestNilClientReturnsErrNilClient(t *testing.T) {
	checks := []error{
		func() error { _, e := ones.GetRailCidrPools(ctx, nil, "f", nil); return e }(),
		func() error { _, e := ones.GetRailCidrPool(ctx, nil, "f", "r", nil); return e }(),
		func() error { _, e := ones.GetIpamConfig(ctx, nil, "f"); return e }(),
		func() error { _, e := ones.UpdateIpamConfig(ctx, nil, "f", true, true, nil); return e }(),
		func() error { _, e := ones.GetTenantComputeIp(ctx, nil, "f", "t", "s"); return e }(),
		func() error { _, e := ones.GetTenantIps(ctx, nil, "f", "t", nil, nil, nil); return e }(),
		func() error { _, e := ones.GetServerIps(ctx, nil, "f", "s"); return e }(),
	}
	for i, err := range checks {
		if !errors.Is(err, ones.ErrNilClient) {
			t.Errorf("call %d: %v", i, err)
		}
	}
}

func TestUnauthenticatedCallIsRejectedByServer(t *testing.T) {
	fm := newFakeFM(t)
	c := ones.NewClientWithAuth(fm.srv.URL, mustJWT(t))
	if c.IPAM == nil || c.Rails == nil {
		t.Fatal("NewClientWithAuth must wire the IPAM and Rails resources")
	}
}

func TestCreateTenantForwardsDhcpCidrPoolOnlyWhenSet(t *testing.T) {
	c, fm := newClient(t)
	// The fake has no POST /tenants route, so the call 404s; only the wire body matters.
	_, _ = ones.CreateTenant(ctx, c, "ns", "t1", nil, nil, nil)
	if strings.Contains(fm.last().Body, "tenantDhcpCidrPool") {
		t.Errorf("unset pool leaked into body: %s", fm.last().Body)
	}
	_, _ = ones.CreateTenant(ctx, c, "ns", "t1", nil, nil, nil, ones.WithTenantDhcpCidrPool("10.9.0.0/24"))
	var body map[string]any
	if err := json.Unmarshal([]byte(fm.last().Body), &body); err != nil || body["tenantDhcpCidrPool"] != "10.9.0.0/24" || body["tenantName"] != "t1" {
		t.Errorf("sync body = %s", fm.last().Body)
	}
	_, _ = ones.CreateTenantAsync(ctx, c, "ns", "t2", nil, nil, nil, ones.WithTenantDhcpCidrPool("10.8.0.0/24"))
	if err := json.Unmarshal([]byte(fm.last().Body), &body); err != nil || body["tenantDhcpCidrPool"] != "10.8.0.0/24" {
		t.Errorf("async body = %s", fm.last().Body)
	}
}

func mustJWT(t *testing.T) ones.AuthProvider {
	t.Helper()
	a, err := ones.NewJWTAuth("access", "refresh", "http://127.0.0.1:1/refresh")
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestCreateTenantValidatesDhcpCidrPoolBeforeAnyRequest(t *testing.T) {
	c, fm := newClient(t)
	for _, bad := range []string{"192.3.0.0", "192.3.0.0/", "192.3.0/24", "192.3.0.256/24", "a.b.c.d/24", "fd00::/64", "192.3.0.0/7", "192.3.0.0/31", "192.3.0.0/32"} {
		_, err := ones.CreateTenant(ctx, c, "ns", "t1", nil, nil, nil, ones.WithTenantDhcpCidrPool(bad))
		if v := wantErr[*ones.ValidationError](t, err); v.Field != "tenantDhcpCidrPool" {
			t.Errorf("%q: field = %q", bad, v.Field)
		}
		_, err = ones.CreateTenantAsync(ctx, c, "ns", "t1", nil, nil, nil, ones.WithTenantDhcpCidrPool(bad))
		wantErr[*ones.ValidationError](t, err)
	}
	if n := len(fm.Requests); n != 0 {
		t.Errorf("%d requests sent for invalid pools", n)
	}
	for _, good := range []string{"192.3.0.0/17", "10.0.0.0/8", "172.16.0.0/30", "192.3.0.5/17", " 192.3.0.0/24 "} {
		_, err := ones.CreateTenant(ctx, c, "ns", "t1", nil, nil, nil, ones.WithTenantDhcpCidrPool(good))
		var v *ones.ValidationError
		if errors.As(err, &v) {
			t.Errorf("%q rejected: %v", good, v)
		}
	}
}

func TestCreateTenantSendsTrimmedDhcpCidrPool(t *testing.T) {
	c, fm := newClient(t)
	_, _ = ones.CreateTenant(ctx, c, "ns", "t1", nil, nil, nil, ones.WithTenantDhcpCidrPool(" 192.3.0.0/24 "))
	if !strings.Contains(fm.last().Body, `"tenantDhcpCidrPool":"192.3.0.0/24"`) {
		t.Errorf("sync body = %s", fm.last().Body)
	}
	_, _ = ones.CreateTenantAsync(ctx, c, "ns", "t2", nil, nil, nil, ones.WithTenantDhcpCidrPool(" 10.0.0.0/24\t"))
	if !strings.Contains(fm.last().Body, `"tenantDhcpCidrPool":"10.0.0.0/24"`) {
		t.Errorf("async body = %s", fm.last().Body)
	}
}
