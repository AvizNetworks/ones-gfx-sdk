package ones_gfx

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadTopology(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/ns_topology.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func extract(t *testing.T, opts NSInterfaceOptions) *NSInterfaces {
	t.Helper()
	n, err := ExtractNSInterfaces(loadTopology(t), opts)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func ifaces(n *NSInterfaces, host string) []NSInterface {
	for _, h := range n.Hosts {
		if h.Host == host {
			return h.NSInterfaces
		}
	}
	return nil
}

func hostNames(n *NSInterfaces) []string {
	var out []string
	for _, h := range n.Hosts {
		out = append(out, h.Host)
	}
	return out
}

func direct(iface, leaf, port string) NSInterface {
	return NSInterface{Interface: iface, Leaf: leaf, LeafPort: port}
}

func viaDPU(iface, dpu, dpuIface, leaf, port string) NSInterface {
	return NSInterface{Interface: iface, Via: dpu, ViaInterface: dpuIface, Leaf: leaf, LeafPort: port}
}

// The Python SDK checks against the same golden file, so the two stay identical.
func TestExtractNSInterfacesMatchesGolden(t *testing.T) {
	want, err := os.ReadFile("testdata/ns_topology_expected.yaml")
	if err != nil {
		t.Fatal(err)
	}
	got, err := extract(t, NSInterfaceOptions{FabricName: "EW"}).YAML()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("YAML differs from golden\n--- got\n%s\n--- want\n%s", got, want)
	}
}

func TestExtractNSInterfacesBothSidesAndDedup(t *testing.T) {
	n := extract(t, NSInterfaceOptions{})
	// h00: each cable is listed under the host and the leaf; reported once.
	if got, want := ifaces(n, "hgx-su00-h00"), []NSInterface{
		direct("enp8s0f1np0", "leaf-cn-00", "swp9"), direct("enp9s0f1np0", "leaf-cn-01", "swp9"),
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("h00 = %v, want %v", got, want)
	}
	// h01: one link has no destination, the other is only on the leaf side.
	if got, want := ifaces(n, "hgx-su00-h01"), []NSInterface{
		direct("enp8s0f1np0", "leaf-cn-00", "swp10"), direct("enp9s0f1np0", "leaf-cn-01", "swp10"),
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("h01 = %v, want %v", got, want)
	}
	if n.HostCount != 4 || n.InterfaceCount != 10 || n.Fabric != nil {
		t.Errorf("counts/fabric = %d %d %v", n.HostCount, n.InterfaceCount, n.Fabric)
	}
}

func TestExtractNSInterfacesDefaultCoversCnFnSn(t *testing.T) {
	n := extract(t, NSInterfaceOptions{})
	if !reflect.DeepEqual(n.LeafPatterns, []string{"^leaf-cn", "^leaf-fn", "^leaf-sn"}) {
		t.Errorf("default patterns = %v", n.LeafPatterns)
	}
	if got, want := ifaces(n, "hgx-su00-h02"), []NSInterface{
		direct("enp12s0f1np0", "leaf-sn-00", "swp5"), direct("enp8s0f1np0", "leaf-fn-00", "swp5"),
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("h02 (front-end + storage) = %v, want %v", got, want)
	}
}

func TestExtractNSInterfacesFollowsDPU(t *testing.T) {
	const dpu = "hgx-su00-h03-dpu-00"
	want := []NSInterface{
		// No PF index in the DPU interface: every N-S uplink.
		viaDPU("hostif", dpu, "mgmt_if", "leaf-cn-00", "swp20"),
		viaDPU("hostif", dpu, "mgmt_if", "leaf-cn-01", "swp20"),
		// pf0 -> uplink p0 -> leaf-cn-00; pf1 -> p1 -> leaf-cn-01.
		viaDPU("vf0", dpu, "pf0vf0_if", "leaf-cn-00", "swp20"),
		viaDPU("vf4", dpu, "pf1vf0_if", "leaf-cn-01", "swp20"),
	}
	if got := ifaces(extract(t, NSInterfaceOptions{}), "hgx-su00-h03"); !reflect.DeepEqual(got, want) {
		t.Errorf("h03 = %v\nwant %v", got, want)
	}
}

func TestExtractNSInterfacesIgnoresNonNS(t *testing.T) {
	for _, h := range extract(t, NSInterfaceOptions{}).Hosts {
		for _, i := range h.NSInterfaces {
			if i.Leaf == "leaf-su00-r0" || strings.HasPrefix(i.Leaf, "leaf-border") {
				t.Errorf("non-N-S leaf reported: %+v", i)
			}
			// hgx-su00-h00's DPU has no uplinks, so its vf0 has no N-S path.
			if i.Via == "hgx-su00-h00-dpu-00" {
				t.Errorf("path through a DPU with no uplinks: %+v", i)
			}
		}
	}
}

func TestExtractNSInterfacesPatterns(t *testing.T) {
	if got := hostNames(extract(t, NSInterfaceOptions{LeafPatterns: []string{"^leaf-fn"}})); !reflect.DeepEqual(got, []string{"hgx-su00-h02"}) {
		t.Errorf("^leaf-fn hosts = %v", got)
	}
	if got := hostNames(extract(t, NSInterfaceOptions{LeafPatterns: []string{"^leaf-cn"}})); !reflect.DeepEqual(got, []string{"hgx-su00-h00", "hgx-su00-h01", "hgx-su00-h03"}) {
		t.Errorf("^leaf-cn hosts = %v", got)
	}
	none := extract(t, NSInterfaceOptions{LeafPatterns: []string{"^nope"}})
	if none.HostCount != 0 || none.Hosts == nil {
		t.Errorf("no-match result = %+v", none)
	}
	if y, _ := none.YAML(); !strings.Contains(string(y), "hosts: []") {
		t.Errorf("empty result YAML = %s", y)
	}
}

func TestExtractNSInterfacesNoCredentials(t *testing.T) {
	y, _ := extract(t, NSInterfaceOptions{FabricName: "EW"}).YAML()
	for _, s := range []string{"secret", "password", "Credentials", "admin", "root"} {
		if strings.Contains(string(y), s) {
			t.Errorf("output contains %q", s)
		}
	}
}

func TestExtractNSInterfacesStringIDs(t *testing.T) {
	doc := `Connectivity:
  Leaf:
    - switchId: "7"
      switchName: leaf-cn-00
  Host:
    - hostId: "3"
      hostName: h3
      Links:
        - link: "G3_eth0 | L7_Ethernet8"
`
	n, err := ExtractNSInterfaces([]byte(doc), NSInterfaceOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := ifaces(n, "h3"); !reflect.DeepEqual(got, []NSInterface{direct("eth0", "leaf-cn-00", "Ethernet8")}) {
		t.Errorf("got %+v", got)
	}
}

func TestExtractNSInterfacesErrors(t *testing.T) {
	for _, in := range []string{"not: [valid", "<!doctype html><html></html>", "Inventory: {}\n", "- a list\n"} {
		_, err := ExtractNSInterfaces([]byte(in), NSInterfaceOptions{})
		var te *TopologyError
		if !errors.As(err, &te) {
			t.Errorf("%q: expected *TopologyError, got %T (%v)", in, err, err)
		}
	}
	if _, err := ExtractNSInterfaces(loadTopology(t), NSInterfaceOptions{LeafPatterns: []string{"("}}); err == nil {
		t.Error("expected an error for an invalid pattern")
	}
}
