package resources

import (
	"context"
	"net/url"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// RMAResource covers RMA staging, triggering and status queries.
type RMAResource struct {
	transport *ones_gfx.Transport
}

// NewRMAResource constructs an RMAResource.
func NewRMAResource(transport *ones_gfx.Transport) *RMAResource {
	return &RMAResource{transport: transport}
}

// RMAInfoItem mirrors Models/RMAInfo.java as a request payload.
type RMAInfoItem struct {
	ReplacingDeviceMac    string `json:"replacingDeviceMac,omitempty"`
	ReplacingDeviceIP     string `json:"replacingDeviceIp,omitempty"`
	ReplacingDeviceSerial string `json:"replacingDeviceSerial,omitempty"`
	ReplacetoDeviceMac    string `json:"replacetoDeviceMac,omitempty"`
	ReplacetoDeviceIP     string `json:"replacetoDeviceIp,omitempty"`
	ReplacetoDeviceSerial string `json:"replacetoDeviceSerial,omitempty"`
	FabricID              string `json:"fabricId,omitempty"`
	State                 string `json:"state,omitempty"`
	MarkedTime            string `json:"markedTime,omitempty"`
	TicketID              string `json:"ticketId,omitempty"`
	BackupSelected        string `json:"backupSelected,omitempty"`
	ScheduledTime         string `json:"scheduledTime,omitempty"`
	IsScheduledForLater   bool   `json:"is_scheduled_for_later,omitempty"`
	TriggerTime           string `json:"triggertime,omitempty"`
}

// RMAInfoRecord mirrors Models/RMAInfo.java JSON (responses).
type RMAInfoRecord struct {
	ID                    int    `json:"id,omitempty"`
	ReplacingDeviceMac    string `json:"replacingDeviceMac,omitempty"`
	ReplacingDeviceIP     string `json:"replacingDeviceIp,omitempty"`
	ReplacingDeviceSerial string `json:"replacingDeviceSerial,omitempty"`
	ReplacetoDeviceMac    string `json:"replacetoDeviceMac,omitempty"`
	ReplacetoDeviceIP     string `json:"replacetoDeviceIp,omitempty"`
	ReplacetoDeviceSerial string `json:"replacetoDeviceSerial,omitempty"`
	FabricID              string `json:"fabricId,omitempty"`
	State                 string `json:"state,omitempty"`
	MarkedTime            string `json:"markedTime,omitempty"`
	TicketID              string `json:"ticketId,omitempty"`
	BackupSelected        string `json:"backupSelected,omitempty"`
	ScheduledTime         string `json:"scheduledTime,omitempty"`
	IsScheduledForLater   bool   `json:"is_scheduled_for_later,omitempty"`
	TriggerTime           string `json:"triggertime,omitempty"`
}

// RMAStatusItem mirrors Models/RMAStatus.java JSON.
type RMAStatusItem struct {
	ID        int    `json:"id,omitempty"`
	RmaInfoID *int   `json:"rmainfoId,omitempty"`
	Task      string `json:"task,omitempty"`
	Status    *int   `json:"status,omitempty"`
	StartTime string `json:"starttime,omitempty"`
	EndTime   string `json:"endtime,omitempty"`
	Logs      string `json:"logs,omitempty"`
	TaskTitle string `json:"tasktitle,omitempty"`
}

// Fill stages RMA replacement records without triggering the swap.
// Maps to POST /fillrmaconfig.
func (r *RMAResource) Fill(ctx context.Context, items []RMAInfoItem) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "fillrmaconfig", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// Trigger executes the RMA swap. Maps to POST /triggerrma.
func (r *RMAResource) Trigger(ctx context.Context, items []RMAInfoItem) (bool, error) {
	res, err := ones_gfx.Call[bool](r.transport, "POST", "triggerrma", items, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return false, err
	}
	return *res, nil
}

// List returns all RMA-marked devices. Maps to GET /getrmainfo.
func (r *RMAResource) List(ctx context.Context) ([]RMAInfoRecord, error) {
	res, err := ones_gfx.Call[[]RMAInfoRecord](r.transport, "GET", "getrmainfo", nil, ones_gfx.OperationModeSynchronous, nil)
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}

// Status returns RMA stage progress for one RMA record.
// Maps to GET /getrmastatus?rmaInfoId=...
func (r *RMAResource) Status(ctx context.Context, rmaInfoID int) ([]RMAStatusItem, error) {
	q := url.Values{}
	q.Set("rmaInfoId", strconv.Itoa(rmaInfoID))
	res, err := ones_gfx.Call[[]RMAStatusItem](r.transport, "GET", "getrmastatus", nil, ones_gfx.OperationModeSynchronous, &ones_gfx.ReqOpts{Query: q})
	if err != nil || res == nil {
		return nil, err
	}
	return *res, nil
}
