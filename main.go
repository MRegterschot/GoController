package main

import (
	"fmt"

	"alexejk.io/go-xmlrpc"
)

func main() {
	client, _ := xmlrpc.NewClient("https://bugzilla.mozilla.org/xmlrpc.cgi")
	defer client.Close()

	result := &struct {
		BugzillaVersion struct {
			Version string
		}
	}{}

	_ = client.Call("Bugzilla.version", nil, result)
	fmt.Printf("Version: %s\n", result.BugzillaVersion.Version)
}
