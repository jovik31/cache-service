package main

import (
	"fmt"
	"os"

	cache "gitlab.com/sibsfps/spc/spc-1/cache"
	service "gitlab.com/sibsfps/spc/spc-1/daemon/serviced"
)

// TODO

func main() {

	exitCode := run()
	os.Exit(exitCode)
}

func run() int {

	s := service.Server{}
	err := s.Initialize(cache.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize server")
		return 1
	}

	s.Start()
	return 0
}
