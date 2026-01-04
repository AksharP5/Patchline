package main

import (
	"os"

	"github.com/AksharP5/Patchline/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
