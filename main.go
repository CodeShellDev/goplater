package main

import (
	"github.com/codeshelldev/goplater/cmd"
	"github.com/codeshelldev/goplater/internals/store"
)

var Version string

func main() {
	store.Version = Version

	cmd.Execute()
}
