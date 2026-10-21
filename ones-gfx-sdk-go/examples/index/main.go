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
	fabrics, err := ones.GetFabrics(ctx, client)
	if err != nil {
		fmt.Println("GetFabrics failed:", err)
	} else {
		fmt.Printf("Fabrics (%d):\n", len(fabrics.Fabrics))
		for _, f := range fabrics.Fabrics {
			fmt.Printf("  - id=%d name=%q description=%q\n", f.ID, f.FabricName, f.Description)
		}
	}

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
