package ones

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// topologyServer serves GET /fabrics/{name} like the ONES gateway and records
// the escaped path and Accept header of the last request.
func topologyServer(t *testing.T, status int, contentType, body string) (*Client, *http.Request) {
	t.Helper()
	last := &http.Request{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*last = *r.Clone(context.Background())
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return newTestClient(t, srv.URL, "tok"), last
}

func readFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("../ones_gfx/testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGetFabricYAMLUsesGatewayRoot(t *testing.T) {
	topo := readFixture(t, "ns_topology.yaml")
	c, last := topologyServer(t, http.StatusOK, "application/x-yaml", topo)

	got, err := GetFabricYAML(context.Background(), c, "full topo check")
	if err != nil {
		t.Fatal(err)
	}
	if got != topo {
		t.Error("body not returned verbatim")
	}
	if last.Method != http.MethodGet || last.URL.EscapedPath() != "/fabrics/full%20topo%20check" {
		t.Errorf("request = %s %s, want GET /fabrics/full%%20topo%%20check (root, not /api/fm)", last.Method, last.URL.EscapedPath())
	}
	if last.Header.Get("authorization") != "tok" || last.Header.Get("Accept") != "application/x-yaml, application/json" {
		t.Errorf("headers = %v", last.Header)
	}
}

func TestGetFabricNSInterfacesEndToEnd(t *testing.T) {
	c, _ := topologyServer(t, http.StatusOK, "application/x-yaml", readFixture(t, "ns_topology.yaml"))
	n, err := GetFabricNSInterfaces(context.Background(), c, "EW", NSInterfaceOptions{})
	if err != nil {
		t.Fatal(err)
	}
	y, _ := n.YAML()
	if string(y) != readFixture(t, "ns_topology_expected.yaml") {
		t.Errorf("YAML differs from golden:\n%s", y)
	}
}

func TestGetFabricYAMLErrors(t *testing.T) {
	ctx := context.Background()

	c, _ := topologyServer(t, http.StatusNotFound, "application/json",
		`{"status":"FAILED","code":404,"errorType":"NOT_FOUND","message":"YAML file for fabric \"nope\" not found."}`)
	_, err := GetFabricYAML(ctx, c, "nope")
	var nf *NotFoundError
	if !errors.As(err, &nf) || nf.Message != `YAML file for fabric "nope" not found.` {
		t.Errorf("404: got %T (%v)", err, err)
	}

	c, _ = topologyServer(t, http.StatusOK, "text/html", "<!doctype html><html></html>")
	_, err = GetFabricNSInterfaces(ctx, c, "EW", NSInterfaceOptions{})
	var te *TopologyError
	if !errors.As(err, &te) {
		t.Errorf("HTML fallback: expected *TopologyError, got %T (%v)", err, err)
	}

	_, err = GetFabricYAML(ctx, c, "a/b")
	var ve *ValidationError
	if !errors.As(err, &ve) || ve.Field != "fabricName" {
		t.Errorf("bad name: got %T (%v)", err, err)
	}
	if _, err := GetFabricYAML(ctx, nil, "EW"); !errors.Is(err, ErrNilClient) {
		t.Errorf("nil client: %v", err)
	}
}
