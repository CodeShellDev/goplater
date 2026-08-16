package debug

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(echoFunc)

// Outputs data to stdout.
//
// @param data []any
var echoFunc = modules.NewFunc("echo", echo)

func echo(_ *templating.Runtime, _ *templating.Context, data ...any) string {
	fmt.Println(data...)

	return ""
}