// Command e2e exercises the 28 Fabric Manager endpoints added to this SDK
// (groups A1-A7) against a live Fabric Manager, mirroring the Python suite.
//
// Safety: it never deletes pre-existing tenants, fabrics or inventory. Writes
// are confined to objects it creates and then removes (network policies,
// webhook registrations). Destructive-by-nature calls run on inputs that
// cannot match real state (A7 uses fabricated hostnames) or are asserted
// against their documented refusal (A5 on an ONES-controlled fabric).
//
// Usage: go run ./examples/e2e -url http://127.0.0.1:PORT -user U -pass P
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

const cb = "https://example.invalid/ones-sdk-e2e-go"

var (
	pass, fail int
	failures   []string
)

// status renders an error as the HTTP status the server returned, so an
// expected refusal reads the same way as in the Python suite.
func status(err error) int {
	var ae *ones_gfx.APIError
	if errors.As(err, &ae) {
		return ae.StatusCode
	}
	var be *ones_gfx.BadRequestError
	if errors.As(err, &be) {
		return be.StatusCode
	}
	var ne *ones_gfx.NotFoundError
	if errors.As(err, &ne) {
		return ne.StatusCode
	}
	var ce *ones_gfx.ConflictError
	if errors.As(err, &ce) {
		return ce.StatusCode
	}
	var se *ones_gfx.ServerError
	if errors.As(err, &se) {
		// ServerError embeds APIError, so the real 5xx is on the embedded
		// struct: a 502 from the DHCP provider must not be flattened to 500.
		return se.StatusCode
	}
	return 0
}

func short(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > n {
		return s[:n]
	}
	return s
}

// check records one call. allow lists HTTP codes that are a documented,
// correct outcome for this environment.
func check(name string, err error, note string, allow ...int) {
	if err == nil {
		pass++
		fmt.Printf("  [PASS] %-34s %s\n", name, note)
		return
	}
	code := status(err)
	for _, a := range allow {
		if code == a {
			pass++
			fmt.Printf("  [PASS] %-34s HTTP %d (expected): %s\n", name, code, short(err.Error(), 90))
			return
		}
	}
	fail++
	msg := fmt.Sprintf("%s: HTTP %d %s", name, code, short(err.Error(), 130))
	failures = append(failures, msg)
	fmt.Printf("  [FAIL] %-34s HTTP %d %s\n", name, code, short(err.Error(), 110))
}

