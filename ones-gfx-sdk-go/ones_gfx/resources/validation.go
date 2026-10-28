package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// Client-side validation. Every check here runs before a request is built, so
// a rejected call never reaches Fabric Manager and returns a
// *ones_gfx.ValidationError naming the offending field.
//
// The checks are deliberately limited to what is wrong for every fabric type:
// empty or path-breaking names, unknown enum values, and an impossible GPU
// quota. Rules that vary by fabric (tenant-name length and character set, for
// example) stay with Fabric Manager, so the SDK never rejects a name FM would
// accept.

func invalid(field, format string, args ...interface{}) error {
	return &ones_gfx.ValidationError{Field: field, Message: fmt.Sprintf(format, args...)}
}

// pathName checks a value interpolated into the URL path. The SDK does not
// escape path segments, so a "/" or "?" would silently address a different
// endpoint, and "." or ".." would be resolved away by the URL parser.
func pathName(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return invalid(field, "must not be empty")
	}
	if v == "." || v == ".." {
		return invalid(field, "%q is not a valid name", v)
	}
	if strings.ContainsAny(v, `/?#\%`) {
		return invalid(field, `%q contains a character that is not allowed in a URL path segment (/ ? # \ %%)`, v)
	}
	return nil
}

// required checks a non-path field that must be set.
func required(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return invalid(field, "must not be empty")
	}
	return nil
}

// fabricTenant is the common fabricName + tenantName pair.
func fabricTenant(fabricName, tenantName string) error {
	return firstErr(pathName("fabricName", fabricName), pathName("tenantName", tenantName))
}

// maxGpus allows nil (omitted), UnlimitedGPUs (-1) or a positive count.
func maxGpus(v *int) error {
	if v == nil || *v == ones_gfx.UnlimitedGPUs || *v > 0 {
		return nil
	}
	return invalid("maxGpusAllowed", "%d is not allowed: use %d for unlimited or a positive number", *v, ones_gfx.UnlimitedGPUs)
}

// gpuAction checks a required action.
func gpuAction(field string, a ones_gfx.GpuAction) error {
	if !a.Valid() {
		return invalid(field, "%q is not valid: want %s or %s", string(a), ones_gfx.GpuActionAdd, ones_gfx.GpuActionDelete)
	}
	return nil
}

// optGpuAction checks an optional action; nil leaves it to the server default.
func optGpuAction(field string, a *ones_gfx.GpuAction) error {
	if a == nil {
		return nil
	}
	return gpuAction(field, *a)
}

// optConfigScope checks an optional scope; nil leaves it to the server default.
func optConfigScope(s *ConfigScope) error {
	if s == nil || s.Valid() {
		return nil
	}
	return invalid("configScope", "%q is not valid: want %s or %s", string(*s), ConfigScopeWholeServer, ConfigScopeParticularGPU)
}

// gpuServers requires at least one server, each with a name.
func gpuServers(servers []GpuServerInfo) error {
	if len(servers) == 0 {
		return invalid("servers", "must list at least one server")
	}
	for i, s := range servers {
		if strings.TrimSpace(s.ServerName) == "" {
			return invalid(fmt.Sprintf("servers[%d].serverName", i), "must not be empty")
		}
	}
	return nil
}

// nonEmptyNames requires at least one non-empty name.
func nonEmptyNames(field string, names []string) error {
	if len(names) == 0 {
		return invalid(field, "must list at least one server")
	}
	for i, n := range names {
		if strings.TrimSpace(n) == "" {
			return invalid(fmt.Sprintf("%s[%d]", field, i), "must not be empty")
		}
	}
	return nil
}

// suidMap requires at least one server index carrying at least one host.
func suidMap(suid SuidMap) error {
	if len(suid) == 0 {
		return invalid("suid", "must map at least one server index to a host and its GPUs")
	}
	for idx, hosts := range suid {
		if _, err := strconv.Atoi(idx); err != nil {
			return invalid(fmt.Sprintf("suid[%s]", idx), "server index must be a whole number")
		}
		if len(hosts) == 0 {
			return invalid(fmt.Sprintf("suid[%s]", idx), "must name at least one host")
		}
		for host := range hosts {
			if strings.TrimSpace(host) == "" {
				return invalid(fmt.Sprintf("suid[%s]", idx), "has an empty host name")
			}
		}
	}
	return nil
}

// gpuIDs rejects negative port ids; an empty list means "every port".
func gpuIDs(ids []int) error {
	for i, id := range ids {
		if id < 0 {
			return invalid(fmt.Sprintf("gpuIds[%d]", i), "%d is negative", id)
		}
	}
	return nil
}

// newPolicyName checks a network policy name sent in a create body.
func newPolicyName(name string) error {
	if err := required("policyName", name); err != nil {
		return err
	}
	if len(name) > MaxNetworkPolicyNameLength {
		return invalid("policyName", "%q is %d characters; the limit is %d", name, len(name), MaxNetworkPolicyNameLength)
	}
	return nil
}

// nmxcDomains requires at least one domain, each with a host.
func nmxcDomains(domains []NmxcDomain) error {
	if len(domains) == 0 {
		return invalid("domains", "must list at least one domain")
	}
	for i, d := range domains {
		if strings.TrimSpace(d.Host) == "" {
			return invalid(fmt.Sprintf("domains[%d].host", i), "must not be empty")
		}
	}
	return nil
}

// peeringFields checks the three body fields every VPC peering call needs.
func peeringFields(name, vpcname, peervpcname string) error {
	return firstErr(required("name", name), required("vpcname", vpcname), required("peervpcname", peervpcname))
}

// firstErr returns the first non-nil error.
func firstErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
