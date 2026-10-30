package ones

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
)

// fmStub is a fake Fabric Manager: it records every request and answers from
// a path -> response table (404 for anything not listed).
type fmStub struct {
	t         *testing.T
	mu        sync.Mutex
	requests  []stubRequest
	responses map[string]stubResponse // key: "METHOD /path"
}

type stubRequest struct {
	Method string
	Path   string
	Query  string
	Body   map[string]interface{}
}

type stubResponse struct {
	status int
	body   string
}

func newFMStub(t *testing.T, responses map[string]stubResponse) (*fmStub, *Client) {
	t.Helper()
	s := &fmStub{t: t, responses: responses}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := stubRequest{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &req.Body); err != nil {
				t.Errorf("%s %s: body is not a JSON object: %s", r.Method, r.URL.Path, raw)
			}
		}
		s.mu.Lock()
		s.requests = append(s.requests, req)
		s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		resp, ok := s.responses[r.Method+" "+r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"status":"error","message":"no stub for `+r.Method+" "+r.URL.Path+`"}`)
			return
		}
		w.WriteHeader(resp.status)
		_, _ = io.WriteString(w, resp.body)
	}))
	t.Cleanup(srv.Close)
	c := newTestClient(t, srv.URL, "tok")
	// Routes the gateway does not proxy (the VF endpoints) go straight to FM,
	// served at the root; point those at the stub too.
	c.GetTransport().SetFMDirectURL(srv.URL)
	return s, c
}

func (s *fmStub) all() []stubRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]stubRequest(nil), s.requests...)
}

func (s *fmStub) last() stubRequest {
	r := s.all()
	if len(r) == 0 {
		s.t.Fatal("no request was sent")
	}
	return r[len(r)-1]
}

func ok(body string) stubResponse { return stubResponse{http.StatusOK, body} }

const msgOK = `{"status":"success","message":"done"}`

// --- validation --------------------------------------------------------------

func TestValidationRejectsBeforeAnyRequest(t *testing.T) {
	stub, c := newFMStub(t, nil)
	ctx := context.Background()
	add, bogus := GpuActionAdd, GpuAction("ADDD")
	scope := ConfigScope("ALL")
	oneServer := []GpuServerInfo{{ServerName: "h0"}}
	suid := SuidMap{"0": {"h0": {Gpus: []string{"G0"}}}}

	cases := []struct {
		name  string
		field string
		call  func() error
	}{
		{"create: empty fabric", "fabricName", func() error { _, err := CreateTenant(ctx, c, "", "t", nil, nil, nil); return err }},
		{"create: empty tenant", "tenantName", func() error { _, err := CreateTenant(ctx, c, "f", " ", nil, nil, nil); return err }},
		{"create: quota 0", "maxGpusAllowed", func() error { _, err := CreateTenant(ctx, c, "f", "t", nil, Ptr(0), nil); return err }},
		{"create async: quota -2", "maxGpusAllowed", func() error { _, err := CreateTenantAsync(ctx, c, "f", "t", nil, Ptr(-2), nil); return err }},
		{"create: slash in name", "tenantName", func() error { _, err := CreateTenant(ctx, c, "f", "a/b", nil, nil, nil); return err }},
		{"get: empty tenant", "tenantName", func() error { _, err := GetTenant(ctx, c, "f", ""); return err }},
		{"list: empty fabric", "fabricName", func() error { _, err := ListTenants(ctx, c, ""); return err }},
		{"delete async: dot-dot", "tenantName", func() error { _, err := DeleteTenantAsync(ctx, c, "f", ".."); return err }},
		{"update: no servers", "servers", func() error { _, err := UpdateTenant(ctx, c, "f", "t", nil, &add); return err }},
		{"update async: bad action", "gpuAction", func() error { _, err := UpdateTenantAsync(ctx, c, "f", "t", oneServer, &bogus); return err }},
		{"modify: bad scope", "configScope", func() error {
			_, err := ModifyGpuAllocations(ctx, c, "f", "t", suid, &add, &scope, nil)
			return err
		}},
		{"modify async: empty suid", "suid", func() error {
			_, err := ModifyGpuAllocationsAsync(ctx, c, "f", "t", SuidMap{}, &add, nil, nil)
			return err
		}},
		{"ports: bad action", "operation", func() error { _, err := AssignGpuPorts(ctx, c, "f", "t", bogus, []string{"h0"}, nil, nil); return err }},
		{"ports: negative id", "gpuIds[0]", func() error {
			_, err := AssignGpuPorts(ctx, c, "f", "t", add, []string{"h0"}, []int{-1}, nil)
			return err
		}},
		{"inventory: empty fabric", "fabricName", func() error { _, err := GetInventoryPorts(ctx, c, ""); return err }},
		{"nmxc inventory: empty fabric", "fabricName", func() error { _, err := GetNmxcInventory(ctx, c, ""); return err }},
		{"vf list: empty server", "serverName", func() error { _, err := ListVfInterfaces(ctx, c, "f", ""); return err }},
		{"vf assign: empty tenant", "tenantName", func() error { _, err := AssignVfInterface(ctx, c, "f", "h0", "pf0vf0_if", ""); return err }},
		{"vf unassign: slash in vf", "vfId", func() error { _, err := UnassignVfInterface(ctx, c, "f", "h0", "pf0/vf0"); return err }},
		{"peering: empty vpcname", "vpcname", func() error { _, err := DeleteVpcPeering(ctx, c, "f", "p", "", "peer"); return err }},
		{"with: no action", "action", func() error {
			_, err := UpdateTenantWith(ctx, c, UpdateTenantServersRequest{FabricName: "f", TenantName: "t", Servers: oneServer})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("expected *ValidationError, got %T (%v)", err, err)
			}
			if ve.Field != tc.field {
				t.Errorf("Field = %q, want %q (%v)", ve.Field, tc.field, err)
			}
			var onesErr ONESError
			if !errors.As(err, &onesErr) {
				t.Error("ValidationError should satisfy ONESError")
			}
		})
	}
	if n := len(stub.all()); n != 0 {
		t.Fatalf("validation failures must not reach the server; %d request(s) were sent", n)
	}
}

func TestValidInputStillReachesServer(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		"POST /api/fm/fabrics/f/tenants": ok(`{"status":"success","message":"created","data":"tenant t created"}`),
	})
	if _, err := CreateTenant(context.Background(), c, "f", "t", nil, Ptr(UnlimitedGPUs), nil); err != nil {
		t.Fatalf("unlimited quota rejected: %v", err)
	}
	if got := stub.last().Body["maxGpusAllowed"]; got != float64(-1) {
		t.Errorf("maxGpusAllowed sent as %v, want -1", got)
	}
}

// --- typed inventory ---------------------------------------------------------

func TestInventoryTypedAndRaw(t *testing.T) {
	ports := `{"fabricName":"U1","success":true,"hostCount":1,"totalPortCount":1,"availablePortCount":0,
		"allocatedPortCount":1,"hosts":[{"hostname":"h0","ufmSystemName":"SYS0","portCount":1,
		"ports":[{"portGuid":"0x1","allocated":true,"tenantId":3,"tenantName":"t1"}]}]}`
	stub, c := newFMStub(t, map[string]stubResponse{
		"GET /api/fm/fabrics/U1/inventoryPorts":   ok(ports),
		"GET /api/fm/fabrics/U1/inventoryHosts":   ok(`{"success":true,"hostCount":1,"source":"UFM_LIVE","hosts":[{"hostname":"SYS0","portCount":8}]}`),
		"POST /api/fm/fabrics/U1/inventorySync":   ok(`{"fabricName":"U1","success":true,"message":"Inventory sync completed successfully"}`),
		"GET /api/fm/fabrics/N1/nmxc/inventory":   ok(`{"fabricType":"NMXC","success":true,"overallHealth":"HEALTHY","summary":{"domainCount":1},"domains":[{"domainId":7,"domainName":"d"}]}`),
		"GET /api/fm/fabrics/U404/inventoryPorts": {http.StatusNotFound, `{"fabricName":"U404","success":false,"error":"FABRIC_NOT_FOUND","message":"Fabric not found: U404"}`},
	})
	ctx := context.Background()

	p, err := GetInventoryPorts(ctx, c, "U1")
	if err != nil {
		t.Fatal(err)
	}
	if h, found := p.Host("SYS0"); !found || *h.Ports[0].TenantID != 3 {
		t.Errorf("Ports decoded wrong: %+v", p)
	}
	raw, err := GetInventoryPortsRaw(ctx, c, "U1")
	if err != nil || raw["hostCount"] != float64(1) {
		t.Errorf("PortsRaw = %v, %v", raw, err)
	}

	h, err := GetInventoryHosts(ctx, c, "U1")
	if err != nil || h.Hosts[0].PortCount != 8 || h.Source != "UFM_LIVE" {
		t.Errorf("UFMHosts = %+v, %v", h, err)
	}
	s, err := InventorySync(ctx, c, "U1")
	if err != nil || !s.Success || stub.last().Method != http.MethodPost {
		t.Errorf("Sync = %+v, %v", s, err)
	}
	n, err := GetNmxcInventory(ctx, c, "N1")
	if err != nil || n.OverallHealth != "HEALTHY" || n.Domains[0].DomainID.String() != "7" {
		t.Errorf("NMXC inventory = %+v, %v", n, err)
	}

	_, err = GetInventoryPorts(ctx, c, "U404")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected *NotFoundError, got %T (%v)", err, err)
	}
	// The transport reports the "error" key as the message; the full body,
	// including the human-readable message, stays in ResponseBody.
	body, _ := nf.ResponseBody.(map[string]interface{})
	if nf.Message != "FABRIC_NOT_FOUND" || body["message"] != "Fabric not found: U404" {
		t.Errorf("NotFoundError = %q, body %v", nf.Message, body)
	}
}

// --- HBN VF interfaces -------------------------------------------------------

func TestVfInterfaceEndpoints(t *testing.T) {
	// A direct-FM route: no /api/fm prefix.
	const base = "/fabrics/f/servers/h0/vf-interfaces"
	stub, c := newFMStub(t, map[string]stubResponse{
		"GET " + base: ok(`{"serverName":"h0","fabricName":"f","dpuCount":1,"dpus":[{"dpuName":"dpu0",
			"vfInterfaces":[{"ifName":"pf0vf0_if","serverIf":"enp1s0","status":"FREE"},
			                {"ifName":"pf0vf1_if","tenantName":"t1","status":"ASSIGNED","vlanId":100}]}]}`),
		"POST " + base + "/pf0vf0_if/assign": ok(`{"status":"success","message":"assigned","data":{"vfId":"pf0vf0_if",
			"serverName":"h0","fabricName":"f","tenantName":"t1","status":"ASSIGNED","vlanId":100,"ipAddress":"10.1.0.5"}}`),
		"DELETE " + base + "/pf0vf0_if/assign": ok(msgOK),
		"POST " + base + "/pf0vf9_if/assign":   {http.StatusNotFound, `{"status":"error","message":"VF pf0vf9_if not found","error":{"code":"VF_NOT_FOUND"}}`},
		"DELETE " + base + "/pf0vf1_if/assign": {http.StatusConflict, `{"status":"error","message":"server left the tenant"}`},
	})
	ctx := context.Background()

	list, err := ListVfInterfaces(ctx, c, "f", "h0")
	if err != nil {
		t.Fatal(err)
	}
	if r := stub.last(); r.Method != http.MethodGet || r.Path != base || r.Body != nil {
		t.Errorf("list request = %+v", r)
	}
	if vfs := list.Flatten(); len(vfs) != 2 || vfs[1].TenantName != "t1" || vfs[1].VlanID != 100 {
		t.Errorf("list decoded wrong: %+v", list)
	}

	a, err := AssignVfInterface(ctx, c, "f", "h0", "pf0vf0_if", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if r := stub.last(); r.Method != http.MethodPost || r.Path != base+"/pf0vf0_if/assign" ||
		!reflect.DeepEqual(r.Body, map[string]interface{}{"tenantName": "t1"}) {
		t.Errorf("assign request = %+v", r)
	}
	if a.TenantName != "t1" || a.IPAddress != "10.1.0.5" {
		t.Errorf("assign decoded wrong: %+v", a)
	}

	u, err := UnassignVfInterface(ctx, c, "f", "h0", "pf0vf0_if")
	if err != nil || u.Status != "success" {
		t.Fatalf("unassign = %+v, %v", u, err)
	}
	if r := stub.last(); r.Method != http.MethodDelete || r.Body != nil {
		t.Errorf("unassign request = %+v", r)
	}

	_, err = AssignVfInterface(ctx, c, "f", "h0", "pf0vf9_if", "t1")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Errorf("expected *NotFoundError, got %T (%v)", err, err)
	}
	_, err = UnassignVfInterface(ctx, c, "f", "h0", "pf0vf1_if")
	var cf *ConflictError
	if !errors.As(err, &cf) || cf.Message != "server left the tenant" {
		t.Errorf("expected *ConflictError, got %T (%v)", err, err)
	}
}

// --- request structs ---------------------------------------------------------

// TestWithMatchesPositional checks each ...With method sends exactly what its
// positional twin sends.
func TestWithMatchesPositional(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		"POST /api/fm/fabrics/f/tenants":                  ok(`{"status":"success","message":"ok","data":"created"}`),
		"PATCH /api/fm/fabrics/f/tenants/t":               ok(msgOK),
		"POST /api/fm/fabrics/f/tenants/t/gpus":           ok(`{"success":true,"serversProcessed":1}`),
		"POST /api/fm/fabrics/f/tenants/t/gpuAllocations": ok(msgOK),
	})
	ctx := context.Background()
	del := GpuActionDelete
	scope := ConfigScopeParticularGPU
	suid := SuidMap{"0": {"h0": {Gpus: []string{"G0", "G1"}}}}

	pairs := []struct {
		name       string
		positional func() error
		with       func() error
	}{
		{"create",
			func() error { _, err := CreateTenant(ctx, c, "f", "t", Ptr("demo"), Ptr(8), Ptr(false)); return err },
			func() error {
				_, err := CreateTenantWith(ctx, c, CreateTenantRequest{FabricName: "f", TenantName: "t",
					Description: "demo", MaxGpusAllowed: 8, Shared: Ptr(false)})
				return err
			}},
		{"create, optionals omitted",
			func() error { _, err := CreateTenant(ctx, c, "f", "t", nil, nil, nil); return err },
			func() error {
				_, err := CreateTenantWith(ctx, c, CreateTenantRequest{FabricName: "f", TenantName: "t"})
				return err
			}},
		{"update",
			func() error {
				_, err := UpdateTenant(ctx, c, "f", "t", []GpuServerInfo{{ServerName: "h0"}}, &del)
				return err
			},
			func() error {
				_, err := UpdateTenantWith(ctx, c, UpdateTenantServersRequest{FabricName: "f", TenantName: "t",
					Servers: []GpuServerInfo{{ServerName: "h0"}}, Action: GpuActionDelete})
				return err
			}},
		{"assign ports",
			func() error {
				_, err := AssignGpuPorts(ctx, c, "f", "t", GpuActionAdd, []string{"h0"}, nil, nil)
				return err
			},
			func() error {
				_, err := AssignGpuPortsWith(ctx, c, AssignPortsRequest{FabricName: "f", TenantName: "t",
					Action: GpuActionAdd, ServerNames: []string{"h0"}})
				return err
			}},
		{"modify allocations",
			func() error { _, err := ModifyGpuAllocations(ctx, c, "f", "t", suid, &del, &scope, nil); return err },
			func() error {
				_, err := ModifyGpuAllocationsWith(ctx, c, ModifyAllocationsRequest{FabricName: "f", TenantName: "t",
					Suid: suid, Action: GpuActionDelete, Scope: ConfigScopeParticularGPU})
				return err
			}},
	}
	for _, p := range pairs {
		t.Run(p.name, func(t *testing.T) {
			if err := p.positional(); err != nil {
				t.Fatalf("positional: %v", err)
			}
			want := stub.last()
			if err := p.with(); err != nil {
				t.Fatalf("With: %v", err)
			}
			if got := stub.last(); !reflect.DeepEqual(got, want) {
				t.Errorf("With sent %+v\npositional sent %+v", got, want)
			}
		})
	}

	// The omitted case must really omit the optionals.
	stub.responses["POST /api/fm/fabrics/f/tenants"] = ok(`{"status":"success","message":"ok","data":"created"}`)
	if _, err := CreateTenantWith(ctx, c, CreateTenantRequest{FabricName: "f", TenantName: "t"}); err != nil {
		t.Fatal(err)
	}
	if b := stub.last().Body; !reflect.DeepEqual(b, map[string]interface{}{"tenantName": "t"}) {
		t.Errorf("zero-value optionals were sent: %v", b)
	}

	// PATCH /tenants/{t} takes the action as "operation" (FM's UpdateTenant
	// body); "gpuAction" belongs to gpuAllocations and FM rejects it here
	// with 400 INVALID_OPERATION.
	upd := UpdateTenantServersRequest{FabricName: "f", TenantName: "t",
		Servers: []GpuServerInfo{{ServerName: "h0"}}, Action: GpuActionAdd}
	if _, err := UpdateTenantWith(ctx, c, upd); err != nil {
		t.Fatal(err)
	}
	if b := stub.last().Body; b["operation"] != "ADD" || b["gpuAction"] != nil {
		t.Errorf("update sent the action under the wrong key: %v", b)
	}
}

// FM reports a failed operation as FAILED, not FAILURE. The wait must stop on
// the first poll and hand back the operation, not poll out its timeout.
func TestWaitForOperationStopsOnFailed(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		"GET /api/fm/operations/op1": ok(`{"id":"op1","type":"GPU_ALLOCATE","status":"FAILED","progress":10,` +
			`"errorMessage":"Server must be attached to the tenant before assigning GPUs."}`),
	})
	start := time.Now()
	op, err := c.Operations.WaitForOperation(context.Background(), "op1",
		resources.WaitOptions{Interval: 50 * time.Millisecond, Timeout: 2 * time.Second})
	if err != nil {
		t.Fatalf("want the failed operation and a nil error, got %v", err)
	}
	if op == nil || op.Status != "FAILED" || op.ErrorMessage == "" {
		t.Fatalf("got %+v", op)
	}
	if n := len(stub.all()); n != 1 || time.Since(start) > time.Second {
		t.Errorf("polled %d times over %s; want one poll", n, time.Since(start))
	}
}
