// Command ns_interfaces prints each host's north-south interfaces as YAML.
//
// From a topology file (as downloaded from GET /fabrics/{name}):
//
//	go run ./examples/ns_interfaces -file Test_yaml.yaml -fabric Test
//
// Or fetched from ONES:
//
//	go run ./examples/ns_interfaces -url https://<ones> -user admin -pass '<pw>' -fabric Test
//
// -leaf-pattern is repeatable and defaults to ^leaf-cn, ^leaf-fn and ^leaf-sn;
// -o writes to a file.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones"
)

type patternList []string

func (p *patternList) String() string     { return strings.Join(*p, ",") }
func (p *patternList) Set(v string) error { *p = append(*p, v); return nil }

func main() {
	var patterns patternList
	file := flag.String("file", "", "topology YAML file")
	fabric := flag.String("fabric", "", "fabric name (required with -url)")
	out := flag.String("o", "", "write YAML here instead of stdout")
	baseURL := flag.String("url", "", "ONES base URL, to fetch the topology instead of reading -file")
	user := flag.String("user", "", "username for -url")
	pass := flag.String("pass", "", "password for -url")
	flag.Var(&patterns, "leaf-pattern", "leaf-name regex; repeatable (default ^leaf-cn ^leaf-fn ^leaf-sn)")
	flag.Parse()

	opts := ones.NSInterfaceOptions{FabricName: *fabric, LeafPatterns: patterns}
	var (
		res *ones.NSInterfaces
		err error
	)
	switch {
	case *baseURL != "":
		if *fabric == "" || *user == "" || *pass == "" {
			fail("-url needs -fabric, -user and -pass")
		}
		client, cerr := ones.InitializeWithCreds(*baseURL, *user, *pass)
		if cerr != nil {
			fail("login: %v", cerr)
		}
		defer client.Close()
		res, err = ones.GetFabricNSInterfaces(context.Background(), client, *fabric, opts)
	case *file != "":
		data, rerr := os.ReadFile(*file)
		if rerr != nil {
			fail("%v", rerr)
		}
		res, err = ones.ExtractNSInterfaces(data, opts)
	default:
		fail("give -file, or -url with -fabric, -user and -pass")
	}
	if err != nil {
		fail("%v", err)
	}

	y, err := res.YAML()
	if err != nil {
		fail("%v", err)
	}
	if *out != "" {
		if err := os.WriteFile(*out, y, 0o644); err != nil {
			fail("%v", err)
		}
		return
	}
	os.Stdout.Write(y)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
