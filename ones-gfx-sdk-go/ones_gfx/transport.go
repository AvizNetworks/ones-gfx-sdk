package ones_gfx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
