package ones

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testFabric = "Fab101"
	testName   = "tenant1-storage-route-leak"
	testVPC    = "tenant1-Fab101-north-south"
	testPeer   = "Fab101-Storage"
)

// recordedRequest is what the stub FM server saw for one call.
type recordedRequest struct {
	Method string
	Path   string
	Auth   string
	Prefer string
	Body   map[string]interface{}
}

// peeringErrorBody is the envelope FM's peeringError() returns.
func peeringErrorBody(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "error",
		"message": message,
		"error":   map[string]interface{}{"code": code, "detail": message},
	}
}

// newPeeringServer serves DELETE .../vpcpeering with the given status and
// body, and answers /api/user/refresh with refreshStatus (issuing "tok2").
func newPeeringServer(t *testing.T, status int, body interface{}, refreshStatus int) (*httptest.Server, *[]recordedRequest) {
	t.Helper()
	var calls []recordedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/user/refresh" {
			w.WriteHeader(refreshStatus)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"token": "tok2", "message": "refreshed"},
			})
			return
		}
		rec := recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Auth:   r.Header.Get("authorization"),
			Prefer: r.Header.Get("Prefer"),
		}
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &rec.Body); err != nil {
				t.Errorf("request body is not JSON: %v", err)
			}
		}
		calls = append(calls, rec)
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &calls
}

func newTestClient(t *testing.T, baseURL, token string) *Client {
	t.Helper()
	c, err := NewClient(baseURL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	c.SetAuthToken(token)
	return c
}

func TestDeleteVpcPeering_Success(t *testing.T) {
	msg := "VPC peering '" + testName + "' deleted successfully."
	srv, calls := newPeeringServer(t, http.StatusOK,
		map[string]interface{}{"status": "success", "message": msg}, http.StatusOK)
	client := newTestClient(t, srv.URL, "tok")

	res, err := DeleteVpcPeering(context.Background(), client, testFabric, testName, testVPC, testPeer)
	if err != nil {
		t.Fatalf("DeleteVpcPeering: unexpected error: %v", err)
	}
	if res == nil || res.Status != "success" || res.Message != msg {
		t.Fatalf("DeleteVpcPeering: got %+v, want status=success message=%q", res, msg)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 request, got %d", len(*calls))
	}
	got := (*calls)[0]
	if got.Method != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", got.Method)
	}
	if want := "/api/fm/fabrics/" + testFabric + "/vpcpeering"; got.Path != want {
		t.Errorf("path = %s, want %s", got.Path, want)
	}
	if got.Auth != "tok" {
		t.Errorf("authorization = %q, want tok", got.Auth)
	}
	if got.Prefer != "respond-sync" {
		t.Errorf("Prefer = %q, want respond-sync", got.Prefer)
	}
	wantBody := map[string]interface{}{"name": testName, "vpcname": testVPC, "peervpcname": testPeer}
	if len(got.Body) != len(wantBody) {
		t.Errorf("body = %v, want %v", got.Body, wantBody)
	}
	for k, v := range wantBody {
		if got.Body[k] != v {
			t.Errorf("body[%s] = %v, want %v", k, got.Body[k], v)
		}
	}
}

func TestDeleteVpcPeering_ResourceHandle(t *testing.T) {
	srv, calls := newPeeringServer(t, http.StatusOK,
		map[string]interface{}{"status": "success", "message": "deleted"}, http.StatusOK)
	client := newTestClient(t, srv.URL, "tok")

	res, err := client.VPCPeering.Delete(context.Background(), testFabric, testName, testVPC, testPeer)
	if err != nil {
		t.Fatalf("VPCPeering.Delete: unexpected error: %v", err)
	}
	if res.Message != "deleted" {
		t.Errorf("message = %q, want deleted", res.Message)
	}
	if len(*calls) != 1 || (*calls)[0].Method != http.MethodDelete {
		t.Errorf("expected a single DELETE, got %+v", *calls)
	}
}

