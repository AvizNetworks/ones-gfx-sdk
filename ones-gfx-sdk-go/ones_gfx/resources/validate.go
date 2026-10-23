package resources

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// pathSegment validates a path parameter and escapes it for use as exactly one
// URL path segment, so a name containing '/', '?' or '#' cannot alter the route.
func pathSegment(value, label string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", &ones_gfx.ValidationError{Field: label, Reason: "is required and must not be blank"}
	}
	return url.PathEscape(value), nil
}

// pathSegments applies pathSegment to label/value pairs in order.
func pathSegments(pairs ...string) ([]string, error) {
	out := make([]string, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		seg, err := pathSegment(pairs[i+1], pairs[i])
		if err != nil {
			return nil, err
		}
		out = append(out, seg)
	}
	return out, nil
}

var poolCIDRRe = regexp.MustCompile(`^[^/]+/[0-9]{1,2}$`)

// validatePoolCIDR checks a tenant IPAM pool the way the server's IpamService does: an IPv4
// network with a prefix length between /8 and /30 (host bits are fine — the server
// canonicalises "192.3.0.5/17" to "192.3.0.0/17"). The Terraform provider enforces the same rule.
func validatePoolCIDR(label, value string) error {
	bad := func(reason string) error {
		return &ones_gfx.ValidationError{Field: label, Reason: reason}
	}
	// Surrounding whitespace is tolerated and stripped by the caller before sending: the server
	// matches the CIDR untrimmed, so a padded value would be rejected there.
	value = strings.TrimSpace(value)
	if !poolCIDRRe.MatchString(value) {
		return bad("must be an IPv4 network with a prefix length, e.g. 192.3.0.0/17 (INVALID_CIDR_FORMAT)")
	}
	parts := strings.Split(value, "/")
	octets := strings.Split(parts[0], ".")
	if len(octets) != 4 {
		return bad("only IPv4 CIDRs are supported (INVALID_CIDR_FORMAT)")
	}
	for _, o := range octets {
		n, err := strconv.Atoi(o)
		if err != nil || n < 0 || n > 255 || o != strings.TrimSpace(o) {
			return bad(fmt.Sprintf("%q is not a valid IPv4 address (INVALID_CIDR_FORMAT)", parts[0]))
		}
	}
	if prefix, _ := strconv.Atoi(parts[1]); prefix < 8 || prefix > 30 {
		return bad(fmt.Sprintf("prefix length must be between /8 and /30, got /%d (CIDR_POOL_TOO_SMALL)", prefix))
	}
	return nil
}
