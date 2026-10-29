package ones_gfx

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// North-south (N-S) host interfaces from a fabric topology YAML.
//
// ONES keeps each fabric's topology as the YAML it was deployed from, served
// by the ONES gateway at GET /fabrics/{fabricName} (Fabrics.TopologyYAML).
// ExtractNSInterfaces reduces it to one question: which interfaces on each
// host are N-S, and where do they plug in.
//
// A host link is N-S when the leaf at its far end has a name matching one of
// the leaf patterns. The defaults cover the three N-S designs ONES generates:
// ^leaf-cn (converged user + storage), ^leaf-fn (front-end network) and
// ^leaf-sn (dedicated storage network). E-W leaves always carry an "su" token
// (leaf-su00-r0, leaf-p0-su00-r1) and never match. Pass your own patterns for
// a reference architecture that names its leaves differently.
//
// How the topology is read:
//   - Leaves come from Connectivity.Leaf (switchId -> switchName), hosts from
//     Connectivity.Host (hostId -> hostName), DPUs from Connectivity.DPU.
//   - Links are strings such as "G1_enp8s0f1np0 | L39_swp9": G<hostId> is a
//     host, L<switchId> a leaf, D<switchId> a DPU, and the rest of each side
//     is the interface.
//   - Links are read from the Host, Leaf and DPU sections alike, so a cable
//     listed on only one side is still found; one listed on both is reported
//     once.
//   - HBN fabrics put a DPU between host and leaf: host vf0 -> DPU
//     pf0vf0_if, and the DPU's uplink p0_if -> leaf. A VF on pf<N> leaves
//     through uplink p<N>, so it is reported with that uplink's leaf port and
//     Via/ViaInterface naming the DPU and its VF interface (the vfId the
//     VF-assign API takes). If the PF cannot be told from the name, every N-S
//     uplink of the DPU is reported.
//   - Endpoints that are not cluster hosts (H<n> external routers on border
//     leaves) are ignored.
//   - Only names and interfaces are copied out; the Credentials blocks in the
//     input are never decoded.
//
// The Python SDK's ones_gfx.extract_ns_interfaces implements the same rules
// and renders the same YAML; both are tested against one golden file.

// DefaultNSLeafPatterns are the leaf-name patterns used when none are given:
// the converged (leaf-cn), front-end (leaf-fn) and dedicated storage
// (leaf-sn) N-S leaves.
var DefaultNSLeafPatterns = []string{"^leaf-cn", "^leaf-fn", "^leaf-sn"}

// NSInterfaceOptions configures ExtractNSInterfaces.
type NSInterfaceOptions struct {
	// FabricName is echoed into the result; empty renders as null.
	FabricName string
	// LeafPatterns are regular expressions matched against leaf names; a
	// link is N-S when its leaf matches any of them. Empty uses
	// DefaultNSLeafPatterns.
	LeafPatterns []string
}

// NSInterface is one N-S path: the host interface and the leaf port it
// reaches. Via and ViaInterface are set when the path goes through a DPU (HBN):
// the DPU's name and the DPU-side VF interface the host interface maps to.
type NSInterface struct {
	Interface    string `yaml:"interface" json:"interface"`
	Via          string `yaml:"via,omitempty" json:"via,omitempty"`
	ViaInterface string `yaml:"viaInterface,omitempty" json:"viaInterface,omitempty"`
	Leaf         string `yaml:"leaf" json:"leaf"`
	LeafPort     string `yaml:"leafPort" json:"leafPort"`
}

// NSHost is a host and its N-S interfaces, sorted by interface.
type NSHost struct {
	Host         string        `yaml:"host" json:"host"`
	NSInterfaces []NSInterface `yaml:"nsInterfaces" json:"nsInterfaces"`
}

// NSInterfaces is the result of ExtractNSInterfaces. Hosts lists only hosts
// with at least one N-S interface, in topology order.
type NSInterfaces struct {
	Fabric         *string  `yaml:"fabric" json:"fabric"`
	LeafPatterns   []string `yaml:"leafPatterns" json:"leafPatterns"`
	HostCount      int      `yaml:"hostCount" json:"hostCount"`
	InterfaceCount int      `yaml:"interfaceCount" json:"interfaceCount"`
	Hosts          []NSHost `yaml:"hosts" json:"hosts"`
}

// YAML renders the result as clean, stable YAML: two-space indent, list items
// indented under their key. Identical to the Python SDK's
// ns_interfaces_to_yaml output.
func (n *NSInterfaces) YAML() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(n); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TopologyError reports input that is not a fabric topology this SDK can
// read: invalid YAML, or YAML without a Connectivity section (an HTML page
// here usually means the gateway did not serve GET /fabrics/{name}).
type TopologyError struct {
	Message string
	Cause   error
}

