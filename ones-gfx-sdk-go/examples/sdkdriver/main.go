// sdkdriver exposes the Go SDK over a JSON-lines protocol on stdin/stdout so
// that ones-gfx-sdk-python/examples/sdk_test_runner.py can run the same test
// steps against both SDKs.
//
// Each request is one line: {"op": "<name>", "args": {...}}
// Each response is one line:
//
//	{"ok": true,  "data": <result>}
//	{"ok": false, "status": <http status or 0>, "body": <server body>, "error": "<message>"}
//
// One driver process holds one logged-in client, so call "login" first.
// Not intended for direct use; build and run it through the test runner.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

type request struct {
	Op   string `json:"op"`
	Args args   `json:"args"`
}

type args struct {
	URL         string                         `json:"url"`
	Username    string                         `json:"username"`
	Password    string                         `json:"password"`
	VerifyTLS   bool                           `json:"verify_tls"`
	Name        string                         `json:"name"`
	Fabric      string                         `json:"fabric"`
	Tenant      string                         `json:"tenant"`
	Description string                         `json:"description"`
	Type        string                         `json:"type"`
	Status      string                         `json:"status"`
	MaxGpus     int                            `json:"max_gpus"`
	Servers     []string                       `json:"servers"`
	Shared      bool                           `json:"shared"`
	Operation   string                         `json:"operation"`
	GpuIDs      []int                          `json:"gpu_ids"`
	Suid        map[string]map[string][]string `json:"suid"`
	ID          string                         `json:"id"`
	Async       bool                           `json:"async"`
	WebhookURL  string                         `json:"webhook_url"`
	Event       string                         `json:"event"`
}

type response struct {
	OK     bool        `json:"ok"`
	Data   interface{} `json:"data,omitempty"`
	Status int         `json:"status,omitempty"`
	Body   interface{} `json:"body,omitempty"`
	Error  string      `json:"error,omitempty"`
}

var client *ones.Client

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 1<<20), 1<<20)
	out := json.NewEncoder(os.Stdout)
	for in.Scan() {
		var req request
		if err := json.Unmarshal(in.Bytes(), &req); err != nil {
			_ = out.Encode(response{Error: "bad request line: " + err.Error()})
			continue
		}
		data, err := dispatch(context.Background(), req.Op, req.Args)
		if err != nil {
			_ = out.Encode(failure(err))
			continue
		}
		_ = out.Encode(response{OK: true, Data: data})
	}
	if client != nil {
		_ = client.Close()
	}
}

// failure flattens the SDK's typed errors into status + server body.
func failure(err error) response {
	var (
		bad      *ones.BadRequestError
		auth     *ones.AuthenticationError
		notFound *ones.NotFoundError
		conflict *ones.ConflictError
		server   *ones.ServerError
		api      *ones.APIError
	)
	var e *ones.APIError
	switch {
	case errors.As(err, &bad):
		e = &bad.APIError
	case errors.As(err, &auth):
		e = &auth.APIError
	case errors.As(err, &notFound):
		e = &notFound.APIError
	case errors.As(err, &conflict):
		e = &conflict.APIError
	case errors.As(err, &server):
		e = &server.APIError
	case errors.As(err, &api):
		e = api
	}
	if e == nil {
		return response{Error: err.Error()}
	}
	return response{Status: e.StatusCode, Body: e.ResponseBody, Error: err.Error()}
}

func gpuAction(op string) ones.GpuAction {
	if op == "DELETE" {
		return ones.GpuActionDelete
	}
	return ones.GpuActionAdd
}

func callOpts(a args) []ones.CallOption {
	if a.WebhookURL == "" {
		return nil
	}
	return []ones.CallOption{ones.WithWebhook(a.WebhookURL, []string{a.Event})}
}

