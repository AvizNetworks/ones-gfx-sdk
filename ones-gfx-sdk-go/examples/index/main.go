// Quick manual test of the ONES Go SDK auth + API flow — the Go twin of
// ones-gfx-sdk-python/index.py.
//
// Run from the SDK root:
//
//	go run ./examples/index
//
// Note this file needs only the one `ones` import: every type, constant and
// option it uses is re-exported there.
package main

import (
	"context"
	"fmt"

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
	ctx := context.Background()

	fmt.Printf("Logging in to %s as %q ...\n", baseURL, username)
	client, err := ones.InitializeWithCreds(baseURL, username, password)
	if err != nil {
		fmt.Println("login failed:", err)
		return
	}
	defer client.Close()
	fmt.Printf("Authenticated: %v (token %d chars)\n",
		client.IsAuthenticated(), len(client.GetAuthToken()))

	fmt.Println("\nFetching all fabrics ...")
	fabrics, err := ones.GetAllFabrics(ctx, client)
	if err != nil {
		fmt.Println("GetAllFabrics failed:", err)
	} else {
		fmt.Printf("Fabrics (%d):\n", len(fabrics))
		for _, f := range fabrics {
			fmt.Printf("  - id=%d name=%q type=%q status=%q description=%q\n",
				f.ID, f.Name, f.Type, f.Status, f.Description)
		}
	}

	fmt.Println("\nAdding a fabric ...")
	msg, err := ones.CreateFabric(ctx, client, "CLI ASN Fabric Go", &ones.FabricCreateArgs{
		Type:        ones.Ptr("DNO ASN"),
		Description: ones.Ptr("fabric created from go cli"),
		Status:      ones.Ptr("draft"),
	})
	if err != nil {
		fmt.Println("CreateFabric failed:", err)
	} else {
		fmt.Println("Add fabric response:", msg)
	}

	// fmt.Println("\nDeleting a fabric ...")
	// delMsg, err := ones.DeleteFabric(ctx, client, "CLI ASN Fabric Go")
	// if err != nil {
	// 	fmt.Println("DeleteFabric failed:", err)
	// } else {
	// 	fmt.Println("Delete fabric response:", delMsg)
	// }

	// fmt.Println("\nRefreshing auth ...")
	// refreshed, err := client.RefreshAuth()
	// if err != nil {
	// 	fmt.Println("refreshAuth failed:", err)
	// } else {
	// 	fmt.Printf("Refresh response: message=%q\n", refreshed.Data.Message)
	// }

	// fmt.Println("\nLogging out ...")
	// out, err := client.Logout()
	// if err != nil {
	// 	fmt.Println("logout failed:", err)
	// } else if out != nil {
	// 	fmt.Printf("Logout response: message=%q\n", out.Data.Message)
	// }
}
