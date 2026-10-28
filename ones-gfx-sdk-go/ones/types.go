package ones

// types.go re-exports the types, constants and options a caller needs, so
// `ones` is the only package you have to import:
//
//	client, err := ones.InitializeWithCreds("https://host:3002", "user", "pass")
//	fabrics, err := ones.GetFabrics(ctx, client)
//	tenants, err := ones.ListTenants(ctx, client, "f1")
//
// These are Go type aliases, so ones.GpuServerInfo and
// resources.GpuServerInfo are the same type and interchangeable.
//
// Generated — edit the generator, not this file.

import (
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx/resources"
)

// Ptr returns a pointer to v — a convenience for the optional (pointer)
// arguments used throughout the API.
func Ptr[T any](v T) *T { return &v }

// --- core types -------------------------------------------------------------

type (
	APIError             = ones_gfx.APIError
	AuthData             = ones_gfx.AuthData
	AuthProvider         = ones_gfx.AuthProvider
	AuthResponse         = ones_gfx.AuthResponse
	AuthenticationError  = ones_gfx.AuthenticationError
	BadRequestError      = ones_gfx.BadRequestError
	CallOption           = ones_gfx.CallOption
	ClientOption         = ones_gfx.ClientOption
	ConfigStatus         = ones_gfx.ConfigStatus
	Tenant               = ones_gfx.Tenant
	VNetInfo             = ones_gfx.VNetInfo
	VNetTenantIPPool     = ones_gfx.VNetTenantIPPool
	NMXCPartition        = ones_gfx.NMXCPartition
	ConflictError        = ones_gfx.ConflictError
	GPUOperation         = ones_gfx.GPUOperation
	GpuAction            = ones_gfx.GpuAction
	JWTAuth              = ones_gfx.JWTAuth
	NotFoundError        = ones_gfx.NotFoundError
	ONESError            = ones_gfx.ONESError
	Operation            = ones_gfx.Operation
	OperationAccepted    = ones_gfx.OperationAccepted
	OperationFailedError = ones_gfx.OperationFailedError
	OperationMode        = ones_gfx.OperationMode
	OperationStatus      = ones_gfx.OperationStatus
	ServerError          = ones_gfx.ServerError
	TokenAuth            = ones_gfx.TokenAuth
	Transport            = ones_gfx.Transport
	TransportError       = ones_gfx.TransportError
	Timestamp            = ones_gfx.Timestamp
	ValidationError      = ones_gfx.ValidationError
)

// --- request/response types -------------------------------------------------

type (
	ApiResponseMessage      = resources.ApiResponseMessage
	AvailableServerResult   = resources.AvailableServerResult
	ConfigScope             = resources.ConfigScope
	FabricsListResponse     = resources.FabricsListResponse
	GpuPortAssignmentResult = resources.GpuPortAssignmentResult
	GpuServerInfo           = resources.GpuServerInfo
	NmxcDomain              = resources.NmxcDomain
	NmxcDomainsUpdateResult = resources.NmxcDomainsUpdateResult
	NmxcFactoryResetResult  = resources.NmxcFactoryResetResult
	NmxcProbeResult         = resources.NmxcProbeResult
	NmxcResetResult         = resources.NmxcResetResult
	OperationStatusItem     = resources.OperationStatusItem
	SuidMap                 = resources.SuidMap
	UfmCredsResult          = resources.UfmCredsResult

	// Typed inventory responses.
	InventoryPortsResult   = resources.InventoryPortsResult
	InventoryHost          = resources.InventoryHost
	InventoryPort          = resources.InventoryPort
	UFMHostsResult         = resources.UFMHostsResult
	UFMHost                = resources.UFMHost
	InventorySyncResult    = resources.InventorySyncResult
	NmxcInventory          = resources.NmxcInventory
	NmxcInventorySummary   = resources.NmxcInventorySummary
	NmxcInventoryDomain    = resources.NmxcInventoryDomain
	NmxcInventoryPartition = resources.NmxcInventoryPartition
	NmxcInventoryServer    = resources.NmxcInventoryServer
	NmxcInventoryGPU       = resources.NmxcInventoryGPU

	// Request structs for the ...With methods.
	CreateTenantRequest        = resources.CreateTenantRequest
	UpdateTenantServersRequest = resources.UpdateTenantServersRequest
	ModifyAllocationsRequest   = resources.ModifyAllocationsRequest
	AssignPortsRequest         = resources.AssignPortsRequest
)

// --- enum constants ---------------------------------------------------------

const (
	GpuActionAdd              = ones_gfx.GpuActionAdd
	GpuActionDelete           = ones_gfx.GpuActionDelete
	OperationAdd              = ones_gfx.OperationAdd
	OperationDelete           = ones_gfx.OperationDelete
	OperationStatusPending    = ones_gfx.OperationStatusPending
	OperationStatusRunning    = ones_gfx.OperationStatusRunning
	OperationStatusSuccess    = ones_gfx.OperationStatusSuccess
	OperationStatusFailure    = ones_gfx.OperationStatusFailure
	OperationModeSynchronous  = ones_gfx.OperationModeSynchronous
	OperationModeAsyncPoll    = ones_gfx.OperationModeAsyncPoll
	OperationModeAsyncWebhook = ones_gfx.OperationModeAsyncWebhook
	ConfigScopeWholeServer    = resources.ConfigScopeWholeServer
	ConfigScopeParticularGPU  = resources.ConfigScopeParticularGPU

	// UnlimitedGPUs is the maxGpusAllowed value meaning "no quota".
	UnlimitedGPUs = ones_gfx.UnlimitedGPUs
)

