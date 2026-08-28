package ones_gfx

// GPUOperation is the ADD/DELETE discriminator for GPU configuration
// (Helper/Enums.java GpuAction). GpuAction is an alias for it.
type GPUOperation string

const (
	// OperationAdd maps GPUs to a tenant.
	OperationAdd GPUOperation = "ADD"

	// OperationDelete removes a GPU mapping from a tenant.
	OperationDelete GPUOperation = "DELETE"
)
