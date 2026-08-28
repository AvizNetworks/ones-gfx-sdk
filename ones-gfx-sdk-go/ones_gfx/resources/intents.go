package resources

import (
	"context"
	"net/url"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// IntentsResource covers UI intent upload/retrieval and Day-1 intent status.
type IntentsResource struct {
	transport *ones_gfx.Transport
}

// NewIntentsResource constructs an IntentsResource.
func NewIntentsResource(transport *ones_gfx.Transport) *IntentsResource {
	return &IntentsResource{transport: transport}
}

// The following types mirror the Intent model graph (Models/*.java). They are
// used only by UploadUIObject and GetUIObject, so they live here.

type DevicefactsItem struct {
	IPAddress    string `json:"ipAddress,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Iv           string `json:"iv,omitempty"`
	AgentVersion string `json:"agentVersion,omitempty"`
	NosVersion   string `json:"nosVersion,omitempty"`
	Config       string `json:"config,omitempty"`
	Isactive     *bool  `json:"isactive,omitempty"`
}

type ConnectionItem struct {
	LocalIP        string `json:"localIp,omitempty"`
	RemoteIP       string `json:"remoteIp,omitempty"`
	LocalIntName   string `json:"localIntName,omitempty"`
	RemoteIntName  string `json:"remoteIntName,omitempty"`
	LocalSwitchID  string `json:"localSwitchId,omitempty"`
	RemoteSwitchID string `json:"remoteSwitchId,omitempty"`
	RemoteAsn      string `json:"remoteAsn,omitempty"`
	LocalAsn       string `json:"localAsn,omitempty"`
	NetworkAddress string `json:"networkAddress,omitempty"`
	IsMclagLink    *bool  `json:"isMclagLink,omitempty"`
	IsHostLink     *bool  `json:"isHostLink,omitempty"`
	LocalPOIp      string `json:"localPOIp,omitempty"`
	RemotePOIp     string `json:"remotePOIp,omitempty"`
	LocalPO        string `json:"localPO,omitempty"`
	RemotePO       string `json:"remotePO,omitempty"`
	Vlans          string `json:"vlans,omitempty"`
	Mode           string `json:"mode,omitempty"`
	MclagPOGroup   string `json:"mclagPOGroup,omitempty"`
	IsMhEnabled    *bool  `json:"isMhEnabled,omitempty"`
}

type StatusItem struct {
	IntentName         string `json:"intentName,omitempty"`
	SubIntentName      string `json:"subIntentName,omitempty"`
	ConfigStatus       string `json:"config_status,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty"`
	ConfigLogs         string `json:"config_logs,omitempty"`
	VerificationLogs   string `json:"verification_logs,omitempty"`
	Isactive           *bool  `json:"isactive,omitempty"`
}

type SvimapperItem struct {
	Vlan     string `json:"vlan,omitempty"`
	Network  string `json:"network,omitempty"`
	Isactive *bool  `json:"isactive,omitempty"`
}

type DeviceItem struct {
	Hostname                string           `json:"hostname,omitempty"`
	SwitchID                *int             `json:"switchId,omitempty"`
	Asn                     string           `json:"asn,omitempty"`
	LoopbackIP              string           `json:"loopbackIp,omitempty"`
	LoopbackIntName         string           `json:"loopbackIntName,omitempty"`
	IsSSpine                *bool            `json:"isSSpine,omitempty"`
	IsSpine                 *bool            `json:"isSpine,omitempty"`
	IsLeaf                  *bool            `json:"isLeaf,omitempty"`
	IsTor                   *bool            `json:"isTor,omitempty"`
	IsHost                  *bool            `json:"isHost,omitempty"`
	IsDpu                   *bool            `json:"isDpu,omitempty"`
	IsMclagEnabled          *bool            `json:"isMclagEnabled,omitempty"`
	MclagKeepaliveVlan      string           `json:"mclagKeepaliveVlan,omitempty"`
	MclagPeerIdentifier     string           `json:"mclagPeerIdentifier,omitempty"`
	IsMclagOverL3           *bool            `json:"isMclagOverL3,omitempty"`
	VtepLoopbackIdentifier  string           `json:"vtepLoopbackIdentifier,omitempty"`
	VtepLoopbackInterface   string           `json:"vtepLoopbackInterface,omitempty"`
	VtepLoopbackInterfaceIp string           `json:"vtepLoopbackInterfaceIp,omitempty"`
	Vrfs                    string           `json:"vrfs,omitempty"`
	Vlans                   string           `json:"vlans,omitempty"`
	IsSflowEnabled          *bool            `json:"isSflowEnabled,omitempty"`
	SFlowAgent              string           `json:"sFlowAgent,omitempty"`
	SFlowPollInterval       string           `json:"sFlowPollInterval,omitempty"`
	Isactive                *bool            `json:"isactive,omitempty"`
	UfmSystemName           string           `json:"ufmSystemName,omitempty"`
	Devicefacts             *DevicefactsItem `json:"devicefacts,omitempty"`
	Connections             []ConnectionItem `json:"connections,omitempty"`
	Statuses                []StatusItem     `json:"statuses,omitempty"`
	Svis                    []SvimapperItem  `json:"svis,omitempty"`
}

type QosmappingItem struct {
	TrafficClass  *int   `json:"traffic_class,omitempty"`
	Dscp          string `json:"dscp,omitempty"`
	Dot1p         string `json:"dot1p,omitempty"`
	Queues        *int   `json:"queues,omitempty"`
	PriorityGroup *int   `json:"priority_group,omitempty"`
}

type SchedulerItem struct {
	Queues      *int   `json:"queues,omitempty"`
	ProfileName string `json:"profile_name,omitempty"`
	Weights     *int   `json:"weights,omitempty"`
}

type QosItem struct {
	DscpTcMap            string `json:"dscp_tc_map,omitempty"`
	Dot1pTcMap           string `json:"dot1p_tc_map,omitempty"`
	TcPgMap              string `json:"tc_pg_map,omitempty"`
	TcQueueMap           string `json:"tc_queue_map,omitempty"`
	PfcEnabledQueues     string `json:"pfc_enabled_queues,omitempty"`
	PfcwdAction          string `json:"pfcwd_action,omitempty"`
	PfcwdDetectionTime   *int   `json:"pfcwd_detection_time,omitempty"`
	PfcwdRestorationTime *int   `json:"pfcwd_restoration_time,omitempty"`
	IsBrsEnabled         *bool  `json:"is_brs_enabled,omitempty"`
	WredPolicy           string `json:"wred_policy,omitempty"`
	EcnMode              string `json:"ecn_mode,omitempty"`
	EcnCnpQueue          *int   `json:"ecn_cnp_queue,omitempty"`
	EcnGmin              *int   `json:"ecn_gmin,omitempty"`
	EcnGmax              *int   `json:"ecn_gmax,omitempty"`
	EcnGmark             *int   `json:"ecn_gmark,omitempty"`
	EcnRmin              *int   `json:"ecn_rmin,omitempty"`
	EcnRmax              *int   `json:"ecn_rmax,omitempty"`
	EcnRmark             *int   `json:"ecn_rmark,omitempty"`
	EcnYmin              *int   `json:"ecn_ymin,omitempty"`
	EcnYmax              *int   `json:"ecn_ymax,omitempty"`
	EcnYmark             *int   `json:"ecn_ymark,omitempty"`
	SchedulerType        string `json:"scheduler_type,omitempty"`
}

type VlanvnimapperItem struct {
	Vlan     string `json:"vlan,omitempty"`
	Vni      string `json:"vni,omitempty"`
	Vrf      string `json:"vrf,omitempty"`
	Network  string `json:"network,omitempty"`
	IsIRB    *bool  `json:"isIRB,omitempty"`
	Isactive *bool  `json:"isactive,omitempty"`
}

type ParametersItem struct {
	Vlan                string              `json:"vlan,omitempty"`
	Vni                 string              `json:"vni,omitempty"`
	AnycastGateway      string              `json:"anycastGateway,omitempty"`
	AnycastMac          string              `json:"anycastMac,omitempty"`
	HostsPerVlan        *int                `json:"hostsPerVlan,omitempty"`
	HostsPerL3Subnet    *int                `json:"hostsPerL3Subnet,omitempty"`
	NetworkAddress      string              `json:"networkAddress,omitempty"`
	PrefixLen           *int                `json:"prefixLen,omitempty"`
	PossibleSubNetworks *int                `json:"possibleSubNetworks,omitempty"`
	PossibleHost        *int                `json:"possibleHost,omitempty"`
	Logs                string              `json:"logs,omitempty"`
	IrbVlans            string              `json:"irbVlans,omitempty"`
	IrbVnis             string              `json:"irbVnis,omitempty"`
	RoutingSymmetric    *bool               `json:"routing_symmetric,omitempty"`
	Isactive            *bool               `json:"isactive,omitempty"`
	VlanVniMapper       []VlanvnimapperItem `json:"vlanVniMapper,omitempty"`
}

// IntentItem mirrors Models/Intent.java JSON (GET /getUIObject).
type IntentItem struct {
	ID                      int              `json:"id,omitempty"`
	Name                    string           `json:"name,omitempty"`
	FabricID                *int             `json:"fabricId,omitempty"`
	OrchestrationMode       string           `json:"orchestrationMode,omitempty"`
	SspineCount             *int             `json:"sspineCount,omitempty"`
	SpineCount              *int             `json:"spineCount,omitempty"`
	LeafCount               *int             `json:"leafCount,omitempty"`
	Fec                     string           `json:"fec,omitempty"`
	Mtu                     string           `json:"mtu,omitempty"`
	AdminStatus             string           `json:"adminStatus,omitempty"`
	AsnSSpine               string           `json:"asnSSpine,omitempty"`
	IsAsnSSpineUnique       *bool            `json:"isAsnSSpineUnique,omitempty"`
	AsnLeaf                 string           `json:"asnLeaf,omitempty"`
	IsAsnLeafUnique         *bool            `json:"isAsnLeafUnique,omitempty"`
	AsnSpine                string           `json:"asnSpine,omitempty"`
	IsAsnSpineUnique        *bool            `json:"isAsnSpineUnique,omitempty"`
	NtpServer               string           `json:"ntpServer,omitempty"`
	Timezone                string           `json:"timezone,omitempty"`
	SysLogServer            string           `json:"sysLogServer,omitempty"`
	SnmpServer              string           `json:"snmpServer,omitempty"`
	IsBGP_U                 *bool            `json:"isBGP_U,omitempty"`
	IsTorPresent            *bool            `json:"isTorPresent,omitempty"`
	ND_RA                   *int             `json:"ND_RA,omitempty"`
	Isactive                *bool            `json:"isactive,omitempty"`
	Logs                    string           `json:"logs,omitempty"`
	Issag                   *bool            `json:"issag,omitempty"`
	DhcpServerIps           string           `json:"dhcpServerIps,omitempty"`
	DhcpRelaySrcInterface   string           `json:"dhcpRelaySrcInterface,omitempty"`
	DhcpRelaySrcInterfaceIP string           `json:"dhcpRelaySrcInterfaceIP,omitempty"`
	Devices                 []DeviceItem     `json:"devices,omitempty"`
	Qosmappings             []QosmappingItem `json:"qosmappings,omitempty"`
	Schedulers              []SchedulerItem  `json:"schedulers,omitempty"`
	QoS                     *QosItem         `json:"QoS,omitempty"`
	Parameters              *ParametersItem  `json:"parameters,omitempty"`
	CreatedAt               string           `json:"createdAt,omitempty"`
	UpdatedAt               string           `json:"updatedAt,omitempty"`
}

// intentPayload is the Intent request body of UploadUIObject.
type intentPayload struct {
	ID                      *int             `json:"id,omitempty"`
	Name                    *string          `json:"name,omitempty"`
	FabricID                *int             `json:"fabricId,omitempty"`
	OrchestrationMode       *string          `json:"orchestrationMode,omitempty"`
	SspineCount             *int             `json:"sspineCount,omitempty"`
	SpineCount              *int             `json:"spineCount,omitempty"`
	LeafCount               *int             `json:"leafCount,omitempty"`
	Fec                     *string          `json:"fec,omitempty"`
	Mtu                     *string          `json:"mtu,omitempty"`
	AdminStatus             *string          `json:"adminStatus,omitempty"`
	AsnSSpine               *string          `json:"asnSSpine,omitempty"`
	IsAsnSSpineUnique       *bool            `json:"isAsnSSpineUnique,omitempty"`
	AsnLeaf                 *string          `json:"asnLeaf,omitempty"`
	IsAsnLeafUnique         *bool            `json:"isAsnLeafUnique,omitempty"`
	AsnSpine                *string          `json:"asnSpine,omitempty"`
	IsAsnSpineUnique        *bool            `json:"isAsnSpineUnique,omitempty"`
	NtpServer               *string          `json:"ntpServer,omitempty"`
	Timezone                *string          `json:"timezone,omitempty"`
	SysLogServer            *string          `json:"sysLogServer,omitempty"`
	SnmpServer              *string          `json:"snmpServer,omitempty"`
	IsBGP_U                 *bool            `json:"isBGP_U,omitempty"`
	IsTorPresent            *bool            `json:"isTorPresent,omitempty"`
	ND_RA                   *int             `json:"ND_RA,omitempty"`
	Isactive                *bool            `json:"isactive,omitempty"`
	Logs                    *string          `json:"logs,omitempty"`
	Issag                   *bool            `json:"issag,omitempty"`
	DhcpServerIps           *string          `json:"dhcpServerIps,omitempty"`
	DhcpRelaySrcInterface   *string          `json:"dhcpRelaySrcInterface,omitempty"`
	DhcpRelaySrcInterfaceIP *string          `json:"dhcpRelaySrcInterfaceIP,omitempty"`
	Devices                 []DeviceItem     `json:"devices,omitempty"`
	Qosmappings             []QosmappingItem `json:"qosmappings,omitempty"`
	Schedulers              []SchedulerItem  `json:"schedulers,omitempty"`
	QoS                     *QosItem         `json:"QoS,omitempty"`
	Parameters              *ParametersItem  `json:"parameters,omitempty"`
}

// UploadUIObjectArgs carries the optional intent fields for UploadUIObject.
type UploadUIObjectArgs struct {
	ID                      *int
	Name                    *string
	FabricID                *int
	OrchestrationMode       *string
	SspineCount             *int
	SpineCount              *int
	LeafCount               *int
	Fec                     *string
	Mtu                     *string
	AdminStatus             *string
	AsnSSpine               *string
	IsAsnSSpineUnique       *bool
	AsnLeaf                 *string
	IsAsnLeafUnique         *bool
	AsnSpine                *string
	IsAsnSpineUnique        *bool
	NtpServer               *string
	Timezone                *string
	SysLogServer            *string
	SnmpServer              *string
	IsBGP_U                 *bool
	IsTorPresent            *bool
	ND_RA                   *int
	Isactive                *bool
	Logs                    *string
	Issag                   *bool
	DhcpServerIps           *string
	DhcpRelaySrcInterface   *string
	DhcpRelaySrcInterfaceIP *string
	Devices                 []DeviceItem
	Qosmappings             []QosmappingItem
	Schedulers              []SchedulerItem
	QoS                     *QosItem
	Parameters              *ParametersItem
}

func (a *UploadUIObjectArgs) payload() *intentPayload {
	return &intentPayload{
		ID: a.ID, Name: a.Name, FabricID: a.FabricID, OrchestrationMode: a.OrchestrationMode,
		SspineCount: a.SspineCount, SpineCount: a.SpineCount, LeafCount: a.LeafCount,
		Fec: a.Fec, Mtu: a.Mtu, AdminStatus: a.AdminStatus, AsnSSpine: a.AsnSSpine,
		IsAsnSSpineUnique: a.IsAsnSSpineUnique, AsnLeaf: a.AsnLeaf, IsAsnLeafUnique: a.IsAsnLeafUnique,
		AsnSpine: a.AsnSpine, IsAsnSpineUnique: a.IsAsnSpineUnique, NtpServer: a.NtpServer,
		Timezone: a.Timezone, SysLogServer: a.SysLogServer, SnmpServer: a.SnmpServer,
		IsBGP_U: a.IsBGP_U, IsTorPresent: a.IsTorPresent, ND_RA: a.ND_RA, Isactive: a.Isactive,
		Logs: a.Logs, Issag: a.Issag, DhcpServerIps: a.DhcpServerIps,
		DhcpRelaySrcInterface: a.DhcpRelaySrcInterface, DhcpRelaySrcInterfaceIP: a.DhcpRelaySrcInterfaceIP,
		Devices: a.Devices, Qosmappings: a.Qosmappings, Schedulers: a.Schedulers,
		QoS: a.QoS, Parameters: a.Parameters,
	}
}

// UploadUIObject stores a UI-built intent and triggers orchestration.
// Maps to POST /uploadUIObject. Use WithRequestOrigin("ones-ui") for internal
// callers. Returns the "Object stored successfully" message.
func (r *IntentsResource) UploadUIObject(ctx context.Context, args *UploadUIObjectArgs, opts ...ones_gfx.CallOption) (string, error) {
	if args == nil {
		args = &UploadUIObjectArgs{}
	}
	rc := ones_gfx.ResolveCallOptions(opts...)
	res, err := ones_gfx.Call[string](r.transport, "POST", "uploadUIObject", args.payload(), ones_gfx.OperationModeSynchronous, rc.ReqOpts())
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// GetUIObject returns an orchestrated intent (by name, or the last one).
// Maps to GET /getUIObject?name=...
func (r *IntentsResource) GetUIObject(ctx context.Context, name *string) (*IntentItem, error) {
	q := url.Values{}
	if name != nil {
		q.Set("name", *name)
	}
	return ones_gfx.Call[IntentItem](r.transport, "GET", "getUIObject", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
}

// LastOrchestratedName returns the name of the last orchestrated intent.
// Maps to GET /getLastOrchestratedIntentName.
func (r *IntentsResource) LastOrchestratedName(ctx context.Context) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, "GET", "getLastOrchestratedIntentName", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Validation returns intent validation results. Dynamic payload.
// Maps to GET /getIntentValidation?intentName=...
func (r *IntentsResource) Validation(ctx context.Context, intentName string) ([]interface{}, error) {
	q := url.Values{}
	q.Set("intentName", intentName)
	res, err := ones_gfx.Call[[]interface{}](r.transport, "GET", "getIntentValidation", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Day1ConfigStatus returns Day-1 config status. Dynamic payload.
// Maps to GET /getDay1ConfigStatus?intentName=...
func (r *IntentsResource) Day1ConfigStatus(ctx context.Context, intentName string) ([]interface{}, error) {
	q := url.Values{}
	q.Set("intentName", intentName)
	res, err := ones_gfx.Call[[]interface{}](r.transport, "GET", "getDay1ConfigStatus", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// DerivationLogs returns intent derivation logs for one device.
// Dynamic payload. Maps to GET /getIntentDerivationLogs?device=...
func (r *IntentsResource) DerivationLogs(ctx context.Context, device string) ([]interface{}, error) {
	q := url.Values{}
	q.Set("device", device)
	res, err := ones_gfx.Call[[]interface{}](r.transport, "GET", "getIntentDerivationLogs", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
