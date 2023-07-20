package main

import (
	"github.com/sgaunet/wms/cli/cmd"
)

var version = "No Version Provided"

func main() {
	// t := time.Now()
	cmd.Execute(version)
	// fmt.Println(time.Since(t))
}
