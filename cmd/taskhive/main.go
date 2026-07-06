package main

import (
	"fmt"
	"os"

	"github.com/mianm12/taskhive/cmd/taskhive/command"
)

func main() {
	if err := command.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