func (e *TopologyError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *TopologyError) Unwrap() error { return e.Cause }
func (e *TopologyError) ONESError()    {}

// flexInt decodes a YAML int or numeric string.
type flexInt struct {
	v  int
	ok bool
}

func (f *flexInt) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return nil
	}
	if v, err := strconv.Atoi(strings.TrimSpace(n.Value)); err == nil {
		f.v, f.ok = v, true
	}
	return nil
}

// topoLink decodes one entry of a Links list. Properties may be null or carry
// keys other than destination.
type topoLink struct {
	Link       string `yaml:"link"`
	Properties *struct {
		Destination string `yaml:"destination"`
	} `yaml:"properties"`
}

type topoLeaf struct {
	SwitchID   flexInt    `yaml:"switchId"`
	SwitchName string     `yaml:"switchName"`
	Links      []topoLink `yaml:"Links"`
}

type topoDPU struct {
	SwitchID   flexInt    `yaml:"switchId"`
	SwitchName string     `yaml:"switchName"`
	Links      []topoLink `yaml:"Links"`
}

type topoHost struct {
	HostID   flexInt    `yaml:"hostId"`
	HostName string     `yaml:"hostName"`
	Links    []topoLink `yaml:"Links"`
}

// topoDoc decodes only what is needed; Credentials and every other section
// are never read.
type topoDoc struct {
	Connectivity *struct {
		Leaf []topoLeaf `yaml:"Leaf"`
		Host []topoHost `yaml:"Host"`
		DPU  []topoDPU  `yaml:"DPU"`
	} `yaml:"Connectivity"`
}

var (
	endpointRE = regexp.MustCompile(`^([A-Za-z]+)(\d+)_(.+)$`)
	pfRE       = regexp.MustCompile(`^pf(\d+)`)
	uplinkRE   = regexp.MustCompile(`^p(\d+)(?:_|$)`)
)

type endpoint struct{ key, port string }

// parseEndpoint turns "G1_enp8s0f1np0" into {"G1", "enp8s0f1np0"}.
func parseEndpoint(s string) (endpoint, bool) {
	m := endpointRE.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return endpoint{}, false
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		return endpoint{}, false
	}
	return endpoint{strings.ToUpper(m[1]) + strconv.Itoa(n), m[3]}, true
}

func linkSides(l topoLink) (endpoint, endpoint, bool) {
	parts := strings.Split(l.Link, "|")
	if len(parts) != 2 {
		return endpoint{}, endpoint{}, false
	}
	a, okA := parseEndpoint(parts[0])
	b, okB := parseEndpoint(parts[1])
	return a, b, okA && okB
}

