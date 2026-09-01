package ones_gfx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Transport wraps HTTP operations and knows about ONES API conventions
// (envelope unwrapping, Prefer header, error mapping). This is exported so
// the resources package can use it.
type Transport struct {
	baseURL    string
	auth       AuthProvider
	httpClient *http.Client
	timeout    time.Duration
}

// NewTransport constructs a Transport with the given configuration.
func NewTransport(baseURL string, auth AuthProvider, cfg clientConfig) *Transport {
	transport := &http.Transport{
		TLSClientConfig: cfg.tlsConfig,
	}
	if !cfg.verifyTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{
		Timeout:   cfg.timeout,
		Transport: transport,
	}

	return &Transport{
		baseURL:    baseURL,
		auth:       auth,
		httpClient: httpClient,
		timeout:    cfg.timeout,
	}
}

// NewTransportWithOptions constructs a Transport using client options.
func NewTransportWithOptions(baseURL string, auth AuthProvider, opts ...ClientOption) *Transport {
	cfg := defaultClientConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return NewTransport(baseURL, auth, cfg)
}

// Get performs a GET request and returns the unwrapped data.
func (t *Transport) Get(path string, timeout *time.Duration) (interface{}, error) {
	return t.request("GET", path, nil, OperationModeSynchronous, timeout)
}

// Post performs a POST request with optional mode and timeout override.
func (t *Transport) Post(path string, body interface{}, mode OperationMode, timeout *time.Duration) (interface{}, error) {
	return t.request("POST", path, body, mode, timeout)
}

// Patch performs a PATCH request with optional mode and timeout override.
func (t *Transport) Patch(path string, body interface{}, mode OperationMode, timeout *time.Duration) (interface{}, error) {
	return t.request("PATCH", path, body, mode, timeout)
}

// Delete performs a DELETE request with optional mode and timeout override.
func (t *Transport) Delete(path string, body interface{}, mode OperationMode, timeout *time.Duration) (interface{}, error) {
	return t.request("DELETE", path, body, mode, timeout)
}

// SetBaseURL repoints the transport at a different base URL. Pass the Fabric
// Manager base (see FMBaseURL), not the root URL.
func (t *Transport) SetBaseURL(baseURL string) {
	t.baseURL = baseURL
}

// Close releases HTTP resources.
func (t *Transport) Close() error {
	t.httpClient.CloseIdleConnections()
	return t.auth.Close()
}

// request is the internal method that handles auth, retry-on-401, and error mapping.
func (t *Transport) request(method, path string, body interface{}, mode OperationMode, timeout *time.Duration) (interface{}, error) {
	// Build URL
	fullURL, err := t.buildURL(path)
	if err != nil {
		return nil, err
	}

	// Proactive refresh if token is near expiry
	if t.auth.NeedsProactiveRefresh() {
		if err := t.auth.Refresh(); err != nil {
			return nil, err
		}
	}

	// Send request
	resp, err := t.send(method, fullURL, body, mode, timeout)
	if err != nil {
		return nil, err
	}

	// Reactive refresh: if we got a 401, refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		if err := t.auth.Refresh(); err != nil {
			return nil, err
		}
		resp, err = t.send(method, fullURL, body, mode, timeout)
		if err != nil {
			return nil, err
		}
	}

	return t.handleResponse(resp)
}

// send builds headers and sends the HTTP request.
func (t *Transport) send(method, fullURL string, body interface{}, mode OperationMode, timeout *time.Duration) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, &TransportError{Message: "failed to marshal request body", Cause: err}
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, &TransportError{Message: "failed to build request", Cause: err}
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	t.auth.Apply(headers)

	// Add Prefer header for sync/async mode
	if mode == OperationModeSynchronous {
		headers["Prefer"] = "respond-sync"
	} else if mode == OperationModeAsyncPoll || mode == OperationModeAsyncWebhook {
		headers["Prefer"] = "respond-async"
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Apply timeout override if provided
	client := t.httpClient
	if timeout != nil {
		client = &http.Client{
			Timeout:   *timeout,
			Transport: t.httpClient.Transport,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &TransportError{Message: fmt.Sprintf("%s %s failed", method, fullURL), Cause: err}
	}

	return resp, nil
}

// handleResponse parses the response and maps errors.
func (t *Transport) handleResponse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &TransportError{Message: "failed to read response body", Cause: err}
	}

	// Parse body as JSON
	var body interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			// Not JSON - treat as plain text
			body = string(bodyBytes)
		}
	}

	// Handle success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return unwrapEnvelope(body), nil
	}

	// Handle errors
	message := extractErrorMessage(body)
	if message == "" {
		message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	// Map status code to typed error
	return nil, classifyError(resp.StatusCode, message, body)
}

