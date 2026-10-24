package main

import (
	"github.com/codeshelldev/goplater/cmd"
	"github.com/codeshelldev/goplater/utils/printer"
)

func main() {
	p := printer.New()
	printer.SetDefault(p)

	cmd.Execute()
}
