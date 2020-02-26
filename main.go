package main

import (
	"fmt"
	"os"
	cmd "proto-parking-services/pkg/cmd/server"
	"proto-parking-services/pkg/services/api/v1/models"
)

func main() {

	go func() {
		models.HandleSocket()
	}()

	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