func TestDeleteVpcPeering_APIErrors(t *testing.T) {
	cases := []struct {
		name   string
		status int
		code   string
		check  func(error) bool
	}{
		{"peering not found", http.StatusNotFound, "VPC_PEERING_NOT_FOUND",
			func(err error) bool { var e *NotFoundError; return errors.As(err, &e) }},
		{"fabric not found", http.StatusNotFound, "FABRIC_NOT_FOUND",
			func(err error) bool { var e *NotFoundError; return errors.As(err, &e) }},
		{"pair mismatch", http.StatusBadRequest, "VPC_PEERING_PAIR_MISMATCH",
			func(err error) bool { var e *BadRequestError; return errors.As(err, &e) }},
		{"missing name", http.StatusBadRequest, "VPC_PEERING_NAME_REQUIRED",
			func(err error) bool { var e *BadRequestError; return errors.As(err, &e) }},
		{"forbidden", http.StatusForbidden, "FORBIDDEN",
			func(err error) bool { var e *AuthenticationError; return errors.As(err, &e) }},
		{"switch failure", http.StatusInternalServerError, "VPC_PEERING_DELETE_FAILED",
			func(err error) bool { var e *ServerError; return errors.As(err, &e) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			message := tc.code + " message"
			srv, _ := newPeeringServer(t, tc.status, peeringErrorBody(tc.code, message), http.StatusOK)
			client := newTestClient(t, srv.URL, "tok")

			res, err := DeleteVpcPeering(context.Background(), client, testFabric, testName, testVPC, testPeer)
			if err == nil {
				t.Fatalf("expected error, got result %+v", res)
			}
			if res != nil {
				t.Errorf("expected nil result on error, got %+v", res)
			}
			if !tc.check(err) {
				t.Fatalf("error type %T (%v) did not match expected type", err, err)
			}
			var onesErr ONESError
			if !errors.As(err, &onesErr) {
				t.Errorf("error %T does not implement ONESError", err)
			}

			// The status, message and machine-readable code reach the caller.
			apiErr := apiErrorOf(err)
			if apiErr == nil {
				t.Fatalf("could not extract APIError from %T", err)
			}
			if apiErr.StatusCode != tc.status {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tc.status)
			}
			if apiErr.Message != message {
				t.Errorf("Message = %q, want %q", apiErr.Message, message)
			}
			body, _ := apiErr.ResponseBody.(map[string]interface{})
			detail, _ := body["error"].(map[string]interface{})
			if detail["code"] != tc.code {
				t.Errorf("ResponseBody error.code = %v, want %s", detail["code"], tc.code)
			}
		})
	}
}

func TestDeleteVpcPeering_UnauthorizedRefreshFails(t *testing.T) {
	srv, calls := newPeeringServer(t, http.StatusUnauthorized,
		peeringErrorBody("UNAUTHORIZED", "token expired"), http.StatusUnauthorized)
	client := newTestClient(t, srv.URL, "tok")

	_, err := DeleteVpcPeering(context.Background(), client, testFabric, testName, testVPC, testPeer)
	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected *AuthenticationError, got %T (%v)", err, err)
	}
	if len(*calls) != 1 {
		t.Errorf("expected no retry after failed refresh, got %d calls", len(*calls))
	}
}

func TestDeleteVpcPeering_UnauthorizedRefreshRetries(t *testing.T) {
	var deletes int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/user/refresh" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"token": "tok2"}})
			return
		}
		deletes++
		if r.Header.Get("authorization") != "tok2" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "deleted"})
	}))
	defer srv.Close()
	client := newTestClient(t, srv.URL, "tok")

	res, err := DeleteVpcPeering(context.Background(), client, testFabric, testName, testVPC, testPeer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Message != "deleted" || deletes != 2 {
		t.Errorf("got %+v after %d attempts, want deleted after 2", res, deletes)
	}
}

func TestDeleteVpcPeering_NotAuthenticated(t *testing.T) {
	srv, _ := newPeeringServer(t, http.StatusUnauthorized,
		peeringErrorBody("UNAUTHORIZED", "missing token"), http.StatusOK)
	client := newTestClient(t, srv.URL, "")

	_, err := DeleteVpcPeering(context.Background(), client, testFabric, testName, testVPC, testPeer)
	if !errors.Is(err, ErrNotAuthenticated) {
		t.Fatalf("expected ErrNotAuthenticated, got %T (%v)", err, err)
	}
}

func TestDeleteVpcPeering_NilClient(t *testing.T) {
	_, err := DeleteVpcPeering(context.Background(), nil, testFabric, testName, testVPC, testPeer)
	if !errors.Is(err, ErrNilClient) {
		t.Fatalf("expected ErrNilClient, got %v", err)
	}
}

// apiErrorOf returns the embedded APIError of any typed SDK error.
func apiErrorOf(err error) *APIError {
	var (
		bad  *BadRequestError
		auth *AuthenticationError
		nf   *NotFoundError
		cf   *ConflictError
		srv  *ServerError
		api  *APIError
	)
	switch {
	case errors.As(err, &bad):
		return &bad.APIError
	case errors.As(err, &auth):
		return &auth.APIError
	case errors.As(err, &nf):
		return &nf.APIError
	case errors.As(err, &cf):
		return &cf.APIError
	case errors.As(err, &srv):
		return &srv.APIError
	case errors.As(err, &api):
		return api
	}
	return nil
}
