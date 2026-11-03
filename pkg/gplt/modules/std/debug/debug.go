package debug

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

var Module = modules.NewModule(echoFunc)

// Outputs data to stdout.
//
// @param data []any
var echoFunc = modules.NewFunc("echo", echo)

func echo(_ *gplt.Runtime, _ *gplt.Context, data ...any) string {
	fmt.Println(data...)

	return ""
}