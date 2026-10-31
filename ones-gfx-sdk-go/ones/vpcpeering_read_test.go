package ones

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// Bodies as FM's ApiResponse.success(message, data) builds them (GFX-5658).
const (
	peeringListBody = `{"status":"success","message":"VPC peerings retrieved successfully.","data":{
		"fabricName":"Fab101",
		"peerings":[
		  {"name":"t1-storage","vpcs":["t1-Fab101-north-south","Fab101-Storage"],"origin":"VPC_PEERING","createdAt":"2026-10-28T10:31:29.123Z"},
		  {"name":"pol-t1-t2","vpcs":["t1-Fab101-north-south","t2-Fab101-north-south"],"origin":"NETWORK_POLICY","createdAt":"2026-10-28T11:00:00+05:30"}],
		"pagination":{"page":1,"size":2,"totalElements":5,"totalPages":3}}}`

	peeringTenantBody = `{"status":"success","message":"VPC peerings retrieved successfully.","data":{
		"fabricName":"Fab101","tenantName":"t1","vpc":"t1-Fab101-north-south",
		"peerings":[{"name":"t1-storage","vpc":"t1-Fab101-north-south","peerVpc":"Fab101-Storage","origin":"VPC_PEERING","createdAt":"2026-10-28T10:31:29Z"}],
		"pagination":{"page":0,"size":50,"totalElements":1,"totalPages":1}}}`

	peeringGetBody = `{"status":"success","message":"VPC peering retrieved successfully.","data":{
		"fabricName":"Fab101",
		"peering":{"name":"t1-storage","vpcs":["t1-Fab101-north-south","Fab101-Storage"],"origin":"VPC_PEERING","createdAt":"2026-10-28T10:31:29Z"}}}`
)

func TestListVpcPeerings(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		// A GET-only direct route: served at FM's root, no /api/fm prefix.
		"GET /fabrics/Fab101/vpcpeering": ok(peeringListBody),
	})
	ctx := context.Background()

	list, err := ListVpcPeerings(ctx, c, "Fab101", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if list.FabricName != "Fab101" || len(list.Peerings) != 2 {
		t.Fatalf("list = %+v", list)
	}
	p0, p1 := list.Peerings[0], list.Peerings[1]
	if p0.Name != "t1-storage" || len(p0.Vpcs) != 2 || p0.Vpcs[1] != "Fab101-Storage" || p0.Vpc != "" {
		t.Errorf("fabric-wide shape decoded wrong: %+v", p0)
	}
	if p0.Origin != VPCPeeringOriginVPCPeering || p1.Origin != VPCPeeringOriginNetworkPolicy {
		t.Errorf("origins = %q, %q", p0.Origin, p1.Origin)
	}
	if _, ok := p1.CreatedAt.Time(); !ok {
		t.Errorf("createdAt with offset did not parse: %q", p1.CreatedAt)
	}
	if pg := list.Pagination; pg.Page != 1 || pg.Size != 2 || pg.TotalElements != 5 || pg.TotalPages != 3 {
		t.Errorf("pagination = %+v", pg)
	}

	if r := stub.last(); r.Method != http.MethodGet || r.Path != "/fabrics/Fab101/vpcpeering" || r.Body != nil {
		t.Errorf("request = %+v", r)
	}
}

func TestListVpcPeeringsQuery(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		"GET /fabrics/F/vpcpeering":            ok(peeringListBody),
		"GET /fabrics/F/tenants/t1/vpcpeering": ok(peeringTenantBody),
	})
	ctx := context.Background()
	cases := []struct {
		page, size int
		want       string
	}{
		{2, 25, "page=2&size=25"},
		{-1, 0, ""},                 // both omitted: FM defaults (page 0, size 50)
		{0, 600, "page=0&size=600"}, // passed through; FM caps size at 500
	}
	for _, tc := range cases {
		if _, err := ListVpcPeerings(ctx, c, "F", tc.page, tc.size); err != nil {
			t.Fatal(err)
		}
		if got := stub.last().Query; got != tc.want {
			t.Errorf("List(page=%d, size=%d) query = %q, want %q", tc.page, tc.size, got, tc.want)
		}
		if _, err := ListTenantVpcPeerings(ctx, c, "F", "t1", tc.page, tc.size); err != nil {
			t.Fatal(err)
		}
		if got := stub.last().Query; got != tc.want {
			t.Errorf("ListForTenant(page=%d, size=%d) query = %q, want %q", tc.page, tc.size, got, tc.want)
		}
	}
}

