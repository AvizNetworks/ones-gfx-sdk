package ones

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
)

const testOperationID = "TOPOAUDIT_Fab101_abc123"

// newRoutedServer serves distinct responses for the trigger POST and the
// operation-poll GET, keyed by exact path.
func newRoutedServer(t *testing.T, routes map[string]struct {
	status int
	body   interface{}
}) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/user/refresh" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"token": "tok2", "message": "refreshed"},
			})
			return
		}
		route, ok := routes[r.URL.Path]
		if !ok {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		w.WriteHeader(route.status)
		if route.body != nil {
			_ = json.NewEncoder(w).Encode(route.body)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestTriggerTopologyAudit_Success(t *testing.T) {
	srv := newRoutedServer(t, map[string]struct {
		status int
		body   interface{}
	}{
		"/api/fm/fabrics/" + testFabric + "/topology-audit": {
			status: http.StatusAccepted,
			body: map[string]interface{}{
				"status":  "success",
				"message": "Topology audit accepted.",
				"data":    map[string]interface{}{"operationId": testOperationID, "status": "PENDING"},
			},
		},
	})
	client := newTestClient(t, srv.URL, "tok")

	res, err := TriggerTopologyAudit(context.Background(), client, testFabric)
	if err != nil {
		t.Fatalf("TriggerTopologyAudit: %v", err)
	}
	if res.OperationID != testOperationID || res.Status != "PENDING" {
		t.Fatalf("got %+v", res)
	}
}

func TestTriggerTopologyAudit_InvalidFabric(t *testing.T) {
	srv := newRoutedServer(t, map[string]struct {
		status int
		body   interface{}
	}{
		"/api/fm/fabrics/" + testFabric + "/topology-audit": {
			status: http.StatusNotFound,
			body: map[string]interface{}{
				"status": "error",
				"error":  map[string]interface{}{"code": "INVALID_FABRIC", "detail": "no fabric or intent found"},
			},
		},
	})
	client := newTestClient(t, srv.URL, "tok")

	_, err := TriggerTopologyAudit(context.Background(), client, testFabric)
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Fatalf("expected *NotFoundError, got %T (%v)", err, err)
	}
}

func TestGetTopologyAuditStatus_IsOperationsGet(t *testing.T) {
	srv := newRoutedServer(t, map[string]struct {
		status int
		body   interface{}
	}{
		"/api/fm/operations/" + testOperationID: {
			status: http.StatusOK,
			body:   map[string]interface{}{"id": testOperationID, "status": "RUNNING", "progress": 10},
		},
	})
	client := newTestClient(t, srv.URL, "tok")

	status, err := GetTopologyAuditStatus(context.Background(), client, testOperationID)
	if err != nil {
		t.Fatalf("GetTopologyAuditStatus: %v", err)
	}
	if status.Status != "RUNNING" {
		t.Fatalf("Status = %q, want RUNNING", status.Status)
	}
}

func opServer(t *testing.T, status string, extra map[string]interface{}) *httptest.Server {
	t.Helper()
	body := map[string]interface{}{"id": testOperationID, "status": status}
	for k, v := range extra {
		body[k] = v
	}
	return newRoutedServer(t, map[string]struct {
		status int
		body   interface{}
	}{
		"/api/fm/operations/" + testOperationID: {status: http.StatusOK, body: body},
	})
}

func TestGetTopologyAuditResult_NotDone(t *testing.T) {
	for _, status := range []string{"PENDING", "RUNNING"} {
		srv := opServer(t, status, nil)
		client := newTestClient(t, srv.URL, "tok")

		_, err := GetTopologyAuditResult(context.Background(), client, testOperationID)
		if !errors.Is(err, ErrTopologyAuditNotDone) {
			t.Fatalf("status=%s: expected ErrTopologyAuditNotDone, got %v", status, err)
		}
	}
}

func TestGetTopologyAuditResult_Failed(t *testing.T) {
	srv := opServer(t, "FAILED", map[string]interface{}{"errorMessage": "No intent found for fabric 'Fab101'."})
	client := newTestClient(t, srv.URL, "tok")

	_, err := GetTopologyAuditResult(context.Background(), client, testOperationID)
	if err == nil || errors.Is(err, ErrTopologyAuditNotDone) {
		t.Fatalf("expected a failure error, got %v", err)
	}
}

func TestGetTopologyAuditResult_Rollups(t *testing.T) {
	cases := []struct {
		name          string
		overallStatus string
		links         []resources.TopologyAuditLink
	}{
		{
			name:          "pass",
			overallStatus: "PASS",
			links: []resources.TopologyAuditLink{
				{Device: "spine1", LocalInterface: "Ethernet1", RemoteDevice: "leaf1", RemoteInterface: "Ethernet1",
					ObservedRemoteDevice: "leaf1", ObservedRemoteInterface: "Ethernet1", Status: "PASS", Message: "OK"},
			},
		},
		{
			name:          "fail",
			overallStatus: "FAIL",
			links: []resources.TopologyAuditLink{
				{Device: "spine1", LocalInterface: "Ethernet3", RemoteDevice: "leaf2", RemoteInterface: "Ethernet1",
					Status: "MISSING_LINK", Message: "No LLDP neighbor observed on spine1:Ethernet3; expected leaf2:Ethernet1"},
			},
		},
		{
			name:          "incomplete",
			overallStatus: "INCOMPLETE",
			links: []resources.TopologyAuditLink{
				{Device: "spine1", LocalInterface: "Ethernet3", Status: "UNVERIFIABLE", Message: "device unreachable"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			report := resources.TopologyAuditReport{
				FabricName:    testFabric,
				IntentVersion: "42:2026-11-10T00:00:00Z",
				OverallStatus: tc.overallStatus,
				Links:         tc.links,
				DeviceIPMap:   map[string]string{"spine1": "10.0.0.1"},
			}
			resultJSON, err := json.Marshal(report)
			if err != nil {
				t.Fatalf("marshal report: %v", err)
			}
			srv := opServer(t, "SUCCESS", map[string]interface{}{"result": string(resultJSON)})
			client := newTestClient(t, srv.URL, "tok")

			got, err := GetTopologyAuditResult(context.Background(), client, testOperationID)
			if err != nil {
				t.Fatalf("GetTopologyAuditResult: %v", err)
			}
			if got.OverallStatus != tc.overallStatus {
				t.Fatalf("OverallStatus = %q, want %q", got.OverallStatus, tc.overallStatus)
			}
			if len(got.Links) != len(tc.links) || got.Links[0].Status != tc.links[0].Status {
				t.Fatalf("Links = %+v, want %+v", got.Links, tc.links)
			}
		})
	}
}

func TestGetTopologyAuditResult_MalformedResult(t *testing.T) {
	srv := opServer(t, "SUCCESS", map[string]interface{}{"result": "{not valid json"})
	client := newTestClient(t, srv.URL, "tok")

	_, err := GetTopologyAuditResult(context.Background(), client, testOperationID)
	if err == nil {
		t.Fatal("expected an error for malformed result JSON")
	}
}

func TestTopologyAuditClientField_NilClient(t *testing.T) {
	if _, err := TriggerTopologyAudit(context.Background(), nil, testFabric); !errors.Is(err, ErrNilClient) {
		t.Fatalf("expected ErrNilClient, got %v", err)
	}
	if _, err := GetTopologyAuditStatus(context.Background(), nil, testOperationID); !errors.Is(err, ErrNilClient) {
		t.Fatalf("expected ErrNilClient, got %v", err)
	}
	if _, err := GetTopologyAuditResult(context.Background(), nil, testOperationID); !errors.Is(err, ErrNilClient) {
		t.Fatalf("expected ErrNilClient, got %v", err)
	}
}
