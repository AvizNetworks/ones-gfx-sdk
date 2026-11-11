// Quick manual test of the Topology Audit SDK surface (GFX-6219/GFX-6225):
// trigger an on-demand audit, then poll until it finishes.
//
// Run from the SDK root:
//
//	go run ./examples/topology_audit --fabric Fab101
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

// --- configuration -----------------------------------------------------------

const (
	baseURL  = "https://localhost:3002"
	username = "superadmin"
	password = "Admin@123456"
)

// -----------------------------------------------------------------------------

func main() {
	fabric := flag.String("fabric", "Fab101", "fabric to audit")
	flag.Parse()

	ctx := context.Background()

	fmt.Printf("Logging in to %s as %q ...\n", baseURL, username)
	client, err := ones.InitializeWithCreds(baseURL, username, password)
	if err != nil {
		fmt.Println("login failed:", err)
		return
	}
	defer client.Close()

	fmt.Printf("\nTriggering topology audit for fabric %q ...\n", *fabric)
	accepted, err := ones.TriggerTopologyAudit(ctx, client, *fabric)
	if err != nil {
		fmt.Println("TriggerTopologyAudit failed:", err)
		return
	}
	fmt.Printf("Accepted: operationId=%s status=%s\n", accepted.OperationID, accepted.Status)

	fmt.Println("\nPolling for the result ...")
	for {
		report, err := ones.GetTopologyAuditResult(ctx, client, accepted.OperationID)
		if errors.Is(err, ones.ErrTopologyAuditNotDone) {
			status, _ := ones.GetTopologyAuditStatus(ctx, client, accepted.OperationID)
			fmt.Printf("  still running (status=%s) ...\n", status.Status)
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			fmt.Println("audit failed:", err)
			return
		}

		fmt.Printf("\nOverall status: %s (intent version %s)\n", report.OverallStatus, report.IntentVersion)
		for _, link := range report.Links {
			fmt.Printf("  %-8s %s:%s -> %s:%s  (observed %s:%s)  %s\n",
				link.Status, link.Device, link.LocalInterface,
				link.RemoteDevice, link.RemoteInterface,
				link.ObservedRemoteDevice, link.ObservedRemoteInterface,
				link.Message)
		}
		return
	}
}
