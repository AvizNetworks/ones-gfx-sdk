package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FabricsResource handles operations on the fabrics collection.
type FabricsResource struct {
	transport *ones_gfx.Transport
}

// NewFabricsResource constructs a FabricsResource.
func NewFabricsResource(transport *ones_gfx.Transport) *FabricsResource {
	return &FabricsResource{transport: transport}
}

// List retrieves all fabrics defined in the ONES system.
// Maps to GET /fabrics.
func (r *FabricsResource) List(ctx context.Context) ([]ones_gfx.Fabric, error) {
	body, err := r.transport.Get("fabrics", nil)
	if err != nil {
		return nil, err
	}

	// Response shape: {"fabrics": [ {...}, {...} ]}
	obj, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", body)
	}

	fabricsRaw, ok := obj["fabrics"]
	if !ok {
		return nil, fmt.Errorf("response missing 'fabrics' field")
	}

	fabricsList, ok := fabricsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'fabrics' field is not an array")
	}

	fabrics := make([]ones_gfx.Fabric, 0, len(fabricsList))
	for _, item := range fabricsList {
		jsonBytes, err := remarshal(item)
		if err != nil {
			continue
		}
		var fabric ones_gfx.Fabric
		if err := unmarshalJSON(jsonBytes, &fabric); err != nil {
			continue
		}
		fabrics = append(fabrics, fabric)
	}

	return fabrics, nil
}

// ModifyGPUAllocations maps or unmaps specific GPUs to a tenant on shared
// fabric servers. It maps to:
//
//	POST /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations
//
// Use req.Operation = ones_gfx.OperationAdd to map GPUs, or
// ones_gfx.OperationDelete to remove the mapping. The Suid field targets one
// or more servers by index and hostname, each with an explicit list of GPU IDs.
//
// Example:
//
//	resp, err := client.Fabrics.ModifyGPUAllocations(ctx, "fabric1", "tenant1",
//	    ones_gfx.GPUAllocationRequest{
//	        Operation: ones_gfx.OperationAdd,
//	        Suid: map[string]map[string]ones_gfx.ServerGPUs{
//	            "0": {"hgx-su00-h00": {GPUs: []string{"G0", "G1", "G2", "G3"}}},
//	        },
//	    })
func (r *FabricsResource) ModifyGPUAllocations(
	ctx context.Context,
	fabricName string,
	tenantName string,
	req ones_gfx.GPUAllocationRequest,
) (*ones_gfx.GPUAllocationResponse, error) {
	if fabricName == "" {
		return nil, fmt.Errorf("fabricName is required")
	}
	if tenantName == "" {
		return nil, fmt.Errorf("tenantName is required")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("fabrics/%s/tenants/%s/gpuAllocations", fabricName, tenantName)
	result, err := r.transport.Post(path, req, ones_gfx.OperationModeSynchronous, nil)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return &ones_gfx.GPUAllocationResponse{Status: "success"}, nil
	}

	if msg, ok := result.(string); ok {
		return &ones_gfx.GPUAllocationResponse{
			Status:  "success",
			Message: msg,
		}, nil
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GPU allocation response: %w", err)
	}

	var resp ones_gfx.GPUAllocationResponse
	if err := json.Unmarshal(jsonBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GPU allocation response: %w", err)
	}

	return &resp, nil
}
