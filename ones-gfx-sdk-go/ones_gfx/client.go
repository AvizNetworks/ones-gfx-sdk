package ones_gfx

// client.go is the Go twin of the Python SDK's ones_gfx/client.py.
//
// Single-token auth model, held per TokenAuth instance: log in (or supply a
// token), and it is attached as a raw `authorization` header on every
// subsequent request. Nothing is stored globally or on disk, so several
// instances can talk to different hosts as different users.
//
// This differs from JWTAuth in auth.go, which uses an access/refresh token pair
// and sends `Authorization: Bearer <token>`.

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ErrNotAuthenticated is returned when an authenticated call is attempted
// without a token.
var ErrNotAuthenticated = errors.New("not authenticated: call Login before calling Fabric Manager APIs")

// AuthData is the "data" object of POST /api/user/login and
// POST /api/user/refresh. Message differs by endpoint ("Login Successful" vs
// "Token refreshed"); IsPwdResetNeeded is returned by login only.
type AuthData struct {
	Message          string `json:"message"`
	Token            string `json:"token"`
	IsPwdResetNeeded bool   `json:"isPwdResetNeeded"`
}

// AuthResponse is the full login/refresh response body. Unlike the resource
// methods, these endpoints are not routed through the {"data": ...} unwrapping
// transport, so the envelope is preserved.
type AuthResponse struct {
	Data AuthData `json:"data"`
}

// TokenAuth implements AuthProvider using the single-token model. Construct it
// with the *root* base URL (e.g. https://host:3002) — the auth endpoints live
// under /api/user/, while the Fabric Manager resources live under /api/fm/
// (see FMBaseURL).
type TokenAuth struct {
	baseURL    string
	httpClient *http.Client

	mu    sync.RWMutex
	token string
}

// NewTokenAuth builds a TokenAuth against the root base URL.
func NewTokenAuth(baseURL string, verifyTLS bool, timeout time.Duration) *TokenAuth {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &TokenAuth{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTLS},
			},
		},
	}
}

// FMBaseURL turns a root base URL into the Fabric Manager base URL expected by
// NewClient. The trailing slash matters: resource methods pass bare paths
// (e.g. "addFabricData") which are resolved against this base, and without the
// slash net/url replaces the last segment (/api/addFabricData).
func FMBaseURL(root string) string {
	return strings.TrimRight(root, "/") + "/api/fm/"
}

// SetBaseURL points this auth at a different host.
func (a *TokenAuth) SetBaseURL(baseURL string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.baseURL = strings.TrimRight(baseURL, "/")
}

// Token returns the held token ("" when not logged in).
func (a *TokenAuth) Token() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.token
}

// SetToken replaces the held token without validation.
func (a *TokenAuth) SetToken(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.token = token
}

// UpdateToken replaces the held token, rejecting an empty value. Prefer this
// over SetToken when rotating a token.
func (a *TokenAuth) UpdateToken(token string) error {
	if token == "" {
		return errors.New("token is required")
	}
	a.SetToken(token)
	return nil
}

// Authenticated reports whether a token is held.
func (a *TokenAuth) Authenticated() bool { return a.Token() != "" }

// Apply attaches the raw token as the `authorization` header. No-op when
// unauthenticated, so the server returns 401 rather than the SDK panicking.
func (a *TokenAuth) Apply(headers map[string]string) {
	if tok := a.Token(); tok != "" {
		headers["authorization"] = tok
	}
}

// NeedsProactiveRefresh reports false: this model has no expiry metadata, so
// refresh happens reactively when the transport sees a 401.
func (a *TokenAuth) NeedsProactiveRefresh() bool { return false }

// Close releases nothing; present to satisfy AuthProvider.
func (a *TokenAuth) Close() error { return nil }

// Login authenticates and stores the returned token.
//
// POST {baseURL}/api/user/login with {"username", "password"}. Returns
// {"data": {"message", "token", "isPwdResetNeeded"}}.
func (a *TokenAuth) Login(username, password string) (*AuthResponse, error) {
	out, err := a.postAuth("/api/user/login", map[string]string{
		"username": username,
		"password": password,
	}, false)
	if err != nil {
		return nil, err
	}
	if out.Data.Token == "" {
		return nil, errors.New("login response carried no token")
	}
	a.SetToken(out.Data.Token)
	return out, nil
}

// Refresh exchanges the current token for a new one and stores it. Satisfies
// AuthProvider, so the transport calls this automatically on a 401.
func (a *TokenAuth) Refresh() error {
	_, err := a.RefreshWithResponse()
	return err
}

// RefreshWithResponse is Refresh but returns the parsed body, mirroring the
// Python client's refreshAuth() which hands back the full response.
func (a *TokenAuth) RefreshWithResponse() (*AuthResponse, error) {
	if !a.Authenticated() {
		return nil, ErrNotAuthenticated
	}
	out, err := a.postAuth("/api/user/refresh", nil, true)
	if err != nil {
		return nil, err
	}
	if out.Data.Token == "" {
		return nil, errors.New("refresh response carried no token")
	}
	a.SetToken(out.Data.Token)
	return out, nil
}

// Logout invalidates the session server-side, then clears the local token. The
// token is cleared regardless of the server's response.
func (a *TokenAuth) Logout() (*AuthResponse, error) {
	out, err := a.postAuth("/api/user/logout", nil, true)
	a.SetToken("")
	return out, err
}

// postAuth performs a POST against an /api/user/ endpoint and decodes the
// {"data": ...} envelope. withAuth attaches the current token.
func (a *TokenAuth) postAuth(path string, body interface{}, withAuth bool) (*AuthResponse, error) {
	var payload []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = b
	}
	a.mu.RLock()
	url := a.baseURL + path
	a.mu.RUnlock()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if withAuth {
		headers := map[string]string{}
		a.Apply(headers)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, &TransportError{Message: fmt.Sprintf("POST %s failed", url), Cause: err}
	}
	defer resp.Body.Close()

	var out AuthResponse
	decErr := json.NewDecoder(resp.Body).Decode(&out)
	if resp.StatusCode/100 != 2 {
		return nil, classifyAuthStatus(resp.StatusCode, path, out.Data.Message)
	}
	if decErr != nil {
		return nil, fmt.Errorf("decode %s response: %w", path, decErr)
	}
	return &out, nil
}

// classifyAuthStatus maps an auth-endpoint HTTP status onto the SDK's error
// taxonomy so callers can use errors.As just like with resource calls.
func classifyAuthStatus(status int, path, message string) error {
	if message == "" {
		message = fmt.Sprintf("POST %s returned HTTP %d", path, status)
	}
	base := APIError{Message: message, StatusCode: status}
	switch {
	case status == http.StatusUnauthorized, status == http.StatusForbidden:
		return &AuthenticationError{APIError: base}
	case status == http.StatusBadRequest:
		return &BadRequestError{APIError: base}
	case status == http.StatusNotFound:
		return &NotFoundError{APIError: base}
	case status == http.StatusConflict:
		return &ConflictError{APIError: base}
	case status >= 500:
		return &ServerError{APIError: base}
	default:
		return &base
	}
}