func dispatch(ctx context.Context, op string, a args) (interface{}, error) {
	if op == "login" {
		c, err := ones.InitializeWithCreds(a.URL, a.Username, a.Password, ones.WithTLSVerify(a.VerifyTLS))
		if err != nil {
			return nil, err
		}
		client = c
		return fmt.Sprintf("logged in, token %d chars", len(c.GetAuthToken())), nil
	}
	if op == "no_token" {
		c, err := ones.NewClient(a.URL, ones.WithTLSVerify(a.VerifyTLS))
		if err != nil {
			return nil, err
		}
		defer c.Close()
		_, err = ones.GetAllFabrics(ctx, c)
		if errors.Is(err, ones.ErrNotAuthenticated) {
			return "ErrNotAuthenticated returned before any request", nil
		}
		return nil, fmt.Errorf("token-less call was not rejected (err=%v)", err)
	}
	if client == nil {
		return nil, errors.New("not logged in: send the login op first")
	}

	switch op {
	case "version":
		return ones.GetControllerVersion(ctx, client)
	case "refresh":
		if _, err := client.RefreshAuth(); err != nil {
			return nil, err
		}
		return "token rotated", nil
	case "logout":
		if _, err := client.Logout(); err != nil {
			return nil, err
		}
		return "logged out", nil

	case "list_fabrics":
		return ones.GetAllFabrics(ctx, client)
	case "get_fabric":
		return ones.GetFabricByName(ctx, client, a.Name)
	case "create_fabric":
		return ones.CreateFabric(ctx, client, a.Name, &ones.FabricCreateArgs{
			Type: ones.Ptr(a.Type), Status: ones.Ptr(a.Status), Description: ones.Ptr(a.Description),
		})
	case "delete_fabric":
		return ones.DeleteFabric(ctx, client, a.Name)
	case "inventory_sync":
		return ones.InventorySync(ctx, client, a.Fabric)

	case "list_tenants":
		return ones.ListTenants(ctx, client, a.Fabric)
	case "get_tenant":
		return ones.GetTenant(ctx, client, a.Fabric, a.Tenant)
	case "available_servers":
		return ones.GetAvailableServers(ctx, client, a.Fabric)
	case "gpus":
		return ones.GetAllGpusList(ctx, client, a.Fabric)
	case "get_operation":
		return ones.GetOperation(ctx, client, a.ID)

	case "create_tenant":
		shared := false
		if a.Async {
			return ones.CreateTenantAsync(ctx, client, a.Fabric, a.Tenant, &a.Description, &a.MaxGpus, &shared, callOpts(a)...)
		}
		return ones.CreateTenant(ctx, client, a.Fabric, a.Tenant, &a.Description, &a.MaxGpus, &shared)
	case "update_tenant":
		specs := make([]ones.GpuServerInfo, 0, len(a.Servers))
		for _, s := range a.Servers {
			var shared *bool
			if a.Shared {
				shared = ones.Ptr(true)
			}
			specs = append(specs, ones.GpuServerInfo{ServerName: s, Shared: shared})
		}
		action := gpuAction(a.Operation)
		if a.Async {
			return ones.UpdateTenantAsync(ctx, client, a.Fabric, a.Tenant, specs, &action, callOpts(a)...)
		}
		return ones.UpdateTenant(ctx, client, a.Fabric, a.Tenant, specs, &action)
	case "delete_tenant":
		if a.Async {
			return ones.DeleteTenantAsync(ctx, client, a.Fabric, a.Tenant, callOpts(a)...)
		}
		return ones.DeleteTenant(ctx, client, a.Fabric, a.Tenant)
	case "assign_ports":
		return ones.AssignGpuPorts(ctx, client, a.Fabric, a.Tenant, gpuAction(a.Operation), a.Servers, a.GpuIDs, nil)
	case "gpu_allocations":
		suid := ones.SuidMap{}
		for idx, hosts := range a.Suid {
			suid[idx] = map[string]ones.DeviceGpus{}
			for host, gpus := range hosts {
				suid[idx][host] = ones.DeviceGpus{Gpus: gpus}
			}
		}
		action := gpuAction(a.Operation)
		if a.Async {
			return ones.ModifyGpuAllocationsAsync(ctx, client, a.Fabric, a.Tenant, suid, &action, nil, nil, callOpts(a)...)
		}
		return ones.ModifyGpuAllocations(ctx, client, a.Fabric, a.Tenant, suid, &action, nil, nil)
	}
	return nil, fmt.Errorf("unknown op %q", op)
}