// ExtractNSInterfaces reduces a fabric topology YAML to each host's N-S
// interfaces. It returns *TopologyError when data is not a fabric topology,
// and an error naming the pattern when a leaf pattern does not compile.
func ExtractNSInterfaces(data []byte, opts NSInterfaceOptions) (*NSInterfaces, error) {
	patterns := opts.LeafPatterns
	if len(patterns) == 0 {
		patterns = DefaultNSLeafPatterns
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid leaf pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	var doc topoDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, &TopologyError{Message: "not a fabric topology (or not valid YAML)", Cause: err}
	}
	if doc.Connectivity == nil {
		return nil, &TopologyError{Message: "not a fabric topology: no 'Connectivity' section " +
			"(an HTML page here usually means the gateway did not serve GET /fabrics/{name})"}
	}
	conn := doc.Connectivity

	hostNames := map[string]string{}
	rank := map[string]int{}
	for _, h := range conn.Host {
		if h.HostID.ok && h.HostName != "" {
			key := "G" + strconv.Itoa(h.HostID.v)
			if _, seen := hostNames[key]; !seen {
				rank[h.HostName] = len(hostNames)
			}
			hostNames[key] = h.HostName
		}
	}
	leafNames := map[string]string{}
	for _, l := range conn.Leaf {
		if l.SwitchID.ok && l.SwitchName != "" {
			leafNames["L"+strconv.Itoa(l.SwitchID.v)] = l.SwitchName
		}
	}
	dpuNames := map[string]string{}
	for _, d := range conn.DPU {
		if d.SwitchID.ok && d.SwitchName != "" {
			dpuNames["D"+strconv.Itoa(d.SwitchID.v)] = d.SwitchName
		}
	}

	isNS := func(leaf string) bool {
		for _, re := range compiled {
			if re.MatchString(leaf) {
				return true
			}
		}
		return false
	}

	direct := map[hop]struct{}{}                   // far = leaf name
	vfs := map[hop]struct{}{}                      // far = DPU key
	uplinks := map[string]map[[3]string]struct{}{} // DPU key -> {dpuIface, leaf, port}

	classify := func(links []topoLink, fromHost bool) {
		for _, l := range links {
			a, b, ok := linkSides(l)
			if !ok {
				continue
			}
			for _, p := range [2][2]endpoint{{a, b}, {b, a}} {
				x, y := p[0], p[1]
				_, xHost := hostNames[x.key]
				_, xDPU := dpuNames[x.key]
				_, yDPU := dpuNames[y.key]
				yLeaf, yKnownLeaf := leafNames[y.key]
				switch {
				case xHost && strings.HasPrefix(y.key, "L"):
					if !yKnownLeaf && fromHost && l.Properties != nil && l.Properties.Destination != "" {
						// In the Host section, destination names the leaf.
						yLeaf, yKnownLeaf = l.Properties.Destination, true
					}
					if yKnownLeaf {
						direct[hop{x.key, x.port, yLeaf, y.port}] = struct{}{}
					}
				case xHost && yDPU:
					vfs[hop{x.key, x.port, y.key, y.port}] = struct{}{}
				case xDPU && yKnownLeaf:
					if uplinks[x.key] == nil {
						uplinks[x.key] = map[[3]string]struct{}{}
					}
					uplinks[x.key][[3]string{x.port, yLeaf, y.port}] = struct{}{}
				}
			}
		}
	}
	for _, h := range conn.Host {
		classify(h.Links, true)
	}
	for _, l := range conn.Leaf {
		classify(l.Links, false)
	}
	for _, d := range conn.DPU {
		classify(d.Links, false)
	}

	found := map[string]map[NSInterface]struct{}{}
	var order []string
	add := func(hostKey string, i NSInterface) {
		if !isNS(i.Leaf) {
			return
		}
		name := hostNames[hostKey]
		if _, ok := found[name]; !ok {
			found[name] = map[NSInterface]struct{}{}
			order = append(order, name)
		}
		found[name][i] = struct{}{}
	}

	for _, h := range sortedHops(direct) {
		add(h.hostKey, NSInterface{Interface: h.iface, Leaf: h.far, LeafPort: h.farPort})
	}
	for _, h := range sortedHops(vfs) {
		var ups [][3]string
		for u := range uplinks[h.far] {
			ups = append(ups, u)
		}
		sort.Slice(ups, func(i, j int) bool { return lessStrings(ups[i][:], ups[j][:]) })
		if m := pfRE.FindStringSubmatch(h.farPort); m != nil {
			var same [][3]string
			for _, u := range ups {
				if um := uplinkRE.FindStringSubmatch(u[0]); um != nil && um[1] == m[1] {
					same = append(same, u)
				}
			}
			if len(same) > 0 {
				ups = same
			}
		}
		for _, u := range ups {
			add(h.hostKey, NSInterface{Interface: h.iface, Via: dpuNames[h.far], ViaInterface: h.farPort, Leaf: u[1], LeafPort: u[2]})
		}
	}

	// Host-section order; hosts reached only from the Leaf side follow in
	// the order they were found.
	sort.SliceStable(order, func(i, j int) bool {
		ri, okI := rank[order[i]]
		rj, okJ := rank[order[j]]
		if !okI {
			ri = len(rank)
		}
		if !okJ {
			rj = len(rank)
		}
		return ri < rj
	})

	out := &NSInterfaces{LeafPatterns: append([]string(nil), patterns...), Hosts: []NSHost{}}
	if opts.FabricName != "" {
		f := opts.FabricName
		out.Fabric = &f
	}
	for _, name := range order {
		ifaces := make([]NSInterface, 0, len(found[name]))
		for i := range found[name] {
			ifaces = append(ifaces, i)
		}
		sort.Slice(ifaces, func(i, j int) bool {
			a, b := ifaces[i], ifaces[j]
			return lessStrings(
				[]string{a.Interface, a.Leaf, a.LeafPort, a.Via, a.ViaInterface},
				[]string{b.Interface, b.Leaf, b.LeafPort, b.Via, b.ViaInterface})
		})
		out.Hosts = append(out.Hosts, NSHost{Host: name, NSInterfaces: ifaces})
		out.InterfaceCount += len(ifaces)
	}
	out.HostCount = len(out.Hosts)
	return out, nil
}

// lessStrings orders two equal-length string tuples lexicographically, the
// way Python orders tuples, so both SDKs sort identically.
func lessStrings(a, b []string) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// hop is one cable from a host: to a leaf (far = leaf name) or to a DPU
// (far = DPU key).
type hop struct{ hostKey, iface, far, farPort string }

// sortedHops returns the set in a fixed order, so results never depend on map
// iteration.
func sortedHops(set map[hop]struct{}) []hop {
	out := make([]hop, 0, len(set))
	for h := range set {
		out = append(out, h)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		return lessStrings([]string{a.hostKey, a.iface, a.far, a.farPort}, []string{b.hostKey, b.iface, b.far, b.farPort})
	})
	return out
}
