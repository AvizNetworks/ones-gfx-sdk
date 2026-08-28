// Create a fabric from the command line.
//
// Usage
// -----
//
//	go run ./examples/createfabric \
//	    -url https://localhost:3002 \
//	    -username superadmin \
//	    -password 'Admin@123456' \
//	    -name "gpu-fabric-1" \
//	    -type "DNO ASN" \
//	    -description "Primary GPU fabric" \
//	    -status draft
//
// Minimal (only the connection details and a fabric name are required):
//
//	go run ./examples/createfabric \
//	    -url https://localhost:3002 -username superadmin -password 'Admin@123456' \
//	    -name "gpu-fabric-1"
//
// List the fabrics afterwards instead of creating one:
//
//	go run ./examples/createfabric \
//	    -url https://localhost:3002 -username superadmin -password 'Admin@123456' \
//	    -name unused -list
//
// All flags:
//
//	-url          ONES base URL, e.g. https://host:3002     (required)
//	-username     login username                            (required)
//	-password     login password                            (required)
//	-name         fabric name                               (required)
//	-type         fabric type, e.g. "DNO ASN"
//	-description  free-text description
//	-status       fabric status, e.g. draft
//	-num-sus      number of SUs                             (int)
//	-max-sus      maximum number of SUs                     (int)
//	-dedicated    mark the fabric dedicated                 (flag)
//	-list         list fabrics instead of creating one
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

type options struct {
	url, username, password string
	name, fabricType        string
	description, status     string
	numSus, maxSus          int
	dedicated, list         bool
}

func parseFlags() (*options, error) {
	o := &options{}
	flag.StringVar(&o.url, "url", "", "ONES base URL (required)")
	flag.StringVar(&o.username, "username", "", "login username (required)")
	flag.StringVar(&o.password, "password", "", "login password (required)")

	flag.StringVar(&o.name, "name", "", "fabric name (required)")
	flag.StringVar(&o.fabricType, "type", "", `fabric type, e.g. "DNO ASN"`)
	flag.StringVar(&o.description, "description", "", "free-text description")
	flag.StringVar(&o.status, "status", "", "fabric status, e.g. draft")
	flag.IntVar(&o.numSus, "num-sus", 0, "number of SUs")
	flag.IntVar(&o.maxSus, "max-sus", 0, "maximum number of SUs")
	flag.BoolVar(&o.dedicated, "dedicated", false, "mark the fabric dedicated")

	flag.BoolVar(&o.list, "list", false, "list fabrics instead of creating one")
	flag.Parse()

	for _, req := range []struct{ val, name string }{
		{o.url, "url"}, {o.username, "username"},
		{o.password, "password"}, {o.name, "name"},
	} {
		if req.val == "" {
			return nil, fmt.Errorf("-%s is required", req.name)
		}
	}
	return o, nil
}

// optStr returns a pointer for a flag the user actually set, nil otherwise, so
// unset fields are omitted from the request body.
func optStr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func optInt(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func optBool(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}

// createFabric creates one fabric from the parsed CLI arguments.
func createFabric(ctx context.Context, client *ones.Client, o *options) (string, error) {
	return ones.CreateFabric(ctx, client, o.name, &ones.FabricCreateArgs{
		Type:        optStr(o.fabricType),
		Description: optStr(o.description),
		Status:      optStr(o.status),
		NumOfSus:    optInt(o.numSus),
		MaxNumOfSus: optInt(o.maxSus),
		Dedicated:   optBool(o.dedicated),
	})
}

func run() error {
	o, err := parseFlags()
	if err != nil {
		flag.Usage()
		return err
	}
	ctx := context.Background()

	client, err := ones.InitializeWithCreds(o.url, o.username, o.password)
	if err != nil {
		return err
	}
	defer client.Close()

	if o.list {
		fabrics, err := ones.GetAllFabrics(ctx, client)
		if err != nil {
			return err
		}
		fmt.Printf("Fabrics (%d):\n", len(fabrics))
		for _, f := range fabrics {
			fmt.Printf("  - id=%d name=%q type=%q status=%q\n",
				f.ID, f.Name, f.Type, f.Status)
		}
		return nil
	}

	fmt.Printf("Creating fabric %q on %s ...\n", o.name, o.url)
	msg, err := createFabric(ctx, client, o)
	if err != nil {
		return err
	}
	fmt.Println("Response:", msg)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
