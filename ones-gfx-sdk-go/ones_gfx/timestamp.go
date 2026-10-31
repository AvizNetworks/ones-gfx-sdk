package ones_gfx

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Timestamp is a time value as Fabric Manager sends it.
//
// FM is not consistent: Instant fields arrive as RFC 3339 with a zone
// ("2026-04-18T14:03:32.632Z"), LocalDateTime fields arrive without one
// ("2026-04-18T14:03:32.632113"), and a few older fields are epoch
// milliseconds. Timestamp keeps the value exactly as sent, so decoding never
// fails on an unexpected format, and Time interprets it on demand.
type Timestamp string

// UnmarshalJSON accepts a string, a number (epoch milliseconds) or null.
// Any other JSON is kept verbatim rather than failing the whole response.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "null" || raw == "" {
		*t = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*t = Timestamp(s)
		return nil
	}
	*t = Timestamp(raw)
	return nil
}

// MarshalJSON writes the value back as a JSON string, or null when empty.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(t))
}

// IsZero reports whether FM sent no value.
func (t Timestamp) IsZero() bool { return t == "" }

// String returns the value as FM sent it.
func (t Timestamp) String() string { return string(t) }

// timestampLayouts are tried in order. Values without a zone are LocalDateTime
// fields, which FM writes in the server's local time; they are read as UTC
// because the server zone is not in the payload.
var timestampLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05.999999999",
}

// Time parses the value. ok is false when the value is empty or in a format
// none of the known layouts match; String still returns it in that case.
func (t Timestamp) Time() (tm time.Time, ok bool) {
	s := strings.TrimSpace(string(t))
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range timestampLayouts {
		if v, err := time.Parse(layout, s); err == nil {
			return v, true
		}
	}
	if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.UnixMilli(ms).UTC(), true
	}
	return time.Time{}, false
}
