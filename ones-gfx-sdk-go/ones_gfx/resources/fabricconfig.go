package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FabricConfigResource covers the ONES UI config service: the endpoints the UI
// (and terraform-provider-fabricapi) use to generate a fabric's topology and
// YAML, deploy it to real devices, and delete it. These live at the root of
// the host (/api/config/..., /fabrics/{name}), not under /api/fm/.
type FabricConfigResource struct {
	transport *ones_gfx.Transport
}

// NewFabricConfigResource constructs a FabricConfigResource. The transport's
// base URL must be the host root (see ones_gfx.ConfigBaseURL).
func NewFabricConfigResource(transport *ones_gfx.Transport) *FabricConfigResource {
	return &FabricConfigResource{transport: transport}
}

// FabricConfigArgs describes a fabric to generate. Type and HostsPerSU are
// required. Zero values take the same defaults as the ONES UI: Status "Draft",
// Instance "fm", NodeType "dgx", SimulationID 1, TenantCtrl "ones",
// NSOperatingSystem "cumulus", NumOfSUs 1, MaxNumOfSUs = NumOfSUs, and HostMap
// derived from HostsPerSU.
type FabricConfigArgs struct {
	Type        string // "NVIDIA SpX RA 1.3", "NVIDIA SpX RA 2.1" or "Aviz RA 1.0"
	Status      string
	Description string

	NumOfSUs     int
	MaxNumOfSUs  int
	HostsPerSU   string            // SU index -> host count, e.g. "{0:2}" or "{0:1,1:2}"
	HostMap      map[string]string // derived from HostsPerSU when empty
	NodeType     string            // dgx, gb200, gb300, b300_32, b300_64, rtxpro_4, rtxpro_8
	TenantCtrl   string
	SimulationID int
	Instance     string

	// East-west (GPU) network. StartingSubnetGPU is required when EnableEW.
	EnableEW          bool
	StartingSubnetGPU string

	// North-south (front-end) network. StartingSubnetCPU is required when
	// EnableNS; StartingSubnetStorage is required when DedicatedStorage.
	// NSOperatingSystem ("cumulus" or "sonic") must match the real
	// front-end/storage switches or switch validation fails at deploy time.
	EnableNS              bool
	NSOperatingSystem     string
	DedicatedStorage      bool
	StartingSubnetCPU     string
	StartingSubnetStorage string
}

// fabricConfigPayload is the POST /api/config/addFabricData body.
type fabricConfigPayload struct {
	Name                  string            `json:"name"`
	Type                  string            `json:"type"`
	Status                string            `json:"status"`
	Description           string            `json:"description"`
	NumOfSUs              int               `json:"numOfSus"`
	MaxNumOfSUs           int               `json:"maxNumOfSus"`
	HostMap               map[string]string `json:"hostMap"`
	StartingSubnetGPU     string            `json:"startingSubnetGpu"`
	SimulationID          int               `json:"simulationId"`
	EnableEW              bool              `json:"enableEW"`
	SUHostCount           string            `json:"suHostCnt"`
	Tenant                string            `json:"tenant"`
	Instance              string            `json:"instance"`
	EnableNS              bool              `json:"enableNS"`
	IsONESControlled      bool              `json:"isOnesControlled"`
	UserAndStorage        bool              `json:"userandstorage"`
	DedicatedStorage      bool              `json:"dedicatedStorage"`
	FrontendStorage       bool              `json:"frontendStorage"`
	StartingSubnetCPU     string            `json:"startingSubnetCpu,omitempty"`
	StartingSubnetStorage string            `json:"startingSubnetStorage,omitempty"`
	NodeType              string            `json:"nodeType"`
	OperatingSystemE      string            `json:"operatingSystemE"`
	OperatingSystemF      string            `json:"operatingSystemF"`
	OperatingSystemS      string            `json:"operatingSystemS"`
}

