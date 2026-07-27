package ones_gfx

import "fmt"

// GPUOperation represents the type of GPU allocation operation.
type GPUOperation string

const (
	// OperationAdd maps GPUs to a tenant on a shared fabric server.
	OperationAdd GPUOperation = "ADD"

	// OperationDelete removes a GPU mapping from a tenant.
	OperationDelete GPUOperation = "DELETE"
)

// ServerGPUs holds the list of GPU identifiers on a single compute node.
//
// Example:
//
//	ServerGPUs{GPUs: []string{"G0", "G1", "G2", "G3"}}
type ServerGPUs struct {
	// GPUs is the list of GPU identifiers (e.g. "G0", "G1") on this node.
	GPUs []string `json:"gpus"`
}

// GPUAllocationRequest is the request body for
// POST /fabrics/{fabricName}/tenants/{tenantName}/gpuAllocations.
//
// Suid maps a server index (e.g. "0", "1") to a map of hostname → ServerGPUs,
// allowing multiple servers to be targeted in a single request.
//
// Example:
//
//	GPUAllocationRequest{
//	    Operation: OperationAdd,
//	    Suid: map[string]map[string]ServerGPUs{
//	        "0": {
//	            "hgx-su00-h00": {GPUs: []string{"G0", "G1", "G2", "G3"}},
//	        },
//	    },
//	}
type GPUAllocationRequest struct {
	// Operation must be OperationAdd or OperationDelete.
	Operation GPUOperation `json:"operation"`

	// Suid maps server index → hostname → ServerGPUs.
	Suid map[string]map[string]ServerGPUs `json:"suid"`
}

// Validate returns an error if the request is missing required fields or
// contains an invalid operation value.
func (r GPUAllocationRequest) Validate() error {
	if r.Operation != OperationAdd && r.Operation != OperationDelete {
		return fmt.Errorf("operation must be %q or %q, got %q", OperationAdd, OperationDelete, r.Operation)
	}
	if len(r.Suid) == 0 {
		return fmt.Errorf("suid must contain at least one server entry")
	}
	for idx, hosts := range r.Suid {
		if len(hosts) == 0 {
			return fmt.Errorf("suid[%s] must contain at least one hostname entry", idx)
		}
		for hostname, srv := range hosts {
			if len(srv.GPUs) == 0 {
				return fmt.Errorf("suid[%s][%s].gpus must not be empty", idx, hostname)
			}
		}
	}
	return nil
}

// GPUAllocationResponse is returned by ModifyGPUAllocations.
// When the server processes the request asynchronously, OperationID will be
// populated and can be passed to Operations.Get to poll for completion.
type GPUAllocationResponse struct {
	// OperationID is the async operation identifier, if applicable.
	OperationID string `json:"operationId,omitempty"`

	// Status is the high-level result string returned by the server.
	Status string `json:"status,omitempty"`

	// Message is an optional human-readable message from the server.
	Message string `json:"message,omitempty"`
}
