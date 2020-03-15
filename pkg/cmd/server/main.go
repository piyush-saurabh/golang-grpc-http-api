package main

import (
	"fmt"
	"os"

	"github.com/piyush-saurabh/golang-grpc-http-api/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
