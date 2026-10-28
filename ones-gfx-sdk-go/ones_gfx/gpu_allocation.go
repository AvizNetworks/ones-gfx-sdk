package ones_gfx

import (
	"fmt"
	"strings"
)

// GPUOperation is the ADD/DELETE discriminator for GPU configuration
// (Helper/Enums.java GpuAction). GpuAction is an alias for it.
type GPUOperation string

const (
	// OperationAdd maps GPUs to a tenant.
	OperationAdd GPUOperation = "ADD"

	// OperationDelete removes a GPU mapping from a tenant.
	OperationDelete GPUOperation = "DELETE"
)

// Valid reports whether o is ADD or DELETE.
func (o GPUOperation) Valid() bool {
	return o == OperationAdd || o == OperationDelete
}

// String returns the wire value.
func (o GPUOperation) String() string { return string(o) }

// ParseGPUOperation converts external input (a CLI flag, a config file) into a
// GPUOperation. It is case-insensitive and ignores surrounding space, so "add"
// and " ADD " both parse; anything other than ADD or DELETE is an error.
func ParseGPUOperation(s string) (GPUOperation, error) {
	op := GPUOperation(strings.ToUpper(strings.TrimSpace(s)))
	if !op.Valid() {
		return "", fmt.Errorf("invalid GPU operation %q: want %s or %s", s, OperationAdd, OperationDelete)
	}
	return op, nil
}