func orDefault(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

// HostMapFromHostsPerSU parses the "{0:1,1:2}" SU host-count shorthand into
// the {"0": "1", "1": "2"} map addFabricData also expects.
func HostMapFromHostsPerSU(raw string) (map[string]string, error) {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{"), "}"))
	hostMap := map[string]string{}
	if trimmed == "" {
		return hostMap, nil
	}
	for _, pair := range strings.Split(trimmed, ",") {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" || strings.TrimSpace(kv[1]) == "" {
			return nil, fmt.Errorf("expected \"su:count\" pairs like \"{0:1,1:2}\", got %q", pair)
		}
		hostMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return hostMap, nil
}

func (a FabricConfigArgs) payload(name string) (*fabricConfigPayload, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("fabric name is required")
	}
	if a.EnableEW && strings.TrimSpace(a.StartingSubnetGPU) == "" {
		return nil, errors.New("StartingSubnetGPU is required when EnableEW is true")
	}
	if a.EnableNS && strings.TrimSpace(a.StartingSubnetCPU) == "" {
		return nil, errors.New("StartingSubnetCPU is required when EnableNS is true")
	}
	if a.DedicatedStorage && strings.TrimSpace(a.StartingSubnetStorage) == "" {
		return nil, errors.New("StartingSubnetStorage is required when DedicatedStorage is true")
	}
	if strings.TrimSpace(a.HostsPerSU) == "" {
		return nil, errors.New("HostsPerSU is required, e.g. \"{0:2}\"")
	}
	numSUs, maxSUs := a.NumOfSUs, a.MaxNumOfSUs
	if numSUs == 0 {
		numSUs = 1
	}
	if maxSUs == 0 {
		maxSUs = numSUs
	}
	hostMap := a.HostMap
	if len(hostMap) == 0 {
		derived, err := HostMapFromHostsPerSU(a.HostsPerSU)
		if err != nil {
			return nil, err
		}
		hostMap = derived
	}
	simulationID := a.SimulationID
	if simulationID == 0 {
		simulationID = 1
	}
	nsOS := orDefault(a.NSOperatingSystem, "cumulus")
	return &fabricConfigPayload{
		Name:              name,
		Type:              a.Type,
		Status:            orDefault(a.Status, "Draft"),
		Description:       a.Description,
		NumOfSUs:          numSUs,
		MaxNumOfSUs:       maxSUs,
		HostMap:           hostMap,
		StartingSubnetGPU: a.StartingSubnetGPU,
		SimulationID:      simulationID,
		EnableEW:          a.EnableEW,
		SUHostCount:       a.HostsPerSU,
		Tenant:            orDefault(a.TenantCtrl, "ones"),
		Instance:          orDefault(a.Instance, "fm"),
		EnableNS:          a.EnableNS,
		// Always ONES-controlled, matching the UI default for RA fabrics.
		IsONESControlled: true,
		// The UI toggles dedicatedStorage and frontendStorage together;
		// userandstorage is the shared (NS on, no dedicated storage) case.
		UserAndStorage:        a.EnableNS && !a.DedicatedStorage,
		DedicatedStorage:      a.DedicatedStorage,
		FrontendStorage:       a.DedicatedStorage,
		StartingSubnetCPU:     a.StartingSubnetCPU,
		StartingSubnetStorage: a.StartingSubnetStorage,
		NodeType:              orDefault(a.NodeType, "dgx"),
		// East-west is always Cumulus in the UI.
		OperatingSystemE: "cumulus",
		OperatingSystemF: nsOS,
		OperatingSystemS: nsOS,
	}, nil
}

func (r *FabricConfigResource) callString(method, path string, body any, ro *ones_gfx.ReqOpts) (string, error) {
	res, err := ones_gfx.Call[string](r.transport, method, path, body, ones_gfx.OperationModeSynchronous, ro)
	if err != nil || res == nil {
		return "", err
	}
	return strings.TrimSpace(*res), nil
}

// Create generates a fabric (topology, YAML and inventory) in Draft state.
// Maps to POST /api/config/addFabricData. Returns the raw response body.
func (r *FabricConfigResource) Create(ctx context.Context, name string, args FabricConfigArgs) (string, error) {
	body, err := args.payload(name)
	if err != nil {
		return "", err
	}
	return r.callString("POST", "api/config/addFabricData", body, nil)
}

