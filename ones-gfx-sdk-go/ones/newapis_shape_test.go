package ones

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// seen records the last request the stub server received.
type seen struct {
	mu     sync.Mutex
	method string
	path   string
	accept string
	body   string
}

func (s *seen) set(r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.method = r.Method
	s.path = r.URL.RequestURI()
	s.accept = r.Header.Get("Accept")
	var buf [4096]byte
	n, _ := r.Body.Read(buf[:])
	s.body = string(buf[:n])
}

// newStub returns a client pointed at a server that records every request.
func newStub(t *testing.T) (*Client, *seen, func()) {
	t.Helper()
	s := &seen{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/login" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"token": "T"}})
			return
		}
		s.set(r)
		w.Header().Set("Content-Type", "application/json")
		// The collection endpoints decode into a slice, so the stub has to
		// answer with one; everything else takes an object.
		if r.Method == http.MethodGet && r.URL.Path == "/api/fm/webhooks" {
			json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"ok": true}})
	}))
	c, err := InitializeWithToken(srv.URL, "T")
	if err != nil {
		srv.Close()
		t.Fatalf("client: %v", err)
	}
	return c, s, srv.Close
}

func TestNewAPIRequestShapes(t *testing.T) {
	c, s, done := newStub(t)
	defer done()
	ctx := context.Background()
	const F, Tn, Sv = "fab1", "ten1", "srv1"

	cases := []struct {
		name       string
		call       func() error
		wantMethod string
		wantPath   string
	}{
		// A1 network policy
		{"ListNetworkPolicies", func() error { _, e := ListNetworkPolicies(ctx, c, F); return e },
			"GET", "/api/fm/fabrics/fab1/networkpolicy"},
		{"CreateNetworkPolicy", func() error {
			_, e := CreateNetworkPolicy(ctx, c, F, "p1",
				NetworkPolicySide{Type: NetworkPolicyTypeTenant, TenantName: "a"},
				NetworkPolicySide{Type: NetworkPolicyTypeService, TenantName: "b", ServiceName: "s"})
			return e
		}, "POST", "/api/fm/fabrics/fab1/networkpolicy"},
		{"PatchNetworkPolicy", func() error {
			_, e := PatchNetworkPolicy(ctx, c, F, "p1",
				NetworkPolicySide{Type: "tenant", TenantName: "a"},
				NetworkPolicySide{Type: "service", TenantName: "b", Endpoints: []string{"10.0.0.2"}})
			return e
		}, "PATCH", "/api/fm/fabrics/fab1/networkpolicy/p1"},
		{"DeleteNetworkPolicy", func() error { _, e := DeleteNetworkPolicy(ctx, c, F, "p1"); return e },
			"DELETE", "/api/fm/fabrics/fab1/networkpolicy/p1"},

		// A2 webhooks — global
		{"RegisterWebhook", func() error {
			_, e := RegisterWebhook(ctx, c, "https://x/cb", []string{WebhookEventTenantCreate}, nil, nil)
			return e
		}, "POST", "/api/fm/webhooks"},
		{"ListWebhooks", func() error { _, e := ListWebhooks(ctx, c); return e }, "GET", "/api/fm/webhooks"},
		{"GetWebhook", func() error { _, e := GetWebhook(ctx, c, "w1"); return e }, "GET", "/api/fm/webhooks/w1"},
		{"UpdateWebhook", func() error { _, e := UpdateWebhook(ctx, c, "w1", nil, nil, nil, nil); return e },
			"PATCH", "/api/fm/webhooks/w1"},
		{"DeleteWebhook", func() error { return DeleteWebhook(ctx, c, "w1") }, "DELETE", "/api/fm/webhooks/w1"},

		// A2 webhooks — fabric
		{"RegisterFabricWebhook", func() error { _, e := RegisterFabricWebhook(ctx, c, F, "https://x/cb", nil, nil); return e },
			"POST", "/api/fm/fabrics/fab1/webhook"},
		{"GetFabricWebhook", func() error { _, e := GetFabricWebhook(ctx, c, F); return e },
			"GET", "/api/fm/fabrics/fab1/webhook"},
		{"UpdateFabricWebhook", func() error { _, e := UpdateFabricWebhook(ctx, c, F, nil, true, nil); return e },
			"PUT", "/api/fm/fabrics/fab1/webhook"},
		{"DeleteFabricWebhook", func() error { return DeleteFabricWebhook(ctx, c, F) },
			"DELETE", "/api/fm/fabrics/fab1/webhook"},

		// A2 webhooks — tenant
		{"RegisterTenantWebhook", func() error { _, e := RegisterTenantWebhook(ctx, c, F, Tn, "https://x", nil, nil); return e },
			"POST", "/api/fm/fabrics/fab1/tenants/ten1/webhook"},
		{"GetTenantWebhook", func() error { _, e := GetTenantWebhook(ctx, c, F, Tn); return e },
			"GET", "/api/fm/fabrics/fab1/tenants/ten1/webhook"},
		{"UpdateTenantWebhook", func() error { _, e := UpdateTenantWebhook(ctx, c, F, Tn, nil, true, nil); return e },
			"PUT", "/api/fm/fabrics/fab1/tenants/ten1/webhook"},
		{"DeleteTenantWebhook", func() error { return DeleteTenantWebhook(ctx, c, F, Tn) },
			"DELETE", "/api/fm/fabrics/fab1/tenants/ten1/webhook"},

		// A3 IP queries
		{"GetServerIps", func() error { _, e := GetServerIps(ctx, c, F, Sv); return e },
			"GET", "/api/fm/fabrics/fab1/computes/srv1/ip-addresses"},
		{"GetTenantIps", func() error { _, e := GetTenantIps(ctx, c, F, Tn, 1, 10, IPAssignStatusAssigned); return e },
			"GET", "/api/fm/fabrics/fab1/tenants/ten1/ip-addresses?page=1&size=10&status=ASSIGNED"},
		{"GetTenantComputeIp", func() error { _, e := GetTenantComputeIp(ctx, c, F, Tn, Sv); return e },
			"GET", "/api/fm/fabrics/fab1/tenants/ten1/computes/srv1/ip-address"},

		// A4 DHCP
		{"ReconcileFabricDhcp", func() error { _, e := ReconcileFabricDhcp(ctx, c, F); return e },
			"POST", "/api/fm/fabrics/fab1/dhcp/reconcile"},
		{"GetTenantDhcp", func() error { _, e := GetTenantDhcp(ctx, c, F, Tn, 0, 100); return e },
			"GET", "/api/fm/fabrics/fab1/tenants/ten1/dhcp?page=0&size=100"},
		{"RetryTenantDhcp", func() error { _, e := RetryTenantDhcp(ctx, c, F, Tn); return e },
			"POST", "/api/fm/fabrics/fab1/tenants/ten1/dhcp/retry"},

		// A5 IPAM config
		{"GetIpamConfig", func() error { _, e := GetIpamConfig(ctx, c, F); return e },
			"GET", "/api/fm/fabrics/fab1/ipam/config"},
		{"UpdateIpamConfig", func() error {
			_, e := UpdateIpamConfig(ctx, c, F, IpamConfig{EnableIpam: true,
				Dhcp: &IpamDhcpConfig{Enabled: true}}, true)
			return e
		}, "PATCH", "/api/fm/fabrics/fab1/ipam/config?force=true"},

		// A6 rails
		{"GetRailCidrPools", func() error { _, e := GetRailCidrPools(ctx, c, F); return e },
			"GET", "/api/fm/fabrics/fab1/rails/cidrpools"},
		{"GetRailCidrPool", func() error { _, e := GetRailCidrPool(ctx, c, F, "rail0"); return e },
			"GET", "/api/fm/fabrics/fab1/rails/rail0/cidrpool"},

		// A7 inventory
		{"DeleteInventoryByHostnames", func() error {
			_, e := DeleteInventoryByHostnames(ctx, c, F, []string{"h1", "h2"})
			return e
		}, "DELETE", "/api/fm/fabrics/fab1/inventory"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err != nil {
				t.Fatalf("call failed: %v", err)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			if s.method != tc.wantMethod || s.path != tc.wantPath {
				t.Errorf("got %s %s, want %s %s", s.method, s.path, tc.wantMethod, tc.wantPath)
			}
		})
	}
}

