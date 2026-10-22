package ones_gfx

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Which Fabric Manager routes the SDK must reach on Fabric Manager directly.
//
// The ONES gateway proxies /api/fm/* through an explicit per-route allowlist
// (fm-router.js). Routes missing from that allowlist do not 404 - they fall
// through to the UI's single-page-app catch-all and come back as HTML with
// HTTP 200, which is far more confusing than a refusal.
//
// Thirty-four routes are in that state as of the 4.4.0 release: the endpoints
// added for network policy, webhook registration, IPAM configuration and
// address queries, tenant DHCP and rail CIDRPools, plus the
// shared-services and VF-interface routes that have been unreachable since
// GFX-5709 merged. The 4.4.0 train has shipped, so the gateway cannot be
// changed for this release and the SDK routes these calls to Fabric Manager
// itself. Everything else keeps going through the gateway.
//
// This file is meant to be deleted. Once a gateway release registers these
// routes, empty DirectFMPaths and every call goes back through the gateway
// with no other change.

// FMDirectPort is the port Fabric Manager serves on. The deployment publishes
// it on the same host as the gateway, so the direct base URL is derived rather
// than configured.
const FMDirectPort = 8443

// DirectFMPaths are path templates, as passed to the transport, that bypass
// the gateway. "{}" stands for exactly one path segment.
var DirectFMPaths = []string{
	// Network policy
	"fabrics/{}/networkpolicy",
	"fabrics/{}/networkpolicy/{}",
	// Webhook registration: global, fabric and tenant scopes
	"webhooks",
	"webhooks/{}",
	"fabrics/{}/webhook",
	"fabrics/{}/tenants/{}/webhook",
	// IPAM address queries
	"fabrics/{}/tenants/{}/ip-addresses",
	"fabrics/{}/computes/{}/ip-addresses",
	"fabrics/{}/tenants/{}/computes/{}/ip-address",
	// IPAM configuration
	"fabrics/{}/ipam/config",
	// Tenant DHCP
	"fabrics/{}/tenants/{}/dhcp",
	"fabrics/{}/tenants/{}/dhcp/retry",
	"fabrics/{}/dhcp/reconcile",
	// Rail CIDRPools
	"fabrics/{}/rails/cidrpools",
	"fabrics/{}/rails/{}/cidrpool",
	// Fabric listing
	"fabrics",
	// Shared services and VF interfaces (GFX-5709)
	"fabrics/{}/shared-services",
	"fabrics/{}/servers/{}/vf-interfaces",
	"fabrics/{}/servers/{}/vf-interfaces/{}/assign",
}

var directPatterns = compileDirectPaths(DirectFMPaths)

func compileDirectPaths(templates []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(templates))
	for _, t := range templates {
		parts := strings.Split(t, "{}")
		for i, p := range parts {
			parts[i] = regexp.QuoteMeta(p)
		}
		// A segment is anything but "/", so a template matches one concrete
		// path and nothing deeper: "fabrics/{}" must not swallow "fabrics/x/y".
		out = append(out, regexp.MustCompile("^"+strings.Join(parts, "[^/]+")+"$"))
	}
	return out
}

// IsDirectPath reports whether path must bypass the gateway and go straight to
// Fabric Manager.
func IsDirectPath(path string) bool {
	p := strings.Trim(path, "/")
	if i := strings.IndexByte(p, '?'); i >= 0 {
		p = p[:i]
	}
	for _, re := range directPatterns {
		if re.MatchString(p) {
			return true
		}
	}
	return false
}

// DeriveFMBaseURL returns Fabric Manager's own base URL, from the gateway's.
//
// The deployment publishes FM on the same host as the gateway, so only the
// port changes. Set a client's direct URL explicitly where that does not hold.
func DeriveFMBaseURL(baseURL string) string {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil || u.Host == "" {
		return strings.TrimRight(baseURL, "/")
	}
	host := u.Hostname()
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}
	return scheme + "://" + net.JoinHostPort(host, strconv.Itoa(FMDirectPort))
}
