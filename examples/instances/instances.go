package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	gcpinstancesinfo "github.com/davidcollom/gcp-instance-info"
)

func main() {

	data, err := gcpinstancesinfo.Data()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// fmt.Println(data.Compute.Instances)

	for _, instance := range data.Compute.Instances {

		spew.Dump(instance)
		fmt.Println("")
	}

}