// TestRailYAMLAcceptHeader covers the content-negotiated rail manifests,
// the only endpoints that ask for a non-JSON Accept.
func TestRailYAMLAcceptHeader(t *testing.T) {
	c, s, done := newStub(t)
	defer done()
	if _, err := GetRailCidrPoolYAML(context.Background(), c, "fab1", "rail0", ""); err != nil {
		t.Fatalf("call failed: %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.accept != RailYamlMediaType {
		t.Errorf("Accept = %q, want %q", s.accept, RailYamlMediaType)
	}
}

// TestIpamConfigDropsFabricName guards the strip in IpamResource.UpdateConfig:
// FM rejects a body whose fabricName disagrees with the URL.
func TestIpamConfigDropsFabricName(t *testing.T) {
	c, s, done := newStub(t)
	defer done()
	cfg := IpamConfig{FabricName: "a-different-fabric", EnableIpam: false}
	if _, err := UpdateIpamConfig(context.Background(), c, "fab1", cfg, false); err != nil {
		t.Fatalf("call failed: %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var got map[string]any
	if err := json.Unmarshal([]byte(s.body), &got); err != nil {
		t.Fatalf("body %q: %v", s.body, err)
	}
	if _, present := got["fabricName"]; present {
		t.Errorf("fabricName was sent: %s", s.body)
	}
}
