package debug

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(echoFunc)

var echoFunc = modules.NewFunc("echo", echo)

func echo(_ *templating.Runtime, _ templating.Context, data ...any) any {
	data = modules.UnpackArgs(data...)

	fmt.Println(data...)

	return nil
}