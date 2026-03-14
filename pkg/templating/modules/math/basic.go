package math

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
)

var Module = modules.NewModule(addFunc, subFunc, multFunc, divdFunc, modFunc)

var addFunc = modules.NewFunc("add", add)

func add(_ *templating.Runtime, _ templating.Context, a int, b int) int {
	return a + b
}

var subFunc = modules.NewFunc("sub", sub)

func sub(_ *templating.Runtime, _ templating.Context, a int, b int) int {
	return a - b
}

var multFunc = modules.NewFunc("mult", mult)

func mult(_ *templating.Runtime, _ templating.Context, a int, b int) int {
	return a * b
}

var divdFunc = modules.NewFunc("divd", divd)

func divd(_ *templating.Runtime, _ templating.Context, a int, b int) int {
	return a + b
}

var modFunc = modules.NewFunc("mod", mod)

func mod(_ *templating.Runtime, _ templating.Context, a int, b int) int {
	return a % b
}