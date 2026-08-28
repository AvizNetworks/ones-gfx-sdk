package ones

// types.go re-exports the types, constants and options a caller needs, so
// `ones` is the only package you have to import:
//
//	client, err := ones.InitializeWithCreds("https://host:3002", "user", "pass")
//	fabrics, err := ones.GetAllFabrics(ctx, client)
//	msg, err := ones.CreateFabric(ctx, client, "f1", &ones.FabricCreateArgs{Type: ones.Ptr("DNO ASN")})
//
// These are Go type aliases, so ones.FabricCreateArgs and
// resources.FabricCreateArgs are the same type and interchangeable.
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
	ConflictError        = ones_gfx.ConflictError
	FabricItem           = ones_gfx.FabricItem
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
)

// --- request/response types -------------------------------------------------

type (
	ApiResponseMessage        = resources.ApiResponseMessage
	AvailableServerResult     = resources.AvailableServerResult
	BootstrapStage            = resources.BootstrapStage
	BootstrapinfoItem         = resources.BootstrapinfoItem
	BootstrapinfoRecord       = resources.BootstrapinfoRecord
	ConfigScope               = resources.ConfigScope
	ControllerVersion         = resources.ControllerVersion
	DeleteFileResult          = resources.DeleteFileResult
	DeviceConfigBackup        = resources.DeviceConfigBackup
	DeviceConfigRestore       = resources.DeviceConfigRestore
	DeviceDetail              = resources.DeviceDetail
	DeviceInventoryItem       = resources.DeviceInventoryItem
	FabricCreateArgs          = resources.FabricCreateArgs
	FabricInventoryUpdateItem = resources.FabricInventoryUpdateItem
	FabricSimItem             = resources.FabricSimItem
	FabricsListResponse       = resources.FabricsListResponse
	FetchBackupFilesResult    = resources.FetchBackupFilesResult
	GetFilesResult            = resources.GetFilesResult
	GpuAllocationHistoryItem  = resources.GpuAllocationHistoryItem
	GpuItem                   = resources.GpuItem
	GpuPortAssignmentResult   = resources.GpuPortAssignmentResult
	GpuServerInfo             = resources.GpuServerInfo
	GpuTenantMappingItem      = resources.GpuTenantMappingItem
	HostAction                = resources.HostAction
	HostItem                  = resources.HostItem
	HosttenantsRecord         = resources.HosttenantsRecord
	ImageUpgradeDetailsItem   = resources.ImageUpgradeDetailsItem
	IntentItem                = resources.IntentItem
	InventoryItem             = resources.InventoryItem
	InventoryRecord           = resources.InventoryRecord
	NetOpsAction              = resources.NetOpsAction
	NetOpsParams              = resources.NetOpsParams
	NmxcDomain                = resources.NmxcDomain
	NmxcDomainsUpdateResult   = resources.NmxcDomainsUpdateResult
	NmxcFactoryResetResult    = resources.NmxcFactoryResetResult
	NmxcProbeResult           = resources.NmxcProbeResult
	NmxcResetResult           = resources.NmxcResetResult
	OperationStatusItem       = resources.OperationStatusItem
	RMAInfoItem               = resources.RMAInfoItem
	RMAInfoRecord             = resources.RMAInfoRecord
	RMAStatusItem             = resources.RMAStatusItem
	SuidMap                   = resources.SuidMap
	TriggerBootstrapResult    = resources.TriggerBootstrapResult
	UfmCredsResult            = resources.UfmCredsResult
	UploadFileResult          = resources.UploadFileResult
	UploadUIObjectArgs        = resources.UploadUIObjectArgs
	WebhookDeliveryStatus     = resources.WebhookDeliveryStatus
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
	HostActionAdd             = resources.HostActionAdd
	HostActionDelete          = resources.HostActionDelete
	NetOpsActionVlan          = resources.NetOpsActionVlan
	NetOpsActionPortChange    = resources.NetOpsActionPortChange
	NetOpsActionVlanMember    = resources.NetOpsActionVlanMember
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