func main() {
	url := flag.String("url", "http://127.0.0.1:8080", "ONES base URL")
	user := flag.String("user", "superadmin", "username")
	pw := flag.String("pass", "", "password")
	flag.Parse()

	c, err := ones.InitializeWithCreds(*url, *user, *pw)
	if err != nil {
		fmt.Println("login failed:", err)
		os.Exit(2)
	}
	ctx := context.Background()

	ver, err := ones.GetControllerVersion(ctx, c)
	if err != nil {
		fmt.Println("version failed:", err)
		os.Exit(2)
	}
	fmt.Printf("\nFM %v\n", ver)

	fabs, err := ones.GetAllFabrics(ctx, c)
	if err != nil {
		fmt.Println("fabrics failed:", err)
		os.Exit(2)
	}
	var F string
	for _, f := range fabs {
		if f.Status == "Deployed" {
			F = f.Name
			break
		}
	}
	if F == "" && len(fabs) > 0 {
		F = fabs[0].Name
	}
	T, S := "pep2", "hgx-su00-h00"
	fmt.Printf("target fabric=%s tenant=%s server=%s\n\n", F, T, S)

	const policy = "sdke2e-go-1"

	// ---------------- A1 network policy ----------------
	fmt.Println("A1  network policy")
	npl, err := ones.ListNetworkPolicies(ctx, c, F)
	if err == nil {
		check("ListNetworkPolicies", nil, fmt.Sprintf("count=%d status=%s", npl.Count, npl.Status))
	} else {
		check("ListNetworkPolicies", err, "")
	}
	src := ones.NetworkPolicySide{Type: ones.NetworkPolicyTypeTenant, TenantName: T}
	dst := ones.NetworkPolicySide{Type: ones.NetworkPolicyTypeService, TenantName: T, ServiceName: "svc", Endpoints: []string{"10.0.0.9"}}
	_, err = ones.CreateNetworkPolicy(ctx, c, F, policy, src, dst)
	check("CreateNetworkPolicy", err, "created", 400, 403, 409, 500)
	_, err = ones.PatchNetworkPolicy(ctx, c, F, policy, src, dst)
	if err != nil && strings.Contains(err.Error(), "policy.name is required") {
		fail++
		failures = append(failures, "PatchNetworkPolicy: body still missing policy.name")
		fmt.Printf("  [FAIL] %-34s body still missing policy.name\n", "PatchNetworkPolicy")
	} else {
		check("PatchNetworkPolicy", err, "patched", 400, 403, 404, 409, 500)
	}
	_, err = ones.DeleteNetworkPolicy(ctx, c, F, policy)
	check("DeleteNetworkPolicy", err, "deleted", 400, 403, 404, 500)

	// ---------------- A2 webhooks ----------------
	fmt.Println("A2  webhook registration")
	whs, err := ones.ListWebhooks(ctx, c)
	check("ListWebhooks", err, fmt.Sprintf("%d global registrations", len(whs)))

	reg, err := ones.RegisterWebhook(ctx, c, cb+"/global",
		[]string{ones.WebhookEventTenantCreate, ones.WebhookEventFabricDeploy}, nil,
		&ones.WebhookFilters{FabricNames: []string{F}})
	var gid string
	if err == nil {
		gid = reg.WebhookID
		check("RegisterWebhook", nil, fmt.Sprintf("webhookId=%s events=%v", gid, reg.EventTypes))
	} else {
		check("RegisterWebhook", err, "")
	}
	aliasReg, err := ones.RegisterWebhook(ctx, c, cb+"/alias", []string{"fabric.tenant_create"}, nil, nil)
	if err == nil {
		ok := len(aliasReg.EventTypes) == 1 && aliasReg.EventTypes[0] == "tenant.create"
		if ok {
			check("RegisterWebhook (alias)", nil, fmt.Sprintf("alias resolved -> %v", aliasReg.EventTypes))
		} else {
			check("RegisterWebhook (alias)", fmt.Errorf("alias not resolved: %v", aliasReg.EventTypes), "")
		}
	} else {
		check("RegisterWebhook (alias)", err, "")
	}
	_, err = ones.RegisterWebhook(ctx, c, cb+"/bad", []string{}, nil, nil)
	check("RegisterWebhook (no events)", err, "", 400)

	if gid != "" {
		got, err := ones.GetWebhook(ctx, c, gid)
		if err == nil {
			check("GetWebhook", nil, "round-tripped "+got.WebhookID)
		} else {
			check("GetWebhook", err, "")
		}
		no := false
		upd, err := ones.UpdateWebhook(ctx, c, gid, nil, &no, nil, nil)
		if err == nil {
			check("UpdateWebhook", nil, fmt.Sprintf("enabled -> %v", upd.Enabled))
		} else {
			check("UpdateWebhook", err, "")
		}
		check("DeleteWebhook", ones.DeleteWebhook(ctx, c, gid), "204 no content")
		_, err = ones.GetWebhook(ctx, c, gid)
		check("GetWebhook (after delete)", err, "", 404)
	}

	fr, err := ones.RegisterFabricWebhook(ctx, c, F, cb+"/fabric", nil, nil)
	if err == nil {
		check("RegisterFabricWebhook", nil, "fabricName="+fr.FabricName)
	} else {
		check("RegisterFabricWebhook", err, "")
	}
	fg, err := ones.GetFabricWebhook(ctx, c, F)
	if err == nil {
		check("GetFabricWebhook", nil, "url="+fg.CallbackURL)
	} else {
		check("GetFabricWebhook", err, "")
	}
	fu, err := ones.UpdateFabricWebhook(ctx, c, F, nil, false, []string{ones.WebhookEventTenantDelete})
	if err == nil {
		check("UpdateFabricWebhook", nil, fmt.Sprintf("enabled=%v events=%v", fu.Enabled, fu.EventTypes))
	} else {
		check("UpdateFabricWebhook", err, "")
	}
	check("DeleteFabricWebhook", ones.DeleteFabricWebhook(ctx, c, F), "204 no content")

	tr, err := ones.RegisterTenantWebhook(ctx, c, F, T, cb+"/tenant", nil, nil)
	if err == nil {
		check("RegisterTenantWebhook", nil, "tenantName="+tr.TenantName)
	} else {
		check("RegisterTenantWebhook", err, "")
	}
	tg, err := ones.GetTenantWebhook(ctx, c, F, T)
	if err == nil {
		check("GetTenantWebhook", nil, "url="+tg.CallbackURL)
	} else {
		check("GetTenantWebhook", err, "")
	}
	tu, err := ones.UpdateTenantWebhook(ctx, c, F, T, nil, true, []string{})
	if err == nil {
		check("UpdateTenantWebhook", nil, fmt.Sprintf("enabled=%v events=%v", tu.Enabled, tu.EventTypes))
	} else {
		check("UpdateTenantWebhook", err, "")
	}
	check("DeleteTenantWebhook", ones.DeleteTenantWebhook(ctx, c, F, T), "204 no content")

	// ---------------- A3 IP queries ----------------
	fmt.Println("A3  IPAM address queries")
	tips, err := ones.GetTenantIps(ctx, c, F, T, 0, 10, "")
	if err == nil {
		check("GetTenantIps", nil, fmt.Sprintf("pool.cidr=%s allocs=%d", tips.Pool.Cidr, len(tips.Allocations)))
	} else {
		check("GetTenantIps", err, "", 404)
	}
	_, err = ones.GetTenantIps(ctx, c, F, T, 0, 10, ones.IPAssignStatusAssigned)
	check("GetTenantIps (status filter)", err, "filtered", 404)
	_, err = ones.GetTenantComputeIp(ctx, c, F, T, S)
	check("GetTenantComputeIp", err, "returned", 404)
	sips, err := ones.GetServerIps(ctx, c, F, S)
	if err == nil {
		check("GetServerIps", nil, "serverName="+sips.ServerName)
	} else {
		check("GetServerIps", err, "", 404)
	}
	_, err = ones.GetServerIps(ctx, c, F, "no-such-server-e2e")
	check("GetServerIps (unknown)", err, "", 404)

	// ---------------- A4 DHCP ----------------
	fmt.Println("A4  tenant DHCP")
	td, err := ones.GetTenantDhcp(ctx, c, F, T, 0, 20)
	if err == nil {
		check("GetTenantDhcp", nil, fmt.Sprintf("leases=%d", len(td.Leases)))
	} else {
		check("GetTenantDhcp", err, "", 404)
	}
	_, err = ones.RetryTenantDhcp(ctx, c, F, T)
	check("RetryTenantDhcp", err, "re-applied", 400, 404, 502)
	_, err = ones.ReconcileFabricDhcp(ctx, c, F)
	check("ReconcileFabricDhcp", err, "reconciled", 400, 404, 502)
	_, err = ones.ReconcileFabricDhcp(ctx, c, "no-such-fabric-e2e")
	check("ReconcileFabricDhcp (unknown)", err, "", 404)

	// ---------------- A5 IPAM config ----------------
	fmt.Println("A5  IPAM configuration")
	cfg, err := ones.GetIpamConfig(ctx, c, F)
	if err == nil {
		check("GetIpamConfig", nil, fmt.Sprintf("enableIpam=%v static=%v dhcp=%v",
			cfg.EnableIpam, cfg.Static != nil, cfg.Dhcp != nil))
	} else {
		check("GetIpamConfig", err, "")
	}
	lease := 3600
	_, err = ones.UpdateIpamConfig(ctx, c, F,
		ones.IpamConfig{EnableIpam: true, Dhcp: &ones.IpamDhcpConfig{Enabled: true, LeaseTimeSeconds: &lease}}, false)
	check("UpdateIpamConfig (ONES-ctl)", err, "", 409, 400)
	_, err = ones.GetIpamConfig(ctx, c, "no-such-fabric-e2e")
	check("GetIpamConfig (unknown)", err, "", 404)

	// ---------------- A6 rails ----------------
	fmt.Println("A6  rail CIDRPools")
	_, err = ones.GetRailCidrPools(ctx, c, F)
	check("GetRailCidrPools", err, "returned", 404)
	_, err = ones.GetRailCidrPoolsYAML(ctx, c, F, ones.RailYamlMediaType)
	check("GetRailCidrPoolsYAML", err, "returned yaml", 404)
	_, err = ones.GetRailCidrPool(ctx, c, F, "no-such-rail-e2e")
	check("GetRailCidrPool (unknown)", err, "", 404)

	// ---------------- A7 inventory ----------------
	fmt.Println("A7  inventory delete by hostname")
	res, err := ones.DeleteInventoryByHostnames(ctx, c, F, []string{"e2e-ghost-1", "e2e-ghost-2"})
	if err == nil {
		check("DeleteInventoryByHostnames", nil, fmt.Sprintf("notFound=%v", res.NotFound))
	} else {
		check("DeleteInventoryByHostnames", err, "", 207, 404)
	}
	_, err = ones.DeleteInventoryByHostnames(ctx, c, F, []string{})
	check("DeleteInventory (empty)", err, "", 400)

	// cleanup: remove any webhook this run left behind
	if left, e := ones.ListWebhooks(ctx, c); e == nil {
		n := 0
		for _, w := range left {
			if strings.Contains(w.CallbackURL, "ones-sdk-e2e-go") {
				if ones.DeleteWebhook(ctx, c, w.WebhookID) == nil {
					n++
				}
			}
		}
		fmt.Printf("\ncleanup: removed %d leftover webhook(s)\n", n)
	}

	fmt.Printf("\n%s\nGo E2E: %d passed, %d failed (%d checks)\n%s\n",
		strings.Repeat("=", 70), pass, fail, pass+fail, strings.Repeat("=", 70))
	if fail > 0 {
		fmt.Println("FAILURES:")
		for _, f := range failures {
			fmt.Println("  -", f)
		}
		os.Exit(1)
	}
}