// buildURL joins baseURL and path.
func (t *Transport) buildURL(path string) (string, error) {
	base, err := url.Parse(t.baseURL)
	if err != nil {
		return "", &TransportError{Message: "invalid base URL", Cause: err}
	}
	rel, err := url.Parse(path)
	if err != nil {
		return "", &TransportError{Message: "invalid path", Cause: err}
	}
	return base.ResolveReference(rel).String(), nil
}

// unwrapEnvelope handles ONES response envelopes.
// Success responses: {"status": "success", "message": "...", "data": ...}
// The data field may be a JSON-encoded string or an object.
func unwrapEnvelope(body interface{}) interface{} {
	obj, ok := body.(map[string]interface{})
	if !ok {
		return body
	}

	// Detect envelope by presence of "status" and "data"/"message"
	status, hasStatus := obj["status"]
	if !hasStatus {
		return body // Not an envelope
	}

	if status != "success" {
		return body
	}

	// Extract data field
	data, hasData := obj["data"]
	if hasData {
		if str, ok := data.(string); ok {
			if str != "" {
				var parsed interface{}
				if err := json.Unmarshal([]byte(str), &parsed); err == nil {
					return parsed
				}
				return str
			}
			// empty string — fall through to message
		} else {
			return data
		}
	}

	// No data field or data was empty string — return message
	if msg, hasMsg := obj["message"]; hasMsg {
		return msg
	}

	return body
}

// extractErrorMessage pulls a useful error message from the response.
func extractErrorMessage(body interface{}) string {
	if str, ok := body.(string); ok {
		return str
	}

	obj, ok := body.(map[string]interface{})
	if !ok {
		return ""
	}

	// Try common error keys
	for _, key := range []string{"error", "message", "detail"} {
		if val, ok := obj[key]; ok {
			if str, ok := val.(string); ok && str != "" {
				return str
			}
		}
	}

	return ""
}

// classifyError maps HTTP status codes to typed errors.
func classifyError(statusCode int, message string, body interface{}) error {
	apiErr := APIError{
		Message:      message,
		StatusCode:   statusCode,
		ResponseBody: body,
	}

	switch statusCode {
	case http.StatusBadRequest:
		return &BadRequestError{APIError: apiErr}
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthenticationError{APIError: apiErr}
	case http.StatusNotFound:
		return &NotFoundError{APIError: apiErr}
	case http.StatusConflict:
		return &ConflictError{APIError: apiErr}
	default:
		if statusCode >= 500 && statusCode < 600 {
			return &ServerError{APIError: apiErr}
		}
		return &apiErr
	}
}

// ─── Typed request helpers ──────────────────────────────────────────────────
//
// The resource methods use Call / CallMultipart so every request body is a
// typed value and every response decodes directly into a typed result —
// no map[string]interface{} building and no remarshal round-trips.

// ReqOpts carries per-request transport options beyond the JSON body.
type ReqOpts struct {
	// Timeout overrides the client-level HTTP timeout for this call.
	Timeout *time.Duration
	// IdempotencyKey is sent as the Idempotency-Key header (async replay).
	IdempotencyKey string
	// RequestOrigin is sent as the x-request-origin header ("ones-ui" = INTERNAL).
	RequestOrigin string
	// Query holds URL query parameters (nil values are dropped).
	Query url.Values
}

// Call executes an HTTP request with a JSON body and decodes the response
// into T. The {"data": ...} envelope is unwrapped before decoding, matching
// the Python SDK client. Returns a nil result when the server sends an empty
// body (e.g. HTTP 204). Non-JSON responses fall back to raw text only when
// T is string.
func Call[T any](tr *Transport, method, path string, body any, mode OperationMode, ro *ReqOpts) (*T, error) {
	var payload []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, &TransportError{Message: "failed to marshal request body", Cause: err}
		}
		payload = b
	}
	raw, err := tr.doRaw(method, path, "application/json", payload, mode, ro)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	raw = unwrapEnvelopeBytes(raw)
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		// Plain-text fallback: Spring returns bare strings (text/plain) for
		// some endpoints (e.g. GET /log/level style handlers).
		if sp, ok := any(&out).(*string); ok {
			*sp = string(raw)
			return &out, nil
		}
		return nil, &TransportError{
			Message: fmt.Sprintf("failed to decode %s %s response into %T", method, path, out),
			Cause:   err,
		}
	}
	return &out, nil
}