// GetYAML returns the fabric's generated YAML. Maps to GET /fabrics/{name}.
func (r *FabricConfigResource) GetYAML(ctx context.Context, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("fabric name is required")
	}
	res, err := ones_gfx.Call[string](r.transport, "GET", "fabrics/"+url.PathEscape(name), nil,
		ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Accept: "*/*"})
	if err != nil || res == nil {
		return "", err
	}
	return *res, nil
}

// Delete deletes a fabric. A fabric that is already gone (404) is not an
// error. Maps to DELETE /api/config/deletefabricbyname/{name}.
//
// Deployed fabrics have no undeploy call: this removes the fabric record, not
// the configuration already pushed to the switches. It also leaves the
// fabric's inventory rows behind, so recreating a fabric with the same device
// IPs can produce duplicate inventory entries.
func (r *FabricConfigResource) Delete(ctx context.Context, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("fabric name is required")
	}
	res, err := r.callString("DELETE", "api/config/deletefabricbyname/"+url.PathEscape(name), nil, nil)
	var notFound *ones_gfx.NotFoundError
	if errors.As(err, &notFound) {
		return "", nil
	}
	return res, err
}

// FabricDevice is one switch, server or DPU in a fabric deploy. The JSON tags
// match the devices.json format used by terraform-provider-fabricapi.
type FabricDevice struct {
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceType string `json:"device_type,omitempty"` // switch, server, dpu, host
	DeviceRole string `json:"device_role,omitempty"` // spine, leaf, host
	// ApplyConfig controls whether ONES configures the device; nil means true.
	ApplyConfig *bool `json:"apply_config,omitempty"`
}

func (d FabricDevice) applyConfig() bool { return d.ApplyConfig == nil || *d.ApplyConfig }

// IsSwitch reports whether the device is a spine/leaf switch.
func (d FabricDevice) IsSwitch() bool {
	t, r := strings.ToLower(strings.TrimSpace(d.DeviceType)), strings.ToLower(strings.TrimSpace(d.DeviceRole))
	return t == "switch" || r == "spine" || r == "leaf"
}

// IsServer reports whether the device is a host, server or DPU.
func (d FabricDevice) IsServer() bool {
	t, r := strings.ToLower(strings.TrimSpace(d.DeviceType)), strings.ToLower(strings.TrimSpace(d.DeviceRole))
	return t == "server" || t == "dpu" || t == "host" || r == "host"
}

// The backend expects "Yes"/"No" strings, as in the UI's device spreadsheet.
func yesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// UploadDeviceIPs records each device's real IP and credentials against the
// fabric. Maps to POST /api/config/uploadip.
func (r *FabricConfigResource) UploadDeviceIPs(ctx context.Context, fabricName string, devices []FabricDevice) (string, error) {
	type entry struct {
		Hostname string `json:"hostname"`
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		// The UI sends both casings; so do we.
		ApplyConfig string `json:"applyConfig"`
		Applyconfig string `json:"applyconfig"`
		Configure   string `json:"configure"`
	}
	entries := make([]entry, 0, len(devices))
	for _, d := range devices {
		yn := yesNo(d.applyConfig())
		entries = append(entries, entry{d.Hostname, d.IP, d.Username, d.Password, yn, yn, yn})
	}
	body := struct {
		UpdatedDevices []entry `json:"updatedDevices"`
		FabricName     string  `json:"fabricName"`
	}{entries, fabricName}
	return r.callString("POST", "api/config/uploadip", body, nil)
}

