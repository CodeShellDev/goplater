package funcs

import "github.com/codeshelldev/goplater/pkg/templating/modules"

var Module = modules.NewModule(readFunc, readArgsFunc, outputToFunc)