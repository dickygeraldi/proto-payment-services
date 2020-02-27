package main

import (
	"fmt"
	"os"
	cmd "proto-parking-services/pkg/cmd/server"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
