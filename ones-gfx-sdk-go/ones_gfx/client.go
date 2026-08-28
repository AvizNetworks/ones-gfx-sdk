package ones_gfx

// client.go is the Go twin of the Python SDK's ones_gfx/client.py.
//
// Single-token auth model: log in once, persist the token to a local
// secrets.json, and attach it as a raw `authorization` header on every
// subsequent request. The token is also exposed as the package-level AuthToken,
// which is populated from secrets.json when a TokenAuth is constructed.
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
	"os"
	"strings"
	"sync"
	"time"
)

// SecretsFile is the local file that persists the auth token between runs.
// Exported as a var so callers can relocate it.
var SecretsFile = "secrets.json"

// AuthToken is the current auth token, populated from SecretsFile by
// NewTokenAuth and refreshed by Login/Refresh. Mirrors the module-level
// auth_token global in client.py. Read it for diagnostics; use TokenAuth's
// methods to change it.
var AuthToken string

// tokenMu guards AuthToken and SecretsFile writes. The transport may invoke
// Refresh concurrently with in-flight requests calling Apply.
var tokenMu sync.RWMutex

// ErrNotAuthenticated is returned when an authenticated call is attempted
// before Login has succeeded.
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
}

// NewTokenAuth builds a TokenAuth against the root base URL and loads any token
// already present in SecretsFile.
func NewTokenAuth(baseURL string, verifyTLS bool, timeout time.Duration) *TokenAuth {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	LoadToken()
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
// sdk.NewClient. The trailing slash matters: resource methods pass bare paths
// (e.g. "addFabricData") which are resolved against this base, and without the
// slash net/url replaces the last segment (/api/addFabricData).
func FMBaseURL(root string) string {
	return strings.TrimRight(root, "/") + "/api/fm/"
}

// LoadToken populates AuthToken from SecretsFile, tolerating a missing file or
// malformed JSON, and returns the resulting token.
func LoadToken() string {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	raw, err := os.ReadFile(SecretsFile)
	if err != nil {
		return AuthToken
	}
	var data map[string]interface{}
	if json.Unmarshal(raw, &data) != nil {
		return AuthToken
	}
	if tok, ok := data["auth_token"].(string); ok {
		AuthToken = tok
	}
	return AuthToken
}

// SaveToken sets AuthToken and persists it to SecretsFile, preserving any other
// keys already in the file. Passing "" clears the stored token.
func SaveToken(token string) error {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	return saveTokenLocked(token)
}

func saveTokenLocked(token string) error {
	AuthToken = token
	data := map[string]interface{}{}
	if raw, err := os.ReadFile(SecretsFile); err == nil {
		_ = json.Unmarshal(raw, &data)
	}
	data["auth_token"] = token
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SecretsFile, out, 0o600)
}

// Apply attaches the raw token as the `authorization` header. No-op when
// unauthenticated, so the server returns 401 rather than the SDK panicking.
func (a *TokenAuth) Apply(headers map[string]string) {
	tokenMu.RLock()
	defer tokenMu.RUnlock()
	if AuthToken != "" {
		headers["authorization"] = AuthToken
	}
}

// NeedsProactiveRefresh reports false: this model has no expiry metadata, so
// refresh happens reactively when the transport sees a 401.
func (a *TokenAuth) NeedsProactiveRefresh() bool { return false }

// Close releases nothing; present to satisfy AuthProvider.
func (a *TokenAuth) Close() error { return nil }

// Token returns the current token.
func (a *TokenAuth) Token() string {
	tokenMu.RLock()
	defer tokenMu.RUnlock()
	return AuthToken
}

// Authenticated reports whether a token is currently held, either from a
// successful Login or loaded from SecretsFile.
func Authenticated() bool {
	tokenMu.RLock()
	defer tokenMu.RUnlock()
	return AuthToken != ""
}

// Login authenticates and persists the returned token.
//
// POST {baseURL}/api/user/login with {"username", "password"}. Returns
// {"data": {"message", "token", "isPwdResetNeeded"}}; stores data.token and
// returns the full parsed response.
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
	if err := SaveToken(out.Data.Token); err != nil {
		return nil, fmt.Errorf("persist token: %w", err)
	}
	return out, nil
}

// Refresh exchanges the current token for a new one and persists it.
//
// POST {baseURL}/api/user/refresh with the current token in the `authorization`
// header. Satisfies AuthProvider, so the transport calls this automatically on
// a 401. Returns ErrNotAuthenticated when no token is held.
func (a *TokenAuth) Refresh() error {
	if a.Token() == "" {
		return ErrNotAuthenticated
	}
	out, err := a.postAuth("/api/user/refresh", nil, true)
	if err != nil {
		return err
	}
	if out.Data.Token == "" {
		return errors.New("refresh response carried no token")
	}
	return SaveToken(out.Data.Token)
}

// RefreshWithResponse is Refresh but returns the parsed body, mirroring the
// Python client's refresh() which hands back the full response.
func (a *TokenAuth) RefreshWithResponse() (*AuthResponse, error) {
	if a.Token() == "" {
		return nil, ErrNotAuthenticated
	}
	out, err := a.postAuth("/api/user/refresh", nil, true)
	if err != nil {
		return nil, err
	}
	if out.Data.Token != "" {
		if err := SaveToken(out.Data.Token); err != nil {
			return nil, fmt.Errorf("persist token: %w", err)
		}
	}
	return out, nil
}

// Logout invalidates the session server-side, then clears the local token.
//
// POST {baseURL}/api/user/logout with the auth header. The local token and the
// auth_token field in SecretsFile are cleared regardless of the server's
// response, matching the Python client.
func (a *TokenAuth) Logout() (*AuthResponse, error) {
	out, err := a.postAuth("/api/user/logout", nil, true)
	if clearErr := SaveToken(""); clearErr != nil && err == nil {
		err = fmt.Errorf("clear token: %w", clearErr)
	}
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
	req, err := http.NewRequest(http.MethodPost, a.baseURL+path, bytes.NewReader(payload))
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
		return nil, &TransportError{Message: fmt.Sprintf("POST %s failed", a.baseURL+path), Cause: err}
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