// CallMultipart executes a multipart/form-data request (file upload endpoints)
// and decodes the response into T, applying the same envelope and error rules
// as Call. Each entry in files maps a form field name to a file path on disk.
func CallMultipart[T any](tr *Transport, method, path string, fields map[string]string, files map[string]string, mode OperationMode, ro *ReqOpts) (*T, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, &TransportError{Message: "failed to write form field", Cause: err}
		}
	}
	for field, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, &TransportError{Message: "failed to open upload file " + filePath, Cause: err}
		}
		part, err := w.CreateFormFile(field, filepath.Base(filePath))
		if err != nil {
			f.Close()
			return nil, &TransportError{Message: "failed to create form file", Cause: err}
		}
		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return nil, &TransportError{Message: "failed to stream upload file", Cause: err}
		}
		f.Close()
	}
	if err := w.Close(); err != nil {
		return nil, &TransportError{Message: "failed to finalize multipart body", Cause: err}
	}
	raw, err := tr.doRaw(method, path, w.FormDataContentType(), buf.Bytes(), mode, ro)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	var out T
	if err := json.Unmarshal(unwrapEnvelopeBytes(raw), &out); err != nil {
		return nil, &TransportError{
			Message: fmt.Sprintf("failed to decode %s %s response into %T", method, path, out),
			Cause:   err,
		}
	}
	return &out, nil
}

// doRaw runs the full request pipeline (URL build, proactive + reactive auth
// refresh, Prefer/idempotency headers) and returns the raw 2xx body bytes.
// A nil slice means the server sent no body.
func (t *Transport) doRaw(method, path, contentType string, payload []byte, mode OperationMode, ro *ReqOpts) ([]byte, error) {
	if ro != nil && len(ro.Query) > 0 {
		path = path + "?" + ro.Query.Encode()
	}
	fullURL, err := t.buildURL(path)
	if err != nil {
		return nil, err
	}

	if t.auth.NeedsProactiveRefresh() {
		if err := t.auth.Refresh(); err != nil {
			return nil, err
		}
	}

	resp, err := t.sendRaw(method, fullURL, contentType, payload, mode, ro)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		if err := t.auth.Refresh(); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp, err = t.sendRaw(method, fullURL, contentType, payload, mode, ro)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &TransportError{Message: "failed to read response body", Cause: err}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if len(bodyBytes) == 0 {
			return nil, nil
		}
		return bodyBytes, nil
	}

	var parsed any
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
			parsed = string(bodyBytes)
		}
	}
	message := extractErrorMessage(parsed)
	if message == "" {
		message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return nil, classifyError(resp.StatusCode, message, parsed)
}

// sendRaw builds and sends one HTTP request with an explicit content type.
func (t *Transport) sendRaw(method, fullURL, contentType string, payload []byte, mode OperationMode, ro *ReqOpts) (*http.Response, error) {
	var reqBody io.Reader
	if payload != nil {
		reqBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, &TransportError{Message: "failed to build request", Cause: err}
	}
	authHeaders := map[string]string{"Content-Type": contentType, "Accept": "application/json"}
	t.auth.Apply(authHeaders)
	for k, v := range authHeaders {
		req.Header.Set(k, v)
	}
	switch mode {
	case OperationModeSynchronous:
		req.Header.Set("Prefer", "respond-sync")
	case OperationModeAsyncPoll, OperationModeAsyncWebhook:
		req.Header.Set("Prefer", "respond-async")
	}
	if ro != nil {
		if ro.IdempotencyKey != "" {
			req.Header.Set("Idempotency-Key", ro.IdempotencyKey)
		}
		if ro.RequestOrigin != "" {
			req.Header.Set("x-request-origin", ro.RequestOrigin)
		}
	}

	client := t.httpClient
	if ro != nil && ro.Timeout != nil {
		client = &http.Client{Timeout: *ro.Timeout, Transport: t.httpClient.Transport}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &TransportError{Message: fmt.Sprintf("%s %s failed", method, fullURL), Cause: err}
	}
	return resp, nil
}

// unwrapEnvelopeBytes applies the ONES envelope rule to raw body bytes: when
// the JSON body is an object containing a "data" key, the value of that key is
// returned (re-encoded); otherwise the body is returned unchanged. This mirrors
// ONESClient.call_api in the Python SDK.
func unwrapEnvelopeBytes(raw []byte) []byte {
	if len(raw) == 0 {
		return raw
	}
	var body any
	if err := json.Unmarshal(raw, &body); err != nil {
		return raw
	}
	obj, ok := body.(map[string]any)
	if !ok {
		return raw
	}
	data, hasData := obj["data"]
	if !hasData {
		return raw
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return raw
	}
	return encoded
}