// Enum parsers turn external input (CLI flags, config files) into checked
// values. Each is case-insensitive and returns an error for unknown input.
var (
	ParseGPUOperation    = ones_gfx.ParseGPUOperation
	ParseOperationStatus = ones_gfx.ParseOperationStatus
	ParseConfigScope     = resources.ParseConfigScope
	ParseConfigStatus    = ones_gfx.ParseConfigStatus
)

// --- errors -----------------------------------------------------------------

// --- options and helpers ----------------------------------------------------

var (
	// WithClientTimeout sets the default per-request timeout.
	WithClientTimeout = ones_gfx.WithClientTimeout

	// WithTLSVerify enables or disables TLS certificate verification.
	WithTLSVerify = ones_gfx.WithTLSVerify

	// WithTLSConfig supplies a custom *tls.Config (e.g. a CA bundle).
	WithTLSConfig = ones_gfx.WithTLSConfig

	// WithTimeout overrides the timeout for a single call.
	WithTimeout = ones_gfx.WithTimeout

	// WithWebhook requests webhook delivery for an async call.
	WithWebhook = ones_gfx.WithWebhook

	// WithIdempotencyKey replays a prior async operation instead of duplicating it.
	WithIdempotencyKey = ones_gfx.WithIdempotencyKey

	// WithRequestOrigin sets the x-request-origin header.
	WithRequestOrigin = ones_gfx.WithRequestOrigin

	// NewJWTAuth builds bearer-token auth for servers using the JWT pair flow.
	NewJWTAuth = ones_gfx.NewJWTAuth

	// FMBaseURL turns a root base URL into the /api/fm/ resource base.
	FMBaseURL = ones_gfx.FMBaseURL
)

// --- network policy, webhooks, IPAM, DHCP and rails -------------------------

type (
	NetworkPolicy       = resources.NetworkPolicy
	NetworkPolicyList   = resources.NetworkPolicyList
	NetworkPolicyName   = resources.NetworkPolicyName
	NetworkPolicyRef    = resources.NetworkPolicyRef
	NetworkPolicySide   = resources.NetworkPolicySide
	NetworkPolicySpec   = resources.NetworkPolicySpec
	WebhookFilters      = resources.WebhookFilters
	WebhookRegistration = resources.WebhookRegistration
	IpamAllocation      = resources.IpamAllocation
	IpamConfig          = resources.IpamConfig
	IpamDhcpConfig      = resources.IpamDhcpConfig
	IpamPagination      = resources.IpamPagination
	IpamPoolSummary     = resources.IpamPoolSummary
	IpamStaticConfig    = resources.IpamStaticConfig
	ServerIPList        = resources.ServerIPList
	TenantIPList        = resources.TenantIPList
	DhcpReapplySummary  = resources.DhcpReapplySummary
	TenantDhcp          = resources.TenantDhcp
	RailCidrPools       = resources.RailCidrPools
	VFInterfaceList     = resources.VFInterfaceList
	VFDpuEntry          = resources.VFDpuEntry
	VFInterface         = resources.VFInterface
	VFAssignment        = resources.VFAssignment
)

// Network policy endpoint types and the longest name FM accepts.
const (
	NetworkPolicyTypeTenant    = resources.NetworkPolicyTypeTenant
	NetworkPolicyTypeService   = resources.NetworkPolicyTypeService
	MaxNetworkPolicyNameLength = resources.MaxNetworkPolicyNameLength
)

// Canonical webhook event types. FM also accepts the fabric.tenant_* aliases
// and resolves them to these on write.
const (
	WebhookEventTenantCreate     = resources.WebhookEventTenantCreate
	WebhookEventTenantDelete     = resources.WebhookEventTenantDelete
	WebhookEventTenantAllocate   = resources.WebhookEventTenantAllocate
	WebhookEventTenantDeallocate = resources.WebhookEventTenantDeallocate
	WebhookEventTenantVPCPeering = resources.WebhookEventTenantVPCPeering
	WebhookEventFabricDeploy     = resources.WebhookEventFabricDeploy
	WebhookEventFabricDelete     = resources.WebhookEventFabricDelete
)

// IP assignment states, used by the status filter on TenantIPs.
const (
	IPAssignStatusAllocatedOnly = resources.IPAssignStatusAllocatedOnly
	IPAssignStatusAssigned      = resources.IPAssignStatusAssigned
	IPAssignStatusFailed        = resources.IPAssignStatusFailed
)

// Media types a rail CIDRPool manifest can be rendered as, and the page and
// hostname limits FM enforces.
const (
	RailYamlMediaType       = resources.RailYamlMediaType
	RailLegacyYamlMediaType = resources.RailLegacyYamlMediaType
	RailTextMediaType       = resources.RailTextMediaType
	MaxDhcpLeasePageSize    = resources.MaxDhcpLeasePageSize
)

// Tenant configuration states, as ConfigStatus values.
const (
	ConfigStatusFail       = ones_gfx.ConfigStatusFail
	ConfigStatusPass       = ones_gfx.ConfigStatusPass
	ConfigStatusNotStarted = ones_gfx.ConfigStatusNotStarted
	ConfigStatusInProgress = ones_gfx.ConfigStatusInProgress
	ConfigStatusBrownField = ones_gfx.ConfigStatusBrownField
	ConfigStatusSkipped    = ones_gfx.ConfigStatusSkipped
)

// WebhookEventTypes lists every canonical event name.
var WebhookEventTypes = resources.WebhookEventTypes
