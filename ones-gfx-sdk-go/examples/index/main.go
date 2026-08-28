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

	// One setup call: store the root base URL. Init wires the /api/user/ auth
	// endpoints and the /api/fm/ resource base, and loads any persisted token.
	ones.Init(baseURL)
	defer ones.Close()

	fmt.Printf("Token on startup: %q\n", ones.Token())

	fmt.Printf("\nLogging in to %s as %q ...\n", baseURL, username)
	login, err := ones.Login(username, password)
	if err != nil {
		fmt.Println("login failed:", err)
		return
	}
	fmt.Printf("Login response: message=%q isPwdResetNeeded=%v\n",
		login.Data.Message, login.Data.IsPwdResetNeeded)

	fmt.Println("\nFetching all fabrics ...")
	fabrics, err := ones.GetAllFabrics(ctx)
	if err != nil {
		fmt.Println("GetAllFabrics failed:", err)
	} else {
		fmt.Printf("Fabrics (%d):\n", len(fabrics))
		for _, f := range fabrics {
			fmt.Printf("  - id=%d name=%q type=%q status=%q description=%q\n",
				f.ID, f.Name, f.Type, f.Status, f.Description)
		}
	}

	// fmt.Println("\nAdding a fabric ...")
	// msg, err := ones.CreateFabric(ctx, "CLI ASN Fabric Go", &ones.FabricCreateArgs{
	// 	Type:        ones.Ptr("DNO ASN"),
	// 	Description: ones.Ptr("fabric created from go cli"),
	// 	Status:      ones.Ptr("draft"),
	// })
	// if err != nil {
	// 	fmt.Println("CreateFabric failed:", err)
	// } else {
	// 	fmt.Println("Add fabric response:", msg)
	// }

	// fmt.Println("\nDeleting a fabric ...")
	// delMsg, err := ones.DeleteFabric(ctx, "CLI ASN Fabric")
	// if err != nil {
	// 	fmt.Println("DeleteFabric failed:", err)
	// } else {
	// 	fmt.Println("Delete fabric response:", delMsg)
	// }

	// fmt.Println("\nRefreshing auth ...")
	// refreshed, err := ones.Refresh()
	// if err != nil {
	// 	fmt.Println("refresh failed:", err)
	// } else {
	// 	fmt.Printf("Refresh response: message=%q\n", refreshed.Data.Message)
	// }

	// fmt.Println("\nLogging out ...")
	// out, err := ones.Logout()
	// if err != nil {
	// 	fmt.Println("logout failed:", err)
	// } else if out != nil {
	// 	fmt.Printf("Logout response: message=%q\n", out.Data.Message)
	// }
}
