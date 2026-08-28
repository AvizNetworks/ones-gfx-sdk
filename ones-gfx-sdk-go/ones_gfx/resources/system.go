package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// SystemResource covers controller info, status logs, streaming control,
// log levels and role info — cross-cutting endpoints without a domain home.
type SystemResource struct {
	transport *ones_gfx.Transport
}

// NewSystemResource constructs a SystemResource.
func NewSystemResource(transport *ones_gfx.Transport) *SystemResource {
	return &SystemResource{transport: transport}
}

// ControllerVersion is the build-info body of GET /getControllerVersion and
// GET /getControllerVersionInternal.
type ControllerVersion struct {
	AppName          string `json:"appName,omitempty"`
	AppArtifactID    string `json:"appArtifactId,omitempty"`
	AppVersion       string `json:"appVersion,omitempty"`
	AppBuildDateTime string `json:"appBuildDateTime,omitempty"`
}

// ControllerVersion returns FM build info. Maps to GET /getControllerVersion.
func (r *SystemResource) ControllerVersion(ctx context.Context) (*ControllerVersion, error) {
	return ones_gfx.Call[ControllerVersion](r.transport, "GET", "getControllerVersion", nil, ones_gfx.OperationModeSynchronous, nil)
}

// ControllerVersionInternal returns FM build info (internal variant).
// Maps to GET /getControllerVersionInternal.
func (r *SystemResource) ControllerVersionInternal(ctx context.Context) (*ControllerVersion, error) {
	return ones_gfx.Call[ControllerVersion](r.transport, "GET", "getControllerVersionInternal", nil, ones_gfx.OperationModeSynchronous, nil)
}

// Status reads the status log of a file. Returns nil when the server
// responds with null. Maps to GET /status?fileName=...
func (r *SystemResource) Status(ctx context.Context, fileName *string) ([]string, error) {
	q := url.Values{}
	if fileName != nil {
		q.Set("fileName", *fileName)
	}
	res, err := ones_gfx.Call[[]string](r.transport, "GET", "status", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// UploadStatus returns the upload status marker.
// Maps to GET /uploadStatus.
func (r *SystemResource) UploadStatus(ctx context.Context) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "GET", "uploadStatus", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// StartStreaming marks a file for streaming.
// Maps to GET /start?filename=...
func (r *SystemResource) StartStreaming(ctx context.Context, filename string) (string, error) {
	q := url.Values{}
	q.Set("filename", filename)
	res, err := ones_gfx.Call[string](r.transport, "GET", "start", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// StopStreaming stops streaming of a file.
// Maps to GET /stop?filename=...
func (r *SystemResource) StopStreaming(ctx context.Context, filename string) (string, error) {
	q := url.Values{}
	q.Set("filename", filename)
	res, err := ones_gfx.Call[string](r.transport, "GET", "stop", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// UpdateRoleInfo renames a role layer. Both fields are required server-side.
// Maps to POST /updateroleinfo.
func (r *SystemResource) UpdateRoleInfo(ctx context.Context, layer int, currentName string) (string, error) {
	body := struct {
		Layer       int    `json:"layer"`
		CurrentName string `json:"current_name"`
	}{layer, currentName}
	res, err := ones_gfx.Call[string](r.transport, "POST", "updateroleinfo", body, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// SetLogLevel changes runtime log levels: logger name → level.
// Maps to POST /log/level.
func (r *SystemResource) SetLogLevel(ctx context.Context, loggers map[string]string) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "POST", "log/level", loggers, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// GetLogLevel returns the current logger → level map.
// Maps to GET /log/level.
func (r *SystemResource) GetLogLevel(ctx context.Context) (map[string]string, error) {
	res, err := ones_gfx.Call[map[string]string](r.transport, "GET", "log/level", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
