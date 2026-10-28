package resources

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

func ptr[T any](v T) *T { return &v }

// assertField checks err is a *ValidationError naming field.
func assertField(t *testing.T, err error, field string) {
	t.Helper()
	var ve *ones_gfx.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError for %s, got %T (%v)", field, err, err)
	}
	if ve.Field != field {
		t.Fatalf("Field = %q, want %q (message %q)", ve.Field, field, ve.Message)
	}
}

func TestPathName(t *testing.T) {
	for _, ok := range []string{"Fab101", "tenant-1", "my_tenant", "a.b", "NMX C"} {
		if err := pathName("fabricName", ok); err != nil {
			t.Errorf("pathName(%q) rejected a valid name: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "   ", ".", "..", "a/b", "a?b", "a#b", `a\b`, "a%2Fb"} {
		assertField(t, pathName("fabricName", bad), "fabricName")
	}
}

func TestMaxGpus(t *testing.T) {
	for _, ok := range []*int{nil, ptr(-1), ptr(1), ptr(64)} {
		if err := maxGpus(ok); err != nil {
			t.Errorf("maxGpus(%v) rejected: %v", *ptrOr(ok), err)
		}
	}
	for _, bad := range []int{0, -2, -100} {
		err := maxGpus(ptr(bad))
		assertField(t, err, "maxGpusAllowed")
		if !strings.Contains(err.Error(), "-1") {
			t.Errorf("message should name the accepted values: %v", err)
		}
	}
}

func ptrOr(p *int) *int {
	if p == nil {
		return ptr(0)
	}
	return p
}

func TestEnumValidators(t *testing.T) {
	assertField(t, gpuAction("action", ""), "action")
	assertField(t, gpuAction("action", "ADDD"), "action")
	if err := optGpuAction("gpuAction", nil); err != nil {
		t.Errorf("nil action should be allowed: %v", err)
	}
	assertField(t, optConfigScope(ptr(ConfigScope("ALL"))), "configScope")
	if err := optConfigScope(ptr(ConfigScopeParticularGPU)); err != nil {
		t.Error(err)
	}

	got, err := ParseConfigScope("particular_gpu")
	if err != nil || got != ConfigScopeParticularGPU {
		t.Errorf("ParseConfigScope = %q, %v", got, err)
	}
	if _, err := ParseConfigScope("ALL"); err == nil {
		t.Error("expected error")
	}
}

func TestCollectionValidators(t *testing.T) {
	assertField(t, gpuServers(nil), "servers")
	assertField(t, gpuServers([]GpuServerInfo{{ServerName: "h0"}, {ServerName: " "}}), "servers[1].serverName")
	assertField(t, nonEmptyNames("serverNames", []string{}), "serverNames")
	assertField(t, gpuIDs([]int{1, -1}), "gpuIds[1]")
	assertField(t, suidMap(SuidMap{}), "suid")
	assertField(t, suidMap(SuidMap{"x": {"h0": {}}}), "suid[x]")
	assertField(t, suidMap(SuidMap{"0": {}}), "suid[0]")
	if err := suidMap(SuidMap{"0": {"h0": {Gpus: []string{"G0"}}}}); err != nil {
		t.Error(err)
	}
	assertField(t, nmxcDomains(nil), "domains")
	assertField(t, nmxcDomains([]NmxcDomain{{Host: ""}}), "domains[0].host")
	assertField(t, newPolicyName(strings.Repeat("p", MaxNetworkPolicyNameLength+1)), "policyName")
	assertField(t, peeringFields("p1", "", "peer"), "vpcname")
}

func TestRequestValidate(t *testing.T) {
	cases := []struct {
		name  string
		err   error
		field string
	}{
		{"create: missing tenant", CreateTenantRequest{FabricName: "f"}.Validate(), "tenantName"},
		{"create: bad quota", CreateTenantRequest{FabricName: "f", TenantName: "t", MaxGpusAllowed: -5}.Validate(), "maxGpusAllowed"},
		{"update: no action", UpdateTenantServersRequest{FabricName: "f", TenantName: "t",
			Servers: []GpuServerInfo{{ServerName: "h0"}}}.Validate(), "action"},
		{"update: no servers", UpdateTenantServersRequest{FabricName: "f", TenantName: "t",
			Action: ones_gfx.GpuActionAdd}.Validate(), "servers"},
		{"modify: no action", ModifyAllocationsRequest{FabricName: "f", TenantName: "t",
			Suid: SuidMap{"0": {"h0": {}}}}.Validate(), "action"},
		{"modify: bad scope", ModifyAllocationsRequest{FabricName: "f", TenantName: "t",
			Suid: SuidMap{"0": {"h0": {}}}, Action: ones_gfx.GpuActionAdd, Scope: "ALL"}.Validate(), "configScope"},
		{"ports: no servers", AssignPortsRequest{FabricName: "f", TenantName: "t",
			Action: ones_gfx.GpuActionAdd}.Validate(), "serverNames"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { assertField(t, c.err, c.field) })
	}

	ok := []error{
		CreateTenantRequest{FabricName: "f", TenantName: "t"}.Validate(),
		CreateTenantRequest{FabricName: "f", TenantName: "t", MaxGpusAllowed: ones_gfx.UnlimitedGPUs}.Validate(),
		UpdateTenantServersRequest{FabricName: "f", TenantName: "t", Action: ones_gfx.GpuActionDelete,
			Servers: []GpuServerInfo{{ServerName: "h0"}}}.Validate(),
		AssignPortsRequest{FabricName: "f", TenantName: "t", Action: ones_gfx.GpuActionAdd,
			ServerNames: []string{"h0"}}.Validate(),
	}
	for i, err := range ok {
		if err != nil {
			t.Errorf("valid request %d rejected: %v", i, err)
		}
	}
}

func TestCreateTenantRequestArgs(t *testing.T) {
	d, m, s := CreateTenantRequest{FabricName: "f", TenantName: "t"}.args()
	if d != nil || m != nil || s != nil {
		t.Errorf("zero optionals should be omitted, got %v %v %v", d, m, s)
	}
	d, m, s = CreateTenantRequest{Description: "x", MaxGpusAllowed: 8, Shared: ptr(false)}.args()
	if *d != "x" || *m != 8 || *s != false {
		t.Errorf("got %v %v %v", *d, *m, *s)
	}
}

const inventoryPortsJSON = `{
  "fabricName": "UFM1", "success": true, "hostCount": 1, "totalPortCount": 2,
  "availablePortCount": 1, "allocatedPortCount": 1, "newSummaryKey": 7,
  "hosts": [{
    "hostname": "hgx-su00-h00", "ufmSystemName": "SA-ZT-NDR-02", "ipAddress": null,
    "systemGuid": "0xabc", "portCount": 2, "availableCount": 1, "allocatedCount": 1,
    "lastConfiguredTenant": "t1",
    "ports": [
      {"portGuid": "0x1", "dname": "mlx5_0", "portName": "p1", "portState": "Active",
       "linkSpeedActive": "NDR", "tenantId": 3, "tenantName": "t1", "pkeyValue": "0x8001",
       "allocated": true, "allocatedAt": "2026-04-18T14:03:32.632113", "lastSyncedAt": null},
      {"portGuid": "0x2", "allocated": false, "tenantId": null, "tenantName": null}
    ]
  }]
}`

func TestInventoryPortsDecode(t *testing.T) {
	var r InventoryPortsResult
	if err := json.Unmarshal([]byte(inventoryPortsJSON), &r); err != nil {
		t.Fatal(err)
	}
	if !r.Success || r.HostCount != 1 || r.AllocatedPortCount != 1 || len(r.Hosts) != 1 {
		t.Fatalf("summary decoded wrong: %+v", r)
	}
	h, ok := r.Host("SA-ZT-NDR-02")
	if !ok || h.Hostname != "hgx-su00-h00" || h.IPAddress != "" || len(h.Ports) != 2 {
		t.Fatalf("host lookup: %+v %v", h, ok)
	}
	p0, p1 := h.Ports[0], h.Ports[1]
	if p0.TenantID == nil || *p0.TenantID != 3 || !p0.Allocated || p0.AllocatedAt.IsZero() || !p0.LastSyncedAt.IsZero() {
		t.Errorf("allocated port decoded wrong: %+v", p0)
	}
	if p1.TenantID != nil || p1.Allocated || p1.TenantName != "" {
		t.Errorf("free port decoded wrong: %+v", p1)
	}
	if string(r.Extra["newSummaryKey"]) != "7" || len(r.Extra) != 1 {
		t.Errorf("Extra = %v", r.Extra)
	}
	if _, ok := r.Host("nope"); ok {
		t.Error("unknown host found")
	}
}

func TestUFMHostsAndSyncDecode(t *testing.T) {
	var h UFMHostsResult
	if err := json.Unmarshal([]byte(`{"fabricName":"UFM1","timestamp":"2026-04-18T14:03:32Z","success":true,
		"hostCount":1,"source":"UFM_LIVE","hosts":[{"systemGuid":"0xabc","hostname":"SA-ZT-NDR-02",
		"ipAddress":"10.0.0.1","state":"active","vendor":"NVIDIA","role":"host","type":"host","portCount":8}]}`), &h); err != nil {
		t.Fatal(err)
	}
	if h.Source != "UFM_LIVE" || len(h.Hosts) != 1 || h.Hosts[0].PortCount != 8 || h.Extra != nil {
		t.Errorf("UFMHostsResult = %+v", h)
	}
	if _, ok := h.Timestamp.Time(); !ok {
		t.Error("timestamp did not parse")
	}

	var s InventorySyncResult
	if err := json.Unmarshal([]byte(`{"fabricName":"UFM1","success":true,"message":"Inventory sync completed successfully"}`), &s); err != nil {
		t.Fatal(err)
	}
	if !s.Success || s.Message == "" {
		t.Errorf("InventorySyncResult = %+v", s)
	}
}

func TestNmxcInventoryDecode(t *testing.T) {
	var n NmxcInventory
	if err := json.Unmarshal([]byte(`{
	  "fabricName":"NMXC","fabricType":"NMXC","success":true,"overallHealth":"HEALTHY",
	  "summary":{"domainCount":1,"totalGpuCount":4,"freeGpuCount":3,"allocatedGpuCount":1,"serverCount":1},
	  "domains":[{"domainId":12,"domainName":"su00-rack00","connectionState":"ACTIVE","nmxcHost":"10.0.0.5",
	    "nmxcPort":9370,"maxGpus":null,"maxPartitions":32,"topologyType":"GB200","lastSyncedAt":"2026-04-18T14:03:32",
	    "gpuCount":4,"freeGpuCount":3,"allocatedGpuCount":1,"serverCount":1,
	    "partitions":[{"partitionId":5,"partitionName":"t1-p","tenantId":3,"tenantName":"t1",
	      "resiliencyMode":"NMX_RESILIENCY_MODE_ADAPTIVE_BANDWIDTH","status":"ACTIVE","gpuCount":1}],
	    "servers":[{"serverName":"su00-rack00-node00","trayIndex":0,"totalGpus":4,"freeGpus":3,"allocatedGpus":1,
	      "tenantId":3,"tenantName":"t1","lastConfiguredTenant":"t1",
	      "gpus":[{"gpuUid":"u0","gpuId":0,"healthStatus":"HEALTHY","partitionId":5,"tenantId":3,"tenantName":"t1","systemUid":"s0"},
	              {"gpuUid":"u1","gpuId":1,"healthStatus":"HEALTHY","partitionId":null,"tenantId":null,"tenantName":null,"systemUid":"s0"}]}]}]}`), &n); err != nil {
		t.Fatal(err)
	}
	if n.Summary.TotalGpuCount != 4 || len(n.Domains) != 1 {
		t.Fatalf("NmxcInventory = %+v", n)
	}
	d := n.Domains[0]
	if d.DomainID.String() != "12" || d.MaxGpus != nil || *d.MaxPartitions != 32 || len(d.Servers[0].Gpus) != 2 {
		t.Errorf("domain decoded wrong: %+v", d)
	}
	if g := d.Servers[0].Gpus[1]; g.PartitionID != nil || g.TenantID != nil || *g.GpuID != 1 {
		t.Errorf("free GPU decoded wrong: %+v", g)
	}
}

func TestNmxcInventoryStringDomainID(t *testing.T) {
	var n NmxcInventory
	if err := json.Unmarshal([]byte(`{"success":true,"domains":[{"domainId":"12"}]}`), &n); err != nil {
		t.Fatalf("a string domainId should decode too: %v", err)
	}
	if n.Domains[0].DomainID.String() != "12" {
		t.Errorf("DomainID = %q", n.Domains[0].DomainID)
	}
}

// Invalid input must fail at once, not poll until the wait times out and
// surface as a WaitTimeoutError. The resources have no transport here, so any
// attempt to send a request would panic.
func TestWaitHelpersRejectInvalidInputImmediately(t *testing.T) {
	ctx := context.Background()
	opts := WaitOptions{Timeout: time.Minute, Interval: time.Second}
	started := time.Now()

	_, err := (&OperationsResource{}).WaitForOperation(ctx, "", opts)
	assertField(t, err, "operationId")
	_, err = (&OperationsResource{}).WaitForOperation(ctx, "op/1", opts)
	assertField(t, err, "operationId")

	_, err = (&TenantsResource{}).WaitForTenantConfigApplied(ctx, "f", "", opts)
	assertField(t, err, "tenantName")

	if waited := time.Since(started); waited > time.Second {
		t.Errorf("validation failures took %s; they should return without polling", waited)
	}
}