func TestListTenantVpcPeerings(t *testing.T) {
	stub, c := newFMStub(t, map[string]stubResponse{
		"GET /fabrics/Fab101/tenants/t1/vpcpeering": ok(peeringTenantBody),
		"GET /fabrics/Fab101/tenants/ghost/vpcpeering": {http.StatusNotFound,
			`{"status":"error","message":"Tenant 'ghost' does not exist in fabric 'Fab101'.","error":{"code":"TENANT_NOT_FOUND","detail":"x"}}`},
	})
	ctx := context.Background()

	list, err := ListTenantVpcPeerings(ctx, c, "Fab101", "t1", -1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if list.TenantName != "t1" || list.Vpc != "t1-Fab101-north-south" || len(list.Peerings) != 1 {
		t.Fatalf("list = %+v", list)
	}
	if p := list.Peerings[0]; p.Vpc != "t1-Fab101-north-south" || p.PeerVpc != "Fab101-Storage" || p.Vpcs != nil {
		t.Errorf("tenant-scoped shape decoded wrong: %+v", p)
	}
	if r := stub.last(); r.Path != "/fabrics/Fab101/tenants/t1/vpcpeering" {
		t.Errorf("path = %s", r.Path)
	}

	_, err = ListTenantVpcPeerings(ctx, c, "Fab101", "ghost", -1, 0)
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected *NotFoundError, got %T (%v)", err, err)
	}
	body, _ := nf.ResponseBody.(map[string]interface{})
	if detail, _ := body["error"].(map[string]interface{}); detail["code"] != "TENANT_NOT_FOUND" {
		t.Errorf("error.code = %v", body["error"])
	}
}

func TestGetVpcPeering(t *testing.T) {
	_, c := newFMStub(t, map[string]stubResponse{
		"GET /fabrics/Fab101/vpcpeering/t1-storage": ok(peeringGetBody),
		"GET /fabrics/Fab101/vpcpeering/nope": {http.StatusNotFound,
			`{"status":"error","message":"VPC peering 'nope' does not exist in fabric 'Fab101'.","error":{"code":"VPC_PEERING_NOT_FOUND","detail":"x"}}`},
		"GET /fabrics/Gone/vpcpeering/x": {http.StatusNotFound,
			`{"status":"error","message":"Fabric 'Gone' not found.","error":{"code":"FABRIC_NOT_FOUND","detail":"x"}}`},
	})
	ctx := context.Background()

	p, err := GetVpcPeering(ctx, c, "Fab101", "t1-storage")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "t1-storage" || len(p.Vpcs) != 2 || p.Origin != VPCPeeringOriginVPCPeering {
		t.Errorf("peering = %+v", p)
	}

	for _, tc := range []struct{ fabric, name string }{{"Fab101", "nope"}, {"Gone", "x"}} {
		_, err := GetVpcPeering(ctx, c, tc.fabric, tc.name)
		var nf *NotFoundError
		if !errors.As(err, &nf) {
			t.Errorf("%s/%s: expected *NotFoundError, got %T (%v)", tc.fabric, tc.name, err, err)
		}
	}
}

func TestVpcPeeringReadsValidateBeforeSending(t *testing.T) {
	stub, c := newFMStub(t, nil)
	ctx := context.Background()
	cases := []struct {
		field string
		call  func() error
	}{
		{"fabricName", func() error { _, err := ListVpcPeerings(ctx, c, "", 0, 0); return err }},
		{"peeringName", func() error { _, err := GetVpcPeering(ctx, c, "F", ""); return err }},
		{"peeringName", func() error { _, err := GetVpcPeering(ctx, c, "F", "a/b"); return err }},
		{"tenantName", func() error { _, err := ListTenantVpcPeerings(ctx, c, "F", " ", 0, 0); return err }},
	}
	for _, tc := range cases {
		var ve *ValidationError
		if err := tc.call(); !errors.As(err, &ve) || ve.Field != tc.field {
			t.Errorf("want ValidationError on %s, got %T (%v)", tc.field, err, err)
		}
	}
	if n := len(stub.all()); n != 0 {
		t.Fatalf("%d request(s) reached the server", n)
	}
}

// The reads bypass the gateway; create and delete on the same path must not.
func TestVpcPeeringRouting(t *testing.T) {
	cases := []struct {
		method, path string
		direct       bool
	}{
		{"GET", "fabrics/F/vpcpeering", true},
		{"GET", "fabrics/F/vpcpeering?page=0&size=50", true},
		{"GET", "fabrics/F/vpcpeering/p1", true},
		{"GET", "fabrics/F/tenants/t1/vpcpeering", true},
		{"POST", "fabrics/F/vpcpeering", false},
		{"DELETE", "fabrics/F/vpcpeering", false},
		{"GET", "fabrics/F/tenants/t1", false},
		{"POST", "fabrics/F/networkpolicy", true}, // path-only table still applies to every method
	}
	for _, tc := range cases {
		if got := ones_gfx.IsDirectRoute(tc.method, tc.path); got != tc.direct {
			t.Errorf("IsDirectRoute(%s, %s) = %v, want %v", tc.method, tc.path, got, tc.direct)
		}
	}

	// End to end: DELETE still goes through the gateway prefix.
	stub, c := newFMStub(t, map[string]stubResponse{
		"DELETE /api/fm/fabrics/F/vpcpeering": ok(`{"status":"success","message":"deleted"}`),
	})
	if _, err := DeleteVpcPeering(context.Background(), c, "F", "p1", "t1-F-north-south", "F-Storage"); err != nil {
		t.Fatal(err)
	}
	if r := stub.last(); r.Path != "/api/fm/fabrics/F/vpcpeering" {
		t.Errorf("DELETE went to %s, want the gateway path", r.Path)
	}
}
