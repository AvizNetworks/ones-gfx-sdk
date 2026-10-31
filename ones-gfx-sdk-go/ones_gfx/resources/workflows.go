package resources

import (
	"context"
	"errors"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// WorkflowsResource holds the multi-call sequences. Every other resource in
// this SDK is one method per endpoint; these are the sequences whose *order*
// is forced by the Fabric Manager, so getting them wrong strands state
// rather than returning an error.
type WorkflowsResource struct {
	transport *ones_gfx.Transport
	tenants   *TenantsResource
	gpu       *GPUResource
	vfs       *VFInterfacesResource
	shared    *SharedServicesResource
}

// NewWorkflowsResource constructs a WorkflowsResource.
func NewWorkflowsResource(transport *ones_gfx.Transport) *WorkflowsResource {
	return &WorkflowsResource{
		transport: transport,
		tenants:   NewTenantsResource(transport),
		gpu:       NewGPUResource(transport),
		vfs:       NewVFInterfacesResource(transport),
		shared:    NewSharedServicesResource(transport),
	}
}

// CleanUpResult is what CleanUpTenant did, and what it could not finish.
// Clean is true when Failures is empty; otherwise the named steps can be
// retried by calling CleanUpTenant again.
type CleanUpResult struct {
	FabricName     string   `json:"fabricName"`
	TenantName     string   `json:"tenantName"`
	Servers        []string `json:"servers"`
	VFsFreed       []string `json:"vfsFreed"`
	SharedServices bool     `json:"sharedServices"`
	Failures       []string `json:"failures"`
	Clean          bool     `json:"clean"`
}

// isNotFound reports whether err is a 404, which every teardown step treats
// as "already gone" rather than a failure.
func isNotFound(err error) bool {
	var nf *ones_gfx.NotFoundError
	return errors.As(err, &nf)
}

// CleanUpTenant tears a DPU tenant down in the one order the FM accepts.
//
// DELETE /tenants/{t} answers 409 GPUS_ALLOCATED while GPUs are still
// attached, and a VF cannot be freed once its server has left the tenant.
// So the order is forced, and is the reverse of provisioning:
//
//  1. unassign every VF the tenant holds
//  2. detach every GPU server from the tenant
//  3. delete the tenant
//
// A shared-services tenant is deleted through its own route instead, which
// also releases its Kea subnet and IPAM pool.
//
// Every step tolerates "already gone", so re-running a half-finished
// teardown finishes it rather than failing. Steps that genuinely fail are
// collected in Failures instead of returning early, so one bad step never
// strands the rest. Pass servers to skip discovery.
func (r *WorkflowsResource) CleanUpTenant(ctx context.Context, fabricName, tenantName string, servers []string, opts ...ones_gfx.CallOption) (*CleanUpResult, error) {
	if err := fabricTenant(fabricName, tenantName); err != nil {
		return nil, err
	}
	result := &CleanUpResult{
		FabricName: fabricName,
		TenantName: tenantName,
		VFsFreed:   []string{},
		Failures:   []string{},
	}

	mappings, err := r.gpu.TenantMappings(ctx, fabricName, &tenantName)
	if err != nil && !isNotFound(err) {
		mappings = nil
	}

	found := servers
	if found == nil {
		seen := map[string]bool{}
		for _, m := range mappings {
			if m.ServerName != "" && !seen[m.ServerName] {
				seen[m.ServerName] = true
				found = append(found, m.ServerName)
			}
		}
	}
	result.Servers = found

	// 1. VFs ----------------------------------------------------------------
	for _, server := range found {
		vfs, err := r.vfs.List(ctx, fabricName, server, opts...)
		if err != nil {
			continue // nothing listed means nothing to free here
		}
		for _, vf := range vfs.Flatten() {
			if vf.TenantName != tenantName {
				continue
			}
			id := vf.IfName
			if id == "" {
				id = vf.ServerIf
			}
			if id == "" {
				continue
			}
			if _, err := r.vfs.Unassign(ctx, fabricName, server, id, opts...); err != nil {
				if !isNotFound(err) {
					result.Failures = append(result.Failures, "vf "+server+"/"+id)
				}
				continue
			}
			result.VFsFreed = append(result.VFsFreed, server+"/"+id)
		}
	}

	// 2. GPU servers --------------------------------------------------------
	// Addressed as fabric -> SU -> server -> port names, so the ports come
	// from the fabric's own GPU inventory, grouped the same way.
	if len(found) > 0 {
		inFound := map[string]bool{}
		for _, s := range found {
			inFound[s] = true
		}
		suid := SuidMap{}
		for _, m := range mappings {
			if !inFound[m.ServerName] || m.LogicalGpuName == "" {
				continue
			}
			su := "0"
			if _, ok := suid[su]; !ok {
				suid[su] = map[string]DeviceGpus{}
			}
			dg := suid[su][m.ServerName]
			dg.Gpus = append(dg.Gpus, m.LogicalGpuName)
			suid[su][m.ServerName] = dg
		}
		if len(suid) > 0 {
			action := ones_gfx.GpuAction("DELETE")
			scope := ConfigScopeWholeServer
			if _, err := r.tenants.ModifyAllocations(ctx, fabricName, tenantName, suid, &action, &scope, nil, opts...); err != nil {
				if !isNotFound(err) {
					result.Failures = append(result.Failures, "gpu detach")
				}
			}
		}
	}

	// 3. Tenant -------------------------------------------------------------
	if ss, err := r.shared.Get(ctx, fabricName, opts...); err == nil && ss != nil && ss.TenantName == tenantName {
		result.SharedServices = true
	}
	if result.SharedServices {
		if _, err := r.shared.Delete(ctx, fabricName, opts...); err != nil && !isNotFound(err) {
			result.Failures = append(result.Failures, "shared-services delete")
		}
	} else {
		if _, err := r.tenants.Delete(ctx, fabricName, tenantName, opts...); err != nil && !isNotFound(err) {
			result.Failures = append(result.Failures, "tenant delete")
		}
	}

	result.Clean = len(result.Failures) == 0
	return result, nil
}