// ValidateDeviceResult is one device's result from ValidateSwitches or
// ValidateServers. A switch passed when Build is set; a server passed when OS
// is set; Error carries the reason for a failure.
type ValidateDeviceResult struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Username string `json:"username"`
	Build    string `json:"build,omitempty"`
	OS       string `json:"os,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (r *FabricConfigResource) validate(path, fabricName string, devices []FabricDevice) ([]ValidateDeviceResult, error) {
	type entry struct {
		Hostname    string `json:"hostname"`
		IP          string `json:"ip"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DeviceType  string `json:"deviceType,omitempty"`
		DeviceRole  string `json:"deviceRole,omitempty"`
		ApplyConfig string `json:"applyConfig,omitempty"`
		FabricName  string `json:"fabricName,omitempty"`
	}
	entries := make([]entry, 0, len(devices))
	for _, d := range devices {
		entries = append(entries, entry{d.Hostname, d.IP, d.Username, d.Password, d.DeviceType, d.DeviceRole,
			yesNo(d.applyConfig()), fabricName})
	}
	res, err := ones_gfx.Call[[]ValidateDeviceResult](r.transport, "POST", path, entries, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// ValidateSwitches logs in to each switch over SSH and checks it is reachable
// and running the expected NOS. Maps to POST /api/config/validateswitch.
func (r *FabricConfigResource) ValidateSwitches(ctx context.Context, fabricName string, devices []FabricDevice) ([]ValidateDeviceResult, error) {
	return r.validate("api/config/validateswitch", fabricName, devices)
}

// ValidateServers logs in to each host/DPU over SSH and checks its OS. Maps to
// POST /api/config/validateserver.
func (r *FabricConfigResource) ValidateServers(ctx context.Context, fabricName string, devices []FabricDevice) ([]ValidateDeviceResult, error) {
	return r.validate("api/config/validateserver", fabricName, devices)
}

// UpdateInventory pushes the device credentials into the Fabric Manager
// inventory. Maps to POST /api/config/updateinventory.
func (r *FabricConfigResource) UpdateInventory(ctx context.Context, fabricName string, devices []FabricDevice) (string, error) {
	type entry struct {
		FabricName string `json:"fabricName"`
		Hostname   string `json:"hostname"`
		IPAddress  string `json:"ipAddress"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Configure  string `json:"configure"`
	}
	entries := make([]entry, 0, len(devices))
	for _, d := range devices {
		entries = append(entries, entry{fabricName, d.Hostname, d.IP, d.Username, d.Password, yesNo(d.applyConfig())})
	}
	return r.callString("POST", "api/config/updateinventory", entries, nil)
}

// PushConfig sends the fabric document (its YAML, as JSON) to be configured on
// the devices. Maps to POST /api/config. Pass the JSON with its keys in the
// YAML's original order (see FabricYAMLToJSON); a Go map would sort them.
func (r *FabricConfigResource) PushConfig(ctx context.Context, fabricName, instance string, data json.RawMessage) (string, error) {
	body := struct {
		Data       json.RawMessage `json:"data"`
		FabricName string          `json:"fabricName"`
		Instance   string          `json:"instance"`
	}{data, fabricName, orDefault(instance, "fm")}
	return r.callString("POST", "api/config", body, nil)
}

// UpdateStatus sets the fabric's status, e.g. "Deployed". deploymentType is
// DEFAULT, PARTIAL_CONFIG or BROWNFIELD. Maps to POST
// /api/config/updatefabricstatus.
func (r *FabricConfigResource) UpdateStatus(ctx context.Context, name, status, description, deploymentType string) (string, error) {
	body := struct {
		Name           string `json:"name"`
		Status         string `json:"status"`
		Description    string `json:"description"`
		DeploymentType string `json:"deploymentType"`
	}{name, status, description, orDefault(deploymentType, "DEFAULT")}
	return r.callString("POST", "api/config/updatefabricstatus", body, nil)
}

// FabricDeployOptions tunes Deploy. All fields are optional.
type FabricDeployOptions struct {
	Description    string
	DeploymentType string // DEFAULT (default), PARTIAL_CONFIG, BROWNFIELD
	Instance       string // default "fm"
	// CustomYAML replaces the server-generated fabric YAML.
	CustomYAML string
}

// FabricDeployResult collects the response of every deploy step.
type FabricDeployResult struct {
	UploadIPs      string
	SwitchResults  []ValidateDeviceResult
	ServerResults  []ValidateDeviceResult
	Inventory      string
	ConfigPushed   string
	StatusResponse string
}

// FabricValidationError is returned by Deploy when switch or server
// validation reports a failed device. Nothing has been configured yet.
type FabricValidationError struct {
	Kind     string // "switch" or "server"
	Failures []string
	Results  []ValidateDeviceResult
}

func (e *FabricValidationError) Error() string {
	return fmt.Sprintf("%s validation failed:\n- %s", e.Kind, strings.Join(e.Failures, "\n- "))
}

func validationFailures(results []ValidateDeviceResult, passed func(ValidateDeviceResult) bool) []string {
	var failures []string
	for _, res := range results {
		switch {
		case strings.TrimSpace(res.Error) != "":
			failures = append(failures, fmt.Sprintf("%s (%s): %s", res.Hostname, res.IP, res.Error))
		case !passed(res):
			failures = append(failures, fmt.Sprintf("%s (%s): validation did not report success", res.Hostname, res.IP))
		}
	}
	return failures
}

// Deploy runs the ONES UI's full "Deploy Fabric" sequence for a fabric made
// with Create:
//
//  1. UploadDeviceIPs for every device
//  2. ValidateSwitches for switches with ApplyConfig (skipped if none)
//  3. ValidateServers for hosts/servers/DPUs (skipped if none)
//  4. UpdateInventory for every device
//  5. fetch the fabric YAML (or use opts.CustomYAML), fill in each device's
//     IP and credentials, and PushConfig it
//  6. UpdateStatus to "Deployed"
//
// Any failure stops the sequence; a validation failure returns a
// *FabricValidationError. The partial result is returned alongside errors.
func (r *FabricConfigResource) Deploy(ctx context.Context, fabricName string, devices []FabricDevice, opts *FabricDeployOptions) (*FabricDeployResult, error) {
	if strings.TrimSpace(fabricName) == "" {
		return nil, errors.New("fabric name is required")
	}
	if len(devices) == 0 {
		return nil, errors.New("devices must list every switch/server/DPU in the fabric")
	}
	if opts == nil {
		opts = &FabricDeployOptions{}
	}
	out := &FabricDeployResult{}
	var err error

	if out.UploadIPs, err = r.UploadDeviceIPs(ctx, fabricName, devices); err != nil {
		return out, fmt.Errorf("upload device IPs: %w", err)
	}

	var switches, servers []FabricDevice
	for _, d := range devices {
		if d.IsSwitch() && d.applyConfig() {
			switches = append(switches, d)
		}
		if d.IsServer() {
			servers = append(servers, d)
		}
	}
	if len(switches) > 0 {
		if out.SwitchResults, err = r.ValidateSwitches(ctx, fabricName, switches); err != nil {
			return out, fmt.Errorf("validate switches: %w", err)
		}
		if f := validationFailures(out.SwitchResults, func(v ValidateDeviceResult) bool { return strings.TrimSpace(v.Build) != "" }); len(f) > 0 {
			return out, &FabricValidationError{Kind: "switch", Failures: f, Results: out.SwitchResults}
		}
	}
	if len(servers) > 0 {
		if out.ServerResults, err = r.ValidateServers(ctx, fabricName, servers); err != nil {
			return out, fmt.Errorf("validate servers: %w", err)
		}
		if f := validationFailures(out.ServerResults, func(v ValidateDeviceResult) bool { return strings.TrimSpace(v.OS) != "" }); len(f) > 0 {
			return out, &FabricValidationError{Kind: "server", Failures: f, Results: out.ServerResults}
		}
	}

	if out.Inventory, err = r.UpdateInventory(ctx, fabricName, devices); err != nil {
		return out, fmt.Errorf("update inventory: %w", err)
	}

	rawYAML := opts.CustomYAML
	if strings.TrimSpace(rawYAML) == "" {
		if rawYAML, err = r.GetYAML(ctx, fabricName); err != nil {
			return out, fmt.Errorf("fetch fabric YAML: %w", err)
		}
	}
	data, err := FabricYAMLToJSON(rawYAML, devices)
	if err != nil {
		return out, err
	}
	if out.ConfigPushed, err = r.PushConfig(ctx, fabricName, opts.Instance, data); err != nil {
		return out, fmt.Errorf("push config: %w", err)
	}

	if out.StatusResponse, err = r.UpdateStatus(ctx, fabricName, "Deployed", opts.Description, opts.DeploymentType); err != nil {
		return out, fmt.Errorf("config was pushed but marking the fabric Deployed failed: %w", err)
	}
	return out, nil
}

// FabricYAMLToJSON parses a fabric YAML document, fills in each listed
// device's ipAddress and Credentials (matched by switchName/hostName under
// Connectivity), and converts it to JSON with keys in their original order.
func FabricYAMLToJSON(rawYAML string, devices []FabricDevice) (json.RawMessage, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(rawYAML), &doc); err != nil {
		return nil, fmt.Errorf("parse fabric YAML: %w", err)
	}
	injectDeviceCredentials(&doc, devices)
	data, err := yamlNodeToOrderedJSON(&doc)
	if err != nil {
		return nil, fmt.Errorf("convert fabric YAML to JSON: %w", err)
	}
	return data, nil
}

func findMappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

func setScalarValue(node *yaml.Node, value string) {
	if node == nil {
		return
	}
	node.Kind, node.Tag, node.Value = yaml.ScalarNode, "!!str", value
}

// injectDeviceCredentials walks Connectivity's tiers (SSpine, Spine, Leaf,
// Tor, Host — each a list of devices) and sets ipAddress and Credentials on
// every device whose switchName/hostName matches one of devices.
func injectDeviceCredentials(doc *yaml.Node, devices []FabricDevice) {
	if doc == nil || len(doc.Content) == 0 {
		return
	}
	connectivity := findMappingValue(doc.Content[0], "Connectivity")
	if connectivity == nil {
		return
	}
	byHostname := make(map[string]FabricDevice, len(devices))
	for _, d := range devices {
		byHostname[strings.ToLower(strings.TrimSpace(d.Hostname))] = d
	}
	for i := 1; i < len(connectivity.Content); i += 2 {
		tier := connectivity.Content[i]
		if tier.Kind != yaml.SequenceNode {
			continue
		}
		for _, device := range tier.Content {
			nameNode := findMappingValue(device, "switchName")
			if nameNode == nil {
				nameNode = findMappingValue(device, "hostName")
			}
			if nameNode == nil {
				continue
			}
			match, ok := byHostname[strings.ToLower(strings.TrimSpace(nameNode.Value))]
			if !ok {
				continue
			}
			setScalarValue(findMappingValue(device, "ipAddress"), match.IP)
			if creds := findMappingValue(device, "Credentials"); creds != nil {
				setScalarValue(findMappingValue(creds, "user"), match.Username)
				setScalarValue(findMappingValue(creds, "password"), match.Password)
			}
		}
	}
}

// yamlNodeToOrderedJSON converts YAML to JSON keeping mapping key order.
func yamlNodeToOrderedJSON(node *yaml.Node) (json.RawMessage, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return json.RawMessage("null"), nil
		}
		return yamlNodeToOrderedJSON(node.Content[0])
	case yaml.MappingNode, yaml.SequenceNode:
		var buf bytes.Buffer
		open, end, step := byte('['), byte(']'), 1
		if node.Kind == yaml.MappingNode {
			open, end, step = '{', '}', 2
		}
		buf.WriteByte(open)
		for i := 0; i+step-1 < len(node.Content); i += step {
			if i > 0 {
				buf.WriteByte(',')
			}
			if step == 2 {
				key, err := json.Marshal(node.Content[i].Value)
				if err != nil {
					return nil, err
				}
				buf.Write(key)
				buf.WriteByte(':')
			}
			val, err := yamlNodeToOrderedJSON(node.Content[i+step-1])
			if err != nil {
				return nil, err
			}
			buf.Write(val)
		}
		buf.WriteByte(end)
		return buf.Bytes(), nil
	case yaml.AliasNode:
		return yamlNodeToOrderedJSON(node.Alias)
	default:
		var v any
		if err := node.Decode(&v); err != nil {
			return nil, err
		}
		return json.Marshal(v)
	}
}
